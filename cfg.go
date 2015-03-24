package parser

// A grammar describes a CFG
type CFG struct {
	Name    string
	Start   string
	Rules   map[string][]Expr
	firsts  map[Sym]map[Sym]struct{}
	follows map[Sym]map[Sym]struct{}
}

func (cfg *CFG) AddRule(name string, expr Expr) {
	if cfg.Rules == nil {
		cfg.Rules = make(map[string][]Expr)
	}
	if rule, exists := cfg.Rules[name]; exists {
		cfg.Rules[name] = append(rule, expr)
	} else {
		cfg.Rules[name] = []Expr{expr}
	}
}

func (cfg *CFG) Firsts() map[Sym]map[Sym]struct{} {
	firsts := map[Sym]map[Sym]struct{}{}
	for name := range cfg.Rules {
		sym := Sym{NonTerm, name}
		firsts[sym] = cfg.first(sym)
	}
	return firsts
}

func (cfg *CFG) First(sym Sym) map[Sym]struct{} {
	return cfg.first(sym)
}

func (cfg *CFG) first(sym Sym) map[Sym]struct{} {
	if cfg.firsts == nil {
		cfg.firsts = make(map[Sym]map[Sym]struct{})
	}
	if set, exists := cfg.firsts[sym]; exists {
		return set
	}
	switch sym.Type {
	case NonTerm:
		set := map[Sym]struct{}{}
		containE := true
		for _, alt := range cfg.Rules[sym.Value] {
			if len(alt) == 0 {
				set[Sym{Type: Empty}] = struct{}{}
			}
			for _, expSym := range alt {
				if expSym.Value == sym.Value {
					break
				}
				expSet := cfg.first(expSym)
				if _, ok := expSet[Sym{Type: Empty}]; !ok {
					containE = false
				} else {
					delete(expSet, Sym{Type: Empty})
				}
				for s, v := range expSet {
					set[s] = v
				}
				break
			}
		}
		if containE {
			set[Sym{Type: Empty}] = struct{}{}
		}
		cfg.firsts[sym] = set
		return set
	default:
		return map[Sym]struct{}{sym: struct{}{}}
	}
}

func (c *CFG) Follow(sym Sym) map[Sym]struct{} {
	return c.follow()[sym]
}

func (c *CFG) Follows() map[Sym]map[Sym]struct{} {
	return c.follow()
}

func (c *CFG) follow() map[Sym]map[Sym]struct{} {
	if c.follows == nil {
		c.follows = make(map[Sym]map[Sym]struct{})
	} else {
		return c.follows
	}
	c.follows[Sym{NonTerm, c.Start}] = map[Sym]struct{}{Sym{End, ""}: struct{}{}}
	for {
		added := 0
		for name, exprs := range c.Rules {
			for _, expr := range exprs {
				for i, sym := range expr {
					var ok bool
					var followB map[Sym]struct{}
					prev := 0
					if followB, ok = c.follows[sym]; ok {
						prev = len(followB)
					} else {
						followB = map[Sym]struct{}{}
					}
					if i == len(expr)-1 {
						if followA, ok := c.follows[Sym{NonTerm, name}]; ok {
							for s, v := range followA {
								followB[s] = v
							}
						}
					} else {
						next := expr[i+1]
						for s, v := range c.first(next) {
							if s.Type != Empty {
								followB[s] = v
							}
						}
					}
					c.follows[sym] = followB
					added += len(followB) - prev
				}
			}
		}
		if added == 0 {
			break
		}
	}
	return c.follows
}

type Expr []Sym

func (e Expr) Equal(o Expr) bool {
	if len(e) != len(o) {
		return false
	}
	for i := range e {
		if e[i] != o[i] {
			return false
		}
	}
	return true
}

type SymType byte

const (
	Term SymType = iota
	NonTerm
	End
	Empty
)

type Sym struct {
	Type  SymType
	Value string
}

func (g *CFG) Validate() bool {
	for _, exprs := range g.Rules {
		for _, expr := range exprs {
			for _, sym := range expr {
				if sym.Type == NonTerm {
					if _, found := g.Rules[sym.Value]; !found {
						println(sym.Value + " not found")
						return false
					}
				}
			}
		}
	}
	return true
}

func (g *CFG) Augment() {
	inRight := false
outer:
	for _, exprs := range g.Rules {
		for _, expr := range exprs {
			for _, sym := range expr {
				if g.Start == sym.Value {
					inRight = true
					break outer
				}
			}
		}
	}
	if inRight {
		g.AddRule("", Expr{Sym{NonTerm, g.Start}})
		g.Start = ""
	}
}

