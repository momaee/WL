package lexer_test

import (
	"strings"
	"testing"

	"github.com/momaee/WL/lexer"
)

func TestScanner_Read(t *testing.T) {
	r := strings.NewReader("<>   this is string [+++]")
	s := lexer.NewScanner(r)
	token := s.Scan()

	if token.Value != "<" {
		t.Errorf("expect < given %q", token.Value)
	}
}

func TestScanner_Scan(t *testing.T) {
	//below string contains long white space and three runes
	r := strings.NewReader("        [+]")
	s := lexer.NewScanner(r)

	// this consumes all white spaces
	s.Scan()

	// read  [ rune
	token := s.Scan()

	if token.Value != "[" {
		t.Errorf("expect [ given %q", token.Value)
	}

	// read  + rune
	token = s.Scan()
	if token.Value != "+" {
		t.Errorf("expect + given %q", token.Value)
	}

	// read the last rune
	token = s.Scan()
	if token.Value != "]" {
		t.Errorf("expect ] given %q", token.Value)
	}

}
