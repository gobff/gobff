package route

import "encoding/json"

type Response struct {
	Data json.RawMessage `json:"data,omitempty"`
	Err  error           `json:"error,omitempty"`
}
