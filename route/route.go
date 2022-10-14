package route

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

type (
	Route interface {
		Run(c *gin.Context)
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
	var response Response
	response.Data, response.Err = res.Run(ctx, input)
	cResponse <- response
}

func closeResponseChannels(responses map[string]chan Response) {
	for _, cResponse := range responses {
		close(cResponse)
	}
}
