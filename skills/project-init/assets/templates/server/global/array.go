package global

import (
	"database/sql/driver"
	"errors"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Arr[T any] []T

func (a *Arr[T]) ensureNonNil() {
	if a == nil {
		return
	}
	if *a == nil {
		*a = make([]T, 0)
	}
}

func (a *Arr[T]) Scan(value interface{}) error {
	if a == nil {
		return errors.New("Arr: Scan on nil receiver")
	}

	switch v := value.(type) {
	case nil:
		*a = make([]T, 0)
		return nil

	case []byte:
		return a.scanBytes(v)

	case string:
		return a.scanBytes([]byte(v))

	default:
		return fmt.Errorf("Arr: unsupported Scan type %T", value)
	}
}

func (a *Arr[T]) scanBytes(b []byte) error {
	// 空值/空串
	if len(b) == 0 {
		*a = make([]T, 0)
		return nil
	}

	// 兼容数据库可能返回 "null"
	if string(b) == "null" {
		*a = make([]T, 0)
		return nil
	}

	// 先清空，避免残留旧数据
	*a = make([]T, 0)

	if err := HZ_JSON.Unmarshal(b, a); err != nil {
		return err
	}

	a.ensureNonNil()
	return nil
}

func (a Arr[T]) Value() (driver.Value, error) {
	// nil slice => []，避免落库成 null
	if a == nil {
		return []byte("[]"), nil
	}
	return HZ_JSON.Marshal(a)
}

// ---------- GORM datatype ----------

// GormDataTypeInterface
func (Arr[T]) GormDataType() string { return "json" }

// 可选：按不同数据库返回更合适的类型
func (Arr[T]) GormDBDataType(db *gorm.DB, _ *schema.Field) string {
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
