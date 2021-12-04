package text_clock

import (
	"bufio"
	"fmt"
	"io"
	"time"
)

func PrintTime(when time.Time, writer io.Writer) {
	timeString := fmt.Sprintf(
		"It's %d minutes past %d\n",
		when.Minute(),
		when.Hour())
	out := bufio.NewWriter(writer)
	out.WriteString(timeString)
	out.Flush()
}
