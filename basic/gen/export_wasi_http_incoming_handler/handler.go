package export_wasi_http_incoming_handler

import (
	httptypes "componentize-basic/gen/wasi_http_types"
	"net/http"
)

var handler http.HandlerFunc

// Handle the specified `Request`, returning a `Response`
func Handle(request *httptypes.IncomingRequest, responseOut *httptypes.ResponseOutparam) {
	req := &http.Request{}
	res := newResponseWriter(responseOut)
	handler(res, req)
}

func HandleFunc(h http.HandlerFunc) {
	handler = h
}
