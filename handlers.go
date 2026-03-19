package main

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

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
				fmt.Printf("ID: %d\nStatus: %s\nDescription: %s\nCreated At: %v\nUpdated At: %v\n",
					task.ID, task.Status, task.Description, task.CreatedAt, task.UpdatedAt)
			}
		}
		return nil
	}
	fmt.Println("Unknown status")

	return nil
}

func handleMvInProg(s *State, args []string) error {
	if len(args) != 1 {
		return errors.New("invalid args for mvInProg")
	}
	if len(s.Tasks) <= 0 {
		return errors.New("empty tasks list")
	}
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}
	index, err := search(id, s.Tasks)
	if err != nil {
		return err
	}
	task := s.Tasks[index]
	task.Status = InProgress
	task.UpdatedAt = time.Now()
	err = saveState(s)
	if err != nil {
		return err
	}
	return nil
}

func handleMvDone(s *State, args []string) error {
	if len(args) != 1 {
		return errors.New("invalid args for mvInProg")
	}
	if len(s.Tasks) <= 0 {
		return errors.New("empty tasks list")
	}
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}
	index, err := search(id, s.Tasks)
	if err != nil {
		return err
	}
	task := s.Tasks[index]
	task.Status = Done
	task.UpdatedAt = time.Now()
	err = saveState(s)
	if err != nil {
		return err
	}
	return nil
}
