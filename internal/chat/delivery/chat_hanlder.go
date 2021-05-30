package delivery

import (
	"encoding/json"
	"net/http"

	"github.com/innovember/real-time-forum/internal/chat"
	"github.com/innovember/real-time-forum/internal/consts"
	"github.com/innovember/real-time-forum/internal/models"
	"github.com/innovember/real-time-forum/internal/mwares"
	"github.com/innovember/real-time-forum/internal/session"
	"github.com/innovember/real-time-forum/pkg/response"
)

type ChatHandler struct {
	roomUsecase  chat.RoomUsecase
	sessionUcase session.SessionUsecase
}

func NewChatHandler(roomUsecase chat.RoomUsecase,
	sessionUcase session.SessionUsecase) *ChatHandler {
	return &ChatHandler{
		roomUsecase:  roomUsecase,
		sessionUcase: sessionUcase,
	}
}

func (ch *ChatHandler) Configure(mux *http.ServeMux, mm *mwares.MiddlewareManager) {
	mux.HandleFunc("/api/v1/chats", mm.CORSConfig(mm.CheckAuth(ch.HandlerGetChats)))
	mux.HandleFunc("/api/v1/room", mm.CORSConfig(mm.CheckCSRF(mm.CheckAuth(ch.HandlerGetRoom))))
	mux.HandleFunc("/api/v1/messages", mm.CORSConfig(mm.CheckAuth(ch.HandlerGetMessages)))
	mux.HandleFunc("/api/v1/message", mm.CORSConfig(mm.CheckCSRF(mm.CheckAuth(ch.HandlerWsSendMessage))))
}

func (ch *ChatHandler) HandlerGetChats(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		cookie, _ := r.Cookie(consts.SessionName)
		session, err := ch.sessionUcase.GetByToken(cookie.Value)
		if err != nil {
			response.JSON(w, false, http.StatusUnauthorized, consts.ErrInvalidSessionToken.Error(), nil)
			return
		}
		chats, err := ch.roomUsecase.GetAllRoomsByUserID(session.UserID)
		if err != nil {
			response.JSON(w, false, http.StatusInternalServerError, err.Error(), nil)
			return
		}
		response.JSON(w, true, http.StatusOK, consts.AllChats, chats)
		return
	default:
		response.JSON(w, false, http.StatusMethodNotAllowed, consts.ErrOnlyGet.Error(), nil)
		return
	}
}

func (ch *ChatHandler) HandlerGetRoom(w http.ResponseWriter, r *http.Request) {

}

func (ch *ChatHandler) HandlerGetMessages(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var (
			input models.InputRoom
		)
		cookie, _ := r.Cookie(consts.SessionName)
		_, err := ch.sessionUcase.GetByToken(cookie.Value)
		if err != nil {
			response.JSON(w, false, http.StatusUnauthorized, consts.ErrInvalidSessionToken.Error(), nil)
			return
		}
		if err = json.NewDecoder(r.Body).Decode(&input); err != nil {
			response.JSON(w, false, http.StatusBadRequest, err.Error(), nil)
			return
		}
		messages, err := ch.roomUsecase.GetMessages(input.RoomID, input.LastMessageID)
		if err != nil {
			response.JSON(w, false, http.StatusBadRequest, err.Error(), nil)
			return
		}
		response.JSON(w, true, http.StatusOK, consts.RoomMessages, messages)
		return
	default:
		response.JSON(w, false, http.StatusMethodNotAllowed, consts.ErrOnlyPOST.Error(), nil)
		return
	}
}

func (ch *ChatHandler) HandlerWsSendMessage(w http.ResponseWriter, r *http.Request) {

}
