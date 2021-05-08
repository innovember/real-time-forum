package repository

import (
	"context"
	"database/sql"

	"github.com/innovember/real-time-forum/internal/helpers"
	"github.com/innovember/real-time-forum/internal/models"
	"github.com/innovember/real-time-forum/internal/user"
)

type UserDBRepository struct {
	dbConn *sql.DB
}

func NewUserDBRepository(conn *sql.DB) user.UserRepository {
	return &UserDBRepository{dbConn: conn}
}

func (ur *UserDBRepository) Insert(user *models.User) (err error) {
	var (
		ctx    context.Context
		tx     *sql.Tx
		result sql.Result
	)
	ctx = context.Background()
	if tx, err = ur.dbConn.BeginTx(ctx, &sql.TxOptions{}); err != nil {
		return err
	}
	if result, err = tx.Exec(`INSERT INTO users
				(nickname,email, password, first_name, last_name, age, gender,created_at, last_active)
			VALUES
				(?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, user.Nickname, user.Email, user.Password,
		user.FirstName, user.LastName, user.Age, user.Gender,
		helpers.GetCurrentUnixTime(), helpers.GetCurrentUnixTime(),
	); err != nil {
		tx.Rollback()
		return err
	}
	if _, err = result.LastInsertId(); err != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (ur *UserDBRepository) SelectByEmail(email string) (*models.User, error) {
	var (
		ctx context.Context
		tx  *sql.Tx
		err error
		u   = &models.User{}
	)
	ctx = context.Background()
	if tx, err = ur.dbConn.BeginTx(ctx, &sql.TxOptions{}); err != nil {
		return nil, err
	}
	if err = tx.QueryRow(`SELECT id, nickname, email, password
						  FROM users
						  WHERE email = ?
	`, email).Scan(&u.ID, &u.Nickname, &u.Email, &u.Password); err != nil {
		if err = tx.Rollback(); err != nil {
			return nil, err
		}
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return u, nil
}

func (ur *UserDBRepository) SelectByNickname(nickname string) (*models.User, error) {
	var (
		ctx context.Context
		tx  *sql.Tx
		err error
		u   = &models.User{}
	)
	ctx = context.Background()
	if tx, err = ur.dbConn.BeginTx(ctx, &sql.TxOptions{}); err != nil {
		return nil, err
	}
	if err = tx.QueryRow(`SELECT id, nickname, email, password
						  FROM users
						  WHERE nickname = ?
	`, nickname).Scan(&u.ID, &u.Nickname, &u.Email, &u.Password); err != nil {
		tx.Rollback()
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return u, nil
}

func (ur *UserDBRepository) SelectByID(userID int64) (*models.User, error) {
	var (
		ctx context.Context
		tx  *sql.Tx
		err error
		u   = &models.User{}
	)
	ctx = context.Background()
	if tx, err = ur.dbConn.BeginTx(ctx, &sql.TxOptions{}); err != nil {
		return nil, err
	}
	if err = tx.QueryRow(`SELECT id, nickname, email, first_name, last_name,
							age, gender,created_at, last_active
						  FROM users
						  WHERE id = ?
	`, userID).Scan(&u.ID, &u.Nickname, &u.Email,
		&u.FirstName, &u.LastName,
		&u.CreatedAt, &u.LastActive); err != nil {
		if err = tx.Rollback(); err != nil {
			return nil, err
		}
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return u, nil
}

func (ur *UserDBRepository) UpdateActivity(userID int64) (err error) {
	var (
		ctx context.Context
		tx  *sql.Tx
	)
	ctx = context.Background()
	if tx, err = ur.dbConn.BeginTx(ctx, &sql.TxOptions{}); err != nil {
		return err
	}
	if _, err = tx.Exec(`UPDATE users
						 SET last_active = ?
						 WHERE id = ?`, helpers.GetCurrentUnixTime(), userID); err != nil {
		tx.Rollback()
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}
