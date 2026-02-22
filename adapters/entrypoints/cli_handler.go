package entrypoints

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/nkasozi/code-for-africa-service-contract-validator/core/dtos"
	"github.com/nkasozi/code-for-africa-service-contract-validator/core/entities"
	"github.com/nkasozi/code-for-africa-service-contract-validator/core/interfaces/orchestrators"
	"github.com/nkasozi/code-for-africa-service-contract-validator/core/interfaces/ports"
)

const (
	EXIT_CODE_SUCCESS           = 0
	EXIT_CODE_VALIDATION_FAILED = 1
	EXIT_CODE_USAGE_ERROR       = 2

	OUTPUT_PREFIX_FAIL      = "FAIL"
	OUTPUT_PREFIX_EXCEPTION = "EXCEPTION"
	OUTPUT_PREFIX_PASS      = "PASS"
)

type CLIHandler struct {
	validator                orchestrators.IServiceContractValidator
	rules_provider           ports.IServiceContractRulesProvider
	rule_exceptions_provider ports.IServiceContractRuleExceptionsProvider
	logger                   ports.ILoggerProvider
}

type CLIRunResult struct {
	ExitCode int
	Output   string
}

func NewCLIHandler(
	validator orchestrators.IServiceContractValidator,
	rules_provider ports.IServiceContractRulesProvider,
	rule_exceptions_provider ports.IServiceContractRuleExceptionsProvider,
	logger ports.ILoggerProvider,
) *CLIHandler {
	return &CLIHandler{
		validator:                validator,
		rules_provider:           rules_provider,
		rule_exceptions_provider: rule_exceptions_provider,
		logger:                   logger,
	}
}

func (h *CLIHandler) Run(args []string) CLIRunResult {
	h.logger.LogInfo("Starting service contract validation")

	parsed_args, parse_error := h.parseArguments(args)
	if parse_error != nil {
		h.logger.LogError(fmt.Sprintf("Failed to parse arguments: %s", parse_error.Error()))
		return CLIRunResult{
			ExitCode: EXIT_CODE_USAGE_ERROR,
			Output:   fmt.Sprintf("Error: %s\n\nUsage: validate --mode=warn|enforce <service.yaml>", parse_error.Error()),
		}
	}

	h.logger.LogInfo(fmt.Sprintf("Validating file: %s in %s mode", parsed_args.FilePath, parsed_args.ValidatorMode.String()))

	service_contract, load_error := dtos.NewUnvalidatedServiceContractFromFile(parsed_args.FilePath)
	if load_error != nil {
		h.logger.LogError(fmt.Sprintf("Failed to load service contract: %s", load_error.Error()))
		return CLIRunResult{
			ExitCode: EXIT_CODE_USAGE_ERROR,
			Output:   fmt.Sprintf("Error: %s", load_error.Error()),
		}
	}

	h.logger.LogInfo(fmt.Sprintf("Loaded service contract for service: %s", service_contract.GetServiceName()))

	validation_command := dtos.ValidateServiceContractCommand{
		Mode:                                  parsed_args.ValidatorMode,
		ServiceContractRulesProvider:          h.rules_provider,
		ServiceContractRuleExceptionsProvider: h.rule_exceptions_provider,
		ServiceContract:                       service_contract,
		Logger:                                h.logger,
	}

	validation_result, validation_error := h.validator.Validate(validation_command)
	if validation_error != nil {
		h.logger.LogError(fmt.Sprintf("Validation failed with error: %s", validation_error.Error()))
		return CLIRunResult{
			ExitCode: EXIT_CODE_VALIDATION_FAILED,
			Output:   fmt.Sprintf("Error: %s", validation_error.Error()),
		}
	}

	output := h.formatValidationOutput(validation_result, service_contract)
	exit_code := h.determineExitCode(validation_result)

	h.logger.LogInfo(fmt.Sprintf("Validation completed with exit code: %d", exit_code))

	return CLIRunResult{
		ExitCode: exit_code,
		Output:   output,
	}
}

type ParsedArguments struct {
	ValidatorMode entities.ServiceValidatorMode
	FilePath      string
}

