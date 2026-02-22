package ports

type ILoggerProvider interface{
	LogInfo(message string)
	LogWarning(message string)
	LogError(message string)
}