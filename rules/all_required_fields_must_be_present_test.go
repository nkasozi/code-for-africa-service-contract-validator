package rules

import "testing"


func TestAllRequiredFieldsMustBePresent_GetRuleName(t *testing.T) {
	rule := &AllRequiredFieldsMustBePresent{}
	expected := "all_required_fields_must_be_present"
	if result := rule.GetRuleName(); result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestAllRequiredFieldsMustBePresent_AllFieldsPresent(t *testing.T) {
	rule := &AllRequiredFieldsMustBePresent{}
	contract := createValidServiceContract()
	if result := rule.IsRuleSatisfied(contract); result != nil {
		t.Errorf("Expected no error, got: %v", result)
	}
}

func TestAllRequiredFieldsMustBePresent_MissingServiceName(t *testing.T) {
	rule := &AllRequiredFieldsMustBePresent{}
	contract := createValidServiceContract()
	contract.ServiceName = ""
	if result := rule.IsRuleSatisfied(contract); result == nil {
		t.Error("Expected error when service_name missing")
	}
}

func TestAllRequiredFieldsMustBePresent_MissingAlerts(t *testing.T) {
	rule := &AllRequiredFieldsMustBePresent{}
	contract := createValidServiceContract()
	contract.Alerts = []string{}
	if result := rule.IsRuleSatisfied(contract); result == nil {
		t.Error("Expected error when alerts missing")
	}
}


//Mocks

type MockServiceContract struct {
	SchemaVersion   string
	ServiceName     string
	Owner           string
	Env             string
	DataSensitivity string
	CostCenter      string
	Alerts          []string
	OtherFields     map[string]interface{}
}

func (m MockServiceContract) GetSchemaVersion() string   { return m.SchemaVersion }
func (m MockServiceContract) GetServiceName() string     { return m.ServiceName }
func (m MockServiceContract) GetOwner() string           { return m.Owner }
func (m MockServiceContract) GetEnv() string             { return m.Env }
func (m MockServiceContract) GetDataSensitivity() string { return m.DataSensitivity }
func (m MockServiceContract) GetCostCenter() string      { return m.CostCenter }
func (m MockServiceContract) GetAlerts() []string        { return m.Alerts }
func (m MockServiceContract) GetOtherFields() map[string]interface{} {
	if m.OtherFields == nil {
		return make(map[string]interface{})
	}
	return m.OtherFields
}

func createValidServiceContract() MockServiceContract {
	return MockServiceContract{
		SchemaVersion:   "1",
		ServiceName:     "test-service",
		Owner:           "team-platform",
		Env:             "prod",
		DataSensitivity: "low",
		CostCenter:      "engineering",
		Alerts:          []string{"5xx_rate_high", "p99_latency_breach"},
		OtherFields:     make(map[string]interface{}),
	}
}
