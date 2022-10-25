package dto

import (
	"encoding/json"
	"net/url"
)

type Request struct {
	Params url.Values
	Body   json.RawMessage
}
