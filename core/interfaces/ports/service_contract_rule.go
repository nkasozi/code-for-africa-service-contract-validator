package ports

type IServiceContractRule interface {
	IsRuleSatisfied(serviceContract IUnvalidatedServiceContract) error
	GetRuleName() string
}
