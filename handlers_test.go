package main

import (
	"bytes"
	"errors"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	m.Run()
	os.Remove(STORAGEFILE)
}

// Helpers
func stateAndFileAreEqual(s *State) error {
	data, err := os.ReadFile(STORAGEFILE)
	if err != nil {
		return err
	}
	stateData, err := tasksToJSON(s)
	if err != nil {
		return err
	}
	if !bytes.Equal(data, stateData) {
		return errors.New("state in memory and in file are not the same")
	}
	return nil
}

// Tests
func TestAddHandler(t *testing.T) {
	s := &State{}
	err := handleAdd(s, []string{"testAddHandler"})
	if len(s.Tasks) != 1 {
		t.Fatal("Not added")
	}
	if err != nil {
		t.Fatal("handler returned error")
	}
	err = stateAndFileAreEqual(s)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDeleteHandler(t *testing.T) {
	s := &State{}
	err := handleAdd(s, []string{"testAddHandler"})
	if err != nil {
		t.Fatal("create handler returned error")
	}

	err = handleDelete(s, []string{"1"})
	if err != nil {
		t.Fatal("handler returned error")
	}
	if len(s.Tasks) != 0 {
		t.Fatal("task still in memory")
	}
	err = stateAndFileAreEqual(s)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDeleteEmptyList(t *testing.T) {
	s := &State{}
	err := handleDelete(s, []string{"1"})
	if err == nil {
		t.Fatal("expected error but it didn't happen")
	}
}

func TestUpdateHandler(t *testing.T) {
	s := &State{}
	err := handleAdd(s, []string{"testAddHandler"})
	if err != nil {
		t.Fatal("create handler returned error")
	}
	err = handleUpdate(s, []string{"1", "newDescription"})
	if err != nil {
		t.Fatal("handler returned error")
	}
	if s.Tasks[0].Description != "newDescription" {
		t.Fatal("Desscription not updated")
	}
	if s.Tasks[0].UpdatedAt.Equal(s.Tasks[0].CreatedAt) || s.Tasks[0].UpdatedAt.Before(s.Tasks[0].CreatedAt) {
		t.Fatal("Updated time not updated or invalid")
	}
	err = stateAndFileAreEqual(s)
	if err != nil {
		t.Fatal(err)
	}
}

func TestUpdateEmptyList(t *testing.T) {
	s := &State{}
	err := handleUpdate(s, []string{"1", "newDiscription"})
	if err == nil {
		t.Fatal("expected error but it didn't happen")
	}
}

func TestMvToInProg(t *testing.T) {
	s := &State{}
	err := handleAdd(s, []string{"testAddHandler"})
	if err != nil {
		t.Fatal("create handler returned error")
	}
	err = handleMvInProg(s, []string{"1"})
	if err != nil {
		t.Fatal("create handler returned error")
	}
	if s.Tasks[0].Status != Status("in-progress") {
		t.Fatal("state not updated")
	}
	err = stateAndFileAreEqual(s)
	if err != nil {
		t.Fatal(err)
	}
}

func TestMvToInProgEmptyList(t *testing.T) {
	s := &State{}
	err := handleMvInProg(s, []string{"1"})
	if err == nil {
		t.Fatal("expected error but it didn't happen")
	}
}

func TestMvToDone(t *testing.T) {
	s := &State{}
	err := handleAdd(s, []string{"testAddHandler"})
	if err != nil {
		t.Fatal("create handler returned error")
	}
	err = handleMvDone(s, []string{"1"})
	if err != nil {
		t.Fatal("create handler returned error")
	}
	if s.Tasks[0].Status != Status("done") {
		t.Fatal("state not updated")
	}
	err = stateAndFileAreEqual(s)
	if err != nil {
		t.Fatal(err)
	}
}

func TestMvToDoneEmptyList(t *testing.T) {
	s := &State{}
	err := handleMvDone(s, []string{"1"})
	if err == nil {
		t.Fatal("expected error but it didn't happen")
	}
}
