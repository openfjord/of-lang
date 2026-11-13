// Package engine runs a package's rules over facts until nothing new appears,
// then checks its constraints and link counts. Every fact it adds keeps the
// reason it was added, so any result can be explained.
package engine

import (
	"slices"
	"strings"

	"oflang/internal/diag"
	"oflang/internal/model"
	"oflang/internal/value"
)

// Edge is one link between two things, always stored in the link's declared
// direction; an inverse name only reads it the other way.
type Edge struct {
	From string      `json:"from"`
	Link model.QName `json:"link"`
	To   string      `json:"to"`
}

// FactType says what a fact states.
type FactType string

// The kinds of fact.
const (
	FactLink  FactType = "link"
	FactKind  FactType = "kind"
	FactField FactType = "field"
)

// Fact is one statement about things: a link between two, a kind one has, or
// the value of one of its fields.
type Fact struct {
	Type  FactType
	Edge  Edge
	Node  string
	Kind  model.QName
	Field string
	Value model.Value
}

// ReasonType says why a fact holds.
type ReasonType string

// The reasons a fact can hold.
const (
	ReasonAsserted ReasonType = "asserted"
	ReasonRule     ReasonType = "rule"
	ReasonUnder    ReasonType = "under"
	ReasonChains   ReasonType = "chains"
)

// Binding is the thing a variable stood for when a rule fired.
type Binding struct {
	Var string `json:"var"`
	ID  string `json:"id"`
}

// Reason records why a fact holds: written in a file, or added by a rule or a
// built-in inference from other facts.
type Reason struct {
	Type     ReasonType
	Pos      diag.Pos
	Rule     model.QName
	Bindings []Binding
	Premises []Fact
}

// Graph is a set of facts with the reason each one holds.
type Graph struct {
	Nodes   []string
	Pos     map[string]diag.Pos
	Kinds   map[string][]model.QName
	Fields  map[string]map[string]model.Value
	Edges   map[Edge]bool
	byLink  map[model.QName][]Edge
	out     map[string]map[model.QName][]string
	in      map[string]map[model.QName][]string
	Reasons map[string]Reason
	Order   []Fact
}

// NewGraph makes an empty graph.
func NewGraph() Graph {
	return Graph{
		Pos:     map[string]diag.Pos{},
		Kinds:   map[string][]model.QName{},
		Fields:  map[string]map[string]model.Value{},
		Edges:   map[Edge]bool{},
		byLink:  map[model.QName][]Edge{},
		out:     map[string]map[model.QName][]string{},
		in:      map[string]map[model.QName][]string{},
		Reasons: map[string]Reason{},
	}
}

// Key identifies a fact, so facts can be looked up and compared across runs.
func Key(f Fact) string {
	switch f.Type {
	case FactLink:
		return strings.Join([]string{"link", f.Edge.From, f.Edge.Link.Pkg, f.Edge.Link.Name, f.Edge.To}, "\x00")
	case FactKind:
		return strings.Join([]string{"kind", f.Node, f.Kind.Pkg, f.Kind.Name}, "\x00")
	default:
		return strings.Join([]string{"field", f.Node, f.Field}, "\x00")
	}
}

// AddNode records a thing, once.
func AddNode(g *Graph, id string, pos diag.Pos) {
	if _, ok := g.Pos[id]; ok {
		return
	}
	g.Nodes = append(g.Nodes, id)
	g.Pos[id] = pos
}

// HasNode says whether a thing is known.
func HasNode(g Graph, id string) bool {
	_, ok := g.Pos[id]
	return ok
}

// Add records a fact with its reason. It returns false when the fact already
// held, keeping the first reason.
func Add(g *Graph, f Fact, r Reason) bool {
	k := Key(f)
	if _, ok := g.Reasons[k]; ok {
		return false
	}
	switch f.Type {
	case FactLink:
		e := f.Edge
		g.Edges[e] = true
		g.byLink[e.Link] = append(g.byLink[e.Link], e)
		if g.out[e.From] == nil {
			g.out[e.From] = map[model.QName][]string{}
		}
		g.out[e.From][e.Link] = append(g.out[e.From][e.Link], e.To)
		if g.in[e.To] == nil {
			g.in[e.To] = map[model.QName][]string{}
		}
		g.in[e.To][e.Link] = append(g.in[e.To][e.Link], e.From)
	case FactKind:
		g.Kinds[f.Node] = append(g.Kinds[f.Node], f.Kind)
	case FactField:
		if g.Fields[f.Node] == nil {
			g.Fields[f.Node] = map[string]model.Value{}
		}
		g.Fields[f.Node][f.Field] = f.Value
	}
	g.Reasons[k] = r
	g.Order = append(g.Order, f)
	return true
}

// Has says whether a fact holds. A field fact holds when the field has that
// value.
func Has(g Graph, f Fact) bool {
	if f.Type == FactField {
		v, ok := g.Fields[f.Node][f.Field]
		return ok && value.Equal(v, f.Value)
	}
	_, ok := g.Reasons[Key(f)]
	return ok
}

// LinkFact makes a link fact.
func LinkFact(from string, link model.QName, to string) Fact {
	return Fact{Type: FactLink, Edge: Edge{From: from, Link: link, To: to}}
}

// KindFact makes a kind fact.
func KindFact(node string, kind model.QName) Fact {
	return Fact{Type: FactKind, Node: node, Kind: kind}
}

// FieldFact makes a field fact.
func FieldFact(node, field string, v model.Value) Fact {
	return Fact{Type: FactField, Node: node, Field: field, Value: v}
}

func sortedBindings(env map[string]string, vars []string) []Binding {
	var out []Binding
	for _, v := range vars {
		if id, ok := env[v]; ok {
			out = append(out, Binding{Var: v, ID: id})
		}
	}
	slices.SortFunc(out, func(a, b Binding) int { return strings.Compare(a.Var, b.Var) })
	return out
}
