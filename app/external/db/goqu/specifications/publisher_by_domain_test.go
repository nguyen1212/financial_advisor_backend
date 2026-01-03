package specifications

import (
	"testing"

	"github.com/doug-martin/goqu/v9"
	"github.com/financial_advisor/app/domain/repository/specifications"
	"github.com/stretchr/testify/assert"
)

func Test_NewPublisherByDomain(t *testing.T) {
	var (
		domain = "example.com"
		got    = NewPublisherByDomain(domain)
		want   = publisherByDomain{
			domain:  domain,
			dialect: goqu.Dialect("mysql"),
		}
	)

	assert.Equal(t, want.dialect, got.(publisherByDomain).dialect)
	assert.Equal(t, want.domain, got.(publisherByDomain).domain)
}

func Test_publisherByDomain_ToGet(t *testing.T) {
	tests := []struct {
		name      string
		spec      publisherByDomain
		wantQuery string
	}{
		{
			name: "basic domain query",
			spec: publisherByDomain{
				domain:  "example.com",
				dialect: goqu.Dialect("mysql"),
			},
			wantQuery: "SELECT * FROM `publishers` WHERE (`domain` = 'example.com') LIMIT 1",
		},
		{
			name: "empty domain",
			spec: publisherByDomain{
				domain:  "",
				dialect: goqu.Dialect("mysql"),
			},
			wantQuery: "SELECT * FROM `publishers` WHERE (`domain` = '') LIMIT 1",
		},
		{
			name: "domain with special characters",
			spec: publisherByDomain{
				domain:  "test-site.co.uk",
				dialect: goqu.Dialect("mysql"),
			},
			wantQuery: "SELECT * FROM `publishers` WHERE (`domain` = 'test-site.co.uk') LIMIT 1",
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

func Test_publisherByDomain_ToCount(t *testing.T) {
	tests := []struct {
		name      string
		spec      publisherByDomain
		wantQuery string
	}{
		{
			name: "basic domain query",
			spec: publisherByDomain{
				domain:  "example.com",
				dialect: goqu.Dialect("mysql"),
			},
			wantQuery: "SELECT COUNT(*) FROM `publishers` WHERE (`domain` = 'example.com')",
		},
		{
			name: "empty domain",
			spec: publisherByDomain{
				domain:  "",
				dialect: goqu.Dialect("mysql"),
			},
			wantQuery: "SELECT COUNT(*) FROM `publishers` WHERE (`domain` = '')",
		},
		{
			name: "domain with special characters",
			spec: publisherByDomain{
				domain:  "test-site.co.uk",
				dialect: goqu.Dialect("mysql"),
			},
			wantQuery: "SELECT COUNT(*) FROM `publishers` WHERE (`domain` = 'test-site.co.uk')",
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

func Test_publisherByDomain_ToFind(t *testing.T) {
	tests := []struct {
		name      string
		spec      publisherByDomain
		paging    specifications.PagingI
		wantQuery string
	}{
		{
			name: "with paging",
			spec: publisherByDomain{
				domain:  "example.com",
				dialect: goqu.Dialect("mysql"),
			},
			paging:    paging{limit: 20, offset: 40},
			wantQuery: "SELECT * FROM `publishers` WHERE (`domain` = 'example.com') LIMIT 20 OFFSET 40",
		},
		{
			name: "without paging",
			spec: publisherByDomain{
				domain:  "example.com",
				dialect: goqu.Dialect("mysql"),
			},
			paging:    nil,
			wantQuery: "SELECT * FROM `publishers` WHERE (`domain` = 'example.com')",
		},
		{
			name: "empty domain with paging",
			spec: publisherByDomain{
				domain:  "",
				dialect: goqu.Dialect("mysql"),
			},
			paging:    paging{limit: 15, offset: 30},
			wantQuery: "SELECT * FROM `publishers` WHERE (`domain` = '') LIMIT 15 OFFSET 30",
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