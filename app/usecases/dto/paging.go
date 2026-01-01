package dto

type PagingRequest struct {
	Page int
	Size int
}

type PagingResponse struct {
	Total int
}
