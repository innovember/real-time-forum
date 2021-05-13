package usecases

import (
	"github.com/innovember/real-time-forum/internal/category"
	"github.com/innovember/real-time-forum/internal/models"
	"github.com/innovember/real-time-forum/internal/post"
)

type PostUsecase struct {
	postRepo     post.PostRepository
	categoryRepo category.CategoryRepository
}

func NewPostUsecase(postRepo post.PostRepository,
	categoryRepo category.CategoryRepository) *PostUsecase {
	return &PostUsecase{
		postRepo:     postRepo,
		categoryRepo: categoryRepo,
	}
}

func (pu *PostUsecase) Create(post *models.Post, categories []string) error {
	newPost, err := pu.postRepo.Insert(post)
	if err != nil {
		return err
	}
	if err = pu.categoryRepo.Insert(newPost.ID, categories); err != nil {
		return err
	}
	return err
}

func (pu *PostUsecase) GetAllPosts() (posts []models.Post, err error) {
	if posts, err = pu.postRepo.SelectAllPosts(); err != nil {
		return nil, err
	}
	return posts, nil
}

func (pu *PostUsecase) GetPostByID(postID int64) (post *models.Post, err error) {
	if post, err = pu.postRepo.SelectPostByID(postID); err != nil {
		return nil, err
	}
	return post, nil
}

func (pu *PostUsecase) GetAllPostsByAuthorID(authorID int64) (posts []models.Post, err error) {
	if posts, err = pu.postRepo.SelectAllPostsByAuthorID(authorID); err != nil {
		return nil, err
	}
	return posts, nil
}

func (pu *PostUsecase) GetAllPostsByCategories(categories []string) (posts []models.Post, err error) {
	if posts, err = pu.postRepo.SelectPostsByCategories(categories); err != nil {
		return nil, err
	}
	return posts, nil
}
