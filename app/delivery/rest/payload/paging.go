package payload

type Paging struct {
	Size int `form:"size"`
	Page int `form:"page"`
}
