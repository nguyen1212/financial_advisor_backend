package specifications

type Paging struct {
	limit  int
	offset int
}

func newPaging(limit, offset int) Paging {
	return Paging{
		limit:  limit,
		offset: offset,
	}
}

func (p Paging) Limit() int {
	return p.limit
}

func (p Paging) Offset() int {
	return p.offset
}

func ToPaging(
	size,
	page int,
) Paging {
	return newPaging(
		size,
		size*(page-1),
	)
}
