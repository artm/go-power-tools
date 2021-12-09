package slow

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"time"
)

type Printer struct {
	reader io.ReadCloser
	writer io.Writer
	delay  time.Duration
}

type option func(*Printer) error

func NewPrinter(options ...option) (*Printer, error) {
	p := &Printer{
		reader: os.Stdin,
		writer: os.Stdout,
		delay:  1 * time.Second,
	}
	for _, o := range options {
		err := o(p)
		if err != nil {
			return nil, err
		}
	}
	return p, nil
}

func WithReader(reader io.Reader) option {
	readCloser, ok := reader.(io.ReadCloser)
	if !ok {
		readCloser = io.NopCloser(reader)
	}
	return func(p *Printer) error {
		p.reader = readCloser
		return nil
	}
}

func WithWriter(writer io.Writer) option {
	return func(p *Printer) error {
		p.writer = writer
		return nil
	}
}

func WithDelay(delay time.Duration) option {
	return func(p *Printer) error {
		p.delay = delay
		return nil
	}
}

func WithArgs(args []string) option {
	return func(p *Printer) error {
		if len(args) < 1 {
			return nil
		}
		var err error
		p.reader, err = os.Open(args[0])
		return err
	}
}

func (p *Printer) Print() error {
	defer p.reader.Close()
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

func Print(options ...option) {
	printer, err := NewPrinter(options...)
	if err != nil {
		panic(err)
	}
	err = printer.Print()
	if err != nil {
		panic(err)
	}
}
