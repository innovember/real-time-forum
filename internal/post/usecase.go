package post

import "github.com/innovember/real-time-forum/internal/models"

type PostUsecase interface {
	Create(*models.Post, []string) error
	GetAllPosts(*models.InputGetPosts) (posts []models.Post, err error)
	GetPostByID(postID int64) (post *models.Post, err error)
	GetAllPostsByAuthorID(*models.InputGetPosts) (posts []models.Post, err error)
	GetAllPostsByCategories(*models.InputGetPosts) (posts []models.Post, err error)
}
