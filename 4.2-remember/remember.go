package remember

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

type Memory struct {
	writer  io.Writer
	memPath string
	args    []string
}

type option func(*Memory) error

func WithWriter(writer io.Writer) option {
	return func(mem *Memory) error {
		mem.writer = writer
		return nil
	}
}

func WithMemPath(memPath string) option {
	return func(mem *Memory) error {
		mem.memPath = memPath
		return nil
	}
}

func WithArgs(args []string) option {
	return func(mem *Memory) error {
		mem.args = make([]string, len(args))
		copy(mem.args, args)
		return nil
	}
}

func Run(options ...option) error {
	mem, err := NewMemory(options...)
	if err != nil {
		return err
	}
	return mem.Run()
}

func NewMemory(options ...option) (*Memory, error) {
	mem := &Memory{
		writer:  os.Stdout,
		memPath: "remember.json",
		args:    os.Args[1:],
	}
	for _, opt := range options {
		if err := opt(mem); err != nil {
			return nil, err
		}
	}
	return mem, nil
}

func (mem *Memory) Run() error {
	reminders, err := mem.readReminders()
	if err != nil {
		return err
	}

	if len(mem.args) == 0 {
		mem.printReminders(reminders)
		return nil
	} else {
		rem := strings.Join(mem.args, " ")
		reminders = append(reminders, rem)
		err = mem.writeReminders(reminders)
		return err
	}
}

func (mem *Memory) readReminders() ([]string, error) {
	var reminders []string
	if _, err := os.Stat(mem.memPath); os.IsNotExist(err) {
		return reminders, nil
	}
	jsonString, err := ioutil.ReadFile(mem.memPath)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(jsonString, &reminders)
	return reminders, err
}

func (mem *Memory) printReminders(reminders []string) {
	for _, rem := range reminders {
		fmt.Fprintln(mem.writer, rem)
	}
}

func (mem *Memory) writeReminders(reminders []string) error {
	jsonString, err := json.Marshal(reminders)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(mem.memPath, jsonString, 0644)
	return err
}
