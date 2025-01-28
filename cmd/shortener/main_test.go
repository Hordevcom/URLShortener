package main

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type Data struct {
	Url string `json:"url"`
}

func TestShortenUrl(t *testing.T) {
	// Создаем новый роутер Gin
	r := gin.Default()
	r.POST("/shorten", shortenUrl)

	// Создаем тестовые кейсы
	tests := []struct {
		name         string
		formData     map[string]string
		expectedCode int
		expectedBody string
	}{
		{
			name:         "valid URL",
			formData:     map[string]string{"url": "http://example.com"},
			expectedCode: http.StatusCreated,
			expectedBody: `{"message":"http://localhost:8080/` + fmt.Sprintf("%x", md5.Sum([]byte("http://example.com")))[:8] + `"}`,
		},
		{
			name:         "empty URL",
			formData:     map[string]string{"url": ""},
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"message":"url param required"}`,
		},
		{
			name:         "missing URL param",
			formData:     map[string]string{},
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"message":"url param required"}`,
		},
	}

	// Пробегаем все тестовые кейсы
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Формируем данные формы как строку
			form := ""
			for key, value := range tt.formData {
				form += fmt.Sprintf("%s=%s", key, value)
			}

			// Создаем новый запрос с телом формы
			req, err := http.NewRequest("POST", "/shorten", strings.NewReader(form))
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			// Создаем ResponseRecorder для захвата ответа
			w := httptest.NewRecorder()

			// Выполняем запрос
			r.ServeHTTP(w, req)

			// Проверяем код ответа
			assert.Equal(t, tt.expectedCode, w.Code)

			// Проверяем тело ответа
			assert.JSONEq(t, tt.expectedBody, w.Body.String())
		})
	}
}
