package metadata

type BaseResponse struct {
	Code   int64  `json:"code"`
	Status string `json:"status"`
}

type PageData struct {
	Limit  int64         `json:"limig"`
	Offset int64         `json:"offset"`
	Info   []interface{} `json:"info"`
}

type ListResponse struct {
	BaseResponse `json:",inline"`
	PageData     `json:",inline"`
}

type Response struct {
	BaseResponse `json:",inline"`
	Info         interface{} `json:"info"`
}
