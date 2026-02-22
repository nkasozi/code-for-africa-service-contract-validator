package rules

import (
	"strings"

	"github.com/nkasozi/code-for-africa-service-contract-validator/core/entities"
	"github.com/nkasozi/code-for-africa-service-contract-validator/core/interfaces/ports"
)

const (
	RULE_NAME_HIGH_SENSITIVITY_REQUIRES_GOVERNANCE = "high_sensitivity_requires_data_governance"

	FIELD_NAME_RETENTION_DAYS = "retention_days"
	FIELD_NAME_DATA_OWNER     = "data_owner"
)

type HighSensitivityRequiresDataGovernance struct{}

func (r *HighSensitivityRequiresDataGovernance) GetRuleName() string {
	return RULE_NAME_HIGH_SENSITIVITY_REQUIRES_GOVERNANCE
}

func (r *HighSensitivityRequiresDataGovernance) IsRuleSatisfied(service_contract ports.IUnvalidatedServiceContract) error {
	raw_sensitivity := service_contract.GetDataSensitivity()

	if strings.ToLower(raw_sensitivity) != entities.SensitivityHigh.String() {
		return nil
	}

	other_fields := service_contract.GetOtherFields()
	missing_fields := r.findMissingGovernanceFields(other_fields)

	if len(missing_fields) > 0 {
		missing_fields_string := strings.Join(missing_fields, ", ")
		return entities.NewRuleValidationFailure(
			RULE_NAME_HIGH_SENSITIVITY_REQUIRES_GOVERNANCE,
			"Missing data governance fields for high sensitivity service",
			"data_sensitivity: high, missing: "+missing_fields_string,
			"High sensitivity services require retention_days and data_owner fields",
			"retention_days: 90, data_owner: compliance-team",
			"Add the missing governance fields: "+missing_fields_string,
		)
	}

	return nil
}

func (r *HighSensitivityRequiresDataGovernance) findMissingGovernanceFields(other_fields map[string]interface{}) []string {
	missing_fields := []string{}

	if !r.hasValidRetentionDays(other_fields) {
		missing_fields = append(missing_fields, FIELD_NAME_RETENTION_DAYS)
	}

	if !r.hasValidDataOwner(other_fields) {
		missing_fields = append(missing_fields, FIELD_NAME_DATA_OWNER)
	}

	return missing_fields
}

func (r *HighSensitivityRequiresDataGovernance) hasValidRetentionDays(other_fields map[string]interface{}) bool {
	value, exists := other_fields[FIELD_NAME_RETENTION_DAYS]
	if !exists {
		return false
	}

	switch typed_value := value.(type) {
	case int:
		return typed_value > 0
	case float64:
		return typed_value > 0
	default:
		return false
	}
}

func (r *HighSensitivityRequiresDataGovernance) hasValidDataOwner(other_fields map[string]interface{}) bool {
	value, exists := other_fields[FIELD_NAME_DATA_OWNER]
	if !exists {
		return false
	}

	string_value, is_string := value.(string)
	if !is_string {
		return false
	}

	return strings.TrimSpace(string_value) != ""
}
