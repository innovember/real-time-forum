package main

import (
	"fmt"
	"log"

	"github.com/innovember/real-time-forum/config"
)

func main() {
	log.Println("Server is starting...")
	config, err := config.LoadConfig("./config/config.json")
	if err != nil {
		log.Fatalln("config error: ", err)
	}
	fmt.Printf("%+v", config)
}
