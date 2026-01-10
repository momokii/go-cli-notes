package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/momokii/go-cli-notes/internal/model"
	"github.com/momokii/go-cli-notes/internal/repository"
	"github.com/momokii/go-cli-notes/internal/util"
)

// NoteService handles note business logic
type NoteService struct {
	noteRepo    repository.NoteRepository
	tagRepo     repository.TagRepository
	linkRepo    repository.LinkRepository
	activityRepo repository.ActivityRepository
	linkParser  *util.LinkParser
}

// NewNoteService creates a new note service
func NewNoteService(
	noteRepo repository.NoteRepository,
	tagRepo repository.TagRepository,
	linkRepo repository.LinkRepository,
	activityRepo repository.ActivityRepository,
	linkParser *util.LinkParser,
) *NoteService {
	return &NoteService{
		noteRepo:    noteRepo,
		tagRepo:     tagRepo,
		linkRepo:    linkRepo,
		activityRepo: activityRepo,
		linkParser:  linkParser,
	}
}

// Create creates a new note
func (s *NoteService) Create(ctx context.Context, userID uuid.UUID, req *model.CreateNoteRequest) (*model.Note, error) {
	// Validate request
	if err := util.ValidateStruct(req); err != nil {
		return nil, fmt.Errorf("%w: %s", model.ErrValidation, err)
	}

	// Set default note type
	noteType := req.NoteType
	if noteType == "" {
		noteType = model.NoteTypeNote
	}

	// Create note
	note := &model.Note{
		UserID:   userID,
		Title:    req.Title,
		Content:  req.Content,
		NoteType: noteType,
		Metadata: make(model.Metadata),
	}

	if err := s.noteRepo.Create(ctx, note); err != nil {
		return nil, fmt.Errorf("create note: %w", err)
	}

	// Extract and create links
	s.processLinks(ctx, userID, note)

	// Log activity
	_ = s.activityRepo.Create(ctx, &model.Activity{
		UserID:  userID,
		NoteID:  &note.ID,
		Action:  model.ActionCreate,
		Metadata: model.ActivityMetadata{
			"title": note.Title,
		},
	})

	return note, nil
}

// GetByID gets a note by ID
func (s *NoteService) GetByID(ctx context.Context, userID, noteID uuid.UUID) (*model.Note, error) {
	note, err := s.noteRepo.FindByID(ctx, userID, noteID)
	if err != nil {
		return nil, fmt.Errorf("find note: %w", err)
	}

	// Update access count
	_ = s.noteRepo.UpdateAccessCount(ctx, userID, noteID)

	// Log activity
	_ = s.activityRepo.Create(ctx, &model.Activity{
		UserID:  userID,
		NoteID:  &note.ID,
		Action:  model.ActionView,
	})

	return note, nil
}

// List lists notes for a user
func (s *NoteService) List(ctx context.Context, userID uuid.UUID, filter model.NoteFilter) ([]*model.Note, int64, error) {
	notes, total, err := s.noteRepo.List(ctx, userID, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("list notes: %w", err)
	}

	return notes, total, nil
}

// Search searches notes using full-text search
func (s *NoteService) Search(ctx context.Context, userID uuid.UUID, filter model.NoteFilter) ([]*model.Note, int64, error) {
	// Search uses the same List method with the Search filter
	notes, total, err := s.noteRepo.List(ctx, userID, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("search notes: %w", err)
	}

	return notes, total, nil
}

// Update updates a note
func (s *NoteService) Update(ctx context.Context, userID, noteID uuid.UUID, req *model.UpdateNoteRequest) (*model.Note, error) {
	// Validate request
	if err := util.ValidateStruct(req); err != nil {
		return nil, fmt.Errorf("%w: %s", model.ErrValidation, err)
	}

	// Get existing note
	note, err := s.noteRepo.FindByID(ctx, userID, noteID)
	if err != nil {
		return nil, fmt.Errorf("find note: %w", err)
	}

	// Update fields
	if req.Title != nil {
		note.Title = *req.Title
	}
	if req.Content != nil {
		note.Content = *req.Content
	}

	// Save changes
	if err := s.noteRepo.Update(ctx, note); err != nil {
		return nil, fmt.Errorf("update note: %w", err)
	}

	// Process links (delete old, create new)
	_ = s.linkRepo.DeleteByNote(ctx, userID, noteID)
	s.processLinks(ctx, userID, note)

	// Log activity
	_ = s.activityRepo.Create(ctx, &model.Activity{
		UserID:  userID,
		NoteID:  &note.ID,
		Action:  model.ActionUpdate,
	})

	return note, nil
}

// Delete soft deletes a note
func (s *NoteService) Delete(ctx context.Context, userID, noteID uuid.UUID) error {
	if err := s.noteRepo.Delete(ctx, userID, noteID); err != nil {
		return fmt.Errorf("delete note: %w", err)
	}

	// Delete associated links
	_ = s.linkRepo.DeleteByNote(ctx, userID, noteID)

	// Log activity
	_ = s.activityRepo.Create(ctx, &model.Activity{
		UserID:  userID,
		NoteID:  &noteID,
		Action:  model.ActionDelete,
	})

	return nil
}

