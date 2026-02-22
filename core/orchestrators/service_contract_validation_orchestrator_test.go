package orchestrators

import (
	"errors"
	"testing"

	"github.com/nkasozi/code-for-africa-service-contract-validator/core/dtos"
	"github.com/nkasozi/code-for-africa-service-contract-validator/core/entities"
	"github.com/nkasozi/code-for-africa-service-contract-validator/core/interfaces/ports"
)



func TestValidate_AllRulesPass_ReturnsNoViolations(t *testing.T) {
	orchestrator := ServiceContractValidationOrchestrator{}
	logger := &MockLogger{}

	rules_provider := &MockRulesProvider{
		rules: []ports.IServiceContractRule{
			&MockRule{name: "rule1", is_satisfied: true},
			&MockRule{name: "rule2", is_satisfied: true},
		},
	}

	exceptions_provider := &MockExceptionsProvider{
		exceptions: []ports.IServiceContractRuleException{},
	}

	command := dtos.ValidateServiceContractCommand{
		Mode:                                  entities.ENFORCE,
		ServiceContractRulesProvider:          rules_provider,
		ServiceContractRuleExceptionsProvider: exceptions_provider,
		ServiceContract: dtos.UnvalidatedServiceContract{
			ServiceName:     "test-service",
			Owner:           "test-owner",
			Env:             "prod",
			DataSensitivity: "internal",
		},
		Logger: logger,
	}

	result, err := orchestrator.Validate(command)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
		return
	}

	if len(result.BrokenRules) != 0 {
		t.Errorf("Expected 0 broken rules, got: %d", len(result.BrokenRules))
	}
}

func TestValidate_RuleFails_ReturnsBrokenRule(t *testing.T) {
	orchestrator := ServiceContractValidationOrchestrator{}
	logger := &MockLogger{}

	rules_provider := &MockRulesProvider{
		rules: []ports.IServiceContractRule{
			&MockRule{name: "rule1", is_satisfied: true},
			&MockRule{name: "required_fields", is_satisfied: false},
		},
	}

	exceptions_provider := &MockExceptionsProvider{
		exceptions: []ports.IServiceContractRuleException{},
	}

	command := dtos.ValidateServiceContractCommand{
		Mode:                                  entities.ENFORCE,
		ServiceContractRulesProvider:          rules_provider,
		ServiceContractRuleExceptionsProvider: exceptions_provider,
		ServiceContract: dtos.UnvalidatedServiceContract{
			ServiceName:     "test-service",
			Owner:           "test-owner",
			Env:             "prod",
			DataSensitivity: "internal",
		},
		Logger: logger,
	}

	result, err := orchestrator.Validate(command)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
		return
	}

	if len(result.BrokenRules) != 1 {
		t.Errorf("Expected 1 broken rule, got: %d", len(result.BrokenRules))
		return
	}

	if result.BrokenRules[0].Rule.GetRuleName() != "required_fields" {
		t.Errorf("Expected broken rule 'required_fields', got: %s", result.BrokenRules[0].Rule.GetRuleName())
	}
}

func TestValidate_WithValidException_RemovesBrokenRule(t *testing.T) {
	orchestrator := ServiceContractValidationOrchestrator{}
	logger := &MockLogger{}

	rules_provider := &MockRulesProvider{
		rules: []ports.IServiceContractRule{
			&MockRule{name: "required_fields", is_satisfied: false},
		},
	}

	exceptions_provider := &MockExceptionsProvider{
		exceptions: []ports.IServiceContractRuleException{
			&MockException{
				rule_name:    "required_fields",
				service_name: "test-service",
				is_expired:   false,
				reason:       "Legacy service migration",
				expires:      "2027-01-01",
				approved_by:  "team-lead",
			},
		},
	}

	command := dtos.ValidateServiceContractCommand{
		Mode:                                  entities.ENFORCE,
		ServiceContractRulesProvider:          rules_provider,
		ServiceContractRuleExceptionsProvider: exceptions_provider,
		ServiceContract: dtos.UnvalidatedServiceContract{
			ServiceName:     "test-service",
			Owner:           "test-owner",
			Env:             "prod",
			DataSensitivity: "internal",
		},
		Logger: logger,
	}

	result, err := orchestrator.Validate(command)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
		return
	}

	if len(result.BrokenRules) != 0 {
		t.Errorf("Expected 0 broken rules after exception, got: %d", len(result.BrokenRules))
	}

	if len(result.Exceptions) != 1 {
		t.Errorf("Expected 1 applied exception, got: %d", len(result.Exceptions))
	}
}

