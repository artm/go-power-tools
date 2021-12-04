package main

import (
	"greeter/greeter"
	"os"
)

func main() {
	name := greeter.AskName(os.Stdout, os.Stdin)
	greeter.Greet(os.Stdout, name)
}
