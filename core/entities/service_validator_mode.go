package entities

type ServiceValidatorMode int

const (
	WARN ServiceValidatorMode = iota
	ENFORCE
)

func (m ServiceValidatorMode) String() string {
	switch m {
	case WARN:
		return "warn"
	case ENFORCE:
		return "enforce"
	default:
		return "unknown"
	}
}

func ParseServiceValidatorMode(mode_string string) (ServiceValidatorMode, bool) {
	switch mode_string {
	case "warn":
		return WARN, true
	case "enforce":
		return ENFORCE, true
	default:
		return WARN, false
	}
}