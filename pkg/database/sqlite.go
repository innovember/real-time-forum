package database

import (
	"context"
	"database/sql"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

func GetDBInstance(driver, dbURI string) (*sql.DB, error) {
	var (
		dbConn            *sql.DB
		err               error
		enableForeignKeys = "?_foreign_keys=on"
	)
	if dbConn, err = sql.Open(driver, dbURI+enableForeignKeys); err != nil {
		return nil, err
	}
	dbConn.SetMaxIdleConns(100)
	if err = dbConn.Ping(); err != nil {
		return nil, err
	}
	return dbConn, nil
}

func FileExist(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}
	return true
}

func CreateDir(path string) error {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func UploadSchemesToDB(dbConn *sql.DB, schemesDir string) error {
	schemes, err := getSchemes(schemesDir)
	if err != nil {
		return err
	}
	tx, err := dbConn.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	for _, scheme := range schemes {
		_, err = tx.Exec(scheme)
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Fatal(rollbackErr.Error())
			}
			return err
		}
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func getSchemes(schemesDir string) ([]string, error) {
	var (
		schemes []string
	)
	files, err := ioutil.ReadDir(schemesDir)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		fileName := filepath.Join(schemesDir, file.Name())
		data, err := ioutil.ReadFile(fileName)
		if err != nil {
			return nil, err
		}
		schemes = append(schemes, string(data))
	}
	return schemes, nil

}
