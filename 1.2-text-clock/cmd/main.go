package main

import (
	"os"
	text_clock "text-clock"
	"time"
)

func main() {
	text_clock.PrintTime(time.Now(), os.Stdout)
}
