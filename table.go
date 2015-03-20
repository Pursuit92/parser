package parser

// An LRTable is a list of states and a set of actions
// to take when a particular symbol is seen. The symbol
// can be a terminal or a non-terminal
type LRTable []map[Sym]Action

type ActType byte

const (
	Shift = iota
	Reduce
	Goto
	Accept
)

type Action struct {
	Type ActType
	Num  int
	New  string
}
