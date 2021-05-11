package comment

type CommentRepository interface {
	// Create(userID int64, comment *models.Comment) (newComment *models.Comment, err error)
	// SelectCommentsByPostID(userID, postID int64) (comments []models.Comment, err error)
	// SelectAuthor(comment *models.Comment) (err error)
	// SelectCommentsByAuthorID(userID, authorID int64) (comments []models.Comment, err error)
	SelectCommentsNumberByPostID(postID int64) (commentsNumber int, err error)
	// SelectCommentByID(userID, commentID int64) (comment *models.Comment, err error)
}
