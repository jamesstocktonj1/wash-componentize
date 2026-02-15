package export_wasi_http_incoming_handler

import (
	httptypes "componentize-basic/gen/wasi_http_types"
	"net/http"

	wittypes "github.com/bytecodealliance/wit-bindgen/wit_types"
)

var handler http.HandlerFunc

// Handle the specified `Request`, returning a `Response`
func Handle(request *httptypes.IncomingRequest, responseOut *httptypes.ResponseOutparam) {
	req, err := newRequest(request)
	if err != nil {
		Err := httptypes.MakeErrorCodeInternalError(wittypes.Some(err.Error()))
		result := wittypes.Err[*httptypes.OutgoingResponse](Err)
		httptypes.ResponseOutparamSet(responseOut, result)

		return
	}
	res := newResponseWriter(responseOut)

	handler(res, req)
}

func HandleFunc(h http.HandlerFunc) {
	handler = h
}
