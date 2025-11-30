package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var notesCmd = &cobra.Command{
	Use:   "note",
	Short: "Manage your notes",
	Run: func(cmd *cobra.Command, args []string) {
		// By default, running "note" will list notes.
		noteListCmd.Run(cmd, args)
	},
}

var noteNewCmd = &cobra.Command{
	Use:   "new [note content]",
	Short: "Create a new note",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		notes, err := readNotes()
		if err != nil {
			fmt.Println("Error reading notes:", err)
			os.Exit(1)
		}

		newID := 1
		if len(notes) > 0 {
			newID = notes[len(notes)-1].ID + 1
		}

		newNote := Note{
			ID:        newID,
			Content:   strings.Join(args, " "),
			CreatedAt: time.Now(),
		}

		notes = append(notes, newNote)
		if err := writeNotes(notes); err != nil {
			fmt.Println("Error writing notes:", err)
			os.Exit(1)
		}
		fmt.Printf("Created note %d.\n", newNote.ID)
	},
}

var noteListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all of your notes",
	Run: func(cmd *cobra.Command, args []string) {
		notes, err := readNotes()
		if err != nil {
			fmt.Println("Error reading notes:", err)
			os.Exit(1)
		}

		if len(notes) == 0 {
			fmt.Println("You have no notes! Add one with 'personalcli note new \"my note\"'")
			return
		}

		fmt.Println("Your notes:")
		for _, note := range notes {
			fmt.Printf("ID: %d | Date: %s\n%s\n---\n", note.ID, note.CreatedAt.Format("2006-01-02 15:04"), note.Content)
		}
	},
}

var noteFindCmd = &cobra.Command{
	Use:   "find [keyword]",
	Short: "Find notes containing a keyword",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		keyword := strings.ToLower(args[0])
		notes, err := readNotes()
		if err != nil {
			fmt.Println("Error reading notes:", err)
			os.Exit(1)
		}

		fmt.Printf("Searching for notes with keyword: \"%s\"\n", keyword)
		found := false
		for _, note := range notes {
			if strings.Contains(strings.ToLower(note.Content), keyword) {
				fmt.Printf("ID: %d | Date: %s\n%s\n---\n", note.ID, note.CreatedAt.Format("2006-01-02 15:04"), note.Content)
				found = true
			}
		}

		if !found {
			fmt.Println("No matching notes found.")
		}
	},
}

func init() {
	rootCmd.AddCommand(notesCmd)
	notesCmd.AddCommand(noteNewCmd)
	notesCmd.AddCommand(noteListCmd)
	notesCmd.AddCommand(noteFindCmd)
}
