// SPDX-License-Identifier: Apache-2.0

package model

// Rules names the rules of the specification that a constraint, a
// derivative and a question share. Each sort numbers its own copy of a rule,
// so a message names the copy that applies.
type Rules struct {
	Sentence    string // the sentence is plain words, with no $ word
	Items       string // each match item is a kind item or a link item
	Ends        string // each end is a variable or a kind
	EndKinds    string // the kind at each end fits the link
	Vars        string // conditions write a variable alone, with a kind item
	Scope       string // a variable of a nested block stays in that block
	Conds       string // each condition has one of the forms
	Types       string // a comparison compares values of one type
	Fields      string // a field is a field of its node's kind
	NodeKinds   string // two nodes compared without fields have related kinds
	Parts       string // part names a part of the field's schema
	Match       string // match has at least one item
	OneKindItem string // each variable has exactly one kind item
	Unused      string // each variable is used outside its kind item
	Count       string // count has one comparison word with a whole number
}

// RulesFor returns the rule names for a sort of statement.
func RulesFor(s Sort) Rules {
	switch s {
	case SortDerivative:
		return Rules{
			Sentence: "DV1", Items: "DV2", Ends: "DV3", EndKinds: "DV4", Vars: "DV5", Scope: "DV6",
			Conds: "DV7", Types: "DV8", Fields: "DV9", NodeKinds: "DV10", Parts: "DV18", Match: "DV19",
			OneKindItem: "DV20", Unused: "DV21", Count: "DV22",
		}
	case SortQuestion:
		return Rules{
			Sentence: "TS4", Items: "TS5", Ends: "TS6", EndKinds: "TS6", Vars: "TS7", Scope: "TS7",
			Conds: "TS16", Types: "TS8", Fields: "TS8", NodeKinds: "TS8", Parts: "TS11", Match: "TS12",
			OneKindItem: "TS13", Unused: "TS14", Count: "TS17",
		}
	default:
		return Rules{
			Sentence: "CN1", Items: "CN2", Ends: "CN3", EndKinds: "CN4", Vars: "CN5", Scope: "CN6",
			Conds: "CN7", Types: "CN8", Fields: "CN9", NodeKinds: "CN10", Parts: "CN15", Match: "CN16",
			OneKindItem: "CN17", Unused: "CN18", Count: "CN19",
		}
	}
}
