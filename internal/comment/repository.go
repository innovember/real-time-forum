package comment

import "github.com/innovember/real-time-forum/internal/models"

type CommentRepository interface {
	Insert(comment *models.Comment) (err error)
	SelectCommentsByPostID(postID int64) (comments []models.Comment, err error)
	SelectCommentsByAuthorID(authorID int64) (comments []models.Comment, err error)
	SelectCommentByID(commentID int64) (comment *models.Comment, err error)
	SelectCommentsNumberByPostID(postID int64) (commentsNumber int64, err error)
}
