package entrypoints

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nkasozi/code-for-africa-service-contract-validator/core/dtos"
	"github.com/nkasozi/code-for-africa-service-contract-validator/core/entities"
	"github.com/nkasozi/code-for-africa-service-contract-validator/core/interfaces/ports"
)



func TestRun_NoFileArgument_ReturnsUsageErrorWithExamples(t *testing.T) {
	handler := createCLITestHandler(t)
	result := handler.Run([]string{})

	if result.ExitCode != EXIT_CODE_USAGE_ERROR {
		t.Errorf("Expected exit code %d, got: %d", EXIT_CODE_USAGE_ERROR, result.ExitCode)
	}
	if !strings.Contains(result.Output, "no service contract file specified") {
		t.Error("Expected error message about missing file")
	}
	if !strings.Contains(result.Output, "Examples:") {
		t.Error("Expected error to include examples")
	}
}

func TestRun_InvalidMode_ReturnsUsageErrorWithExamples(t *testing.T) {
	handler := createCLITestHandler(t)
	temp_file := createCLITempServiceFile(t)
	result := handler.Run([]string{"--mode=invalid", temp_file})

	if result.ExitCode != EXIT_CODE_USAGE_ERROR {
		t.Errorf("Expected exit code %d, got: %d", EXIT_CODE_USAGE_ERROR, result.ExitCode)
	}
	if !strings.Contains(result.Output, "invalid mode 'invalid'") {
		t.Error("Expected error message about invalid mode")
	}
	if !strings.Contains(result.Output, "must be 'warn' or 'enforce'") {
		t.Error("Expected error to show valid mode options")
	}
}

func TestRun_NonExistentFile_ReturnsUsageError(t *testing.T) {
	handler := createCLITestHandler(t)
	result := handler.Run([]string{"--mode=warn", "nonexistent-file.yaml"})

	if result.ExitCode != EXIT_CODE_USAGE_ERROR {
		t.Errorf("Expected exit code %d, got: %d", EXIT_CODE_USAGE_ERROR, result.ExitCode)
	}
}

func TestRun_InvalidYAML_ReturnsUsageError(t *testing.T) {
	handler := createCLITestHandler(t)
	temp_file := createCLITempFileWithContent(t, "invalid: yaml: [[[")
	result := handler.Run([]string{"--mode=warn", temp_file})

	if result.ExitCode != EXIT_CODE_USAGE_ERROR {
		t.Errorf("Expected exit code %d, got: %d", EXIT_CODE_USAGE_ERROR, result.ExitCode)
	}
}

func TestRun_WarnMode_ParsesCorrectly(t *testing.T) {
	logger := &MockCLILogger{}
	validator := &MockValidator{
		result: dtos.ValidateServiceContractResult{
			Mode:        entities.WARN,
			BrokenRules: []dtos.RuleValidationError{},
			Exceptions:  []ports.IServiceContractRuleException{},
		},
	}
	handler := NewCLIHandler(validator, &MockRulesProvider{}, &MockExceptionsProvider{}, logger)
	temp_file := createCLITempServiceFile(t)
	result := handler.Run([]string{"--mode=warn", temp_file})

	if result.ExitCode != EXIT_CODE_SUCCESS {
		t.Errorf("Expected exit code %d, got: %d", EXIT_CODE_SUCCESS, result.ExitCode)
	}
	if !strings.Contains(result.Output, "PASS:") {
		t.Error("Expected PASS output for valid service contract")
	}
}

func TestRun_EnforceMode_ParsesCorrectly(t *testing.T) {
	logger := &MockCLILogger{}
	validator := &MockValidator{
		result: dtos.ValidateServiceContractResult{
			Mode:        entities.ENFORCE,
			BrokenRules: []dtos.RuleValidationError{},
			Exceptions:  []ports.IServiceContractRuleException{},
		},
	}
	handler := NewCLIHandler(validator, &MockRulesProvider{}, &MockExceptionsProvider{}, logger)
	temp_file := createCLITempServiceFile(t)
	result := handler.Run([]string{"--mode=enforce", temp_file})

	if result.ExitCode != EXIT_CODE_SUCCESS {
		t.Errorf("Expected exit code %d, got: %d", EXIT_CODE_SUCCESS, result.ExitCode)
	}
}

func TestRun_WarnModeWithBrokenRules_ReturnsSuccessExitCode(t *testing.T) {
	logger := &MockCLILogger{}
	mock_rule := &MockCLIRule{name: "test_rule"}
	validator := &MockValidator{
		result: dtos.ValidateServiceContractResult{
			Mode: entities.WARN,
			BrokenRules: []dtos.RuleValidationError{
				{Rule: mock_rule, Error: entities.NewRuleValidationFailure("test_rule", "found", "need", "examples", "fix")},
			},
			Exceptions: []ports.IServiceContractRuleException{},
		},
	}
	handler := NewCLIHandler(validator, &MockRulesProvider{}, &MockExceptionsProvider{}, logger)
	temp_file := createCLITempServiceFile(t)
	result := handler.Run([]string{"--mode=warn", temp_file})

	if result.ExitCode != EXIT_CODE_SUCCESS {
		t.Errorf("Expected exit code %d in warn mode even with broken rules, got: %d", EXIT_CODE_SUCCESS, result.ExitCode)
	}
	if !strings.Contains(result.Output, "FAIL:") {
		t.Error("Expected FAIL output for broken rule")
	}
}

