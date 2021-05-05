package mwares

import (
	"net/http"

	"github.com/innovember/real-time-forum/internal/user"
)

type MiddlewareManager struct {
	userUsecase user.UserUsecase
}

func NewMiddlewareManager(userUsecase user.UserUsecase) *MiddlewareManager {
	return &MiddlewareManager{
		userUsecase: userUsecase,
	}
}

func (mm *MiddlewareManager) CORSConfig(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if r.Method == "OPTIONS" {
			w.WriteHeader(200)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		next(w, r)
	}
}
