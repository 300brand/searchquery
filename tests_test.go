package query

var tests = []struct {
	Input  string
	Query  Query
	String string
}{
	{
		Input:  "a b",
		String: ":a :b",
		Query: Query{
			Optional: []SubQuery{
				SubQuery{
					Quote:    QuoteNone,
					Operator: OperatorField,
					Phrase:   "a",
				},
				SubQuery{
					Quote:    QuoteNone,
					Operator: OperatorField,
					Phrase:   "b",
				},
			},
		},
	},
	{
		Input:  "a OR b",
		String: ":a :b",
		Query: Query{
			Optional: []SubQuery{
				SubQuery{
					Quote:    QuoteNone,
					Operator: OperatorField,
					Phrase:   "a",
				},
				SubQuery{
					Quote:    QuoteNone,
					Operator: OperatorField,
					Phrase:   "b",
				},
			},
		},
	},
	{
		Input:  "a AND b",
		String: "+:a +:b",
		Query: Query{
			Required: []SubQuery{
				SubQuery{
					Quote:    QuoteNone,
					Operator: OperatorField,
					Phrase:   "a",
				},
				SubQuery{
					Quote:    QuoteNone,
					Operator: OperatorField,
					Phrase:   "b",
				},
			},
		},
	},
	{
		Input:  "txt~'^foo.*' date>='01.01.2001' date<='02.02.2002'",
		String: "txt~'^foo.*' date>='01.01.2001' date<='02.02.2002'",
		Query: Query{
			Optional: []SubQuery{
				SubQuery{
					Quote:    QuoteSingle,
					Operator: OperatorRegex,
					Field:    "txt",
					Phrase:   "^foo.*",
				},
				SubQuery{
					Quote:    QuoteSingle,
					Operator: OperatorRelGTE,
					Field:    "date",
					Phrase:   "01.01.2001",
				},
				SubQuery{
					Quote:    QuoteSingle,
					Operator: OperatorRelLTE,
					Field:    "date",
					Phrase:   "02.02.2002",
				},
			},
		},
	},
	{
		Input:  "a AND (b OR c) AND NOT d",
		String: "+:a +(:b :c) -:d",
		Query: Query{
			Excluded: []SubQuery{
				SubQuery{
					Quote:    QuoteNone,
					Operator: OperatorField,
					Phrase:   "d",
				},
			},
			Required: []SubQuery{
				SubQuery{
					Quote:    QuoteNone,
					Operator: OperatorField,
					Phrase:   "a",
				},
				SubQuery{
					Quote:    QuoteNone,
					Operator: OperatorSubquery,
					Query: &Query{
						Optional: []SubQuery{
							SubQuery{
								Quote:    QuoteNone,
								Operator: OperatorField,
								Phrase:   "b",
							},
							SubQuery{
								Quote:    QuoteNone,
								Operator: OperatorField,
								Phrase:   "c",
							},
						},
					},
				},
			},
		},
	},
	{
		Input:  "+a +(b c) -d",
		String: "+:a +(:b :c) -:d",
		Query: Query{
			Excluded: []SubQuery{
				SubQuery{
					Quote:    QuoteNone,
					Operator: OperatorField,
					Phrase:   "d",
				},
			},
			Required: []SubQuery{
				SubQuery{
					Quote:    QuoteNone,
					Operator: OperatorField,
					Phrase:   "a",
				},
				SubQuery{
					Quote:    QuoteNone,
					Operator: OperatorSubquery,
					Query: &Query{
						Optional: []SubQuery{
							SubQuery{
								Quote:    QuoteNone,
								Operator: OperatorField,
								Phrase:   "b",
							},
							SubQuery{
								Quote:    QuoteNone,
								Operator: OperatorField,
								Phrase:   "c",
							},
						},
					},
				},
			},
		},
	},
	{
		Input:  "Id#123,444,555,666 AND (b OR c)",
		String: "+Id#123,444,555,666 +(:b :c)",
		Query: Query{
			Required: []SubQuery{
				SubQuery{
					Quote:    QuoteNone,
					Operator: OperatorCSV,
					Field:    "Id",
					Phrase:   "123,444,555,666",
				},
				SubQuery{
					Quote:    QuoteNone,
					Operator: OperatorSubquery,
					Query: &Query{
						Optional: []SubQuery{
							SubQuery{
								Quote:    QuoteNone,
								Operator: OperatorField,
								Phrase:   "b",
							},
							SubQuery{
								Quote:    QuoteNone,
								Operator: OperatorField,
								Phrase:   "c",
							},
						},
					},
				},
			},
		},
	},
	{
		Input:  "+mandatoryWord -excludedWord +field:word \"exact phrase\"",
		String: "+:mandatoryWord +field:word :\"exact phrase\" -:excludedWord",
		Query: Query{
			Excluded: []SubQuery{
				SubQuery{
					Quote:    QuoteNone,
					Operator: OperatorField,
					Phrase:   "excludedWord",
				},
			},
			Optional: []SubQuery{
				SubQuery{
					Quote:    QuoteDouble,
					Operator: OperatorField,
					Phrase:   "exact phrase",
				},
			},
			Required: []SubQuery{
				SubQuery{
					Quote:    QuoteNone,
					Operator: OperatorField,
					Phrase:   "mandatoryWord",
				},
				SubQuery{
					Quote:    QuoteNone,
					Operator: OperatorField,
					Field:    "field",
					Phrase:   "word",
				},
			},
		},
	},
	{
		Input:  "\"Red Hat\" AND Google",
		String: "+:\"Red Hat\" +:Google",
		Query: Query{
			Required: []SubQuery{
				SubQuery{
					Quote:    QuoteDouble,
					Operator: OperatorField,
					Phrase:   "Red Hat",
				},
				SubQuery{
					Quote:    QuoteNone,
					Operator: OperatorField,
					Phrase:   "Google",
				},
			},
		},
	},
	{
		Input:  "Google AND NOT \"Red Hat\"",
		String: "+:Google -:\"Red Hat\"",
		Query: Query{
			Required: []SubQuery{
				SubQuery{
					Quote:    QuoteNone,
					Operator: OperatorField,
					Phrase:   "Google",
				},
			},
			Excluded: []SubQuery{
				SubQuery{
					Quote:    QuoteDouble,
					Operator: OperatorField,
					Phrase:   "Red Hat",
				},
			},
		},
	},
	{
		Input:  "\"Red Hat\" OR \"Fusion IO\"",
		String: ":\"Red Hat\" :\"Fusion IO\"",
		Query: Query{
			Optional: []SubQuery{
				SubQuery{
					Quote:    QuoteDouble,
					Operator: OperatorField,
					Phrase:   "Red Hat",
				},
				SubQuery{
					Quote:    QuoteDouble,
					Operator: OperatorField,
					Phrase:   "Fusion IO",
				},
			},
		},
	},
	{
		Input:  "(\"Cloud Computing\" AND \"Red Hat\") (\"Cloud Computing\" AND \"Fusion IO\")",
		String: "(+:\"Cloud Computing\" +:\"Red Hat\") (+:\"Cloud Computing\" +:\"Fusion IO\")",
		Query: Query{
			Optional: []SubQuery{
				SubQuery{
					Quote:    QuoteNone,
					Operator: OperatorSubquery,
					Query: &Query{
						Required: []SubQuery{
							SubQuery{
								Quote:    QuoteDouble,
								Operator: OperatorField,
								Phrase:   "Cloud Computing",
							},
							SubQuery{
								Quote:    QuoteDouble,
								Operator: OperatorField,
								Phrase:   "Red Hat",
							},
						},
					},
				},
				SubQuery{
					Quote:    QuoteNone,
					Operator: OperatorSubquery,
					Query: &Query{
						Required: []SubQuery{
							SubQuery{
								Quote:    QuoteDouble,
								Operator: OperatorField,
								Phrase:   "Cloud Computing",
							},
							SubQuery{
								Quote:    QuoteDouble,
								Operator: OperatorField,
								Phrase:   "Fusion IO",
							},
						},
					},
				},
			},
		},
	},
	{
		Input:  "\"Cloud Computing\" AND (\"Red Hat\" OR \"Fusion IO\")",
		String: "+:\"Cloud Computing\" +(:\"Red Hat\" :\"Fusion IO\")",
		Query: Query{
			Required: []SubQuery{
				SubQuery{
					Quote:    QuoteDouble,
					Operator: OperatorField,
					Phrase:   "Cloud Computing",
				},
				SubQuery{
					Quote:    QuoteNone,
					Operator: OperatorSubquery,
					Query: &Query{
						Optional: []SubQuery{
							SubQuery{
								Quote:    QuoteDouble,
								Operator: OperatorField,
								Phrase:   "Red Hat",
							},
							SubQuery{
								Quote:    QuoteDouble,
								Operator: OperatorField,
								Phrase:   "Fusion IO",
							},
						},
					},
				},
			},
		},
	},
	{
		Input:  "\"Colon:In the Tech\" AND \"Red Hat\"",
		String: "+:\"Colon:In the Tech\" +:\"Red Hat\"",
		Query: Query{
			Required: []SubQuery{
				SubQuery{
					Quote:    QuoteDouble,
					Operator: OperatorField,
					Phrase:   "Colon:In the Tech",
				},
				SubQuery{
					Quote:    QuoteDouble,
					Operator: OperatorField,
					Phrase:   "Red Hat",
				},
			},
		},
	},
}
