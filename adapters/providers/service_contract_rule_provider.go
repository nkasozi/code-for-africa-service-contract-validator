package providers

import "github.com/nkasozi/code-for-africa-service-contract-validator/core/interfaces/ports"

type ServiceContractRuleProvider struct {
	allServiceContractRules []ports.IServiceContractRule
}

func NewServiceContractRuleProvider(allServiceContractRules []ports.IServiceContractRule) *ServiceContractRuleProvider {
	return &ServiceContractRuleProvider{
		allServiceContractRules: allServiceContractRules,
	}
}

func (p *ServiceContractRuleProvider) LoadServiceContractRules() ([]ports.IServiceContractRule, error) {
	return p.allServiceContractRules, nil
}
