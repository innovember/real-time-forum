package delivery

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/innovember/real-time-forum/internal/consts"
	"github.com/innovember/real-time-forum/internal/models"
	"github.com/innovember/real-time-forum/internal/mwares"
	"github.com/innovember/real-time-forum/internal/session"
	"github.com/innovember/real-time-forum/internal/tools/cookie"
	"github.com/innovember/real-time-forum/internal/tools/csrf"
	"github.com/innovember/real-time-forum/internal/user"
	"github.com/innovember/real-time-forum/pkg/response"
)

type SessionHandler struct {
	SessionUcase session.SessionUsecase
	UserUcase    user.UserUsecase
}

func NewSessionHandler(sessionUcase session.SessionUsecase,
	userUcase user.UserUsecase) *SessionHandler {
	return &SessionHandler{
		UserUcase:    userUcase,
		SessionUcase: sessionUcase,
	}
}

func (sh *SessionHandler) Configure(mux *http.ServeMux, mm *mwares.MiddlewareManager) {
	mux.HandleFunc("/api/v1/session/login", mm.CORSConfig(sh.HandlerLogin))
	mux.HandleFunc("/api/v1/session/logout", mm.CORSConfig(mm.CheckCSRF(mm.CheckAuth(sh.HandlerLogout))))
}

func (sh *SessionHandler) HandlerLogin(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var (
			input models.InputUserSignIn
		)
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			response.JSON(w, false, http.StatusBadRequest, err.Error(), nil)
			return
		}
		user, err := sh.UserUcase.GetByEmailOrNickname(input.Nickname)
		if err != nil {
			switch err {
			case consts.ErrNoData:
				response.JSON(w, false, http.StatusBadRequest, consts.ErrUserNotExist.Error(), nil)
				return
			default:
				response.JSON(w, false, http.StatusInternalServerError, err.Error(), nil)
				return
			}
		}
		if err = sh.UserUcase.CheckPassword(&input); err != nil {
			response.JSON(w, false, http.StatusBadRequest, err.Error(), nil)
			return
		}
		session := models.NewSession(user.ID)
		if err = sh.SessionUcase.CreateSession(session); err != nil {
			response.JSON(w, false, http.StatusInternalServerError, err.Error(), nil)
			return
		}
		token, err := csrf.NewCSRFToken(session)
		if err != nil {
			response.JSON(w, false, http.StatusInternalServerError, err.Error(), nil)
			return
		}
		cookie := cookie.BuildCookie(session)
		w.Header().Set("Set-Cookie", cookie.String())
		w.Header().Set("X-CSRF-TOKEN", token)
		if err = sh.SessionUcase.UpdateStatus(session.UserID, consts.StatusOnline); err != nil {
			response.JSON(w, false, http.StatusInternalServerError, err.Error(), nil)
			return
		}
		response.JSON(w, true, http.StatusOK, consts.LoginSuccess, nil)
		return
	default:
		response.JSON(w, false, http.StatusMethodNotAllowed, consts.ErrOnlyPOST.Error(), nil)
		return
	}
}

func (sh *SessionHandler) HandlerLogout(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodDelete:
		cookie, _ := r.Cookie(consts.SessionName)
		session, err := sh.SessionUcase.GetByToken(cookie.Value)
		if err != nil {
			response.JSON(w, false, http.StatusUnauthorized, consts.ErrInvalidSessionToken.Error(), nil)
			return
		}
		err = sh.SessionUcase.DeleteSession(session.Token)
		if err != nil {
			response.JSON(w, false, http.StatusInternalServerError, err.Error(), nil)
			return
		}
		cookie.Path = "/"
		cookie.Expires = time.Now().AddDate(0, 0, -2)
		w.Header().Set("Set-Cookie", cookie.String())
		if err = sh.SessionUcase.UpdateStatus(session.UserID, consts.StatusOffline); err != nil {
			response.JSON(w, false, http.StatusInternalServerError, err.Error(), nil)
			return
		}
		response.JSON(w, true, http.StatusOK, consts.LogoutSuccess, nil)
		return
	default:
		response.JSON(w, false, http.StatusMethodNotAllowed, consts.ErrOnlyDelete.Error(), nil)
		return
	}
}
