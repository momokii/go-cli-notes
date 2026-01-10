package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/momokii/go-cli-notes/internal/model"
	"github.com/momokii/go-cli-notes/internal/service"
)

type TagHandler struct {
	tagService any // TagService interface
}

// ListTags handles GET /api/v1/tags/
func (h *TagHandler) ListTags(c *fiber.Ctx) error {
	userIDStr, ok := getUserID(c)
	if !ok {
		return sendError(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return sendError(c, fiber.StatusBadRequest, "Invalid user ID")
	}

	// Parse pagination params
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 20)

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	svc, ok := h.tagService.(*service.TagService)
	if !ok {
		return sendError(c, fiber.StatusInternalServerError, "Service error")
	}

	tags, total, err := svc.List(c.Context(), userID, page, limit)
	if err != nil {
		return sendError(c, fiber.StatusInternalServerError, "Failed to list tags")
	}

	totalPages := int(total) / limit
	if int(total)%limit != 0 {
		totalPages++
	}

	response := struct {
		Tags       []*model.TagWithCount `json:"tags"`
		Pagination *model.Pagination     `json:"pagination"`
	}{
		Tags: tags,
		Pagination: &model.Pagination{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}

	return sendJSON(c, fiber.StatusOK, response)
}

// GetTag handles GET /api/v1/tags/:id
func (h *TagHandler) GetTag(c *fiber.Ctx) error {
	userIDStr, ok := getUserID(c)
	if !ok {
		return sendError(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return sendError(c, fiber.StatusBadRequest, "Invalid user ID")
	}

	tagID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return sendError(c, fiber.StatusBadRequest, "Invalid tag ID")
	}

	svc, ok := h.tagService.(*service.TagService)
	if !ok {
		return sendError(c, fiber.StatusInternalServerError, "Service error")
	}

	tag, err := svc.GetByID(c.Context(), userID, tagID)
	if err != nil {
		return sendError(c, fiber.StatusNotFound, "Tag not found")
	}

	return sendJSON(c, fiber.StatusOK, tag)
}

// CreateTag handles POST /api/v1/tags/
func (h *TagHandler) CreateTag(c *fiber.Ctx) error {
	userIDStr, ok := getUserID(c)
	if !ok {
		return sendError(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return sendError(c, fiber.StatusBadRequest, "Invalid user ID")
	}

	var req model.CreateTagRequest
	if err := c.BodyParser(&req); err != nil {
		return sendError(c, fiber.StatusBadRequest, "Invalid request body")
	}

	svc, ok := h.tagService.(*service.TagService)
	if !ok {
		return sendError(c, fiber.StatusInternalServerError, "Service error")
	}

	tag, err := svc.Create(c.Context(), userID, &req)
	if err != nil {
		if err.Error() == "tag with name '"+req.Name+"' already exists" {
			return sendError(c, fiber.StatusConflict, "Tag with this name already exists")
		}
		return sendError(c, fiber.StatusInternalServerError, "Failed to create tag")
	}

	return sendJSON(c, fiber.StatusCreated, tag)
}

// UpdateTag handles PUT /api/v1/tags/:id
func (h *TagHandler) UpdateTag(c *fiber.Ctx) error {
	userIDStr, ok := getUserID(c)
	if !ok {
		return sendError(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return sendError(c, fiber.StatusBadRequest, "Invalid user ID")
	}

	tagID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return sendError(c, fiber.StatusBadRequest, "Invalid tag ID")
	}

	var req model.UpdateTagRequest
	if err := c.BodyParser(&req); err != nil {
		return sendError(c, fiber.StatusBadRequest, "Invalid request body")
	}

	svc, ok := h.tagService.(*service.TagService)
	if !ok {
		return sendError(c, fiber.StatusInternalServerError, "Service error")
	}

	tag, err := svc.Update(c.Context(), userID, tagID, &req)
	if err != nil {
		if err.Error() == "tag with name '"+*req.Name+"' already exists" {
			return sendError(c, fiber.StatusConflict, "Tag with this name already exists")
		}
		return sendError(c, fiber.StatusInternalServerError, "Failed to update tag")
	}

	return sendJSON(c, fiber.StatusOK, tag)
}

// DeleteTag handles DELETE /api/v1/tags/:id
func (h *TagHandler) DeleteTag(c *fiber.Ctx) error {
	userIDStr, ok := getUserID(c)
	if !ok {
		return sendError(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return sendError(c, fiber.StatusBadRequest, "Invalid user ID")
	}

	tagID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return sendError(c, fiber.StatusBadRequest, "Invalid tag ID")
	}

	svc, ok := h.tagService.(*service.TagService)
	if !ok {
		return sendError(c, fiber.StatusInternalServerError, "Service error")
	}

	if err := svc.Delete(c.Context(), userID, tagID); err != nil {
		return sendError(c, fiber.StatusNotFound, "Tag not found")
	}

	return sendJSON(c, fiber.StatusNoContent, nil)
}

// AddTagToNote handles POST /api/v1/notes/:id/tags/:tag_id
func (h *TagHandler) AddTagToNote(c *fiber.Ctx) error {
	userIDStr, ok := getUserID(c)
	if !ok {
		return sendError(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return sendError(c, fiber.StatusBadRequest, "Invalid user ID")
	}

	noteID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return sendError(c, fiber.StatusBadRequest, "Invalid note ID")
	}

	tagID, err := uuid.Parse(c.Params("tag_id"))
	if err != nil {
		return sendError(c, fiber.StatusBadRequest, "Invalid tag ID")
	}

	svc, ok := h.tagService.(*service.TagService)
	if !ok {
		return sendError(c, fiber.StatusInternalServerError, "Service error")
	}

	if err := svc.AddToNote(c.Context(), userID, noteID, tagID); err != nil {
		return sendError(c, fiber.StatusInternalServerError, "Failed to add tag to note")
	}

	return sendJSON(c, fiber.StatusOK, fiber.Map{"message": "Tag added to note"})
}

// RemoveTagFromNote handles DELETE /api/v1/notes/:id/tags/:tag_id
func (h *TagHandler) RemoveTagFromNote(c *fiber.Ctx) error {
	userIDStr, ok := getUserID(c)
	if !ok {
		return sendError(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return sendError(c, fiber.StatusBadRequest, "Invalid user ID")
	}

	noteID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return sendError(c, fiber.StatusBadRequest, "Invalid note ID")
	}

	tagID, err := uuid.Parse(c.Params("tag_id"))
	if err != nil {
		return sendError(c, fiber.StatusBadRequest, "Invalid tag ID")
	}

	svc, ok := h.tagService.(*service.TagService)
	if !ok {
		return sendError(c, fiber.StatusInternalServerError, "Service error")
	}

	if err := svc.RemoveFromNote(c.Context(), userID, noteID, tagID); err != nil {
		return sendError(c, fiber.StatusNotFound, "Tag not associated with note")
	}

	return sendJSON(c, fiber.StatusOK, fiber.Map{"message": "Tag removed from note"})
}

// GetNoteTags handles GET /api/v1/notes/:id/tags
func (h *TagHandler) GetNoteTags(c *fiber.Ctx) error {
	userIDStr, ok := getUserID(c)
	if !ok {
		return sendError(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return sendError(c, fiber.StatusBadRequest, "Invalid user ID")
	}

	noteID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return sendError(c, fiber.StatusBadRequest, "Invalid note ID")
	}

	svc, ok := h.tagService.(*service.TagService)
	if !ok {
		return sendError(c, fiber.StatusInternalServerError, "Service error")
	}

	tags, err := svc.GetByNote(c.Context(), userID, noteID)
	if err != nil {
		return sendError(c, fiber.StatusNotFound, "Note not found")
	}

	return sendJSON(c, fiber.StatusOK, fiber.Map{"tags": tags})
}

// Helper for string to int conversion
func strconvParseInt(s string) (int, error) {
	return strconv.Atoi(s)
}
