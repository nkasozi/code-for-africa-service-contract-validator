package orchestrators

import "github.com/nkasozi/code-for-africa-service-contract-validator/core/dtos"

type IServiceContractValidator interface {
	Validate(command dtos.ValidateServiceContractCommand) (dtos.ValidateServiceContractResult, error)
}