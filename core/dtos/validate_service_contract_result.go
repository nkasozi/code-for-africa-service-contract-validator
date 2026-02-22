package dtos

import (
	"github.com/nkasozi/code-for-africa-service-contract-validator/core/entities"
	"github.com/nkasozi/code-for-africa-service-contract-validator/core/interfaces/ports"
)

type RuleValidationError struct {
	Rule  ports.IServiceContractRule
	Error error
}

type ValidateServiceContractResult struct {
	Mode             entities.ServiceValidatorMode
	BrokenRules      []RuleValidationError
	Exceptions       []ports.IServiceContractRuleException
}