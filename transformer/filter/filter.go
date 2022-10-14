package filter

import (
	"fmt"
	"github.com/gobff/gobff/transformer/parser"
	"reflect"
)

func ResolveFields(data any, fields []parser.Field) (any, error) {
	switch reflect.TypeOf(data).Kind() {
	case reflect.Array, reflect.Slice:
		return resolveFieldsInArray(data.([]any), fields)
	case reflect.Map:
		return resolveFieldsInMap(data.(map[string]any), fields)
	default:
		return nil, fmt.Errorf("invalid kind")
	}
}

func resolveFieldsInArray(data []any, fields []parser.Field) ([]any, error) {
	arrLen := len(data)
	result := make([]any, arrLen)
	for i := 0; i < arrLen; i++ {
		value, err := ResolveFields(data[i], fields)
		if err != nil {
			return nil, err
		}
		result[i] = value
	}
	return result, nil
}

func resolveFieldsInMap(data map[string]any, fields []parser.Field) (map[string]any, error) {
	var (
		err    error
		result = make(map[string]any)
	)
	for _, field := range fields {
		value, found := data[field.Key]
		if !found {
			return nil, fmt.Errorf("key not found: %s", value)
		}
		if len(field.Children) > 0 {
			value, err = ResolveFields(value, field.Children)
			if err != nil {
				return nil, err
			}
		}
		result[field.Key] = value
	}
	return result, nil
}
