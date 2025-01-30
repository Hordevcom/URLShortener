package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRedirect(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	urlStore = make(map[string]string)

	r.GET("/:id", redirect)

	urlStore["abcdef12"] = "https://example.com"

	tests := []struct {
		name           string
		param          string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "valid short URL",
			param:          "abcdef12",
			expectedStatus: http.StatusTemporaryRedirect,
			expectedBody:   "{}",
		},
		{
			name:           "invalid short URL",
			param:          "nonexistent",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "{}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/"+tt.param, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Contains(t, w.Body.String(), tt.expectedBody)
		})
	}
}
