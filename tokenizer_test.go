package main

import (
	"testing"
	"unicode"
)

func TestTokenizer_trimLeft(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected string
	}{
		{
			name:     "Trim leading spaces",
			content:  "   hello",
			expected: "hello",
		},
		{
			name:     "No spaces to trim",
			content:  "hello",
			expected: "hello",
		},
		{
			name:     "All spaces",
			content:  "     ",
			expected: "",
		},
		{
			name:     "Empty string",
			content:  "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenizer := NewTokenizer([]rune(tt.content))
			tokenizer.trimLeftSpaces()
			result := string(tokenizer.data)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestTokenizer_removeFirst(t *testing.T) {
	tests := []struct {
		name      string
		content   string
		n         int
		expected  string
		remaining string
	}{
		{
			name:      "Remove first 3 characters",
			content:   "abcdef",
			n:         3,
			expected:  "abc",
			remaining: "def",
		},
		{
			name:      "Remove zero characters",
			content:   "abcdef",
			n:         0,
			expected:  "",
			remaining: "abcdef",
		},
		{
			name:      "Remove first character",
			content:   "hello",
			n:         1,
			expected:  "h",
			remaining: "ello",
		},
		{
			name:      "Remove more than available",
			content:   "hi",
			n:         5,
			expected:  "",
			remaining: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenizer := NewTokenizer([]rune(tt.content))
			result := string(tokenizer.removeFirst(tt.n))
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
			if string(tokenizer.data) != tt.remaining {
				t.Errorf("expected remaining content %q, got %q", tt.remaining, string(tokenizer.data))
			}
		})
	}
}

func TestTokenizer_removeIf(t *testing.T) {
	tests := []struct {
		name      string
		content   string
		predicate func(rune) bool
		expected  string
		remaining string
	}{
		{
			name:    "Remove while letters",
			content: "hello123",
			predicate: func(r rune) bool {
				return unicode.IsLetter(r)
			},
			expected:  "hello",
			remaining: "123",
		},
		{
			name:    "Remove while digits",
			content: "123abc",
			predicate: func(r rune) bool {
				return unicode.IsDigit(r)
			},
			expected:  "123",
			remaining: "abc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenizer := NewTokenizer([]rune(tt.content))
			result := string(tokenizer.removeIf(tt.predicate))
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
			if string(tokenizer.data) != tt.remaining {
				t.Errorf("expected remaining content %q, got %q", tt.remaining, string(tokenizer.data))
			}
		})
	}
}

func TestTokenizer_NextToken(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected string
		ok       bool
	}{
		{
			name:     "Next token - numeric",
			content:  "123abc",
			expected: "123",
			ok:       true,
		},
		{
			name:     "Next token - alphabetic",
			content:  "hello world",
			expected: "hello",
			ok:       true,
		},
		{
			name:     "Next token - alphabetic (with digits at the end)",
			content:  "hello123",
			expected: "hello123",
			ok:       true,
		},
		{
			name:     "Next token - single character",
			content:  "!",
			expected: "!",
			ok:       true,
		},
		{
			name:     "Next token - empty content",
			content:  "",
			expected: "",
			ok:       false,
		},
		{
			name:     "Next token - leading spaces",
			content:  "   hello",
			expected: "hello",
			ok:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenizer := NewTokenizer([]rune(tt.content))
			result, ok := tokenizer.NextToken()
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
			if ok != tt.ok {
				t.Errorf("expected ok = %v, got ok = %v", tt.ok, ok)
			}
		})
	}
}
