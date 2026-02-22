package ports

type IServiceContractRuleExceptionsProvider interface {
	LoadServiceRuleExceptions() ([]IServiceContractRuleException, error)
}
