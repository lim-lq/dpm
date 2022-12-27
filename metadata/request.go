package metadata

type PageRequest struct {
	PageNo   int64   `json:"pageNo"`
	PageSize int64   `json:"pageSize"`
	Filters  Filters `json:"filters"`
}
