package main

import (
	"io"
	"net/http"
)

func handler(w http.ResponseWriter, h *http.Request) {
	io.WriteString(w, "http.ResponseWriter sample")
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
