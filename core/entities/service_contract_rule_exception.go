package entities

import (
	"strings"
	"time"
)

type ServiceContractRuleException struct {
	Rule       string `yaml:"rule"`
	Service    string `yaml:"service"`
	Reason     string `yaml:"reason"`
	Expires    string `yaml:"expires"`
	ApprovedBy string `yaml:"approved_by"`
}

func (e *ServiceContractRuleException) AppliesToRuleAndService(rule_name string, service_name string) bool {
	rule_matches := strings.EqualFold(e.Rule, rule_name)
	service_matches := strings.EqualFold(e.Service, service_name)
	return rule_matches && service_matches
}

func (e *ServiceContractRuleException) IsExpired() bool {
	if strings.TrimSpace(e.Expires) == "" {
		return true
	}

	expiry_date, parse_error := time.Parse(DATE_FORMAT_EXPIRY, e.Expires)
	if parse_error != nil {
		return true
	}

	return time.Now().After(expiry_date)
}

func (e *ServiceContractRuleException) GetRule() string {
	return e.Rule
}

func (e *ServiceContractRuleException) GetService() string {
	return e.Service
}

func (e *ServiceContractRuleException) GetReason() string {
	return e.Reason
}

func (e *ServiceContractRuleException) GetExpires() string {
	return e.Expires
}

func (e *ServiceContractRuleException) GetApprovedBy() string {
	return e.ApprovedBy
}
