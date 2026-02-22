package providers

import (
	"errors"
	"testing"

	"github.com/nkasozi/code-for-africa-service-contract-validator/core/interfaces/ports"
)

func TestLoadServiceContractRules_ReturnsInjectedRules(t *testing.T) {
	rules := []ports.IServiceContractRule{
		&MockTestRule{name: "rule1"},
		&MockTestRule{name: "rule2"},
		&MockTestRule{name: "rule3"},
	}

	provider := NewServiceContractRuleProvider(rules)

	loaded_rules, err := provider.LoadServiceContractRules()

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
		return
	}

	if len(loaded_rules) != 3 {
		t.Errorf("Expected 3 rules, got: %d", len(loaded_rules))
		return
	}

	expected_names := []string{"rule1", "rule2", "rule3"}
	for i, rule := range loaded_rules {
		if rule.GetRuleName() != expected_names[i] {
			t.Errorf("Expected rule %d to be '%s', got: '%s'", i, expected_names[i], rule.GetRuleName())
		}
	}
}

func TestLoadServiceContractRules_WithEmptyList_ReturnsEmpty(t *testing.T) {
	provider := NewServiceContractRuleProvider([]ports.IServiceContractRule{})

	loaded_rules, err := provider.LoadServiceContractRules()

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
		return
	}

	if len(loaded_rules) != 0 {
		t.Errorf("Expected 0 rules, got: %d", len(loaded_rules))
	}
}

func TestLoadServiceContractRules_PreservesRuleOrder(t *testing.T) {
	rules := []ports.IServiceContractRule{
		&MockTestRule{name: "alpha"},
		&MockTestRule{name: "beta"},
		&MockTestRule{name: "gamma"},
		&MockTestRule{name: "delta"},
	}

	provider := NewServiceContractRuleProvider(rules)

	loaded_rules, _ := provider.LoadServiceContractRules()

	for i, rule := range loaded_rules {
		if rule.GetRuleName() != rules[i].GetRuleName() {
			t.Errorf("Rule order not preserved at index %d: expected '%s', got '%s'",
				i, rules[i].GetRuleName(), rule.GetRuleName())
		}
	}
}

func TestLoadServiceContractRules_MultipleCallsReturnsSameRules(t *testing.T) {
	rules := []ports.IServiceContractRule{
		&MockTestRule{name: "consistent_rule"},
	}

	provider := NewServiceContractRuleProvider(rules)

	first_load, _ := provider.LoadServiceContractRules()
	second_load, _ := provider.LoadServiceContractRules()

	if len(first_load) != len(second_load) {
		t.Errorf("Inconsistent rule counts: first=%d, second=%d", len(first_load), len(second_load))
		return
	}

	if first_load[0].GetRuleName() != second_load[0].GetRuleName() {
		t.Errorf("Rule names differ between calls: first='%s', second='%s'",
			first_load[0].GetRuleName(), second_load[0].GetRuleName())
	}
}

type FailingTestRule struct {
	name string
}

func (f *FailingTestRule) GetRuleName() string {
	return f.name
}

func (f *FailingTestRule) IsRuleSatisfied(contract ports.IUnvalidatedServiceContract) error {
	return errors.New("rule always fails")
}

func TestLoadServiceContractRules_SupportsRulesWithDifferentBehaviors(t *testing.T) {
	rules := []ports.IServiceContractRule{
		&MockTestRule{name: "passing_rule"},
		&FailingTestRule{name: "failing_rule"},
	}

	provider := NewServiceContractRuleProvider(rules)

	loaded_rules, err := provider.LoadServiceContractRules()

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
		return
	}

	if len(loaded_rules) != 2 {
		t.Errorf("Expected 2 rules, got: %d", len(loaded_rules))
		return
	}

	passing_result := loaded_rules[0].IsRuleSatisfied(nil)
	if passing_result != nil {
		t.Errorf("First rule should pass, got error: %v", passing_result)
	}

	failing_result := loaded_rules[1].IsRuleSatisfied(nil)
	if failing_result == nil {
		t.Error("Second rule should fail, got nil")
	}
}

// Mocks
type MockTestRule struct {
	name string
}

func (m *MockTestRule) GetRuleName() string {
	return m.name
}

func (m *MockTestRule) IsRuleSatisfied(contract ports.IUnvalidatedServiceContract) error {
	return nil
}
