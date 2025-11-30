package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Note represents a single note item.
type Note struct {
	ID        int       `json:"id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// notesFilePath is the path to the JSON file where notes are stored.
var notesFilePath string

func init() {
	// Get home directory
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error finding home directory:", err)
		os.Exit(1)
	}

	// Define the path for the notes file
	configDir := filepath.Join(home, ".config", "personalcli")
	notesFilePath = filepath.Join(configDir, "notes.json")

	// Create the directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0755); err != nil {
		fmt.Println("Error creating config directory:", err)
		os.Exit(1)
	}
}

// readNotes reads all notes from the notes.json file.
func readNotes() ([]Note, error) {
	data, err := os.ReadFile(notesFilePath)
	if err != nil {
		// If the file doesn't exist, return an empty list of notes.
		if os.IsNotExist(err) {
			return []Note{}, nil
		}
		return nil, err
	}

	var notes []Note
	err = json.Unmarshal(data, &notes)
	if err != nil {
		return nil, err
	}
	return notes, nil
}

// writeNotes writes a list of notes to the notes.json file.
func writeNotes(notes []Note) error {
	data, err := json.MarshalIndent(notes, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(notesFilePath, data, 0644)
}
