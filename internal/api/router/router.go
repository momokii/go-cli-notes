package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/momokii/go-cli-notes/internal/api/handler"
	"github.com/momokii/go-cli-notes/internal/api/middleware"
)

// Setup configures all routes for the API
func Setup(app *fiber.App, h *handler.Handlers, jwtManager any) {
	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// API v1 routes
	v1 := app.Group("/api/v1")

	// Auth routes (public)
	auth := v1.Group("/auth")
	auth.Post("/register", h.Auth.Register)
	auth.Post("/login", h.Auth.Login)
	auth.Post("/refresh", h.Auth.RefreshToken)
	auth.Post("/logout", middleware.Auth(jwtManager), h.Auth.Logout)

	// Tag routes (authenticated)
	tags := v1.Group("/tags")
	tags.Use(middleware.Auth(jwtManager))
	tags.Get("/", h.Tag.ListTags)
	tags.Post("/", h.Tag.CreateTag)
	tags.Get("/:id/notes", h.Tag.GetTagNotes)
	tags.Get("/:id", h.Tag.GetTag)
	tags.Put("/:id", h.Tag.UpdateTag)
	tags.Delete("/:id", h.Tag.DeleteTag)

	// Note routes (authenticated)
	notes := v1.Group("/notes")
	notes.Use(middleware.Auth(jwtManager))

	// Define specific routes BEFORE parameterized routes
	notes.Get("/graph", h.Link.GetLinkGraph)
	notes.Get("/daily/:date", h.Note.GetOrCreateDailyNote)
	notes.Get("/trending", h.Activity.GetTrendingNotes)
	notes.Get("/forgotten", h.Activity.GetForgottenNotes)

	// General note routes
	notes.Post("/", h.Note.Create)
	notes.Get("/", h.Note.List)
	notes.Get("/:id", h.Note.GetByID)
	notes.Put("/:id", h.Note.Update)
	notes.Delete("/:id", h.Note.Delete)

	// Note-Tag association routes
	notes.Get("/:id/tags", h.Tag.GetNoteTags)
	notes.Post("/:id/tags/:tag_id", h.Tag.AddTagToNote)
	notes.Delete("/:id/tags/:tag_id", h.Tag.RemoveTagFromNote)

	// Note-Link association routes
	notes.Get("/:id/links", h.Link.GetOutgoingLinks)
	notes.Get("/:id/backlinks", h.Link.GetBacklinks)

	// Search routes (authenticated)
	search := v1.Group("/search")
	search.Use(middleware.Auth(jwtManager))
	search.Get("/", h.Search.Search)

	// Activity routes (authenticated)
	activity := v1.Group("/activity")
	activity.Use(middleware.Auth(jwtManager))
	activity.Get("/recent", h.Activity.GetRecentActivity)

	// Stats routes (authenticated)
	stats := v1.Group("/stats")
	stats.Use(middleware.Auth(jwtManager))
	stats.Get("/", h.Activity.GetUserStats)
}
