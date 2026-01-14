package specifications

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/financial_advisor/app/domain/repository/specifications"
)

type jobByUUID struct {
	base
	uuid    string
	dialect goqu.DialectWrapper
}

func JobByUUID(uuid string) specifications.I {
	q := jobByUUID{
		uuid:    uuid,
		dialect: goqu.Dialect("mysql"),
	}

	q.base.buildQuery = q.buildQuery

	return q
}

func (q jobByUUID) buildQuery() *goqu.SelectDataset {
	query := q.dialect.From("jobs").Where(goqu.Ex{
		"uuid": q.uuid,
	})

	return query
}
