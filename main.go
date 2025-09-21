package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Task struct {
	ID          int `json:"id"`
	Description string `json:"description"`
	Completed   bool `json:"completed"`
}

var taskFile = "tasks.json"

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go [list|add <description>|clear <task_id>]")
		return
	}
	command := os.Args[1]
	switch command {
	case "list":
		listTasks()

	case "add":
		if len(os.Args) < 3 {
			fmt.Println("Usage: go run main.go add <description>")
			return
		}
		description := strings.Join(os.Args[2:], " ")
		addTask(description)

	case "clear":
		if len(os.Args) == 3 {
			taskID := os.Args[2]
			err := clearTask(taskID)
			if err != nil {
				fmt.Println("Error clearing task:", err)
			}	
		} else {
			err := clearCompletedTasks()
			if err != nil {
				fmt.Println("Error clearing completed tasks:", err)
			}
		}
		

	default:
		fmt.Println("Unknown command:", command)
	}
	
}

func listTasks() {
	tasks, err := loadTasks()
	if err != nil {
		fmt.Println("Error loading tasks:", err)
		return
	}
	print(tasks)
}

func print(tasks []Task) {
	for _, task := range tasks {
		status := " "
		if task.Completed {
			status = "x"
		}
		fmt.Printf("[%s] %d: %s\n", status, task.ID, task.Description)
	}
}

func addTask(description string) {
	tasks, err := loadTasks()
	if err != nil {
		fmt.Println("Error loading tasks:", err)
		return
	}
	newID := len(tasks) + 1
	newTask := Task{ID: newID, Description: description, Completed: false}
	tasks = append(tasks, newTask)
	err = saveTasks(tasks)
	if err != nil {
		fmt.Println("Error saving tasks:", err)
	}
	print(tasks)
}

func loadTasks() (tasks []Task, err error) {
	data, err := os.ReadFile(taskFile)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return nil, err
	}

	if len(data) == 0 {
		return []Task{}, nil
	}

	err = json.Unmarshal(data, &tasks)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return nil, err
	}
	return tasks, nil
}

func saveTasks(tasks []Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return err
	}	
	err = os.WriteFile(taskFile, data, 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return err
	}
	return nil
}

func clearTask(taskID string) error {
	tasks, err := loadTasks()
	if err != nil {
		fmt.Println("Error loading tasks:", err)
		return err
	}	
	found := false
	for i, task := range tasks {
		if strconv.Itoa(task.ID) == taskID {
			tasks[i].Completed = true
			found = true
			break
		}
	}
	if !found {
		fmt.Println("Task ID not found:", taskID)
		return nil
	}
	err = saveTasks(tasks)
	if err != nil {
		fmt.Println("Error saving tasks:", err)
		return err
	}
	print(tasks)
	return nil
}

func clearCompletedTasks() error {
	tasks, err := loadTasks()
	if err != nil {
		fmt.Println("Error loading tasks:", err)
		return err
	}
	var activeTasks []Task
	count := 0
	for _, task := range tasks {
		if !task.Completed {
			task.ID = count + 1
			count++
			activeTasks = append(activeTasks, task)
		}
	}
	err = saveTasks(activeTasks)
	if err != nil {
		fmt.Println("Error saving tasks:", err)
		return err
	}
	print(activeTasks)
	return nil
}