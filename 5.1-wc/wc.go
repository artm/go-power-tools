package wc

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type Wc struct {
	output                            io.Writer
	bytes, chars, lines, width, words bool
	paths                             []string
}

type option func(*Wc) error

func NewWc(options ...option) (*Wc, error) {
	wc := &Wc{
		output: os.Stdout,
	}
	for _, option := range options {
		if err := option(wc); err != nil {
			return nil, err
		}
	}
	return wc, nil
}

func WithArgs(args []string) option {
	return func(wc *Wc) error {
		fset := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
		bytes := fset.Bool("c", false, "print the byte counts")
		chars := fset.Bool("m", false, "print the character counts")
		lines := fset.Bool("l", false, "print the newline counts")
		width := fset.Bool("L", false, "print the maximum display width")
		words := fset.Bool("w", false, "print the word counts")
		err := fset.Parse(args)
		if err != nil {
			return err
		}
		wc.bytes = *bytes
		wc.chars = *chars
		wc.lines = *lines
		wc.width = *width
		wc.words = *words
		if wc.flagCount() == 0 {
			wc.lines = true
			wc.words = true
			wc.bytes = true
		}
		wc.paths = fset.Args()
		return nil
	}
}

func WithOutput(output io.Writer) option {
	return func(wc *Wc) error {
		wc.output = output
		return nil
	}
}

func (wc *Wc) Count() error {
	for _, path := range wc.paths {
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		wc.countIn(f, path)
	}
	return nil
}

func Count() error {
	wc, err := NewWc(WithArgs(os.Args[1:]))
	if err != nil {
		return err
	}
	err = wc.Count()
	if err != nil {
		return err
	}
	return nil
}

func (wc *Wc) countIn(reader io.Reader, path string) {
	var lines, bytes, chars, width, words int
	streader := bufio.NewReader(reader)
	for {
		line, err := streader.ReadString('\n')
		if err == io.EOF {
			break
		}
		lines++
		bytes += len(line)
		lineChars := len([]rune(line))
		chars += lineChars
		if lineChars > width {
			width = lineChars - 1
		}
		words += len(strings.Split(line, " "))
	}
	tokens := []string{}
	if wc.flagCount() == 1 {
		var count int
		switch {
		case wc.bytes:
			count = bytes
		case wc.chars:
			count = chars
		case wc.lines:
			count = lines
		case wc.width:
			count = width
		case wc.words:
			count = words
		}
		tokens = append(tokens, strconv.Itoa(count))
	} else {
		order := []struct {
			on    bool
			count int
		}{
			{on: wc.lines, count: lines},
			{on: wc.words, count: words},
			{on: wc.chars, count: chars},
			{on: wc.bytes, count: bytes},
			{on: wc.width, count: width},
		}
		for _, rec := range order {
			if rec.on {
				tokens = append(tokens, fmt.Sprintf("%2d", rec.count))
			}
		}
	}
	tokens = append(tokens, path)
	fmt.Fprintf(wc.output, "%s\n", strings.Join(tokens, " "))
}

func (wc *Wc) flagCount() int {
	flags := []bool{
		wc.bytes,
		wc.chars,
		wc.lines,
		wc.width,
		wc.words,
	}
	var count int
	for _, flag := range flags {
		if flag {
			count++
		}
	}
	return count
}
