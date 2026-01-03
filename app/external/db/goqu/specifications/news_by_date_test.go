package specifications

import (
	"testing"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/financial_advisor/app/domain/entity"
	"github.com/financial_advisor/app/domain/repository/specifications"
	"github.com/stretchr/testify/assert"
)

func Test_NewsByDate(t *testing.T) {
	var (
		startDate = time.Now()
		endDate   = startDate.Add(24 * time.Hour)
		status    = entity.NewsStatusAdded

		got  = NewsByDate(startDate, endDate, &status)
		want = &newsByDate{
			startDate: startDate,
			endDate:   endDate,
			status:    &status,
			dialect:   goqu.Dialect("mysql"),
		}
	)

	assert.Equal(t, want.dialect, got.(*newsByDate).dialect)
	assert.Equal(t, want.startDate, got.(*newsByDate).startDate)
	assert.Equal(t, want.endDate, got.(*newsByDate).endDate)
	assert.Equal(t, want.status, got.(*newsByDate).status)
}

func Test_newsByDate_ToGet(t *testing.T) {
	tests := []struct {
		name      string
		spec      newsByDate
		wantQuery string
	}{
		{
			name: "with all fields",
			spec: newsByDate{
				startDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				endDate:   time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC),
				status:    &[]entity.NewsStatus{entity.NewsStatusAdded}[0],
				dialect:   goqu.Dialect("mysql"),
			},
			wantQuery: "SELECT * FROM `news` WHERE ((`created_at` >= '2024-01-01 00:00:00') AND (`created_at` <= '2024-01-31 23:59:59') AND (`status` = 1)) LIMIT 1",
		},
		{
			name: "start date is zero",
			spec: newsByDate{
				startDate: time.Time{},
				endDate:   time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC),
				status:    &[]entity.NewsStatus{entity.NewsStatusAdded}[0],
				dialect:   goqu.Dialect("mysql"),
			},
			wantQuery: "SELECT * FROM `news` WHERE ((`created_at` <= '2024-01-31 23:59:59') AND (`status` = 1)) LIMIT 1",
		},
		{
			name: "start date is zero",
			spec: newsByDate{
				startDate: time.Time{},
				endDate:   time.Time{},
				status:    &[]entity.NewsStatus{entity.NewsStatusAdded}[0],
				dialect:   goqu.Dialect("mysql"),
			},
			wantQuery: "SELECT * FROM `news` WHERE (`status` = 1) LIMIT 1",
		},
		{
			name: "status is nil",
			spec: newsByDate{
				startDate: time.Time{},
				endDate:   time.Time{},
				status:    nil,
				dialect:   goqu.Dialect("mysql"),
			},
			wantQuery: "SELECT * FROM `news` LIMIT 1",
		},
	}

	for i := range tests {
		t.Run(tests[i].name, func(t *testing.T) {
			t.Parallel()

			tests[i].spec.base.buildQuery = tests[i].spec.buildQuery

			gotQuery, err := tests[i].spec.ToGet()
			if err != nil {
				t.Errorf("ToGet() error = %v", err)

				return
			}

			assert.Equal(t, tests[i].wantQuery, gotQuery)
		})
	}
}

