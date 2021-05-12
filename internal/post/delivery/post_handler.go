package delivery

import (
	"encoding/json"
	"net/http"

	"github.com/innovember/real-time-forum/internal/consts"
	"github.com/innovember/real-time-forum/internal/models"
	"github.com/innovember/real-time-forum/internal/mwares"
	"github.com/innovember/real-time-forum/internal/post"
	"github.com/innovember/real-time-forum/internal/user"
	"github.com/innovember/real-time-forum/pkg/response"
)

type PostHandler struct {
	postUcase post.PostUsecase
	userUcase user.UserUsecase
}

func NewPostHandler(postUcase post.PostUsecase, userUcase user.UserUsecase) *PostHandler {
	return &PostHandler{
		postUcase: postUcase,
		userUcase: userUcase,
	}
}

func (ph *PostHandler) Configure(mux *http.ServeMux, mm *mwares.MiddlewareManager) {
	mux.HandleFunc("/api/v1/post", mm.CORSConfig(mm.CheckCSRF(mm.CheckAuth(ph.HandlerCreatePost))))
}

func (ph *PostHandler) HandlerCreatePost(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var (
			input  models.InputPost
			err    error
			userID int64
		)
		if err = json.NewDecoder(r.Body).Decode(&input); err != nil {
			response.JSON(w, false, http.StatusBadRequest, err.Error(), nil)
			return
		}
		if r.Context().Value(consts.ConstAuthedUserParam) != nil {
			userID = r.Context().Value(consts.ConstAuthedUserParam).(int64)
		}
		user, err := ph.userUcase.GetByID(userID)
		if err != nil {
			switch err {
			case consts.ErrNoData:
				response.JSON(w, false, http.StatusUnauthorized, consts.ErrUserNotExist.Error(), nil)
				return
			default:
				response.JSON(w, false, http.StatusInternalServerError, err.Error(), nil)
				return
			}
		}
		post := models.Post{
			AuthorID: user.ID,
			Title:    input.Title,
			Content:  input.Content,
		}
		err = ph.postUcase.Create(&post, input.Categories)
		if err != nil {
			response.JSON(w, true, http.StatusInternalServerError, err.Error(), nil)
			return
		}
		response.JSON(w, true, http.StatusCreated, consts.PostCreated, nil)
		return
	default:
		response.JSON(w, false, http.StatusMethodNotAllowed, consts.ErrOnlyPOST.Error(), nil)
		return
	}
}
