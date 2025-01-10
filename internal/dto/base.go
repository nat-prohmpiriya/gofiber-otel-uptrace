package dto

type BaseQuery struct {
	Page     int    `query:"page" json:"page" validate:"min=1"`
	PageSize int    `query:"page_size" json:"page_size" validate:"min=1,max=100"`
	SortBy   string `query:"sort_by" json:"sort_by"`
	SortDir  string `query:"sort_dir" json:"sort_dir" validate:"oneof=asc desc"`
}

type PaginationResponse struct {
	CurrentPage int   `json:"current_page"`
	PageSize    int   `json:"page_size"`
	TotalItems  int64 `json:"total_items"`
	TotalPages  int   `json:"total_pages"`
	HasNext     bool  `json:"has_next"`
	HasPrevious bool  `json:"has_previous"`
}
