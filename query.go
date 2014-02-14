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

var reNot = `NOT|PAS|NICHT|NON`

var Regexes = struct {
	Term, Field, Operator, OperatorNoField, And, Or, Not, Sign, SignWord *regexp.Regexp
}{
	Term:            regexp.MustCompile(`[^\s()]+`),
	Field:           regexp.MustCompile(`[\w]+`),
	Operator:        regexp.MustCompile(`~\d+|==|<=|>=|!=|!:|=~|!~|[:=<>~#]`), // Longest ops first
	OperatorNoField: regexp.MustCompile(`=~|!~|[~:#]`),                        // Ops that admit an empty left operand
	And:             regexp.MustCompile(`\&|AND|ET|UND|E`),
	Or:              regexp.MustCompile(`\||OR|OU|ODER|O`),
	Not:             regexp.MustCompile(reNot),
	SignWord:        regexp.MustCompile(`^(` + reNot + `)\b\s*`),
}

func Parse(s string) (q *Query, err error) {
	q, _, err = parse(s, false, "", OperatorField)
	return
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

func parse(s string, implicitPlus bool, parentField, parentOperator string) (q *Query, remaining string, err error) {
	q = new(Query)
	loops := 0
	for s != "" {
		fmt.Printf("Loop: %d\n", loops)
		loops++

		prefix := PrefixOptional
		if implicitPlus {
			prefix = PrefixRequired
		}
		//field := parentField
		//operator := parentOperator

		if s[0] == ')' {
			break
		}
		// Parse prefix ('+', '-' or 'NOT')
		if sign := s[0]; sign == '+' || sign == '-' {
			prefix = s[0:1]
		} else if sm := Regexes.SignWord.FindStringSubmatch(s); len(sm) > 0 {
			prefix = PrefixExcluded
		}

		fmt.Println("Prefix", prefix)
		break
	}

	if len(q.Required) == 0 && len(q.Optional) == 0 {
		err = fmt.Errorf("No positive value in query: %s", s)
	}
	return
}