func Test_newsByDate_ToCount(t *testing.T) {
	tests := []struct {
		name      string
		spec      newsByDate
		wantQuery string
	}{
		{
			name: "with all fields",
			spec: newsByDate{
				startDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				endDate:   time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC),
				status:    &[]entity.NewsStatus{entity.NewsStatusAdded}[0],
				dialect:   goqu.Dialect("mysql"),
			},
			wantQuery: "SELECT COUNT(*) FROM `news` WHERE ((`created_at` >= '2024-01-01 00:00:00') AND (`created_at` <= '2024-01-31 23:59:59') AND (`status` = 1))",
		},
		{
			name: "start date is zero",
			spec: newsByDate{
				startDate: time.Time{},
				endDate:   time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC),
				status:    &[]entity.NewsStatus{entity.NewsStatusAdded}[0],
				dialect:   goqu.Dialect("mysql"),
			},
			wantQuery: "SELECT COUNT(*) FROM `news` WHERE ((`created_at` <= '2024-01-31 23:59:59') AND (`status` = 1))",
		},
		{
			name: "start date is zero",
			spec: newsByDate{
				startDate: time.Time{},
				endDate:   time.Time{},
				status:    &[]entity.NewsStatus{entity.NewsStatusAdded}[0],
				dialect:   goqu.Dialect("mysql"),
			},
			wantQuery: "SELECT COUNT(*) FROM `news` WHERE (`status` = 1)",
		},
		{
			name: "status is nil",
			spec: newsByDate{
				startDate: time.Time{},
				endDate:   time.Time{},
				status:    nil,
				dialect:   goqu.Dialect("mysql"),
			},
			wantQuery: "SELECT COUNT(*) FROM `news`",
		},
	}

	for i := range tests {
		t.Run(tests[i].name, func(t *testing.T) {
			t.Parallel()

			tests[i].spec.base.buildQuery = tests[i].spec.buildQuery

			gotQuery, err := tests[i].spec.ToCount()
			if err != nil {
				t.Errorf("ToCount() error = %v", err)

				return
			}

			assert.Equal(t, tests[i].wantQuery, gotQuery)
		})
	}
}

func Test_newsByDate_ToFind(t *testing.T) {
	tests := []struct {
		name      string
		spec      newsByDate
		paging    specifications.PagingI
		wantQuery string
	}{
		{
			name: "with all fields",
			spec: newsByDate{
				startDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				endDate:   time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC),
				status:    &[]entity.NewsStatus{entity.NewsStatusAdded}[0],
				dialect:   goqu.Dialect("mysql"),
			},
			paging:    paging{limit: 30, offset: 60},
			wantQuery: "SELECT * FROM `news` WHERE ((`created_at` >= '2024-01-01 00:00:00') AND (`created_at` <= '2024-01-31 23:59:59') AND (`status` = 1)) LIMIT 30 OFFSET 60",
		},
		{
			name: "start date is zero",
			spec: newsByDate{
				startDate: time.Time{},
				endDate:   time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC),
				status:    &[]entity.NewsStatus{entity.NewsStatusAdded}[0],
				dialect:   goqu.Dialect("mysql"),
			},
			paging:    paging{limit: 30, offset: 60},
			wantQuery: "SELECT * FROM `news` WHERE ((`created_at` <= '2024-01-31 23:59:59') AND (`status` = 1)) LIMIT 30 OFFSET 60",
		},
		{
			name: "start date is zero",
			spec: newsByDate{
				startDate: time.Time{},
				endDate:   time.Time{},
				status:    &[]entity.NewsStatus{entity.NewsStatusAdded}[0],
				dialect:   goqu.Dialect("mysql"),
			},
			paging:    paging{limit: 30, offset: 60},
			wantQuery: "SELECT * FROM `news` WHERE (`status` = 1) LIMIT 30 OFFSET 60",
		},
		{
			name: "status is nil",
			spec: newsByDate{
				startDate: time.Time{},
				endDate:   time.Time{},
				status:    nil,
				dialect:   goqu.Dialect("mysql"),
			},
			paging:    paging{limit: 30, offset: 60},
			wantQuery: "SELECT * FROM `news` LIMIT 30 OFFSET 60",
		},
		{
			name: "status is nil",
			spec: newsByDate{
				startDate: time.Time{},
				endDate:   time.Time{},
				status:    nil,
				dialect:   goqu.Dialect("mysql"),
			},
			paging:    nil,
			wantQuery: "SELECT * FROM `news`",
		},
	}

	for i := range tests {
		t.Run(tests[i].name, func(t *testing.T) {
			t.Parallel()

			tests[i].spec.base.buildQuery = tests[i].spec.buildQuery

			gotQuery, err := tests[i].spec.ToFind(tests[i].paging)
			if err != nil {
				t.Errorf("ToFind() error = %v", err)

				return
			}

			assert.Equal(t, tests[i].wantQuery, gotQuery)
		})
	}
}
