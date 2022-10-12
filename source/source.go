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

type Source interface {
	Run(ctx context.Context, input json.RawMessage) (output json.RawMessage, err error)
}

func GetSource(kind string, config yaml.Node) (Source, error) {
	switch kind {
	case KindHttp:
		return newSourceHTTP(config)
	default:
		return nil, fmt.Errorf("source not found: %s", kind)
	}
}
