package global

import (
	"database/sql/driver"
	"errors"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Data[T any] struct {
	V *T
}

func NewData[T any](v T) Data[T] { return Data[T]{V: &v} }
func (d Data[T]) IsNil() bool    { return d.V == nil }
func (d Data[T]) Ptr() *T        { return d.V }

func (d Data[T]) Get() (T, bool) {
	if d.V == nil {
		var zero T
		return zero, false
	}
	return *d.V, true
}

func (d *Data[T]) Set(v T) {
	if d == nil {
		return
	}
	d.V = &v
}

func (d *Data[T]) Clear() {
	if d == nil {
		return
	}
	d.V = nil
}

// ---------- JSON ----------

func (d Data[T]) MarshalJSON() ([]byte, error) {
	if d.V == nil {
		return []byte("null"), nil
	}
	return HZ_JSON.Marshal(*d.V)
}

func (d *Data[T]) UnmarshalJSON(b []byte) error {
	if d == nil {
		return errors.New("Data: UnmarshalJSON on nil receiver")
	}
	if len(b) == 0 || string(b) == "null" {
		d.V = nil
		return nil
	}
	var v T
	if err := HZ_JSON.Unmarshal(b, &v); err != nil {
		return err
	}
	d.V = &v
	return nil
}

// ---------- SQL (Valuer / Scanner) ----------

func (d Data[T]) Value() (driver.Value, error) {
	if d.V == nil {
		return nil, nil
	}
	// 作为 JSON 单值存储（适合 JSON/JSONB 列）
	return HZ_JSON.Marshal(*d.V)
}

func (d *Data[T]) Scan(value interface{}) error {
	if d == nil {
		return errors.New("Data: Scan on nil receiver")
	}

	switch v := value.(type) {
	case nil:
		d.V = nil
		return nil
	case []byte:
		return d.scanBytes(v)
	case string:
		return d.scanBytes([]byte(v))
	default:
		return fmt.Errorf("Data: unsupported Scan type %T", value)
	}
}

func (d *Data[T]) scanBytes(b []byte) error {
	if len(b) == 0 || string(b) == "null" {
		d.V = nil
		return nil
	}
	var v T
	if err := HZ_JSON.Unmarshal(b, &v); err != nil {
		return err
	}
	d.V = &v
	return nil
}

// ---------- GORM datatype ----------

func (Data[T]) GormDataType() string { return "json" }

func (Data[T]) GormDBDataType(db *gorm.DB, _ *schema.Field) string {
	switch db.Dialector.Name() {
	case "mysql":
		return "JSON"
	case "postgres":
		return "JSONB"
	case "sqlite":
		return "JSON"
	default:
		return "JSON"
	}
}
