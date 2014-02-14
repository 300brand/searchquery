package query

import (
	"testing"
)

func TestString(t *testing.T) {
	for i, test := range tests {
		if got := test.Query.String(); got != test.String {
			t.Errorf("[%d] Exp: %s", i, test.String)
			t.Errorf("[%d] Got: %s", i, got)
		}
	}
}
