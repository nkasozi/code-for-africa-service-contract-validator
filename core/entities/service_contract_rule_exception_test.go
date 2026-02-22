package entities

import (
	"testing"
	"time"
)

func TestAppliesToRuleAndService_ExactMatch_ReturnsTrue(t *testing.T) {
	exception := &ServiceContractRuleException{
		Rule:    "AllRequiredFieldsMustBePresent",
		Service: "payment-service",
	}

	result := exception.AppliesToRuleAndService("AllRequiredFieldsMustBePresent", "payment-service")

	if !result {
		t.Error("Expected exception to apply when rule and service match exactly")
	}
}

func TestAppliesToRuleAndService_CaseInsensitive_ReturnsTrue(t *testing.T) {
	exception := &ServiceContractRuleException{
		Rule:    "AllRequiredFieldsMustBePresent",
		Service: "payment-service",
	}

	test_cases := []struct {
		rule_name    string
		service_name string
	}{
		{"allrequiredfieldsmustbepresent", "payment-service"},
		{"ALLREQUIREDFIELDSMUSTBEPRESENT", "payment-service"},
		{"AllRequiredFieldsMustBePresent", "PAYMENT-SERVICE"},
		{"allrequiredfieldsmustbepresent", "PAYMENT-SERVICE"},
	}

	for _, tc := range test_cases {
		result := exception.AppliesToRuleAndService(tc.rule_name, tc.service_name)
		if !result {
			t.Errorf("Expected case-insensitive match for rule='%s', service='%s'", tc.rule_name, tc.service_name)
		}
	}
}

func TestAppliesToRuleAndService_DifferentRule_ReturnsFalse(t *testing.T) {
	exception := &ServiceContractRuleException{
		Rule:    "AllRequiredFieldsMustBePresent",
		Service: "payment-service",
	}

	result := exception.AppliesToRuleAndService("DifferentRule", "payment-service")

	if result {
		t.Error("Expected exception NOT to apply when rule name differs")
	}
}

func TestAppliesToRuleAndService_DifferentService_ReturnsFalse(t *testing.T) {
	exception := &ServiceContractRuleException{
		Rule:    "AllRequiredFieldsMustBePresent",
		Service: "payment-service",
	}

	result := exception.AppliesToRuleAndService("AllRequiredFieldsMustBePresent", "different-service")

	if result {
		t.Error("Expected exception NOT to apply when service name differs")
	}
}

func TestAppliesToRuleAndService_BothDifferent_ReturnsFalse(t *testing.T) {
	exception := &ServiceContractRuleException{
		Rule:    "AllRequiredFieldsMustBePresent",
		Service: "payment-service",
	}

	result := exception.AppliesToRuleAndService("DifferentRule", "different-service")

	if result {
		t.Error("Expected exception NOT to apply when both rule and service differ")
	}
}

func TestIsExpired_EmptyExpires_ReturnsTrue(t *testing.T) {
	exception := &ServiceContractRuleException{
		Rule:    "test_rule",
		Service: "test-service",
		Expires: "",
	}

	result := exception.IsExpired()

	if !result {
		t.Error("Expected empty expires to be treated as expired")
	}
}

func TestIsExpired_WhitespaceOnlyExpires_ReturnsTrue(t *testing.T) {
	exception := &ServiceContractRuleException{
		Rule:    "test_rule",
		Service: "test-service",
		Expires: "   ",
	}

	result := exception.IsExpired()

	if !result {
		t.Error("Expected whitespace-only expires to be treated as expired")
	}
}

func TestIsExpired_InvalidDateFormat_ReturnsTrue(t *testing.T) {
	invalid_dates := []string{
		"not-a-date",
		"2030/01/01",
		"01-01-2030",
		"2030-1-1",
		"January 1, 2030",
	}

	for _, invalid_date := range invalid_dates {
		exception := &ServiceContractRuleException{
			Rule:    "test_rule",
			Service: "test-service",
			Expires: invalid_date,
		}

		result := exception.IsExpired()

		if !result {
			t.Errorf("Expected invalid date format '%s' to be treated as expired", invalid_date)
		}
	}
}

