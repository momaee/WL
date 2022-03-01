package token

import "fmt"

const (
	IllegalToken      Type = iota
	LeftToken              // <
	RightToken             // >
	PlusToken              // +
	MinusToken             // -
	PrintToken             // .
	ReadToken              // ,
	LeftBracketToken       // [
	RightBracketToken      // ]
	WhitespaceToken
	UserDefinedToken
)

// Memory capacity
const MemorySize int = 3000

// special var to indicate end of stream.
var EOF = rune(-1)

// Type represents a lexical token type.
type Type int

// Token represents a lexical tokens.
type Token struct {
	// the type of token.
	Tok Type

	// The literal value of the token(as parsed).
	Value string

	Operator Operator
}

type Memory struct {
	Cell   [MemorySize]int
	Cursor int
}

// C is complementary information about instruction like position or counts of occurrence
// Incase of opening loop, C is the index of the closing loop and vice versa
type Operator func(c int, memory *Memory)

var (
	AllTokens = map[rune]*Token{
		'<': {Tok: LeftToken, Value: "<", Operator: seekBwd},
		'>': {Tok: RightToken, Value: ">", Operator: seekFwd},
		'+': {Tok: PlusToken, Value: "+", Operator: inc},
		'-': {Tok: MinusToken, Value: "-", Operator: dec},
		'.': {Tok: PrintToken, Value: "."},
		',': {Tok: ReadToken, Value: ","},
		'[': {Tok: LeftBracketToken, Value: "["},
		']': {Tok: RightBracketToken, Value: "]"},
	}
)

// dec method decrements the value of the current Cell in memory by v.
var dec Operator = func(c int, memory *Memory) {
	if memory.Cell[memory.Cursor]-c >= 0 {
		memory.Cell[memory.Cursor] -= c
	} else {
		memory.Cell[memory.Cursor] = 256 + memory.Cell[memory.Cursor] - c
	}
}

// inc method increments the value of the current Cell in memory by v.
var inc Operator = func(c int, memory *Memory) {
	memory.Cell[memory.Cursor] = (memory.Cell[memory.Cursor] + c) % 255
}

// seekFwd method moves the cursor in the memory forward by c.
// this move is relative to current cursor position
var seekFwd Operator = func(c int, memory *Memory) {
	memory.Cursor += c
}

// seekFwd method moves the cursor in the memory backward by c.
// this move is relative to current cursor position
var seekBwd Operator = func(c int, memory *Memory) {
	memory.Cursor -= c
}

func AddOperator(symbol rune, operator Operator) error {
	if _, ok := AllTokens[symbol]; ok {
		return fmt.Errorf("symbol %v already exists", symbol)
	}
	AllTokens[symbol] = &Token{Tok: UserDefinedToken, Operator: operator, Value: string(symbol)}
	return nil
}

func RemoveOperator(symbol rune) error {
	fmt.Println("symbol", symbol)
	fmt.Println("all tokens:", AllTokens)
	if _, ok := AllTokens[symbol]; !ok {
		return fmt.Errorf("symbol %v does not exist", symbol)
	}
	delete(AllTokens, symbol)
	return nil
}

func (t *Token) HasOperator() bool {
	return t.Operator != nil
}
