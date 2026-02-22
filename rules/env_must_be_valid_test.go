package rules

import "testing"

func TestEnvMustBeValid_GetRuleName(t *testing.T) {
	rule := &EnvMustBeValid{}
	expected := "env_must_be_valid"
	if result := rule.GetRuleName(); result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestEnvMustBeValid_Dev(t *testing.T) {
	rule := &EnvMustBeValid{}
	contract := createValidServiceContract()
	contract.Env = "dev"
	if result := rule.IsRuleSatisfied(contract); result != nil {
		t.Errorf("Expected no error for dev, got: %v", result)
	}
}

func TestEnvMustBeValid_Staging(t *testing.T) {
	rule := &EnvMustBeValid{}
	contract := createValidServiceContract()
	contract.Env = "staging"
	if result := rule.IsRuleSatisfied(contract); result != nil {
		t.Errorf("Expected no error for staging, got: %v", result)
	}
}

func TestEnvMustBeValid_Prod(t *testing.T) {
	rule := &EnvMustBeValid{}
	contract := createValidServiceContract()
	contract.Env = "prod"
	if result := rule.IsRuleSatisfied(contract); result != nil {
		t.Errorf("Expected no error for prod, got: %v", result)
	}
}

func TestEnvMustBeValid_UppercaseValid(t *testing.T) {
	rule := &EnvMustBeValid{}
	for _, env := range []string{"DEV", "STAGING", "PROD"} {
		contract := createValidServiceContract()
		contract.Env = env
		if result := rule.IsRuleSatisfied(contract); result != nil {
			t.Errorf("Expected no error for %s, got: %v", env, result)
		}
	}
}

func TestEnvMustBeValid_InvalidEnv(t *testing.T) {
	rule := &EnvMustBeValid{}
	contract := createValidServiceContract()
	contract.Env = "production"
	if result := rule.IsRuleSatisfied(contract); result == nil {
		t.Error("Expected error for invalid env 'production'")
	}
}
