package metadata

type BaseResponse struct {
	Code   int64  `json:"code"`
	Status string `json:"status"`
}

type PageData struct {
	PageNo     int64       `json:"pageNo"`
	PageSize   int64       `json:"pageSize"`
	TotalCount int64       `json:"totalCount"`
	Data       interface{} `json:"data"`
}

type PageListResponse struct {
	BaseResponse `json:",inline"`
	Info         PageData `json:"info"`
}

type Response struct {
	BaseResponse `json:",inline"`
	Info         interface{} `json:"info"`
}
