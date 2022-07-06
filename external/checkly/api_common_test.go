package external

import (
	"testing"
)

func TestCheckValueString(t *testing.T) {
	defaultValue := "foo"
	overrideValue := "bar"

	testValue := checkValueString(overrideValue, defaultValue)

	if testValue != overrideValue {
		t.Errorf("Expected %s, got %s", overrideValue, testValue)
	}

	overrideValue = ""

	testValue = checkValueString(overrideValue, defaultValue)
	if testValue != defaultValue {
		t.Errorf("Expected %s, got %s", overrideValue, testValue)
	}

}

func TestCheckValueInt(t *testing.T) {
	defaultValue := 1
	overrideValue := 2

	testValue := checkValueInt(overrideValue, defaultValue)

	if testValue != overrideValue {
		t.Errorf("Expected %d, got %d", overrideValue, testValue)
	}

	overrideValue = 0

	testValue = checkValueInt(overrideValue, defaultValue)
	if testValue != defaultValue {
		t.Errorf("Expected %d, got %d", overrideValue, testValue)
	}

}

func TestCheckValueArray(t *testing.T) {
	defaultValue := []string{"foo"}
	overrideValue := []string{"foo", "bar"}

	testValue := checkValueArray(overrideValue, defaultValue)

	if len(testValue) != len(overrideValue) {
		t.Errorf("Expected %d, got %d", len(overrideValue), len(testValue))
	}

	overrideValue = []string{}

	testValue = checkValueArray(overrideValue, defaultValue)
	if len(testValue) != len(defaultValue) {
		t.Errorf("Expected %d, got %d", len(overrideValue), len(testValue))
	}

}
