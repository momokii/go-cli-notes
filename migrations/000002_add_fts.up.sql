-- 000002_add_fts.up.sql
-- Add full-text search support to notes

-- Add tsvector column for full-text search
ALTER TABLE notes ADD COLUMN content_tsv tsvector;

-- Create function to generate tsvector from title and content
CREATE OR REPLACE FUNCTION notes_tsv_trigger() RETURNS trigger AS $$
BEGIN
    NEW.content_tsv :=
        setweight(to_tsvector('english', coalesce(NEW.title, '')), 'A') ||
        setweight(to_tsvector('english', coalesce(NEW.content, '')), 'B');
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger to auto-update tsvector on insert/update
CREATE TRIGGER notes_tsv_update BEFORE INSERT OR UPDATE ON notes
    FOR EACH ROW EXECUTE FUNCTION notes_tsv_trigger();

-- Create GIN index for fast full-text search
CREATE INDEX idx_notes_content_tsv ON notes USING gin(content_tsv);

-- Create index for combined text search (without tsvector, for alternative search)
CREATE INDEX idx_notes_title_trgm ON notes USING gin(title gin_trgm_ops);
CREATE INDEX idx_notes_content_trgm ON notes USING gin(content gin_trgm_ops);

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

-- Create trigger for word count and reading time
CREATE TRIGGER notes_metadata_update BEFORE INSERT OR UPDATE ON notes
    FOR EACH ROW EXECUTE FUNCTION notes_metadata_trigger();

-- Enable pg_trgm extension for trigram-based search (useful for fuzzy matching)
CREATE EXTENSION IF NOT EXISTS pg_trgm;
