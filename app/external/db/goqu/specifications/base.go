package specifications

import (
	"errors"

	"github.com/doug-martin/goqu/v9"
	"github.com/financial_advisor/app/domain/repository/specifications"
)

type base struct {
	buildQuery func() *goqu.SelectDataset
}

func (q base) ToSQL() (string, []any, error) {
	if q.buildQuery == nil {
		return "", nil, errors.New("spec not implemented")
	}

	query := q.buildQuery()
	if query == nil {
		return "", nil, errors.New("invalid spec")
	}

	return query.ToSQL()
}

func (q base) ToCount() (string, error) {
	if q.buildQuery == nil {
		return "", errors.New("spec not implemented")
	}

	query := q.buildQuery()
	if query == nil {
		return "", errors.New("invalid spec")
	}

	sql, _, err := query.Select(goqu.COUNT("*")).ToSQL()

	return sql, err
}

func (q base) ToGet() (string, error) {
	if q.buildQuery == nil {
		return "", errors.New("spec not implemented")
	}

	query := q.buildQuery()
	if query == nil {
		return "", errors.New("invalid spec")
	}

	sql, _, err := query.Limit(1).ToSQL()

	return sql, err
}

func (q base) ToFind(paging specifications.PagingI) (string, error) {
	if q.buildQuery == nil {
		return "", errors.New("spec not implemented")
	}

	query := q.buildQuery()
	if query == nil {
		return "", errors.New("invalid spec")
	}

	if paging != nil {
		query = query.Limit(uint(paging.Limit())).Offset(uint(paging.Offset()))
	}

	sql, _, err := query.ToSQL()
	if err != nil {
		return "", err
	}

	return sql, nil
}
