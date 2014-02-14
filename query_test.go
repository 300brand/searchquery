package query

import (
	"testing"
)

func TestParse(t *testing.T) {
	for i, test := range tests {
		parsed, err := Parse(test.Input)
		if err != nil {
			t.Errorf("[%d] Error parsing %s: %s", i, test.Input, err)
			continue
		}
		if got := parsed.String(); got != test.String {
			t.Errorf("[%d] Exp: %s", i, test.String)
			t.Errorf("[%d] Got: %s", i, got)
		}
	}
}

func TestString(t *testing.T) {
	for i, test := range tests {
		if got := test.Query.String(); got != test.String {
			t.Errorf("[%d] Exp: %s", i, test.String)
			t.Errorf("[%d] Got: %s", i, got)
		}
	}
}