func TestRun_EnforceModeWithBrokenRules_ReturnsFailedExitCode(t *testing.T) {
	logger := &MockCLILogger{}
	mock_rule := &MockCLIRule{name: "test_rule"}
	validator := &MockValidator{
		result: dtos.ValidateServiceContractResult{
			Mode: entities.ENFORCE,
			BrokenRules: []dtos.RuleValidationError{
				{Rule: mock_rule, Error: entities.NewRuleValidationFailure("test_rule", "found", "need", "examples", "fix")},
			},
			Exceptions: []ports.IServiceContractRuleException{},
		},
	}
	handler := NewCLIHandler(validator, &MockRulesProvider{}, &MockExceptionsProvider{}, logger)
	temp_file := createCLITempServiceFile(t)
	result := handler.Run([]string{"--mode=enforce", temp_file})

	if result.ExitCode != EXIT_CODE_VALIDATION_FAILED {
		t.Errorf("Expected exit code %d in enforce mode with broken rules, got: %d", EXIT_CODE_VALIDATION_FAILED, result.ExitCode)
	}
}

func TestRun_DefaultsToWarnMode_WhenModeNotSpecified(t *testing.T) {
	logger := &MockCLILogger{}
	validator := &MockValidator{
		result: dtos.ValidateServiceContractResult{
			Mode:        entities.WARN,
			BrokenRules: []dtos.RuleValidationError{},
			Exceptions:  []ports.IServiceContractRuleException{},
		},
	}
	handler := NewCLIHandler(validator, &MockRulesProvider{}, &MockExceptionsProvider{}, logger)
	temp_file := createCLITempServiceFile(t)
	result := handler.Run([]string{temp_file})

	if result.ExitCode != EXIT_CODE_SUCCESS {
		t.Errorf("Expected exit code %d when mode defaults to warn, got: %d", EXIT_CODE_SUCCESS, result.ExitCode)
	}
}

func TestRun_FailureOutput_IncludesActionableDetails(t *testing.T) {
	logger := &MockCLILogger{}
	mock_rule := &MockCLIRule{name: "test_rule"}
	validator := &MockValidator{
		result: dtos.ValidateServiceContractResult{
			Mode: entities.WARN,
			BrokenRules: []dtos.RuleValidationError{
				{Rule: mock_rule, Error: entities.NewRuleValidationFailure("test_rule", "what was found", "what is needed", "example values", "how to fix it")},
			},
			Exceptions: []ports.IServiceContractRuleException{},
		},
	}
	handler := NewCLIHandler(validator, &MockRulesProvider{}, &MockExceptionsProvider{}, logger)
	temp_file := createCLITempServiceFile(t)
	result := handler.Run([]string{"--mode=warn", temp_file})

	if !strings.Contains(result.Output, "Found:") {
		t.Error("Expected output to include 'Found:' section")
	}
	if !strings.Contains(result.Output, "Need:") {
		t.Error("Expected output to include 'Need:' section")
	}
	if !strings.Contains(result.Output, "Examples:") {
		t.Error("Expected output to include 'Examples:' section")
	}
	if !strings.Contains(result.Output, "Fix:") {
		t.Error("Expected output to include 'Fix:' section")
	}
}



type MockValidator struct {
	result dtos.ValidateServiceContractResult
	err    error
}

func (m *MockValidator) Validate(command dtos.ValidateServiceContractCommand) (dtos.ValidateServiceContractResult, error) {
	return m.result, m.err
}

type MockRulesProvider struct {
	rules []ports.IServiceContractRule
	err   error
}

func (m *MockRulesProvider) LoadServiceContractRules() ([]ports.IServiceContractRule, error) {
	return m.rules, m.err
}

type MockExceptionsProvider struct {
	exceptions []ports.IServiceContractRuleException
	err        error
}

func (m *MockExceptionsProvider) LoadServiceRuleExceptions() ([]ports.IServiceContractRuleException, error) {
	return m.exceptions, m.err
}

type MockCLILogger struct {
	info_logs    []string
	error_logs   []string
	warning_logs []string
}

func (m *MockCLILogger) LogInfo(message string) {
	m.info_logs = append(m.info_logs, message)
}

func (m *MockCLILogger) LogError(message string) {
	m.error_logs = append(m.error_logs, message)
}

func (m *MockCLILogger) LogWarning(message string) {
	m.warning_logs = append(m.warning_logs, message)
}

type MockCLIRule struct {
	name string
}

func (m *MockCLIRule) GetRuleName() string {
	return m.name
}

func (m *MockCLIRule) IsRuleSatisfied(contract ports.IUnvalidatedServiceContract) error {
	return nil
}

func createCLITestHandler(_ *testing.T) *CLIHandler {
	logger := &MockCLILogger{}
	validator := &MockValidator{
		result: dtos.ValidateServiceContractResult{
			BrokenRules: []dtos.RuleValidationError{},
			Exceptions:  []ports.IServiceContractRuleException{},
		},
	}
	return NewCLIHandler(validator, &MockRulesProvider{}, &MockExceptionsProvider{}, logger)
}

func createCLITempServiceFile(t *testing.T) string {
	content := `schema_version: "1"
service_name: test-service
owner: team:platform
env: prod
data_sensitivity: low
cost_center: engineering
alerts:
  - high_error_rate
`
	return createCLITempFileWithContent(t, content)
}

func createCLITempFileWithContent(t *testing.T, content string) string {
	temp_dir := t.TempDir()
	temp_file := filepath.Join(temp_dir, "service.yaml")
	err := os.WriteFile(temp_file, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	return temp_file
}
