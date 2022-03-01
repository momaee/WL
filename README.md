# Brainfucker Interpreter

## How to interpret a brainfuck code

Please read this [article](https://levelup.gitconnected.com/write-your-own-brainfuck-interpreter-98e828c72854)

## How to use

1. First of all please see the test files.

2. typical usage

    ```go
    // create new io.Reader from inputs
    code := strings.NewReader("++++++++[>++++[>++>+++>+++>+<<<<-]>+>+>->>+[<]<-]>>.>---.+++++++..+++.>>.<-.<.+++.------.--------.>>+.>++.") //Hello World!\n

    // Standards interface to io
    input := new(bytes.Buffer)
    output := new(bytes.Buffer)

    // initialize the machine
    bfm := interpreter.NewInterpreter(input, output, code)

    // Store the result in output interface 
    err := bfm.Run()
    if err != nil {
        //handle err
    }

    // print the result 
    fmt.Println (output.String()) 
    ```

3. Add custom operators

    ```go
    // create new io.Reader from inputs
    code := strings.NewReader("+++***") // (1+1+1)*(2*2*2) = 24

    // Standards interface to io
    input := new(bytes.Buffer)
    output := new(bytes.Buffer)

    // initialize the machine
    bfm := interpreter.NewInterpreter(input, output, code)

    // Add your own operator
    // C is complementary information about instructionlike position or counts of occurrence
    // Incase of opening loop, C is the index of the closing loop and vice versa
    err := bfm.AddOperator('*', func(c int, memory *interpreter.Memory) {
        memory.Cell[memory.Cursor] = (memory.Cell[memory.Cursor] * int(math.Pow(2, float64(c)))) % 255
    })

    // Store the result in output interface 
    err := bfm.Run()
    if err != nil {
        //handle err
    }

    // print the result of arbitrary cell in the memory
    fmt.Println (bfm.GetValueInMemory(0))
    ```

## Run tests

In the root of the project run ```go test ./...```
