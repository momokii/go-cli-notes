package util

import (
	"regexp"
	"strings"
)

// LinkParser parses wiki-style links from note content
type LinkParser struct {
	// Regex pattern matches [[Note Title]] or [[Note Title|Display Text]]
	linkPattern *regexp.Regexp
}

// NewLinkParser creates a new link parser
func NewLinkParser() *LinkParser {
	// Pattern matches [[...]] with optional display text
	// [[Note Title]] or [[Note Title|Display Text]]
	pattern := regexp.MustCompile(`\[\[([^\]|]+)(?:\|([^\]]+))?\]\]`)
	return &LinkParser{
		linkPattern: pattern,
	}
}

// Link represents a parsed wiki-style link
type Link struct {
	Title      string   // The note title to link to
	Display    string   // The display text (defaults to title)
	Context    string   // Surrounding text context
	StartPos   int      // Position in content
	EndPos     int      // End position in content
}

// ExtractLinks extracts all wiki-style links from note content
func (p *LinkParser) ExtractLinks(content string) []*Link {
	if content == "" {
		return nil
	}

	matches := p.linkPattern.FindAllStringSubmatchIndex(content, -1)
	if len(matches) == 0 {
		return nil
	}

	links := make([]*Link, 0, len(matches))

	for _, match := range matches {
		fullMatch := content[match[0]:match[1]]

		// Extract title and display text using the pattern
		submatches := p.linkPattern.FindStringSubmatch(fullMatch)
		if len(submatches) < 2 {
			continue
		}

		title := strings.TrimSpace(submatches[1])
		display := title
		if len(submatches) > 2 && submatches[2] != "" {
			display = strings.TrimSpace(submatches[2])
		}

		// Get context (50 chars before and after)
		contextStart := match[0] - 50
		if contextStart < 0 {
			contextStart = 0
		}
		contextEnd := match[1] + 50
		if contextEnd > len(content) {
			contextEnd = len(content)
		}
		context := strings.TrimSpace(content[contextStart:contextEnd])

		link := &Link{
			Title:    title,
			Display:  display,
			Context:  context,
			StartPos: match[0],
			EndPos:   match[1],
		}

		links = append(links, link)
	}

	return links
}

// ReplaceLinks replaces wiki-style links with markdown links
func (p *LinkParser) ReplaceLinks(content string, replacer func(title, display string) string) string {
	if content == "" {
		return content
	}

	return p.linkPattern.ReplaceAllStringFunc(content, func(match string) string {
		submatches := p.linkPattern.FindStringSubmatch(match)
		if len(submatches) < 2 {
			return match
		}

		title := strings.TrimSpace(submatches[1])
		display := title
		if len(submatches) > 2 && submatches[2] != "" {
			display = strings.TrimSpace(submatches[2])
		}

		return replacer(title, display)
	})
}

// NormalizeTitle normalizes a note title for matching
func (p *LinkParser) NormalizeTitle(title string) string {
	// Convert to lowercase and trim spaces
	title = strings.ToLower(strings.TrimSpace(title))
	// Remove special characters (keep alphanumeric, spaces, hyphens, underscores)
	title = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == ' ' || r == '-' || r == '_' {
			return r
		}
		return -1
	}, title)
	// Replace multiple spaces with single space
	title = strings.Join(strings.Fields(title), " ")
	return title
}

// StripLinks removes all wiki-style links from content
func (p *LinkParser) StripLinks(content string) string {
	if content == "" {
		return content
	}
	return p.linkPattern.ReplaceAllString(content, "$2")
}
