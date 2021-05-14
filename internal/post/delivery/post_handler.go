package delivery

import (
	"encoding/json"
	"net/http"
	"strconv"

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
	mux.HandleFunc("/api/v1/posts/", mm.CORSConfig(ph.HandlerGetPost))
	mux.HandleFunc("/api/v1/posts", mm.CORSConfig(ph.HandlerGetPosts))
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
			response.JSON(w, false, http.StatusInternalServerError, err.Error(), nil)
			return
		}
		response.JSON(w, true, http.StatusCreated, consts.PostCreated, nil)
		return
	default:
		response.JSON(w, false, http.StatusMethodNotAllowed, consts.ErrOnlyPOST.Error(), nil)
		return
	}
}

func (ph *PostHandler) HandlerGetPost(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		_id := r.URL.Path[len("/api/v1/posts/"):]
		postID, err := strconv.Atoi(_id)
		if err != nil {
			response.JSON(w, false, http.StatusBadRequest, consts.ErrPostNotExist.Error(), nil)
			return
		}
		post, err := ph.postUcase.GetPostByID(int64(postID))
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
		response.JSON(w, true, http.StatusOK, consts.PostByIDSuccess, post)
		return
	default:
		response.JSON(w, false, http.StatusMethodNotAllowed, consts.ErrOnlyGet.Error(), nil)
		return
	}
}

func (ph *PostHandler) HandlerGetPosts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var (
			input models.InputGetPosts
			err   error
			posts []models.Post
		)
		if err = json.NewDecoder(r.Body).Decode(&input); err != nil {
			response.JSON(w, false, http.StatusBadRequest, err.Error(), nil)
			return
		}
		switch input.Option {
		case "all":
			posts, err = ph.postUcase.GetAllPosts(&input)
			if err != nil {
				response.JSON(w, false, http.StatusInternalServerError, err.Error(), nil)
				return
			}
		case "author":
			_, err := ph.userUcase.GetByID(input.AuthorID)
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
			posts, err = ph.postUcase.GetAllPostsByAuthorID(&input)
			if err != nil {
				response.JSON(w, false, http.StatusInternalServerError, err.Error(), nil)
				return
			}
		case "categories":
			posts, err = ph.postUcase.GetAllPostsByCategories(&input)
			if err != nil {
				response.JSON(w, false, http.StatusInternalServerError, err.Error(), nil)
				return
			}
		}
		response.JSON(w, true, http.StatusOK, consts.Posts, posts)
		return
	default:
		response.JSON(w, false, http.StatusMethodNotAllowed, consts.ErrOnlyPOST.Error(), nil)
		return
	}
}
