package main

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"google.golang.org/api/calendar/v3"
)

var calendarCmd = &cobra.Command{
	Use:   "calendar",
	Short: "Interact with your Google Calendar",
	Run: func(cmd *cobra.Command, args []string) {
		// By default, running "calendar" will list upcoming events.
		calendarEventsCmd.Run(cmd, args)
	},
}

var calendarEventsCmd = &cobra.Command{
	Use:   "events",
	Short: "List upcoming events from your primary calendar",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := getClient()
		if err != nil {
			fmt.Printf("Unable to get calendar client: %v\n", err)
			os.Exit(1)
		}

		srv, err := calendar.New(client)
		if err != nil {
			fmt.Printf("Unable to retrieve Calendar client: %v\n", err)
			os.Exit(1)
		}

		t := time.Now().Format(time.RFC3339)
		events, err := srv.Events.List("primary").ShowDeleted(false).
			SingleEvents(true).TimeMin(t).MaxResults(10).OrderBy("startTime").Do()
		if err != nil {
			fmt.Printf("Unable to retrieve next ten of the user's events: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Upcoming events:")
		if len(events.Items) == 0 {
			fmt.Println("No upcoming events found.")
		} else {
			for _, item := range events.Items {
				dateStr := item.Start.DateTime
				if dateStr == "" {
					dateStr = item.Start.Date
				}
				formattedDate := formatEventDate(dateStr)
				fmt.Printf("- %s (%s)\n", item.Summary, formattedDate)
			}
		}
	},
}

func formatEventDate(dateStr string) string {
	// Handle full-day events (date only)
	if len(dateStr) == 10 { // "YYYY-MM-DD"
		t, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			return dateStr // return original on error
		}
		return t.Format("January 2")
	}

	// Handle events with a specific time
	t, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return dateStr // return original on error
	}
	return t.Format("January 2 - 3:04PM")
}

func init() {
	rootCmd.AddCommand(calendarCmd)
	calendarCmd.AddCommand(calendarEventsCmd)
}
