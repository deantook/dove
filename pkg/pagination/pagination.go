package pagination

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// PageRequest 分页请求参数
type PageRequest struct {
	Page      int    `json:"page" form:"page" binding:"min=1" example:"1"`                    // 页码，从1开始
	PageSize  int    `json:"page_size" form:"page_size" binding:"min=1,max=100" example:"10"` // 每页大小，最大100
	SortBy    string `json:"sort_by" form:"sort_by" example:"created_at"`                     // 排序字段
	SortOrder string `json:"sort_order" form:"sort_order" example:"desc"`                     // 排序方向：asc, desc
	Keyword   string `json:"keyword" form:"keyword" example:"john"`                           // 搜索关键词
	SearchBy  string `json:"search_by" form:"search_by" example:"username"`                   // 搜索字段
}

// PageResponse 分页响应
type PageResponse struct {
	Data       interface{} `json:"data"`        // 数据列表
	Total      int64       `json:"total"`       // 总记录数
	Page       int         `json:"page"`        // 当前页码
	PageSize   int         `json:"page_size"`   // 每页大小
	TotalPages int         `json:"total_pages"` // 总页数
	HasNext    bool        `json:"has_next"`    // 是否有下一页
	HasPrev    bool        `json:"has_prev"`    // 是否有上一页
}

// ParsePageRequest 从 gin.Context 解析分页请求
func ParsePageRequest(c *gin.Context) *PageRequest {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	sortBy := c.DefaultQuery("sort_by", "")
	sortOrder := c.DefaultQuery("sort_order", "desc")
	keyword := c.DefaultQuery("keyword", "")
	searchBy := c.DefaultQuery("search_by", "")

	// 设置默认值和限制
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	// 验证排序方向
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}

	return &PageRequest{
		Page:      page,
		PageSize:  pageSize,
		SortBy:    sortBy,
		SortOrder: sortOrder,
		Keyword:   keyword,
		SearchBy:  searchBy,
	}
}

// NewPageResponse 创建分页响应
func NewPageResponse(data interface{}, total int64, page, pageSize int) *PageResponse {
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	return &PageResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}
}

// GetOffset 获取数据库查询的偏移量
func (p *PageRequest) GetOffset() int {
	return (p.Page - 1) * p.PageSize
}

// GetLimit 获取数据库查询的限制数量
func (p *PageRequest) GetLimit() int {
	return p.PageSize
}

// GetSortBy 获取排序字段
func (p *PageRequest) GetSortBy() string {
	return p.SortBy
}

// GetSortOrder 获取排序方向
func (p *PageRequest) GetSortOrder() string {
	return p.SortOrder
}

// HasSort 检查是否有排序参数
func (p *PageRequest) HasSort() bool {
	return p.SortBy != ""
}

// ValidateSortField 验证排序字段是否有效
func (p *PageRequest) ValidateSortField(allowedFields []string) bool {
	if p.SortBy == "" {
		return true
	}
	for _, field := range allowedFields {
		if field == p.SortBy {
			return true
		}
	}
	return false
}

// GetKeyword 获取搜索关键词
func (p *PageRequest) GetKeyword() string {
	return p.Keyword
}

// GetSearchBy 获取搜索字段
func (p *PageRequest) GetSearchBy() string {
	return p.SearchBy
}

// HasSearch 检查是否有搜索参数
func (p *PageRequest) HasSearch() bool {
	return p.Keyword != ""
}

// ValidateSearchField 验证搜索字段是否有效
func (p *PageRequest) ValidateSearchField(allowedFields []string) bool {
	if p.SearchBy == "" {
		return true
	}
	for _, field := range allowedFields {
		if field == p.SearchBy {
			return true
		}
	}
	return false
}
