package query

import (
	"testing"
)

var testSets = map[string]struct {
	tests []testType
	f     func(string) (*Query, error)
}{
	"normal": {parseTests, Parse},
	"greedy": {parseGreedyTests, ParseGreedy},
}

func TestParse(t *testing.T) {
	for p, set := range testSets {
		for i, test := range set.tests {
			parsed, err := set.f(test.Input)
			if err != nil {
				t.Errorf("[%s.%d] Error parsing %s: %s", p, i, test.Input, err)
				continue
			}
			if got := parsed.String(); got != test.String {
				t.Errorf("[%s.%d] Exp: %s", p, i, test.String)
				t.Errorf("[%s.%d] Got: %s", p, i, got)
			}
		}
	}
}

func TestString(t *testing.T) {
	for p, set := range testSets {
		for i, test := range set.tests {
			if got := test.Query.String(); got != test.String {
				t.Errorf("[%s.%d] Exp: %s", p, i, test.String)
				t.Errorf("[%s.%d] Got: %s", p, i, got)
			}
		}
	}
}
