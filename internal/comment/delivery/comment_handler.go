package delivery

import (
	"net/http"

	"github.com/innovember/real-time-forum/internal/comment"
	"github.com/innovember/real-time-forum/internal/mwares"
	"github.com/innovember/real-time-forum/internal/post"
	"github.com/innovember/real-time-forum/internal/user"
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

func (uh *CommentHandler) Configure(mux *http.ServeMux, mm *mwares.MiddlewareManager) {

}
