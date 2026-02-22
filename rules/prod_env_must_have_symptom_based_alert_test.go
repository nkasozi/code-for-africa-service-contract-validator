package rules

import "testing"

func TestProdEnvMustHaveSymptomBasedAlert_GetRuleName(t *testing.T) {
	rule := NewProdEnvMustHaveSymptomBasedAlert()
	expected := "prod_env_must_have_symptom_based_alert"
	if result := rule.GetRuleName(); result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestProdEnvMustHaveSymptomBasedAlert_ProdWith5xx(t *testing.T) {
	rule := NewProdEnvMustHaveSymptomBasedAlert()
	contract := createValidServiceContract()
	contract.Env = "prod"
	contract.Alerts = []string{"5xx_rate_high"}
	if result := rule.IsRuleSatisfied(contract); result != nil {
		t.Errorf("Expected no error with 5xx alert, got: %v", result)
	}
}

func TestProdEnvMustHaveSymptomBasedAlert_ProdWithLatency(t *testing.T) {
	rule := NewProdEnvMustHaveSymptomBasedAlert()
	contract := createValidServiceContract()
	contract.Env = "prod"
	contract.Alerts = []string{"p99_latency_breach"}
	if result := rule.IsRuleSatisfied(contract); result != nil {
		t.Errorf("Expected no error with latency alert, got: %v", result)
	}
}

func TestProdEnvMustHaveSymptomBasedAlert_ProdWithHealthCheck(t *testing.T) {
	rule := NewProdEnvMustHaveSymptomBasedAlert()
	contract := createValidServiceContract()
	contract.Env = "prod"
	contract.Alerts = []string{"health-check-failure"}
	if result := rule.IsRuleSatisfied(contract); result != nil {
		t.Errorf("Expected no error with health-check alert, got: %v", result)
	}
}

func TestProdEnvMustHaveSymptomBasedAlert_ProdWithOnlyInfraAlerts(t *testing.T) {
	rule := NewProdEnvMustHaveSymptomBasedAlert()
	contract := createValidServiceContract()
	contract.Env = "prod"
	contract.Alerts = []string{"cpu_high", "disk_full"}
	if result := rule.IsRuleSatisfied(contract); result == nil {
		t.Error("Expected error for prod with only infra alerts")
	}
}

func TestProdEnvMustHaveSymptomBasedAlert_DevWithNoSymptom(t *testing.T) {
	rule := NewProdEnvMustHaveSymptomBasedAlert()
	contract := createValidServiceContract()
	contract.Env = "dev"
	contract.Alerts = []string{"cpu_high"}
	if result := rule.IsRuleSatisfied(contract); result != nil {
		t.Errorf("Expected no error for dev env, got: %v", result)
	}
}

func TestProdEnvMustHaveSymptomBasedAlert_ProdWithMixedAlerts(t *testing.T) {
	rule := NewProdEnvMustHaveSymptomBasedAlert()
	contract := createValidServiceContract()
	contract.Env = "prod"
	contract.Alerts = []string{"cpu_high", "5xx_rate_high", "disk_full"}
	if result := rule.IsRuleSatisfied(contract); result != nil {
		t.Errorf("Expected no error with mixed alerts, got: %v", result)
	}
}
