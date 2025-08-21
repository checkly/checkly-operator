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

package checkly

import (
	"errors"
	"testing"
)

func TestIsNotFoundError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
		{
			name:     "404 error",
			err:      errors.New("404 Not Found"),
			expected: true,
		},
		{
			name:     "not found error",
			err:      errors.New("resource not found"),
			expected: true,
		},
		{
			name:     "does not exist error",
			err:      errors.New("check does not exist"),
			expected: true,
		},
		{
			name:     "mixed case 404",
			err:      errors.New("Error: 404 - Resource Not Found"),
			expected: true,
		},
		{
			name:     "mixed case not found",
			err:      errors.New("Item NOT FOUND in database"),
			expected: true,
		},
		{
			name:     "other error",
			err:      errors.New("internal server error"),
			expected: false,
		},
		{
			name:     "500 error",
			err:      errors.New("500 Internal Server Error"),
			expected: false,
		},
		{
			name:     "authorization error",
			err:      errors.New("401 Unauthorized"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isNotFoundError(tt.err)
			if result != tt.expected {
				t.Errorf("isNotFoundError(%v) = %v, expected %v", tt.err, result, tt.expected)
			}
		})
	}
}
