package cmd

import (
	"testing"
)

func TestIsValidVersion(t *testing.T) {
	tests := []struct {
		version string
		valid   bool
	}{
		{"21", true},
		{"17.0.8", true},
		{"11.0.20", true},
		{"8", true},
		{"21.0", true},
		{"17.0.8.1", false},  // Not a standard JDK version format (too many parts)
		{"", false},
		{"invalid", false},
		{"21.invalid", false},
		{"21.0.invalid", false},
		{"v21", false},
		{"21-ea", false},
	}

	for _, test := range tests {
		// Call the isValidVersion function from install.go
		result := isValidVersion(test.version) 
		if result != test.valid {
			t.Errorf("isValidVersion(%s) = %v, expected %v", test.version, result, test.valid)
		}
	}
}
