package cmd

import (
	"testing"
)

func TestIsLTSVersion(t *testing.T) {
	tests := []struct {
		version int
		isLTS   bool
	}{
		{8, true},
		{11, true},
		{17, true},
		{21, true},
		{24, true},  // Future LTS (21 + 3*1)
		{27, true},  // Future LTS (21 + 3*2)
		{30, true},  // Future LTS (21 + 3*3)
		{9, false},
		{10, false},
		{12, false},
		{13, false},
		{14, false},
		{15, false},
		{16, false},
		{18, false},
		{19, false},
		{20, false},
		{22, false},
		{23, false},  // Not LTS (21 + 2, not divisible by 3)
		{25, false},  // Not LTS (21 + 4, not divisible by 3)
		{26, false},  // Not LTS (21 + 5, not divisible by 3)
	}

	for _, test := range tests {
		result := isLTSVersion(test.version)
		if result != test.isLTS {
			t.Errorf("isLTSVersion(%d) = %v, expected %v", test.version, result, test.isLTS)
		}
	}
}
