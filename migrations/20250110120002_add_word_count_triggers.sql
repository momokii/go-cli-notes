-- +goose Up
-- Add word count and reading time calculation triggers

-- Function to calculate word count and reading time
CREATE OR REPLACE FUNCTION calculate_note_metrics()
RETURNS TRIGGER AS $$
DECLARE
    word_count INT;
    reading_time INT;
    cleaned_content TEXT;
BEGIN
    -- Use content if exists, otherwise empty string
    cleaned_content := COALESCE(NEW.content, '');

    -- Calculate word count using regexp_split_to_array
    -- This handles NULL content and counts words properly
    word_count := array_length(regexp_split_to_array(cleaned_content, '\s+'), 1);

    -- Handle NULL/empty cases
    IF word_count IS NULL THEN
        word_count := 0;
    END IF;

    -- Filter out empty strings from the split result
    IF cleaned_content = '' OR cleaned_content IS NULL THEN
        word_count := 0;
    END IF;

    -- Calculate reading time: assume 200 words per minute (standard)
    -- Minimum 1 minute if there's any content
    IF word_count = 0 THEN
        reading_time := 0;
    ELSE
        reading_time := (word_count + 199) / 200; -- Ceiling division
    END IF;

    -- Set the calculated values
    NEW.word_count := word_count;
    NEW.reading_time_minutes := reading_time;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger for INSERT
DROP TRIGGER IF EXISTS calculate_note_metrics_on_insert ON notes;
CREATE TRIGGER calculate_note_metrics_on_insert
    BEFORE INSERT ON notes
    FOR EACH ROW
    EXECUTE FUNCTION calculate_note_metrics();

-- Trigger for UPDATE (only when content changes)
DROP TRIGGER IF EXISTS calculate_note_metrics_on_update ON notes;
CREATE TRIGGER calculate_note_metrics_on_update
    BEFORE UPDATE OF content ON notes
    FOR EACH ROW
    WHEN (OLD.content IS DISTINCT FROM NEW.content)
    EXECUTE FUNCTION calculate_note_metrics();

-- +goose Down
-- Remove word count and reading time calculation triggers

DROP TRIGGER IF EXISTS calculate_note_metrics_on_insert ON notes;
DROP TRIGGER IF EXISTS calculate_note_metrics_on_update ON notes;
DROP FUNCTION IF EXISTS calculate_note_metrics();
