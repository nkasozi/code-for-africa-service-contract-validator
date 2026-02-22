package providers

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nkasozi/code-for-africa-service-contract-validator/core/entities"
	"github.com/nkasozi/code-for-africa-service-contract-validator/core/interfaces/ports"
)

func TestLoadServiceRuleExceptions_WithHardcodedOnly_ReturnsHardcodedExceptions(t *testing.T) {
	hardcoded_exceptions := []ports.IServiceContractRuleException{
		&entities.ServiceContractRuleException{
			Rule:       "test_rule",
			Service:    "test-service",
			Reason:     "hardcoded reason",
			Expires:    "2030-01-01",
			ApprovedBy: "admin",
		},
	}

	provider := NewServiceContractRuleExceptionProvider(
		"nonexistent_file.yaml",
		hardcoded_exceptions,
	)

	exceptions, err := provider.LoadServiceRuleExceptions()

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
		return
	}

	if len(exceptions) != 1 {
		t.Errorf("Expected 1 exception, got: %d", len(exceptions))
		return
	}

	if exceptions[0].GetRule() != "test_rule" {
		t.Errorf("Expected rule 'test_rule', got: %s", exceptions[0].GetRule())
	}
}

func TestLoadServiceRuleExceptions_WithValidFile_ReturnsFileExceptions(t *testing.T) {
	temp_dir := t.TempDir()
	temp_file := filepath.Join(temp_dir, "exceptions.yaml")

	file_content := `- rule: file_rule
  service: file-service
  reason: from file
  expires: "2030-01-01"
  approved_by: file_admin
`
	err := os.WriteFile(temp_file, []byte(file_content), 0644)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	provider := NewServiceContractRuleExceptionProvider(
		temp_file,
		[]ports.IServiceContractRuleException{},
	)

	exceptions, err := provider.LoadServiceRuleExceptions()

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
		return
	}

	if len(exceptions) != 1 {
		t.Errorf("Expected 1 exception, got: %d", len(exceptions))
		return
	}

	if exceptions[0].GetRule() != "file_rule" {
		t.Errorf("Expected rule 'file_rule', got: %s", exceptions[0].GetRule())
	}

	if exceptions[0].GetService() != "file-service" {
		t.Errorf("Expected service 'file-service', got: %s", exceptions[0].GetService())
	}
}

func TestLoadServiceRuleExceptions_CombinesHardcodedAndFileExceptions(t *testing.T) {
	temp_dir := t.TempDir()
	temp_file := filepath.Join(temp_dir, "exceptions.yaml")

	file_content := `- rule: file_rule
  service: file-service
  reason: from file
  expires: "2030-01-01"
  approved_by: file_admin
`
	err := os.WriteFile(temp_file, []byte(file_content), 0644)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	hardcoded_exceptions := []ports.IServiceContractRuleException{
		&entities.ServiceContractRuleException{
			Rule:       "hardcoded_rule",
			Service:    "hardcoded-service",
			Reason:     "hardcoded reason",
			Expires:    "2030-01-01",
			ApprovedBy: "admin",
		},
	}

	provider := NewServiceContractRuleExceptionProvider(
		temp_file,
		hardcoded_exceptions,
	)

	exceptions, err := provider.LoadServiceRuleExceptions()

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
		return
	}

	if len(exceptions) != 2 {
		t.Errorf("Expected 2 exceptions (1 hardcoded + 1 file), got: %d", len(exceptions))
		return
	}

	found_hardcoded := false
	found_file := false
	for _, exc := range exceptions {
		if exc.GetRule() == "hardcoded_rule" {
			found_hardcoded = true
		}
		if exc.GetRule() == "file_rule" {
			found_file = true
		}
	}

	if !found_hardcoded {
		t.Error("Expected to find hardcoded exception")
	}
	if !found_file {
		t.Error("Expected to find file exception")
	}
}