func TestIsExpired_FutureDate_ReturnsFalse(t *testing.T) {
	future_date := time.Now().AddDate(1, 0, 0).Format("2006-01-02")

	exception := &ServiceContractRuleException{
		Rule:    "test_rule",
		Service: "test-service",
		Expires: future_date,
	}

	result := exception.IsExpired()

	if result {
		t.Errorf("Expected future date '%s' to NOT be expired", future_date)
	}
}

func TestIsExpired_PastDate_ReturnsTrue(t *testing.T) {
	past_date := time.Now().AddDate(-1, 0, 0).Format("2006-01-02")

	exception := &ServiceContractRuleException{
		Rule:    "test_rule",
		Service: "test-service",
		Expires: past_date,
	}

	result := exception.IsExpired()

	if !result {
		t.Errorf("Expected past date '%s' to be expired", past_date)
	}
}

func TestIsExpired_Yesterday_ReturnsTrue(t *testing.T) {
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

	exception := &ServiceContractRuleException{
		Rule:    "test_rule",
		Service: "test-service",
		Expires: yesterday,
	}

	result := exception.IsExpired()

	if !result {
		t.Error("Expected yesterday's date to be expired")
	}
}

func TestIsExpired_Tomorrow_ReturnsFalse(t *testing.T) {
	tomorrow := time.Now().AddDate(0, 0, 1).Format("2006-01-02")

	exception := &ServiceContractRuleException{
		Rule:    "test_rule",
		Service: "test-service",
		Expires: tomorrow,
	}

	result := exception.IsExpired()

	if result {
		t.Error("Expected tomorrow's date to NOT be expired")
	}
}

func TestGetters_ReturnCorrectValues(t *testing.T) {
	exception := &ServiceContractRuleException{
		Rule:       "TestRule",
		Service:    "test-service",
		Reason:     "Test reason for exception",
		Expires:    "2030-06-15",
		ApprovedBy: "manager@company.com",
	}

	if exception.GetRule() != "TestRule" {
		t.Errorf("Expected GetRule() to return 'TestRule', got: %s", exception.GetRule())
	}

	if exception.GetService() != "test-service" {
		t.Errorf("Expected GetService() to return 'test-service', got: %s", exception.GetService())
	}

	if exception.GetReason() != "Test reason for exception" {
		t.Errorf("Expected GetReason() to return 'Test reason for exception', got: %s", exception.GetReason())
	}

	if exception.GetExpires() != "2030-06-15" {
		t.Errorf("Expected GetExpires() to return '2030-06-15', got: %s", exception.GetExpires())
	}

	if exception.GetApprovedBy() != "manager@company.com" {
		t.Errorf("Expected GetApprovedBy() to return 'manager@company.com', got: %s", exception.GetApprovedBy())
	}
}

func TestGetters_WithEmptyFields_ReturnEmptyStrings(t *testing.T) {
	exception := &ServiceContractRuleException{}

	if exception.GetRule() != "" {
		t.Errorf("Expected GetRule() to return empty string, got: %s", exception.GetRule())
	}

	if exception.GetService() != "" {
		t.Errorf("Expected GetService() to return empty string, got: %s", exception.GetService())
	}

	if exception.GetReason() != "" {
		t.Errorf("Expected GetReason() to return empty string, got: %s", exception.GetReason())
	}

	if exception.GetExpires() != "" {
		t.Errorf("Expected GetExpires() to return empty string, got: %s", exception.GetExpires())
	}

	if exception.GetApprovedBy() != "" {
		t.Errorf("Expected GetApprovedBy() to return empty string, got: %s", exception.GetApprovedBy())
	}
}
