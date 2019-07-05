package main

import (
	"context"
	"flag"
	"fmt"
	"go-api/controllers"
	"go-api/models"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/pangpanglabs/echoswagger"
	configutil "github.com/pangpanglabs/goutils/config"
	"github.com/pangpanglabs/goutils/echomiddleware"

	"go-api/factory"
)

var (
	handleWithFilter func(handlerFunc echo.HandlerFunc, c echo.Context) error
)

func main() {
	appEnv := flag.String("app-env", os.Getenv("APP_ENV"), "app env")
	fruitConnEnv := flag.String("FRUIT_CONN", os.Getenv("FRUIT_CONN"), "FRUIT_CONN")
	sqlEnv := flag.String("SQL_DRIVER", os.Getenv("SQL_DRIVER"), "SQL_DRIVER")
	jwtEnv := flag.String("JWT_SECRET", os.Getenv("JWT_SECRET"), "JWT_SECRET")
	flag.Parse()

	var c Config
	if err := configutil.Read(*appEnv, &c); err != nil {
		panic(err)
	}

	fmt.Println("Config===", c)
	db, err := initDB(*sqlEnv, *fruitConnEnv)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	e := echo.New()

	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	//for subdomain  just like https://staging.p2shop.com.cn/fruit/docs
	middlePath := ""
	if *appEnv != "" {
		middlePath = "fruit"
	}
	r := echoswagger.New(e, middlePath, "docs", &echoswagger.Info{
		Title:       "Sample Fruit API",
		Description: "This is docs for fruit service",
		Version:     "1.0.0",
	})
	r.AddSecurityAPIKey("Authorization", "JWT token", echoswagger.SecurityInHeader)
	r.SetUI(echoswagger.UISetting{
		HideTop: true,
	})
	controllers.FruitApiController{}.Init(r.Group("fruits", "/fruits"))
	controllers.FruitApiController{}.Init(r.Group("v1/fruits", "/v1/fruits"))
	e.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(*jwtEnv),
		Skipper: func(c echo.Context) bool {
			ignore := []string{
				"/ping",
				"/fruits",
				"/sign",
				"/docs",
			}

			for _, i := range ignore {
				if strings.HasPrefix(c.Request().URL.Path, i) {
					return true
				}
			}

			return false
		},
	}))

	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.Use(middleware.RequestID())
	e.Use(echomiddleware.ContextLogger())
	e.Use(echomiddleware.ContextDB(c.Service, db, echomiddleware.KafkaConfig(c.Logger.Kafka)))
	e.Use(echomiddleware.BehaviorLogger(c.Service, echomiddleware.KafkaConfig(c.BehaviorLog.Kafka)))

	e.Validator = &Validator{}

	e.Debug = c.Debug

	configMap := map[string]interface{}{
		"key": "123",
	}
	setContextValueMiddleware := setContextValue(&configMap)
	handleWithFilter = func(handlerFunc echo.HandlerFunc, c echo.Context) error {
		return setContextValueMiddleware(handlerFunc)(c)
	}
	e.Use(setContextValueMiddleware)
	if err := e.Start(":" + c.HttpPort); err != nil {
		log.Println(err)
	}
}

func setContextValue(configMap *map[string]interface{}) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			reqContext := context.WithValue(req.Context(), factory.ContextConfigName, configMap)
			c.SetRequest(req.WithContext(reqContext))
			return next(c)
		}
	}
}

func initDB(driver, connection string) (*xorm.Engine, error) {
	db, err := xorm.NewEngine(driver, connection)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxLifetime(10 * time.Second) //How long will the idle connection in the connection pool remain? 连接池中的空闲连接保持多久
	db.SetMaxOpenConns(30)                  //Maximum number of connections allowed to open 允许打开的最大连接数
	db.SetMaxIdleConns(10)                  //Several idle connections remain in the connection pool when all connections are idle 所有连接空闲时，连接池中保留几个空闲连接
	db.Sync2(new(models.Fruit))
	return db, nil
}

type Config struct {
	Logger struct {
		Kafka echomiddleware.KafkaConfig
	}
	BehaviorLog struct {
		Kafka echomiddleware.KafkaConfig
	}
	Trace struct {
		Zipkin echomiddleware.ZipkinConfig
	}

	Debug    bool
	Service  string
	HttpPort string
}

type Validator struct{}

func (v *Validator) Validate(i interface{}) error {
	_, err := govalidator.ValidateStruct(i)
	return err
}
