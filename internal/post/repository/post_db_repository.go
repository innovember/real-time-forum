package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/innovember/real-time-forum/internal/helpers"
	"github.com/innovember/real-time-forum/internal/models"
	"github.com/innovember/real-time-forum/internal/post"
	"github.com/innovember/real-time-forum/internal/user"
)

type PostDBRepository struct {
	dbConn   *sql.DB
	userRepo user.UserRepository
}

func NewPostDBRepository(conn *sql.DB, userRepo user.UserRepository) post.PostRepository {
	return &PostDBRepository{
		dbConn:   conn,
		userRepo: userRepo,
	}
}

func (pr *PostDBRepository) Insert(post *models.Post) (*models.Post, error) {
	var (
		ctx    context.Context
		tx     *sql.Tx
		result sql.Result
		err    error
	)
	ctx = context.Background()
	if tx, err = pr.dbConn.BeginTx(ctx, &sql.TxOptions{}); err != nil {
		return nil, err
	}
	if result, err = tx.Exec(`
	INSERT INTO posts(author_id,title, content, created_at)
	VALUES(?,?,?,?)`, post.AuthorID, post.Title,
		post.Content, helpers.GetCurrentUnixTime()); err != nil {
		tx.Rollback()
		return nil, err
	}
	if post.ID, err = result.LastInsertId(); err != nil {
		tx.Rollback()
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return post, nil
}

func (pr *PostDBRepository) SelectAllPosts() ([]models.Post, error) {
	var (
		rows  *sql.Rows
		ctx   context.Context
		tx    *sql.Tx
		err   error
		posts []models.Post
	)
	ctx = context.Background()
	if tx, err = pr.dbConn.BeginTx(ctx, &sql.TxOptions{}); err != nil {
		return nil, err
	}
	if rows, err = tx.Query(`
		SELECT *
		FROM posts
		ORDER BY created_at DESC
		`); err != nil {
		tx.Rollback()
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			p    models.Post
			user *models.User
		)
		rows.Scan(&p.ID, &p.AuthorID, &p.Title, &p.Content,
			&p.CreatedAt)
		user, err = pr.userRepo.SelectByID(p.AuthorID)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		p.Author = user
		if err = pr.SelectCategories(&p); err != nil {
			tx.Rollback()
			return nil, err
		}
		// if p.CommentsNumber, err = commentRepo.GetCommentsNumberByPostID(p.ID); err != nil {
		// 	tx.Rollback()
		// 	return nil, err
		// }
		posts = append(posts, p)
	}
	err = rows.Err()
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return posts, nil
}

func (pr *PostDBRepository) SelectPostByID(postID int64) (post *models.Post, err error) {
	var (
		p    models.Post
		ctx  context.Context
		tx   *sql.Tx
		user *models.User
	)
	ctx = context.Background()
	if tx, err = pr.dbConn.BeginTx(ctx, &sql.TxOptions{}); err != nil {
		return nil, err
	}
	if err = tx.QueryRow(`
	SELECT * FROM posts WHERE id = ?`, postID,
	).Scan(&p.ID, &p.AuthorID, &p.Title,
		&p.Content, &p.CreatedAt); err != nil {
		tx.Rollback()
		return nil, err
	}
	user, err = pr.userRepo.SelectByID(p.AuthorID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	p.Author = user
	if err = pr.SelectCategories(&p); err != nil {
		tx.Rollback()
		return nil, err
	}
	// if p.CommentsNumber, err = commentRepo.GetCommentsNumberByPostID(p.ID); err != nil {
	// 	return nil, err
	// }
	return &p, nil
}

func (pr *PostDBRepository) SelectCategories(post *models.Post) (err error) {
	var (
		rows       *sql.Rows
		ctx        context.Context
		tx         *sql.Tx
		categories []models.Category
	)
	ctx = context.Background()
	if tx, err = pr.dbConn.BeginTx(ctx, &sql.TxOptions{}); err != nil {
		return err
	}
	if rows, err = tx.Query(`
		SELECT c.id,c.name
		FROM categories c
		LEFT JOIN posts_categories pc
		ON pc.post_id = ?
		WHERE c.id = pc.category_id`,
		post.ID); err != nil {
		tx.Rollback()
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var c models.Category
		rows.Scan(&c.ID, &c.Name)
		categories = append(categories, c)
	}
	err = rows.Err()
	if err != nil {
		tx.Rollback()
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	post.Categories = categories
	return nil
}

func (pr *PostDBRepository) SelectPostsByCategories(categories []string) (posts []models.Post, err error) {
	var (
		rows           *sql.Rows
		ctx            context.Context
		tx             *sql.Tx
		categoriesList = strings.Join(categories, ", ")
	)
	ctx = context.Background()
	if tx, err = pr.dbConn.BeginTx(ctx, &sql.TxOptions{}); err != nil {
		return nil, err
	}
	query := fmt.Sprintf(`
		SELECT p.*
		FROM posts_categories as pc
		INNER JOIN posts as p
		ON p.id = pc.post_id
		INNER JOIN categories as c
		ON c.id=pc.category_id
		WHERE c.name in (%s)
		GROUP BY p.id
		HAVING COUNT(DISTINCT c.id) = %d
		ORDER BY p.created_at DESC`, categoriesList, len(categories))
	if rows, err = tx.Query(query); err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			p    models.Post
			user *models.User
		)
		rows.Scan(&p.ID, &p.AuthorID, &p.Title, &p.Content,
			&p.CreatedAt)
		user, err := pr.userRepo.SelectByID(p.AuthorID)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		p.Author = user
		if err = pr.SelectCategories(&p); err != nil {
			tx.Rollback()
			return nil, err
		}
		// if p.CommentsNumber, err = commentRepo.GetCommentsNumberByPostID(p.ID); err != nil {
		// 	return nil, err
		// }
		posts = append(posts, p)
	}
	err = rows.Err()
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return posts, nil
}

func (pr *PostDBRepository) SelectAllPostsByAuthorID(authorID int64) (posts []models.Post, err error) {
	var (
		rows *sql.Rows
		ctx  context.Context
		tx   *sql.Tx
	)
	ctx = context.Background()
	if tx, err = pr.dbConn.BeginTx(ctx, &sql.TxOptions{}); err != nil {
		return nil, err
	}
	if rows, err = tx.Query(`
		SELECT *
		FROM posts
		WHERE author_id = ?
		ORDER BY created_at DESC
		`, authorID); err != nil {
		tx.Rollback()
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			p    models.Post
			user *models.User
		)
		rows.Scan(&p.ID, &p.AuthorID, &p.Title, &p.Content,
			&p.CreatedAt)
		user, err := pr.userRepo.SelectByID(p.AuthorID)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		p.Author = user
		if err = pr.SelectCategories(&p); err != nil {
			tx.Rollback()
			return nil, err
		}
		// if p.CommentsNumber, err = commentRepo.GetCommentsNumberByPostID(p.ID); err != nil {
		// 	return nil, status, err
		// }
		posts = append(posts, p)
	}
	err = rows.Err()
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return posts, nil
}
