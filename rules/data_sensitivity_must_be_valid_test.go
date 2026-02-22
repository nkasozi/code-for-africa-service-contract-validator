package rules

import "testing"

func TestDataSensitivityMustBeValid_GetRuleName(t *testing.T) {
	rule := &DataSensitivityMustBeValid{}
	expected := "data_sensitivity_must_be_valid"
	if result := rule.GetRuleName(); result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestDataSensitivityMustBeValid_Low(t *testing.T) {
	rule := &DataSensitivityMustBeValid{}
	contract := createValidServiceContract()
	contract.DataSensitivity = "low"
	if result := rule.IsRuleSatisfied(contract); result != nil {
		t.Errorf("Expected no error for low, got: %v", result)
	}
}

func TestDataSensitivityMustBeValid_Medium(t *testing.T) {
	rule := &DataSensitivityMustBeValid{}
	contract := createValidServiceContract()
	contract.DataSensitivity = "medium"
	if result := rule.IsRuleSatisfied(contract); result != nil {
		t.Errorf("Expected no error for medium, got: %v", result)
	}
}

func TestDataSensitivityMustBeValid_High(t *testing.T) {
	rule := &DataSensitivityMustBeValid{}
	contract := createValidServiceContract()
	contract.DataSensitivity = "high"
	if result := rule.IsRuleSatisfied(contract); result != nil {
		t.Errorf("Expected no error for high, got: %v", result)
	}
}

func TestDataSensitivityMustBeValid_UppercaseValid(t *testing.T) {
	rule := &DataSensitivityMustBeValid{}
	for _, sens := range []string{"LOW", "MEDIUM", "HIGH"} {
		contract := createValidServiceContract()
		contract.DataSensitivity = sens
		if result := rule.IsRuleSatisfied(contract); result != nil {
			t.Errorf("Expected no error for %s, got: %v", sens, result)
		}
	}
}

func TestDataSensitivityMustBeValid_Invalid(t *testing.T) {
	rule := &DataSensitivityMustBeValid{}
	contract := createValidServiceContract()
	contract.DataSensitivity = "critical"
	if result := rule.IsRuleSatisfied(contract); result == nil {
		t.Error("Expected error for invalid sensitivity")
	}
}
