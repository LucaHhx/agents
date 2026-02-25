package common

type FileField struct {
	Label string `json:"label"`
	Key   string `json:"key"`
	Type  string `json:"type"`
	DBKey string `json:"dbKey"`
}

type ExportReq struct {
	Entity   string                     `json:"entity"`
	Filename string                     `json:"filename"`
	Fields   []FileField                `json:"fields"`
	Search   map[string]SearchCondition `json:"search"`
	Sorteds  []OrderCondition           `json:"sorteds"`
}

type ImportReq struct {
	Entity string                   `json:"entity" form:"entity"`
	Reset  bool                     `json:"reset" form:"reset"`
	Rows   []map[string]interface{} `json:"rows" form:"rows"`
}

type ImportRowError struct {
	Row     int    `json:"row"`
	Message string `json:"message"`
}

type ImportResp struct {
	Entity      string           `json:"entity"`
	Mode        string           `json:"mode"`
	TotalRows   int              `json:"totalRows"`
	SuccessRows int              `json:"successRows"`
	FailedRows  int              `json:"failedRows"`
	Errors      []ImportRowError `json:"errors"`
}
