package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/momokii/go-cli-notes/internal/model"
	"github.com/spf13/cobra"
)

var noteCmd = &cobra.Command{
	Use:   "note",
	Short: "Manage notes",
}

// noteListCmd lists all notes
var noteListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all notes",
	RunE: func(cmd *cobra.Command, args []string) error {
		page, _ := cmd.Flags().GetInt("page")
		limit, _ := cmd.Flags().GetInt("limit")
		search, _ := cmd.Flags().GetString("search")
		tag, _ := cmd.Flags().GetString("tag")

		filter := model.NoteFilter{
			Page:   page,
			Limit:  limit,
			Search: search,
		}

		// Handle tag filtering - support both tag ID and tag name
		if tag != "" {
			// Try to parse as UUID first
			if _, err := uuid.Parse(tag); err == nil {
				filter.TagID = &tag
			} else {
				// Not a UUID, search for tag by name
				tags, err := apiClient.GetTags()
				if err != nil {
					return fmt.Errorf("get tags: %w", err)
				}

				var foundTagID *string
				for _, t := range tags {
					if strings.EqualFold(t.Name, tag) {
						tagID := t.ID.String()
						foundTagID = &tagID
						break
					}
				}

				if foundTagID == nil {
					return fmt.Errorf("tag '%s' not found", tag)
				}
				filter.TagID = foundTagID
			}
		}

		notes, total, err := apiClient.ListNotes(filter)
		if err != nil {
			return fmt.Errorf("list notes: %w", err)
		}

		if total == 0 {
			fmt.Println("No notes found")
			return nil
		}

		fmt.Printf("Found %d note(s):\n\n", total)
		for _, note := range notes {
			fmt.Printf("ID: %s\n", note.ID)
			fmt.Printf("Title: %s\n", note.Title)
			fmt.Printf("Type: %s\n", note.NoteType)
			fmt.Printf("Words: %d\n", note.WordCount)
			fmt.Printf("Created: %s\n", note.CreatedAt.Format("2006-01-02 15:04"))
			fmt.Println("---")
		}

		return nil
	},
}

// noteGetCmd gets a single note
var noteGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get a note by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := uuid.Parse(args[0])
		if err != nil {
			return fmt.Errorf("invalid note ID: %w", err)
		}

		note, err := apiClient.GetNote(id)
		if err != nil {
			return fmt.Errorf("get note: %w", err)
		}

		fmt.Printf("Title: %s\n", note.Title)
		fmt.Printf("Type: %s\n", note.NoteType)
		fmt.Printf("Word Count: %d\n", note.WordCount)
		fmt.Printf("Reading Time: %d min\n", note.ReadingTimeMinutes)
		fmt.Printf("Created: %s\n", note.CreatedAt.Format("2006-01-02 15:04:05"))
		fmt.Printf("Updated: %s\n", note.UpdatedAt.Format("2006-01-02 15:04:05"))
		fmt.Println("\nContent:")
		fmt.Println("---")
		fmt.Println(note.Content)

		return nil
	},
}

// noteCreateCmd creates a new note
var noteCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new note",
	RunE: func(cmd *cobra.Command, args []string) error {
		title, _ := cmd.Flags().GetString("title")
		content, _ := cmd.Flags().GetString("content")
		noteType, _ := cmd.Flags().GetString("type")

		if title == "" {
			return fmt.Errorf("title is required (use --title flag)")
		}

		noteTypeEnum := model.NoteType(noteType)
		req := &model.CreateNoteRequest{
			Title:    title,
			Content:  content,
			NoteType: noteTypeEnum,
		}

		note, err := apiClient.CreateNote(req)
		if err != nil {
			return fmt.Errorf("create note: %w", err)
		}

		fmt.Printf("Note created successfully!\n")
		fmt.Printf("ID: %s\n", note.ID)
		fmt.Printf("Title: %s\n", note.Title)

		return nil
	},
}

// noteSearchCmd searches notes
var noteSearchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search notes",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := args[0]
		page, _ := cmd.Flags().GetInt("page")
		limit, _ := cmd.Flags().GetInt("limit")

		result, err := apiClient.SearchNotes(query, page, limit)
		if err != nil {
			return fmt.Errorf("search notes: %w", err)
		}

		if len(result.Results) == 0 {
			fmt.Println("No results found")
			return nil
		}

		fmt.Printf("Found %d result(s) for '%s':\n\n", len(result.Results), query)
		for _, r := range result.Results {
			fmt.Printf("ID: %s\n", r.Note.ID)
			fmt.Printf("Title: %s\n", r.Note.Title)
			fmt.Printf("Snippet: %s\n", r.Snippet)
			fmt.Println("---")
		}

		return nil
	},
}

