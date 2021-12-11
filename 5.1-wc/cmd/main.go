package main

import (
	"fmt"
	"os"
	"wc"
)

func main() {
	if err := wc.Count(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
