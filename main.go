package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/relax-space/go-api/config"
	"github.com/relax-space/go-api/controllers"
	"github.com/relax-space/go-api/models"

	"github.com/hublabs/common/api"
	
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
	"github.com/go-xorm/xorm"
	"github.com/labstack/echo/middleware"
	"github.com/pangpanglabs/echoswagger"
	"github.com/pangpanglabs/goutils/behaviorlog"
	"github.com/pangpanglabs/goutils/echomiddleware"
	"github.com/sirupsen/logrus"
)

func main() {
	e,c,db:=initEcho()
	defer db.Close()
	if err := e.Start(":" + c.HTTPPort); err != nil {
		log.Println(err)
	}
}

func ping(c echo.Context) error{
	return c.String(http.StatusOK, "pong")
}
func initEcho()(*echo.Echo,config.C,*xorm.Engine){
	c := config.Init(os.Getenv("APP_ENV"))

	fmt.Println("Config===", c)
	db, err := models.InitDB(c.Database.Driver, c.Database.Connection)
	if err != nil {
		panic(err)
	}

	if err := models.InitTable(db); err != nil {
		panic(err)
	}

	e := echo.New()

	e.GET("/ping", ping)

	r := echoswagger.New(e, "docs", &echoswagger.Info{
		Title:       "Sample Fruit API",
		Description: "This is docs for fruit service",
		Version:     "1.0.0",
	})
	r.AddSecurityAPIKey("Authorization", "JWT token", echoswagger.SecurityInHeader)
	r.SetUI(echoswagger.UISetting{
		HideTop: true,
	})
	controllers.FruitAPIController{}.Init(r.Group("fruits", "v1/fruits"))

	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.Use(middleware.RequestID())
	e.Use(echomiddleware.ContextLogger())
	e.Use(echomiddleware.ContextDB(c.ServiceName, db, c.Database.Logger.Kafka))
	e.Use(echomiddleware.BehaviorLogger(c.ServiceName, c.BehaviorLog.Kafka))
	if !strings.HasSuffix(c.Appenv, "production") {
		behaviorlog.SetLogLevel(logrus.InfoLevel)
	}

	api.SetErrorMessagePrefix(c.ServiceName)

	e.Debug = c.Debug
	return e,c,db
}

