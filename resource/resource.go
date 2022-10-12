package resource

import (
	"context"
	"encoding/json"
	"github.com/gondalf/gondalf/source"
)

type (
	Resource interface {
		Run(ctx context.Context, input json.RawMessage) (json.RawMessage, error)
	}
	resource struct {
		source source.Source
	}
)

func NewResource(source source.Source) Resource {
	r := &resource{source: source}
	return r
}

func (r resource) Run(ctx context.Context, input json.RawMessage) (json.RawMessage, error) {
	return r.source.Run(ctx, input)
}