// noteDailyCmd gets or creates a daily note
var noteDailyCmd = &cobra.Command{
	Use:   "daily [date]",
	Short: "Get or create a daily note (YYYY-MM-DD format, or today if omitted)",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		date := "today"
		if len(args) > 0 {
			date = args[0]
		}

		note, isCreated, err := apiClient.GetDailyNote(date)
		if err != nil {
			return fmt.Errorf("get daily note: %w", err)
		}

		if isCreated {
			fmt.Printf("Created new daily note for %s\n", date)
		} else {
			fmt.Printf("Found existing daily note for %s\n", date)
		}

		fmt.Printf("ID: %s\n", note.ID)
		fmt.Printf("Title: %s\n", note.Title)
		fmt.Printf("Content:\n%s\n", note.Content)

		return nil
	},
}

// noteUpdateCmd updates an existing note
var noteUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update a note (interactive by default, or use flags for automation)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := uuid.Parse(args[0])
		if err != nil {
			return fmt.Errorf("invalid note ID: %w", err)
		}

		title, _ := cmd.Flags().GetString("title")
		content, _ := cmd.Flags().GetString("content")

		// If flags provided, use flag-based update (for automation)
		if title != "" || content != "" {
			req := &model.UpdateNoteRequest{}
			if title != "" {
				req.Title = &title
			}
			if content != "" {
				req.Content = &content
			}

			if err := apiClient.UpdateNote(id, req); err != nil {
				return fmt.Errorf("update note: %w", err)
			}

			fmt.Println("Note updated successfully!")
			return nil
		}

		// Interactive mode (default when no flags provided)
		// Get current note first
		note, err := apiClient.GetNote(id)
		if err != nil {
			return fmt.Errorf("get note: %w", err)
		}

		fmt.Printf("Updating note: %s\n", note.Title)
		fmt.Println("Current values shown - leave empty to keep existing value")
		fmt.Println()

		// Prompt for title update
		fmt.Printf("Current title: %s\n", note.Title)
		fmt.Print("New title (press Enter to keep current): ")
		reader := bufio.NewReader(os.Stdin)
		newTitle, _ := reader.ReadString('\n')
		newTitle = strings.TrimSpace(newTitle)

		// Prompt for content update using editor
		fmt.Print("\nEdit content? (Y/n): ")
		confirmEdit, _ := reader.ReadString('\n')
		confirmEdit = strings.TrimSpace(strings.ToLower(confirmEdit))

		var newContent string
		if confirmEdit != "n" {
			// Create temp file with current content
			tmpFile := filepath.Join(os.TempDir(), fmt.Sprintf("kg-cli-note-%s.md", id.String()))
			if err := os.WriteFile(tmpFile, []byte(note.Content), 0644); err != nil {
				return fmt.Errorf("create temp file: %w", err)
			}
			defer os.Remove(tmpFile)

			// Get editor from environment or use default
			editor := os.Getenv("EDITOR")
			if editor == "" {
				editor = "vi" // Default to vi
			}

			// Open editor
			editorCmd := exec.Command(editor, tmpFile)
			editorCmd.Stdin = os.Stdin
			editorCmd.Stdout = os.Stdout
			editorCmd.Stderr = os.Stderr

			if err := editorCmd.Run(); err != nil {
				return fmt.Errorf("editor failed: %w", err)
			}

			// Read edited content
			contentBytes, err := os.ReadFile(tmpFile)
			if err != nil {
				return fmt.Errorf("read edited content: %w", err)
			}
			newContent = string(contentBytes)
		} else {
			newContent = note.Content
		}

		// Build request
		req := &model.UpdateNoteRequest{}
		if newTitle != "" && newTitle != note.Title {
			req.Title = &newTitle
		}
		if newContent != note.Content {
			req.Content = &newContent
		}

		// Check if anything changed
		if req.Title == nil && req.Content == nil {
			fmt.Println("\nNo changes made.")
			return nil
		}

		// Confirm update
		fmt.Print("\nSave changes? (Y/n): ")
		confirmSave, _ := reader.ReadString('\n')
		confirmSave = strings.TrimSpace(strings.ToLower(confirmSave))

		if confirmSave == "n" {
			fmt.Println("Update cancelled.")
			return nil
		}

		// Perform update
		if err := apiClient.UpdateNote(id, req); err != nil {
			return fmt.Errorf("update note: %w", err)
		}

		fmt.Println("Note updated successfully!")
		return nil
	},
}

// noteDeleteCmd deletes a note
var noteDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a note",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := uuid.Parse(args[0])
		if err != nil {
			return fmt.Errorf("invalid note ID: %w", err)
		}

		// Confirm deletion
		fmt.Printf("Are you sure you want to delete this note? (y/N): ")
		var confirm string
		fmt.Scanln(&confirm)

		if confirm != "y" && confirm != "Y" {
			fmt.Println("Deletion cancelled")
			return nil
		}

		if err := apiClient.DeleteNote(id); err != nil {
			return fmt.Errorf("delete note: %w", err)
		}

		fmt.Println("Note deleted successfully!")
		return nil
	},
}

