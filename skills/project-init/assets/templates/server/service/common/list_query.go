package common

import (
	"fmt"
	"regexp"
	"server/global"
	modelcommon "server/models/common"
	"strings"

	"gorm.io/gorm"
)

const (
	defaultPageIndex = 1
	defaultPageSize  = 20
	maxPageSize      = 200
)

var columnPattern = regexp.MustCompile(`^[a-z0-9_]+$`)

func NormalizePaging(req modelcommon.PagingReq) (page int, pageSize int) {
	page = req.PageIndex
	if page <= 0 {
		page = defaultPageIndex
	}

	pageSize = req.PageSize
	if pageSize <= 0 {
		pageSize = defaultPageSize
	}
	if pageSize > maxPageSize {
		pageSize = maxPageSize
	}

	return page, pageSize
}

func ApplySearchConditions(q *gorm.DB, search map[string]modelcommon.SearchCondition) *gorm.DB {
	for key, s := range search {
		col := normalizeColumnName(key)
		if col == "" {
			continue
		}

		switch s.Operator {
		case modelcommon.OperatorEqual:
			if start, end, ok := global.ParseTimeRangeFromFilterValue(s.Value); ok {
				q = q.Where(fmt.Sprintf("(%s >= ? AND %s < ?)", col, col), start, end)
			} else {
				q = q.Where(fmt.Sprintf("%s = ?", col), s.Value)
			}

		case modelcommon.OperatorNotEqual, "!=":
			if start, end, ok := global.ParseTimeRangeFromFilterValue(s.Value); ok {
				q = q.Where(fmt.Sprintf("(%s < ? OR %s >= ?)", col, col), start, end)
			} else {
				q = q.Where(fmt.Sprintf("%s <> ?", col), s.Value)
			}

		case modelcommon.OperatorGreaterThan:
			q = q.Where(fmt.Sprintf("%s > ?", col), s.Value)
		case modelcommon.OperatorGreaterThanOrEqual:
			q = q.Where(fmt.Sprintf("%s >= ?", col), s.Value)
		case modelcommon.OperatorLessThan:
			q = q.Where(fmt.Sprintf("%s < ?", col), s.Value)
		case modelcommon.OperatorLessThanOrEqual:
			q = q.Where(fmt.Sprintf("%s <= ?", col), s.Value)

		case modelcommon.OperatorLike:
			q = q.Where(fmt.Sprintf("%s LIKE ?", col), fmt.Sprintf("%%%v%%", s.Value))
		case "startsWith":
			q = q.Where(fmt.Sprintf("%s LIKE ?", col), fmt.Sprintf("%v%%", s.Value))
		case "endsWith":
			q = q.Where(fmt.Sprintf("%s LIKE ?", col), fmt.Sprintf("%%%v", s.Value))

		case modelcommon.OperatorIn:
			q = q.Where(fmt.Sprintf("%s IN ?", col), s.Value)
		case modelcommon.OperatorNotIn, "notIn":
			q = q.Where(fmt.Sprintf("%s NOT IN ?", col), s.Value)

		case modelcommon.OperatorBetween:
			if arr, ok := toSlice2(s.Value); ok {
				q = q.Where(fmt.Sprintf("%s BETWEEN ? AND ?", col), arr[0], arr[1])
			}
		case "notBetween":
			if arr, ok := toSlice2(s.Value); ok {
				q = q.Where(fmt.Sprintf("%s NOT BETWEEN ? AND ?", col), arr[0], arr[1])
			}

		case modelcommon.OperatorIsNull:
			q = q.Where(fmt.Sprintf("%s IS NULL", col))
		case modelcommon.OperatorIsNotNull:
			q = q.Where(fmt.Sprintf("%s IS NOT NULL", col))
		}
	}

	return q
}

func ApplySortConditions(q *gorm.DB, sorteds []modelcommon.OrderCondition, defaultOrder string) *gorm.DB {
	for _, o := range sorteds {
		col := normalizeColumnName(o.Field)
		if col == "" {
			continue
		}
		if o.Order == modelcommon.OrderDesc {
			q = q.Order(col + " DESC")
		} else {
			q = q.Order(col + " ASC")
		}
	}

	if len(sorteds) == 0 && strings.TrimSpace(defaultOrder) != "" {
		q = q.Order(defaultOrder)
	}

	return q
}

func normalizeColumnName(input string) string {
	col := strings.ToLower(strings.TrimSpace(input))
	if !columnPattern.MatchString(col) {
		return ""
	}
	return col
}

func toSlice2(v interface{}) ([2]interface{}, bool) {
	var out [2]interface{}
	switch x := v.(type) {
	case []interface{}:
		if len(x) >= 2 {
			out[0], out[1] = x[0], x[1]
			return out, true
		}
	case [2]interface{}:
		return x, true
	default:
		return out, false
	}
	return out, false
}
