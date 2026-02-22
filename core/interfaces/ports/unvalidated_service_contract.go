package ports

type IUnvalidatedServiceContract interface {
	GetSchemaVersion() string
	GetServiceName() string
	GetOwner() string
	GetEnv() string
	GetDataSensitivity() string
	GetCostCenter() string
	GetAlerts() []string
	GetOtherFields() map[string]interface{}
}
