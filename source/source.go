package source

import (
	"context"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
)

const (
	KindHttp = "http"
)

type (
	Params map[string]string
	Source interface {
		Run(ctx context.Context, params Params, input json.RawMessage) (output json.RawMessage, err error)
		ValidateParams(params Params) error
	}
)

func GetSource(kind string, config yaml.Node) (Source, error) {
	switch kind {
	case KindHttp:
		return newSourceHTTP(config)
	default:
		return nil, fmt.Errorf("source not found: %s", kind)
	}
}
