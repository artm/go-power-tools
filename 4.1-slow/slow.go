package slow

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"time"
)

type Printer struct {
	readers []io.ReadCloser
	writer  io.Writer
	delay   time.Duration
}

type option func(*Printer) error

func NewPrinter(options ...option) (*Printer, error) {
	p := &Printer{
		readers: []io.ReadCloser{os.Stdin},
		writer:  os.Stdout,
		delay:   1 * time.Second,
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
		p.readers = []io.ReadCloser{readCloser}
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
		if len(args) == 0 {
			return nil
		}
		p.readers = []io.ReadCloser{}
		for _, arg := range args {
			reader, err := os.Open(arg)
			if err != nil {
				return err
			}
			p.readers = append(p.readers, reader)
		}
		return nil
	}
}

func (p *Printer) Print() error {
	for _, reader := range p.readers {
		defer reader.Close()
		bufReader := bufio.NewReader(reader)
		for {
			rune, _, err := bufReader.ReadRune()
			if err == io.EOF {
				break
			} else if err != nil {
				return err
			}
			fmt.Fprint(p.writer, string(rune))
			time.Sleep(p.delay)
		}
	}
	return nil
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
