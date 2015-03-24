package parser

import (
	"io/ioutil"
	"log"
	"regexp"
	"strings"
)

var (
	rule   = regexp.MustCompile(`^(?P<name>[^ ]*) -> (?P<def>.*);$`)
	strlit = regexp.MustCompile(`"(?P<val>(\"|[^"])*)"`)
	repl   = strings.NewReplacer("\\\\", "\\", "\\\"", "\"", "\\n", "\n", "\\t", "\t")
)

func parseCFG(file, start string) *CFG {
	g := new(CFG)
	g.Name = "CFG"
	bs, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(string(bs), "\n")
	for _, l := range lines {
		if len(l) == 0 {
			continue
		}
		name := rule.ReplaceAllString(l, "${name}")
		expr := make([]Sym, 0)
		def := rule.ReplaceAllString(l, "${def}")
		syms := strings.Split(def, " ")
		for _, sym := range syms {
			if len(sym) == 0 {
				continue
			}
			var newSym Sym
			if strlit.MatchString(sym) {
				newSym.Type = Term
				newSym.Value = repl.Replace(strlit.ReplaceAllString(sym, "${val}"))
			} else {
				newSym.Type = NonTerm
				newSym.Value = sym
			}

			expr = append(expr, newSym)

		}

		g.AddRule(name, expr)
	}

	g.Start = start
	g.Augment()
	return g

}
