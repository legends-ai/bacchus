package handlers

import (
	"net/http"

	"github.com/simplyianm/httputil"
	"github.com/simplyianm/inject"
)

// SpiderRequest is a request to spider.
type SpiderRequest struct {
	Region string `json:"region"`
}

// SpiderResponse is the response.
type SpiderResponse struct {
	Success bool `json:""`
}

// Validate validates if the request should be processed.
func (r *SpiderRequest) Validate() error {
	// TODO(simplyianm): check if region is valid
	return nil
}

// Handle handles the request
func (r *SpiderRequest) Handle(w http.ResponseWriter) {
	httputil.WriteJSON(w, resp)
}

func newSpiderRequest(injector inject.Injector) func() httputil.Request {
	return func() httputil.Request {
		ret := &SpiderRequest{}
		injector.Apply(ret)
		return ret
	}
}
