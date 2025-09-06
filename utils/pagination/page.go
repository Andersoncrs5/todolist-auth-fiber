package pagination

type Page[T any] struct {
	Items     []T   `json:"items"`
	Total     int64 `json:"total"`
	PageIndex int   `json:"pageIndex"`
	PageSize  int   `json:"pageSize"`
}
