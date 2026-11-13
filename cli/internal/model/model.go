// Package model holds what the language's files mean once their names are
// resolved: kinds, links, value types, statements and facts. It is inert data;
// loading, checking and running live in other packages.
package model

import (
	"math/big"

	"oflang/internal/diag"
)

// QName names a definition by its package and its name, so names from
// different packages never collide once imports are resolved.
type QName struct {
	Pkg  string `json:"package"`
	Name string `json:"name"`
}

// String renders a QName for messages.
func (q QName) String() string { return q.Name }

// Sort says what kind of entry a definition is: the value of its is: line.
type Sort string

// The sorts of entry.
const (
	SortKind       Sort = "Kind"
	SortLink       Sort = "Link"
	SortValueType  Sort = "ValueType"
	SortConstraint Sort = "Constraint"
	SortRule       Sort = "Rule"
	SortMapping    Sort = "Mapping"
	SortFact       Sort = "Fact"
	SortTest       Sort = "Test"
	SortQuestion   Sort = "Question"
)

// DefinitionSorts are the sorts that may appear outside tests.
var DefinitionSorts = []Sort{SortKind, SortLink, SortValueType, SortConstraint, SortRule, SortMapping}

// TestSorts are the sorts that may appear in a package's tests.
var TestSorts = []Sort{SortFact, SortTest, SortQuestion}

// Ref is what a name in a file's scope stands for. A link's inverse name is a
// Ref to the link with Inverted set.
type Ref struct {
	Sort     Sort
	Q        QName
	Inverted bool
}

// ValueKind is the JSON Schema type of a literal value.
type ValueKind string

// The kinds of value.
const (
	ValueString  ValueKind = "string"
	ValueInteger ValueKind = "integer"
	ValueNumber  ValueKind = "number"
	ValueBoolean ValueKind = "boolean"
	ValueArray   ValueKind = "array"
	ValueObject  ValueKind = "object"
)

// Value is a literal field value. Numbers are exact.
type Value struct {
	Kind  ValueKind        `json:"kind"`
	Text  string           `json:"text,omitzero"`
	Num   *big.Rat         `json:"-"`
	Bool  bool             `json:"bool,omitzero"`
	Items []Value          `json:"items,omitzero"`
	Props map[string]Value `json:"props,omitzero"`
}

// Schema describes a field's value with a subset of JSON Schema.
type Schema struct {
	Pos        diag.Pos
	Type       string // a built-in type, or "" when Ref names a value type
	Ref        QName  // a value type, when set
	HasRef     bool
	Enum       []Value
	Const      *Value
	Default    *Value
	Pattern    string
	MinLength  *int
	MaxLength  *int
	Minimum    *big.Rat
	Maximum    *big.Rat
	Format     string
	Items      *Schema
	Properties []Field
}

// Field is one field of a kind, or one part of an object value.
type Field struct {
	Name     string
	Pos      diag.Pos
	Schema   Schema
	Optional bool
}

// Count says how many of a link a thing must have.
type Count struct {
	Min int
	Max int // -1 when there is no upper limit
}

// LinkCount is one entry of a kind's links block.
type LinkCount struct {
	Ref   Ref
	Name  string
	Pos   diag.Pos
	Count Count
}

// Kind is a kind of thing.
type Kind struct {
	Q           QName
	Pos         diag.Pos
	Definition  string
	Extends     *QName
	Examples    []string
	NotExamples []string
	Fields      []Field
	Links       []LinkCount
}

// Link is a named, directed relationship between two kinds.
type Link struct {
	Q          QName
	Pos        diag.Pos
	From       QName
	To         QName
	Reads      string
	Definition string
	Under      *QName
	Inverse    string
	Chains     bool
}

// ValueType is a named, reusable schema.
type ValueType struct {
	Q          QName
	Pos        diag.Pos
	Definition string
	Schema     Schema
}

// PatternLink is one link a pattern must find. Link is always the declared
// link; Inverted is set when the pattern names it by its inverse.
type PatternLink struct {
	Pos      diag.Pos
	From     string
	To       string
	Name     string
	Link     QName
	Inverted bool
	Repeated bool
}

