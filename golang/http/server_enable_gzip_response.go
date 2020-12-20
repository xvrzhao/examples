package http

import (
	"compress/gzip"
	"log"
	"net/http"
)

// RunServerEnableGzip demonstrates how a HTTP server responses with gzip encoding.
func RunServerEnableGzip(addr string) {
	serveMux.HandleFunc("/test/gzip", func(writer http.ResponseWriter, request *http.Request) {
		msg := []byte("Hello, Gzip!")
		// 0. the following operation assume that the client supports the gzip encoding, so it is necessary
		//	  to determine whether to enable gzip response by retrieving request's `accept-encoding` header.

		// 1. set header: content encoding and type, can't set one without the other
		writer.Header().Set("content-encoding", "gzip")
		writer.Header().Set("content-type", http.DetectContentType(msg))

		// 2. wrap the writer
		gzipWriter := gzip.NewWriter(writer)
		defer gzipWriter.Close() // remember to close because this will also flush any unwritten data to underlying writer

		// 3. write data by wrapper
		_, err := gzipWriter.Write(msg)
		if err != nil {
			http.Error(writer, "server error", http.StatusInternalServerError)
			log.Printf("failed to write compressed stream: %v", err)
			return
		}
	})

	err := http.ListenAndServe(addr, serveMux)
	if err != nil {
		log.Printf("failed to listen: %v", err)
		return
	}
}
