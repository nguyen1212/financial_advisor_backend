// Package specifications holds implementation to build a raw query by goqu lib
package specifications

import (
	"time"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
	"github.com/financial_advisor/app/domain/entity"
	"github.com/financial_advisor/app/domain/repository/specifications"
)

type newsByDate struct {
	base
	startDate time.Time
	endDate   time.Time
	status    *entity.NewsStatus
	dialect   goqu.DialectWrapper
}

func NewsByDate(
	startDate,
	endDate time.Time,
	status *entity.NewsStatus,
) specifications.I {
	spec := &newsByDate{
		dialect:   goqu.Dialect("mysql"),
		startDate: startDate,
		endDate:   endDate,
		status:    status,
	}

	spec.base.buildQuery = spec.buildQuery

	return spec
}

func (q newsByDate) buildQuery() *goqu.SelectDataset {
	query := q.dialect.From("news")

	if !q.startDate.IsZero() {
		query = query.Where(
			goqu.C("created_at").Gte(q.startDate),
		)
	}

	if !q.endDate.IsZero() {
		query = query.Where(
			goqu.C("created_at").Lte(q.endDate),
		)
	}

	if q.status != nil {
		query = query.Where(
			goqu.C("status").Eq(*q.status),
		)
	}

	return query
}
