package parser

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/Pursuit92/log"
)

var (
	TestItems = []Item{
		Item{
			Name: "S",
			Expr: []Sym{Sym{NonTerm, "A"}},
			Dot:  0,
			LA: Sym{
				Type: End,
			},
		},
		Item{
			Name: "A",
			Expr: []Sym{
				Sym{NonTerm, "A"},
				Sym{Term, "+"},
				Sym{NonTerm, "B"},
			},
			Dot: 0,
			LA: Sym{
				Type: End,
			},
		},
		Item{
			Name: "A",
			Expr: []Sym{Sym{Term, "a"}},
			Dot:  0,
			LA: Sym{
				Type: End,
			},
		},
		Item{
			Name: "A",
			Expr: []Sym{
				Sym{NonTerm, "A"},
				Sym{Term, "+"},
				Sym{NonTerm, "B"},
			},
			Dot: 0,
			LA:  Sym{Term, "+"},
		},
		Item{
			Name: "A",
			Expr: []Sym{Sym{Term, "a"}},
			Dot:  0,
			LA:   Sym{Term, "+"},
		},
	}
)

// func TestItems(t *testing.T) {
// 	items := GetItems(TestRule)

// 	if len(items) != 5 {
// 		t.Fatalf("Wrong number of items, expecting 5, got %d", len(items))
// 	}
// }

func TestLoad(t *testing.T) {
	file, err := os.Open("grammar.json")
	if err != nil {
		t.Fatal(err)
	}
	grammar := new(CFG)
	err = json.NewDecoder(file).Decode(grammar)
	if err != nil {
		t.Fatal(err)
	}
	if !grammar.Validate() {
		t.Fatal("failed to validate")
	}
}

func TestTrans(t *testing.T) {
	//var TestTrans [][]Item
	_ = getTrans(TestItems)
	//spew.Dump(trans)
	// if reflect.DeepEqual(trans, TestTrans) != true {
	// 	t.Fatal("no match")
	// }
}

var (
	TestCFG = CFG{
		Name:  "Test",
		Start: "S",
		Rules: map[string][]Expr{
			"S": []Expr{
				Expr{
					Sym{
						Type:  NonTerm,
						Value: "A",
					},
				},
			},
			"A": []Expr{
				Expr{
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
				Expr{
					Sym{
						Type:  Term,
						Value: "a",
					},
				},
			},
			"B": []Expr{
				Expr{
					Sym{
						Type:  Term,
						Value: "b",
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
	input := []byte("a+ b\n+b")
	ast, err := TestTable.Parse(&scanner{0, true, input})
	if err != nil {
		t.Fatal(err)
	}
	log.Normal(ast.SEXP())
}

// func TestClosure(t *testing.T) {
// 	clos := getClosure(&TestCFG, (&TestCFG).getStart())
// 	trans := getTrans(clos)
// 	state1 := trans[0]
// 	for i, v := range state1 {
// 		state1[i] = v.Advance()
// 	}
// 	_ = getClosure(&TestCFG, state1)
// 	//spew.Dump(getTrans(clos))
// }

func TestMakeTable(t *testing.T) {
	log.SetLevel(log.LogDebug)
	tab := (&TestCFG).MakeTable()
	log.SetLevel(log.LogNormal)
	input := []byte("a+b+b")
	ast, err := tab.Parse(&scanner{0, true, input})
	if err != nil {
		t.Fatal(err)
	}
	log.Normal(ast.SEXP())
	//spew.Dump(tab)
}

func BenchmarkGrammar(b *testing.B) {
	b.StopTimer()
	cfg := parseCFG("grammar.cfg", "Grammar")
	if !cfg.Validate() {
		b.Fatal("validate fail")
	}

	log.Normal("building table")
	b.StartTimer()
	tab := cfg.MakeTable()
	b.StopTimer()
	log.Normal("done, %d states", len(tab))
	//spew.Dump(tab)
	b.StartTimer()

}

// moment of truth
func TestParseGrammar(t *testing.T) {
	log.SetLevel(log.LogNormal)
	cfg := parseCFG("grammar.cfg", "Grammar")
	if !cfg.Validate() {
		t.Fatal("validate fail")
	}

	log.Normal("building table")
	tab := cfg.MakeTable()
	log.Normal("done, %d states", len(tab))
	//spew.Dump(tab)

	gbs := []byte("this_is_a_sentence") //, err := ioutil.ReadFile("grammar.cfg")
	// if err != nil {
	// 	t.Fatal(err)
	// }

	ast, err := tab.Parse(&scanner{0, false, gbs})
	if err != nil {
		//spew.Dump(tab[0])
		t.Fatal(err)
	}

	fmt.Println(ast.SEXP())
}

func TestFirsts(t *testing.T) {
	firsts := (&TestCFG).Firsts()
	fmt.Println("Firsts:")
	for s, f := range firsts {
		fmt.Printf("  %s:\n", s.Value)
		for s, _ := range f {
			fmt.Printf("    %s\n", s.Value)
		}
	}
}
func TestFollows(t *testing.T) {
	follows := (&TestCFG).Follows()
	fmt.Println("Follows:")
	for s, f := range follows {
		fmt.Printf("  %s:\n", s.Value)
		for s, _ := range f {
			fmt.Printf("    %s\n", s.Value)
		}
	}
}

func BenchmarkFollows(b *testing.B) {

	b.StopTimer()
	log.SetLevel(log.LogDebug)
	cfg := parseCFG("grammar.cfg", "Grammar")
	if !cfg.Validate() {
		b.Fatal("validate fail")
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		cfg.Follows()
	}
}

func BenchmarkFirsts(b *testing.B) {

	b.StopTimer()
	log.SetLevel(log.LogDebug)
	cfg := parseCFG("grammar.cfg", "Grammar")
	if !cfg.Validate() {
		b.Fatal("validate fail")
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		cfg.Firsts()
	}
}
