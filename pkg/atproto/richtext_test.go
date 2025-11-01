package atproto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDetectHashtags(t *testing.T) {
	tests := []struct {
		name          string
		text          string
		expectedCount int
		expectedTags  []string
		expectedStart []int
		expectedEnd   []int
	}{
		{
			name:          "Single hashtag",
			text:          "Hello #GitHub world",
			expectedCount: 1,
			expectedTags:  []string{"GitHub"},
			expectedStart: []int{6},
			expectedEnd:   []int{13},
		},
		{
			name:          "Multiple hashtags",
			text:          "Check out #GitHub and #OpenSource projects",
			expectedCount: 2,
			expectedTags:  []string{"GitHub", "OpenSource"},
			expectedStart: []int{10, 22},
			expectedEnd:   []int{17, 33},
		},
		{
			name:          "Hashtag at start",
			text:          "#GitHub is awesome",
			expectedCount: 1,
			expectedTags:  []string{"GitHub"},
			expectedStart: []int{0},
			expectedEnd:   []int{7},
		},
		{
			name:          "Hashtag at end",
			text:          "Posted on #GitHub",
			expectedCount: 1,
			expectedTags:  []string{"GitHub"},
			expectedStart: []int{10},
			expectedEnd:   []int{17},
		},
		{
			name:          "Hashtag with trailing punctuation",
			text:          "Check #GitHub! It's great.",
			expectedCount: 1,
			expectedTags:  []string{"GitHub"},
			expectedStart: []int{6},
			expectedEnd:   []int{13},
		},
		{
			name:          "No hashtags",
			text:          "Just plain text without tags",
			expectedCount: 0,
			expectedTags:  []string{},
		},
		{
			name:          "Hashtag starting with digit (invalid)",
			text:          "This #123test should not match",
			expectedCount: 0,
			expectedTags:  []string{},
		},
		{
			name:          "Multiple hashtags with newlines",
			text:          "First line #GitHub\n\nSecond line #OpenSource",
			expectedCount: 2,
			expectedTags:  []string{"GitHub", "OpenSource"},
			expectedStart: []int{11, 32},
			expectedEnd:   []int{18, 43},
		},
		{
			name:          "Hashtags in thread format",
			text:          "🧵 1/2 Some content\n\n#GitHub #OpenSource",
			expectedCount: 2,
			expectedTags:  []string{"GitHub", "OpenSource"},
			// Note: 🧵 emoji is 4 bytes in UTF-8
			expectedStart: []int{23, 31},
			expectedEnd:   []int{30, 42},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			facets := DetectHashtags(tt.text)

			assert.Equal(t, tt.expectedCount, len(facets), "Number of facets should match")

			for i, facet := range facets {
				if i < len(tt.expectedTags) {
					assert.Equal(t, "app.bsky.richtext.facet#tag", facet.Features[0].Type)
					assert.Equal(t, tt.expectedTags[i], facet.Features[0].Tag)
				}

				if i < len(tt.expectedStart) {
					assert.Equal(t, tt.expectedStart[i], facet.Index.ByteStart,
						"ByteStart for tag #%d should match", i)
				}

				if i < len(tt.expectedEnd) {
					assert.Equal(t, tt.expectedEnd[i], facet.Index.ByteEnd,
						"ByteEnd for tag #%d should match", i)
				}
			}
		})
	}
}

func TestStripTrailingPunctuation(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "No punctuation",
			input:    "#GitHub",
			expected: "#GitHub",
		},
		{
			name:     "Exclamation mark",
			input:    "#GitHub!",
			expected: "#GitHub",
		},
		{
			name:     "Multiple punctuation",
			input:    "#GitHub!?",
			expected: "#GitHub",
		},
		{
			name:     "Period",
			input:    "#GitHub.",
			expected: "#GitHub",
		},
		{
			name:     "Comma",
			input:    "#GitHub,",
			expected: "#GitHub",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := stripTrailingPunctuation(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
