package metadata

type Page struct {
	Limit  int64 `json:"limit"`
	Offset int64 `json:"offset"`
}

type Condition struct {
	Page    `json:",inline"`
	Filters map[string]interface{} `json:"filters"`
}
