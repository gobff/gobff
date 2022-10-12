package route

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gondalf/gondalf/resource"
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

	results := make(map[string]chan Response)
	for name, res := range r.resources {
		name, res := name, res
		cResponse := make(chan Response)
		go func() {
			var response Response
			response.Data, response.Err = res.Run(c, input)
			cResponse <- response
		}()
		results[name] = cResponse
	}

	response := make(map[string]Response)
	for name, cResult := range results {
		response[r.resources[name].As] = <-cResult
	}
	c.IndentedJSON(http.StatusOK, response)
}
