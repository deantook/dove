package pagination

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestParsePageRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name             string
		query            string
		expectedPage     int
		expectedSize     int
		expectedSortBy   string
		expectedOrder    string
		expectedKeyword  string
		expectedSearchBy string
	}{
		{"default values", "", 1, 10, "", "desc", "", ""},
		{"custom values", "?page=2&page_size=20", 2, 20, "", "desc", "", ""},
		{"with sort", "?page=1&page_size=10&sort_by=created_at&sort_order=asc", 1, 10, "created_at", "asc", "", ""},
		{"with search", "?page=1&page_size=10&keyword=john&search_by=username", 1, 10, "", "desc", "john", "username"},
		{"with sort and search", "?page=1&page_size=10&sort_by=created_at&sort_order=desc&keyword=test&search_by=email", 1, 10, "created_at", "desc", "test", "email"},
		{"invalid sort order", "?page=1&page_size=10&sort_by=created_at&sort_order=invalid", 1, 10, "created_at", "desc", "", ""},
		{"page too small", "?page=0&page_size=10", 1, 10, "", "desc", "", ""},
		{"page size too small", "?page=1&page_size=0", 1, 10, "", "desc", "", ""},
		{"page size too large", "?page=1&page_size=200", 1, 100, "", "desc", "", ""},
		{"invalid values", "?page=abc&page_size=def", 1, 10, "", "desc", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			req, _ := http.NewRequest("GET", "/test"+tt.query, nil)
			c.Request = req

			result := ParsePageRequest(c)

			if result.Page != tt.expectedPage {
				t.Errorf("ParsePageRequest() page = %v, want %v", result.Page, tt.expectedPage)
			}
			if result.PageSize != tt.expectedSize {
				t.Errorf("ParsePageRequest() pageSize = %v, want %v", result.PageSize, tt.expectedSize)
			}
			if result.SortBy != tt.expectedSortBy {
				t.Errorf("ParsePageRequest() sortBy = %v, want %v", result.SortBy, tt.expectedSortBy)
			}
			if result.SortOrder != tt.expectedOrder {
				t.Errorf("ParsePageRequest() sortOrder = %v, want %v", result.SortOrder, tt.expectedOrder)
			}
			if result.Keyword != tt.expectedKeyword {
				t.Errorf("ParsePageRequest() keyword = %v, want %v", result.Keyword, tt.expectedKeyword)
			}
			if result.SearchBy != tt.expectedSearchBy {
				t.Errorf("ParsePageRequest() searchBy = %v, want %v", result.SearchBy, tt.expectedSearchBy)
			}
		})
	}
}

func TestNewPageResponse(t *testing.T) {
	data := []string{"item1", "item2", "item3"}

	tests := []struct {
		name      string
		total     int64
		page      int
		pageSize  int
		wantPages int
		wantNext  bool
		wantPrev  bool
	}{
		{"first page", 25, 1, 10, 3, true, false},
		{"middle page", 25, 2, 10, 3, true, true},
		{"last page", 25, 3, 10, 3, false, true},
		{"single page", 5, 1, 10, 1, false, false},
		{"exact pages", 20, 2, 10, 2, false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewPageResponse(data, tt.total, tt.page, tt.pageSize)

			if result.Total != tt.total {
				t.Errorf("NewPageResponse() total = %v, want %v", result.Total, tt.total)
			}
			if result.Page != tt.page {
				t.Errorf("NewPageResponse() page = %v, want %v", result.Page, tt.page)
			}
			if result.PageSize != tt.pageSize {
				t.Errorf("NewPageResponse() pageSize = %v, want %v", result.PageSize, tt.pageSize)
			}
			if result.TotalPages != tt.wantPages {
				t.Errorf("NewPageResponse() totalPages = %v, want %v", result.TotalPages, tt.wantPages)
			}
			if result.HasNext != tt.wantNext {
				t.Errorf("NewPageResponse() hasNext = %v, want %v", result.HasNext, tt.wantNext)
			}
			if result.HasPrev != tt.wantPrev {
				t.Errorf("NewPageResponse() hasPrev = %v, want %v", result.HasPrev, tt.wantPrev)
			}
		})
	}
}

