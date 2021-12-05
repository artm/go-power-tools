package greeter

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type Greeter struct {
	In  io.Reader
	Out io.Writer
}

func NewGreeter() Greeter {
	return Greeter{
		In:  os.Stdin,
		Out: os.Stdout,
	}
}

func (greeter Greeter) Greet() {
	fmt.Fprint(greeter.Out, "What's your name? ")

	bufin := bufio.NewReader(greeter.In)
	name, _ := bufin.ReadString('\n')
	name = strings.TrimRight(name, "\r\n")

	fmt.Fprintf(greeter.Out, "Hello, %s!\n", name)
}

func Greet() {
	NewGreeter().Greet()
}
