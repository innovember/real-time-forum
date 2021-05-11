package usecases

import "github.com/innovember/real-time-forum/internal/comment"

type CommentUsecase struct {
	commentRepo comment.CommentRepository
}

func NewCommentUsecase(commentRepo comment.CommentRepository) *CommentUsecase {
	return &CommentUsecase{
		commentRepo: commentRepo,
	}
}
