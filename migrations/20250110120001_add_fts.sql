-- +goose Up
-- Add full-text search support to notes
-- NOTE: This migration is idempotent and can be safely re-run

-- Add tsvector column for full-text search (idempotent using DO block)
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'notes' AND column_name = 'content_tsv'
    ) THEN
        ALTER TABLE notes ADD COLUMN content_tsv tsvector;
    END IF;
END $$;

-- Create function to generate tsvector from title and content
CREATE OR REPLACE FUNCTION notes_tsv_trigger() RETURNS trigger AS $$
BEGIN
    NEW.content_tsv :=
        setweight(to_tsvector('english', coalesce(NEW.title, '')), 'A') ||
        setweight(to_tsvector('english', coalesce(NEW.content, '')), 'B');
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger to auto-update tsvector on insert/update (idempotent)
DROP TRIGGER IF EXISTS notes_tsv_update ON notes;
CREATE TRIGGER notes_tsv_update BEFORE INSERT OR UPDATE ON notes
    FOR EACH ROW EXECUTE FUNCTION notes_tsv_trigger();

-- Create GIN index for fast full-text search (idempotent)
CREATE INDEX IF NOT EXISTS idx_notes_content_tsv ON notes USING gin(content_tsv);

-- Create index for combined text search (idempotent)
CREATE INDEX IF NOT EXISTS idx_notes_title_trgm ON notes USING gin(title gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_notes_content_trgm ON notes USING gin(content gin_trgm_ops);

-- Create function to calculate word count
CREATE OR REPLACE FUNCTION calculate_word_count(text_content TEXT) RETURNS INT AS $$
BEGIN
    IF text_content IS NULL OR text_content = '' THEN
        RETURN 0;
    END IF;
    RETURN array_length(regexp_split_to_array(text_content, '\s+'), 1);
END;
$$ LANGUAGE plpgsql;

-- Create function to calculate reading time (assuming 200 words per minute)
CREATE OR REPLACE FUNCTION calculate_reading_time(word_count INT) RETURNS INT AS $$
BEGIN
    IF word_count IS NULL OR word_count = 0 THEN
        RETURN 0;
    END IF;
    RETURN CEIL(word_count::float / 200.0);
END;
$$ LANGUAGE plpgsql;

-- Create trigger to auto-calculate word count and reading time
CREATE OR REPLACE FUNCTION notes_metadata_trigger() RETURNS trigger AS $$
BEGIN
    -- Calculate word count from content
    NEW.word_count := calculate_word_count(NEW.content);

    -- Calculate reading time
    NEW.reading_time_minutes := calculate_reading_time(NEW.word_count);

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger for word count and reading time (idempotent)
DROP TRIGGER IF EXISTS notes_metadata_update ON notes;
CREATE TRIGGER notes_metadata_update BEFORE INSERT OR UPDATE ON notes
    FOR EACH ROW EXECUTE FUNCTION notes_metadata_trigger();

-- Enable pg_trgm extension for trigram-based search (useful for fuzzy matching)
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- +goose Down
-- Rollback full-text search support

-- Drop triggers
DROP TRIGGER IF EXISTS notes_metadata_update ON notes;
DROP TRIGGER IF EXISTS notes_tsv_update ON notes;

-- Drop functions
DROP FUNCTION IF EXISTS notes_metadata_trigger();
DROP FUNCTION IF EXISTS notes_tsv_trigger();
DROP FUNCTION IF EXISTS calculate_word_count(TEXT);
DROP FUNCTION IF EXISTS calculate_reading_time(INT);

-- Drop indexes
DROP INDEX IF EXISTS idx_notes_content_trgm;
DROP INDEX IF EXISTS idx_notes_title_trgm;
DROP INDEX IF EXISTS idx_notes_content_tsv;

-- Drop column
ALTER TABLE notes DROP COLUMN IF EXISTS content_tsv;

-- Drop extension
DROP EXTENSION IF EXISTS pg_trgm;
