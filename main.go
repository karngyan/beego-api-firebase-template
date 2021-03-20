package main

import (
	"beego-api-firebase-template/controllers/middleware"
	"beego-api-firebase-template/modules/global"
	"fmt"
	"strings"

	_ "beego-api-firebase-template/routers"
	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/filter/cors"
	_ "github.com/lib/pq"
)

func main() {

	if err := orm.RegisterDriver("postgres", orm.DRPostgres); err != nil {
		fmt.Println("register driver failed")
		panic(err)
	}

	sqlConnString, _ := beego.AppConfig.String("sqlconn")
	if err := orm.RegisterDataBase("default", "postgres", strings.TrimSpace(sqlConnString)); err != nil {
		fmt.Println("register database failed")
		panic(err)
	}

	global.Init()

	beego.InsertFilter("/api/*", beego.BeforeRouter, middleware.AuthFilter, beego.WithReturnOnOutput(false))
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowAllOrigins:  true,
		AllowCredentials: true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Origin", "Authorization", "Access-Control-Allow-Origin", "X-User-ApiKey", "X-Token"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin"},
	}), beego.WithReturnOnOutput(false))

	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"

		if err := orm.RunSyncdb("default", false, true); err != nil {
			fmt.Println("run sync db fudged up", err.Error())
		}
	}

	beego.Run()
}
