package export_wasi_http_incoming_handler

import (
	httptypes "componentize-basic/gen/wasi_http_types"
	"net/http"
)

func newRequest(request *httptypes.IncomingRequest) (*http.Request, error) {
	return nil, nil
}
