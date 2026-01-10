package model

// SearchResult represents a single search result
type SearchResult struct {
	Note    *Note  `json:"note"`
	Rank    float64 `json:"rank"`    // Relevance score (0-1)
	Snippet string  `json:"snippet"` // Highlighted text excerpt
}

// SearchRequest represents a search request
type SearchRequest struct {
	Query    string `query:"q" validate:"required,min=1,max=500"`
	Page     int    `query:"page" validate:"min=1"`
	Limit    int    `query:"limit" validate:"min=1,max=100"`
	NoteType *NoteType `query:"type"`
	TagID    *string `query:"tag_id"`
}

// SearchResponse represents a search response
type SearchResponse struct {
	Query      string          `json:"query"`
	Results    []*SearchResult `json:"results"`
	Pagination *Pagination     `json:"pagination"`
}

// SuggestionRequest represents an autocomplete suggestion request
type SuggestionRequest struct {
	Query string `query:"q" validate:"required,min=1,max=100"`
	Limit int    `query:"limit" validate:"min=1,max=20"`
}

// Suggestion represents an autocomplete suggestion
type Suggestion struct {
	Type  string `json:"type"` // "note" or "tag"
	ID    string `json:"id"`
	Title string `json:"title"`
}

// DailyNoteRequest represents a daily note request
type DailyNoteRequest struct {
	Date string `param:"date" validate:"required,datetime=2006-01-02"` // YYYY-MM-DD format
}

// DailyNoteResponse represents a daily note response
type DailyNoteResponse struct {
	Note       *Note       `json:"note"`
	PrevDate   *string     `json:"prev_date,omitempty"` // YYYY-MM-DD
	NextDate   *string     `json:"next_date,omitempty"` // YYYY-MM-DD
	RelatedNotes []*Note   `json:"related_notes,omitempty"` // Notes from same week
}
