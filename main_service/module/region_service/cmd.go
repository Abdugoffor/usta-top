package region_cmd

import (
	region_handler "main_service/module/region_service/handler"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/julienschmidt/httprouter"
)

func Cmd(router *httprouter.Router, db *pgxpool.Pool) {
	group := "/api/v1"
	{
		region_handler.NewRegionHandler(router, group, db)
	}
}
