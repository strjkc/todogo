package main

import (
	"errors"
	"fmt"
	"os"
	"time"
)

type State struct {
	Counter int     `json:"counter"`
	Tasks   []*Task `json:"tasks"`
}

type Task struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	Status      Status    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type Status string

const (
	ToDo       Status = "todo"
	InProgress Status = "in-progress"
	Done       Status = "done"
)

const STORAGEFILE = "storage.json"

var (
	TasksMap map[Status][]*Task
	Commands map[string]func(s *State, args []string) error
)

func initCommands() {
	Commands = make(map[string]func(s *State, args []string) error)
	Commands["add"] = handleAdd
	Commands["delete"] = handleDelete
	Commands["update"] = handleUpdate
	Commands["list"] = handleList
	Commands["mark-in-progress"] = handleMvInProg
	Commands["mark-done"] = handleMvDone
}

func main() {
	state := &State{}
	TasksMap = make(map[Status][]*Task)
	initCommands()
	err := loadTasks(state)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		panic(err)
	}
	cliArgs := os.Args[1:]
	command := cliArgs[0]
	commandArgs := cliArgs[1:]
	fun := Commands[command]
	err = fun(state, commandArgs)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
}
