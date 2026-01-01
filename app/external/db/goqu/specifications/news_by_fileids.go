package specifications

import (
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/financial_advisor/app/domain/repository/specifications"
)

type newsWithFullTextByFileIDs struct {
	base
	fileIDs      []uint64
	previewLimit uint64
	dialect      goqu.DialectWrapper
}

func NewNewsWithFullTextByFileIDs(
	fileIDs []uint64,
	previewLimit uint64,
) specifications.I {
	spec := newsWithFullTextByFileIDs{
		fileIDs:      fileIDs,
		previewLimit: previewLimit,
		dialect:      goqu.Dialect("mysql"),
	}

	spec.base.buildQuery = spec.buildQuery

	return spec
}

func (q newsWithFullTextByFileIDs) buildQuery() *goqu.SelectDataset {
	query := q.dialect.From("news").
		Where(goqu.C("id").In(q.fileIDs))

	if q.previewLimit == 0 {
		q.previewLimit = 256
	}

	return query.Select(goqu.L(fmt.Sprintf("HIGHLIGHT({limit=%d}, content) AS content, id, news_id", q.previewLimit)))
}
