package common

type Paging struct {
	Page       int    `json:"page" form:"page"`
	Limit      int    `json:"limit" form:"limit"`
	Total      int64  `json:"total"`
	FakeCursor string `json:"cursor" form:"cursor"`
	NextCursor string `json:"next_cursor"`
}

func (p *Paging) Process() error {
	if p.Page <= 0 {
		p.Page = 1
	}

	if p.Limit <= 0 || p.Limit > 100 {
		p.Limit = 20
	}

	return nil
}
