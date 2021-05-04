package main

import (
	"log"
	"path/filepath"

	"github.com/innovember/real-time-forum/config"
	"github.com/innovember/real-time-forum/pkg/database"
)

func main() {
	log.Println("Server is starting...")

	config, err := config.LoadConfig("./config/config.json")
	if err != nil {
		log.Fatalln("config error: ", err)
	}

	if !database.FileExist(filepath.Join(config.GetDBPath(), config.GetDBFilename())) {
		err := database.CreateDir(config.GetDBPath())
		if err != nil {
			log.Fatal("dbDir err: ", err)
		}
	}
	dbConn, err := database.GetDBInstance(config.GetDBDriver(), config.GetProdDBConnString())
	if err != nil {
		log.Fatal("dbConn err: ", err)
	}
	defer dbConn.Close()
	if err := database.UploadSchemesToDB(dbConn, config.GetDBSchemesDir()); err != nil {
		log.Fatal("upload schemes err: ", err)
	}
}
