package searchquery

type testType struct {
	Input  string
	Query  Query
	String string
}

var parseTests = []testType{
	{
		Input:  "a b",
		String: ":a :b",
		Query: Query{
			Optional: []SubQuery{
				SubQuery{
					Quote:    QuoteNone,
					Operator: OperatorField,
					Value:    "a",
				},
				SubQuery{
					Quote:    QuoteNone,
					Operator: OperatorField,
					Value:    "b",
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
					Value:    "a",
				},
				SubQuery{
					Quote:    QuoteNone,
					Operator: OperatorField,
					Value:    "b",
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
					Value:    "a",
				},
				SubQuery{
					Quote:    QuoteNone,
					Operator: OperatorField,
					Value:    "b",
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
					Value:    "^foo.*",
				},
				SubQuery{
					Quote:    QuoteSingle,
					Operator: OperatorRelGTE,
					Field:    "date",
					Value:    "01.01.2001",
				},
				SubQuery{
					Quote:    QuoteSingle,
					Operator: OperatorRelLTE,
					Field:    "date",
					Value:    "02.02.2002",
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
					Value:    "d",
				},
			},
			Required: []SubQuery{
				SubQuery{
					Quote:    QuoteNone,
					Operator: OperatorField,
					Value:    "a",
				},
				SubQuery{
					Quote:    QuoteNone,
					Operator: OperatorSubquery,
					Query: &Query{
						Optional: []SubQuery{
							SubQuery{
								Quote:    QuoteNone,
								Operator: OperatorField,
								Value:    "b",
							},
							SubQuery{
								Quote:    QuoteNone,
								Operator: OperatorField,
								Value:    "c",
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
			Required: []SubQuery{
				SubQuery{
					Quote:    QuoteNone,
					Operator: OperatorField,
					Value:    "a",
				},
				SubQuery{
					Quote:    QuoteNone,
					Operator: OperatorSubquery,
					Query: &Query{
						Optional: []SubQuery{
							SubQuery{
								Quote:    QuoteNone,
								Operator: OperatorField,
								Value:    "b",
							},
							SubQuery{
								Quote:    QuoteNone,
								Operator: OperatorField,
								Value:    "c",
							},
						},
					},
				},
			},
			Excluded: []SubQuery{
				SubQuery{
					Quote:    QuoteNone,
					Operator: OperatorField,
					Value:    "d",
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
					Value:    "123,444,555,666",
				},
				SubQuery{
					Quote:    QuoteNone,
					Operator: OperatorSubquery,
					Query: &Query{
						Optional: []SubQuery{
							SubQuery{
								Quote:    QuoteNone,
								Operator: OperatorField,
								Value:    "b",
							},
							SubQuery{
								Quote:    QuoteNone,
								Operator: OperatorField,
								Value:    "c",
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
			Required: []SubQuery{
				SubQuery{
					Quote:    QuoteNone,
					Operator: OperatorField,
					Value:    "mandatoryWord",
				},
				SubQuery{
					Quote:    QuoteNone,
					Operator: OperatorField,
					Field:    "field",
					Value:    "word",
				},
			},
			Excluded: []SubQuery{
				SubQuery{
					Quote:    QuoteNone,
					Operator: OperatorField,
					Value:    "excludedWord",
				},
			},
			Optional: []SubQuery{
				SubQuery{
					Quote:    QuoteDouble,
					Operator: OperatorField,
					Value:    "exact phrase",
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
					Value:    "Red Hat",
				},
				SubQuery{
					Quote:    QuoteNone,
					Operator: OperatorField,
					Value:    "Google",
				},
			},
		},
	},
	{
		Input:  "Google AND NOT \"Red Hat\"",
		String: "+:Google -:\"Red Hat\"",
		Query: Query{
			Excluded: []SubQuery{
				SubQuery{
					Quote:    QuoteDouble,
					Operator: OperatorField,
					Value:    "Red Hat",
				},
			},
			Required: []SubQuery{
				SubQuery{
					Quote:    QuoteNone,
					Operator: OperatorField,
					Value:    "Google",
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
					Value:    "Red Hat",
				},
				SubQuery{
					Quote:    QuoteDouble,
					Operator: OperatorField,
					Value:    "Fusion IO",
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
								Value:    "Cloud Computing",
							},
							SubQuery{
								Quote:    QuoteDouble,
								Operator: OperatorField,
								Value:    "Red Hat",
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
								Value:    "Cloud Computing",
							},
							SubQuery{
								Quote:    QuoteDouble,
								Operator: OperatorField,
								Value:    "Fusion IO",
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
					Value:    "Cloud Computing",
				},
				SubQuery{
					Quote:    QuoteNone,
					Operator: OperatorSubquery,
					Query: &Query{
						Optional: []SubQuery{
							SubQuery{
								Quote:    QuoteDouble,
								Operator: OperatorField,
								Value:    "Red Hat",
							},
							SubQuery{
								Quote:    QuoteDouble,
								Operator: OperatorField,
								Value:    "Fusion IO",
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
					Value:    "Colon:In the Tech",
				},
				SubQuery{
					Quote:    QuoteDouble,
					Operator: OperatorField,
					Value:    "Red Hat",
				},
			},
		},
	},
}
var parseGreedyTests = []testType{
	{
		Input:  "a b",
		String: "+:a +:b",
		Query: Query{
			Required: []SubQuery{
				SubQuery{
					Quote:    QuoteNone,
					Operator: OperatorField,
					Value:    "a",
				},
				SubQuery{
					Quote:    QuoteNone,
					Operator: OperatorField,
					Value:    "b",
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
					Value:    "a",
				},
				SubQuery{
					Quote:    QuoteNone,
					Operator: OperatorField,
					Value:    "b",
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
					Value:    "a",
				},
				SubQuery{
					Quote:    QuoteNone,
					Operator: OperatorField,
					Value:    "b",
				},
			},
		},
	},
	{
		Input:  "txt~'^foo.*' date>='01.01.2001' date<='02.02.2002'",
		String: "+txt~'^foo.*' +date>='01.01.2001' +date<='02.02.2002'",
		Query: Query{
			Required: []SubQuery{
				SubQuery{
					Quote:    QuoteSingle,
					Operator: OperatorRegex,
					Field:    "txt",
					Value:    "^foo.*",
				},
				SubQuery{
					Quote:    QuoteSingle,
					Operator: OperatorRelGTE,
					Field:    "date",
					Value:    "01.01.2001",
				},
				SubQuery{
					Quote:    QuoteSingle,
					Operator: OperatorRelLTE,
					Field:    "date",
					Value:    "02.02.2002",
				},
			},
		},
	},
	{
		Input:  "a AND (b OR c) AND NOT d",
		String: "+:a +(:b :c) -:d",
		Query: Query{
			Required: []SubQuery{
				SubQuery{
					Quote:    QuoteNone,
					Operator: OperatorField,
					Value:    "a",
				},
				SubQuery{
					Quote:    QuoteNone,
					Operator: OperatorSubquery,
					Query: &Query{
						Optional: []SubQuery{
							SubQuery{
								Quote:    QuoteNone,
								Operator: OperatorField,
								Value:    "b",
							},
							SubQuery{
								Quote:    QuoteNone,
								Operator: OperatorField,
								Value:    "c",
							},
						},
					},
				},
			},
			Excluded: []SubQuery{
				SubQuery{
					Quote:    QuoteNone,
					Operator: OperatorField,
					Value:    "d",
				},
			},
		},
	},
	{
		Input:  "+a +(b c) -d",
		String: "+:a +(+:b +:c) -:d",
		Query: Query{
			Excluded: []SubQuery{
				SubQuery{
					Quote:    QuoteNone,
					Operator: OperatorField,
					Value:    "d",
				},
			},
			Required: []SubQuery{
				SubQuery{
					Quote:    QuoteNone,
					Operator: OperatorField,
					Value:    "a",
				},
				SubQuery{
					Quote:    QuoteNone,
					Operator: OperatorSubquery,
					Query: &Query{
						Required: []SubQuery{
							SubQuery{
								Quote:    QuoteNone,
								Operator: OperatorField,
								Value:    "b",
							},
							SubQuery{
								Quote:    QuoteNone,
								Operator: OperatorField,
								Value:    "c",
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
					Value:    "123,444,555,666",
				},
				SubQuery{
					Quote:    QuoteNone,
					Operator: OperatorSubquery,
					Query: &Query{
						Optional: []SubQuery{
							SubQuery{
								Quote:    QuoteNone,
								Operator: OperatorField,
								Value:    "b",
							},
							SubQuery{
								Quote:    QuoteNone,
								Operator: OperatorField,
								Value:    "c",
							},
						},
					},
				},
			},
		},
	},
	{
		Input:  "+mandatoryWord -excludedWord +field:word \"exact phrase\"",
		String: "+:mandatoryWord +field:word +:\"exact phrase\" -:excludedWord",
		Query: Query{
			Required: []SubQuery{
				SubQuery{
					Quote:    QuoteNone,
					Operator: OperatorField,
					Value:    "mandatoryWord",
				},
				SubQuery{
					Quote:    QuoteNone,
					Operator: OperatorField,
					Field:    "field",
					Value:    "word",
				},
				SubQuery{
					Quote:    QuoteDouble,
					Operator: OperatorField,
					Value:    "exact phrase",
				},
			},
			Excluded: []SubQuery{
				SubQuery{
					Quote:    QuoteNone,
					Operator: OperatorField,
					Value:    "excludedWord",
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
					Value:    "Red Hat",
				},
				SubQuery{
					Quote:    QuoteNone,
					Operator: OperatorField,
					Value:    "Google",
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
					Value:    "Google",
				},
			},
			Excluded: []SubQuery{
				SubQuery{
					Quote:    QuoteDouble,
					Operator: OperatorField,
					Value:    "Red Hat",
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
					Value:    "Red Hat",
				},
				SubQuery{
					Quote:    QuoteDouble,
					Operator: OperatorField,
					Value:    "Fusion IO",
				},
			},
		},
	},
	{
		Input:  "(\"Cloud Computing\" AND \"Red Hat\") (\"Cloud Computing\" AND \"Fusion IO\")",
		String: "+(+:\"Cloud Computing\" +:\"Red Hat\") +(+:\"Cloud Computing\" +:\"Fusion IO\")",
		Query: Query{
			Required: []SubQuery{
				SubQuery{
					Quote:    QuoteNone,
					Operator: OperatorSubquery,
					Query: &Query{
						Required: []SubQuery{
							SubQuery{
								Quote:    QuoteDouble,
								Operator: OperatorField,
								Value:    "Cloud Computing",
							},
							SubQuery{
								Quote:    QuoteDouble,
								Operator: OperatorField,
								Value:    "Red Hat",
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
								Value:    "Cloud Computing",
							},
							SubQuery{
								Quote:    QuoteDouble,
								Operator: OperatorField,
								Value:    "Fusion IO",
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
					Value:    "Cloud Computing",
				},
				SubQuery{
					Quote:    QuoteNone,
					Operator: OperatorSubquery,
					Query: &Query{
						Optional: []SubQuery{
							SubQuery{
								Quote:    QuoteDouble,
								Operator: OperatorField,
								Value:    "Red Hat",
							},
							SubQuery{
								Quote:    QuoteDouble,
								Operator: OperatorField,
								Value:    "Fusion IO",
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
					Value:    "Colon:In the Tech",
				},
				SubQuery{
					Quote:    QuoteDouble,
					Operator: OperatorField,
					Value:    "Red Hat",
				},
			},
		},
	},
}
