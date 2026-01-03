package specifications

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/financial_advisor/app/domain/repository/specifications"
)

type newsByIDs struct {
	base
	ids     []uint64
	dialect goqu.DialectWrapper
}

func NewNewsByIDs(ids []uint64) specifications.I {
	spec := newsByIDs{ids: ids, dialect: goqu.Dialect("mysql")}

	spec.base.buildQuery = spec.buildQuery

	return spec
}

func (q newsByIDs) buildQuery() *goqu.SelectDataset {
	query := q.dialect.From("news").Where(goqu.Ex{
		"id": q.ids,
	})

	return query
}
