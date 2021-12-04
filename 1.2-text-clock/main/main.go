package main

import (
	"os"
	"text-clock/text_clock"
	"time"
)

func main() {
	text_clock.PrintTime(time.Now(), os.Stdout)
}
