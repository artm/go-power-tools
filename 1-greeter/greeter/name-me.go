package greeter

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

func AskName(out io.Writer, in io.Reader) string {
	bufout := bufio.NewWriter(out)
	bufout.WriteString("What's your name? ")
	bufout.Flush()

	bufin := bufio.NewReader(in)
	name, _ := bufin.ReadString('\n')
	name = strings.TrimRight(name, "\r\n")
	return name
}

func Greet(out io.Writer, name string) {
	bufout := bufio.NewWriter(out)
	bufout.WriteString(fmt.Sprintf("Hello, %s!\n", name))
	bufout.Flush()
}
