package scopes

import (
	"github.com/jinzhu/gorm"
)

type Pagination struct {
	CurrentPage int `json:"current_page"`
	PageSize    int `json:"page_size"`
	Total       int `json:"total"`
	IsLastPage  int `json:"is_last_page"`
}

func (p *Pagination) init(currentPage int, pageSize int, total int) {
	if currentPage == 0 {
		pageSize = 0
	} else if pageSize == 0 {
		pageSize = 10
	}

	p.CurrentPage = currentPage
	p.PageSize = pageSize
	p.Total = total

	isLastPage := 0
	if p.CurrentPage != 0 && p.PageSize != 0 {
		if p.CurrentPage*p.PageSize >= p.Total {
			isLastPage = 1
		}
	} else {
		isLastPage = 1
	}

	p.IsLastPage = isLastPage
}

func PaginationScope(p *Pagination, currentPage int, pageSize ...int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		var (
			total     int = 0
			offset    int = 0
			_pageSize int = 10
		)

		if len(pageSize) > 0 {
			_pageSize = pageSize[0]
		}

		db.Count(&total)

		p.init(currentPage, _pageSize, total)

		if p.CurrentPage != 0 && p.PageSize != 0 {
			offset = (p.CurrentPage - 1) * p.PageSize
			return db.Offset(offset).Limit(p.PageSize)
		}
		return db
	}
}
