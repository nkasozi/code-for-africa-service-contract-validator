package rules

import (
	"strings"

	"github.com/nkasozi/code-for-africa-service-contract-validator/core/entities"
	"github.com/nkasozi/code-for-africa-service-contract-validator/core/interfaces/ports"
)

const (
	RULE_NAME_VALID_ENV = "env_must_be_valid"
)

type EnvMustBeValid struct {
}

func (r *EnvMustBeValid) GetRuleName() string {
	return RULE_NAME_VALID_ENV
}

func (r *EnvMustBeValid) IsRuleSatisfied(service_contract ports.IUnvalidatedServiceContract) error {
	raw_environment := service_contract.GetEnv()

	_, parse_error := entities.ParseServiceEnvironment(raw_environment)

	if parse_error != nil {
		allowed_options := entities.GetAllowedEnvironmentStrings()
		return entities.NewRuleValidationFailure(
			RULE_NAME_VALID_ENV,
			"Invalid environment value",
			"env: "+raw_environment,
			"env must be one of: "+strings.Join(allowed_options, ", "),
			"env: dev, env: staging, env: prod",
			"Update env field to use a valid environment value",
		)
	}

	return nil
}
