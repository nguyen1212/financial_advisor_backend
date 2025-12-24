package specifications

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Test_newsByID_Query(t *testing.T) {
	t.Parallel()

	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{DryRun: true})
	if err != nil {
		t.Errorf("open db test connection failed: %v", err)
		return
	}

	tests := []struct {
		name    string
		spec    newsByID
		wantSQL string
	}{
		{
			name: "filter by ID without preloaders",
			spec: newsByID{
				id:         123,
				preloaders: []string{},
			},
			wantSQL: "SELECT * FROM `news` WHERE id = 123",
		},
		{
			name: "filter by ID with single preloader",
			spec: newsByID{
				id:         456,
				preloaders: []string{"Publisher"},
			},
			wantSQL: "SELECT * FROM `news` WHERE id = 456",
		},
		{
			name: "filter by ID with multiple preloaders",
			spec: newsByID{
				id:         789,
				preloaders: []string{"Publisher", "Category"},
			},
			wantSQL: "SELECT * FROM `news` WHERE id = 789",
		},
		{
			name: "zero ID",
			spec: newsByID{
				id:         0,
				preloaders: []string{},
			},
			wantSQL: "SELECT * FROM `news` WHERE id = 0",
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

func TestNewNewsByID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		id          uint64
		preloaders  []string
		want        newsByID
	}{
		{
			name:       "create spec without preloaders",
			id:         123,
			preloaders: []string{},
			want: newsByID{
				id:         123,
				preloaders: []string{},
			},
		},
		{
			name:       "create spec with single preloader",
			id:         456,
			preloaders: []string{"Publisher"},
			want: newsByID{
				id:         456,
				preloaders: []string{"Publisher"},
			},
		},
		{
			name:       "create spec with multiple preloaders",
			id:         789,
			preloaders: []string{"Publisher", "Category"},
			want: newsByID{
				id:         789,
				preloaders: []string{"Publisher", "Category"},
			},
		},
		{
			name:       "create spec with nil preloaders",
			id:         100,
			preloaders: nil,
			want: newsByID{
				id:         100,
				preloaders: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := NewNewsByID(tt.id, tt.preloaders...).(newsByID)
			assert.Equal(t, tt.want, result)
		})
	}
}