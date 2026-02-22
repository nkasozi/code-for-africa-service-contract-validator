package main

import (
	"github.com/nkasozi/code-for-africa-service-contract-validator/core/interfaces/ports"
	"github.com/nkasozi/code-for-africa-service-contract-validator/rules"
)

const (
	DEFAULT_EXCEPTIONS_FILE_PATH = "exceptions.yaml"
)

var AllServiceContractRulesToApply = []ports.IServiceContractRule{
	&rules.AllRequiredFieldsMustBePresent{},
	&rules.EnvMustBeValid{},
	&rules.DataSensitivityMustBeValid{},
	rules.NewProdEnvMustHaveSymptomBasedAlert(),
	&rules.HighSensitivityRequiresDataGovernance{},
}

var AllServiceContractRuleExceptionsToApply = []ports.IServiceContractRuleException{}
