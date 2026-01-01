package specifications

import (
	"testing"

	"github.com/doug-martin/goqu/v9"
	"github.com/financial_advisor/app/domain/repository/specifications"
	"github.com/stretchr/testify/assert"
)

func Test_NewPublishersByNone(t *testing.T) {
	var (
		got  = NewPublishersByNone()
		want = publishersByNone{
			dialect: goqu.Dialect("mysql"),
		}
	)

	assert.Equal(t, want.dialect, got.(publishersByNone).dialect)
}

func Test_publishersByNone_ToGet(t *testing.T) {
	tests := []struct {
		name      string
		spec      publishersByNone
		wantQuery string
	}{
		{
			name: "basic query with no filters",
			spec: publishersByNone{
				dialect: goqu.Dialect("mysql"),
			},
			wantQuery: "SELECT * FROM `publishers` LIMIT 1",
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

func Test_publishersByNone_ToCount(t *testing.T) {
	tests := []struct {
		name      string
		spec      publishersByNone
		wantQuery string
	}{
		{
			name: "basic count query with no filters",
			spec: publishersByNone{
				dialect: goqu.Dialect("mysql"),
			},
			wantQuery: "SELECT COUNT(*) FROM `publishers`",
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

func Test_publishersByNone_ToFind(t *testing.T) {
	tests := []struct {
		name      string
		spec      publishersByNone
		paging    specifications.PagingI
		wantQuery string
	}{
		{
			name: "with paging",
			spec: publishersByNone{
				dialect: goqu.Dialect("mysql"),
			},
			paging:    paging{limit: 50, offset: 100},
			wantQuery: "SELECT * FROM `publishers` LIMIT 50 OFFSET 100",
		},
		{
			name: "without paging",
			spec: publishersByNone{
				dialect: goqu.Dialect("mysql"),
			},
			paging:    nil,
			wantQuery: "SELECT * FROM `publishers`",
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