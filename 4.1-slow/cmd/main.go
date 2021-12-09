package main

import (
	"os"
	"slow"
)

func main() {
	slow.Print(slow.WithArgs(os.Args[1:]))
}
