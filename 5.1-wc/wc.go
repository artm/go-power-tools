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
	output    io.Writer
	input     io.Reader
	stats     []stat
	paths     []string
	printPath bool
}

type option func(*Wc) error

type lineCount struct{}
type wordCount struct{}
type charCount struct{}
type byteCount struct{}
type widthStat struct{}

var defaultStats = []stat{lineCount{}, wordCount{}, byteCount{}}
var wordRe = regexp.MustCompile(`\pL+`)

type resultRow struct {
	numbers []int
	path    string
}

func NewWc(options ...option) (*Wc, error) {
	wc := &Wc{
		output:    os.Stdout,
		printPath: false,
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
		statsOrder := []struct {
			enabled *bool
			stat    stat
		}{
			{fset.Bool("l", false, "print the newline counts"), lineCount{}},
			{fset.Bool("w", false, "print the word counts"), wordCount{}},
			{fset.Bool("m", false, "print the character counts"), charCount{}},
			{fset.Bool("c", false, "print the byte counts"), byteCount{}},
			{fset.Bool("L", false, "print the maximum display width"), widthStat{}},
		}
		if err := fset.Parse(args); err != nil {
			return err
		}
		for _, rec := range statsOrder {
			if *rec.enabled {
				wc.stats = append(wc.stats, rec.stat)
			}
		}
		if len(wc.stats) == 0 {
			wc.stats = defaultStats
		}
		wc.paths = fset.Args()
		if len(wc.paths) == 0 {
			wc.paths = []string{"-"}
		} else {
			wc.printPath = true
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
		var file io.ReadCloser
		var err error
		if path == "-" {
			file = wc.inputCloser()
		} else {
			file, err = os.Open(path)
			if err != nil {
				return err
			}
		}
		defer file.Close()
		result := resultRow{
			path: path,
		}
		result.numbers, err = wc.countIn(file)
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

func (wc *Wc) inputCloser() io.ReadCloser {
	file, ok := wc.input.(io.ReadCloser)
	if !ok {
		file = io.NopCloser(wc.input)
	}
	return file
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
				(result.path == "-") &&
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
		if wc.printPath {
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
