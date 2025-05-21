package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Hordevcom/URLShortener/internal/middleware/jwtgen"
)

// ShortenOrigURLs структура для серриализации JSON
type ShortenOrigURLs struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

// GetUserUrls по этому хендлеру получаем все урлы конкретного пользователя
func (h *ShortenHandler) GetUserUrls(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	var UserID int
	var ShorigURLs []ShortenOrigURLs

	if err != nil {
		token, _ := jwtgen.BuildJWTString()
		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    token,
			HttpOnly: true,
		})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
		fmt.Println("Token created at GetUserUrls!")
		return
	}
	fmt.Println(cookie.Value)
	if err := cookie.Valid(); err == nil {
		UserID = jwtgen.GetUserID(cookie.Value)
		fmt.Println(UserID)
		fmt.Println("UserID collected from cookie.Value")
	}

	URLs, ok := h.DB.GetWithUserID(r.Context(), UserID)

	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	for key, value := range URLs {
		ShorigURLs = append(ShorigURLs, ShortenOrigURLs{
			ShortURL:    h.Config.Host + "/" + key,
			OriginalURL: value,
		})
	}
	var ShorigURLs1 []ShortenOrigURLs
	ShorigURLs1 = append(ShorigURLs1, ShorigURLs[len(ShorigURLs)-2])
	fmt.Println(ShorigURLs)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(ShorigURLs1)
	if err != nil {
		panic(err)
	}

}
