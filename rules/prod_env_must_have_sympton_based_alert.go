package rules

import (
	"regexp"
	"strings"

	"github.com/nkasozi/code-for-africa-service-contract-validator/core/entities"
	"github.com/nkasozi/code-for-africa-service-contract-validator/core/interfaces/ports"
)

const (
	RULE_NAME_PROD_SYMPTOM_ALERT = "prod_env_must_have_symptom_based_alert"

	PATTERN_5XX_ERRORS     = `5\d{2}|5xx`
	PATTERN_LATENCY        = `p\d{2}`
	PATTERN_HEALTH_CHECK   = `health[-_]?check`
	PATTERN_GENERIC_ERRORS = `fail|error|breach|timeout`
)

var patterns = []string{
	PATTERN_5XX_ERRORS,
	PATTERN_LATENCY,
	PATTERN_HEALTH_CHECK,
	PATTERN_GENERIC_ERRORS,
}

type ProdEnvMustHaveSymptomBasedAlert struct {
	symptom_regexes []*regexp.Regexp
}

func NewProdEnvMustHaveSymptomBasedAlert() *ProdEnvMustHaveSymptomBasedAlert {
	compiled := make([]*regexp.Regexp, 0, len(patterns))
	for _, p := range patterns {
		compiled = append(compiled, regexp.MustCompile("(?i)"+p))
	}

	return &ProdEnvMustHaveSymptomBasedAlert{
		symptom_regexes: compiled,
	}
}

func (r *ProdEnvMustHaveSymptomBasedAlert) GetRuleName() string {
	return RULE_NAME_PROD_SYMPTOM_ALERT
}

func (r *ProdEnvMustHaveSymptomBasedAlert) IsRuleSatisfied(service_contract ports.IUnvalidatedServiceContract) error {
	if strings.ToLower(service_contract.GetEnv()) != entities.EnvironmentProduction.String() {
		return nil
	}

	alerts := service_contract.GetAlerts()
	found_alerts := formatAlertsList(alerts)

	for _, alert_name := range alerts {
		if r.matchesSymptom(alert_name) {
			return nil
		}
	}

	return entities.NewRuleValidationFailure(
		RULE_NAME_PROD_SYMPTOM_ALERT,
		"alerts: "+found_alerts,
		"Production services require at least one symptom-based alert",
		"high_error_rate, p99_latency_breach, health_check_failed",
		"Add a symptom-based alert (5xx errors, latency metrics, or health checks) to your alerts list",
	)
}

func (r *ProdEnvMustHaveSymptomBasedAlert) matchesSymptom(alert_name string) bool {
	for _, re := range r.symptom_regexes {
		if re.MatchString(alert_name) {
			return true
		}
	}
	return false
}

func formatAlertsList(alerts []string) string {
	if len(alerts) == 0 {
		return "[none]"
	}
	return "[" + strings.Join(alerts, ", ") + "]"
}