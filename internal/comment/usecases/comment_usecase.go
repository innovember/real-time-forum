package usecases

import (
	"github.com/innovember/real-time-forum/internal/comment"
	"github.com/innovember/real-time-forum/internal/models"
)

type CommentUsecase struct {
	commentRepo comment.CommentRepository
}

func NewCommentUsecase(commentRepo comment.CommentRepository) *CommentUsecase {
	return &CommentUsecase{
		commentRepo: commentRepo,
	}
}

func (cu *CommentUsecase) Create(comment *models.Comment) (err error) {
	if err = cu.commentRepo.Insert(comment); err != nil {
		return err
	}
	return err
}
func (cu *CommentUsecase) GetCommentsByPostID(postID int64) (comments []models.Comment, err error) {
	if comments, err = cu.commentRepo.SelectCommentsByPostID(postID); err != nil {
		return nil, err
	}
	return comments, err
}

func (cu *CommentUsecase) GetCommentsByAuthorID(authorID int64) (comments []models.Comment, err error) {
	if comments, err = cu.commentRepo.SelectCommentsByAuthorID(authorID); err != nil {
		return nil, err
	}
	return comments, err
}

func (cu *CommentUsecase) GetCommentByID(commentID int64) (comment *models.Comment, err error) {
	if comment, err = cu.commentRepo.SelectCommentByID(commentID); err != nil {
		return nil, err
	}
	return comment, err
}
