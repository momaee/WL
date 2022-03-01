package parser

import (
	"github.com/momaee/WL/lexer"
	"github.com/momaee/WL/stack"
	"github.com/momaee/WL/token"
)

// RuneParser will parse tokens and pack them in instructions
// initial state of the RuneParser is Parse method
type RuneParser interface {
	Parse() []*Inst
}

// Inst is an abstraction for an operation which machine can understand
// T is one single instruction
// C is complementary information about instruction like position or counts of occurrence
// Incase of opening loop, C is the index of the closing loop and vice versa
type Inst struct {
	T *token.Token
	C int
}

// parser builds AST (abstract structure tree).
// parser uses Stack to keep track of loops
// it contains the Scanner to tokenize the data from input
// buf is an internal struct to process input at a time of scan
// inst is an slice, which every member is one single instruction
type parser struct {
	l    lexer.LexScanner
	inst []*Inst
	buf  struct {
		tok     *token.Token // last read token
		tokbufn bool         // whether the token buffer is in use.
	}
	stack stack.Stack
}

// NewParser creates new parser using given LexScanner.
func NewParser(l lexer.LexScanner) RuneParser {
	return &parser{l: l}
}

func (p *parser) Parse() []*Inst {
	for {
		tok := p.scan()
		if tok.Tok == token.IllegalToken {
			break
		}

		for _, t := range token.AllTokens {
			if tok == t {
				if tok.Tok == token.LeftBracketToken {
					openLoop := p.buildInst(tok, 0)
					p.stack.Push(openLoop)
				} else if tok.Tok == token.RightBracketToken {
					openLoop := p.stack.Pop().(int)
					closeLoop := p.buildInst(tok, openLoop)
					p.inst[openLoop].C = closeLoop
				} else {
					p.addInst(tok)
				}
			}
		}
	}
	return p.inst
}

// scan returns next token unit.
func (p *parser) scan() *token.Token {
	// there is a token on the buffer
	if p.buf.tokbufn {
		p.buf.tokbufn = false
		return p.buf.tok
	}
	// read the next token from s
	tok := p.l.Scan()
	p.buf.tok = tok
	return tok
}

// unscan sends the already consumed token back to buff.
func (p *parser) unscan() {
	p.buf.tokbufn = true
}

// addInst adds instructions to []*inst of parser
// for efficiency, if there are multiple occurrences of the
// same token consecutively, we will fold it.
func (p *parser) addInst(t *token.Token) int {
	// token occurrence count
	c := 1
	for {
		next := p.scan()
		if next.Tok != t.Tok {
			p.unscan()
			break
		}
		c++
	}
	return p.buildInst(t, c)
}

// buildInst creates a instruction from the given literals.
func (p *parser) buildInst(t *token.Token, c int) int {
	// build instruction
	inst := &Inst{
		T: t,
		C: c,
	}
	// add inst to instruction list
	p.inst = append(p.inst, inst)
	return len(p.inst) - 1
}
