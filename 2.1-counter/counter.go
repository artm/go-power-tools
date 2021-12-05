package counter

import (
	"fmt"
	"io"
	"os"
	"time"
)

type Counter struct {
	counter int
	Writer  io.Writer
	Stop    chan bool
	Delay   time.Duration
}

func NewCounter() Counter {
	return Counter{
		Writer: os.Stdout,
		Stop:   make(chan bool),
		Delay:  10 * time.Minute,
	}
}

func (counter *Counter) Next() int {
	current := counter.counter
	counter.counter++
	return current
}

func (counter *Counter) Run() {
	for {
		select {
		case <-counter.Stop:
			return
		default:
			time.Sleep(counter.Delay)
			fmt.Fprintln(counter.Writer, counter.Next())
		}
	}
}
