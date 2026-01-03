package specifications

import (
	"fmt"
	"strings"

	"github.com/doug-martin/goqu/v9"
	"github.com/financial_advisor/app/domain/repository/specifications"
)

type FullTextSearchOp int

const (
	FullTextSearchOpAnd FullTextSearchOp = iota
	FullTextSearchOpProximity
	FullTextSearchOpQuorum
)

type newsWithFullTextByKeywords struct {
	base
	keywords         []string
	previewLimit     uint64
	fullTextSearchOp FullTextSearchOp
	dialect          goqu.DialectWrapper
	isCount          bool
}

func NewNewsWithFullTextByKeywords(
	keywords []string,
	previewLimit uint64,
	fullTextSearchOp FullTextSearchOp,
) specifications.I {
	spec := newsWithFullTextByKeywords{
		keywords:         keywords,
		previewLimit:     previewLimit,
		fullTextSearchOp: fullTextSearchOp,
		dialect:          goqu.Dialect("mysql"),
	}

	spec.base.buildQuery = spec.buildQuery

	return spec
}

func (q newsWithFullTextByKeywords) buildQuery() *goqu.SelectDataset {
	var (
		inputs = strings.Join(q.keywords, " ")
		query  = q.dialect.From("news")
	)

	switch q.fullTextSearchOp {
	case FullTextSearchOpProximity:
		inputs = `"` + inputs + `"` + "~10"
	case FullTextSearchOpQuorum:
		// we set default proximity to 10 for quorum search
		inputs = `"` + inputs + `"` + "/0.3"
	default:
	}

	query = query.Where(goqu.L("MATCH(?)", "@content "+inputs))

	if q.previewLimit == 0 {
		q.previewLimit = 256
	}

	if q.isCount {
		return query
	}

	query = query.Select(goqu.L(fmt.Sprintf(
		"HIGHLIGHT({limit=%d,before_match='<mark>',after_match='</mark>'}, content) AS content, news_id, id",
		q.previewLimit,
	)))

	return query
}

func (q newsWithFullTextByKeywords) ToFind(paging specifications.PagingI) (string, error) {
	q.isCount = false
	q.base.buildQuery = q.buildQuery

	return q.base.ToFind(paging)
}

func (q newsWithFullTextByKeywords) ToCount() (string, error) {
	q.isCount = true
	q.base.buildQuery = q.buildQuery

	return q.base.ToCount()
}
