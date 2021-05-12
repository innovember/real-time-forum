package repository

import (
	"context"
	"database/sql"

	"github.com/innovember/real-time-forum/internal/category"
	"github.com/innovember/real-time-forum/internal/models"
)

type CategoryDBRepository struct {
	dbConn *sql.DB
}

func NewCategoryDBRepository(conn *sql.DB) category.CategoryRepository {
	return &CategoryDBRepository{dbConn: conn}
}

func (cr *CategoryDBRepository) Insert(postID int64, categories []string) (err error) {
	var (
		ctx        context.Context
		tx         *sql.Tx
		result     sql.Result
		categoryID int64
		isExist    bool
	)
	ctx = context.Background()
	if tx, err = cr.dbConn.BeginTx(ctx, &sql.TxOptions{}); err != nil {
		return err
	}
	for _, category := range categories {
		if isExist, err = cr.IsCategoryExist(category, tx); err != nil {
			tx.Rollback()
			return err
		}
		if !isExist {
			if result, err = tx.Exec(`INSERT INTO categories(name) VALUES(?)`, category); err != nil {
				tx.Rollback()
				return err
			}
			if categoryID, err = result.LastInsertId(); err != nil {
				tx.Rollback()
				return err
			}
		} else {
			if categoryID, err = cr.SelectByName(category, tx); err != nil {
				tx.Rollback()
				return err
			}
		}
		if _, err = tx.Exec(
			`INSERT INTO posts_categories (post_id, category_id)
			VALUES (?, ?)`,
			postID, categoryID,
		); err != nil {
			tx.Rollback()
			return err
		}
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (cr *CategoryDBRepository) SelectAllCategories() (categories []models.Category, err error) {
	var (
		ctx  context.Context
		tx   *sql.Tx
		rows *sql.Rows
	)
	ctx = context.Background()
	if tx, err = cr.dbConn.BeginTx(ctx, &sql.TxOptions{}); err != nil {
		return nil, err
	}
	if rows, err = tx.Query(`SELECT * FROM categories`); err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var c models.Category
		err = rows.Scan(&c.ID, &c.Name)
		if err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return categories, nil
}

func (cr *CategoryDBRepository) SelectByName(name string, tx *sql.Tx) (id int64, err error) {
	if err = tx.QueryRow(`SELECT id FROM categories WHERE name=?`, name).Scan(
		&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (cr *CategoryDBRepository) IsCategoryExist(category string, tx *sql.Tx) (bool, error) {
	var (
		err error
		id  int64
	)
	if err = tx.QueryRow(`SELECT id FROM categories WHERE name=?`, category).Scan(
		&id); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
