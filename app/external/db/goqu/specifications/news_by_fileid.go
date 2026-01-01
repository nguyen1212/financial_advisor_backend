package specifications

import (
	"fmt"
	"strings"

	"github.com/doug-martin/goqu/v9"
	"github.com/financial_advisor/app/domain/repository/specifications"
)

type newsWithFullTextByFileID struct {
	base
	fileID           uint64
	previewLimit     int64
	higlightKeywords []string
	dialect          goqu.DialectWrapper
}

func NewNewsWithFullTextByFileID(
	fileID uint64,
	previewLimit int64,
	higlightKeywords []string,
) specifications.I {
	spec := newsWithFullTextByFileID{
		fileID:           fileID,
		previewLimit:     previewLimit,
		higlightKeywords: higlightKeywords,
		dialect:          goqu.Dialect("mysql"),
	}

	spec.base.buildQuery = spec.buildQuery

	return spec
}

func (q newsWithFullTextByFileID) buildQuery() *goqu.SelectDataset {
	var (
		inputs = strings.Join(q.higlightKeywords, " ")
		query  = q.dialect.From("news")
	)

	if len(inputs) > 0 {
		query = query.Where(goqu.L("MATCH(?) AND id = ?", "@content "+inputs, q.fileID))
	} else {
		query = query.Where(goqu.C("id").Eq(q.fileID))
	}

	query = query.Select(goqu.L(fmt.Sprintf(
		"HIGHLIGHT({limit=%d,before_match='<mark>',after_match='</mark>'}, content) AS content, news_id, id",
		q.previewLimit,
	)))

	return query
}
