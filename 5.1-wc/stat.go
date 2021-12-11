package wc

import (
	"strings"
	"unicode/utf8"
)

type stat interface {
	update(*int, string)
	aggregate(*int, int)
}

type lineCount struct{}
type wordCount struct{}
type charCount struct{}
type byteCount struct{}
type widthStat struct{}

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
