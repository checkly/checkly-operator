/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
