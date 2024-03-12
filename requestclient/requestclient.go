package requestclient

import "net/http"

type RequestClient interface {
	Do(*http.Request) (*http.Response, error)
}
