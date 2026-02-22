package rules

import "testing"

func TestHighSensitivityRequiresDataGovernance_GetRuleName(t *testing.T) {
	rule := &HighSensitivityRequiresDataGovernance{}
	expected := "high_sensitivity_requires_data_governance"
	if result := rule.GetRuleName(); result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestHighSensitivityRequiresDataGovernance_HighWithAllFields(t *testing.T) {
	rule := &HighSensitivityRequiresDataGovernance{}
	contract := createValidServiceContract()
	contract.DataSensitivity = "high"
	contract.OtherFields = map[string]interface{}{
		"retention_days": 90,
		"data_owner":     "data-team",
	}
	if result := rule.IsRuleSatisfied(contract); result != nil {
		t.Errorf("Expected no error with governance fields, got: %v", result)
	}
}

func TestHighSensitivityRequiresDataGovernance_HighMissingRetentionDays(t *testing.T) {
	rule := &HighSensitivityRequiresDataGovernance{}
	contract := createValidServiceContract()
	contract.DataSensitivity = "high"
	contract.OtherFields = map[string]interface{}{
		"data_owner": "data-team",
	}
	if result := rule.IsRuleSatisfied(contract); result == nil {
		t.Error("Expected error when retention_days missing")
	}
}

func TestHighSensitivityRequiresDataGovernance_HighMissingDataOwner(t *testing.T) {
	rule := &HighSensitivityRequiresDataGovernance{}
	contract := createValidServiceContract()
	contract.DataSensitivity = "high"
	contract.OtherFields = map[string]interface{}{
		"retention_days": 90,
	}
	if result := rule.IsRuleSatisfied(contract); result == nil {
		t.Error("Expected error when data_owner missing")
	}
}

func TestHighSensitivityRequiresDataGovernance_LowSensitivity(t *testing.T) {
	rule := &HighSensitivityRequiresDataGovernance{}
	contract := createValidServiceContract()
	contract.DataSensitivity = "low"
	contract.OtherFields = map[string]interface{}{}
	if result := rule.IsRuleSatisfied(contract); result != nil {
		t.Errorf("Expected no error for low sensitivity, got: %v", result)
	}
}

func TestHighSensitivityRequiresDataGovernance_HighWithZeroRetention(t *testing.T) {
	rule := &HighSensitivityRequiresDataGovernance{}
	contract := createValidServiceContract()
	contract.DataSensitivity = "high"
	contract.OtherFields = map[string]interface{}{
		"retention_days": 0,
		"data_owner":     "data-team",
	}
	if result := rule.IsRuleSatisfied(contract); result == nil {
		t.Error("Expected error for zero retention_days")
	}
}

func TestHighSensitivityRequiresDataGovernance_HighWithEmptyDataOwner(t *testing.T) {
	rule := &HighSensitivityRequiresDataGovernance{}
	contract := createValidServiceContract()
	contract.DataSensitivity = "high"
	contract.OtherFields = map[string]interface{}{
		"retention_days": 90,
		"data_owner":     "",
	}
	if result := rule.IsRuleSatisfied(contract); result == nil {
		t.Error("Expected error for empty data_owner")
	}
}
