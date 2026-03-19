package main

import (
	"encoding/json"
	"fmt"
	"os"
)

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

func tasksToJSON(state *State) ([]byte, error) {
	b, err := json.Marshal(state)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func loadTasks(s *State) error {
	b, err := os.ReadFile(STORAGEFILE)
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
	file, err := os.Create(STORAGEFILE)
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
