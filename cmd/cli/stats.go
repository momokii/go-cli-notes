package main

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show user statistics",
	RunE: func(cmd *cobra.Command, args []string) error {
		if !authState.IsAuthenticated() {
			return fmt.Errorf("not authenticated. Please run 'kg-cli login' first")
		}

		stats, err := apiClient.GetStats()
		if err != nil {
			return fmt.Errorf("get stats: %w", err)
		}

		fmt.Println("Knowledge Garden Statistics")
		fmt.Println("==========================")
		fmt.Printf("Total Notes: %d\n", stats.TotalNotes)
		fmt.Printf("Total Tags: %d\n", stats.TotalTags)
		fmt.Printf("Total Links: %d\n", stats.TotalLinks)
		fmt.Printf("Total Words: %d\n", stats.TotalWords)
		fmt.Printf("Notes Created Today: %d\n", stats.NotesCreatedToday)
		fmt.Printf("Notes Created This Week: %d\n", stats.NotesCreatedWeek)

		if stats.LastActivity != nil {
			fmt.Printf("Last Activity: %s\n", stats.LastActivity.Format(time.RFC1123))
		} else {
			fmt.Println("Last Activity: Never")
		}

		return nil
	},
}

var activityCmd = &cobra.Command{
	Use:   "activity",
	Short: "Show recent activity",
	RunE: func(cmd *cobra.Command, args []string) error {
		if !authState.IsAuthenticated() {
			return fmt.Errorf("not authenticated. Please run 'kg-cli login' first")
		}

		limit, _ := cmd.Flags().GetInt("limit")

		activities, err := apiClient.GetRecentActivity(limit)
		if err != nil {
			return fmt.Errorf("get activity: %w", err)
		}

		if len(activities) == 0 {
			fmt.Println("No recent activity")
			return nil
		}

		fmt.Printf("Recent Activity (last %d):\n\n", limit)
		for _, activity := range activities {
			fmt.Printf("%s: %s", activity.CreatedAt.Format("2006-01-02 15:04"), activity.Action)
			if activity.NoteID != nil {
				fmt.Printf(" (%s)", *activity.NoteID)
			}
			fmt.Println()
		}

		return nil
	},
}

var trendingCmd = &cobra.Command{
	Use:   "trending",
	Short: "Show trending notes",
	RunE: func(cmd *cobra.Command, args []string) error {
		if !authState.IsAuthenticated() {
			return fmt.Errorf("not authenticated. Please run 'kg-cli login' first")
		}

		limit, _ := cmd.Flags().GetInt("limit")

		notes, err := apiClient.GetTrendingNotes(limit)
		if err != nil {
			return fmt.Errorf("get trending notes: %w", err)
		}

		if len(notes) == 0 {
			fmt.Println("No trending notes")
			return nil
		}

		fmt.Printf("Trending Notes (top %d):\n\n", limit)
		for _, t := range notes {
			fmt.Printf("Title: %s\n", t.Note.Title)
			fmt.Printf("Access Count: %d\n", t.AccessCount)
			fmt.Println("---")
		}

		return nil
	},
}

func init() {
	activityCmd.Flags().IntP("limit", "l", 10, "Number of activities to show")
	trendingCmd.Flags().IntP("limit", "l", 5, "Number of trending notes to show")

	rootCmd.AddCommand(statsCmd)
	rootCmd.AddCommand(activityCmd)
	rootCmd.AddCommand(trendingCmd)
}
