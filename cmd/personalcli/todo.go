package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var todoCmd = &cobra.Command{
	Use:   "todo",
	Short: "Manage your todo list",
	Run: func(cmd *cobra.Command, args []string) {
		// By default, running "todo" will list tasks.
		listCmd.Run(cmd, args)
	},
}

var addCmd = &cobra.Command{
	Use:   "add [task description]",
	Short: "Add a new task to your todo list",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		tasks, err := readTasks()
		if err != nil {
			fmt.Println("Error reading tasks:", err)
			os.Exit(1)
		}

		newID := 1
		if len(tasks) > 0 {
			newID = tasks[len(tasks)-1].ID + 1
		}

		newTask := Task{
			ID:          newID,
			Description: strings.Join(args, " "),
			Completed:   false,
		}

		tasks = append(tasks, newTask)
		if err := writeTasks(tasks); err != nil {
			fmt.Println("Error writing tasks:", err)
			os.Exit(1)
		}
		fmt.Printf("Added task: \"%s\"\n", newTask.Description)
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all of your tasks",
	Run: func(cmd *cobra.Command, args []string) {
		tasks, err := readTasks()
		if err != nil {
			fmt.Println("Error reading tasks:", err)
			os.Exit(1)
		}

		if len(tasks) == 0 {
			fmt.Println("You have no tasks! Add one with 'personalcli todo add \"my task\"'")
			return
		}

		fmt.Println("Your tasks:")
		for _, task := range tasks {
			status := " "
			if task.Completed {
				status = "✔"
			}
			fmt.Printf("[%s] %d: %s\n", status, task.ID, task.Description)
		}
	},
}

var doneCmd = &cobra.Command{
	Use:   "done [task_id]",
	Short: "Mark a task as completed",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		taskID, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Invalid task ID. Please provide a number.")
			os.Exit(1)
		}

		tasks, err := readTasks()
		if err != nil {
			fmt.Println("Error reading tasks:", err)
			os.Exit(1)
		}

		taskFound := false
		for i := range tasks {
			if tasks[i].ID == taskID {
				tasks[i].Completed = true
				taskFound = true
				break
			}
		}

		if !taskFound {
			fmt.Println("Task ID not found.")
			os.Exit(1)
		}

		if err := writeTasks(tasks); err != nil {
			fmt.Println("Error writing tasks:", err)
			os.Exit(1)
		}

		fmt.Printf("Marked task %d as completed.\n", taskID)
	},
}

func init() {
	// Add a command to clear all tasks
	var clearCmd = &cobra.Command{
		Use:   "clear",
		Short: "Clear all tasks from the list",
		Run: func(cmd *cobra.Command, args []string) {
			if err := writeTasks([]Task{}); err != nil {
				fmt.Println("Error clearing tasks:", err)
				os.Exit(1)
			}
			fmt.Println("All tasks cleared.")
		},
	}

	rootCmd.AddCommand(todoCmd)
	todoCmd.AddCommand(addCmd)
	todoCmd.AddCommand(listCmd)
	todoCmd.AddCommand(doneCmd)
	todoCmd.AddCommand(clearCmd)
}