func TestPageRequest_GetOffset(t *testing.T) {
	tests := []struct {
		name     string
		page     int
		pageSize int
		want     int
	}{
		{"first page", 1, 10, 0},
		{"second page", 2, 10, 10},
		{"third page", 3, 10, 20},
		{"custom page size", 2, 20, 20},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &PageRequest{
				Page:     tt.page,
				PageSize: tt.pageSize,
			}
			if got := req.GetOffset(); got != tt.want {
				t.Errorf("PageRequest.GetOffset() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPageRequest_GetLimit(t *testing.T) {
	req := &PageRequest{Page: 1, PageSize: 25}
	if got := req.GetLimit(); got != 25 {
		t.Errorf("PageRequest.GetLimit() = %v, want %v", got, 25)
	}
}

func TestPageRequest_SortMethods(t *testing.T) {
	req := &PageRequest{
		Page:      1,
		PageSize:  10,
		SortBy:    "created_at",
		SortOrder: "desc",
	}

	if got := req.GetSortBy(); got != "created_at" {
		t.Errorf("PageRequest.GetSortBy() = %v, want %v", got, "created_at")
	}

	if got := req.GetSortOrder(); got != "desc" {
		t.Errorf("PageRequest.GetSortOrder() = %v, want %v", got, "desc")
	}

	if !req.HasSort() {
		t.Error("PageRequest.HasSort() should return true when sort_by is set")
	}

	req.SortBy = ""
	if req.HasSort() {
		t.Error("PageRequest.HasSort() should return false when sort_by is empty")
	}
}

func TestPageRequest_ValidateSortField(t *testing.T) {
	allowedFields := []string{"id", "created_at", "updated_at", "name"}

	tests := []struct {
		name          string
		sortBy        string
		expectedValid bool
	}{
		{"valid field", "created_at", true},
		{"valid field 2", "name", true},
		{"invalid field", "invalid_field", false},
		{"empty field", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &PageRequest{SortBy: tt.sortBy}
			if got := req.ValidateSortField(allowedFields); got != tt.expectedValid {
				t.Errorf("PageRequest.ValidateSortField() = %v, want %v", got, tt.expectedValid)
			}
		})
	}
}

func TestPageRequest_SearchMethods(t *testing.T) {
	req := &PageRequest{
		Page:     1,
		PageSize: 10,
		Keyword:  "john",
		SearchBy: "username",
	}

	if got := req.GetKeyword(); got != "john" {
		t.Errorf("PageRequest.GetKeyword() = %v, want %v", got, "john")
	}

	if got := req.GetSearchBy(); got != "username" {
		t.Errorf("PageRequest.GetSearchBy() = %v, want %v", got, "username")
	}

	if !req.HasSearch() {
		t.Error("PageRequest.HasSearch() should return true when keyword is set")
	}

	req.Keyword = ""
	if req.HasSearch() {
		t.Error("PageRequest.HasSearch() should return false when keyword is empty")
	}
}

func TestPageRequest_ValidateSearchField(t *testing.T) {
	allowedFields := []string{"username", "email", "nickname", "name"}

	tests := []struct {
		name          string
		searchBy      string
		expectedValid bool
	}{
		{"valid field", "username", true},
		{"valid field 2", "email", true},
		{"invalid field", "invalid_field", false},
		{"empty field", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &PageRequest{SearchBy: tt.searchBy}
			if got := req.ValidateSearchField(allowedFields); got != tt.expectedValid {
				t.Errorf("PageRequest.ValidateSearchField() = %v, want %v", got, tt.expectedValid)
			}
		})
	}
}
