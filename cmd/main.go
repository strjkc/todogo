package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"
)

type State struct {
	Counter int     `json:"counter"`
	Tasks   []*Task `json:"tasks"`
}

type Status string

const (
	ToDo       Status = "todo"
	InProgress Status = "in-progress"
	Done       Status = "done"
)

type Task struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	Status      Status    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

var TasksMap map[Status][]*Task

var Commands map[string]func(s *State, args []string) error

func initCommands() {
	Commands = make(map[string]func(s *State, args []string) error)
	Commands["add"] = handleAdd
	Commands["delete"] = handleDelete
	Commands["update"] = handleUpdate
	Commands["list"] = handleList
	Commands["mark-in-progress"] = handleMvInProg
	Commands["mark-done"] = handleMvDone
}

func handleAdd(s *State, args []string) error {
	if len(args) != 1 {
		return errors.New("invalid argument for add")
	}
	arg := args[0]
	s.Counter += 1
	task := Task{ID: s.Counter, Status: ToDo, Description: arg, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	s.Tasks = append(s.Tasks, &task)
	err := saveState(s)
	if err != nil {
		s.Counter -= 1
		return err
	}
	return nil
}

func handleDelete(s *State, args []string) error {
	if len(args) != 1 {
		return errors.New("invalid argument for delete")
	}
	if len(s.Tasks) < 1 {
		return errors.New("no tasks to delete")
	}
	i, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}
	index, err := search(i, s.Tasks)
	if err != nil {
		return err
	}
	s.Tasks = append(s.Tasks[:index], s.Tasks[index+1:]...)
	err = saveState(s)
	if err != nil {
		return err
	}
	return nil
}

func saveState(s *State) error {
	b, err := tasksToJSON(s)
	if err != nil {
		fmt.Printf("Error converting to json: %v", err)
		return err
	}
	err = storeTasks(b)
	if err != nil {
		fmt.Printf("Error storing json: %v", err)
		return err
	}
	return nil
}

func search(num int, lst []*Task) (int, error) {
	i := 0
	j := len(lst)
	for i <= j {
		mid := (i + j) / 2
		if lst[mid].ID > num {
			j = mid - 1
		} else if lst[mid].ID < num {
			i = mid + 1
		} else {
			return mid, nil
		}
	}
	return -1, errors.New("index not found")
}

func handleUpdate(s *State, args []string) error {
	if len(args) != 2 {
		return errors.New("invalid argument for update")
	}
	if len(s.Tasks) == 0 {
		return errors.New("nothing to update")
	}
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}
	newDescr := args[1]
	index, err := search(id, s.Tasks)
	if err != nil {
		return err
	}
	task := s.Tasks[index]
	task.Description = newDescr
	task.UpdatedAt = time.Now()
	err = saveState(s)
	if err != nil {
		return err
	}
	return nil
}

func handleList(s *State, args []string) error {
	if len(args) != 1 {
		return errors.New("invalid args for list")
	}
	status := Status(args[0])
	if status == ToDo || status == InProgress || status == Done {
		for _, task := range s.Tasks {
			if task.Status == status {
				fmt.Printf("ID: %d\nStatus: %s\nDescription: %s\n Created At: %v\nUpdated At: %v\n",
					task.ID, task.Status, task.Description, task.CreatedAt, task.UpdatedAt)
			}
		}
		return nil
	}
	fmt.Println("Unknown status")

	return nil
}

func handleMvInProg(s *State, args []string) error { return nil }
func handleMvDone(s *State, args []string) error   { return nil }

func tasksToJSON(state *State) ([]byte, error) {
	b, err := json.Marshal(state)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// TODO storage.json should be constant
func loadTasks(s *State) error {
	b, err := os.ReadFile("storage.json")
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, s)
	if err != nil {
		return err
	}

	return nil
}

func storeTasks(b []byte) error {
	file, err := os.Create("storage.json")
	if err != nil {
		return err
	}
	bWritten := 0
	for bWritten != len(b) {
		i, err := file.Write(b)
		if err != nil {
			return err
		}
		bWritten += i
	}
	return nil
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
