package auth

import (
	"encoding/json"
	"log"
	"messenger/internal/cookies"
	"net/http"
	"time"
)

func RegisterHandler(auth *Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := auth.Register(r.Context(), req.Username, req.Password); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusCreated)
	})
}

func LoginHandler(auth *Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		sessionId, err := auth.Login(r.Context(), req.Username, req.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     cookies.SessionIDCookie,
			Value:    sessionId.String(),
			Path:     "/",
			HttpOnly: true,
			Expires:  time.Now().Add(sessionTime),
		})
		log.Printf("Saved %v for cookie %q\n", sessionId, cookies.SessionIDCookie)

		w.WriteHeader(http.StatusOK)
	})
}
