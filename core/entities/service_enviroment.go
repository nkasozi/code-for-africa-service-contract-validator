package entities

import (
	"fmt"
	"strings"
)

const (
	ENV_STRING_DEV        = "dev"
	ENV_STRING_STAGING    = "staging"
	ENV_STRING_PRODUCTION = "prod"
	ENV_STRING_UNKNOWN    = "unknown"
)

type ServiceEnvironment string

const (
	EnvironmentDev        ServiceEnvironment = ENV_STRING_DEV
	EnvironmentStaging    ServiceEnvironment = ENV_STRING_STAGING
	EnvironmentProduction ServiceEnvironment = ENV_STRING_PRODUCTION
)

func (e ServiceEnvironment) String() string {
	return string(e)
}

func GetAllowedEnvironments() []ServiceEnvironment {
	return []ServiceEnvironment{
		EnvironmentDev,
		EnvironmentStaging,
		EnvironmentProduction,
	}
}

func GetAllowedEnvironmentStrings() []string {
	allowed := GetAllowedEnvironments()
	result := make([]string, len(allowed))

	for i, env := range allowed {
		result[i] = env.String()
	}

	return result
}

func ParseServiceEnvironment(raw_value string) (ServiceEnvironment, error) {
	normalized := strings.ToLower(strings.TrimSpace(raw_value))

	switch normalized {
	case ENV_STRING_DEV:
		return EnvironmentDev, nil
	case ENV_STRING_STAGING:
		return EnvironmentStaging, nil
	case ENV_STRING_PRODUCTION:
		return EnvironmentProduction, nil
	default:
		return "", fmt.Errorf("unknown environment: %s", raw_value)
	}
}
