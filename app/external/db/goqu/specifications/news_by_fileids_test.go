package specifications

import (
	"testing"

	"github.com/doug-martin/goqu/v9"
	"github.com/financial_advisor/app/domain/repository/specifications"
	"github.com/stretchr/testify/assert"
)

func Test_NewNewsWithFullTextByFileIDs(t *testing.T) {
	var (
		fileIDs      = []uint64{1, 2, 3}
		previewLimit = uint64(1024)
		got          = NewNewsWithFullTextByFileIDs(fileIDs, previewLimit)
		want         = newsWithFullTextByFileIDs{
			fileIDs:      fileIDs,
			previewLimit: previewLimit,
			dialect:      goqu.Dialect("mysql"),
		}
	)

	actualSpec := got.(newsWithFullTextByFileIDs)
	assert.Equal(t, want.dialect, actualSpec.dialect)
	assert.Equal(t, want.fileIDs, actualSpec.fileIDs)
	assert.Equal(t, want.previewLimit, actualSpec.previewLimit)
}

func Test_newsWithFullTextByFileIDs_ToGet(t *testing.T) {
	tests := []struct {
		name      string
		spec      newsWithFullTextByFileIDs
		wantQuery string
	}{
		{
			name: "with multiple file IDs and custom preview limit",
			spec: newsWithFullTextByFileIDs{
				fileIDs:      []uint64{1, 2, 3},
				previewLimit: 512,
				dialect:      goqu.Dialect("mysql"),
			},
			wantQuery: "SELECT HIGHLIGHT({limit=512}, content) AS content, id, news_id FROM `news` WHERE (`id` IN (1, 2, 3)) LIMIT 1",
		},
		{
			name: "with single file ID and zero preview limit (should default to 256)",
			spec: newsWithFullTextByFileIDs{
				fileIDs:      []uint64{123},
				previewLimit: 0,
				dialect:      goqu.Dialect("mysql"),
			},
			wantQuery: "SELECT HIGHLIGHT({limit=256}, content) AS content, id, news_id FROM `news` WHERE (`id` IN (123)) LIMIT 1",
		},
		{
			name: "with empty file IDs",
			spec: newsWithFullTextByFileIDs{
				fileIDs:      []uint64{},
				previewLimit: 256,
				dialect:      goqu.Dialect("mysql"),
			},
			wantQuery: "SELECT HIGHLIGHT({limit=256}, content) AS content, id, news_id FROM `news` WHERE (`id` IN ()) LIMIT 1",
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

func Test_newsWithFullTextByFileIDs_ToCount(t *testing.T) {
	tests := []struct {
		name      string
		spec      newsWithFullTextByFileIDs
		wantQuery string
	}{
		{
			name: "with multiple file IDs",
			spec: newsWithFullTextByFileIDs{
				fileIDs:      []uint64{1, 2, 3, 4, 5},
				previewLimit: 512,
				dialect:      goqu.Dialect("mysql"),
			},
			wantQuery: "SELECT COUNT(*) FROM `news` WHERE (`id` IN (1, 2, 3, 4, 5))",
		},
		{
			name: "with single file ID",
			spec: newsWithFullTextByFileIDs{
				fileIDs:      []uint64{999},
				previewLimit: 256,
				dialect:      goqu.Dialect("mysql"),
			},
			wantQuery: "SELECT COUNT(*) FROM `news` WHERE (`id` IN (999))",
		},
		{
			name: "with empty file IDs",
			spec: newsWithFullTextByFileIDs{
				fileIDs:      []uint64{},
				previewLimit: 256,
				dialect:      goqu.Dialect("mysql"),
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

func Test_newsWithFullTextByFileIDs_ToFind(t *testing.T) {
	tests := []struct {
		name      string
		spec      newsWithFullTextByFileIDs
		paging    specifications.PagingI
		wantQuery string
	}{
		{
			name: "with file IDs and paging",
			spec: newsWithFullTextByFileIDs{
				fileIDs:      []uint64{10, 20, 30},
				previewLimit: 512,
				dialect:      goqu.Dialect("mysql"),
			},
			paging:    paging{limit: 15, offset: 30},
			wantQuery: "SELECT HIGHLIGHT({limit=512}, content) AS content, id, news_id FROM `news` WHERE (`id` IN (10, 20, 30)) LIMIT 15 OFFSET 30",
		},
		{
			name: "with file IDs without paging",
			spec: newsWithFullTextByFileIDs{
				fileIDs:      []uint64{100, 200},
				previewLimit: 256,
				dialect:      goqu.Dialect("mysql"),
			},
			paging:    nil,
			wantQuery: "SELECT HIGHLIGHT({limit=256}, content) AS content, id, news_id FROM `news` WHERE (`id` IN (100, 200))",
		},
		{
			name: "with zero preview limit and paging (should default to 256)",
			spec: newsWithFullTextByFileIDs{
				fileIDs:      []uint64{1, 2},
				previewLimit: 0,
				dialect:      goqu.Dialect("mysql"),
			},
			paging:    paging{limit: 5, offset: 0},
			wantQuery: "SELECT HIGHLIGHT({limit=256}, content) AS content, id, news_id FROM `news` WHERE (`id` IN (1, 2)) LIMIT 5",
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