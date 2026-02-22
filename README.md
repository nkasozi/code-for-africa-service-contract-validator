# Service Contract Validator

A small CLI tool that tries to validate `service.yaml` files against a set of platform rules. The idea is to catch common issues early before deployments get messy. It's fairly straightforward - no guarantee it covers every edge case, but it handles the basics reasonably well.

## Prerequisites

This project requires **Go 1.25.0** or later.

### macOS

```bash
brew install go@1.25
```

Or download from https://go.dev/dl/

### Linux

```bash
wget https://go.dev/dl/go1.25.0.linux-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.25.0.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
```

### Windows

Download the installer from https://go.dev/dl/go1.25.0.windows-amd64.msi

### Verify Installation

```bash
go version
# should output: go version go1.25.0 ...
```

## Getting Started

```bash
go build -o validate .

./validate --mode=warn service.yaml

./validate --mode=enforce service.yaml
```

## What It Checks

### The Baseline Stuff

| Rule                                     | What It Does                                                                                                                |
| ---------------------------------------- | --------------------------------------------------------------------------------------------------------------------------- |
| `all_required_fields_must_be_present`    | Makes sure the essential fields are there (schema_version, service_name, owner, env, data_sensitivity, cost_center, alerts) |
| `env_must_be_valid`                      | Environment should be one of: `dev`, `staging`, `prod`                                                                      |
| `data_sensitivity_must_be_valid`         | Data sensitivity should be one of: `low`, `medium`, `high`                                                                  |
| `prod_env_must_have_symptom_based_alert` | Production services need at least one useful alert (5xx errors, latency stuff, health checks)                               |

### The Extra Policy

| Rule                                        | What It Does                                                            |
| ------------------------------------------- | ----------------------------------------------------------------------- |
| `high_sensitivity_requires_data_governance` | High sensitivity services need `retention_days` and `data_owner` fields |

I went with this extra rule because when you're handling sensitive data in a civic-tech context, it seems sensible to know who owns the data and how long it sticks around. Might be overkill for some situations, but better to have clear accountability than not.

## Warn vs Enforce

Two modes, nothing fancy:

- **warn**: Tells you what's wrong but doesn't block anything. Good for easing teams into it.
- **enforce**: Same feedback but exits with code 1 when things fail. Use this once people have had time to fix their stuff.

```bash
./validate --mode=warn service.yaml
./validate --mode=enforce service.yaml
```

## Exceptions

Sometimes rules need to be bent temporarily. Put exceptions in an `exceptions.yaml`:

```yaml
- rule: prod_env_must_have_symptom_based_alert
  service: legacy-api
  reason: Migrating to new monitoring stack
  expires: "2026-06-01"
  approved_by: platform-team
```

When there's a matching exception, the validator notes it but doesn't fail. Expired exceptions are ignored.

## Schema Changes

If you ever need to bump to `schema_version: 2`:

1. Add new fields as optional first
2. Run in warn mode for any new v2 rules
3. Give people time to migrate
4. Eventually, after enough adoption, enforce the new rules
5. Use exceptions for stragglers

Nothing revolutionary here - just gradual, boring rollout.

## Rolling Out to Multiple Repos

Here's roughly how I'd approach it:

### Week 1 - Just Watch

- Turn it on in warn mode everywhere
- Don't block anything
- See what fails, get a baseline

### Weeks 2-3 - Fix the Easy Stuff

- Let teams know what needs fixing
- Set up exceptions for the tricky legacy cases
- Most services probably just need a few fields added
- Inform Teams of the deadline when we will switch to enforcement
- End of every week, talk to teams that still have failures

### Week 4+ - Start Enforcing

- Again, Inform Teams that we are going to start enforcing the rules
- Switch repos with clean runs to enforce mode first
- Keep an eye on whether people are struggling
- Switch repos with some failing runs to enforce mode
- Inform all Teams that switch has happened such that they now know deployment succeded
- Adjust if needed after discussions with Teams

### Week 5 o 6 - Postmortem

What went right or wrong. How can we do it better next time.

### What to Watch For

- Are most repos eventually passing? Good sign
- Are the Same rules failing everywhere? Maybe the error messages need work
- Too many exceptions? Rules might be too strict for where we are
- Lots of questions? Maybe Docs might need improving, can we automate some of the setup stuff

## Example Output

When everything passes:

```
PASS: All validation rules passed for service 'api-gateway'
```

When something fails:

```
FAIL: all_required_fields_must_be_present
  Service: incomplete-service (prod)
  Issue: Required fields are missing or empty

  Found: cost_center: [empty], alerts: [empty]
  Need: All required fields must be present and non-empty

  Examples: schema_version: '1', service_name: my-service, owner: team:platform

  Fix: Add the missing fields to your service.yaml: cost_center, alerts
```

