package specifications

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Test_newsByDate_ToSQL(t *testing.T) {
	t.Parallel()

	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{DryRun: true})
	if err != nil {
		t.Errorf("open db test connection failed: %v", err)

		return
	}

	tests := []struct {
		name    string
		spec    *newsByDate
		wantSQL string
	}{
		{
			name:    "no date provided",
			spec:    &newsByDate{},
			wantSQL: "SELECT * FROM `news`",
		},
		{
			name: "from date provided",
			spec: &newsByDate{
				startDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			wantSQL: "SELECT * FROM `news` WHERE created_at >= \"2024-01-01 00:00:00\"",
		},
		{
			name: "to date provided",
			spec: &newsByDate{
				endDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			wantSQL: "SELECT * FROM `news` WHERE created_at <= \"2024-01-01 00:00:00\"",
		},
	}

	for i := range tests {
		t.Run(tests[i].name, func(t *testing.T) {
			t.Parallel()

			stmt := tests[i].spec.Query(db).Take(nil).Statement
			sql := db.Explain(stmt.SQL.String(), stmt.Vars...)
			sql = strings.Trim(sql, "LIMIT 1")

			assert.Equal(t, tests[i].wantSQL, sql)
		})
	}
}

func TestNewNewsByDate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		startDate time.Time
		endDate   time.Time
		want      *newsByDate
	}{
		{
			name:      "with both dates",
			startDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			endDate:   time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC),
			want: &newsByDate{
				startDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				endDate:   time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC),
			},
		},
		{
			name:      "with zero dates",
			startDate: time.Time{},
			endDate:   time.Time{},
			want: &newsByDate{
				startDate: time.Time{},
				endDate:   time.Time{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := NewNewsByDate(tt.startDate, tt.endDate).(*newsByDate)
			assert.Equal(t, tt.want, result)
		})
	}
}
