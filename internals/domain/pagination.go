package domain

type Pagination struct {
	Page  int `json:"page"`
	Size  int `json:"size"`
	Total int `json:"total"`
}

func (p *Pagination) GetOffset() int {
	return (p.Page - 1) * p.Size
}

func (p *Pagination) SetTotal(total int) {
	p.Total = total
}