func TestValidate_WithExpiredException_DoesNotRemoveBrokenRule(t *testing.T) {
	orchestrator := ServiceContractValidationOrchestrator{}
	logger := &MockLogger{}

	rules_provider := &MockRulesProvider{
		rules: []ports.IServiceContractRule{
			&MockRule{name: "required_fields", is_satisfied: false},
		},
	}

	exceptions_provider := &MockExceptionsProvider{
		exceptions: []ports.IServiceContractRuleException{
			&MockException{
				rule_name:    "required_fields",
				service_name: "test-service",
				is_expired:   true,
			},
		},
	}

	command := dtos.ValidateServiceContractCommand{
		Mode:                                  entities.ENFORCE,
		ServiceContractRulesProvider:          rules_provider,
		ServiceContractRuleExceptionsProvider: exceptions_provider,
		ServiceContract: dtos.UnvalidatedServiceContract{
			ServiceName:     "test-service",
			Owner:           "test-owner",
			Env:             "prod",
			DataSensitivity: "internal",
		},
		Logger: logger,
	}

	result, err := orchestrator.Validate(command)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
		return
	}

	if len(result.BrokenRules) != 1 {
		t.Errorf("Expired exception should not remove broken rule, got: %d broken rules", len(result.BrokenRules))
	}

	if len(result.Exceptions) != 0 {
		t.Errorf("Expired exception should not be applied, got: %d exceptions", len(result.Exceptions))
	}
}

func TestValidate_ExceptionForDifferentService_DoesNotApply(t *testing.T) {
	orchestrator := ServiceContractValidationOrchestrator{}
	logger := &MockLogger{}

	rules_provider := &MockRulesProvider{
		rules: []ports.IServiceContractRule{
			&MockRule{name: "required_fields", is_satisfied: false},
		},
	}

	exceptions_provider := &MockExceptionsProvider{
		exceptions: []ports.IServiceContractRuleException{
			&MockException{
				rule_name:    "required_fields",
				service_name: "different-service",
				is_expired:   false,
			},
		},
	}

	command := dtos.ValidateServiceContractCommand{
		Mode:                                  entities.ENFORCE,
		ServiceContractRulesProvider:          rules_provider,
		ServiceContractRuleExceptionsProvider: exceptions_provider,
		ServiceContract: dtos.UnvalidatedServiceContract{
			ServiceName:     "test-service",
			Owner:           "test-owner",
			Env:             "prod",
			DataSensitivity: "internal",
		},
		Logger: logger,
	}

	result, err := orchestrator.Validate(command)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
		return
	}

	if len(result.BrokenRules) != 1 {
		t.Errorf("Exception for different service should not apply, got: %d broken rules", len(result.BrokenRules))
	}
}

func TestValidate_RulesProviderFails_ReturnsError(t *testing.T) {
	orchestrator := ServiceContractValidationOrchestrator{}
	logger := &MockLogger{}

	rules_provider := &MockRulesProvider{
		should_fail: true,
	}

	exceptions_provider := &MockExceptionsProvider{}

	command := dtos.ValidateServiceContractCommand{
		Mode:                                  entities.ENFORCE,
		ServiceContractRulesProvider:          rules_provider,
		ServiceContractRuleExceptionsProvider: exceptions_provider,
		ServiceContract: dtos.UnvalidatedServiceContract{
			ServiceName:     "test-service",
			Owner:           "test-owner",
			Env:             "prod",
			DataSensitivity: "internal",
		},
		Logger: logger,
	}

	_, err := orchestrator.Validate(command)

	if err == nil {
		t.Error("Expected error when rules provider fails, got nil")
	}
}

func TestValidate_ExceptionsProviderFails_ReturnsError(t *testing.T) {
	orchestrator := ServiceContractValidationOrchestrator{}
	logger := &MockLogger{}

	rules_provider := &MockRulesProvider{
		rules: []ports.IServiceContractRule{},
	}

	exceptions_provider := &MockExceptionsProvider{
		should_fail: true,
	}

	command := dtos.ValidateServiceContractCommand{
		Mode:                                  entities.ENFORCE,
		ServiceContractRulesProvider:          rules_provider,
		ServiceContractRuleExceptionsProvider: exceptions_provider,
		ServiceContract: dtos.UnvalidatedServiceContract{
			ServiceName:     "test-service",
			Owner:           "test-owner",
			Env:             "prod",
			DataSensitivity: "internal",
		},
		Logger: logger,
	}

	_, err := orchestrator.Validate(command)

	if err == nil {
		t.Error("Expected error when exceptions provider fails, got nil")
	}
}

