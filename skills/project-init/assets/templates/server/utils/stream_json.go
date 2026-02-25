package utils

import (
	"io"
	"server/global"
	"server/models/common"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetNDJSONHeaders(c *gin.Context) {
	c.Header("Content-Type", "application/x-ndjson; charset=utf-8")
	c.Header("Cache-Control", "no-store")
	c.Header("X-Accel-Buffering", "no")
}

func WriteJSONLine(w io.Writer, v any) error {
	data, err := global.HZ_JSON.Marshal(v)
	if err != nil {
		return err
	}
	data = append(data, '\n')
	_, err = w.Write(data)
	return err
}

func StructToMap(v any) (map[string]any, error) {
	data, err := global.HZ_JSON.Marshal(v)
	if err != nil {
		return nil, err
	}
	out := make(map[string]any)
	if err := global.HZ_JSON.Unmarshal(data, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func StreamQueryAsNDJSON[T any](q *gorm.DB, fields []common.FileField, writer io.Writer) error {
	rows, err := q.Rows()
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var item T
		if err := q.ScanRows(rows, &item); err != nil {
			return err
		}
		m, err := StructToMap(item)
		if err != nil {
			return err
		}
		row := BuildRowByFields(m, fields)
		if err := WriteJSONLine(writer, row); err != nil {
			return err
		}
	}

	return nil
}
