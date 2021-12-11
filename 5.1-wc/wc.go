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
	input                             io.Reader
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
		if !(*bytes || *chars || *lines || *width || *words) {
			wc.lines = true
			wc.words = true
			wc.bytes = true
		}
		wc.paths = fset.Args()
		if len(wc.paths) == 0 {
			wc.paths = []string{""}
		}
		return nil
	}
}

func WithOutput(output io.Writer) option {
	return func(wc *Wc) error {
		wc.output = output
		return nil
	}
}

func WithInput(input io.Reader) option {
	return func(wc *Wc) error {
		wc.input = input
		return nil
	}
}

func (wc *Wc) Count() error {
	results := make([][]string, 0)
	for _, path := range wc.paths {
		var f io.ReadCloser
		var err error
		if path == "-" || path == "" {
			var ok bool
			f, ok = wc.input.(io.ReadCloser)
			if !ok {
				f = io.NopCloser(wc.input)
			}
		} else {
			f, err = os.Open(path)
			if err != nil {
				return err
			}
		}
		result, err := wc.countIn(f)
		if err != nil {
			return err
		}
		result = append(result, path)
		results = append(results, result)
	}
	results = wc.CalcTotals(results)
	wc.Print(results)
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

func (wc *Wc) countIn(reader io.Reader) ([]string, error) {
	var lines, bytes, chars, width, words int
	streader := bufio.NewReader(reader)
	for {
		line, err := streader.ReadString('\n')
		if line != "" {
			bytes += len(line)
			lineChars := len([]rune(line))
			chars += lineChars
			if lineChars > width {
				width = lineChars
				if strings.HasSuffix(line, "\n") {
					width--
					lines++
				}
			}
			words += len(strings.Split(line, " "))
		}
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
	}
	results := make([]string, 0)
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
			results = append(results, strconv.Itoa(rec.count))
		}
	}
	return results, nil
}

func (wc *Wc) CalcTotals(results [][]string) [][]string {
	if len(results) < 2 {
		return results
	}
	colNum := len(results[0])
	totals := make([]int, colNum-1)
	for _, result := range results {
		for i, count := range result[:colNum-1] {
			iCount, _ := strconv.Atoi(count)
			totals[i] += iCount
		}
	}
	sTotals := make([]string, colNum)
	for i, v := range totals {
		sTotals[i] = strconv.Itoa(v)
	}
	sTotals[colNum-1] = "total"
	results = append(results, sTotals)
	return results
}

func (wc *Wc) Print(results [][]string) {
	colWidth := 1
	for _, result := range results {
		for _, count := range result[:len(result)-1] {
			if colWidth < len(count) {
				colWidth = len(count)
			}
			path := result[len(result)-1]
			if len(result) > 2 &&
				(path == "-" || path == "") &&
				colWidth < 7 {
				colWidth = 7
			}
		}
	}

	colFmt := fmt.Sprintf("%%%ds", colWidth)
	for _, result := range results {
		var row []string
		for _, count := range result[:len(result)-1] {
			row = append(row, fmt.Sprintf(colFmt, count))
		}
		if path := result[len(result)-1]; path != "" {
			row = append(row, path)
		}
		fmt.Fprintln(wc.output, strings.Join(row, " "))
	}
}
