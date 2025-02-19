package v1

type Page struct {
	Size  int `json:"size,omitempty"`
	Index int `json:"index,omitempty"`
}

type Paginator struct {
	Total int64       `json:"total,omitempty"`
	Data  interface{} `json:"data,omitempty"`
}

func (p Page) GetOffset() int {
	index := p.Index
	if p.Index <= 0 {
		index = 1
	}
	size := p.Size
	if p.Size <= 0 {
		size = 20
	}
	return (index - 1) * size
}
