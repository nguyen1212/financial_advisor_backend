package specifications

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/financial_advisor/app/domain/repository/specifications"
)

type publisherByID struct {
	base
	id      uint64
	dialect goqu.DialectWrapper
}

func NewPublisherByID(id uint64) specifications.I {
	spec := publisherByID{
		id:      id,
		dialect: goqu.Dialect("mysql"),
	}

	spec.base.buildQuery = spec.buildQuery

	return spec
}

func (q publisherByID) buildQuery() *goqu.SelectDataset {
	query := q.dialect.From("publishers")

	query = query.Where(goqu.C("id").Eq(q.id))

	return query
}
