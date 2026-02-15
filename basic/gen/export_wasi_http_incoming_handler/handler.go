package export_wasi_http_incoming_handler

import (
	. "componentize-basic/gen/wasi_http_types"
	"net/http"

	. "github.com/bytecodealliance/wit-bindgen/wit_types"
)

// Handle the specified `Request`, returning a `Response`
func Handle(request *IncomingRequest, responseOut *ResponseOutparam) {
	response := MakeOutgoingResponse(MakeFields())

	body := response.Body()

	ResponseOutparamSet(responseOut, Ok[*OutgoingResponse, ErrorCode](response))

	message := []byte("Hello, world!")

	if body.IsOk() {
		stream := body.Ok()
		if writeResult := stream.Write(); writeResult.IsOk() {
			writeResult.Ok().BlockingWriteAndFlush(message)
		}
	}
}

func HandleFunc(h http.HandlerFunc) {

}
