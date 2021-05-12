package comment

import "github.com/innovember/real-time-forum/internal/models"

type CommentUsecase interface {
	Create(comment *models.Comment) (err error)
	GetCommentsByPostID(postID int64) (comments []models.Comment, err error)
	GetCommentsByAuthorID(authorID int64) (comments []models.Comment, err error)
	GetCommentByID(commentID int64) (comment *models.Comment, err error)
}
