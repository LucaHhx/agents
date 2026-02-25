// datetime.go
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
	DateTimeLayout        = "2006-01-02 15:04:05"
	dateTimeTMinuteLayout = "2006-01-02T15:04"
	dateTimeMinuteLayout  = "2006-01-02 15:04"
	dateTimeTSecondLayout = "2006-01-02T15:04:05"
)

type DateTime struct {
	time.Time
}

func NewDateTime(t time.Time) *DateTime { return &DateTime{Time: t} }
func NowDateTime() DateTime             { return DateTime{Time: time.Now()} }
func NowDateTimePtr() *DateTime         { return &DateTime{Time: time.Now()} }

func (t *DateTime) GetTime() time.Time {
	if t == nil {
		return time.Time{}
	}
	return t.Time
}

func (t *DateTime) AddDays(days int) *DateTime {
	if t == nil {
		return nil
	}
	return &DateTime{Time: t.Time.AddDate(0, 0, days)}
}

func (t *DateTime) Format() string {
	if t == nil || t.Time.IsZero() {
		return ""
	}
	return t.Time.Format(DateTimeLayout)
}

// ---------- JSON ----------

func (t *DateTime) UnmarshalJSON(b []byte) error {
	if t == nil {
		return fmt.Errorf("DateTime: UnmarshalJSON on nil receiver")
	}
	s := strings.Trim(string(b), `"`)
	if s == "" || s == "null" {
		t.Time = time.Time{}
		return nil
	}

	tt, err := parseDateTimeValue(s)
	if err != nil {
		return err
	}
	t.Time = tt
	return nil
}

func (t DateTime) MarshalJSON() ([]byte, error) {
	if t.Time.IsZero() {
		return []byte(`null`), nil
	}
	return []byte(fmt.Sprintf(`"%s"`, t.Time.Format(DateTimeLayout))), nil
}

// ---------- SQL (Valuer / Scanner) ----------

func (t DateTime) Value() (driver.Value, error) {
	if t.Time.IsZero() {
		return nil, nil
	}
	// 直接返回 time.Time 最兼容（驱动自己处理）
	return t.Time, nil
}

func (t *DateTime) Scan(v interface{}) error {
	if t == nil {
		return fmt.Errorf("DateTime: Scan on nil receiver")
	}

	switch val := v.(type) {
	case time.Time:
		t.Time = val
		return nil
	case []byte:
		s := string(val)
		if s == "" {
			t.Time = time.Time{}
			return nil
		}
		tt, err := parseDateTimeValue(s)
		if err != nil {
			return err
		}
		t.Time = tt
		return nil
	case string:
		if val == "" {
			t.Time = time.Time{}
			return nil
		}
		tt, err := parseDateTimeValue(val)
		if err != nil {
			return err
		}
		t.Time = tt
		return nil
	case nil:
		t.Time = time.Time{}
		return nil
	default:
		return fmt.Errorf("DateTime: cannot scan type %T", v)
	}
}

func parseDateTimeValue(value string) (time.Time, error) {
	v := strings.TrimSpace(value)
	if v == "" {
		return time.Time{}, fmt.Errorf("DateTime: empty datetime")
	}

	if tt, err := time.ParseInLocation(DateTimeLayout, v, time.Local); err == nil {
		return tt, nil
	}
	if tt, err := time.ParseInLocation(dateTimeTSecondLayout, v, time.Local); err == nil {
		return tt, nil
	}
	if tt, err := time.ParseInLocation(dateTimeTMinuteLayout, v, time.Local); err == nil {
		return tt, nil
	}
	if tt, err := time.ParseInLocation(dateTimeMinuteLayout, v, time.Local); err == nil {
		return tt, nil
	}
	if tt, err := time.ParseInLocation(DateLayout, v, time.Local); err == nil {
		return tt, nil
	}
	if tt, err := time.Parse(time.RFC3339, v); err == nil {
		return tt.In(time.Local), nil
	}

	return time.Time{}, fmt.Errorf("DateTime: invalid datetime format: %q", value)
}

// ParseTimeRangeFromFilterValue parses date/datetime filter values and returns
// a half-open range [start, end). It supports:
// - 2006-01-02 -> +1 day
// - 2006-01-02T15:04 / 2006-01-02 15:04 -> +1 minute
// - 2006-01-02T15:04:05 / 2006-01-02 15:04:05 / RFC3339 -> +1 second
func ParseTimeRangeFromFilterValue(value any) (start time.Time, end time.Time, ok bool) {
	raw, ok := value.(string)
	if !ok {
		return time.Time{}, time.Time{}, false
	}

	v := strings.TrimSpace(raw)
	if v == "" {
		return time.Time{}, time.Time{}, false
	}

	if tt, err := time.ParseInLocation(DateLayout, v, time.Local); err == nil {
		return tt, tt.AddDate(0, 0, 1), true
	}
	if tt, err := time.ParseInLocation(dateTimeTMinuteLayout, v, time.Local); err == nil {
		return tt, tt.Add(time.Minute), true
	}
	if tt, err := time.ParseInLocation(dateTimeMinuteLayout, v, time.Local); err == nil {
		return tt, tt.Add(time.Minute), true
	}
	if tt, err := time.ParseInLocation(dateTimeTSecondLayout, v, time.Local); err == nil {
		return tt, tt.Add(time.Second), true
	}
	if tt, err := time.ParseInLocation(DateTimeLayout, v, time.Local); err == nil {
		return tt, tt.Add(time.Second), true
	}
	if tt, err := time.Parse(time.RFC3339, v); err == nil {
		lt := tt.In(time.Local)
		return lt, lt.Add(time.Second), true
	}

	return time.Time{}, time.Time{}, false
}

// ---------- GORM datatype ----------

// GormDataTypeInterface
func (DateTime) GormDataType() string { return "datetime" }

// 可选：更精确地按不同 DB 返回类型（推荐）
func (DateTime) GormDBDataType(db *gorm.DB, _ *schema.Field) string {
	switch db.Dialector.Name() {
	case "mysql":
		return "DATETIME"
	case "postgres":
		return "TIMESTAMP"
	case "sqlite":
		return "DATETIME"
	default:
		return "DATETIME"
	}
}
