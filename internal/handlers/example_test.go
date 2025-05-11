package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/Hordevcom/URLShortener/internal/config"
	"github.com/Hordevcom/URLShortener/internal/storage"
	"github.com/go-chi/chi/v5"
)

func ExampleShortenHandler_ShortenURL() {
	// Подготовка к выполнению запроса
	m1 := storage.NewMapStorage()
	m1.Set(context.Background(), "abc123", "https://example.com", 0)
	conf := config.NewConfig()
	app := &ShortenHandler{
		Storage: m1,
		Config:  conf,
	}
	// инициализируем body запроса
	originalURL := "https://example.com"
	reqBody := []byte(originalURL)

	req := httptest.NewRequest("POST", "/shorten", bytes.NewReader(reqBody))
	// Добавляем куки
	req.AddCookie(&http.Cookie{
		Name:  "token",
		Value: "2",
	})
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(app.ShortenURL)

	// Вызываем хендлер
	handler.ServeHTTP(rr, req)
	resp := rr.Result()

	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	fmt.Println(resp.StatusCode)
	fmt.Println(string(body))

	// Output:
	// 201
	// http://localhost:8080/c984d06a
}

func ExampleShortenHandler_ShortenURLJSON() {
	// Подготовка к выполнению запроса
	m1 := storage.NewMapStorage()
	m1.Set(context.Background(), "abc123", "https://example.com", 0)
	app := &ShortenHandler{
		Storage:     storage.NewMapStorage(),
		JSONStorage: *storage.NewJSONStorage(),
		Config:      config.Config{Host: "http://localhost"},
	}
	// инициализируем body запроса
	requestBody := `{"url":"https://example.com"}`

	// Подготовка запроса POST /api/shorten
	req := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewBufferString(requestBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.ShortenURLJSON)

	// Вызываем хендлер
	handler.ServeHTTP(rr, req)
	resp := rr.Result()

	var res Response
	json.NewDecoder(resp.Body).Decode(&res)

	// Выводим код ответа и Body ответа
	fmt.Println(resp.StatusCode)
	fmt.Println(res.Result)

	// Output:
	// 201
	// http://localhost/c984d06a
}

func ExampleShortenHandler_Redirect() {
	// Подготовка к выполнению запроса
	m1 := storage.NewMapStorage()
	m1.Set(context.Background(), "c984d06a", "https://example.com", 0)
	app := &ShortenHandler{Storage: m1}

	// Используем chi роутер, чтобы брать id запроса
	r := chi.NewRouter()
	r.Get("/{id}", app.Redirect)

	// Подготовка запроса Get /c984d06a
	req := httptest.NewRequest(http.MethodGet, "/c984d06a", nil)
	rr := httptest.NewRecorder()

	// Вызываем хендлер
	r.ServeHTTP(rr, req)

	// Выводим код ответа и значение заголовка "Location"
	fmt.Println(rr.Code)
	fmt.Println(rr.Header().Get("Location"))

	// Output:
	// 307
	// https://example.com
}
