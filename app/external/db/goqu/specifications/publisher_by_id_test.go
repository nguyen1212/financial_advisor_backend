package specifications

import (
	"testing"

	"github.com/doug-martin/goqu/v9"
	"github.com/financial_advisor/app/domain/repository/specifications"
	"github.com/stretchr/testify/assert"
)

func Test_NewPublisherByID(t *testing.T) {
	var (
		id   = uint64(456)
		got  = NewPublisherByID(id)
		want = publisherByID{
			id:      id,
			dialect: goqu.Dialect("mysql"),
		}
	)

	assert.Equal(t, want.dialect, got.(publisherByID).dialect)
	assert.Equal(t, want.id, got.(publisherByID).id)
}

func Test_publisherByID_ToGet(t *testing.T) {
	tests := []struct {
		name      string
		spec      publisherByID
		wantQuery string
	}{
		{
			name: "basic query",
			spec: publisherByID{
				id:      456,
				dialect: goqu.Dialect("mysql"),
			},
			wantQuery: "SELECT * FROM `publishers` WHERE (`id` = 456) LIMIT 1",
		},
		{
			name: "zero id",
			spec: publisherByID{
				id:      0,
				dialect: goqu.Dialect("mysql"),
			},
			wantQuery: "SELECT * FROM `publishers` WHERE (`id` = 0) LIMIT 1",
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

func Test_publisherByID_ToCount(t *testing.T) {
	tests := []struct {
		name      string
		spec      publisherByID
		wantQuery string
	}{
		{
			name: "basic query",
			spec: publisherByID{
				id:      456,
				dialect: goqu.Dialect("mysql"),
			},
			wantQuery: "SELECT COUNT(*) FROM `publishers` WHERE (`id` = 456)",
		},
		{
			name: "zero id",
			spec: publisherByID{
				id:      0,
				dialect: goqu.Dialect("mysql"),
			},
			wantQuery: "SELECT COUNT(*) FROM `publishers` WHERE (`id` = 0)",
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

func Test_publisherByID_ToFind(t *testing.T) {
	tests := []struct {
		name      string
		spec      publisherByID
		paging    specifications.PagingI
		wantQuery string
	}{
		{
			name: "with paging",
			spec: publisherByID{
				id:      456,
				dialect: goqu.Dialect("mysql"),
			},
			paging:    paging{limit: 25, offset: 50},
			wantQuery: "SELECT * FROM `publishers` WHERE (`id` = 456) LIMIT 25 OFFSET 50",
		},
		{
			name: "without paging",
			spec: publisherByID{
				id:      456,
				dialect: goqu.Dialect("mysql"),
			},
			paging:    nil,
			wantQuery: "SELECT * FROM `publishers` WHERE (`id` = 456)",
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