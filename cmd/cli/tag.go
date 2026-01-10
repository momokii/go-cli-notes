package main

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var tagCmd = &cobra.Command{
	Use:   "tag",
	Short: "Manage tags",
}

// tagListCmd lists all tags
var tagListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tags",
	RunE: func(cmd *cobra.Command, args []string) error {
		tags, err := apiClient.GetTags()
		if err != nil {
			return fmt.Errorf("list tags: %w", err)
		}

		if len(tags) == 0 {
			fmt.Println("No tags found")
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

// tagCreateCmd creates a new tag
var tagCreateCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a new tag",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		tag, err := apiClient.CreateTag(name)
		if err != nil {
			return fmt.Errorf("create tag: %w", err)
		}

		fmt.Printf("Tag created successfully!\n")
		fmt.Printf("ID: %s\n", tag.ID)
		fmt.Printf("Name: %s\n", tag.Name)

		return nil
	},
}

// tagUpdateCmd updates an existing tag
var tagUpdateCmd = &cobra.Command{
	Use:   "update <id> <new-name>",
	Short: "Update a tag name",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := uuid.Parse(args[0])
		if err != nil {
			return fmt.Errorf("invalid tag ID: %w", err)
		}

		newName := args[1]

		tag, err := apiClient.UpdateTag(id, newName)
		if err != nil {
			return fmt.Errorf("update tag: %w", err)
		}

		fmt.Printf("Tag updated successfully!\n")
		fmt.Printf("ID: %s\n", tag.ID)
		fmt.Printf("New Name: %s\n", tag.Name)

		return nil
	},
}

// tagDeleteCmd deletes a tag
var tagDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a tag",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := uuid.Parse(args[0])
		if err != nil {
			return fmt.Errorf("invalid tag ID: %w", err)
		}

		// Confirm deletion
		fmt.Printf("Are you sure you want to delete this tag? (y/N): ")
		var confirm string
		fmt.Scanln(&confirm)

		if strings.ToLower(confirm) != "y" {
			fmt.Println("Deletion cancelled")
			return nil
		}

		if err := apiClient.DeleteTag(id); err != nil {
			return fmt.Errorf("delete tag: %w", err)
		}

		fmt.Println("Tag deleted successfully!")

		return nil
	},
}

// tagAddCmd adds a tag to a note
var tagAddCmd = &cobra.Command{
	Use:   "add <tag-id-or-name> <note-id>",
	Short: "Add a tag to a note",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		tagIdentifier := args[0]
		noteID, err := uuid.Parse(args[1])
		if err != nil {
			return fmt.Errorf("invalid note ID: %w", err)
		}

		// Try to parse as UUID first, if that fails, search by name
		tagID, err := uuid.Parse(tagIdentifier)
		if err != nil {
			// Not a UUID, search for tag by name
			tags, err := apiClient.GetTags()
			if err != nil {
				return fmt.Errorf("get tags: %w", err)
			}

			var foundTag *uuid.UUID
			for _, tag := range tags {
				if strings.EqualFold(tag.Name, tagIdentifier) {
					foundTag = &tag.ID
					break
				}
			}

			if foundTag == nil {
				return fmt.Errorf("tag '%s' not found", tagIdentifier)
			}
			tagID = *foundTag
		}

		if err := apiClient.AddTagToNote(noteID, tagID); err != nil {
			return fmt.Errorf("add tag to note: %w", err)
		}

		fmt.Println("Tag added to note successfully!")

		return nil
	},
}

// tagRemoveCmd removes a tag from a note
var tagRemoveCmd = &cobra.Command{
	Use:   "remove <tag-id-or-name> <note-id>",
	Short: "Remove a tag from a note",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		tagIdentifier := args[0]
		noteID, err := uuid.Parse(args[1])
		if err != nil {
			return fmt.Errorf("invalid note ID: %w", err)
		}

		// Try to parse as UUID first, if that fails, search by name
		tagID, err := uuid.Parse(tagIdentifier)
		if err != nil {
			// Not a UUID, search for tag by name
			tags, err := apiClient.GetTags()
			if err != nil {
				return fmt.Errorf("get tags: %w", err)
			}

			var foundTag *uuid.UUID
			for _, tag := range tags {
				if strings.EqualFold(tag.Name, tagIdentifier) {
					foundTag = &tag.ID
					break
				}
			}

			if foundTag == nil {
				return fmt.Errorf("tag '%s' not found", tagIdentifier)
			}
			tagID = *foundTag
		}

		if err := apiClient.RemoveTagFromNote(noteID, tagID); err != nil {
			return fmt.Errorf("remove tag from note: %w", err)
		}

		fmt.Println("Tag removed from note successfully!")

		return nil
	},
}

func init() {
	tagCmd.AddCommand(tagListCmd)
	tagCmd.AddCommand(tagCreateCmd)
	tagCmd.AddCommand(tagUpdateCmd)
	tagCmd.AddCommand(tagDeleteCmd)
	tagCmd.AddCommand(tagAddCmd)
	tagCmd.AddCommand(tagRemoveCmd)
	rootCmd.AddCommand(tagCmd)
}
