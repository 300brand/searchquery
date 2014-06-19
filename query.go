// Ported from
// http://search.cpan.org/~karman/Search-Query-0.23/lib/Search/Query/Parser.pm
package searchquery

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
	OperatorNone              = ``
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

var R = struct {
	Prefix, PrefixWord                   *regexp.Regexp
	FieldOpQQ, FieldOpQ, FieldOp, OpOnly *regexp.Regexp
	TermQQ, TermQ, SingleTerm            *regexp.Regexp
	OpenParen, CloseParen                *regexp.Regexp
	BoolAnd, BoolOr                      *regexp.Regexp
}{
	Prefix:     regexp.MustCompile(`^(\+|-)\s*`),
	PrefixWord: regexp.MustCompile(`^(` + reNot + `)\b\s*`),
	FieldOpQQ:  regexp.MustCompile(`^"(` + reField + `)"\s*(` + reOperator + `)\s*`),
	FieldOpQ:   regexp.MustCompile(`^'(` + reField + `)'\s*(` + reOperator + `)\s*`),
	FieldOp:    regexp.MustCompile(`^(` + reField + `)\s*(` + reOperator + `)\s*`),
	OpOnly:     regexp.MustCompile(`^()(` + reOperatorNoField + `)\s*`),
	TermQQ:     regexp.MustCompile(`^(")([^"]*?)"\s*`),
	TermQ:      regexp.MustCompile(`^(')([^']*?)'\s*`),
	SingleTerm: regexp.MustCompile(`^()(` + reTerm + `)\s*`),
	OpenParen:  regexp.MustCompile(`^\(\s*`),
	CloseParen: regexp.MustCompile(`^\)\s*`),
	BoolAnd:    regexp.MustCompile(`^(` + reAnd + `)\b\s*`),
	BoolOr:     regexp.MustCompile(`^(` + reOr + `)\b\s*`),
}

var (
	reTerm            = `[^\s()]+`
	reAnd             = `\&|AND|ET|UND|E`
	reOr              = `\||OR|OU|ODER|O`
	reNot             = `NOT|PAS|NICHT|NON`
	reField           = `[\w]+`
	reOperator        = `~\d+|==|<=|>=|!=|!:|=~|!~|[:=<>~#]`
	reOperatorNoField = `=~|!~|[~:#]`
)

var (
	fieldOperators = []*regexp.Regexp{
		R.FieldOpQQ,
		R.FieldOpQ,
		R.FieldOp,
		R.OpOnly,
	}
	terms = []*regexp.Regexp{
		R.TermQQ,
		R.TermQ,
		R.SingleTerm,
	}
)

func Parse(s string) (q *Query, err error) {
	q, _, err = parse(s, PrefixOptional, "", OperatorField)
	return
}

func ParseGreedy(s string) (q *Query, err error) {
	q, _, err = parse(s, PrefixRequired, "", OperatorField)
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

func parse(s string, defaultPrefix string, parentField string, parentOperator Operator) (q *Query, remaining string, err error) {
	q = new(Query)
	preBool := ""
	for s != "" {
		prefix := defaultPrefix
		subQuery := SubQuery{
			Field:    parentField,
			Operator: parentOperator,
		}
		var sm []string

		// return from recursive call if meeting a ')'
		if s[0] == ')' {
			break
		}

		// Parse prefix ('+', '-' or 'NOT')
		if sm = R.Prefix.FindStringSubmatch(s); len(sm) > 0 {
			prefix = sm[1]
			s = s[len(sm[0]):]
		} else if sm = R.PrefixWord.FindStringSubmatch(s); len(sm) > 0 {
			prefix = PrefixExcluded
			s = s[len(sm[0]):]
		}

		// Parse field name and operator
		for _, re := range fieldOperators {
			sm = re.FindStringSubmatch(s)
			if len(sm) == 0 {
				continue
			}
			subQuery.Field, subQuery.Operator = sm[1], Operator(sm[2])
			if parentField != "" {
				err = fmt.Errorf("Field '%s' inside '%s'", subQuery.Field, parentField)
			}
			s = s[len(sm[0]):]
			break
		}

		// Term matching
		for _, re := range terms {
			sm = re.FindStringSubmatch(s)
			if len(sm) == 0 {
				continue
			}
			subQuery.Quote = Quote(sm[1])
			subQuery.Value = sm[2]
			s = s[len(sm[0]):]
			goto BooleanOperators
		}

		// Parenthesis matching
		if sm = R.OpenParen.FindStringSubmatch(s); len(sm) > 0 {
			subQuery.Query, s, err = parse(s[len(sm[0]):], defaultPrefix, subQuery.Field, subQuery.Operator)
			// Important not to pass OperatorSubquery into the sub-parse
			subQuery.Operator = OperatorSubquery
			if err != nil {
				return
			}
			p := R.CloseParen.FindString(s)
			if p == "" {
				err = fmt.Errorf("No matching )")
			}
			s = s[len(p):]
		}

		if subQuery.Operator == OperatorNone {
			err = fmt.Errorf("Unexpected string in query: %s", s)
			return
		}

		// Boolean Operators
	BooleanOperators:

		postBool := ""
		if and, or := R.BoolAnd.FindString(s), R.BoolOr.FindString(s); and != "" {
			postBool = "AND"
			s = s[len(and):]
		} else if or != "" {
			postBool = "OR"
			s = s[len(or):]
		}
		if preBool != "" && postBool != "" && preBool != postBool {
			err = fmt.Errorf("Cannot mix AND/OR; use parenthesis")
			return
		}
		Bool := preBool
		if preBool == "" {
			Bool = postBool
		}
		// Set for next loop:
		preBool = postBool

		// Insert SubQuery into Query struct
		switch {
		case prefix == PrefixRequired && Bool == "OR":
			prefix = PrefixOptional
		case prefix == PrefixOptional && Bool == "AND":
			prefix = PrefixRequired
		case prefix == PrefixExcluded && Bool == "OR":
			err = fmt.Errorf("Operands of OR cannot have - or NOT prefix")
			return
		}
		switch prefix {
		case PrefixRequired:
			q.Required = append(q.Required, subQuery)
		case PrefixOptional:
			q.Optional = append(q.Optional, subQuery)
		case PrefixExcluded:
			q.Excluded = append(q.Excluded, subQuery)
		default:
			err = fmt.Errorf("Invalid prefix: %s", prefix)
			return
		}
	}

	if len(q.Required) == 0 && len(q.Optional) == 0 {
		err = fmt.Errorf("No positive value in query: %s", s)
	}
	remaining = s
	return
}
