package route

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gobff/gobff/resource"
	"io"
	"net/http"
)

type (
	Route interface {
		Run(c *gin.Context)
	}
	Resource struct {
		resource.Resource
		As string
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
		name, res := name, res
		cResponse := make(chan Response)
		go func() {
			var response Response
			response.Data, response.Err = res.Run(ctx, input)
			cResponse <- response
		}()
		responses[name] = cResponse
	}
	return responses
}

func closeResponseChannels(responses map[string]chan Response) {
	for _, cResponse := range responses {
		close(cResponse)
	}
}
