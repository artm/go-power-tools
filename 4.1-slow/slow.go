package slow

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"time"
)

type Printer struct {
	reader io.Reader
	writer io.Writer
	delay  time.Duration
}

type option func(*Printer)

func NewPrinter(options ...option) *Printer {
	p := &Printer{
		reader: os.Stdin,
		writer: os.Stdout,
		delay:  1 * time.Second,
	}
	for _, o := range options {
		o(p)
	}
	return p
}

func WithReader(reader io.Reader) option {
	return func(p *Printer) {
		p.reader = reader
	}
}

func WithWriter(writer io.Writer) option {
	return func(p *Printer) {
		p.writer = writer
	}
}

func WithDelay(delay time.Duration) option {
	return func(p *Printer) {
		p.delay = delay
	}
}

func (p *Printer) Print() error {
	reader := bufio.NewReader(p.reader)
	for {
		rune, size, err := reader.ReadRune()
		if err != nil && err != io.EOF {
			return err
		}
		if size > 0 {
			fmt.Fprint(p.writer, string(rune))
		}
		if err == io.EOF {
			return nil
		}
		time.Sleep(p.delay)
	}
}
