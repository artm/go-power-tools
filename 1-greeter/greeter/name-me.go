package greeter

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

func Greet(out io.Writer, in io.Reader) {
	bufout := bufio.NewWriter(out)
	bufout.WriteString("What's your name? ")
	bufout.Flush()

	bufin := bufio.NewReader(in)
	name, _ := bufin.ReadString('\n')
	name = strings.TrimRight(name, "\r\n")

	bufout.WriteString(fmt.Sprintf("Hello, %s!\n", name))
	bufout.Flush()
}
