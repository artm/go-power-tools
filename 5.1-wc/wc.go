package wc

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

type stat interface {
	update(*int, string)
	aggregate(*int, int)
}

type Wc struct {
	output io.Writer
	input  io.Reader
	stats  []stat
	paths  []string
}

type option func(*Wc) error

type lineCount struct{}
type wordCount struct{}
type charCount struct{}
type byteCount struct{}
type widthStat struct{}

type resultRow struct {
	numbers []int
	path    string
}

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
		lines := fset.Bool("l", false, "print the newline counts")
		words := fset.Bool("w", false, "print the word counts")
		chars := fset.Bool("m", false, "print the character counts")
		bytes := fset.Bool("c", false, "print the byte counts")
		width := fset.Bool("L", false, "print the maximum display width")
		err := fset.Parse(args)
		if err != nil {
			return err
		}
		if !(*bytes || *chars || *lines || *width || *words) {
			*lines = true
			*words = true
			*bytes = true
		}
		if *lines {
			wc.stats = append(wc.stats, lineCount{})
		}
		if *words {
			wc.stats = append(wc.stats, wordCount{})
		}
		if *chars {
			wc.stats = append(wc.stats, charCount{})
		}
		if *bytes {
			wc.stats = append(wc.stats, byteCount{})
		}
		if *width {
			wc.stats = append(wc.stats, widthStat{})
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
	results := make([]resultRow, 0)
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
		result := resultRow{
			path: path,
		}
		result.numbers, err = wc.countIn(f)
		if err != nil {
			return err
		}
		results = append(results, result)
	}
	if len(results) > 1 {
		total := wc.CalcTotal(results)
		results = append(results, total)
	}
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

func (wc *Wc) countIn(reader io.Reader) ([]int, error) {
	counts := make([]int, len(wc.stats))
	streader := bufio.NewReader(reader)
	for {
		line, err := streader.ReadString('\n')
		if len(line) > 0 {
			for i, stat := range wc.stats {
				stat.update(&counts[i], line)
			}
		}
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
	}
	return counts, nil
}

func (wc *Wc) CalcTotal(results []resultRow) resultRow {
	total := resultRow{
		path:    "total",
		numbers: make([]int, len(wc.stats)),
	}
	for _, result := range results {
		for i, stat := range wc.stats {
			stat.aggregate(&total.numbers[i], result.numbers[i])
		}
	}
	return total
}

func (wc *Wc) Print(results []resultRow) {
	statsCount := len(wc.stats)
	wideStdin := len(results) > 1 || statsCount > 1
	colWidth := 1
	if len(results) > 1 && statsCount == 1 {
		colWidth = 4
	}
	for _, result := range results {
		for _, count := range result.numbers {
			width := len(strconv.Itoa(count))
			if colWidth < width {
				colWidth = width
			}
			if wideStdin &&
				(result.path == "-" || result.path == "") &&
				colWidth < 7 {
				colWidth = 7
			}
		}
	}

	colFmt := fmt.Sprintf("%%%dd", colWidth)
	for _, result := range results {
		var row []string
		for _, count := range result.numbers {
			row = append(row, fmt.Sprintf(colFmt, count))
		}
		if result.path != "" {
			row = append(row, result.path)
		}
		fmt.Fprintln(wc.output, strings.Join(row, " "))
	}
}

func (lineCount) update(count *int, line string) {
	if strings.HasSuffix(line, "\n") {
		*count++
	}
}

func (lineCount) aggregate(agg *int, count int) {
	*agg += count
}

var wordRe = regexp.MustCompile(`\pL+`)

func (wordCount) update(count *int, line string) {
	*count += len(wordRe.FindAllString(line, -1))
}

func (wordCount) aggregate(agg *int, count int) {
	*agg += count
}

func (charCount) update(count *int, line string) {
	*count += utf8.RuneCountInString(line)
}

func (charCount) aggregate(agg *int, count int) {
	*agg += count
}

func (byteCount) update(count *int, line string) {
	*count += len(line)
}

func (byteCount) aggregate(agg *int, count int) {
	*agg += count
}

func (ws widthStat) update(width *int, line string) {
	line = strings.TrimRight(line, "\n")
	runeCount := utf8.RuneCountInString(line)
	ws.aggregate(width, runeCount)
}

func (widthStat) aggregate(agg *int, width int) {
	if *agg < width {
		*agg = width
	}
}
