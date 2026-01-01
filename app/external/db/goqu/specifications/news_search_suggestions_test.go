package specifications

import (
	"testing"

	"github.com/doug-martin/goqu/v9"
	"github.com/financial_advisor/app/domain/repository/specifications"
	"github.com/stretchr/testify/assert"
)

func Test_NewNewsSearchSuggestions(t *testing.T) {
	var (
		keywords  = []string{"search", "suggestion"}
		fuzziness = MediumFuzziness
		got       = NewNewsSearchSuggestions(keywords, fuzziness)
		want      = newsSearchSuggestions{
			keywords: keywords,
			fuziness: fuzziness,
			dialect:  goqu.Dialect("mysql"),
		}
	)

	actualSpec := got.(newsSearchSuggestions)
	assert.Equal(t, want.dialect, actualSpec.dialect)
	assert.Equal(t, want.keywords, actualSpec.keywords)
	assert.Equal(t, want.fuziness, actualSpec.fuziness)
}

func Test_newsSearchSuggestions_ToGet(t *testing.T) {
	tests := []struct {
		name      string
		spec      newsSearchSuggestions
		wantQuery string
	}{
		{
			name: "disabled fuzziness",
			spec: newsSearchSuggestions{
				keywords: []string{"golang", "programming"},
				fuziness: DisabledFuzziness,
				dialect:  goqu.Dialect("mysql"),
			},
			wantQuery: "SELECT * FROM `news` WHERE CALL AUTOCOMPLETE('golang programming', 'news', 0 AS fuzziness, 1 AS preserve) LIMIT 1",
		},
		{
			name: "strong fuzziness",
			spec: newsSearchSuggestions{
				keywords: []string{"machine", "learning"},
				fuziness: StrongFuzziness,
				dialect:  goqu.Dialect("mysql"),
			},
			wantQuery: "SELECT * FROM `news` WHERE CALL AUTOCOMPLETE('machine learning', 'news', 1 AS fuzziness, 1 AS preserve) LIMIT 1",
		},
		{
			name: "medium fuzziness",
			spec: newsSearchSuggestions{
				keywords: []string{"artificial", "intelligence"},
				fuziness: MediumFuzziness,
				dialect:  goqu.Dialect("mysql"),
			},
			wantQuery: "SELECT * FROM `news` WHERE CALL AUTOCOMPLETE('artificial intelligence', 'news', 2 AS fuzziness, 1 AS preserve) LIMIT 1",
		},
		{
			name: "single keyword with medium fuzziness",
			spec: newsSearchSuggestions{
				keywords: []string{"blockchain"},
				fuziness: MediumFuzziness,
				dialect:  goqu.Dialect("mysql"),
			},
			wantQuery: "SELECT * FROM `news` WHERE CALL AUTOCOMPLETE('blockchain', 'news', 2 AS fuzziness, 1 AS preserve) LIMIT 1",
		},
		{
			name: "empty keywords",
			spec: newsSearchSuggestions{
				keywords: []string{},
				fuziness: StrongFuzziness,
				dialect:  goqu.Dialect("mysql"),
			},
			wantQuery: "SELECT * FROM `news` WHERE CALL AUTOCOMPLETE('', 'news', 1 AS fuzziness, 1 AS preserve) LIMIT 1",
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

func Test_newsSearchSuggestions_ToCount(t *testing.T) {
	tests := []struct {
		name      string
		spec      newsSearchSuggestions
		wantQuery string
	}{
		{
			name: "disabled fuzziness count",
			spec: newsSearchSuggestions{
				keywords: []string{"golang", "programming"},
				fuziness: DisabledFuzziness,
				dialect:  goqu.Dialect("mysql"),
			},
			wantQuery: "SELECT COUNT(*) FROM `news` WHERE CALL AUTOCOMPLETE('golang programming', 'news', 0 AS fuzziness, 1 AS preserve)",
		},
		{
			name: "strong fuzziness count",
			spec: newsSearchSuggestions{
				keywords: []string{"machine", "learning"},
				fuziness: StrongFuzziness,
				dialect:  goqu.Dialect("mysql"),
			},
			wantQuery: "SELECT COUNT(*) FROM `news` WHERE CALL AUTOCOMPLETE('machine learning', 'news', 1 AS fuzziness, 1 AS preserve)",
		},
		{
			name: "medium fuzziness count",
			spec: newsSearchSuggestions{
				keywords: []string{"artificial", "intelligence"},
				fuziness: MediumFuzziness,
				dialect:  goqu.Dialect("mysql"),
			},
			wantQuery: "SELECT COUNT(*) FROM `news` WHERE CALL AUTOCOMPLETE('artificial intelligence', 'news', 2 AS fuzziness, 1 AS preserve)",
		},
		{
			name: "single keyword count",
			spec: newsSearchSuggestions{
				keywords: []string{"blockchain"},
				fuziness: MediumFuzziness,
				dialect:  goqu.Dialect("mysql"),
			},
			wantQuery: "SELECT COUNT(*) FROM `news` WHERE CALL AUTOCOMPLETE('blockchain', 'news', 2 AS fuzziness, 1 AS preserve)",
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

func Test_newsSearchSuggestions_ToFind(t *testing.T) {
	tests := []struct {
		name      string
		spec      newsSearchSuggestions
		paging    specifications.PagingI
		wantQuery string
	}{
		{
			name: "with paging and medium fuzziness",
			spec: newsSearchSuggestions{
				keywords: []string{"tech", "news"},
				fuziness: MediumFuzziness,
				dialect:  goqu.Dialect("mysql"),
			},
			paging:    paging{limit: 10, offset: 20},
			wantQuery: "CALL AUTOCOMPLETE('tech news', 'news', 2 AS fuzziness, 1 AS preserve) LIMIT 10 OFFSET 20",
		},
		{
			name: "without paging and strong fuzziness",
			spec: newsSearchSuggestions{
				keywords: []string{"finance", "market"},
				fuziness: StrongFuzziness,
				dialect:  goqu.Dialect("mysql"),
			},
			paging:    nil,
			wantQuery: "CALL AUTOCOMPLETE('finance market', 'news', 1 AS fuzziness, 1 AS preserve)",
		},
		{
			name: "with paging and disabled fuzziness",
			spec: newsSearchSuggestions{
				keywords: []string{"economy"},
				fuziness: DisabledFuzziness,
				dialect:  goqu.Dialect("mysql"),
			},
			paging:    paging{limit: 5, offset: 0},
			wantQuery: "CALL AUTOCOMPLETE('economy', 'news', 0 AS fuzziness, 1 AS preserve) LIMIT 5",
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