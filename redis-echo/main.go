package main

import (
	"fmt"
	wasihttp "redis-echo/gen/export_wasi_http_incoming_handler"
	store "redis-echo/gen/wasi_keyvalue_store"
	_ "redis-echo/gen/wit_exports"

	"net/http"
)

func init() {
	wasihttp.HandleFunc(greeting)
}

func greeting(w http.ResponseWriter, r *http.Request) {
	openRes := store.Open("new")
	if openRes.IsErr() {
		http.Error(w, fmt.Sprintf("failed to open keyvalue store - %+v", openRes.Err()), http.StatusInternalServerError)
		return
	}
	bucket := openRes.Ok()

	setRes := bucket.Set("name", []uint8("World"))
	if setRes.IsErr() {
		http.Error(w, fmt.Sprintf("failed to set key - %+v", setRes.Err()), http.StatusInternalServerError)
		return
	}

	getRes := bucket.Get("name")
	if getRes.IsErr() {
		http.Error(w, fmt.Sprintf("failed to get key - %+v", getRes.Err()), http.StatusInternalServerError)
		return
	}
	value := getRes.Ok()
	if value.IsNone() {
		http.Error(w, "key not found", http.StatusNotFound)
		return
	}

	fmt.Fprintf(w, "Hello, %s!", string(value.Some()))
}

func main() {}
