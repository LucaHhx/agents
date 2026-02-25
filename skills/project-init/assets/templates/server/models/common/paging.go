package common

type Operator string

const (
	OperatorEqual              Operator = "="
	OperatorLike               Operator = "like"
	OperatorIn                 Operator = "in"
	OperatorNotIn              Operator = "not in"
	OperatorBetween            Operator = "between"
	OperatorNotEqual           Operator = "<>"
	OperatorGreaterThan        Operator = ">"
	OperatorGreaterThanOrEqual Operator = ">="
	OperatorLessThan           Operator = "<"
	OperatorLessThanOrEqual    Operator = "<="
	OperatorIsNull             Operator = "is null"
	OperatorIsNotNull          Operator = "is not null"
)

type Order string

const (
	OrderAsc  Order = "asc"
	OrderDesc Order = "desc"
)

type SearchCondition struct {
	Operator Operator `json:"operator"` // 操作符
	Value    any      `json:"value"`    // 值
}

type OrderCondition struct {
	Field string `json:"field"`
	Order Order  `json:"order"`
}

type PagingReq struct {
	PageIndex int                        `json:"pageIndex"`
	PageSize  int                        `json:"pageSize"`
	Search    map[string]SearchCondition `json:"search"`
	Sorteds   []OrderCondition           `json:"sorteds"`
}

type PagingResp[T any] struct {
	List     []T `json:"list"`
	Total    int `json:"total"`
	Page     int `json:"pageIndex"`
	PageSize int `json:"pageSize"`
}
