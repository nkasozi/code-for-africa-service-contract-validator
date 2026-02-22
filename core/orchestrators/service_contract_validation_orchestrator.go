package orchestrators

import (
	"github.com/nkasozi/code-for-africa-service-contract-validator/core/dtos"
	"github.com/nkasozi/code-for-africa-service-contract-validator/core/interfaces/ports"
)

type ServiceContractValidationOrchestrator struct{}

func (o *ServiceContractValidationOrchestrator) Validate(command dtos.ValidateServiceContractCommand) (dtos.ValidateServiceContractResult, error) {
	command.Logger.LogInfo("Starting validation orchestration")

	rules, rules_load_error := command.ServiceContractRulesProvider.LoadServiceContractRules()
	if rules_load_error != nil {
		command.Logger.LogError("Failed to load rules: " + rules_load_error.Error())
		return dtos.ValidateServiceContractResult{}, rules_load_error
	}

	command.Logger.LogInfo("Loaded rules successfully, checking for violations")
	broken_rules := o.findBrokenRules(rules, command.ServiceContract, command.Logger)

	exceptions, exceptions_load_error := command.ServiceContractRuleExceptionsProvider.LoadServiceRuleExceptions()
	if exceptions_load_error != nil {
		command.Logger.LogError("Failed to load exceptions: " + exceptions_load_error.Error())
		return dtos.ValidateServiceContractResult{}, exceptions_load_error
	}

	command.Logger.LogInfo("Loaded exceptions successfully, applying exceptions to broken rules")
	service_name := command.ServiceContract.GetServiceName()
	applied_exceptions, remaining_broken_rules := o.applyExceptionsToBrokenRules(
		broken_rules,
		exceptions,
		service_name,
		command.Logger,
	)

	return dtos.ValidateServiceContractResult{
		Mode:        command.Mode,
		BrokenRules: remaining_broken_rules,
		Exceptions:  applied_exceptions,
	}, nil
}

func (o *ServiceContractValidationOrchestrator) findBrokenRules(
	rules []ports.IServiceContractRule,
	contract ports.IUnvalidatedServiceContract,
	logger ports.ILoggerProvider,
) []dtos.RuleValidationError {
	broken_rules := []dtos.RuleValidationError{}

	for _, rule := range rules {
		rule_error := rule.IsRuleSatisfied(contract)
		if rule_error != nil {
			logger.LogInfo("Rule violation found: " + rule.GetRuleName())
			broken_rules = append(broken_rules, dtos.RuleValidationError{
				Rule:  rule,
				Error: rule_error,
			})
		}
	}

	return broken_rules
}

func (o *ServiceContractValidationOrchestrator) applyExceptionsToBrokenRules(
	broken_rules []dtos.RuleValidationError,
	exceptions []ports.IServiceContractRuleException,
	service_name string,
	logger ports.ILoggerProvider,
) ([]ports.IServiceContractRuleException, []dtos.RuleValidationError) {
	applied_exceptions := []ports.IServiceContractRuleException{}
	remaining_broken_rules := []dtos.RuleValidationError{}

	for _, broken_rule := range broken_rules {
		matching_exception := o.findMatchingException(broken_rule, exceptions, service_name, logger)

		if matching_exception != nil {
			logger.LogInfo("Exception applied for rule: " + broken_rule.Rule.GetRuleName())
			applied_exceptions = append(applied_exceptions, matching_exception)
			continue
		}

		remaining_broken_rules = append(remaining_broken_rules, broken_rule)
	}

	return applied_exceptions, remaining_broken_rules
}

func (o *ServiceContractValidationOrchestrator) findMatchingException(
	broken_rule dtos.RuleValidationError,
	exceptions []ports.IServiceContractRuleException,
	service_name string,
	logger ports.ILoggerProvider,
) ports.IServiceContractRuleException {
	rule_name := broken_rule.Rule.GetRuleName()

	for _, exception := range exceptions {
		if !exception.AppliesToRuleAndService(rule_name, service_name) {
			continue
		}

		if exception.IsExpired() {
			logger.LogWarning("Exception for rule '" + rule_name + "' has expired")
			continue
		}

		return exception
	}

	return nil
}