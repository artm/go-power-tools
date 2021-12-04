package main

import (
	"greeter/greeter"
	"os"
)

func main() {
	greeter.Greet(os.Stdin, os.Stdout)
}
