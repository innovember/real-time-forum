package post

import "github.com/innovember/real-time-forum/internal/models"

type PostUsecase interface {
	Create(*models.Post, []string) error
	GetAllPosts() (posts []models.Post, err error)
	GetPostByID(postID int64) (post *models.Post, err error)
	GetAllPostsByAuthorID(authorID int64) (posts []models.Post, err error)
	GetAllPostsByCategories(categories []string) (posts []models.Post, err error)
}
