package specifications

import (
	"testing"

	"github.com/doug-martin/goqu/v9"
	"github.com/financial_advisor/app/domain/repository/specifications"
	"github.com/stretchr/testify/assert"
)

func Test_NewNewsByID(t *testing.T) {
	var (
		id   = uint64(123)
		got  = NewNewsByID(id)
		want = newsByID{
			id:      id,
			dialect: goqu.Dialect("mysql"),
		}
	)

	assert.Equal(t, want.dialect, got.(newsByID).dialect)
	assert.Equal(t, want.id, got.(newsByID).id)
}

func Test_newsByID_ToGet(t *testing.T) {
	tests := []struct {
		name      string
		spec      newsByID
		wantQuery string
	}{
		{
			name: "basic query",
			spec: newsByID{
				id:      123,
				dialect: goqu.Dialect("mysql"),
			},
			wantQuery: "SELECT * FROM `news` WHERE (`id` = 123) LIMIT 1",
		},
		{
			name: "zero id",
			spec: newsByID{
				id:      0,
				dialect: goqu.Dialect("mysql"),
			},
			wantQuery: "SELECT * FROM `news` WHERE (`id` = 0) LIMIT 1",
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

func Test_newsByID_ToCount(t *testing.T) {
	tests := []struct {
		name      string
		spec      newsByID
		wantQuery string
	}{
		{
			name: "basic query",
			spec: newsByID{
				id:      123,
				dialect: goqu.Dialect("mysql"),
			},
			wantQuery: "SELECT COUNT(*) FROM `news` WHERE (`id` = 123)",
		},
		{
			name: "zero id",
			spec: newsByID{
				id:      0,
				dialect: goqu.Dialect("mysql"),
			},
			wantQuery: "SELECT COUNT(*) FROM `news` WHERE (`id` = 0)",
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

func Test_newsByID_ToFind(t *testing.T) {
	tests := []struct {
		name      string
		spec      newsByID
		paging    specifications.PagingI
		wantQuery string
	}{
		{
			name: "with paging",
			spec: newsByID{
				id:      123,
				dialect: goqu.Dialect("mysql"),
			},
			paging:    paging{limit: 30, offset: 60},
			wantQuery: "SELECT * FROM `news` WHERE (`id` = 123) LIMIT 30 OFFSET 60",
		},
		{
			name: "without paging",
			spec: newsByID{
				id:      123,
				dialect: goqu.Dialect("mysql"),
			},
			paging:    nil,
			wantQuery: "SELECT * FROM `news` WHERE (`id` = 123)",
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