// Restore restores a soft deleted note
func (s *NoteService) Restore(ctx context.Context, userID, noteID uuid.UUID) error {
	if err := s.noteRepo.Restore(ctx, userID, noteID); err != nil {
		return fmt.Errorf("restore note: %w", err)
	}

	return nil
}

// GetOrCreateDailyNote gets or creates a daily note for a given date
func (s *NoteService) GetOrCreateDailyNote(ctx context.Context, userID uuid.UUID, dateStr string) (*model.Note, bool, error) {
	// Try to find existing daily note for this date
	title := "Daily Note - " + dateStr
	note, err := s.noteRepo.FindByTitle(ctx, userID, title)

	// If found, return it
	if err == nil {
		return note, false, nil
	}

	// If not found (and not a different error), create it
	if err == repository.ErrNotFound {
		noteType := model.NoteTypeDaily
		req := &model.CreateNoteRequest{
			Title:    title,
			Content:  "",
			NoteType: noteType,
		}

		note, err = s.Create(ctx, userID, req)
		if err != nil {
			return nil, false, fmt.Errorf("create daily note: %w", err)
		}

		return note, true, nil
	}

	return nil, false, fmt.Errorf("find daily note: %w", err)
}

// GetOutgoingLinks gets all outgoing links from a note
func (s *NoteService) GetOutgoingLinks(ctx context.Context, userID, noteID uuid.UUID) ([]*model.Link, error) {
	// Verify note exists and belongs to user
	_, err := s.noteRepo.FindByID(ctx, userID, noteID)
	if err != nil {
		return nil, fmt.Errorf("note not found: %w", err)
	}

	// Get outgoing links
	links, err := s.linkRepo.GetBySource(ctx, userID, noteID)
	if err != nil {
		return nil, fmt.Errorf("get outgoing links: %w", err)
	}

	// Populate target note details
	for _, link := range links {
		targetNote, err := s.noteRepo.FindByID(ctx, userID, link.TargetNoteID)
		if err == nil {
			link.TargetNote = targetNote
		}
	}

	return links, nil
}

// GetBacklinks gets all incoming links to a note
func (s *NoteService) GetBacklinks(ctx context.Context, userID, noteID uuid.UUID) ([]*model.Link, error) {
	// Verify note exists and belongs to user
	_, err := s.noteRepo.FindByID(ctx, userID, noteID)
	if err != nil {
		return nil, fmt.Errorf("note not found: %w", err)
	}

	// Get backlinks
	links, err := s.linkRepo.GetByTarget(ctx, userID, noteID)
	if err != nil {
		return nil, fmt.Errorf("get backlinks: %w", err)
	}

	// Populate source note details
	for _, link := range links {
		sourceNote, err := s.noteRepo.FindByID(ctx, userID, link.SourceNoteID)
		if err == nil {
			link.SourceNote = sourceNote
		}
	}

	return links, nil
}

// GetLinkGraph gets the full knowledge graph for a user
func (s *NoteService) GetLinkGraph(ctx context.Context, userID uuid.UUID) (*model.GraphResponse, error) {
	// Get all notes for the user
	notes, _, err := s.noteRepo.List(ctx, userID, model.NoteFilter{
		Page:  1,
		Limit: 1000, // Get all notes
	})
	if err != nil {
		return nil, fmt.Errorf("get notes: %w", err)
	}

	// Build nodes map
	nodeMap := make(map[uuid.UUID]*model.GraphNode)
	nodes := make([]*model.GraphNode, 0, len(notes))
	for _, note := range notes {
		node := &model.GraphNode{
			ID:    note.ID,
			Title: note.Title,
			Type:  note.NoteType,
		}
		nodeMap[note.ID] = node
		nodes = append(nodes, node)
	}

	// Get all links for the user
	edges := make([]*model.GraphEdge, 0)
	for _, note := range notes {
		links, err := s.linkRepo.GetBySource(ctx, userID, note.ID)
		if err != nil {
			continue
		}

		for _, link := range links {
			// Only include edges where both nodes exist
			if _, sourceExists := nodeMap[link.SourceNoteID]; sourceExists {
				if _, targetExists := nodeMap[link.TargetNoteID]; targetExists {
					edges = append(edges, &model.GraphEdge{
						Source:    link.SourceNoteID,
						Target:    link.TargetNoteID,
						Context:   link.LinkContext,
						CreatedAt: link.CreatedAt,
					})
				}
			}
		}
	}

	return &model.GraphResponse{
		Nodes: nodes,
		Edges: edges,
	}, nil
}

// processLinks extracts wiki-style links and creates them in the database
func (s *NoteService) processLinks(ctx context.Context, userID uuid.UUID, note *model.Note) {
	links := s.linkParser.ExtractLinks(note.Content)
	for _, link := range links {
		// Try to find target note by title
		targetNote, err := s.noteRepo.FindByTitle(ctx, userID, link.Title)
		if err != nil {
			// Target note doesn't exist, skip
			continue
		}

		// Create link
		_ = s.linkRepo.Create(ctx, &model.Link{
			UserID:       userID,
			SourceNoteID: note.ID,
			TargetNoteID: targetNote.ID,
			LinkContext:  &link.Context,
		})
	}
}
