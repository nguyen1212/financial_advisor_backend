package specifications

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/financial_advisor/app/domain/repository/specifications"
)

type newsByHashedURL struct {
	base
	hashedURL []byte
	dialect   goqu.DialectWrapper
}

func NewNewsByHashedURL(hashedURL []byte) specifications.I {
	spec := newsByHashedURL{
		hashedURL: hashedURL,
		dialect:   goqu.Dialect("mysql"),
	}

	spec.base.buildQuery = spec.buildQuery

	return spec
}

func (q newsByHashedURL) buildQuery() *goqu.SelectDataset {
	query := q.dialect.From("news")

	if q.hashedURL != nil {
		query = query.Where(
			goqu.C("hashed_url").Eq(q.hashedURL),
		)
	}

	return query
}
