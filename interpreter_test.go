package interpreter_test

import (
	"bytes"
	"fmt"
	"math"
	"strings"
	"testing"

	interpreter "github.com/momaee/WL"
	"github.com/stretchr/testify/assert"
)

func TestBrainFuckMachine_LoopOperation(t *testing.T) {
	code := strings.NewReader("----[---->+<]>++.+.+.+.")

	input := new(bytes.Buffer)

	output := new(bytes.Buffer)

	bfm := interpreter.NewInterpreter(input, output, code)

	err := bfm.Run()
	assert.NoError(t, err)

	if output.String() != "ABCD" {
		t.Errorf("wrong output, got %s", output.String())
	}
}

func TestBrainFuckMachine_PrintHelloWorld(t *testing.T) {
	code := strings.NewReader(`++++++++[>++++[>++>+++>+++>+<<<<-]>+>+>->>+[<]<-]>>.>---.+++++++..+++.>>.<-.<.+++.------.--------.>>+.>++.`)

	i := new(bytes.Buffer)

	o := new(bytes.Buffer)

	bfm := interpreter.NewInterpreter(i, o, code)

	err := bfm.Run()
	assert.NoError(t, err)

	if o.String() != "Hello World!\n" {
		t.Errorf("wrong output, got %s", o.String())
	}

}

func TestGetValueInMemory(t *testing.T) {
	code := strings.NewReader("+")

	i := new(bytes.Buffer)

	o := new(bytes.Buffer)

	bfm := interpreter.NewInterpreter(i, o, code)

	err := bfm.Run()
	assert.NoError(t, err)

	if bfm.GetValueInMemory(0) != 1 {
		t.Errorf("wrong value, got %d", bfm.GetValueInMemory(0))
	}
}

func TestAddOperatorError(t *testing.T) {
	code := strings.NewReader("+")

	i := new(bytes.Buffer)

	o := new(bytes.Buffer)

	bfm := interpreter.NewInterpreter(i, o, code)

	err := bfm.AddOperator('+', func(c int, memory *interpreter.Memory) {
		memory.Cell[memory.Cursor] = (memory.Cell[memory.Cursor] * 2 * c) % 255
	})
	assert.Error(t, err)

	err = bfm.Run()
	assert.NoError(t, err)

	if bfm.GetValueInMemory(0) != 1 {
		t.Errorf("wrong value, got %d", bfm.GetValueInMemory(0))
	}
}

func TestAddOperatorNoError(t *testing.T) {
	code := strings.NewReader("+++***") // (1+1+1)*(2*2*2)

	i := new(bytes.Buffer)

	o := new(bytes.Buffer)

	bfm := interpreter.NewInterpreter(i, o, code)

	err := bfm.AddOperator('*', func(c int, memory *interpreter.Memory) {
		memory.Cell[memory.Cursor] = (memory.Cell[memory.Cursor] * int(math.Pow(2, float64(c)))) % 255
	})
	assert.NoError(t, err)

	err = bfm.Run()
	assert.NoError(t, err)

	if bfm.GetValueInMemory(0) != 24 {
		t.Errorf("wrong value, got %d", bfm.GetValueInMemory(0))
	}
}

func TestRemoveOperatorError(t *testing.T) {
	t.Run("1", func(t *testing.T) {
		code := strings.NewReader("+")

		i := new(bytes.Buffer)

		o := new(bytes.Buffer)

		bfm := interpreter.NewInterpreter(i, o, code)

		err := bfm.RemoveOperator('/')
		fmt.Println("err", err)
		assert.Error(t, err)

		err = bfm.Run()
		assert.NoError(t, err)

		if bfm.GetValueInMemory(0) != 1 {
			t.Errorf("wrong value, got %d", bfm.GetValueInMemory(0))
		}
	})
}

func TestRemoveOperatorNoError(t *testing.T) {
	t.Run("1", func(t *testing.T) {
		code := strings.NewReader("+++###") // (1+1+1) because # would be removed

		i := new(bytes.Buffer)

		o := new(bytes.Buffer)

		bfm := interpreter.NewInterpreter(i, o, code)

		err := bfm.AddOperator('#', func(c int, memory *interpreter.Memory) {
			memory.Cell[memory.Cursor] = (memory.Cell[memory.Cursor] * int(math.Pow(2, float64(c)))) % 255
		})
		assert.NoError(t, err)

		err = bfm.RemoveOperator('#')
		assert.NoError(t, err)

		err = bfm.Run()
		assert.NoError(t, err)

		if bfm.GetValueInMemory(0) != 3 {
			t.Errorf("wrong value, got %d", bfm.GetValueInMemory(0))
		}
	})

	t.Run("2", func(t *testing.T) {
		code := strings.NewReader("++++")

		i := new(bytes.Buffer)

		o := new(bytes.Buffer)

		bfm := interpreter.NewInterpreter(i, o, code)

		err := bfm.RemoveOperator('+')
		assert.NoError(t, err)

		err = bfm.Run()
		assert.NoError(t, err)

		if bfm.GetValueInMemory(0) != 0 {
			t.Errorf("wrong value, got %d", bfm.GetValueInMemory(0))
		}

	})
}
