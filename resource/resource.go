package resource

import (
	"context"
	"encoding/json"
	"github.com/gondalf/gondalf/source"
	"github.com/jellydator/ttlcache/v3"
	"sync"
)

type (
	Result struct {
		Data  json.RawMessage
		Error error
	}
	ChanResult chan Result
	Resource   interface {
		Run(ctx context.Context, input json.RawMessage) ChanResult
	}
	resource struct {
		source  source.Source
		mutex   *sync.Mutex
		cache   ttlcache.Cache[string, json.RawMessage]
		running bool
	}
)

func NewResource(source source.Source) Resource {
	r := &resource{source: source}
	return r
}

func (r resource) Run(ctx context.Context, input json.RawMessage) ChanResult {
	cResult := make(ChanResult)
	go func() {
		var result Result
		result.Data, result.Error = r.source.Run(ctx, input)
		cResult <- result
	}()
	return cResult

}
