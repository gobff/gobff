package route

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gobff/gobff/resource"
	"github.com/gobff/gobff/transformer"
	"io"
	"net/http"
)

type (
	Route interface {
		Run(c *gin.Context)
	}
	Resource struct {
		resource.Resource
		Transformer transformer.Transformer
		As          string
	}
	Resources map[string]Resource
	routeImpl struct {
		path      string
		resources Resources
	}
)

func New(path string, resources Resources) Route {
	return &routeImpl{
		path:      path,
		resources: resources,
	}
}

func (r *routeImpl) Run(c *gin.Context) {
	var input json.RawMessage
	if err := c.Bind(&input); err != nil && err != io.EOF {
		c.JSON(http.StatusInternalServerError, Response{Err: err})
		return
	}

	mapChanResponse := r.executeResources(c, input)
	defer closeResponseChannels(mapChanResponse)

	requestResponse := make(map[string]Response)
	for name, chanResponse := range mapChanResponse {
		requestResponse[r.resources[name].As] = <-chanResponse
	}
	c.IndentedJSON(http.StatusOK, requestResponse)
}

func (r *routeImpl) executeResources(ctx context.Context, input json.RawMessage) map[string]chan Response {
	responses := make(map[string]chan Response)
	for name, res := range r.resources {
		name, res, cResponse := name, res, make(chan Response)
		go executeResource(ctx, cResponse, res, input)
		responses[name] = cResponse
	}
	return responses
}

func executeResource(ctx context.Context, cResponse chan Response, res Resource, input json.RawMessage) {
	data, err := res.Run(ctx, input)
	if err != nil {
		cResponse <- Response{Err: err}
		return
	}

	if res.Transformer != nil {
		data, err = res.Transformer.Transform(data)
		if err != nil {
			cResponse <- Response{Err: err}
			return
		}
	}

	cResponse <- Response{
		Data: data,
	}
}

func closeResponseChannels(responses map[string]chan Response) {
	for _, cResponse := range responses {
		close(cResponse)
	}
}
