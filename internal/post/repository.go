package post

import "github.com/innovember/real-time-forum/internal/models"

type PostRepository interface {
	Insert(post *models.Post) (newPost *models.Post, err error)
	SelectAllPosts(*models.InputGetPosts) (posts []models.Post, err error)
	SelectPostByID(postID int64) (post *models.Post, err error)
	SelectCategories(post *models.Post) (err error)
	SelectPostsByCategories(*models.InputGetPosts) (posts []models.Post, err error)
	SelectAllPostsByAuthorID(*models.InputGetPosts) (posts []models.Post, err error)
}
