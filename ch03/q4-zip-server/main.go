package main

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/zip")                              // 必須
	w.Header().Set("Content-Disposition", "attachment; filename=ascii_sample.zip") // オプション

	zipWriter := zip.NewWriter(w)
	defer zipWriter.Close()
	file, err := zipWriter.Create(strings.Repeat("a", 70000))
	if err != nil {
		panic(err)
	}
	io.Copy(file, strings.NewReader("あああ\naaa"))
}

func main() {
	fmt.Println("server start")
	http.HandleFunc("/", handler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
