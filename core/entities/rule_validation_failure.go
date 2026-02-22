package entities

type RuleValidationFailure struct {
	RuleName      string
	Issue         string
	Found         string
	Need          string
	Examples      string
	FixSuggestion string
}

func NewRuleValidationFailure(
	rule_name string,
	issue string,
	found string,
	need string,
	examples string,
	fix_suggestion string,
) *RuleValidationFailure {
	return &RuleValidationFailure{
		RuleName:      rule_name,
		Issue:         issue,
		Found:         found,
		Need:          need,
		Examples:      examples,
		FixSuggestion: fix_suggestion,
	}
}

func (e *RuleValidationFailure) Error() string {
	return e.RuleName + ": " + e.Issue
}

func (e *RuleValidationFailure) GetIssue() string {
	return e.Issue
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
