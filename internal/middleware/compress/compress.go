package compress

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

// Write запись данных
func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

// CompressMiddleware прослойка для компресии данных
func CompressMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		gzWriter := gzip.NewWriter(w)
		defer gzWriter.Close()

		w.Header().Set("Content-Encoding", "gzip")
		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gzWriter}, r)
	})
}

// DecompressMiddleware прослойка для декомпресии данных
func DecompressMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if !strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		gzipReader, err := gzip.NewReader(r.Body)

		if err != nil {
			panic(err)
		}

		defer gzipReader.Close()

		body, err := io.ReadAll(gzipReader)

		if err != nil {
			panic(err)
		}

		r.Body = io.NopCloser(bytes.NewReader(body))
		r.ContentLength = int64(len(body))

		next.ServeHTTP(w, r)
	})
}
