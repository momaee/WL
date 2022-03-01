package interpreter

import (
	"fmt"
	"io"

	"github.com/momaee/WL/lexer"
	"github.com/momaee/WL/parser"
	"github.com/momaee/WL/token"
)

// interface for an interpreter
// Run method executes created instructions by Parser
type Interpreter interface {
	Run() error
	AddOperator(symbol rune, operator Operator) error
	RemoveOperator(symbol rune) error
	GetValueInMemory(position int) int
}

// brainFuck is an implementation of the Interpreter
// it has internal parser which builds instructions from the input
// result is written into w
// memory struct keeps memory data and cursor to move between memory cells and update their values
// err != nil if any error happen during the print/read operation
type brainFuck struct {
	p      parser.RuneParser
	w      io.Writer
	i      io.Reader
	buf    []byte
	ip     int
	err    error
	memory Memory
}

type Memory = token.Memory

type Operator = token.Operator

func (b *brainFuck) execute(c int, op Operator) {
	op(c, &b.memory)
}

// NewInterpreter creates new Interpreter instance and initialize it's internal parser.
// i is used to read input from io
// w is used to write output to io
// code is used to read instructions from io
func NewInterpreter(i io.Reader, w io.Writer, code io.Reader) Interpreter {
	return &brainFuck{
		p:   parser.NewParser(lexer.NewScanner(code)),
		w:   w,
		i:   i,
		buf: make([]byte, 1),
	}
}

// Run method executes the instructions
// err != nil if error happen during read/print operations
// output returns in format of bytes
func (b *brainFuck) Run() error {
	inst := b.p.Parse()
	for b.ip < len(inst) {

		for _, t := range token.AllTokens {
			if inst[b.ip].T == t {
				c := inst[b.ip].C
				if t.HasOperator() {
					b.execute(c, t.Operator)
				} else {
					switch inst[b.ip].T.Tok {
					case token.PrintToken:
						b.execute(c, b.write())

					case token.ReadToken:
						b.execute(c, b.read())

					case token.LeftBracketToken:
						if b.val() == 0 {
							b.execute(c, b.jump())
							continue
						}

					case token.RightBracketToken:
						if b.val() != 0 {
							b.execute(c, b.jump())
							continue
						}

					default:
						b.err = fmt.Errorf("unknown token %v", inst[b.ip].T.Tok)
						return b.err
					}
				}
			}
		}

		b.ip++
	}

	return b.err
}

// cur method returns the position of current cursor in the memory
func (b *brainFuck) cur() int {
	return b.memory.Cursor
}

// jump method forwards the cursor to position p.
func (b *brainFuck) jump() Operator {
	return func(p int, memory *Memory) {
		b.ip = p
	}
}

// read reads input from io
// if any error happen during the Read operation err property will be set.
func (b *brainFuck) read() Operator {
	return func(times int, memory *Memory) {
		for i := 0; i < times; i++ {
			if _, err := b.i.Read(b.buf); err != nil {
				b.err = err
				return
			}
			memory.Cell[memory.Cursor] = int(b.buf[0])
		}
	}
}

// write method prints the value in current Cell of the memory
// if any error happen during the Write operation err property will be set.
func (b *brainFuck) write() Operator {
	return func(times int, memory *Memory) {
		b.buf[0] = byte(memory.Cell[memory.Cursor])
		for i := 0; i < times; i++ {
			if _, err := b.w.Write(b.buf); err != nil {
				b.err = err
				return
			}
		}
	}
}

// val method returns current value of which cursor is pointing.
func (b *brainFuck) val() int {
	return b.memory.Cell[b.cur()]
}

// AddOperator adds new Operator to the lexer
func (b *brainFuck) AddOperator(symbol rune, operator Operator) error {
	return token.AddOperator(symbol, operator)
}

// RemoveOperator removes Operator from the lexer
func (b *brainFuck) RemoveOperator(symbol rune) error {
	return token.RemoveOperator(symbol)
}

func (b *brainFuck) GetValueInMemory(position int) int {
	if position < 0 || position > len(b.memory.Cell) {
		return 0
	}
	return b.memory.Cell[position]
}
