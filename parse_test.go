package parser

import (
	"testing"

	"github.com/Pursuit92/log"
)

var (
	TestCFG = CFG{
		Name:  "Test",
		Start: "S",
		Rules: []Rule{
			Rule{
				Name: "S",
				Alts: []Expr{
					[]Sym{
						Sym{
							Type:  NonTerm,
							Value: "A",
						},
					},
				},
			},
			Rule{
				Name: "A",
				Alts: []Expr{
					[]Sym{
						Sym{
							Type:  NonTerm,
							Value: "A",
						},
						Sym{
							Type:  Term,
							Value: "+",
						},
						Sym{
							Type:  NonTerm,
							Value: "B",
						},
					},
				},
			},
			Rule{
				Name: "A",
				Alts: []Expr{
					[]Sym{
						Sym{
							Type:  Term,
							Value: "a",
						},
					},
				},
			},
			Rule{
				Name: "B",
				Alts: []Expr{
					[]Sym{
						Sym{
							Type:  Term,
							Value: "b",
						},
					},
				},
			},
		},
	}
	TestTable LRTable = []map[Sym]Action{
		0: map[Sym]Action{
			Sym{
				Type:  Term,
				Value: "a",
			}: Action{
				Type: Shift,
				Num:  2,
			},
			Sym{
				Type:  NonTerm,
				Value: "A",
			}: Action{
				Type: Goto,
				Num:  1,
			},
		},
		1: map[Sym]Action{
			Sym{
				Type:  Term,
				Value: "+",
			}: Action{
				Type: Shift,
				Num:  3,
			},
			Sym{
				Type: End,
			}: Action{
				Type: Accept,
				Num:  1,
				New:  "S",
			},
		},
		2: map[Sym]Action{
			Sym{
				Type:  Term,
				Value: "+",
			}: Action{
				Type: Reduce,
				Num:  1,
				New:  "A",
			},
			Sym{
				Type: End,
			}: Action{
				Type: Reduce,
				Num:  1,
				New:  "A",
			},
		},
		3: map[Sym]Action{
			Sym{
				Type:  Term,
				Value: "b",
			}: Action{
				Type: Shift,
				Num:  5,
			},
			Sym{
				Type:  NonTerm,
				Value: "B",
			}: Action{
				Type: Goto,
				Num:  4,
			},
		},
		4: map[Sym]Action{
			Sym{
				Type:  Term,
				Value: "+",
			}: Action{
				Type: Reduce,
				Num:  3,
				New:  "A",
			},
			Sym{
				Type: End,
			}: Action{
				Type: Reduce,
				Num:  3,
				New:  "A",
			},
		},
		5: map[Sym]Action{
			Sym{
				Type:  Term,
				Value: "+",
			}: Action{
				Type: Reduce,
				Num:  1,
				New:  "B",
			},
			Sym{
				Type: End,
			}: Action{
				Type: Reduce,
				Num:  1,
				New:  "B",
			},
		},
	}
)

func TestValidate(t *testing.T) {
	if !TestCFG.Validate() {
		t.Fatal("grammar failed validation")
	}
}

func TestParse(t *testing.T) {
	log.SetLevel(log.LogNormal)
	input := []byte("a+b+b")
	ast, err := TestTable.Parse(&scanner{0, input})
	if err != nil {
		t.Fatal(err)
	}
	log.Normal(ast.SEXP())
}
