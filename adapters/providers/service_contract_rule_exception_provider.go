package providers

import (
	"fmt"
	"os"

	"github.com/nkasozi/code-for-africa-service-contract-validator/core/entities"
	"github.com/nkasozi/code-for-africa-service-contract-validator/core/interfaces/ports"
	"gopkg.in/yaml.v3"
)

type ServiceContractRuleExceptionProvider struct {
	exceptions_file_path string
	hardcoded_exceptions []ports.IServiceContractRuleException
}

func NewServiceContractRuleExceptionProvider(
	exceptions_file_path string,
	hardcoded_exceptions []ports.IServiceContractRuleException,
) *ServiceContractRuleExceptionProvider {
	return &ServiceContractRuleExceptionProvider{
		exceptions_file_path: exceptions_file_path,
		hardcoded_exceptions: hardcoded_exceptions,
	}
}

func (p *ServiceContractRuleExceptionProvider) LoadServiceRuleExceptions() ([]ports.IServiceContractRuleException, error) {
	file_exceptions, load_error := p.loadExceptionsFromFile()
	if load_error != nil {
		return p.hardcoded_exceptions, nil
	}

	all_exceptions := append(p.hardcoded_exceptions, file_exceptions...)
	return all_exceptions, nil
}

func (p *ServiceContractRuleExceptionProvider) loadExceptionsFromFile() ([]ports.IServiceContractRuleException, error) {
	file_contents, read_error := os.ReadFile(p.exceptions_file_path)
	if read_error != nil {
		return nil, fmt.Errorf("failed to read exceptions file '%s': %w", p.exceptions_file_path, read_error)
	}

	var raw_exceptions []entities.ServiceContractRuleException
	parse_error := yaml.Unmarshal(file_contents, &raw_exceptions)
	if parse_error != nil {
		return nil, fmt.Errorf("failed to parse exceptions YAML: %w", parse_error)
	}

	exceptions := make([]ports.IServiceContractRuleException, len(raw_exceptions))
	for i := range raw_exceptions {
		exceptions[i] = &raw_exceptions[i]
	}

	return exceptions, nil
}