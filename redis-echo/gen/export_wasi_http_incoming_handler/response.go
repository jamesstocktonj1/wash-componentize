package export_wasi_http_incoming_handler

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	httptypes "redis-echo/gen/wasi_http_types"
	streams "redis-echo/gen/wasi_io_streams"
	"sync"

	wittypes "github.com/bytecodealliance/wit-bindgen/wit_types"
)

func newResponseWriter(resp *httptypes.ResponseOutparam) http.ResponseWriter {
	return &responseWriter{
		outparam:   resp,
		header:     make(http.Header),
		statusCode: http.StatusOK,
	}
}

type responseWriter struct {
	outparam *httptypes.ResponseOutparam
	response *httptypes.OutgoingResponse

	body       *httptypes.OutgoingBody
	stream     *streams.OutputStream
	wasiHeader *httptypes.Headers

	header     http.Header
	headerOnce sync.Once
	err        error

	statusCode int
}

var _ http.ResponseWriter = (*responseWriter)(nil)

func (r *responseWriter) Header() http.Header {
	return r.header
}

func (r *responseWriter) Write(data []byte) (int, error) {
	r.headerOnce.Do(r.flush)
	if r.err != nil {
		return 0, r.err
	}

	writeRes := r.stream.Write(data)
	if writeRes.IsErr() {
		switch writeRes.Err().Tag() {
		case streams.StreamErrorClosed:
			return 0, io.EOF
		case streams.StreamErrorLastOperationFailed:
			return 0, errors.New(writeRes.Err().LastOperationFailed().ToDebugString())
		default:
			return 0, fmt.Errorf("unknown error writing to stream: %+v", writeRes.Err())
		}
	}

	r.stream.BlockingFlush()
	return len(data), nil
}

func (r *responseWriter) WriteHeader(statusCode int) {
	r.headerOnce.Do(func() {
		r.statusCode = statusCode
		r.flush()
	})
}

func (r *responseWriter) Close() error {
	return nil
}

func (r *responseWriter) flushHeaders() error {
	// for key, value := range r.header {
	// 	r.wasiHeader.Set()
	// 	values := []uint8{}
	// }

	r.wasiHeader = httptypes.MakeFields()
	return nil
}

func (r *responseWriter) flush() {
	if err := r.flushHeaders(); err != nil {
		r.err = err
		return
	}

	r.response = httptypes.MakeOutgoingResponse(r.wasiHeader)
	r.response.SetStatusCode(uint16(r.statusCode))

	bodyRes := r.response.Body()
	if bodyRes.IsErr() {
		r.err = fmt.Errorf("failed to open response body: %+v", bodyRes.Err())
		return
	}
	r.body = bodyRes.Ok()

	writeRes := r.body.Write()
	if writeRes.IsErr() {
		r.err = fmt.Errorf("failed to open response body stream: %+v", writeRes.Err())
		return
	}
	r.stream = writeRes.Ok()

	result := wittypes.Ok[*httptypes.OutgoingResponse, httptypes.ErrorCode](r.response)
	httptypes.ResponseOutparamSet(r.outparam, result)
}
