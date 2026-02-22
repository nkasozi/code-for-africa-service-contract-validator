package dtos

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type UnvalidatedServiceContract struct {
	SchemaVersion   string                 `yaml:"schema_version"`
	ServiceName     string                 `yaml:"service_name"`
	Owner           string                 `yaml:"owner"`
	Env             string                 `yaml:"env"`
	DataSensitivity string                 `yaml:"data_sensitivity"`
	CostCenter      string                 `yaml:"cost_center"`
	Alerts          []string               `yaml:"alerts"`
	OtherFields     map[string]interface{} `yaml:",inline"`
	FilePath        string                 `yaml:"-"`
}

func NewUnvalidatedServiceContractFromFile(file_path string) (UnvalidatedServiceContract, error) {
	file_contents, read_error := os.ReadFile(file_path)
	if read_error != nil {
		return UnvalidatedServiceContract{}, fmt.Errorf("failed to read service contract file '%s': %w", file_path, read_error)
	}

	var contract UnvalidatedServiceContract
	parse_error := yaml.Unmarshal(file_contents, &contract)
	if parse_error != nil {
		return UnvalidatedServiceContract{}, fmt.Errorf("failed to parse service contract YAML in '%s': %w", file_path, parse_error)
	}

	contract.FilePath = file_path
	return contract, nil
}

func NewUnvalidatedServiceContractFromYAML(yaml_content []byte) (UnvalidatedServiceContract, error) {
	var contract UnvalidatedServiceContract
	parse_error := yaml.Unmarshal(yaml_content, &contract)
	if parse_error != nil {
		return UnvalidatedServiceContract{}, fmt.Errorf("failed to parse service contract YAML: %w", parse_error)
	}

	return contract, nil
}

func (u UnvalidatedServiceContract) GetSchemaVersion() string {
	return u.SchemaVersion
}

func (u UnvalidatedServiceContract) GetServiceName() string {
	return u.ServiceName
}

func (u UnvalidatedServiceContract) GetOwner() string {
	return u.Owner
}

func (u UnvalidatedServiceContract) GetEnv() string {
	return u.Env
}

func (u UnvalidatedServiceContract) GetDataSensitivity() string {
	return u.DataSensitivity
}

func (u UnvalidatedServiceContract) GetCostCenter() string {
	return u.CostCenter
}

func (u UnvalidatedServiceContract) GetAlerts() []string {
	return u.Alerts
}

func (u UnvalidatedServiceContract) GetOtherFields() map[string]interface{} {
	if u.OtherFields == nil {
		return make(map[string]interface{})
	}
	return u.OtherFields
}
