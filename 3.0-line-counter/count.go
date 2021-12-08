package count

import (
	"bufio"
	"errors"
	"io"
	"os"
	"strings"
)

type counter struct {
	input   io.Reader
	output  io.Writer
	pattern string
	fold    func(string) string
}

type option func(*counter) error

func NewCounter(opts ...option) (counter, error) {
	c := counter{
		input:  os.Stdin,
		output: os.Stdout,
		fold:   func(a string) string { return a },
	}
	for _, opt := range opts {
		err := opt(&c)
		if err != nil {
			return counter{}, err
		}
	}
	c.pattern = c.fold(c.pattern)
	return c, nil
}

func WithInput(input io.Reader) option {
	return func(c *counter) error {
		if input == nil {
			return errors.New("nil input reader")
		}
		c.input = input
		return nil
	}
}

func WithOutput(output io.Writer) option {
	return func(c *counter) error {
		if output == nil {
			return errors.New("nil output writer")
		}
		c.output = output
		return nil
	}
}

func WithPattern(pattern string) option {
	return func(c *counter) error {
		c.pattern = pattern
		return nil
	}
}

func IgnoreCase() option {
	return func(c *counter) error {
		c.fold = strings.ToLower
		return nil
	}
}

func (c counter) Lines() int {
	lines := 0
	scanner := bufio.NewScanner(c.input)
	for scanner.Scan() {
		line := c.fold(scanner.Text())
		if strings.Contains(line, c.pattern) {
			lines++
		}
	}
	return lines
}

func Lines() int {
	c, err := NewCounter()
	if err != nil {
		panic("internal error")
	}
	return c.Lines()
}
