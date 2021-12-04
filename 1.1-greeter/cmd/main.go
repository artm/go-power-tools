package main

import (
	"greeter"
	"os"
)

func main() {
	greeter.Greet(os.Stdin, os.Stdout)
}
