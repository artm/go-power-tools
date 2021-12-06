package grep

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type Searcher struct {
	reader io.Reader
	writer io.Writer
}

type option func(*Searcher)

func NewSearcher(options ...option) Searcher {
	searcher := Searcher{
		reader: os.Stdin,
		writer: os.Stdout,
	}
	for _, optfunc := range options {
		optfunc(&searcher)
	}
	return searcher
}

func WithWriter(writer io.Writer) option {
	return func(searcher *Searcher) {
		searcher.writer = writer
	}
}

func WithReader(reader io.Reader) option {
	return func(searcher *Searcher) {
		searcher.reader = reader
	}
}

func (searcher *Searcher) Search(what string) {
	scanner := bufio.NewScanner(searcher.reader)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, what) {
			fmt.Fprintln(searcher.writer, line)
		}
	}
}
