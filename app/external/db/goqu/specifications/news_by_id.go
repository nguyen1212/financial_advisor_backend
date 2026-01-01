package specifications

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/financial_advisor/app/domain/repository/specifications"
)

type newsByID struct {
	base
	id      uint64
	dialect goqu.DialectWrapper
}

func NewNewsByID(id uint64) specifications.I {
	spec := newsByID{
		id:      id,
		dialect: goqu.Dialect("mysql"),
	}

	spec.base.buildQuery = spec.buildQuery

	return spec
}

func (q newsByID) buildQuery() *goqu.SelectDataset {
	query := q.dialect.From("news")

	query = query.Where(goqu.C("id").Eq(q.id))

	return query
}
