package text_clock

import (
	"fmt"
	"io"
	"os"
	"time"
)

type Printer struct {
	Writer io.Writer
}

func Print() {
	NewPrinter().Print(time.Now())
}

func NewPrinter() *Printer {
	return &Printer{
		Writer: os.Stdout,
	}
}

func (printer *Printer) Print(when time.Time) {
	fmt.Fprintf(
		printer.Writer,
		"It's %d minutes past %d\n",
		when.Minute(),
		when.Hour(),
	)
}
