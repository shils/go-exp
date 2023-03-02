package operators

import "testing"

type booleanOperands struct {
	l bool
	r bool
}

var booleanPredicateCases = []booleanOperands{
	{false, false},
	{false, true},
	{true, false},
	{true, true},
}

func Test_Or(t *testing.T) {
	for _, test := range booleanPredicateCases {
		expected := test.l || test.r
		if Or(test.l, test.r) != expected {
			t.Errorf("Or(%v, %v) should be %v", test.l, test.r, expected)
		}
	}
}

func Test_And(t *testing.T) {
	for _, test := range booleanPredicateCases {
		expected := test.l && test.r
		if And(test.l, test.r) != expected {
			t.Errorf("And(%v, %v) should be %v", test.l, test.r, expected)
		}
	}
}

func Test_BXor(t *testing.T) {
	for _, test := range booleanPredicateCases {
		expected := test.l != test.r
		if BXor(test.l, test.r) != expected {
			t.Errorf("Xor(%v, %v) should be %v", test.l, test.r, expected)
		}
	}
}
