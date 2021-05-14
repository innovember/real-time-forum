package comment

import "github.com/innovember/real-time-forum/internal/models"

type CommentRepository interface {
	Insert(comment *models.Comment) (err error)
	SelectCommentsByPostID(*models.InputGetComments) (comments []models.Comment, err error)
	SelectCommentsByAuthorID(*models.InputGetComments) (comments []models.Comment, err error)
	SelectCommentByID(commentID int64) (comment *models.Comment, err error)
	SelectCommentsNumberByPostID(postID int64) (commentsNumber int64, err error)
}
