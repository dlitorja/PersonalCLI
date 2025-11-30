package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Task represents a single todo item.
type Task struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

// tasksFilePath is the path to the JSON file where tasks are stored.
var tasksFilePath string

func init() {
	// Get home directory
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error finding home directory:", err)
		os.Exit(1)
	}
	
	// Define the path for the tasks file
	configDir := filepath.Join(home, ".config", "personalcli")
	tasksFilePath = filepath.Join(configDir, "tasks.json")

	// Create the directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0755); err != nil {
		fmt.Println("Error creating config directory:", err)
		os.Exit(1)
	}
}

// readTasks reads all tasks from the tasks.json file.
func readTasks() ([]Task, error) {
	data, err := os.ReadFile(tasksFilePath)
	if err != nil {
		// If the file doesn't exist, return an empty list of tasks.
		if os.IsNotExist(err) {
			return []Task{}, nil
		}
		return nil, err
	}

	var tasks []Task
	err = json.Unmarshal(data, &tasks)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

// writeTasks writes a list of tasks to the tasks.json file.
func writeTasks(tasks []Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(tasksFilePath, data, 0644)
}
