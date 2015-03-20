package parser

import (
	"encoding/json"
	"os"
	"testing"
)

var (
	TestRule = Rule{
		Name: "test",
		Alts: []Expr{
			[]Sym{
				Sym{
					Type:  Term,
					Value: "a",
				},
				Sym{
					Type:  NonTerm,
					Value: "test",
				},
			},
			[]Sym{
				Sym{
					Type:  Term,
					Value: "a",
				},
			},
		},
	}
)

func TestItems(t *testing.T) {
	items := GetItems(TestRule)

	if len(items) != 5 {
		t.Fatalf("Wrong number of items, expecting 5, got %d", len(items))
	}
}

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
}
