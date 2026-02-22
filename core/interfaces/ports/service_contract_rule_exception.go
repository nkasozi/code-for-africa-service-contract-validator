package ports

type IServiceContractRuleException interface {
	AppliesToRuleAndService(rule_name string, service_name string) bool
	IsExpired() bool
	GetRule() string
	GetService() string
	GetReason() string
	GetExpires() string
	GetApprovedBy() string
}