package cookie

import (
	"net/http"

	"github.com/innovember/real-time-forum/internal/models"
)

func BuildCookie(session *models.Session) *http.Cookie {
	return &http.Cookie{
		Value:    session.Token,
		Name:     session.Name,
		Expires:  session.Expires,
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
	}
}
