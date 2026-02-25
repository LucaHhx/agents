package utils

import (
	"errors"
	"server/global"
	"server/models/common"
	"strings"
)

func NormalizeFileFields(fields []common.FileField) []common.FileField {
	out := make([]common.FileField, 0, len(fields))
	seen := make(map[string]struct{}, len(fields))
	for _, field := range fields {
		key := strings.TrimSpace(field.Key)
		if key == "" || key == "file" {
			continue
		}
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		field.Key = key
		out = append(out, field)
	}
	return out
}

func BuildRowByFields(source map[string]any, fields []common.FileField) map[string]any {
	row := make(map[string]any, len(fields))
	for _, field := range fields {
		row[field.Key] = normalizeCell(source[field.Key])
	}
	return row
}

func MapToStruct(row map[string]any, out any) error {
	data, err := global.HZ_JSON.Marshal(row)
	if err != nil {
		return err
	}
	if err := global.HZ_JSON.Unmarshal(data, out); err != nil {
		return errors.New("invalid row data: " + err.Error())
	}
	return nil
}

func MapToModel[T any](row map[string]any, ignoreKeys ...string) (T, error) {
	cloned := make(map[string]any, len(row))
	for k, v := range row {
		cloned[k] = v
	}
	DeleteKeys(cloned, ignoreKeys...)
	var out T
	if err := MapToStruct(cloned, &out); err != nil {
		return out, err
	}
	return out, nil
}

func DeleteKeys(m map[string]any, keys ...string) {
	for _, key := range keys {
		delete(m, key)
	}
}

func normalizeCell(v any) any {
	if v == nil {
		return ""
	}
	switch x := v.(type) {
	case map[string]any, []any:
		data, err := global.HZ_JSON.Marshal(x)
		if err != nil {
			return ""
		}
		return string(data)
	default:
		return x
	}
}
