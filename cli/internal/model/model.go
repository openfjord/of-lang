// SPDX-License-Identifier: Apache-2.0

// Package model holds what the language's files mean once their names are
// resolved: kinds, links, value types, constraints, derivatives and questions.
// It is inert data; loading, checking and running live in other packages.
package model

import (
	"math/big"

	"github.com/openfjord/of-lang/cli/internal/diag"
)

// QName names a definition by its package and its name, so names from
// different packages never collide once imports are resolved.
type QName struct {
	Pkg  string `json:"package"`
	Name string `json:"name"`
}

// String renders a QName for messages.
func (q QName) String() string { return q.Name }

// Sort says what a statement is: the value of its is: key.
type Sort string

// The sorts of statement.
const (
	SortKind       Sort = "Kind"
	SortLink       Sort = "Link"
	SortValueType  Sort = "ValueType"
	SortConstraint Sort = "Constraint"
	SortDerivative Sort = "Derivative"
	SortDoc        Sort = "Doc"
	SortFact       Sort = "Fact"
	SortTest       Sort = "Test"
	SortQuestion   Sort = "Question"
)

// DefinitionSorts are the sorts that may appear in defs/.
var DefinitionSorts = []Sort{SortKind, SortLink, SortValueType, SortConstraint, SortDerivative, SortDoc}

// TestSorts are the sorts that may appear in tests/.
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
	Ref        QName  // a value type, when HasRef is set
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
	OneOf      []Schema
	// Keywords lists every keyword the schema wrote, so a check can tell
	// what a field adds to the value type it names.
	Keywords []string
}

// Field is one field of a kind, or one part of an object value.
type Field struct {
	Name     string
	Pos      diag.Pos
	Schema   Schema
	Optional bool
}

// Kind is a kind of thing. A kind with OneOf is a union of the listed kinds.
type Kind struct {
	Q           QName
	Pos         diag.Pos
	Definition  string
	Extends     *QName
	OneOf       []QName
	Examples    []string
	NotExamples []string
	Fields      []Field
	Tags        []Tag
	// Docs are the docs of the package that tell more about the kind.
	Docs []QName
}

// Doc names a Markdown file in the docs/ folder of its package.
type Doc struct {
	Q   QName
	Pos diag.Pos
	// Path is the absolute path of the Markdown file.
	Path string
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
	Tags       []Tag
}

// ValueType is a named, reusable schema.
type ValueType struct {
	Q          QName
	Pos        diag.Pos
	Definition string
	Schema     Schema
	Tags       []Tag
}

// Tag is one tag a definition carries, where it is written.
type Tag struct {
	Name string
	Pos  diag.Pos
}

// KindItem gives a variable to the things of a kind.
type KindItem struct {
	Pos  diag.Pos
	Var  string
	Kind QName
}

// End is one end of a link item: a variable, or a kind alone for a thing that
// nothing else refers to.
type End struct {
	Pos  diag.Pos
	Var  string
	Kind QName // set when Var is empty
}

// LinkItem is one link a match must find. Link is always the declared link;
// Inverted is set when the item names it by its inverse.
type LinkItem struct {
	Pos      diag.Pos
	From     End
	To       End
	Name     string
	Link     QName
	Inverted bool
	Repeated bool
}

// Block is a match and its filter: the statement's own, or a nested block's.
type Block struct {
	Pos    diag.Pos
	Kinds  []KindItem
	Links  []LinkItem
	Filter []Condition
}

// Operand is what a comparison compares with: a literal value, a list of
// them, or a variable and possibly one of its fields.
type Operand struct {
	Pos    diag.Pos
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

// Condition is one entry of a filter or an unless list.
type Condition struct {
	Type  CondType
	Pos   diag.Pos
	Node  string
	Field string
	Part  []string
	Op    Op
	Than  Operand
	Kind  QName
	// Block is the nested block of at least one, none and count.
	Block *Block
	// Any holds the conditions of any of.
	Any []Condition
}

// DeriveType says what a derive item adds.
type DeriveType string

// The kinds of derive item.
const (
	DeriveLink  DeriveType = "link"
	DeriveKind  DeriveType = "kind"
	DeriveField DeriveType = "field"
)

// Derivation is one item of a derivative's derive list.
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

// Statement is a constraint, a derivative or a question: a match with what to
// do with each row.
type Statement struct {
	Q       QName
	Sort    Sort
	Pos     diag.Pos
	Says    string
	Tags    []Tag
	Match   Block
	Unless  []Condition
	Message string
	Derive  []Derivation
}

// Exports lists what a package lets others import.
type Exports struct {
	Kinds      []string
	Links      []string
	ValueTypes []string
}

// OpenPoint is an extension point a package opens to others.
type OpenPoint struct {
	Pos    diag.Pos
	Extend string
	Under  string
}

// TagDecl is one tag a manifest declares.
type TagDecl struct {
	Pos         diag.Pos
	Name        string
	Description string
}

// Manifest is a package's package.yaml.
type Manifest struct {
	Pos          diag.Pos
	Package      string
	Version      string
	Owner        string
	RDFNamespace string
	DependsOn    map[string]string
	Exports      Exports
	OpenFor      []OpenPoint
	Tags         []TagDecl
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
	Tests    []string
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
	Docs       map[QName]Doc
	Statements map[QName]Statement
	// Order is every constraint and derivative in the order it was read, for
	// stable output.
	Order []QName
	// Inverses maps a link's inverse name to the link.
	Inverses map[QName]QName
}