func TestLoadServiceRuleExceptions_WithInvalidYAML_ReturnsHardcodedOnly(t *testing.T) {
	temp_dir := t.TempDir()
	temp_file := filepath.Join(temp_dir, "exceptions.yaml")

	invalid_content := `this is not valid yaml: [[[`
	err := os.WriteFile(temp_file, []byte(invalid_content), 0644)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	hardcoded_exceptions := []ports.IServiceContractRuleException{
		&entities.ServiceContractRuleException{
			Rule:       "hardcoded_rule",
			Service:    "hardcoded-service",
			Reason:     "hardcoded reason",
			Expires:    "2030-01-01",
			ApprovedBy: "admin",
		},
	}

	provider := NewServiceContractRuleExceptionProvider(
		temp_file,
		hardcoded_exceptions,
	)

	exceptions, err := provider.LoadServiceRuleExceptions()

	if err != nil {
		t.Errorf("Expected no error (should fall back to hardcoded), got: %v", err)
		return
	}

	if len(exceptions) != 1 {
		t.Errorf("Expected 1 hardcoded exception, got: %d", len(exceptions))
	}
}

func TestLoadServiceRuleExceptions_WithMultipleFileExceptions_LoadsAll(t *testing.T) {
	temp_dir := t.TempDir()
	temp_file := filepath.Join(temp_dir, "exceptions.yaml")

	file_content := `- rule: rule1
  service: service1
  reason: reason1
  expires: "2030-01-01"
  approved_by: admin1
- rule: rule2
  service: service2
  reason: reason2
  expires: "2030-06-15"
  approved_by: admin2
- rule: rule3
  service: service3
  reason: reason3
  expires: "2030-12-31"
  approved_by: admin3
`
	err := os.WriteFile(temp_file, []byte(file_content), 0644)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	provider := NewServiceContractRuleExceptionProvider(
		temp_file,
		[]ports.IServiceContractRuleException{},
	)

	exceptions, err := provider.LoadServiceRuleExceptions()

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
		return
	}

	if len(exceptions) != 3 {
		t.Errorf("Expected 3 exceptions, got: %d", len(exceptions))
	}
}

func TestNewServiceContractRuleExceptionProvider_StoresInjectedFilePath(t *testing.T) {
	expected_file_path := "custom_exceptions.yaml"
	provider := NewServiceContractRuleExceptionProvider(expected_file_path, []ports.IServiceContractRuleException{})

	if provider.exceptions_file_path != expected_file_path {
		t.Errorf("Expected file path '%s', got: %s", expected_file_path, provider.exceptions_file_path)
	}
}

func TestLoadServiceRuleExceptions_PreservesExceptionFields(t *testing.T) {
	temp_dir := t.TempDir()
	temp_file := filepath.Join(temp_dir, "exceptions.yaml")

	file_content := `- rule: AllRequiredFieldsMustBePresent
  service: legacy-payment-service
  reason: Migration in progress from legacy system
  expires: "2027-06-01"
  approved_by: jane.smith@company.com
`
	err := os.WriteFile(temp_file, []byte(file_content), 0644)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	provider := NewServiceContractRuleExceptionProvider(
		temp_file,
		[]ports.IServiceContractRuleException{},
	)

	exceptions, err := provider.LoadServiceRuleExceptions()

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
		return
	}

	if len(exceptions) != 1 {
		t.Errorf("Expected 1 exception, got: %d", len(exceptions))
		return
	}

	exception := exceptions[0]

	if exception.GetRule() != "AllRequiredFieldsMustBePresent" {
		t.Errorf("Expected rule 'AllRequiredFieldsMustBePresent', got: %s", exception.GetRule())
	}

	if exception.GetService() != "legacy-payment-service" {
		t.Errorf("Expected service 'legacy-payment-service', got: %s", exception.GetService())
	}

	if exception.GetReason() != "Migration in progress from legacy system" {
		t.Errorf("Expected reason 'Migration in progress from legacy system', got: %s", exception.GetReason())
	}

	if exception.GetExpires() != "2027-06-01" {
		t.Errorf("Expected expires '2027-06-01', got: %s", exception.GetExpires())
	}

	if exception.GetApprovedBy() != "jane.smith@company.com" {
		t.Errorf("Expected approved_by 'jane.smith@company.com', got: %s", exception.GetApprovedBy())
	}
}
