package specifications

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/financial_advisor/app/domain/repository/specifications"
)

type publishersByNone struct {
	base
	dialect goqu.DialectWrapper
}

func NewPublishersByNone() specifications.I {
	spec := publishersByNone{
		dialect: goqu.Dialect("mysql"),
	}

	spec.base.buildQuery = spec.buildQuery

	return spec
}

func (q publishersByNone) buildQuery() *goqu.SelectDataset {
	query := q.dialect.From("publishers")

	return query
}
