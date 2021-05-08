package mwares

import (
	"context"
	"log"
	"net/http"

	"github.com/innovember/real-time-forum/internal/consts"
	"github.com/innovember/real-time-forum/internal/session"
	"github.com/innovember/real-time-forum/internal/tools/csrf"
	"github.com/innovember/real-time-forum/internal/user"
	"github.com/innovember/real-time-forum/pkg/response"
)

type MiddlewareManager struct {
	origins      []string
	userUcase    user.UserUsecase
	sessionUcase session.SessionUsecase
}

func NewMiddlewareManager(userUcase user.UserUsecase,
	sessionUcase session.SessionUsecase) *MiddlewareManager {
	return &MiddlewareManager{
		origins:      []string{"http://localhost:3000", "localhost", "http://localhost:8081"},
		userUcase:    userUcase,
		sessionUcase: sessionUcase,
	}
}

func (mm *MiddlewareManager) isAllowedOrigin(origin string) bool {
	for _, allowed := range mm.origins {
		if string(allowed) == origin {
			return true
		}
	}
	return false
}

func (mm *MiddlewareManager) CORSConfig(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		origin := w.Header().Get("Origin")
		allowOrigin := ""
		log.Println(origin)
		if mm.isAllowedOrigin(origin) {
			allowOrigin = origin
		}
		w.Header().Set("Access-Control-Allow-Origin", allowOrigin)
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, X-Csrf-Token")
		w.Header().Set("Access-Control-Expose-Headers", "X-Csrf-Token")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		next(w, r)
	}
}

func (mm *MiddlewareManager) CheckAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			cookie *http.Cookie
			err    error
		)
		if cookie, err = r.Cookie(consts.SessionName); err != nil {
			response.JSON(w, false, http.StatusUnauthorized, consts.ErrSessionTokenNotFound.Error(), nil)
			return
		}
		session, err := mm.sessionUcase.GetByToken(cookie.Value)
		if err != nil {
			response.JSON(w, false, http.StatusUnauthorized, consts.ErrInvalidSessionToken.Error(), nil)
			return
		}

		user, err := mm.userUcase.GetByID(session.UserID)
		if err != nil {
			response.JSON(w, false, http.StatusUnauthorized, consts.ErrUserNotExist.Error(), nil)
			return
		}
		ctx := context.WithValue(r.Context(), consts.ConstAuthedUserParam, user.ID)
		next(w, r.WithContext(ctx))
	}
}

func (mm *MiddlewareManager) CheckCSRF(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			cookie *http.Cookie
			err    error
		)
		if cookie, err = r.Cookie(consts.SessionName); err != nil {
			response.JSON(w, false, http.StatusUnauthorized, consts.ErrSessionTokenNotFound.Error(), nil)
			return
		}
		session, err := mm.sessionUcase.GetByToken(cookie.Value)
		if err != nil {
			response.JSON(w, false, http.StatusUnauthorized, consts.ErrInvalidSessionToken.Error(), nil)
			return
		}
		csrfToken := r.Header.Get("X-Csrf-Token")
		if csrfToken == "" {
			response.JSON(w, false, http.StatusBadRequest, consts.ErrCSRF.Error(), nil)
			return
		}
		err = csrf.ValidateCSRFToken(session, csrfToken)
		if err != nil {
			response.JSON(w, false, http.StatusBadRequest, consts.ErrCSRF.Error(), nil)
			return
		}
		next(w, r)
	}
}

func (mm *MiddlewareManager) XSS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		next(w, r)
	}
}
