package parser

import (
	"bytes"
	"container/list"
	"errors"
	"fmt"

	"github.com/Pursuit92/log"
)

type AST struct {
	Leaf     bool
	Parent   *AST
	Children []*AST
	Name     string
}

func (a *AST) AddChild(n *AST) {
	n.Parent = a
	if a.Children == nil {
		a.Children = []*AST{n}
		a.Leaf = false
		return
	}
	a.Children = append([]*AST{n}, a.Children...)
}

type scanner struct {
	i  uint64
	bs []byte
}

func (s *scanner) next() Sym {
	if s.i == uint64(len(s.bs)) {
		return Sym{Type: End}
	} else {
		s.i++
		return Sym{Type: Term, Value: string([]byte{s.bs[s.i-1]})}
	}
}

func (pt LRTable) Parse(tokens *scanner) (*AST, error) {
	stack := new(list.List)
	stack.PushFront(0)
	tok := tokens.next()
	for {
		s := stack.Front().Value.(int) // current state
		log.Debug("State %d, token: %v\n", s, tok)
		action, ok := pt[s][tok]
		if !ok {
			goto parseErr
		}
		switch action.Type {
		case Shift:
			log.Debug("Action is shift %d\n", action.Num)
			stack.PushFront(tok)
			stack.PushFront(action.Num)
			tok = tokens.next()
		case Reduce, Accept:
			log.Debug("Action is reduce %d to %s\n", action.Num, action.New)
			newNode := &AST{
				Name: action.New,
				Leaf: false,
			}
			for i := 0; i < action.Num; i++ {
				stack.Remove(stack.Front())          // state
				front := stack.Remove(stack.Front()) // token
				var newChild *AST
				switch front := front.(type) {
				case *AST:
					newChild = front
				case Sym:
					newChild = &AST{
						Leaf: true,
						Name: front.Value,
					}
				}
				newNode.AddChild(newChild)
			}
			if action.Type == Accept {
				return newNode, nil
			}
			s := stack.Front().Value.(int)
			newSym := Sym{
				Type:  NonTerm,
				Value: action.New,
			}
			stack.PushFront(newNode)
			stack.PushFront(pt[s][newSym].Num)
		}
	}

parseErr:
	return nil, errors.New("Parse Error")
}

func (ast *AST) SEXP() string {
	buf := &bytes.Buffer{}
	open := "("
	close := ")"
	sep := " "
	if ast.Leaf {
		open, close = "\"", "\""
		sep = ""
	}
	fmt.Fprint(buf, open)
	fmt.Fprint(buf, ast.Name)
	fmt.Fprint(buf, sep)
	for i, v := range ast.Children {
		fmt.Fprint(buf, v.SEXP())
		if i < len(ast.Children)-1 {
			fmt.Fprint(buf, " ")
		}
	}
	fmt.Fprint(buf, close)
	return buf.String()
}