type Item struct {
	Name string
	Expr Expr
	Dot  int
	LA   Sym
}

func (it Item) Advance() Item {
	if it.Dot < len(it.Expr) {
		it.Dot++
	}
	return it
}

func (it Item) Equal(ot Item) bool {
	return (it.Name == ot.Name) && it.Expr.Equal(ot.Expr) && (it.Dot == ot.Dot) && (it.LA == ot.LA)
}

func getTrans(its []Item) map[Sym][]Item {
	tSet := make(map[Sym][]Item)
	// build the set of transitions
	for _, v := range its {
		if v.Dot < len(v.Expr) {
			sym := v.Expr[v.Dot]
			if set, ok := tSet[sym]; !ok {
				tSet[sym] = []Item{v}
			} else {
				tSet[sym] = append(set, v)
			}
		}
	}
	return tSet
}

func getReduct(its []Item) map[Sym]Action {
	ret := make(map[Sym]Action)
	for _, v := range its {
		if v.Dot == len(v.Expr) {
			ret[v.LA] = Action{
				Type: Reduce,
				Num:  v.Dot,
				New:  v.Name,
			}
		}
	}

	return ret
}

func inClosure(it Item, clos []Item) bool {
	for _, v := range clos {
		if v.Equal(it) {
			return true
		}
	}
	return false
}

func getClosure(g *CFG, its []Item) []Item {
	newItems := its
	closure := make([]Item, 0)
	var tmpItems []Item
	for {
		tmpItems = make([]Item, 0)
		for _, it := range newItems {
			if it.Dot < len(it.Expr) {
				sym := it.Expr[it.Dot]
				for la := range g.Follow(sym) {
					if sym.Type == NonTerm {
						for name, exprs := range g.Rules {
							for _, expr := range exprs {
								if sym.Value == name {
									newIt := Item{
										Name: name,
										Expr: expr,
										LA:   la,
										Dot:  0,
									}
									if !(inClosure(newIt, closure) ||
										inClosure(newIt, newItems) ||
										inClosure(newIt, tmpItems)) {
										tmpItems = append(tmpItems, newIt)
									}
								}
							}
						}
					}
				}
			}
		}
		closure = append(closure, newItems...)

		if len(tmpItems) == 0 {
			return closure
		}
		newItems = tmpItems
	}
}

func (c *CFG) getStart() []Item {
	exprs, ok := c.Rules[c.Start]
	if !ok {
		return nil
	}
	its := make([]Item, 0, len(exprs))
	for _, v := range exprs {
		its = append(its, Item{
			Name: c.Start,
			Expr: v,
			Dot:  0,
			LA:   Sym{Type: End},
		})
	}
	return its
}

func (cfg *CFG) MakeTable() LRTable {
	table := make(LRTable, 0)
	cfg.Augment()

	start := cfg.getStart()
	table.makeState(cfg, start)
	return table
}

func (table *LRTable) makeState(cfg *CFG, into []Item) {
	current := len(*table)
	items := getClosure(cfg, into)
	*table = append(*table, getReduct(items))
	trans := getTrans(items)
	// log.Debug("Current State: %d", current)
	// log.Debug("  Incoming: %v", into)
	// log.Debug("  Closure: %v", items)
	// log.Debug("  Trans on: ")
	// for _, v := range trans {
	// 	log.Debug("    %v,", v[0].Expr[v[0].Dot])
	// }
	// log.Debug("  Reduce on: ")
	// for k, v := range getReduct(items) {
	// 	log.Debug("    %v -> %s, ", k, v.New)
	// }
	for _, v := range trans {
		offset := len(*table)
		actType := Goto
		sym := v[0].Expr[v[0].Dot]
		if sym.Type == Term {
			actType = Shift
		}
		// if e, conflict := (*table)[current][sym]; conflict {
		// 	fmt.Println("Conflict: ", current, sym, e.Type, actType)
		// }
		(*table)[current][sym] = Action{
			Type: actType,
			Num:  offset,
		}
		for i, w := range v {
			v[i] = w.Advance()
		}
		table.makeState(cfg, v)
	}
	for k, v := range (*table)[current] {
		if v.Type == Reduce && v.New == cfg.Start {
			newAct := (*table)[current][k]
			newAct.Type = Accept
			(*table)[current][k] = newAct
		}
	}
}
