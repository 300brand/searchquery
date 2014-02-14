// Ported from
// http://search.cpan.org/~karman/Search-Query-0.23/lib/Search/Query/Parser.pm
package query

import (
	"fmt"
	"regexp"
	"strings"
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
	Value    string
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

	PrefixExcluded string = `-`
	PrefixOptional        = ``
	PrefixRequired        = `+`
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

func (q Query) String() string {
	buf := make([]string, 0, len(q.Required)+len(q.Optional)+len(q.Excluded))
	for _, sq := range q.Required {
		buf = append(buf, PrefixRequired+sq.String())
	}
	for _, sq := range q.Optional {
		buf = append(buf, PrefixOptional+sq.String())
	}
	for _, sq := range q.Excluded {
		buf = append(buf, PrefixExcluded+sq.String())
	}
	return strings.Join(buf, " ")
}

func (sq SubQuery) String() string {
	if sq.Operator == OperatorSubquery {
		return "(" + sq.Query.String() + ")"
	}
	return fmt.Sprintf("%s%s%s%s%s", sq.Field, sq.Operator, sq.Quote, sq.Value, sq.Quote)
}

func parse(q, implicitPlus, parentField, parentOp string) {
	for q != "" {

	}
}
