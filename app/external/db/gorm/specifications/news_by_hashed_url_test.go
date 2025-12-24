package specifications

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Test_newsByHashedURL_Query(t *testing.T) {
	t.Parallel()

	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{DryRun: true})
	if err != nil {
		t.Errorf("open db test connection failed: %v", err)
		return
	}

	tests := []struct {
		name      string
		spec      newsByHashedURL
		wantSQL   string
	}{
		{
			name: "filter by hashed URL",
			spec: newsByHashedURL{
				hashedURL: []byte("test-hash"),
			},
			wantSQL: "SELECT * FROM `news` WHERE hashed_url = \"test-hash\"",
		},
		{
			name: "empty hashed URL",
			spec: newsByHashedURL{
				hashedURL: []byte{},
			},
			wantSQL: "SELECT * FROM `news` WHERE hashed_url = \"\"",
		},
		{
			name: "complex hashed URL",
			spec: newsByHashedURL{
				hashedURL: []byte("complex-hash-with-special-chars-123"),
			},
			wantSQL: "SELECT * FROM `news` WHERE hashed_url = \"complex-hash-with-special-chars-123\"",
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

func TestNewNewsByHashedURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		hashedURL []byte
		want      newsByHashedURL
	}{
		{
			name:      "create spec with hashed URL",
			hashedURL: []byte("test-hash"),
			want: newsByHashedURL{
				hashedURL: []byte("test-hash"),
			},
		},
		{
			name:      "create spec with empty hashed URL",
			hashedURL: []byte{},
			want: newsByHashedURL{
				hashedURL: []byte{},
			},
		},
		{
			name:      "create spec with nil hashed URL",
			hashedURL: nil,
			want: newsByHashedURL{
				hashedURL: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := NewNewsByHashedURL(tt.hashedURL).(newsByHashedURL)
			assert.Equal(t, tt.want, result)
		})
	}
}