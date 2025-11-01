package atproto

import (
	"regexp"
	"unicode/utf8"

	"github.com/think-root/bluesky-connector/internal/models"
)

// DetectHashtags finds hashtags in text and returns facets for them
// Hashtags must:
// - Start with # followed by a non-digit, non-space character
// - Be max 66 characters (including #)
// - Not start with a digit
func DetectHashtags(text string) []models.RichTextFacet {
	var facets []models.RichTextFacet

	// Regex pattern for hashtags
	// (?:^|\s) - start of string or whitespace (non-capturing)
	// (#[^\d\s]\S*) - # followed by non-digit, non-space, then any non-space chars
	re := regexp.MustCompile(`(?:^|\s)(#[^\d\s]\S*)`)

	matches := re.FindAllStringSubmatchIndex(text, -1)

	for _, match := range matches {
		if len(match) < 4 {
			continue
		}

		// match[2] and match[3] are the start and end of the hashtag (group 1)
		tagStart := match[2]
		tagEnd := match[3]

		if tagStart == -1 || tagEnd == -1 {
			continue
		}

		tag := text[tagStart:tagEnd]

		// Strip trailing punctuation
		tag = stripTrailingPunctuation(tag)
		tagEnd = tagStart + len(tag)

		// Max length check (including #)
		if len(tag) > 66 {
			continue
		}

		// Remove the # from the tag value
		tagValue := tag[1:]
		if len(tagValue) == 0 {
			continue
		}

		// Convert string indices to UTF-8 byte offsets
		byteStart := utf8.RuneCountInString(text[:tagStart])
		byteEnd := utf8.RuneCountInString(text[:tagEnd])

		// Actually we need byte offsets, not rune counts
		byteStart = len(text[:tagStart])
		byteEnd = len(text[:tagEnd])

		facets = append(facets, models.RichTextFacet{
			Index: models.ByteSlice{
				ByteStart: byteStart,
				ByteEnd:   byteEnd,
			},
			Features: []models.Feature{
				{
					Type: "app.bsky.richtext.facet#tag",
					Tag:  tagValue,
				},
			},
		})
	}

	return facets
}

// stripTrailingPunctuation removes trailing punctuation from a string
func stripTrailingPunctuation(s string) string {
	// Remove common trailing punctuation
	punctuation := []rune{'.', ',', ';', '!', '?', ':', ')', ']', '}'}

	for len(s) > 0 {
		lastRune, size := utf8.DecodeLastRuneInString(s)
		isPunct := false
		for _, p := range punctuation {
			if lastRune == p {
				isPunct = true
				break
			}
		}
		if !isPunct {
			break
		}
		s = s[:len(s)-size]
	}

	return s
}
