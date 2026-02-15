package export_wasi_http_incoming_handler

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	httptypes "redis-echo/gen/wasi_http_types"
	streams "redis-echo/gen/wasi_io_streams"
	"runtime"
	"sync"
)

func newRequest(request *httptypes.IncomingRequest) (*http.Request, error) {
	method := parseMethod(request.Method())
	if method == "" {
		return nil, errors.New("unknown method type")
	}

	authority := "localhost"
	if request.Authority().IsSome() {
		authority = request.Authority().Some()
	}

	path := "/"
	if request.PathWithQuery().IsSome() {
		path = request.PathWithQuery().Some()
	}

	body, err := newRequestBody(request)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("http://%s%s", authority, path)
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	// TODO: copy headers
	headers := request.Headers()
	defer headers.Drop()

	req.Host = authority
	req.URL.Host = authority
	req.RequestURI = path

	return req, nil
}

func parseMethod(method httptypes.Method) string {
	switch method.Tag() {
	case httptypes.MethodGet:
		return http.MethodGet
	case httptypes.MethodHead:
		return http.MethodHead
	case httptypes.MethodPost:
		return http.MethodPost
	case httptypes.MethodPut:
		return http.MethodPut
	case httptypes.MethodDelete:
		return http.MethodDelete
	case httptypes.MethodConnect:
		return http.MethodConnect
	case httptypes.MethodOptions:
		return http.MethodOptions
	case httptypes.MethodTrace:
		return http.MethodTrace
	case httptypes.MethodPatch:
		return http.MethodPatch
	case httptypes.MethodOther:
		method.Other()
		return method.Other()
	default:
		return ""
	}
}

func newRequestBody(request *httptypes.IncomingRequest) (io.ReadCloser, error) {
	consumeRes := request.Consume()
	if consumeRes.IsErr() {
		return nil, fmt.Errorf("failed to consume incoming request: %+v", consumeRes.Err())
	}
	body := consumeRes.Ok()

	streamRes := body.Stream()
	if streamRes.IsErr() {
		return nil, fmt.Errorf("failed to open request body stream: %+v", streamRes.Err())
	}
	stream := streamRes.Ok()

	return &requestBody{
		body:   body,
		stream: stream,
	}, nil
}

type requestBody struct {
	body   *httptypes.IncomingBody
	stream *streams.InputStream

	flushOnce sync.Once
}

var _ io.ReadCloser = (*requestBody)(nil)

func (r *requestBody) Read(p []byte) (int, error) {
	pollable := r.stream.Subscribe()
	for !pollable.Ready() {
		runtime.Gosched()
	}
	defer pollable.Drop()

	readRes := r.stream.Read(uint64(len(p)))
	if readRes.IsErr() {
		switch readRes.Err().Tag() {
		case streams.StreamErrorClosed:
			r.Close()
			return 0, io.EOF
		case streams.StreamErrorLastOperationFailed:
			return 0, errors.New(readRes.Err().LastOperationFailed().ToDebugString())
		default:
			return 0, fmt.Errorf("unknown error reading from stream: %+v", readRes.Err())
		}
	}

	data := readRes.Ok()
	copy(p, data)
	return len(data), nil
}

func (r *requestBody) Close() error {
	r.flushOnce.Do(r.flush)

	if r.stream != nil {
		r.stream.Drop()
		r.stream = nil
	}

	if r.body != nil {
		r.body.Drop()
		r.body = nil
	}

	return nil
}

func (r *requestBody) flush() {
	futureTrailers := httptypes.IncomingBodyFinish(r.body)
	defer futureTrailers.Drop()

	futureTrailers.Get()
	r.body = nil
}
