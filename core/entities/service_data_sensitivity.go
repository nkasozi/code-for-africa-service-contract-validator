package entities

import (
	"fmt"
	"strings"
)

const (
	SENSITIVITY_LOW     = "low"
	SENSITIVITY_MEDIUM  = "medium"
	SENSITIVITY_HIGH    = "high"
)

type DataSensitivity string

const (
	SensitivityLow    DataSensitivity = SENSITIVITY_LOW
	SensitivityMedium DataSensitivity = SENSITIVITY_MEDIUM
	SensitivityHigh   DataSensitivity = SENSITIVITY_HIGH
)

func (s DataSensitivity) String() string {
	return string(s)
}

func GetAllowedSensitivityStrings() []string {
	return []string{
		SENSITIVITY_LOW,
		SENSITIVITY_MEDIUM,
		SENSITIVITY_HIGH,
	}
}

func ParseDataSensitivity(raw_value string) (DataSensitivity, error) {
	normalized := strings.ToLower(strings.TrimSpace(raw_value))

	switch normalized {
	case SENSITIVITY_LOW:
		return SensitivityLow, nil
	case SENSITIVITY_MEDIUM:
		return SensitivityMedium, nil
	case SENSITIVITY_HIGH:
		return SensitivityHigh, nil
	default:
		return "", fmt.Errorf("invalid sensitivity")
	}
}