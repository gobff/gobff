package source

import (
	"context"
	"encoding/json"
	"gopkg.in/yaml.v3"
)

type (
	Source interface {
		Run(ctx context.Context, input json.RawMessage) (output json.RawMessage, err error)
	}
	FactoryFunc func(config yaml.Node) (Source, error)
)
