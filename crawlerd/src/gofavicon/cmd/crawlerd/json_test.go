package main

import (
	"testing"
)

func TestBalancedJson(t *testing.T) {
	var input = map[string]bool {
		"{}": true,
		"[]": true,
		"{[],[],[]}": true,
		"{{}}": true,
		"{{{[]}}}": true,
		"{{[}}": false,
		"{{}": false,
	}

	for s, e := range input {
		v := isBalancedBrackets([]byte(s))
		if v != e {
			t.Errorf("isBalancedBracked failed on value -- %s. Expected: %t, actual: %t", s, e, v)
		}
	}
}