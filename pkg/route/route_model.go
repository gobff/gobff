package route

import "encoding/json"

type (
	Result struct {
		Error error                       `json:"error,omitempty"`
		Data  map[string]ResourceResponse `json:"data,omitempty"`
	}
	ResourceResultData struct {
		OriginData json.RawMessage
		OutputData json.RawMessage
	}
	ResourceResult struct {
		*ResourceResultData
		Error error
		Alias string
		Omit  bool
	}
	ResourceResponse struct {
		Data  json.RawMessage `json:"data,omitempty"`
		Error error           `json:"error,omitempty"`
	}
)
