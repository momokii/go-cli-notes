package handler

import (
	"github.com/gofiber/fiber/v2"
)

// AuthHandler handles authentication HTTP requests
type AuthHandler struct {
	authService any // AuthService interface
}

// NoteHandler handles note HTTP requests
type NoteHandler struct {
	noteService any // NoteService interface
}

// Handlers holds all handlers
type Handlers struct {
	Auth     *AuthHandler
	Note     *NoteHandler
	Tag      *TagHandler
	Search   *SearchHandler
	Link     *LinkHandler
	Activity *ActivityHandler
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService any) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// NewNoteHandler creates a new note handler
func NewNoteHandler(noteService any) *NoteHandler {
	return &NoteHandler{
		noteService: noteService,
	}
}

// NewTagHandler creates a new tag handler
func NewTagHandler(tagService any) *TagHandler {
	return &TagHandler{
		tagService: tagService,
	}
}

// NewSearchHandler creates a new search handler
func NewSearchHandler(noteService any) *SearchHandler {
	return &SearchHandler{
		noteService: noteService,
	}
}

// NewLinkHandler creates a new link handler
func NewLinkHandler(noteService any) *LinkHandler {
	return &LinkHandler{
		noteService: noteService,
	}
}

// NewActivityHandler creates a new activity handler
func NewActivityHandler(activityRepo any, noteService any) *ActivityHandler {
	return &ActivityHandler{
		activityRepo: activityRepo,
		noteService:  noteService,
	}
}

// sendJSON sends a JSON response
func sendJSON(c *fiber.Ctx, status int, data any) error {
	return c.Status(status).JSON(data)
}

// sendError sends an error response
func sendError(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(fiber.Map{
		"error": message,
	})
}

// getUserID gets the user ID from the request context
func getUserID(c *fiber.Ctx) (string, bool) {
	userID := c.Locals("user_id")
	if userID == nil {
		return "", false
	}
	return userID.(string), true
}
