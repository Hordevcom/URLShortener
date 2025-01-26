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
		name         string
		method       string
		shortUrl     string
		expectedCode int
		expectedBody string
	}{
		{
			name:         "Valid short URL",
			method:       http.MethodGet,
			shortUrl:     "e9db20b2",
			expectedCode: http.StatusTemporaryRedirect,
			expectedBody: "Location: https://yandex.ru",
		},
		{
			name:         "Invalid short URL",
			method:       http.MethodGet,
			shortUrl:     "asddgfs",
			expectedCode: http.StatusBadRequest,
			expectedBody: "",
		},
		{
			name:         "Invalid method",
			method:       http.MethodPost,
			shortUrl:     "asdsfffas",
			expectedCode: http.StatusBadRequest,
			expectedBody: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "http://localhost:8080/"+tt.shortUrl, nil)
			resRec := httptest.NewRecorder()

			handler := http.HandlerFunc(redirect)
			handler.ServeHTTP(resRec, req)

			if resRec.Code != tt.expectedCode {
				t.Errorf("expected status %d, got %d", tt.expectedCode, resRec.Code)
			}

			if resRec.Body.String() != tt.expectedBody {
				t.Errorf("expected body %v got %v", tt.expectedBody, resRec.Body.String())
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
