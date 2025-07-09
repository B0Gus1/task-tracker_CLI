package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

type ToDoList struct {
	Tasks map[int]Task `json:"tasks"`
}

type Task struct {
	Description string     `json:"description"`
	Status      TaskStatus `json:"status"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}

type TaskStatus int

const (
	Todo TaskStatus = iota + 1
	InProgress
	Done
)

func (t *ToDoList) Add(description string) {
	var newID int
	for i := 1; ; i++ {
		if _, exists := t.Tasks[i]; !exists {
			newID = i
			break
		}
	}
	t.Tasks[newID] = Task{Description: description, Status: Todo, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	fmt.Printf("Task added successfully (ID: %d)\n", newID)
}

func (t *ToDoList) Update(id int, description string) {
	t.Tasks[id] = Task{Description: description, Status: t.Tasks[id].Status, CreatedAt: t.Tasks[id].CreatedAt, UpdatedAt: time.Now()}
}

func (t *ToDoList) Delete(id int) {
	delete(t.Tasks, id)
}

func (t *ToDoList) MarkInProgress(id int) {
	t.Tasks[id] = Task{Description: t.Tasks[id].Description, Status: InProgress, CreatedAt: t.Tasks[id].CreatedAt, UpdatedAt: time.Now()}
}

func (t *ToDoList) MarkDone(id int) {
	t.Tasks[id] = Task{Description: t.Tasks[id].Description, Status: Done, CreatedAt: t.Tasks[id].CreatedAt, UpdatedAt: time.Now()}
}

func (t *ToDoList) List() {
	for id := range t.Tasks {
		fmt.Println(t.Tasks[id].Description)
	}
}

func (t *ToDoList) ListDone() {
	for id := range t.Tasks {
		if t.Tasks[id].Status == Done {
			fmt.Println(t.Tasks[id].Description)
		}
	}
}

func (t *ToDoList) ListToDo() {
	for id := range t.Tasks {
		if t.Tasks[id].Status == Todo {
			fmt.Println(t.Tasks[id].Description)
		}
	}
}

func (t *ToDoList) ListInProgress() {
	for id := range t.Tasks {
		if t.Tasks[id].Status == InProgress {
			fmt.Println(t.Tasks[id].Description)
		}
	}
}

func main() {
	filename := "todolist.json"
	var toDoList ToDoList

	buf, err := os.ReadFile(filename)
	if err != nil {
		if !os.IsNotExist(err) {
			log.Fatalf("Error reading file: %v\n", err)
		}
		toDoList = ToDoList{
			Tasks: make(map[int]Task),
		}
	} else {
		if len(buf) == 0 {
			toDoList = ToDoList{
				Tasks: make(map[int]Task),
			}
		} else {
			err = json.Unmarshal(buf, &toDoList)
			if err != nil {
				log.Fatalf("Error parsing json: %v\n", err)
			}
		}
	}

	if len(os.Args) == 1 {
		log.Fatal("Write a command\n")
	}

	switch os.Args[1] {
	case "add":
		if len(os.Args) != 3 {
			log.Fatalf("add: Invalid number of arguments: %d, expected: 1\n", len(os.Args)-2)
		}
		toDoList.Add(os.Args[2])
	case "update":
		if len(os.Args) != 4 {
			log.Fatalf("update: Invalid number of arguments: %d, expected: 2\n", len(os.Args)-2)
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatalf("update: Invalid id: %s\n", os.Args[2])
		}
		if _, exists := toDoList.Tasks[id]; !exists {
			log.Fatalf("update: Id does not exist: %d\n", id)
		}
		toDoList.Update(id, os.Args[3])
	case "delete":
		if len(os.Args) != 3 {
			log.Fatalf("delete: Invalid number of arguments: %d, expected: 1\n", len(os.Args)-2)
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatalf("delete: Invalid id: %s\n", os.Args[2])
		}
		if _, exists := toDoList.Tasks[id]; !exists {
			log.Fatalf("delete: Id does not exist: %d\n", id)
		}
		toDoList.Delete(id)
	case "mark-in-progress":
		if len(os.Args) != 3 {
			log.Fatalf("mark-in-progress: Invalid number of arguments: %d, expected: 1\n", len(os.Args)-2)
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatalf("mark-in-progress: Invalid id: %s\n", os.Args[2])
		}
		if _, exists := toDoList.Tasks[id]; !exists {
			log.Fatalf("mark-in-progress: Id does not exist: %d\n", id)
		}
		toDoList.MarkInProgress(id)
	case "mark-done":
		if len(os.Args) != 3 {
			log.Fatalf("mark-done: Invalid number of arguments: %d, expected: 1\n", len(os.Args)-2)
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatalf("mark-done: Invalid id: %s\n", os.Args[2])
		}
		if _, exists := toDoList.Tasks[id]; !exists {
			log.Fatalf("mark-done: Id does not exist: %d\n", id)
		}
		toDoList.MarkDone(id)
	case "list":
		switch {
		case len(os.Args) == 2:
			toDoList.List()
		case len(os.Args) > 3:
			log.Fatalf("Invalid number of arguments: %d, expected: 1\n", len(os.Args)-2)
		default:
			switch os.Args[2] {
			case "done":
				toDoList.ListDone()
			case "todo":
				toDoList.ListToDo()
			case "in-progress":
				toDoList.ListInProgress()
			default:
				log.Fatalf("Invalid argument: %s\n", os.Args[2])
			}
		}
	default:
		log.Fatalf("Unknown command: %s\n", os.Args[1])
	}

	buf, err = json.Marshal(toDoList)
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	err = os.WriteFile(filename, buf, 0644)
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	return
}