// Operand is what a comparison compares against: a literal value, or a
// variable and possibly one of its fields.
type Operand struct {
	Value  *Value
	Values []Value
	Node   string
	Field  string
}

// CondType says which kind of condition a Condition is.
type CondType string

// The kinds of condition.
const (
	CondCompare    CondType = "compare"
	CondIsA        CondType = "is a"
	CondIsSet      CondType = "is set"
	CondIsNotSet   CondType = "is not set"
	CondAtLeastOne CondType = "at least one"
	CondNone       CondType = "none"
	CondCount      CondType = "count"
	CondAnyOf      CondType = "any of"
)

// Op is a comparison word.
type Op string

// The comparison words.
const (
	OpEquals       Op = "equals"
	OpDoesNotEqual Op = "does not equal"
	OpIsAbove      Op = "is above"
	OpIsBelow      Op = "is below"
	OpIsAtLeast    Op = "is at least"
	OpIsAtMost     Op = "is at most"
	OpIsOneOf      Op = "is one of"
)

// Ops lists every comparison word, in the order messages list them.
var Ops = []Op{OpEquals, OpDoesNotEqual, OpIsAbove, OpIsBelow, OpIsAtLeast, OpIsAtMost, OpIsOneOf}

// Condition is one entry of a where or require list.
type Condition struct {
	Type  CondType
	Pos   diag.Pos
	Node  string
	Field string
	Op    Op
	Than  Operand
	Kind  QName
	// For at least one, none and count.
	Links []PatternLink
	Where []Condition
	// For any of.
	Any []Condition
}

// DeriveType says what a derivation adds.
type DeriveType string

// The kinds of derivation.
const (
	DeriveLink  DeriveType = "link"
	DeriveKind  DeriveType = "kind"
	DeriveField DeriveType = "field"
)

// Derivation is one entry of a rule's derive list.
type Derivation struct {
	Type  DeriveType
	Pos   diag.Pos
	From  string
	To    string
	Name  string
	Link  QName
	Inv   bool
	Node  string
	Kind  QName
	Field string
	SetTo Operand
}

// Statement is a constraint, a rule or a question: a pattern with what to do
// with each match.
type Statement struct {
	Q       QName
	Sort    Sort
	Pos     diag.Pos
	Says    string
	Vars    []string
	Links   []PatternLink
	Where   []Condition
	Require []Condition
	Message string
	Derive  []Derivation
}

// Mapping relates a definition to a term in an outside vocabulary.
type Mapping struct {
	Q          QName
	Pos        diag.Pos
	Term       Ref
	Equivalent string
}

// Exports lists what a package lets others import.
type Exports struct {
	Kinds      []string
	Links      []string
	ValueTypes []string
}

// Manifest is a package's package.yaml.
type Manifest struct {
	Pos       diag.Pos
	Package   string
	Version   string
	Namespace string
	DependsOn map[string]string
	Exports   Exports
	OpenFor   []OpenPoint
}

// OpenPoint is an extension point a package opens to others.
type OpenPoint struct {
	Extend string
	Under  string
}

// Definition records where a name is defined in a package.
type Definition struct {
	Sort Sort
	Pos  diag.Pos
	File string
	// For a link's inverse name, the link it is the inverse of.
	InverseOf string
}

// Package is one loaded package.
type Package struct {
	Name     string
	Dir      string
	Manifest Manifest
	Files    []string
	Docs     []string
	Defs     map[string]Definition
}

// World is every loaded package with its definitions resolved.
type World struct {
	Root       string
	Packages   map[string]Package
	Kinds      map[QName]Kind
	Links      map[QName]Link
	ValueTypes map[QName]ValueType
	Statements map[QName]Statement
	Mappings   map[QName]Mapping
	// Order is every statement in the order it was read, for stable output.
	Order []QName
	// Inverses maps a link's inverse name to the link.
	Inverses map[QName]QName
}
