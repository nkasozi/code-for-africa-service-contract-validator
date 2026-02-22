package rules

import (
	"strings"

	"github.com/nkasozi/code-for-africa-service-contract-validator/core/entities"
	"github.com/nkasozi/code-for-africa-service-contract-validator/core/interfaces/ports"
)

const (
	RULE_NAME_VALID_SENSITIVITY = "data_sensitivity_must_be_valid"
)

type DataSensitivityMustBeValid struct {
}

func (r *DataSensitivityMustBeValid) GetRuleName() string {
	return RULE_NAME_VALID_SENSITIVITY
}

func (r *DataSensitivityMustBeValid) IsRuleSatisfied(service_contract ports.IUnvalidatedServiceContract) error {
	raw_sensitivity := service_contract.GetDataSensitivity()

	_, parse_error := entities.ParseDataSensitivity(raw_sensitivity)

	if parse_error != nil {
		allowed_options := entities.GetAllowedSensitivityStrings()
		return entities.NewRuleValidationFailure(
			RULE_NAME_VALID_SENSITIVITY,
			"Invalid data sensitivity value",
			"data_sensitivity: "+raw_sensitivity,
			"data_sensitivity must be one of: "+strings.Join(allowed_options, ", "),
			"data_sensitivity: low, data_sensitivity: medium, data_sensitivity: high",
			"Update data_sensitivity field to use a valid value",
		)
	}

	return nil
}
