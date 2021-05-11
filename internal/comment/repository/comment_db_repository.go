package repository

import (
	"context"
	"database/sql"

	"github.com/innovember/real-time-forum/internal/comment"
)

type CommentDBRepository struct {
	dbConn *sql.DB
}

func NewCommentDBRepository(conn *sql.DB) comment.CommentRepository {
	return &CommentDBRepository{dbConn: conn}
}

func (cr *CommentDBRepository) SelectCommentsNumberByPostID(postID int64) (commentsNumber int, err error) {
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
		return 0, err
	}
	if err = tx.Commit(); err != nil {
		return 0, err
	}
	return commentsNumber, nil
}
