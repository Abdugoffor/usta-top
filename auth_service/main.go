package main

import (
	"auth_service/config"
	"auth_service/helper"
	user_cmd "auth_service/module/user_service"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func main() {
	helper.LoadEnv()

	db := config.DBConnect()

	router := httprouter.New()
	{
		user_cmd.Cmd(router, db)
	}

	port := helper.ENV("APP_PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("🚀 Auth Server started on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
