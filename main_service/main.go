package main

import (
	"fmt"
	"log"
	"main_service/config"
	"main_service/helper"
	"main_service/middleware"
	categorya_cmd "main_service/module/categorya_service"
	comment_cmd "main_service/module/comment_service"
	country_cmd "main_service/module/country_service"
	language_cmd "main_service/module/language_service"
	resume_cmd "main_service/module/resume_service"
	translations_cmd "main_service/module/translations_service"
	upload_cmd "main_service/module/upload_service"
	user_cmd "main_service/module/user_service"
	vacancy_cmd "main_service/module/vacancy_service"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
)

func main() {
	helper.LoadEnv()

	if helper.ENV("JWT_KEY") == "" {
		log.Fatal("❌ JWT_KEY environment variable is required and must not be empty")
	}

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "migrate:create":
			if len(os.Args) < 3 {
				fmt.Println("Usage:   go run . migrate:create <name>")
				fmt.Println("Example: go run . migrate:create add_code_to_regions")
				os.Exit(1)
			}
			config.MigrateCreate(os.Args[2])
			return
		}
	}

	db := config.DBConnect()

	router := httprouter.New()
	{
		user_cmd.Cmd(router, db)
		country_cmd.Cmd(router, db)
		language_cmd.Cmd(router, db)
		translations_cmd.Cmd(router, db)
		categorya_cmd.Cmd(router, db)
		vacancy_cmd.Cmd(router, db)
		resume_cmd.Cmd(router, db)
		comment_cmd.Cmd(router, db)
		upload_cmd.Cmd(router)
	}

	router.ServeFiles("/uploads/*filepath", http.Dir("uploads"))

	port := helper.ENV("APP_PORT")
	{
		if port == "" {
			port = "8080"
		}
	}

	// 4 MB body limit (3MB upload + overhead), 30 req/s per IP with burst 60
	handler := middleware.SecurityHeaders(
		middleware.CORS(
			middleware.RateLimit(30, 60)(
				http.MaxBytesHandler(router, 4<<20),
			),
		),
	)

	log.Printf("🚀 Server started on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