func (h *CLIHandler) parseArguments(args []string) (ParsedArguments, error) {
	flag_set := flag.NewFlagSet("validate", flag.ContinueOnError)
	mode_flag := flag_set.String("mode", entities.WARN.String(), "Validation mode: warn or enforce")

	parse_error := flag_set.Parse(args)
	if parse_error != nil {
		return ParsedArguments{}, parse_error
	}

	remaining_args := flag_set.Args()
	if len(remaining_args) == 0 {
		return ParsedArguments{}, fmt.Errorf("no service contract file specified\n\nExamples:\n  validate --mode=warn service.yaml\n  validate --mode=enforce ./path/to/service.yaml")
	}

	file_path := remaining_args[0]
	mode_string := strings.ToLower(*mode_flag)

	validator_mode, is_valid_mode := entities.ParseServiceValidatorMode(mode_string)
	if !is_valid_mode {
		return ParsedArguments{}, fmt.Errorf("invalid mode '%s': must be 'warn' or 'enforce'\n\nExamples:\n  validate --mode=warn service.yaml\n  validate --mode=enforce service.yaml", mode_string)
	}

	return ParsedArguments{
		ValidatorMode: validator_mode,
		FilePath:      file_path,
	}, nil
}

func (h *CLIHandler) formatValidationOutput(
	result dtos.ValidateServiceContractResult,
	contract dtos.UnvalidatedServiceContract,
) string {
	var output_builder strings.Builder
	service_name := contract.GetServiceName()
	environment := contract.GetEnv()

	for _, exception := range result.Exceptions {
		output_builder.WriteString(h.formatExceptionOutput(exception, service_name, environment))
		output_builder.WriteString("\n")
	}

	for index, broken_rule := range result.BrokenRules {
		if index > 0 {
			output_builder.WriteString("\n\n")
		}
		output_builder.WriteString(h.formatFailureOutput(broken_rule.Rule, broken_rule.Error, service_name, environment))
	}

	has_unexcepted_failures := len(result.BrokenRules) > 0
	if !has_unexcepted_failures {
		output_builder.WriteString(fmt.Sprintf("%s: All validation rules passed for service '%s'\n", OUTPUT_PREFIX_PASS, service_name))
	}

	return output_builder.String()
}

func (h *CLIHandler) formatFailureOutput(
	rule ports.IServiceContractRule,
	validation_error error,
	service_name string,
	environment string,
) string {
	failure, is_detailed := validation_error.(*entities.RuleValidationFailure)
	if is_detailed {
		return fmt.Sprintf(
			"%s: %s\n"+
			"Service: %s (%s)\n"+
			"Issue: %s\n"+
			"\n"+
			"Found: %s\n"+
			"Need: %s\n"+
			"\n"+
			"Examples: %s\n"+
			"\n"+
			"Fix: %s\n",
			OUTPUT_PREFIX_FAIL,
			rule.GetRuleName(),
			service_name,
			environment,
			failure.GetIssue(),
			failure.GetFound(),
			failure.GetNeed(),
			failure.GetExamples(),
			failure.GetFixSuggestion(),
		)
	}

	return fmt.Sprintf(
		"%s: %s\n"+
		"Service: %s (%s)\n"+
		"Issue: %s\n",
		OUTPUT_PREFIX_FAIL,
		rule.GetRuleName(),
		service_name,
		environment,
		validation_error.Error(),
	)
}

func (h *CLIHandler) formatExceptionOutput(
	exception ports.IServiceContractRuleException,
	service_name string,
	environment string,
) string {
	return fmt.Sprintf(
		"[%s]: %s\n"+
		"Service: %s (%s)\n"+
		"Reason: %s\n"+
		"Expires: %s\n"+
		"Approved by: %s\n",
		OUTPUT_PREFIX_EXCEPTION,
		exception.GetRule(),
		service_name,
		environment,
		exception.GetReason(),
		exception.GetExpires(),
		exception.GetApprovedBy(),
	)
}

func (h *CLIHandler) determineExitCode(result dtos.ValidateServiceContractResult) int {
	has_unexcepted_failures := len(result.BrokenRules) > 0
	is_enforce_mode := result.Mode == entities.ENFORCE

	if has_unexcepted_failures && is_enforce_mode {
		return EXIT_CODE_VALIDATION_FAILED
	}

	return EXIT_CODE_SUCCESS
}

func ExecuteCLI(
	validator orchestrators.IServiceContractValidator,
	rules_provider ports.IServiceContractRulesProvider,
	rule_exceptions_provider ports.IServiceContractRuleExceptionsProvider,
	logger ports.ILoggerProvider,
) {
	handler := NewCLIHandler(validator, rules_provider, rule_exceptions_provider, logger)
	result := handler.Run(os.Args[1:])
	fmt.Print(result.Output)
	os.Exit(result.ExitCode)
}
