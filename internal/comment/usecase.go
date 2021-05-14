package comment

import "github.com/innovember/real-time-forum/internal/models"

type CommentUsecase interface {
	Create(comment *models.Comment) (err error)
	GetCommentsByPostID(*models.InputGetComments) (comments []models.Comment, err error)
	GetCommentsByAuthorID(*models.InputGetComments) (comments []models.Comment, err error)
	GetCommentByID(commentID int64) (comment *models.Comment, err error)
}
