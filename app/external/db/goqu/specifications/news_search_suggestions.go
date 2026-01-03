package specifications

import (
	"fmt"
	"strings"

	"github.com/doug-martin/goqu/v9"
	"github.com/financial_advisor/app/domain/repository/specifications"
)

type Fuzziness int

const (
	DisabledFuzziness Fuzziness = iota
	StrongFuzziness
	MediumFuzziness
)

type newsSearchSuggestions struct {
	base
	keywords []string
	fuziness Fuzziness
	dialect  goqu.DialectWrapper
}

func NewNewsSearchSuggestions(
	keywords []string,
	fuzziness Fuzziness,
) specifications.I {
	spec := newsSearchSuggestions{
		keywords: keywords,
		fuziness: fuzziness,
		dialect:  goqu.Dialect("mysql"),
	}

	spec.base.buildQuery = spec.buildQuery

	return spec
}

func (q newsSearchSuggestions) buildQuery() *goqu.SelectDataset {
	var (
		inputs = strings.Join(q.keywords, " ")
		query  = q.dialect.From("news").Where(goqu.L(
			fmt.Sprintf("CALL AUTOCOMPLETE(?, 'news', %d AS fuzziness, 1 AS preserve)", q.fuziness),
			inputs,
		))
	)

	return query
}

func (q newsSearchSuggestions) ToFind(paging specifications.PagingI) (string, error) {
	q.base.buildQuery = q.buildQuery

	sql, err := q.base.ToFind(paging)
	if err != nil {
		return "", err
	}

	sql = strings.TrimPrefix(sql, "SELECT * FROM `news` WHERE ")

	return sql, nil
}
