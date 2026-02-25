// date.go
package global

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

const (
	DateLayout = "2006-01-02"
)

type Date struct {
	time.Time
}

func NewDate(t time.Time) *Date {
	d := &Date{}
	d.Set(t)
	return d
}

func Today() Date {
	d := Date{}
	d.Set(time.Now())
	return d
}

func TodayPtr() *Date {
	d := &Date{}
	d.Set(time.Now())
	return d
}

// Set 会把时间归一化到当天 00:00:00（本地时区）
func (d *Date) Set(t time.Time) {
	if d == nil {
		return
	}
	if t.IsZero() {
		d.Time = time.Time{}
		return
	}
	tt := t.In(time.Local)
	d.Time = time.Date(tt.Year(), tt.Month(), tt.Day(), 0, 0, 0, 0, time.Local)
}

func (d *Date) GetTime() time.Time {
	if d == nil {
		return time.Time{}
	}
	return d.Time
}

func (d *Date) AddDays(days int) *Date {
	if d == nil {
		return nil
	}
	nd := &Date{}
	nd.Set(d.Time.AddDate(0, 0, days))
	return nd
}

func (d *Date) Format() string {
	if d == nil || d.Time.IsZero() {
		return ""
	}
	return d.Time.Format(DateLayout)
}

// ---------- JSON ----------

func (d *Date) UnmarshalJSON(b []byte) error {
	if d == nil {
		return fmt.Errorf("Date: UnmarshalJSON on nil receiver")
	}
	s := strings.Trim(string(b), `"`)
	if s == "" || s == "null" {
		d.Time = time.Time{}
		return nil
	}

	// 优先按 date 解析
	if tt, err := time.ParseInLocation(DateLayout, s, time.Local); err == nil {
		d.Set(tt)
		return nil
	}

	// 兼容传了 datetime 的情况（只取日期部分）
	if tt, err := time.ParseInLocation(DateTimeLayout, s, time.Local); err == nil {
		d.Set(tt)
		return nil
	}

	// 兼容 RFC3339
	if tt, err := time.Parse(time.RFC3339, s); err == nil {
		d.Set(tt.In(time.Local))
		return nil
	}

	return fmt.Errorf("Date: invalid date format: %q", s)
}

func (d Date) MarshalJSON() ([]byte, error) {
	if d.Time.IsZero() {
		return []byte(`null`), nil
	}
	return []byte(fmt.Sprintf(`"%s"`, d.Time.Format(DateLayout))), nil
}

// ---------- SQL (Valuer / Scanner) ----------

func (d Date) Value() (driver.Value, error) {
	if d.Time.IsZero() {
		return nil, nil
	}
	// 对 DATE 字段，返回 YYYY-MM-DD 字符串更直观，也更不容易被时区影响
	return d.Time.Format(DateLayout), nil
}

func (d *Date) Scan(v interface{}) error {
	if d == nil {
		return fmt.Errorf("Date: Scan on nil receiver")
	}

	switch val := v.(type) {
	case time.Time:
		d.Set(val)
		return nil
	case []byte:
		s := string(val)
		if s == "" {
			d.Time = time.Time{}
			return nil
		}
		// MySQL DATE 通常就是 YYYY-MM-DD
		if tt, err := time.ParseInLocation(DateLayout, s, time.Local); err == nil {
			d.Set(tt)
			return nil
		}
		// 兼容有些驱动返回 datetime
		if tt, err := time.ParseInLocation(DateTimeLayout, s, time.Local); err == nil {
			d.Set(tt)
			return nil
		}
		return fmt.Errorf("Date: cannot parse %q", s)
	case string:
		if val == "" {
			d.Time = time.Time{}
			return nil
		}
		if tt, err := time.ParseInLocation(DateLayout, val, time.Local); err == nil {
			d.Set(tt)
			return nil
		}
		if tt, err := time.ParseInLocation(DateTimeLayout, val, time.Local); err == nil {
			d.Set(tt)
			return nil
		}
		return fmt.Errorf("Date: cannot parse %q", val)
	case nil:
		d.Time = time.Time{}
		return nil
	default:
		return fmt.Errorf("Date: cannot scan type %T", v)
	}
}

// ---------- GORM datatype ----------

// GormDataTypeInterface
func (Date) GormDataType() string { return "date" }

func (Date) GormDBDataType(db *gorm.DB, _ *schema.Field) string {
	switch db.Dialector.Name() {
	case "mysql":
		return "DATE"
	case "postgres":
		return "DATE"
	case "sqlite":
		return "DATE"
	default:
		return "DATE"
	}
}