// noteLinksCmd shows outgoing links from a note
var noteLinksCmd = &cobra.Command{
	Use:   "links <id>",
	Short: "Show outgoing links from a note",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("requires note ID, e.g., kg-cli note links 123e4567-e89b-12d3-a456-426614174000")
		}
		id, err := uuid.Parse(args[0])
		if err != nil {
			return fmt.Errorf("invalid note ID format")
		}

		links, err := apiClient.GetLinks(id)
		if err != nil {
			return fmt.Errorf("get links: %w", err)
		}

		if len(links) == 0 {
			fmt.Println("No outgoing links found")
			return nil
		}

		fmt.Printf("Found %d outgoing link(s):\n\n", len(links))
		for _, link := range links {
			if link.TargetNote != nil {
				fmt.Printf("To: %s (ID: %s)\n", link.TargetNote.Title, link.TargetNote.ID)
			} else {
				fmt.Printf("To: (deleted note)\n")
			}
			if link.LinkContext != nil {
				fmt.Printf("Context: %s\n", *link.LinkContext)
			}
			fmt.Printf("Created: %s\n", link.CreatedAt.Format("2006-01-02 15:04"))
			fmt.Println("---")
		}

		return nil
	},
}

// noteBacklinksCmd shows backlinks to a note
var noteBacklinksCmd = &cobra.Command{
	Use:   "backlinks <id>",
	Short: "Show backlinks to a note",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("requires note ID, e.g., kg-cli note backlinks 123e4567-e89b-12d3-a456-426614174000")
		}
		id, err := uuid.Parse(args[0])
		if err != nil {
			return fmt.Errorf("invalid note ID format")
		}

		backlinks, err := apiClient.GetBacklinks(id)
		if err != nil {
			return fmt.Errorf("get backlinks: %w", err)
		}

		if len(backlinks) == 0 {
			fmt.Println("No backlinks found")
			return nil
		}

		fmt.Printf("Found %d backlink(s):\n\n", len(backlinks))
		for _, link := range backlinks {
			if link.SourceNote != nil {
				fmt.Printf("From: %s (ID: %s)\n", link.SourceNote.Title, link.SourceNote.ID)
			} else {
				fmt.Printf("From: (deleted note)\n")
			}
			if link.LinkContext != nil {
				fmt.Printf("Context: %s\n", *link.LinkContext)
			}
			fmt.Printf("Created: %s\n", link.CreatedAt.Format("2006-01-02 15:04"))
			fmt.Println("---")
		}

		return nil
	},
}

// noteTagsCmd shows tags on a note
var noteTagsCmd = &cobra.Command{
	Use:   "tags <id>",
	Short: "Show tags on a note",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("requires note ID, e.g., kg-cli note tags 123e4567-e89b-12d3-a456-426614174000")
		}
		id, err := uuid.Parse(args[0])
		if err != nil {
			return fmt.Errorf("invalid note ID format")
		}

		tags, err := apiClient.GetNoteTags(id)
		if err != nil {
			return fmt.Errorf("get note tags: %w", err)
		}

		if len(tags) == 0 {
			fmt.Println("No tags found on this note")
			return nil
		}

		fmt.Printf("Found %d tag(s):\n\n", len(tags))
		for _, tag := range tags {
			fmt.Printf("ID: %s\n", tag.ID)
			fmt.Printf("Name: %s\n", tag.Name)
			fmt.Println("---")
		}

		return nil
	},
}

func init() {
	// Add flags to noteListCmd
	noteListCmd.Flags().IntP("page", "p", 1, "Page number")
	noteListCmd.Flags().IntP("limit", "l", 20, "Notes per page")
	noteListCmd.Flags().StringP("search", "s", "", "Search query")
	noteListCmd.Flags().StringP("tag", "t", "", "Filter by tag name or ID")

	// Add flags to noteCreateCmd
	noteCreateCmd.Flags().StringP("title", "t", "", "Note title (required)")
	noteCreateCmd.Flags().StringP("content", "c", "", "Note content")
	noteCreateCmd.Flags().StringP("type", "T", "note", "Note type (note, daily, meeting, idea)")

	// Add flags to noteSearchCmd
	noteSearchCmd.Flags().IntP("page", "p", 1, "Page number")
	noteSearchCmd.Flags().IntP("limit", "l", 20, "Results per page")

	// Add flags to noteUpdateCmd
	noteUpdateCmd.Flags().StringP("title", "t", "", "New note title")
	noteUpdateCmd.Flags().StringP("content", "c", "", "New note content")

	// Add subcommands to noteCmd
	noteCmd.AddCommand(noteListCmd)
	noteCmd.AddCommand(noteGetCmd)
	noteCmd.AddCommand(noteCreateCmd)
	noteCmd.AddCommand(noteUpdateCmd)
	noteCmd.AddCommand(noteDeleteCmd)
	noteCmd.AddCommand(noteSearchCmd)
	noteCmd.AddCommand(noteDailyCmd)
	noteCmd.AddCommand(noteLinksCmd)
	noteCmd.AddCommand(noteBacklinksCmd)
	noteCmd.AddCommand(noteTagsCmd)

	// Add noteCmd to rootCmd
	rootCmd.AddCommand(noteCmd)
}