func TestValidate_MultipleBrokenRulesWithPartialExceptions_HandlesCorrectly(t *testing.T) {
	orchestrator := ServiceContractValidationOrchestrator{}
	logger := &MockLogger{}

	rules_provider := &MockRulesProvider{
		rules: []ports.IServiceContractRule{
			&MockRule{name: "rule1", is_satisfied: false},
			&MockRule{name: "rule2", is_satisfied: false},
			&MockRule{name: "rule3", is_satisfied: false},
		},
	}

	exceptions_provider := &MockExceptionsProvider{
		exceptions: []ports.IServiceContractRuleException{
			&MockException{
				rule_name:    "rule1",
				service_name: "test-service",
				is_expired:   false,
			},
			&MockException{
				rule_name:    "rule3",
				service_name: "test-service",
				is_expired:   false,
			},
		},
	}

	command := dtos.ValidateServiceContractCommand{
		Mode:                                  entities.ENFORCE,
		ServiceContractRulesProvider:          rules_provider,
		ServiceContractRuleExceptionsProvider: exceptions_provider,
		ServiceContract: dtos.UnvalidatedServiceContract{
			ServiceName:     "test-service",
			Owner:           "test-owner",
			Env:             "prod",
			DataSensitivity: "internal",
		},
		Logger: logger,
	}

	result, err := orchestrator.Validate(command)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
		return
	}

	if len(result.BrokenRules) != 1 {
		t.Errorf("Expected 1 broken rule (rule2), got: %d", len(result.BrokenRules))
		return
	}

	if result.BrokenRules[0].Rule.GetRuleName() != "rule2" {
		t.Errorf("Expected broken rule 'rule2', got: %s", result.BrokenRules[0].Rule.GetRuleName())
	}

	if len(result.Exceptions) != 2 {
		t.Errorf("Expected 2 applied exceptions, got: %d", len(result.Exceptions))
	}
}

func TestValidate_PreservesValidatorMode(t *testing.T) {
	orchestrator := ServiceContractValidationOrchestrator{}
	logger := &MockLogger{}

	rules_provider := &MockRulesProvider{
		rules: []ports.IServiceContractRule{},
	}

	exceptions_provider := &MockExceptionsProvider{
		exceptions: []ports.IServiceContractRuleException{},
	}

	warn_command := dtos.ValidateServiceContractCommand{
		Mode:                                  entities.WARN,
		ServiceContractRulesProvider:          rules_provider,
		ServiceContractRuleExceptionsProvider: exceptions_provider,
		ServiceContract: dtos.UnvalidatedServiceContract{
			ServiceName:     "test-service",
			Owner:           "test-owner",
			Env:             "prod",
			DataSensitivity: "internal",
		},
		Logger: logger,
	}

	result, _ := orchestrator.Validate(warn_command)

	if result.Mode != entities.WARN {
		t.Errorf("Expected WARN mode in result, got: %v", result.Mode)
	}

	enforce_command := dtos.ValidateServiceContractCommand{
		Mode:                                  entities.ENFORCE,
		ServiceContractRulesProvider:          rules_provider,
		ServiceContractRuleExceptionsProvider: exceptions_provider,
		ServiceContract: dtos.UnvalidatedServiceContract{
			ServiceName:     "test-service",
			Owner:           "test-owner",
			Env:             "prod",
			DataSensitivity: "internal",
		},
		Logger: logger,
	}

	result, _ = orchestrator.Validate(enforce_command)

	if result.Mode != entities.ENFORCE {
		t.Errorf("Expected ENFORCE mode in result, got: %v", result.Mode)
	}
}





//Mocks
type MockRule struct {
	name         string
	is_satisfied bool
}

func (m *MockRule) GetRuleName() string {
	return m.name
}

func (m *MockRule) IsRuleSatisfied(contract ports.IUnvalidatedServiceContract) error {
	if m.is_satisfied {
		return nil
	}
	return errors.New("rule not satisfied")
}

type MockRulesProvider struct {
	rules       []ports.IServiceContractRule
	should_fail bool
}

func (m *MockRulesProvider) LoadServiceContractRules() ([]ports.IServiceContractRule, error) {
	if m.should_fail {
		return nil, errors.New("failed to load rules")
	}
	return m.rules, nil
}

type MockException struct {
	rule_name    string
	service_name string
	is_expired   bool
	reason       string
	expires      string
	approved_by  string
}

func (m *MockException) AppliesToRuleAndService(rule_name string, service_name string) bool {
	return m.rule_name == rule_name && m.service_name == service_name
}

func (m *MockException) IsExpired() bool {
	return m.is_expired
}

func (m *MockException) GetRule() string {
	return m.rule_name
}

func (m *MockException) GetService() string {
	return m.service_name
}

func (m *MockException) GetReason() string {
	return m.reason
}

func (m *MockException) GetExpires() string {
	return m.expires
}

func (m *MockException) GetApprovedBy() string {
	return m.approved_by
}

type MockExceptionsProvider struct {
	exceptions  []ports.IServiceContractRuleException
	should_fail bool
}

func (m *MockExceptionsProvider) LoadServiceRuleExceptions() ([]ports.IServiceContractRuleException, error) {
	if m.should_fail {
		return nil, errors.New("failed to load exceptions")
	}
	return m.exceptions, nil
}

type MockLogger struct {
	info_messages    []string
	error_messages   []string
	warning_messages []string
}

func (m *MockLogger) LogInfo(message string) {
	m.info_messages = append(m.info_messages, message)
}

func (m *MockLogger) LogError(message string) {
	m.error_messages = append(m.error_messages, message)
}

func (m *MockLogger) LogWarning(message string) {
	m.warning_messages = append(m.warning_messages, message)
}
