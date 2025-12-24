package specifications

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPaging(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		limit  int
		offset int
		want   Paging
	}{
		{
			name:   "basic paging",
			limit:  10,
			offset: 20,
			want: Paging{
				limit:  10,
				offset: 20,
			},
		},
		{
			name:   "zero values",
			limit:  0,
			offset: 0,
			want: Paging{
				limit:  0,
				offset: 0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := newPaging(tt.limit, tt.offset)
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestPaging_Limit(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		paging Paging
		want   int
	}{
		{
			name:   "return limit value",
			paging: Paging{limit: 25, offset: 50},
			want:   25,
		},
		{
			name:   "zero limit",
			paging: Paging{limit: 0, offset: 10},
			want:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := tt.paging.Limit()
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestPaging_Offset(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		paging Paging
		want   int
	}{
		{
			name:   "return offset value",
			paging: Paging{limit: 25, offset: 50},
			want:   50,
		},
		{
			name:   "zero offset",
			paging: Paging{limit: 10, offset: 0},
			want:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := tt.paging.Offset()
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestToPaging(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		size int
		page int
		want Paging
	}{
		{
			name: "first page",
			size: 10,
			page: 1,
			want: Paging{
				limit:  10,
				offset: 0,
			},
		},
		{
			name: "second page",
			size: 10,
			page: 2,
			want: Paging{
				limit:  10,
				offset: 10,
			},
		},
		{
			name: "third page with different size",
			size: 20,
			page: 3,
			want: Paging{
				limit:  20,
				offset: 40,
			},
		},
		{
			name: "zero page",
			size: 10,
			page: 0,
			want: Paging{
				limit:  10,
				offset: -10,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := ToPaging(tt.size, tt.page)
			assert.Equal(t, tt.want, result)
		})
	}
}

