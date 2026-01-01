package specifications

import (
	"testing"

	"github.com/doug-martin/goqu/v9"
	"github.com/financial_advisor/app/domain/repository/specifications"
	"github.com/stretchr/testify/assert"
)

func Test_NewNewsByIDs(t *testing.T) {
	var (
		ids  = []uint64{1, 2, 3}
		got  = NewNewsByIDs(ids)
		want = newsByIDs{
			ids:     ids,
			dialect: goqu.Dialect("mysql"),
		}
	)

	assert.Equal(t, want.dialect, got.(newsByIDs).dialect)
	assert.Equal(t, want.ids, got.(newsByIDs).ids)
}

func Test_newsByIDs_ToGet(t *testing.T) {
	tests := []struct {
		name      string
		spec      newsByIDs
		wantQuery string
	}{
		{
			name: "multiple ids",
			spec: newsByIDs{
				ids:     []uint64{1, 2, 3},
				dialect: goqu.Dialect("mysql"),
			},
			wantQuery: "SELECT * FROM `news` WHERE (`id` IN (1, 2, 3)) LIMIT 1",
		},
		{
			name: "single id",
			spec: newsByIDs{
				ids:     []uint64{123},
				dialect: goqu.Dialect("mysql"),
			},
			wantQuery: "SELECT * FROM `news` WHERE (`id` IN (123)) LIMIT 1",
		},
		{
			name: "empty ids",
			spec: newsByIDs{
				ids:     []uint64{},
				dialect: goqu.Dialect("mysql"),
			},
			wantQuery: "SELECT * FROM `news` WHERE (`id` IN ()) LIMIT 1",
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

func Test_newsByIDs_ToCount(t *testing.T) {
	tests := []struct {
		name      string
		spec      newsByIDs
		wantQuery string
	}{
		{
			name: "multiple ids",
			spec: newsByIDs{
				ids:     []uint64{1, 2, 3},
				dialect: goqu.Dialect("mysql"),
			},
			wantQuery: "SELECT COUNT(*) FROM `news` WHERE (`id` IN (1, 2, 3))",
		},
		{
			name: "single id",
			spec: newsByIDs{
				ids:     []uint64{123},
				dialect: goqu.Dialect("mysql"),
			},
			wantQuery: "SELECT COUNT(*) FROM `news` WHERE (`id` IN (123))",
		},
		{
			name: "empty ids",
			spec: newsByIDs{
				ids:     []uint64{},
				dialect: goqu.Dialect("mysql"),
			},
			wantQuery: "SELECT COUNT(*) FROM `news` WHERE (`id` IN ())",
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

func Test_newsByIDs_ToFind(t *testing.T) {
	tests := []struct {
		name      string
		spec      newsByIDs
		paging    specifications.PagingI
		wantQuery string
	}{
		{
			name: "multiple ids with paging",
			spec: newsByIDs{
				ids:     []uint64{1, 2, 3},
				dialect: goqu.Dialect("mysql"),
			},
			paging:    paging{limit: 30, offset: 60},
			wantQuery: "SELECT * FROM `news` WHERE (`id` IN (1, 2, 3)) LIMIT 30 OFFSET 60",
		},
		{
			name: "multiple ids without paging",
			spec: newsByIDs{
				ids:     []uint64{1, 2, 3},
				dialect: goqu.Dialect("mysql"),
			},
			paging:    nil,
			wantQuery: "SELECT * FROM `news` WHERE (`id` IN (1, 2, 3))",
		},
		{
			name: "single id with paging",
			spec: newsByIDs{
				ids:     []uint64{123},
				dialect: goqu.Dialect("mysql"),
			},
			paging:    paging{limit: 10, offset: 20},
			wantQuery: "SELECT * FROM `news` WHERE (`id` IN (123)) LIMIT 10 OFFSET 20",
		},
		{
			name: "empty ids without paging",
			spec: newsByIDs{
				ids:     []uint64{},
				dialect: goqu.Dialect("mysql"),
			},
			paging:    nil,
			wantQuery: "SELECT * FROM `news` WHERE (`id` IN ())",
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