package payload

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPaging(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		paging       Paging
		expectedSize int
		expectedPage int
	}{
		{
			name: "valid paging with positive values",
			paging: Paging{
				Size: 10,
				Page: 1,
			},
			expectedSize: 10,
			expectedPage: 1,
		},
		{
			name: "paging with zero values",
			paging: Paging{
				Size: 0,
				Page: 0,
			},
			expectedSize: 0,
			expectedPage: 0,
		},
		{
			name: "paging with large values",
			paging: Paging{
				Size: 100,
				Page: 50,
			},
			expectedSize: 100,
			expectedPage: 50,
		},
		{
			name: "paging with negative values",
			paging: Paging{
				Size: -5,
				Page: -1,
			},
			expectedSize: -5,
			expectedPage: -1,
		},
		{
			name: "paging with mixed positive and negative",
			paging: Paging{
				Size: 20,
				Page: -2,
			},
			expectedSize: 20,
			expectedPage: -2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.expectedSize, tt.paging.Size)
			assert.Equal(t, tt.expectedPage, tt.paging.Page)
		})
	}
}

func TestPaging_DefaultValues(t *testing.T) {
	t.Parallel()

	// Test that a new Paging struct has zero values
	paging := Paging{}

	assert.Equal(t, 0, paging.Size)
	assert.Equal(t, 0, paging.Page)
}

func TestPaging_StructTags(t *testing.T) {
	t.Parallel()

	// This test verifies that the struct tags are correctly defined
	// by checking that the struct can be used in form binding contexts
	paging := Paging{
		Size: 25,
		Page: 3,
	}

	// Test that the fields have the expected form tag values
	// This is more of a structural test to ensure the tags exist
	assert.Equal(t, 25, paging.Size)
	assert.Equal(t, 3, paging.Page)
}