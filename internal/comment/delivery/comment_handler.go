package delivery

import (
	"encoding/json"
	"net/http"

	"github.com/innovember/real-time-forum/internal/comment"
	"github.com/innovember/real-time-forum/internal/consts"
	"github.com/innovember/real-time-forum/internal/models"
	"github.com/innovember/real-time-forum/internal/mwares"
	"github.com/innovember/real-time-forum/internal/post"
	"github.com/innovember/real-time-forum/internal/user"
	"github.com/innovember/real-time-forum/pkg/response"
)

type CommentHandler struct {
	userUcase    user.UserUsecase
	postUcase    post.PostUsecase
	commentUcase comment.CommentUsecase
}

func NewCommentHandler(
	userUcase user.UserUsecase,
	postUcase post.PostUsecase,
	commentUcase comment.CommentUsecase,
) *CommentHandler {
	return &CommentHandler{
		userUcase:    userUcase,
		postUcase:    postUcase,
		commentUcase: commentUcase,
	}
}

func (ch *CommentHandler) Configure(mux *http.ServeMux, mm *mwares.MiddlewareManager) {
	mux.HandleFunc("/api/v1/comment", mm.CORSConfig(mm.CheckCSRF(mm.CheckAuth(ch.HandlerCreateComment))))
	mux.HandleFunc("/api/v1/comments", mm.CORSConfig(ch.HandlertGetComments))
}

func (ch *CommentHandler) HandlerCreateComment(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var (
			input  models.InputComment
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
		user, err := ch.userUcase.GetByID(userID)
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
		_, err = ch.postUcase.GetPostByID(int64(input.PostID))
		if err != nil {
			switch err {
			case consts.ErrNoData:
				response.JSON(w, false, http.StatusNotFound, consts.ErrPostNotExist.Error(), nil)
				return
			default:
				response.JSON(w, false, http.StatusInternalServerError, err.Error(), nil)
				return
			}
		}
		comment := models.Comment{
			AuthorID: user.ID,
			PostID:   input.PostID,
			Content:  input.Content,
		}
		err = ch.commentUcase.Create(&comment)
		if err != nil {
			response.JSON(w, false, http.StatusInternalServerError, err.Error(), nil)
			return
		}
		response.JSON(w, true, http.StatusCreated, consts.CommentCreated, nil)
		return
	default:
		response.JSON(w, false, http.StatusMethodNotAllowed, consts.ErrOnlyPOST.Error(), nil)
		return
	}
}

func (ch *CommentHandler) HandlertGetComments(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var (
			input    models.InputGetComments
			err      error
			comments []models.Comment
		)
		if err = json.NewDecoder(r.Body).Decode(&input); err != nil {
			response.JSON(w, false, http.StatusBadRequest, err.Error(), nil)
			return
		}
		switch input.Option {
		case "post":
			_, err := ch.postUcase.GetPostByID(int64(input.PostID))
			if err != nil {
				switch err {
				case consts.ErrNoData:
					response.JSON(w, true, http.StatusNotFound, consts.ErrPostNotExist.Error(), nil)
					return
				default:
					response.JSON(w, true, http.StatusInternalServerError, err.Error(), nil)
					return
				}
			}
			comments, err = ch.commentUcase.GetCommentsByPostID(&input)
			if err != nil {
				response.JSON(w, false, http.StatusInternalServerError, err.Error(), nil)
				return
			}
		case "user":
			_, err := ch.userUcase.GetByID(input.UserID)
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
			comments, err = ch.commentUcase.GetCommentsByAuthorID(&input)
			if err != nil {
				response.JSON(w, false, http.StatusInternalServerError, err.Error(), nil)
				return
			}
		}
		response.JSON(w, true, http.StatusOK, consts.Comments, comments)
		return
	default:
		response.JSON(w, false, http.StatusMethodNotAllowed, consts.ErrOnlyPOST.Error(), nil)
		return
	}
}
