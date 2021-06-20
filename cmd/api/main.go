package main

import (
	"log"
	"net/http"
	"path/filepath"

	"github.com/innovember/real-time-forum/internal/models"
	"github.com/innovember/real-time-forum/internal/mwares"
	userDelivery "github.com/innovember/real-time-forum/internal/user/delivery"
	userRepo "github.com/innovember/real-time-forum/internal/user/repository"
	userUsecase "github.com/innovember/real-time-forum/internal/user/usecases"

	sessionDelivery "github.com/innovember/real-time-forum/internal/session/delivery"
	sessionRepo "github.com/innovember/real-time-forum/internal/session/repository"
	sessionUsecase "github.com/innovember/real-time-forum/internal/session/usecases"

	categoryDelivery "github.com/innovember/real-time-forum/internal/category/delivery"
	categoryRepo "github.com/innovember/real-time-forum/internal/category/repository"
	categoryUsecase "github.com/innovember/real-time-forum/internal/category/usecases"

	postDelivery "github.com/innovember/real-time-forum/internal/post/delivery"
	postRepo "github.com/innovember/real-time-forum/internal/post/repository"
	postUsecase "github.com/innovember/real-time-forum/internal/post/usecases"

	commentDelivery "github.com/innovember/real-time-forum/internal/comment/delivery"
	commentRepo "github.com/innovember/real-time-forum/internal/comment/repository"
	commentUsecase "github.com/innovember/real-time-forum/internal/comment/usecases"

	chatDelivery "github.com/innovember/real-time-forum/internal/chat/delivery"
	chatRepo "github.com/innovember/real-time-forum/internal/chat/repository"
	chatUsecase "github.com/innovember/real-time-forum/internal/chat/usecases"

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
	hubs := models.NewRoomHubs()
	userRepository := userRepo.NewUserDBRepository(dbConn)
	sessionRepository := sessionRepo.NewSessionDBRepository(dbConn)
	categoryRepository := categoryRepo.NewCategoryDBRepository(dbConn)
	commentRepository := commentRepo.NewCommentDBRepository(dbConn, userRepository)
	postRepository := postRepo.NewPostDBRepository(dbConn, userRepository, commentRepository)
	roomRepository := chatRepo.NewRoomRepository(dbConn)

	hubRepository := chatRepo.NewHubRepository(hubs)

	userUsecase := userUsecase.NewUserUsecase(userRepository)
	sessionUsecase := sessionUsecase.NewSessionUsecase(sessionRepository)
	categoryUsecase := categoryUsecase.NewCategoryUsecase(categoryRepository)
	postUsecase := postUsecase.NewPostUsecase(postRepository, categoryRepository)
	commentUsecase := commentUsecase.NewCommentUsecase(commentRepository)
	roomUsecase := chatUsecase.NewRoomUsecase(roomRepository, userRepository)
	hubUsecase := chatUsecase.NewHubUsecase(hubRepository, roomRepository)

	go sessionUsecase.DeleteExpiredSessions()

	mux := http.NewServeMux()
	mm := mwares.NewMiddlewareManager(userUsecase, sessionUsecase)

	userHandler := userDelivery.NewUserHandler(userUsecase)
	userHandler.Configure(mux, mm)

	sessionHandler := sessionDelivery.NewSessionHandler(sessionUsecase, userUsecase)
	sessionHandler.Configure(mux, mm)

	categoryHandler := categoryDelivery.NewCategoryHandler(categoryUsecase)
	categoryHandler.Configure(mux, mm)

	postHandler := postDelivery.NewPostHandler(postUsecase, userUsecase)
	postHandler.Configure(mux, mm)

	commentHandler := commentDelivery.NewCommentHandler(userUsecase, postUsecase, commentUsecase)
	commentHandler.Configure(mux, mm)

	chatHandler := chatDelivery.NewChatHandler(roomUsecase, sessionUsecase, hubUsecase)
	chatHandler.Configure(mux, mm)

	log.Println("Server is listening", config.GetLocalServerPath())
	err = http.ListenAndServe(config.GetPort(), mux)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
