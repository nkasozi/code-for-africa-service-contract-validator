package entities

import "testing"

func TestServiceValidatorMode_String_ReturnsWarnForWarnMode(t *testing.T) {
	mode := WARN

	result := mode.String()

	if result != "warn" {
		t.Errorf("Expected 'warn', got: %s", result)
	}
}

func TestServiceValidatorMode_String_ReturnsEnforceForEnforceMode(t *testing.T) {
	mode := ENFORCE

	result := mode.String()

	if result != "enforce" {
		t.Errorf("Expected 'enforce', got: %s", result)
	}
}

func TestServiceValidatorMode_String_ReturnsUnknownForInvalidMode(t *testing.T) {
	mode := ServiceValidatorMode(999)

	result := mode.String()

	if result != "unknown" {
		t.Errorf("Expected 'unknown', got: %s", result)
	}
}

func TestParseServiceValidatorMode_ParsesWarnCorrectly(t *testing.T) {
	mode, is_valid := ParseServiceValidatorMode("warn")

	if !is_valid {
		t.Error("Expected valid=true for 'warn'")
	}
	if mode != WARN {
		t.Errorf("Expected WARN mode, got: %v", mode)
	}
}

func TestParseServiceValidatorMode_ParsesEnforceCorrectly(t *testing.T) {
	mode, is_valid := ParseServiceValidatorMode("enforce")

	if !is_valid {
		t.Error("Expected valid=true for 'enforce'")
	}
	if mode != ENFORCE {
		t.Errorf("Expected ENFORCE mode, got: %v", mode)
	}
}

func TestParseServiceValidatorMode_ReturnsFalseForInvalidMode(t *testing.T) {
	_, is_valid := ParseServiceValidatorMode("invalid")

	if is_valid {
		t.Error("Expected valid=false for 'invalid'")
	}
}

func TestParseServiceValidatorMode_IsCaseSensitive(t *testing.T) {
	_, is_valid := ParseServiceValidatorMode("WARN")

	if is_valid {
		t.Error("Expected valid=false for uppercase 'WARN' (case-sensitive)")
	}
}
