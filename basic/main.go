package main

import (
	wasihttp "componentize-basic/gen/export_wasi_http_incoming_handler"
	_ "componentize-basic/gen/wit_exports"

	"fmt"
	"net/http"
)

func init() {
	wasihttp.HandleFunc(greeting)
}

func greeting(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

func main() {}
