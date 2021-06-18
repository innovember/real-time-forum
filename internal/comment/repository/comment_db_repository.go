package repository

import (
	"context"
	"database/sql"

	"github.com/innovember/real-time-forum/internal/comment"
	"github.com/innovember/real-time-forum/internal/consts"
	"github.com/innovember/real-time-forum/internal/helpers"
	"github.com/innovember/real-time-forum/internal/models"
	"github.com/innovember/real-time-forum/internal/user"
)

type CommentDBRepository struct {
	dbConn   *sql.DB
	userRepo user.UserRepository
}

func NewCommentDBRepository(conn *sql.DB, userRepo user.UserRepository) comment.CommentRepository {
	return &CommentDBRepository{
		dbConn:   conn,
		userRepo: userRepo,
	}
}

func (cr *CommentDBRepository) SelectCommentsNumberByPostID(postID int64) (commentsNumber int64, err error) {
	var (
		ctx context.Context
		tx  *sql.Tx
	)
	ctx = context.Background()
	if tx, err = cr.dbConn.BeginTx(ctx, &sql.TxOptions{}); err != nil {
		return 0, err
	}
	if err = tx.QueryRow(`
	SELECT COUNT(id)
	FROM comments
	WHERE post_id = ?`, postID).Scan(&commentsNumber); err != nil {
		tx.Rollback()
		if err == consts.ErrNoData {
			return 0, nil
		}
		return 0, err
	}
	if err = tx.Commit(); err != nil {
		return 0, err
	}
	return commentsNumber, nil
}

func (cr *CommentDBRepository) Insert(comment *models.Comment) (err error) {
	var (
		ctx    context.Context
		tx     *sql.Tx
		result sql.Result
	)
	ctx = context.Background()
	if tx, err = cr.dbConn.BeginTx(ctx, &sql.TxOptions{}); err != nil {
		return err
	}
	if result, err = tx.Exec(`
	INSERT INTO comments(author_id,post_id,content, created_at)
	VALUES(?,?,?,?)`, comment.AuthorID, comment.PostID, comment.Content, helpers.GetCurrentUnixTime()); err != nil {
		tx.Rollback()
		return err
	}
	if _, err = result.LastInsertId(); err != nil {
		tx.Rollback()
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (cr *CommentDBRepository) SelectCommentsByPostID(input *models.InputGetComments) (comments []models.Comment, err error) {
	var (
		ctx   context.Context
		tx    *sql.Tx
		rows  *sql.Rows
		total int
	)
	ctx = context.Background()
	if tx, err = cr.dbConn.BeginTx(ctx, &sql.TxOptions{}); err != nil {
		return nil, err
	}
	if err = tx.QueryRow(`SELECT count(id) AS total
	 					FROM comments;
						 `).Scan(
		&total); err != nil {
		tx.Rollback()
		return nil, err
	}
	if input.LastCommentID == 0 {
		input.LastCommentID = total + 1
	}
	if rows, err = tx.Query(`
		SELECT *
		FROM comments
		WHERE id < ?
                AND post_id = ?
		ORDER BY created_at DESC
		LIMIT ?
		`,
		input.LastCommentID,
                input.PostID,
		input.Limit); err != nil {
		tx.Rollback()
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			c    models.Comment
			user *models.User
		)
		rows.Scan(&c.ID, &c.AuthorID,
			&c.PostID, &c.Content,
			&c.CreatedAt)
		user, err = cr.userRepo.SelectByID(c.AuthorID)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		c.Author = user
		comments = append(comments, c)
	}
	err = rows.Err()
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return comments, nil
}

func (cr *CommentDBRepository) SelectCommentsByAuthorID(input *models.InputGetComments) (comments []models.Comment, err error) {
	var (
		ctx   context.Context
		tx    *sql.Tx
		rows  *sql.Rows
		total int
	)
	ctx = context.Background()
	if tx, err = cr.dbConn.BeginTx(ctx, &sql.TxOptions{}); err != nil {
		return nil, err
	}
	if err = tx.QueryRow(`SELECT count(id) AS total
	 					FROM comments;
						 `).Scan(
		&total); err != nil {
		tx.Rollback()
		return nil, err
	}
	if input.LastCommentID == 0 {
		input.LastCommentID = total + 1
	}
	if rows, err = tx.Query(`
		SELECT *
		FROM comments
		WHERE author_id = ?
		AND id < ?
		ORDER BY created_at DESC
		LIMIT ?
		`,
		input.UserID,
		input.LastCommentID,
		input.Limit); err != nil {
		tx.Rollback()
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			c    models.Comment
			user *models.User
		)
		rows.Scan(&c.ID, &c.AuthorID,
			&c.PostID, &c.Content,
			&c.CreatedAt)
		user, err = cr.userRepo.SelectByID(c.AuthorID)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		c.Author = user
		comments = append(comments, c)
	}
	err = rows.Err()
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return comments, nil
}

func (cr *CommentDBRepository) SelectCommentByID(commentID int64) (comment *models.Comment, err error) {
	var (
		p    models.Post
		ctx  context.Context
		tx   *sql.Tx
		user *models.User
		c    models.Comment
	)
	ctx = context.Background()
	if tx, err = cr.dbConn.BeginTx(ctx, &sql.TxOptions{}); err != nil {
		return nil, err
	}
	if err = tx.QueryRow(`
	SELECT * FROM comments
	 WHERE id = ?`, commentID,
	).Scan(&c.ID, &c.AuthorID,
		&c.PostID, &c.Content,
		&c.CreatedAt); err != nil {
		tx.Rollback()
		return nil, err
	}
	user, err = cr.userRepo.SelectByID(p.AuthorID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	c.Author = user
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return &c, nil
}
