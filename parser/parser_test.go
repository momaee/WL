package parser_test

import (
	"strings"
	"testing"

	"github.com/momaee/WL/lexer"
	"github.com/momaee/WL/parser"
	"github.com/momaee/WL/token"
)

func TestParser_Parse(t *testing.T) {
	input := strings.NewReader("+++++ -- [-]")

	lexer := lexer.NewScanner(input)

	p := parser.NewParser(lexer)

	instructions := p.Parse()

	// since we are folding instructions
	// there is one instruction +, but 4 times
	if len(instructions) != 5 {
		t.Errorf("wrong length, expected 5 got %+v", len(instructions))
	}
	expected := []*parser.Inst{
		{T: &token.Token{Tok: token.PlusToken, Value: "+"}, C: 5},
		{T: &token.Token{Tok: token.MinusToken, Value: "-"}, C: 2},
		{T: &token.Token{Tok: token.LeftBracketToken, Value: "["}, C: 4},
		{T: &token.Token{Tok: token.MinusToken, Value: "-"}, C: 1},
		{T: &token.Token{Tok: token.RightBracketToken, Value: "]"}, C: 2},
	}
	for i, v := range expected {
		if v.T.Tok != instructions[i].T.Tok || v.C != instructions[i].C || v.T.Value != instructions[i].T.Value {
			t.Errorf("incorrect instruction. expected %+v got %+v", *v, *instructions[i])
		}
	}
}

func TestInnerLoops(t *testing.T) {
	input := strings.NewReader("-[--[+]--]")
	// how to index above input:0122345667

	lexer := lexer.NewScanner(input)

	p := parser.NewParser(lexer)

	instructions := p.Parse()

	expected := []*parser.Inst{
		{T: &token.Token{Tok: token.MinusToken, Value: "-"}, C: 1},
		{T: &token.Token{Tok: token.LeftBracketToken, Value: "["}, C: 7}, // index of respective closing bracket is 7
		{T: &token.Token{Tok: token.MinusToken, Value: "-"}, C: 2},
		{T: &token.Token{Tok: token.LeftBracketToken, Value: "["}, C: 5},
		{T: &token.Token{Tok: token.PlusToken, Value: "+"}, C: 1},
		{T: &token.Token{Tok: token.RightBracketToken, Value: "]"}, C: 3}, // index of respective opening bracket is 3
		{T: &token.Token{Tok: token.MinusToken, Value: "-"}, C: 2},
		{T: &token.Token{Tok: token.RightBracketToken, Value: "]"}, C: 1},
	}

	for i, v := range expected {
		if v.T.Tok != instructions[i].T.Tok || v.C != instructions[i].C || v.T.Value != instructions[i].T.Value {
			t.Errorf("incorrect instruction. expected %+v got %+v", *v, *instructions[i])
		}
	}
}

func Test_MoveBetweenCells(t *testing.T) {
	input := strings.NewReader("+>>>+++++++>>+++ --<<")

	lexer := lexer.NewScanner(input)

	p := parser.NewParser(lexer)

	instructions := p.Parse()

	expected := []*parser.Inst{
		{T: &token.Token{Tok: token.PlusToken, Value: "+"}, C: 1},
		{T: &token.Token{Tok: token.RightToken, Value: ">"}, C: 3},
		{T: &token.Token{Tok: token.PlusToken, Value: "+"}, C: 7},
		{T: &token.Token{Tok: token.RightToken, Value: ">"}, C: 2},
		{T: &token.Token{Tok: token.PlusToken, Value: "+"}, C: 3},
		{T: &token.Token{Tok: token.MinusToken, Value: "-"}, C: 2},
		{T: &token.Token{Tok: token.LeftToken, Value: "<"}, C: 2},
	}

	for i, v := range expected {
		if v.T.Tok != instructions[i].T.Tok || v.C != instructions[i].C || v.T.Value != instructions[i].T.Value {
			t.Errorf("incorrect instruction. expected %+v got %+v", *v, *instructions[i])
		}
	}

}
