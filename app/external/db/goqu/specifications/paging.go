package specifications

import "github.com/financial_advisor/app/domain/repository/specifications"

type paging struct {
	limit  int
	offset int
}

func ToPaging(size, page int) specifications.PagingI {
	return paging{
		limit:  size,
		offset: (page - 1) * size,
	}
}

func (p paging) Limit() int {
	return p.limit
}

func (p paging) Offset() int {
	return p.offset
}