With an exception:

```
[EXCEPTION]: prod_env_must_have_symptom_based_alert
  Service: legacy-api (prod)
  Reason: Migrating to new monitoring stack
  Expires: 2026-06-01
  Approved by: platform-team
PASS: All validation rules passed for service 'legacy-api'
```

## Running Tests

```bash
go test ./... -v
```

## Adding a New Rule

To add a new validation rule, create a struct that implements the `IServiceContractRule` interface:

```go
type IServiceContractRule interface {
    IsRuleSatisfied(serviceContract IUnvalidatedServiceContract) error
    GetRuleName() string
}
```

### Step 1: Create the Rule File

Create a new file in `rules/` (e.g., `rules/my_new_rule.go`):

```go
package rules

import (
    "github.com/nkasozi/code-for-africa-service-contract-validator/core/entities"
    "github.com/nkasozi/code-for-africa-service-contract-validator/core/interfaces/ports"
)

const RULE_NAME_MY_NEW_RULE = "my_new_rule"

type MyNewRule struct{}

func (r *MyNewRule) GetRuleName() string {
    return RULE_NAME_MY_NEW_RULE
}

func (r *MyNewRule) IsRuleSatisfied(contract ports.IUnvalidatedServiceContract) error {
    if !someCondition(contract) {
        return entities.NewRuleValidationFailure(
            RULE_NAME_MY_NEW_RULE,
            "brief issue description",
            "what we found: actual value",
            "what we need: expected value",
            "example: valid_value",
            "how to fix: do this thing",
        )
    }
    return nil
}
```

### Step 2: Register the Rule

Add your rule to the list in `global_constants.go`:

```go
var AllServiceContractRulesToApply = []ports.IServiceContractRule{
    &rules.AllRequiredFieldsMustBePresent{},
    &rules.EnvMustBeValid{},
    &rules.MyNewRule{},  // add here
}
```

### Step 3: Write Tests First

Create `rules/my_new_rule_test.go` with test cases before implementing the rule logic.

## Adding Rule Exceptions

### Via exceptions.yaml (Recommended)

Create or edit `exceptions.yaml` in your project root:

```yaml
- rule: my_new_rule
  service: legacy-service
  reason: Migration in progress
  expires: "2026-12-01"
  approved_by: platform-team
```

Fields:

- `rule`: Must match the rule's `GetRuleName()` return value (case-insensitive)
- `service`: Must match the service name in `service.yaml` (case-insensitive)
- `reason`: Why this exception exists
- `expires`: Date in `YYYY-MM-DD` format (expired exceptions are ignored)
- `approved_by`: Who approved this exception

### Via Code (Hardcoded)

For permanent exceptions, add to `global_constants.go`:

```go
var AllServiceContractRuleExceptionsToApply = []ports.IServiceContractRuleException{
    &entities.ServiceContractRuleException{
        Rule:       "my_new_rule",
        Service:    "special-service",
        Reason:     "Permanent exception for legacy system",
        Expires:    "2099-12-31",
        ApprovedBy: "platform-team",
    },
}
```

## Architecture

This project follows a hexagonal (ports & adapters) architecture:

```
├── main.go                     # Entry point
├── app.go                      # Composition root (wires dependencies)
├── global_constants.go         # App-wide constants and rule registration
│
├── core/                       # THE CORE (Business Logic)
│   ├── entities/               # Domain objects (ServiceContract, Alert, etc.)
│   ├── interfaces/
│   │   ├── ports/              # Contracts for external dependencies
│   │   └── orchestrators/      # Contracts for business use-cases
│   ├── orchestrators/          # Business logic implementation
│   ├── dtos/                   # Command/Result objects
│   └── shared/                 # Domain helpers
│
├── adapters/                   # THE ADAPTERS (Infrastructure)
│   ├── entrypoints/            # CLI handler
│   ├── providers/              # Rule and exception loaders
│   └── persistence/            # (future: file storage)
│
└── rules/                      # Validation rule implementations
```

### Key Principles

**Dependency Inversion**: Core business logic depends only on interfaces (ports), never on concrete implementations. The `app.go` composition root wires everything together.

**Testability**: Every component receives its dependencies via constructor injection. No global state manipulation in business logic.

**Separation of Concerns**:

- `core/orchestrators/` - Coordinates validation flow
- `rules/` - Individual rule implementations
- `adapters/` - Infrastructure concerns (CLI, file loading)

### Data Flow

```
CLI Args → CLIHandler → ValidationOrchestrator → Rules → Result → CLI Output
                              ↓
                    RulesProvider (loads rules)
                    ExceptionsProvider (loads exceptions)
```

## AI Disclosure

Am not sure if its full blown AI but I used the generate unit test functionality for help in generating unit test cases, which I reviewed and adjusted before including them.
