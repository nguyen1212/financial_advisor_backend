package specifications

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Test_publisherByID_Query(t *testing.T) {
	t.Parallel()

	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{DryRun: true})
	if err != nil {
		t.Errorf("open db test connection failed: %v", err)
		return
	}

	tests := []struct {
		name    string
		spec    *publisherByID
		wantSQL string
	}{
		{
			name: "filter by ID",
			spec: &publisherByID{
				id: 123,
			},
			wantSQL: "SELECT * FROM `publishers` WHERE id = 123",
		},
		{
			name: "filter by different ID",
			spec: &publisherByID{
				id: 456,
			},
			wantSQL: "SELECT * FROM `publishers` WHERE id = 456",
		},
		{
			name: "zero ID",
			spec: &publisherByID{
				id: 0,
			},
			wantSQL: "SELECT * FROM `publishers` WHERE id = 0",
		},
		{
			name: "large ID",
			spec: &publisherByID{
				id: 999999999,
			},
			wantSQL: "SELECT * FROM `publishers` WHERE id = 999999999",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			stmt := tt.spec.Query(db).Take(nil).Statement
			sql := db.Explain(stmt.SQL.String(), stmt.Vars...)
			sql = strings.Trim(sql, "LIMIT 1")

			assert.Equal(t, tt.wantSQL, sql)
		})
	}
}

func TestNewPublisherByID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		id   uint64
		want *publisherByID
	}{
		{
			name: "create spec with ID",
			id:   123,
			want: &publisherByID{
				id: 123,
			},
		},
		{
			name: "create spec with zero ID",
			id:   0,
			want: &publisherByID{
				id: 0,
			},
		},
		{
			name: "create spec with large ID",
			id:   999999999,
			want: &publisherByID{
				id: 999999999,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := NewPublisherByID(tt.id).(*publisherByID)
			assert.Equal(t, tt.want, result)
		})
	}
}