package main

import (
	"github.com/nkasozi/code-for-africa-service-contract-validator/adapters/entrypoints"
	"github.com/nkasozi/code-for-africa-service-contract-validator/adapters/providers"
	"github.com/nkasozi/code-for-africa-service-contract-validator/core/orchestrators"
)

type Application struct {
	debug_enabled bool
}

func NewApplication(debug_enabled bool) *Application {
	return &Application{
		debug_enabled: debug_enabled,
	}
}

func (a *Application) Run() {
	logger := providers.NewConsoleLogger(a.debug_enabled)

	rules_provider := providers.NewServiceContractRuleProvider(AllServiceContractRulesToApply)
	exceptions_provider := providers.NewServiceContractRuleExceptionProvider(
		DEFAULT_EXCEPTIONS_FILE_PATH,
		AllServiceContractRuleExceptionsToApply,
	)

	validator := &orchestrators.ServiceContractValidationOrchestrator{}

	entrypoints.ExecuteCLI(validator, rules_provider, exceptions_provider, logger)
}
