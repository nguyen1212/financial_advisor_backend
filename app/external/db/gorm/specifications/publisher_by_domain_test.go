package specifications

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Test_publisherByDomain_Query(t *testing.T) {
	t.Parallel()

	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{DryRun: true})
	if err != nil {
		t.Errorf("open db test connection failed: %v", err)
		return
	}

	tests := []struct {
		name    string
		spec    *publisherByDomain
		wantSQL string
	}{
		{
			name: "filter by domain",
			spec: &publisherByDomain{
				domain: "example.com",
			},
			wantSQL: "SELECT * FROM `publishers` WHERE domain = \"example.com\"",
		},
		{
			name: "empty domain",
			spec: &publisherByDomain{
				domain: "",
			},
			wantSQL: "SELECT * FROM `publishers` WHERE domain = \"\"",
		},
		{
			name: "complex domain",
			spec: &publisherByDomain{
				domain: "https://www.reuters.com",
			},
			wantSQL: "SELECT * FROM `publishers` WHERE domain = \"https://www.reuters.com\"",
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

func TestNewPublisherByDomain(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		domain string
		want   *publisherByDomain
	}{
		{
			name:   "create spec with domain",
			domain: "example.com",
			want: &publisherByDomain{
				domain: "example.com",
			},
		},
		{
			name:   "create spec with empty domain",
			domain: "",
			want: &publisherByDomain{
				domain: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := NewPublisherByDomain(tt.domain).(*publisherByDomain)
			assert.Equal(t, tt.want, result)
		})
	}
}