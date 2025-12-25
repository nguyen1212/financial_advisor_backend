package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ToNewsCategory(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		input    string
		expected NewsCategory
	}{
		{
			name:     "Finance category",
			input:    "finance",
			expected: NewsCategoryFinance,
		},
		{
			name:     "Military category",
			input:    "military",
			expected: NewsCategoryMilitary,
		},
		{
			name:     "Unknown category",
			input:    "sports",
			expected: NewsCategoryUnknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := ToNewsCategory(tt.input)

			assert.Equal(t, result, tt.expected)
		})
	}
}

func Test_NewsStatus_String(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		status   NewsStatus
		expected string
	}{
		{
			name:     "Added status",
			status:   NewsStatusAdded,
			expected: "added",
		},
		{
			name:     "Synced status",
			status:   NewsStatusSynced,
			expected: "synced",
		},
		{
			name:     "Failed status",
			status:   NewsStatusFailed,
			expected: "failed",
		},
		{
			name:     "Unknown status",
			status:   NewsStatusUnknown,
			expected: "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := tt.status.String()

			assert.Equal(t, tt.expected, result)
		})
	}
}

func Test_NewsCategory_String(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		category NewsCategory
		expected string
	}{
		{
			name:     "Finance category",
			category: NewsCategoryFinance,
			expected: "finance",
		},
		{
			name:     "Military category",
			category: NewsCategoryMilitary,
			expected: "military",
		},
		{
			name:     "Unknown category",
			category: NewsCategoryUnknown,
			expected: "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := tt.category.String()

			assert.Equal(t, tt.expected, result)
		})
	}
}
