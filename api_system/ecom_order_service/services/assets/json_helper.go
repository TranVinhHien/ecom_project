package assets_services

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// NullableJSON is a custom type that can handle NULL JSON values from database
type NullableJSON struct {
	Data  json.RawMessage
	Valid bool // Valid is true if Data is not NULL
}

// Scan implements the Scanner interface for database/sql
func (nj *NullableJSON) Scan(value interface{}) error {
	if value == nil {
		nj.Data, nj.Valid = json.RawMessage("[]"), false
		return nil
	}
	nj.Valid = true
	switch v := value.(type) {
	case []byte:
		nj.Data = json.RawMessage(v)
		return nil
	case string:
		nj.Data = json.RawMessage(v)
		return nil
	default:
		return errors.New("failed to scan NullableJSON")
	}
}

// Value implements the driver Valuer interface
func (nj NullableJSON) Value() (driver.Value, error) {
	if !nj.Valid {
		return nil, nil
	}
	return []byte(nj.Data), nil
}

// MarshalJSON implements json.Marshaler
func (nj NullableJSON) MarshalJSON() ([]byte, error) {
	if !nj.Valid || len(nj.Data) == 0 {
		return []byte("[]"), nil
	}
	return nj.Data, nil
}

// UnmarshalJSON implements json.Unmarshaler
func (nj *NullableJSON) UnmarshalJSON(data []byte) error {
	nj.Valid = true
	nj.Data = json.RawMessage(data)
	return nil
}
