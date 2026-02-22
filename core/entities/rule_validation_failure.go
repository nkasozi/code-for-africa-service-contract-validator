package entities

type RuleValidationFailure struct {
	RuleName    string
	Found       string
	Need        string
	Examples    string
	FixSuggestion string
}

func NewRuleValidationFailure(
	rule_name string,
	found string,
	need string,
	examples string,
	fix_suggestion string,
) *RuleValidationFailure {
	return &RuleValidationFailure{
		RuleName:      rule_name,
		Found:         found,
		Need:          need,
		Examples:      examples,
		FixSuggestion: fix_suggestion,
	}
}

func (e *RuleValidationFailure) Error() string {
	return e.RuleName + ": validation failed"
}

func (e *RuleValidationFailure) GetFound() string {
	return e.Found
}

func (e *RuleValidationFailure) GetNeed() string {
	return e.Need
}

func (e *RuleValidationFailure) GetExamples() string {
	return e.Examples
}

func (e *RuleValidationFailure) GetFixSuggestion() string {
	return e.FixSuggestion
}
