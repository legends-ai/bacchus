package config

import "testing"

func TestIsValidKey(t *testing.T) {
	for _, test := range []struct {
		Key      string
		Expected bool
	}{
		{"abcdefab-1234-1234-1234-1234567890ab", true},
		{"abcdefab-1234-1234-123h-1234567890ab", false},
	} {
		res := isValidKey(test.Key)
		if test.Expected != res {
			t.Errorf("Expected %v, got %v", test.Expected, res)
		}
	}
}
