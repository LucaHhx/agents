package global

import (
	"database/sql/driver"
	"errors"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Map[V comparable, T any] map[V]T

func (m *Map[V, T]) ensureNonNil() {
	if m == nil {
		return
	}
	if *m == nil {
		*m = make(map[V]T)
	}
}

// ---------- SQL Scanner ----------

func (m *Map[V, T]) Scan(value interface{}) error {
	if m == nil {
		return errors.New("Map: Scan on nil receiver")
	}

	switch v := value.(type) {
	case nil:
		*m = make(map[V]T)
		return nil

	case []byte:
		return m.scanBytes(v)

	case string:
		return m.scanBytes([]byte(v))

	default:
		return fmt.Errorf("Map: unsupported Scan type %T", value)
	}
}

func (m *Map[V, T]) scanBytes(b []byte) error {
	// 空值
	if len(b) == 0 {
		*m = make(map[V]T)
		return nil
	}

	// 数据库可能返回 "null"
	if string(b) == "null" {
		*m = make(map[V]T)
		return nil
	}

	// 清空旧数据，避免脏数据
	*m = make(map[V]T)

	if err := HZ_JSON.Unmarshal(b, m); err != nil {
		return err
	}

	m.ensureNonNil()
	return nil
}

// ---------- SQL Valuer ----------

func (m Map[V, T]) Value() (driver.Value, error) {
	// nil map => {}
	if m == nil {
		return []byte("{}"), nil
	}
	return HZ_JSON.Marshal(m)
}

// ---------- GORM datatype ----------

// GormDataTypeInterface
func (Map[V, T]) GormDataType() string { return "json" }

// 按数据库类型返回更准确的 JSON 类型
func (Map[V, T]) GormDBDataType(db *gorm.DB, _ *schema.Field) string {
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
