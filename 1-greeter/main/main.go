package main

import (
	"greeter/greeter"
	"os"
)

func main() {
	greeter.Greet(os.Stdout, os.Stdin)
}
