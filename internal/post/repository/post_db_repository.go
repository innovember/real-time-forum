package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/innovember/real-time-forum/internal/comment"
	"github.com/innovember/real-time-forum/internal/helpers"
	"github.com/innovember/real-time-forum/internal/models"
	"github.com/innovember/real-time-forum/internal/post"
	"github.com/innovember/real-time-forum/internal/user"
)

type PostDBRepository struct {
	dbConn      *sql.DB
	userRepo    user.UserRepository
	commentRepo comment.CommentRepository
}

func NewPostDBRepository(
	conn *sql.DB,
	userRepo user.UserRepository,
	commentRepo comment.CommentRepository,
) post.PostRepository {
	return &PostDBRepository{
		dbConn:      conn,
		userRepo:    userRepo,
		commentRepo: commentRepo,
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

func (pr *PostDBRepository) SelectAllPosts(input *models.InputGetPosts) ([]models.Post, error) {
	var (
		rows  *sql.Rows
		ctx   context.Context
		tx    *sql.Tx
		err   error
		posts []models.Post
		total int
	)
	ctx = context.Background()
	if tx, err = pr.dbConn.BeginTx(ctx, &sql.TxOptions{}); err != nil {
		return nil, err
	}
	if err = tx.QueryRow(`SELECT count(id) AS total
	 					FROM posts;
						 `).Scan(
		&total); err != nil {
		tx.Rollback()
		return nil, err
	}
	if input.LastPostID == 0 {
		input.LastPostID = total + 1
	}
	if rows, err = tx.Query(`
		SELECT *
		FROM posts
		WHERE id < ?
		ORDER BY created_at DESC
		LIMIT ?
		`,
		input.LastPostID,
		input.Limit); err != nil {
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
		if p.CommentsNumber, err = pr.commentRepo.SelectCommentsNumberByPostID(p.ID); err != nil {
			tx.Rollback()
			return nil, err
		}
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
	p.CommentsNumber, err = pr.commentRepo.SelectCommentsNumberByPostID(p.ID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
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

func (pr *PostDBRepository) SelectPostsByCategories(input *models.InputGetPosts) (posts []models.Post, err error) {
	var (
		rows           *sql.Rows
		ctx            context.Context
		tx             *sql.Tx
		categoriesList = strings.Join(input.Categories, ", ")
		total          int
	)
	ctx = context.Background()
	if tx, err = pr.dbConn.BeginTx(ctx, &sql.TxOptions{}); err != nil {
		return nil, err
	}
	if err = tx.QueryRow(`SELECT count(id) AS total
	 					FROM posts;
						 `).Scan(
		&total); err != nil {
		tx.Rollback()
		return nil, err
	}
	if input.LastPostID == 0 {
		input.LastPostID = total + 1
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
		WHERE p.id < %d
		ORDER BY p.created_at DESC
		LIMIT %d`,
		categoriesList,
		len(input.Categories),
		input.LastPostID,
		input.Limit)
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
		if p.CommentsNumber, err = pr.commentRepo.SelectCommentsNumberByPostID(p.ID); err != nil {
			tx.Rollback()
			return nil, err
		}
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

func (pr *PostDBRepository) SelectAllPostsByAuthorID(input *models.InputGetPosts) (posts []models.Post, err error) {
	var (
		rows  *sql.Rows
		ctx   context.Context
		tx    *sql.Tx
		total int
	)
	ctx = context.Background()
	if tx, err = pr.dbConn.BeginTx(ctx, &sql.TxOptions{}); err != nil {
		return nil, err
	}
	if err = tx.QueryRow(`SELECT count(id) AS total
	 					FROM posts;
						 `).Scan(
		&total); err != nil {
		tx.Rollback()
		return nil, err
	}
	if input.LastPostID == 0 {
		input.LastPostID = total + 1
	}
	if rows, err = tx.Query(`
		SELECT *
		FROM posts
		WHERE author_id = ?
		AND id < ?
		LIMIT ?
		ORDER BY created_at DESC
		`, input.AuthorID,
		input.LastPostID,
		input.Limit); err != nil {
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
		if p.CommentsNumber, err = pr.commentRepo.SelectCommentsNumberByPostID(p.ID); err != nil {
			tx.Rollback()
			return nil, err
		}
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
