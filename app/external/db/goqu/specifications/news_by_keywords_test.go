package specifications

import (
	"testing"

	"github.com/doug-martin/goqu/v9"
	"github.com/financial_advisor/app/domain/repository/specifications"
	"github.com/stretchr/testify/assert"
)

func Test_NewNewsWithFullTextByKeywords(t *testing.T) {
	var (
		keywords         = []string{"test", "search"}
		previewLimit     = uint64(512)
		fullTextSearchOp = FullTextSearchOpAnd
		got              = NewNewsWithFullTextByKeywords(keywords, previewLimit, fullTextSearchOp)
		want             = newsWithFullTextByKeywords{
			keywords:         keywords,
			previewLimit:     previewLimit,
			fullTextSearchOp: fullTextSearchOp,
			dialect:          goqu.Dialect("mysql"),
		}
	)

	actualSpec := got.(newsWithFullTextByKeywords)
	assert.Equal(t, want.dialect, actualSpec.dialect)
	assert.Equal(t, want.keywords, actualSpec.keywords)
	assert.Equal(t, want.previewLimit, actualSpec.previewLimit)
	assert.Equal(t, want.fullTextSearchOp, actualSpec.fullTextSearchOp)
}

func Test_newsWithFullTextByKeywords_ToFind(t *testing.T) {
	tests := []struct {
		name      string
		spec      newsWithFullTextByKeywords
		paging    specifications.PagingI
		wantQuery string
	}{
		{
			name: "AND operation with paging",
			spec: newsWithFullTextByKeywords{
				keywords:         []string{"golang", "programming"},
				previewLimit:     512,
				fullTextSearchOp: FullTextSearchOpAnd,
				dialect:          goqu.Dialect("mysql"),
			},
			paging:    paging{limit: 10, offset: 20},
			wantQuery: "SELECT HIGHLIGHT({limit=512,before_match='<mark>',after_match='</mark>'}, content) AS content, news_id, id FROM `news` WHERE MATCH('@content golang programming') LIMIT 10 OFFSET 20",
		},
		{
			name: "Proximity operation without paging",
			spec: newsWithFullTextByKeywords{
				keywords:         []string{"machine", "learning"},
				previewLimit:     256,
				fullTextSearchOp: FullTextSearchOpProximity,
				dialect:          goqu.Dialect("mysql"),
			},
			paging:    nil,
			wantQuery: "SELECT HIGHLIGHT({limit=256,before_match='<mark>',after_match='</mark>'}, content) AS content, news_id, id FROM `news` WHERE MATCH('@content \\\"machine learning\\\"~10')",
		},
		{
			name: "Quorum operation with paging",
			spec: newsWithFullTextByKeywords{
				keywords:         []string{"artificial", "intelligence"},
				previewLimit:     1024,
				fullTextSearchOp: FullTextSearchOpQuorum,
				dialect:          goqu.Dialect("mysql"),
			},
			paging:    paging{limit: 5, offset: 15},
			wantQuery: "SELECT HIGHLIGHT({limit=1024,before_match='<mark>',after_match='</mark>'}, content) AS content, news_id, id FROM `news` WHERE MATCH('@content \\\"artificial intelligence\\\"/0.3') LIMIT 5 OFFSET 15",
		},
		{
			name: "Zero preview limit should default to 256",
			spec: newsWithFullTextByKeywords{
				keywords:         []string{"technology"},
				previewLimit:     0,
				fullTextSearchOp: FullTextSearchOpAnd,
				dialect:          goqu.Dialect("mysql"),
			},
			paging:    nil,
			wantQuery: "SELECT HIGHLIGHT({limit=256,before_match='<mark>',after_match='</mark>'}, content) AS content, news_id, id FROM `news` WHERE MATCH('@content technology')",
		},
		{
			name: "Single keyword with AND operation",
			spec: newsWithFullTextByKeywords{
				keywords:         []string{"blockchain"},
				previewLimit:     128,
				fullTextSearchOp: FullTextSearchOpAnd,
				dialect:          goqu.Dialect("mysql"),
			},
			paging:    paging{limit: 20, offset: 0},
			wantQuery: "SELECT HIGHLIGHT({limit=128,before_match='<mark>',after_match='</mark>'}, content) AS content, news_id, id FROM `news` WHERE MATCH('@content blockchain') LIMIT 20",
		},
	}

	for i := range tests {
		t.Run(tests[i].name, func(t *testing.T) {
			t.Parallel()

			gotQuery, err := tests[i].spec.ToFind(tests[i].paging)
			if err != nil {
				t.Errorf("ToFind() error = %v", err)

				return
			}

			assert.Equal(t, tests[i].wantQuery, gotQuery)
		})
	}
}

func Test_newsWithFullTextByKeywords_ToCount(t *testing.T) {
	tests := []struct {
		name      string
		spec      newsWithFullTextByKeywords
		wantQuery string
	}{
		{
			name: "AND operation count",
			spec: newsWithFullTextByKeywords{
				keywords:         []string{"golang", "programming"},
				previewLimit:     512,
				fullTextSearchOp: FullTextSearchOpAnd,
				dialect:          goqu.Dialect("mysql"),
			},
			wantQuery: "SELECT COUNT(*) FROM `news` WHERE MATCH('@content golang programming')",
		},
		{
			name: "Proximity operation count",
			spec: newsWithFullTextByKeywords{
				keywords:         []string{"machine", "learning"},
				previewLimit:     256,
				fullTextSearchOp: FullTextSearchOpProximity,
				dialect:          goqu.Dialect("mysql"),
			},
			wantQuery: "SELECT COUNT(*) FROM `news` WHERE MATCH('@content \\\"machine learning\\\"~10')",
		},
		{
			name: "Quorum operation count",
			spec: newsWithFullTextByKeywords{
				keywords:         []string{"artificial", "intelligence"},
				previewLimit:     1024,
				fullTextSearchOp: FullTextSearchOpQuorum,
				dialect:          goqu.Dialect("mysql"),
			},
			wantQuery: "SELECT COUNT(*) FROM `news` WHERE MATCH('@content \\\"artificial intelligence\\\"/0.3')",
		},
		{
			name: "Single keyword count",
			spec: newsWithFullTextByKeywords{
				keywords:         []string{"blockchain"},
				previewLimit:     128,
				fullTextSearchOp: FullTextSearchOpAnd,
				dialect:          goqu.Dialect("mysql"),
			},
			wantQuery: "SELECT COUNT(*) FROM `news` WHERE MATCH('@content blockchain')",
		},
	}

	for i := range tests {
		t.Run(tests[i].name, func(t *testing.T) {
			t.Parallel()

			gotQuery, err := tests[i].spec.ToCount()
			if err != nil {
				t.Errorf("ToCount() error = %v", err)

				return
			}

			assert.Equal(t, tests[i].wantQuery, gotQuery)
		})
	}
}