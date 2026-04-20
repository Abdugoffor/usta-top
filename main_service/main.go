// @title           Usta Top API
// @version         1.0
// @description     Usta Top - vakansiya va resume platformasi API
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  support@usta-top.uz

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

package main

import (
	"fmt"
	"log"
	"main_service/config"
	_ "main_service/docs"
	"main_service/helper"
	categorya_cmd "main_service/module/categorya_service"
	comment_cmd "main_service/module/comment_service"
	country_cmd "main_service/module/country_service"
	region_cmd "main_service/module/region_service"
	resume_cmd "main_service/module/resume_service"
	user_cmd "main_service/module/user_service"
	vacancy_cmd "main_service/module/vacancy_service"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	helper.LoadEnv()

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
		region_cmd.Cmd(router, db)
		country_cmd.Cmd(router, db)
		categorya_cmd.Cmd(router, db)
		vacancy_cmd.Cmd(router, db)
		resume_cmd.Cmd(router, db)
		comment_cmd.Cmd(router, db)
	}

	router.HandlerFunc(http.MethodGet, "/swagger/*filepath", httpSwagger.WrapHandler)

	port := helper.ENV("APP_PORT")
	{
		if port == "" {
			port = "8080"
		}
	}

	log.Printf("🚀 Server started on :%s", port)
	log.Printf("📖 Swagger UI: http://localhost:%s/swagger/index.html", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
