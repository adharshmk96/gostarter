package utils

type PaginationParams struct {
	Page int `json:"page"`
	Size int `json:"size"`
}

func (p *PaginationParams) GetOffset() int {
	return (p.Page - 1) * p.Size
}
