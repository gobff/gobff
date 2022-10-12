package resource

import (
	"context"
	"encoding/json"
)

type Resource interface {
	Run(ctx context.Context, input json.RawMessage) (output json.RawMessage, err error)
}
