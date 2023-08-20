package isbn

import (
	"database/sql/driver"
	"reflect"
)

// into SQL
func (isbn ISBN) Value() (driver.Value, error) {
	return isbn[:], nil
}

// from SQL
func (isbn *ISBN) Scan(v any) (err error) {
	switch u := v.(type) {
	case []byte: // sqlite3 & pgx
		*isbn, err = ParseBytes(u)
	case string: // driver.Value is []byte therefore this won't run
		*isbn, err = Parse(u)
	default:
		return invalidTypeError{reflect.TypeOf(v)}
	}
	return err
}
