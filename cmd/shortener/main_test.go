package main

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_redirect(t *testing.T) {
	urlStore["e9db20b2"] = "https://yandex.ru"

	tests := []struct {
		name           string
		method         string
		shortUrl       string
		expectedCode   int
		expectedHeader string
	}{
		{
			name:           "Valid short URL",
			method:         http.MethodGet,
			shortUrl:       "e9db20b2",
			expectedCode:   http.StatusTemporaryRedirect,
			expectedHeader: "https://yandex.ru",
		},
		{
			name:           "Invalid short URL",
			method:         http.MethodGet,
			shortUrl:       "asddgfs",
			expectedCode:   http.StatusBadRequest,
			expectedHeader: "",
		},
		{
			name:           "Invalid method",
			method:         http.MethodPost,
			shortUrl:       "asdsfffas",
			expectedCode:   http.StatusBadRequest,
			expectedHeader: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "localhost:8080"+tt.shortUrl, nil)
			resRec := httptest.NewRecorder()

			redirect(resRec, req)

			if resRec.Code != tt.expectedCode {
				t.Errorf("expected status %d, got %d", tt.expectedCode, resRec.Code)
			}

			locationHeader := resRec.Header().Get("Location")
			if locationHeader != tt.expectedHeader {
				t.Errorf("expected header %v got %v", tt.expectedHeader, locationHeader)
			}
		})
	}
}

func Test_shortenUrl(t *testing.T) {
	tests := []struct {
		name         string
		method       string
		url          string
		expectedCode int
		expectedBody string
	}{
		{
			name:         "Simple POST request",
			method:       http.MethodPost,
			url:          "https://yandex.ru",
			expectedCode: http.StatusCreated,
			expectedBody: "http://localhost:8080" + fmt.Sprintf("%x", md5.Sum([]byte("https://yandex.ru")))[:8],
		},
		{
			name:         "Invalid request",
			method:       http.MethodGet,
			url:          "",
			expectedCode: http.StatusBadRequest,
			expectedBody: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req *http.Request
			if tt.method == http.MethodPost {
				req = httptest.NewRequest(tt.method, "http://localhost:8080/shorten", bytes.NewBufferString("url="+tt.url))
			} else {
				req = httptest.NewRequest(tt.method, "http://localhost:8080/shorten", nil)
			}

			rr := httptest.NewRecorder()

			handler := http.HandlerFunc(shortenUrl)
			handler.ServeHTTP(rr, req)

			if rr.Code != tt.expectedCode {
				t.Errorf("expected status %d, got %d", tt.expectedCode, rr.Code)
			}
		})
	}
}
