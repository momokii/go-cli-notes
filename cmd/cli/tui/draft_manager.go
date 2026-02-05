package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

// DraftManager handles automatic saving of note drafts
type DraftManager struct {
	draftsDir    string
	saveInterval time.Duration
	lastSave     time.Time
	noteID       uuid.UUID
	title        string
	content      string
	tags         []string
}

// NewDraftManager creates a new draft manager
func NewDraftManager(saveInterval time.Duration) (*DraftManager, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("get home dir: %w", err)
	}

	draftsDir := filepath.Join(homeDir, ".config", "kg-cli", "drafts")

	// Create drafts directory if it doesn't exist
	if err := os.MkdirAll(draftsDir, 0700); err != nil {
		return nil, fmt.Errorf("create drafts dir: %w", err)
	}

	return &DraftManager{
		draftsDir:    draftsDir,
		saveInterval: saveInterval,
		lastSave:     time.Time{},
	}, nil
}

// SetNote sets the current note being edited
func (dm *DraftManager) SetNote(noteID uuid.UUID, title, content string, tags []string) {
	dm.noteID = noteID
	dm.title = title
	dm.content = content
	dm.tags = tags
}

// UpdateContent updates the current note content
func (dm *DraftManager) UpdateContent(title, content string, tags []string) {
	dm.title = title
	dm.content = content
	dm.tags = tags
}

// ShouldSave returns true if enough time has passed since last save
func (dm *DraftManager) ShouldSave() bool {
	if dm.lastSave.IsZero() {
		return true
	}
	return time.Since(dm.lastSave) >= dm.saveInterval
}

// Save saves the current draft to disk
func (dm *DraftManager) Save() error {
	if dm.noteID == uuid.Nil {
		// No note ID yet, don't save
		return nil
	}

	draftPath := dm.getDraftPath(dm.noteID)

	// Create draft content
	draftContent := fmt.Sprintf("Title: %s\nTags: %v\n\n%s\n",
		dm.title,
		dm.tags,
		dm.content)

	// Write draft file
	if err := os.WriteFile(draftPath, []byte(draftContent), 0600); err != nil {
		return fmt.Errorf("write draft: %w", err)
	}

	dm.lastSave = time.Now()
	return nil
}

// LoadDraft loads a draft for the given note ID
// Returns (title, content, tags, found, error)
func (dm *DraftManager) LoadDraft(noteID uuid.UUID) (string, string, []string, bool, error) {
	draftPath := dm.getDraftPath(noteID)

	// Check if draft exists
	if _, err := os.Stat(draftPath); os.IsNotExist(err) {
		return "", "", nil, false, nil
	}

	// Read draft file
	data, err := os.ReadFile(draftPath)
	if err != nil {
		return "", "", nil, false, fmt.Errorf("read draft: %w", err)
	}

	// Parse draft (simple format: Title: xxx\nTags: xxx\n\nContent)
	// For now, just return the raw content as the note content
	title := "Draft"
	content := string(data)
	tags := []string{}

	dm.noteID = noteID
	dm.lastSave = time.Now()

	return title, content, tags, true, nil
}

// ClearDraft removes the draft file for the given note ID
func (dm *DraftManager) ClearDraft(noteID uuid.UUID) error {
	draftPath := dm.getDraftPath(noteID)

	if err := os.Remove(draftPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("remove draft: %w", err)
	}

	if dm.noteID == noteID {
		dm.noteID = uuid.Nil
		dm.lastSave = time.Time{}
	}

	return nil
}

// getDraftPath returns the path to the draft file for a note
func (dm *DraftManager) getDraftPath(noteID uuid.UUID) string {
	return filepath.Join(dm.draftsDir, noteID.String()+".md")
}

// HasDraft returns true if a draft exists for the given note ID
func (dm *DraftManager) HasDraft(noteID uuid.UUID) bool {
	draftPath := dm.getDraftPath(noteID)
	_, err := os.Stat(draftPath)
	return err == nil
}
