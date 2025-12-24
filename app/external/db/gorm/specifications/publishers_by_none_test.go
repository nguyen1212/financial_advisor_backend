package specifications

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Test_publishersByNone_Query(t *testing.T) {
	t.Parallel()

	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{DryRun: true})
	if err != nil {
		t.Errorf("open db test connection failed: %v", err)
		return
	}

	tests := []struct {
		name    string
		spec    *publishersByNone
		wantSQL string
	}{
		{
			name:    "no filter applied",
			spec:    &publishersByNone{},
			wantSQL: "SELECT * FROM `publishers`",
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

func TestNewPublishersByNone(t *testing.T) {
	t.Parallel()

	result := NewPublishersByNone().(*publishersByNone)
	expected := &publishersByNone{}

	assert.Equal(t, expected, result)
}