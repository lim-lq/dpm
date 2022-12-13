package metadata

type Filters map[string]interface{}

type Page struct {
	Limit  int64 `json:"limit"`
	Offset int64 `json:"offset"`
}

type Condition struct {
	Page    `json:",inline"`
	Filters Filters `json:"filters"`
}
