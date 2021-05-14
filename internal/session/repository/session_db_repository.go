package repository

import (
	"context"
	"database/sql"

	"github.com/innovember/real-time-forum/internal/consts"
	"github.com/innovember/real-time-forum/internal/helpers"
	"github.com/innovember/real-time-forum/internal/models"
	"github.com/innovember/real-time-forum/internal/session"
)

type SessionDBRepository struct {
	dbConn *sql.DB
}

func NewSessionDBRepository(dbConn *sql.DB) session.SessionRepository {
	return &SessionDBRepository{
		dbConn: dbConn,
	}
}

func (sr *SessionDBRepository) Insert(session *models.Session) error {
	var (
		ctx    context.Context
		tx     *sql.Tx
		result sql.Result
		err    error
	)
	ctx = context.Background()
	if tx, err = sr.dbConn.BeginTx(ctx, &sql.TxOptions{}); err != nil {
		return err
	}
	if result, err = tx.Exec(`INSERT INTO sessions(user_id, token, expires_at)
								VALUES(?,?,?)`,
		session.UserID,
		session.Token,
		session.ExpiresAt); err != nil {
		tx.Rollback()
		return err
	}
	if _, err = result.LastInsertId(); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (sr *SessionDBRepository) Delete(token string) error {
	var (
		ctx context.Context
		tx  *sql.Tx
		err error
	)
	ctx = context.Background()
	if tx, err = sr.dbConn.BeginTx(ctx, &sql.TxOptions{}); err != nil {
		return err
	}
	if _, err = tx.Exec(`DELETE FROM sessions
		WHERE token = ?`, token); err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (sr *SessionDBRepository) SelectByToken(token string) (*models.Session, error) {
	var (
		ctx     context.Context
		tx      *sql.Tx
		err     error
		session = &models.Session{}
	)
	ctx = context.Background()
	if tx, err = sr.dbConn.BeginTx(ctx, &sql.TxOptions{}); err != nil {
		return nil, err
	}
	if err = tx.QueryRow(`SELECT user_id, token, expires_at
	FROM sessions
	WHERE token = ?`, token).Scan(&session.UserID, &session.Token, &session.ExpiresAt); err != nil {
		tx.Rollback()
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return session, nil
}

func (sr *SessionDBRepository) DeleteTokens() error {
	var (
		ctx   context.Context
		tx    *sql.Tx
		rows  *sql.Rows
		users []int64
		err   error
	)
	ctx = context.Background()
	if tx, err = sr.dbConn.BeginTx(ctx, &sql.TxOptions{}); err != nil {
		return err
	}
	if rows, err = tx.Query(`SELECT user_id
							 FROM sessions
							 WHERE expires_at < ?
		`, helpers.GetCurrentUnixTime()); err != nil && err != consts.ErrNoData {
		tx.Rollback()
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var u int64
		err = rows.Scan(
			&u)
		if err != nil {
			return err
		}
		users = append(users, u)
	}
	err = rows.Err()
	if err != nil {
		return err
	}
	if len(users) > 0 {
		for _, id := range users {
			if _, err = tx.Exec(`UPDATE users
						 SET status = ?
						 WHERE id = ?`, consts.StatusOffline, id); err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	if _, err = tx.Exec(`DELETE FROM sessions
		WHERE expires_at < ?`, helpers.GetCurrentUnixTime()); err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (sr *SessionDBRepository) UpdateStatus(userID int64, status string) (err error) {
	var (
		ctx context.Context
		tx  *sql.Tx
	)
	ctx = context.Background()
	if tx, err = sr.dbConn.BeginTx(ctx, &sql.TxOptions{}); err != nil {
		return err
	}
	if _, err = tx.Exec(`UPDATE users
						 SET status = ?
						 WHERE id = ?`, status, userID); err != nil {
		tx.Rollback()
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}
