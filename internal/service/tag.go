package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/momokii/go-cli-notes/internal/model"
	"github.com/momokii/go-cli-notes/internal/repository"
	"github.com/momokii/go-cli-notes/internal/util"
)

// TagService handles tag business logic
type TagService struct {
	tagRepo     repository.TagRepository
	noteRepo    repository.NoteRepository
	activityRepo repository.ActivityRepository
}

// NewTagService creates a new tag service
func NewTagService(
	tagRepo repository.TagRepository,
	noteRepo repository.NoteRepository,
	activityRepo repository.ActivityRepository,
) *TagService {
	return &TagService{
		tagRepo:     tagRepo,
		noteRepo:    noteRepo,
		activityRepo: activityRepo,
	}
}

// Create creates a new tag
func (s *TagService) Create(ctx context.Context, userID uuid.UUID, req *model.CreateTagRequest) (*model.Tag, error) {
	// Validate request
	if err := util.ValidateStruct(req); err != nil {
		return nil, fmt.Errorf("%w: %s", model.ErrValidation, err)
	}

	// Check if tag already exists
	_, err := s.tagRepo.FindByName(ctx, userID, req.Name)
	if err == nil {
		return nil, fmt.Errorf("tag with name '%s' already exists", req.Name)
	}

	tag := &model.Tag{
		UserID: userID,
		Name:   req.Name,
		Color:  req.Color,
	}

	if err := s.tagRepo.Create(ctx, tag); err != nil {
		return nil, fmt.Errorf("create tag: %w", err)
	}

	return tag, nil
}

// GetByID gets a tag by ID
func (s *TagService) GetByID(ctx context.Context, userID, tagID uuid.UUID) (*model.Tag, error) {
	tag, err := s.tagRepo.FindByID(ctx, userID, tagID)
	if err != nil {
		return nil, fmt.Errorf("find tag: %w", err)
	}
	return tag, nil
}

// List lists tags for a user with pagination
func (s *TagService) List(ctx context.Context, userID uuid.UUID, page, limit int) ([]*model.TagWithCount, int64, error) {
	tags, total, err := s.tagRepo.ListWithNoteCount(ctx, userID, page, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("list tags: %w", err)
	}
	return tags, total, nil
}

// Update updates a tag
func (s *TagService) Update(ctx context.Context, userID, tagID uuid.UUID, req *model.UpdateTagRequest) (*model.Tag, error) {
	// Validate request
	if err := util.ValidateStruct(req); err != nil {
		return nil, fmt.Errorf("%w: %s", model.ErrValidation, err)
	}

	// Get existing tag
	tag, err := s.tagRepo.FindByID(ctx, userID, tagID)
	if err != nil {
		return nil, fmt.Errorf("find tag: %w", err)
	}

	// Check name uniqueness if updating name
	if req.Name != nil && *req.Name != tag.Name {
		existing, _ := s.tagRepo.FindByName(ctx, userID, *req.Name)
		if existing != nil {
			return nil, fmt.Errorf("tag with name '%s' already exists", *req.Name)
		}
		tag.Name = *req.Name
	}

	if req.Color != nil {
		tag.Color = req.Color
	}

	if err := s.tagRepo.Update(ctx, tag); err != nil {
		return nil, fmt.Errorf("update tag: %w", err)
	}

	return tag, nil
}

// Delete deletes a tag
func (s *TagService) Delete(ctx context.Context, userID, tagID uuid.UUID) error {
	if err := s.tagRepo.Delete(ctx, userID, tagID); err != nil {
		return fmt.Errorf("delete tag: %w", err)
	}
	return nil
}

// AddToNote adds a tag to a note
func (s *TagService) AddToNote(ctx context.Context, userID, noteID, tagID uuid.UUID) error {
	// Verify note ownership
	_, err := s.noteRepo.FindByID(ctx, userID, noteID)
	if err != nil {
		return fmt.Errorf("note not found: %w", err)
	}

	// Verify tag ownership
	_, err = s.tagRepo.FindByID(ctx, userID, tagID)
	if err != nil {
		return fmt.Errorf("tag not found: %w", err)
	}

	if err := s.tagRepo.AddToNote(ctx, noteID, tagID); err != nil {
		return fmt.Errorf("add tag to note: %w", err)
	}

	return nil
}

// RemoveFromNote removes a tag from a note
func (s *TagService) RemoveFromNote(ctx context.Context, userID, noteID, tagID uuid.UUID) error {
	// Verify note ownership
	_, err := s.noteRepo.FindByID(ctx, userID, noteID)
	if err != nil {
		return fmt.Errorf("note not found: %w", err)
	}

	if err := s.tagRepo.RemoveFromNote(ctx, noteID, tagID); err != nil {
		return fmt.Errorf("remove tag from note: %w", err)
	}

	return nil
}

// GetByNote gets all tags for a note
func (s *TagService) GetByNote(ctx context.Context, userID, noteID uuid.UUID) ([]*model.Tag, error) {
	// Verify note ownership
	_, err := s.noteRepo.FindByID(ctx, userID, noteID)
	if err != nil {
		return nil, fmt.Errorf("note not found: %w", err)
	}

	tags, err := s.tagRepo.GetByNote(ctx, noteID)
	if err != nil {
		return nil, fmt.Errorf("get tags by note: %w", err)
	}

	return tags, nil
}

// GetNotesByTag gets all notes for a tag
func (s *TagService) GetNotesByTag(ctx context.Context, userID, tagID uuid.UUID) ([]*model.Note, error) {
	// Verify tag ownership
	_, err := s.tagRepo.FindByID(ctx, userID, tagID)
	if err != nil {
		return nil, fmt.Errorf("tag not found: %w", err)
	}

	notes, err := s.tagRepo.GetNotesByTag(ctx, userID, tagID)
	if err != nil {
		return nil, fmt.Errorf("get notes by tag: %w", err)
	}

	return notes, nil
}
