-- 000002_add_fts.down.sql
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
