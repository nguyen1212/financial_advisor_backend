package specifications

import (
	"testing"

	"github.com/doug-martin/goqu/v9"
	"github.com/financial_advisor/app/domain/repository/specifications"
	"github.com/stretchr/testify/assert"
)

func Test_NewNewsByHashedURL(t *testing.T) {
	var (
		hashedURL = []byte("test-hash")
		got       = NewNewsByHashedURL(hashedURL)
		want      = newsByHashedURL{
			hashedURL: hashedURL,
			dialect:   goqu.Dialect("mysql"),
		}
	)

	assert.Equal(t, want.dialect, got.(newsByHashedURL).dialect)
	assert.Equal(t, want.hashedURL, got.(newsByHashedURL).hashedURL)
}

func Test_newsByHashedURL_ToGet(t *testing.T) {
	tests := []struct {
		name      string
		spec      newsByHashedURL
		wantQuery string
	}{
		{
			name: "with hashed URL",
			spec: newsByHashedURL{
				hashedURL: []byte("test-hash"),
				dialect:   goqu.Dialect("mysql"),
			},
			wantQuery: "SELECT * FROM `news` WHERE (`hashed_url` = 'test-hash') LIMIT 1",
		},
		{
			name: "with nil hashed URL",
			spec: newsByHashedURL{
				hashedURL: nil,
				dialect:   goqu.Dialect("mysql"),
			},
			wantQuery: "SELECT * FROM `news` LIMIT 1",
		},
		{
			name: "with empty hashed URL",
			spec: newsByHashedURL{
				hashedURL: []byte{},
				dialect:   goqu.Dialect("mysql"),
			},
			wantQuery: "SELECT * FROM `news` WHERE (`hashed_url` = '') LIMIT 1",
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

func Test_newsByHashedURL_ToCount(t *testing.T) {
	tests := []struct {
		name      string
		spec      newsByHashedURL
		wantQuery string
	}{
		{
			name: "with hashed URL",
			spec: newsByHashedURL{
				hashedURL: []byte("test-hash"),
				dialect:   goqu.Dialect("mysql"),
			},
			wantQuery: "SELECT COUNT(*) FROM `news` WHERE (`hashed_url` = 'test-hash')",
		},
		{
			name: "with nil hashed URL",
			spec: newsByHashedURL{
				hashedURL: nil,
				dialect:   goqu.Dialect("mysql"),
			},
			wantQuery: "SELECT COUNT(*) FROM `news`",
		},
		{
			name: "with empty hashed URL",
			spec: newsByHashedURL{
				hashedURL: []byte{},
				dialect:   goqu.Dialect("mysql"),
			},
			wantQuery: "SELECT COUNT(*) FROM `news` WHERE (`hashed_url` = '')",
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

func Test_newsByHashedURL_ToFind(t *testing.T) {
	tests := []struct {
		name      string
		spec      newsByHashedURL
		paging    specifications.PagingI
		wantQuery string
	}{
		{
			name: "with hashed URL and paging",
			spec: newsByHashedURL{
				hashedURL: []byte("test-hash"),
				dialect:   goqu.Dialect("mysql"),
			},
			paging:    paging{limit: 30, offset: 60},
			wantQuery: "SELECT * FROM `news` WHERE (`hashed_url` = 'test-hash') LIMIT 30 OFFSET 60",
		},
		{
			name: "with hashed URL without paging",
			spec: newsByHashedURL{
				hashedURL: []byte("test-hash"),
				dialect:   goqu.Dialect("mysql"),
			},
			paging:    nil,
			wantQuery: "SELECT * FROM `news` WHERE (`hashed_url` = 'test-hash')",
		},
		{
			name: "with nil hashed URL and paging",
			spec: newsByHashedURL{
				hashedURL: nil,
				dialect:   goqu.Dialect("mysql"),
			},
			paging:    paging{limit: 30, offset: 60},
			wantQuery: "SELECT * FROM `news` LIMIT 30 OFFSET 60",
		},
		{
			name: "with nil hashed URL without paging",
			spec: newsByHashedURL{
				hashedURL: nil,
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