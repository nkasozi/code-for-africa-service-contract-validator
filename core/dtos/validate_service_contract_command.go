package dtos

import (
	"github.com/nkasozi/code-for-africa-service-contract-validator/core/entities"
	"github.com/nkasozi/code-for-africa-service-contract-validator/core/interfaces/ports"
)

type ValidateServiceContractCommand struct{
	Mode entities.ServiceValidatorMode
	ServiceContractRulesProvider ports.IServiceContractRulesProvider
	ServiceContractRuleExceptionsProvider ports.IServiceContractRuleExceptionsProvider
	ServiceContract UnvalidatedServiceContract
	Logger ports.ILoggerProvider
}