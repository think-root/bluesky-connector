package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitTextIntoParts(t *testing.T) {
	client := &BlueSkyClient{}

	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{
			name:     "Short text",
			input:    "Hello, world!",
			expected: 1,
		},
		{
			name:     "Exactly max length",
			input:    string(make([]byte, MaxPostLength)),
			expected: 1,
		},
		{
			name:     "Slightly over max length",
			input:    string(make([]byte, MaxPostLength+10)),
			expected: 2,
		},
		{
			name:     "Long text with spaces",
			input:    "This is a very long text that should be split into multiple parts because it exceeds the maximum post length limit. " + string(make([]byte, MaxPostLength)),
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parts := client.splitTextIntoParts(tt.input)
			assert.Equal(t, tt.expected, len(parts))
			
			// Verify all parts are within limits
			for i, part := range parts {
				assert.LessOrEqual(t, len(part), MaxPostLength, "Part %d exceeds max length", i)
			}
			
			// Verify all parts combined equal original (minus whitespace)
			combined := ""
			for _, part := range parts {
				combined += part + " "
			}
			// This is a basic check - in real implementation we'd need more sophisticated comparison
			assert.Contains(t, combined, tt.input[:min(50, len(tt.input))])
		})
	}
}

func TestMinMax(t *testing.T) {
	assert.Equal(t, 5, min(5, 10))
	assert.Equal(t, 5, min(10, 5))
	assert.Equal(t, 10, max(5, 10))
	assert.Equal(t, 10, max(10, 5))
}