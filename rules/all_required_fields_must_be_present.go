package rules

import (
	"strings"

	"github.com/nkasozi/code-for-africa-service-contract-validator/core/entities"
	"github.com/nkasozi/code-for-africa-service-contract-validator/core/interfaces/ports"
)

const (
	FIELD_SCHEMA_VERSION      = "schema_version"
	FIELD_SERVICE_NAME        = "service_name"
	FIELD_OWNER               = "owner"
	FIELD_ENV                 = "env"
	FIELD_DATA_SENSITIVITY    = "data_sensitivity"
	FIELD_COST_CENTER         = "cost_center"
	FIELD_ALERTS              = "alerts"
	RULE_NAME_REQUIRED_FIELDS = "all_required_fields_must_be_present"
)

type AllRequiredFieldsMustBePresent struct {
}

func (r *AllRequiredFieldsMustBePresent) GetRuleName() string {
	return RULE_NAME_REQUIRED_FIELDS
}

func (r *AllRequiredFieldsMustBePresent) IsRuleSatisfied(service_contract ports.IUnvalidatedServiceContract) error {
	var missing_fields []string

	append_if_missing := func(field_name string, value string) {
		if strings.TrimSpace(value) == "" {
			missing_fields = append(missing_fields, field_name)
		}
	}

	append_if_missing(FIELD_SCHEMA_VERSION, service_contract.GetSchemaVersion())
	append_if_missing(FIELD_SERVICE_NAME, service_contract.GetServiceName())
	append_if_missing(FIELD_OWNER, service_contract.GetOwner())
	append_if_missing(FIELD_ENV, service_contract.GetEnv())
	append_if_missing(FIELD_DATA_SENSITIVITY, service_contract.GetDataSensitivity())
	append_if_missing(FIELD_COST_CENTER, service_contract.GetCostCenter())

	if len(service_contract.GetAlerts()) == 0 {
		missing_fields = append(missing_fields, FIELD_ALERTS)
	}

	if len(missing_fields) > 0 {
		found_description := buildFoundDescription(missing_fields)
		return entities.NewRuleValidationFailure(
			RULE_NAME_REQUIRED_FIELDS,
			found_description,
			"All required fields must be present and non-empty",
			"schema_version: '1', service_name: my-service, owner: team:platform",
			"Add the missing fields to your service.yaml: "+strings.Join(missing_fields, ", "),
		)
	}

	return nil
}

func buildFoundDescription(missing_fields []string) string {
	descriptions := make([]string, len(missing_fields))
	for i, field := range missing_fields {
		descriptions[i] = field + ": [empty]"
	}
	return strings.Join(descriptions, ", ")
}