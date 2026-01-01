package specifications

import (
	"testing"

	"github.com/doug-martin/goqu/v9"
	"github.com/financial_advisor/app/domain/repository/specifications"
	"github.com/stretchr/testify/assert"
)

func Test_NewNewsWithFullTextByFileID(t *testing.T) {
	var (
		fileID           = uint64(789)
		previewLimit     = int64(512)
		highlightKeywords = []string{"test", "keyword"}
		got              = NewNewsWithFullTextByFileID(fileID, previewLimit, highlightKeywords)
		want             = newsWithFullTextByFileID{
			fileID:           fileID,
			previewLimit:     previewLimit,
			higlightKeywords: highlightKeywords,
			dialect:          goqu.Dialect("mysql"),
		}
	)

	actualSpec := got.(newsWithFullTextByFileID)
	assert.Equal(t, want.dialect, actualSpec.dialect)
	assert.Equal(t, want.fileID, actualSpec.fileID)
	assert.Equal(t, want.previewLimit, actualSpec.previewLimit)
	assert.Equal(t, want.higlightKeywords, actualSpec.higlightKeywords)
}

func Test_newsWithFullTextByFileID_ToGet(t *testing.T) {
	tests := []struct {
		name      string
		spec      newsWithFullTextByFileID
		wantQuery string
	}{
		{
			name: "with highlight keywords",
			spec: newsWithFullTextByFileID{
				fileID:           789,
				previewLimit:     512,
				higlightKeywords: []string{"test", "keyword"},
				dialect:          goqu.Dialect("mysql"),
			},
			wantQuery: "SELECT HIGHLIGHT({limit=512,before_match='<mark>',after_match='</mark>'}, content) AS content, news_id, id FROM `news` WHERE MATCH('@content test keyword') AND id = 789 LIMIT 1",
		},
		{
			name: "with empty keywords",
			spec: newsWithFullTextByFileID{
				fileID:           789,
				previewLimit:     256,
				higlightKeywords: []string{},
				dialect:          goqu.Dialect("mysql"),
			},
			wantQuery: "SELECT HIGHLIGHT({limit=256,before_match='<mark>',after_match='</mark>'}, content) AS content, news_id, id FROM `news` WHERE (`id` = 789) LIMIT 1",
		},
		{
			name: "with single keyword",
			spec: newsWithFullTextByFileID{
				fileID:           123,
				previewLimit:     128,
				higlightKeywords: []string{"search"},
				dialect:          goqu.Dialect("mysql"),
			},
			wantQuery: "SELECT HIGHLIGHT({limit=128,before_match='<mark>',after_match='</mark>'}, content) AS content, news_id, id FROM `news` WHERE MATCH('@content search') AND id = 123 LIMIT 1",
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

func Test_newsWithFullTextByFileID_ToCount(t *testing.T) {
	tests := []struct {
		name      string
		spec      newsWithFullTextByFileID
		wantQuery string
	}{
		{
			name: "with highlight keywords",
			spec: newsWithFullTextByFileID{
				fileID:           789,
				previewLimit:     512,
				higlightKeywords: []string{"test", "keyword"},
				dialect:          goqu.Dialect("mysql"),
			},
			wantQuery: "SELECT COUNT(*) FROM `news` WHERE MATCH('@content test keyword') AND id = 789",
		},
		{
			name: "with empty keywords",
			spec: newsWithFullTextByFileID{
				fileID:           789,
				previewLimit:     256,
				higlightKeywords: []string{},
				dialect:          goqu.Dialect("mysql"),
			},
			wantQuery: "SELECT COUNT(*) FROM `news` WHERE (`id` = 789)",
		},
		{
			name: "with single keyword",
			spec: newsWithFullTextByFileID{
				fileID:           123,
				previewLimit:     128,
				higlightKeywords: []string{"search"},
				dialect:          goqu.Dialect("mysql"),
			},
			wantQuery: "SELECT COUNT(*) FROM `news` WHERE MATCH('@content search') AND id = 123",
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

func Test_newsWithFullTextByFileID_ToFind(t *testing.T) {
	tests := []struct {
		name      string
		spec      newsWithFullTextByFileID
		paging    specifications.PagingI
		wantQuery string
	}{
		{
			name: "with keywords and paging",
			spec: newsWithFullTextByFileID{
				fileID:           789,
				previewLimit:     512,
				higlightKeywords: []string{"test", "keyword"},
				dialect:          goqu.Dialect("mysql"),
			},
			paging:    paging{limit: 10, offset: 20},
			wantQuery: "SELECT HIGHLIGHT({limit=512,before_match='<mark>',after_match='</mark>'}, content) AS content, news_id, id FROM `news` WHERE MATCH('@content test keyword') AND id = 789 LIMIT 10 OFFSET 20",
		},
		{
			name: "with keywords without paging",
			spec: newsWithFullTextByFileID{
				fileID:           789,
				previewLimit:     256,
				higlightKeywords: []string{"search"},
				dialect:          goqu.Dialect("mysql"),
			},
			paging:    nil,
			wantQuery: "SELECT HIGHLIGHT({limit=256,before_match='<mark>',after_match='</mark>'}, content) AS content, news_id, id FROM `news` WHERE MATCH('@content search') AND id = 789",
		},
		{
			name: "without keywords with paging",
			spec: newsWithFullTextByFileID{
				fileID:           123,
				previewLimit:     128,
				higlightKeywords: []string{},
				dialect:          goqu.Dialect("mysql"),
			},
			paging:    paging{limit: 5, offset: 10},
			wantQuery: "SELECT HIGHLIGHT({limit=128,before_match='<mark>',after_match='</mark>'}, content) AS content, news_id, id FROM `news` WHERE (`id` = 123) LIMIT 5 OFFSET 10",
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