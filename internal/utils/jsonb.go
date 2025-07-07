package utils

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type JSONB[T any] struct {
	Data T
}

func (j JSONB[T]) Value() (driver.Value, error) {
	return json.Marshal(j.Data)
}

func (j *JSONB[T]) Scan(value interface{}) error {
	if value == nil {
		var zero T
		j.Data = zero
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan JSONB: type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, &j.Data)
}
