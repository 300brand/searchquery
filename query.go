// Ported from
// http://search.cpan.org/~karman/Search-Query-0.23/lib/Search/Query/Parser.pm
package query

import (
	"regexp"
)

type Query struct {
	Excluded []SubQuery
	Optional []SubQuery
	Required []SubQuery
}

type SubQuery struct {
	Quote    Quote
	Operator Operator
	Field    string
	Phrase   string
	Query    *Query // non-nil when Op == ()
}

type Quote string

type Operator string

const (
	QuoteNone   Quote = ``
	QuoteSingle       = `'`
	QuoteDouble       = `"`

	OperatorCSV      Operator = `#`
	OperatorField             = `:`
	OperatorRegex             = `~`
	OperatorRegexNeg          = `!~`
	OperatorRelE              = `==`
	OperatorRelGT             = `>`
	OperatorRelGTE            = `>=`
	OperatorRelLT             = `<`
	OperatorRelLTE            = `<=`
	OperatorRelNE             = `!=`
	OperatorSubquery          = `()`
)

var Regexes = struct {
	Term, Field, Operator, OperatorNoField, And, Or, Not *regexp.Regexp
}{
	Term:            regexp.MustCompile(`[^\s()]+`),
	Field:           regexp.MustCompile(`[\w]+`),
	Operator:        regexp.MustCompile(`~\d+|==|<=|>=|!=|!:|=~|!~|[:=<>~#]`), // Longest ops first
	OperatorNoField: regexp.MustCompile(`=~|!~|[~:#]`),                        // Ops that admit an empty left operand
	And:             regexp.MustCompile(`\&|AND|ET|UND|E`),
	Or:              regexp.MustCompile(`\||OR|OU|ODER|O`),
	Not:             regexp.MustCompile(`NOT|PAS|NICHT|NON`),
}

func Parse(q string) *Query {
	return &Query{}
}

func parse(q, implicitPlus, parentField, parentOp string) {
	for q != "" {

	}
}
