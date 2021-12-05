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
	bufout := bufio.NewWriter(greeter.Out)
	bufout.WriteString("What's your name? ")
	bufout.Flush()

	bufin := bufio.NewReader(greeter.In)
	name, _ := bufin.ReadString('\n')
	name = strings.TrimRight(name, "\r\n")

	bufout.WriteString(fmt.Sprintf("Hello, %s!\n", name))
	bufout.Flush()
}

func Greet() {
	NewGreeter().Greet()
}
