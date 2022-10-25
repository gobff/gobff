package transformer

import (
	"encoding/json"
	"github.com/carlosrodriguesf/gobff/tool/transformer/filter"
	"github.com/carlosrodriguesf/gobff/tool/transformer/parser"
)

type (
	Transformer interface {
		Transform(input json.RawMessage) (output json.RawMessage, err error)
	}
	transformer struct {
		fields []parser.Field
	}
)

func New(pattern string) (Transformer, error) {
	fields, err := parser.GetFieldsFromPattern(pattern)
	if err != nil {
		return nil, err
	}
	return transformer{
		fields: fields,
	}, nil
}

func (t transformer) Transform(input json.RawMessage) (json.RawMessage, error) {
	var root any
	if err := json.Unmarshal(input, &root); err != nil {
		return nil, err
	}

	result, err := filter.ResolveFields(root, t.fields)
	if err != nil {
		return nil, err
	}

	input, err = json.Marshal(result)
	if err != nil {
		return nil, err
	}

	return input, nil
}
