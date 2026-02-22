package ports

type IServiceContractRulesProvider interface {
	LoadServiceContractRules() ([]IServiceContractRule, error)
}
