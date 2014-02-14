#!/usr/bin/perl -w

use strict;

use Switch;
use Search::QueryParser;
use Data::Dumper;

my $header = <<HEADER;
package query

var tests = []struct {
	Input  string
	Query  Query
	String string
}{
HEADER
my $fmt = <<FMT;
	{
		Input:  %s,
		String: %s,
		Query: %s	},
FMT
my $footer = <<FOOTER;
}
FOOTER

my @tests = (
	{ query => 'a b' },
	{ query => 'a OR b' },
	{ query => 'a AND b' },
	{ query => "txt~'^foo.*' date>='01.01.2001' date<='02.02.2002'" },
	# Following two should be equivalent:
	{ query => "a AND (b OR c) AND NOT d" },
	{ query => "+a +(b c) -d" },
	{ query => "Id#123,444,555,666 AND (b OR c)" },
	{ query => '+mandatoryWord -excludedWord +field:word "exact phrase"' },
	{ query => '"Red Hat" AND Google' },
	{ query => 'Google AND NOT "Red Hat"' },
	{ query => '"Red Hat" OR "Fusion IO"' },
	{ query => '("Cloud Computing" AND "Red Hat") ("Cloud Computing" AND "Fusion IO")' },
	{ query => '"Cloud Computing" AND ("Red Hat" OR "Fusion IO")' },
	{ query => '"Colon:In the Tech" AND "Red Hat"' },
	# Should error: cannot mix AND/OR in requests; use parentheses
	# { query => '"Cloud Computing" AND "Red Hat" OR "Cloud Computing" AND "Fusion IO"' },
);

$Data::Dumper::Useqq = 1;
$Data::Dumper::Varname = '';

print($header);
foreach (@tests) {
	my $s = "$_->{query}";
	my $qp = new Search::QueryParser;
	my $query;
	$query = $qp->parse($s) or $query = "Error in query : " . $qp->err;
	printf($fmt, sdump($s), sdump($qp->unparse($query)), go_Query($query, "\t\t"));
}
print($footer);

sub sdump {
	return substr(Dumper(shift), 5, -2)
}

sub go_Query {
	my $q = shift;
	my $indent = shift;
	my %keymap = (
		Required => '+',
		Optional => '',
		Excluded => '-',
	);

	my @subQ;
	push @subQ, "Query{\n";

	while ((my $name, my $prefix) = each %keymap) {
		next if not $q->{$prefix};

		push @subQ, "$indent\t$name: []SubQuery{\n";
		push @subQ, go_SubQuery($_, "$indent\t\t") foreach @{$q->{$prefix}};
		push @subQ, "$indent\t},\n";
	}
	push @subQ, "$indent},\n";
	return join "", @subQ;
}

sub go_SubQuery {
	my $subQ = shift;
	my $indent = shift;
	my $str = "${indent}SubQuery{\n";

	# Convert quote to a constant
	my $quote;
	switch ($subQ->{quote} || "") {
		case '"' { $quote = "QuoteDouble"; }
		case "'" { $quote = "QuoteSingle"; }
		else     { $quote = "QuoteNone";   }
	}
	$str .= "$indent\tQuote:    $quote,\n";

	# Convert operator to a constant
	my $op;
	switch ($subQ->{op}) {
		case "#"  { $op = "OperatorCSV";          }
		case ":"  { $op = "OperatorField";        }
		case "~"  { $op = "OperatorRegex";        }
		case "!~" { $op = "OperatorRegexNeg";     }
		case "==" { $op = "OperatorRelE";         }
		case ">"  { $op = "OperatorRelGT";        }
		case ">=" { $op = "OperatorRelGTE";       }
		case "<"  { $op = "OperatorRelLT";        }
		case "<=" { $op = "OperatorRelLTE";       }
		case "!=" { $op = "OperatorRelNE";        }
		case "()" { $op = "OperatorSubquery";     }
		else      { print "wtf is $subQ->{op}\n"; }
	}
	$str .= "$indent\tOperator: $op,\n";

	# Add field name if applicable
	if ($subQ->{field}) {
		$str .= "$indent\tField:    " . sdump($subQ->{field}) . ",\n";
	}

	if ($subQ->{op} eq '()') {
		$str .= "$indent\tQuery: &" . go_Query($subQ->{value}, "$indent\t");
	} else {
		$str .= "$indent\tPhrase:   " . sdump($subQ->{value}) . ",\n";
	}

	# return  "(" . go_Query($subQ->{value}) . ")"  if $subQ->{op} eq '()';
	# my $quote = $subQ->{quote} || "";
	# return "$subQ->{field}$subQ->{op}$quote$subQ->{value}$quote";

	$str .= "$indent},\n";
}
