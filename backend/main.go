package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

var taskFile = "../tasks.json"

func main() {
	var mode = flag.String("mode", "server", "Mode to run: cli or server")
	flag.Parse()
	if *mode == "cli" {
		cli()
	} else {
		startServer()
	}
}

func cli() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go [list|add <description>|clear <task_id>]")
		return
	}
	command := os.Args[1]
	switch command {
	case "list":
		tasks, err := fetchTasks()
		if err != nil {
			fmt.Println("Error loading tasks:", err)
			return
		}
		print(tasks)
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("Usage: go run main.go add <description>")
			return
		}
		description := strings.Join(os.Args[2:], " ")
		tasks := addTask(description)
		print(tasks)

	case "clear":
		if len(os.Args) == 3 {
			taskID := os.Args[2]
			tasks, err := clearTask(taskID)
			if err != nil {
				fmt.Println("Error clearing task:", err)
			}	
			print(tasks)
		} else {
			tasks, err := clearCompletedTasks()
			if err != nil {
				fmt.Println("Error clearing completed tasks:", err)
			}
			print(tasks)
		}
	default:
		fmt.Println("Unknown command:", command)
	}
}

func startServer() {
	fmt.Println("Starting server on :8080")
}
	