package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

func print(tasks []Task) {
	if len(tasks) == 0 {
		fmt.Println("No tasks found.")
		return
	}
	for _, task := range tasks {
		status := " "
		if task.Completed {
			status = "x"
		}
		fmt.Printf("[%s] %d: %s\n", status, task.ID, task.Description)
	}
}

func addTask(description string) []Task{
	tasks, err := fetchTasks()
	if err != nil {
		fmt.Println("Error loading tasks:", err)
		return nil
	}
	newID := len(tasks) + 1
	newTask := Task{ID: newID, Description: description, Completed: false}
	tasks = append(tasks, newTask)
	err = saveTasks(tasks)
	if err != nil {
		fmt.Println("Error saving tasks:", err)
		return nil
	}

	return tasks
}

func fetchTasks() (tasks []Task, err error) {
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

func updateTaskStatus(taskID string, isCompleted bool) ([]Task, error) {
	tasks, err := fetchTasks()
	if err != nil {
		fmt.Println("Error loading tasks:", err)
		return nil, err
	}	
	found := false
	for i, task := range tasks {
		if strconv.Itoa(task.ID) == taskID {
			tasks[i].Completed = isCompleted
			found = true
			break
		}
	}
	if !found {
		fmt.Println("Task ID not found:", taskID)
		return tasks, nil
	}
	return tasks, nil
}

func clearTask(taskID string) ([]Task, error) {
	tasks, err := updateTaskStatus(taskID, true)
	if err != nil {
		fmt.Println("Error updating task status:", err)
		return nil, err
	}

	err = saveTasks(tasks)
	if err != nil {
		fmt.Println("Error saving tasks:", err)
		return nil, err
	}
	return tasks, nil
}

func clearCompletedTasks() ([]Task, error) {
	tasks, err := fetchTasks()
	if err != nil {
		fmt.Println("Error loading tasks:", err)
		return nil, err
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
		return nil, err
	}
	return activeTasks, nil
}