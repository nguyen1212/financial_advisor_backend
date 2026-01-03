package specifications

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/financial_advisor/app/domain/repository/specifications"
)

type publisherByDomain struct {
	base
	domain  string
	dialect goqu.DialectWrapper
}

func NewPublisherByDomain(domain string) specifications.I {
	spec := publisherByDomain{
		domain:  domain,
		dialect: goqu.Dialect("mysql"),
	}

	spec.base.buildQuery = spec.buildQuery

	return spec
}

func (q publisherByDomain) buildQuery() *goqu.SelectDataset {
	query := q.dialect.From("publishers")

	query = query.Where(goqu.C("domain").Eq(q.domain))

	return query
}

// func (q PublisherByDomain) ToSQL() (string, []any, error) {
// 	return q.buildQuery().ToSQL()
// }
//
// func (q PublisherByDomain) ToCount() (string, error) {
// 	query := q.buildQuery()
//
// 	sql, _, err := query.Select(goqu.COUNT("*")).ToSQL()
//
// 	return sql, err
// }
//
// func (q PublisherByDomain) ToGet() (string, error) {
// 	query := q.buildQuery()
//
// 	sql, _, err := query.Select().Limit(1).ToSQL()
//
// 	return sql, err
// }
