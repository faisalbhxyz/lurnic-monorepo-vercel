package utils

import (
	"encoding/json"
	"reflect"
	"strings"
	"time"

	"gorm.io/datatypes"
)

func EmptyStringToNil(s *string) *string {
	if s != nil && strings.TrimSpace(*s) == "" {
		return nil
	}
	return s
}

func ZeroToNil[T any](v *T) *T {
	if v == nil {
		return nil
	}

	switch val := any(*v).(type) {
	case string:
		if strings.TrimSpace(val) == "" {
			return nil
		}
	case bool:
		return v
	case float32, float64, int, int64, uint, uint64:
		zero := reflect.Zero(reflect.TypeOf(val)).Interface()
		if reflect.DeepEqual(val, zero) {
			return nil
		}
	case time.Time:
		if val.IsZero() {
			return nil
		}
	case []string:
		if len(val) == 0 {
			return nil
		}
	default:
		// For unsupported types, just return as-is
		return v
	}

	return v
}

func NormalizeTags(rawTags *[]string) (datatypes.JSON, error) {
	if rawTags == nil || len(*rawTags) == 0 {
		return datatypes.JSON([]byte("[]")), nil
	}

	// The first element is a JSON string itself, so unmarshal it:
	var tags []string
	err := json.Unmarshal([]byte((*rawTags)[0]), &tags)
	if err != nil {
		// fallback: marshal rawTags normally if it fails
		b, _ := json.Marshal(rawTags)
		return datatypes.JSON(b), err
	}

	// Now marshal tags back to JSON
	b, err := json.Marshal(tags)
	if err != nil {
		return datatypes.JSON(b), err
	}

	return datatypes.JSON(b), nil
}
