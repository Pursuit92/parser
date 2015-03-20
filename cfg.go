package parser

// A grammar describes a CFG
type CFG struct {
	Name  string
	Start string
	Rules []Rule
}

// A Rule is a rule of a CFG. It only has basic
// support for alternatives
type Rule struct {
	Name string
	Alts []Expr
}

type Expr []Sym

type SymType byte

const (
	Term SymType = iota
	NonTerm
	End
)

type Sym struct {
	Type  SymType
	Value string
}

// SplitAlts takes a rule and splits it into multiple rules based
// on its Alts slice. Each Rule in the slice will have one Alt.
// It also returns the rule name for convenience.
func (r Rule) SplitAlts() []Rule {
	rules := make([]Rule, len(r.Alts))
	for i, v := range r.Alts {
		rules[i].Name = r.Name
		rules[i].Alts = []Expr{v}
	}
	return rules
}

func (g *CFG) Flatten() {
	alts := make([][]Rule, len(g.Rules))
	numRules := 0
	for i, v := range g.Rules {
		alts[i] = v.SplitAlts()
		numRules += len(alts[i])
	}
	rules := make([]Rule, 0, numRules)
	for _, v := range alts {
		rules = append(rules, v...)
	}
	g.Rules = rules
}

func (g *CFG) Validate() bool {
	for _, rule := range g.Rules {
		for _, alt := range rule.Alts {
			for _, sym := range alt {
				if sym.Type == NonTerm {
					found := false
					for _, rule := range g.Rules {
						if sym.Value == rule.Name {
							found = true
							break
						}
					}
					if !found {
						return false
					}
				}
			}
		}
	}
	return true
}

type Item struct {
	Name       string
	Seen, Pred Expr
}

func GetItems(r Rule) []Item {
	itemLen := 0
	for _, v := range r.Alts {
		itemLen += len(v) + 1
	}
	items := make([]Item, 0, itemLen)
	for _, v := range r.Alts {
		for i := range v {
			items = append(items, Item{
				Name: r.Name,
				Seen: v[:i],
				Pred: v[i:],
			})
		}
		items = append(items, Item{
			Name: r.Name,
			Seen: v,
			Pred: Expr{},
		})
	}
	return items
}
