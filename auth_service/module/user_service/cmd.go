package user_cmd

import (
	user_handler "auth_service/module/user_service/handler"
	user_service "auth_service/module/user_service/service"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/julienschmidt/httprouter"
)

func Cmd(router *httprouter.Router, db *pgxpool.Pool) {
	svc := user_service.NewUserService(db)
	h := user_handler.NewUserHandler(svc)

	router.POST("/api/v1/auth/register", h.Register)
	router.POST("/api/v1/auth/login", h.Login)
}
