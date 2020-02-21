package controllers_test

import (
	"context"
	"nomni/utils/validator"
	"os"
	"testing"

	"go-api/config"
	"go-api/models"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/labstack/echo"
	"github.com/pangpanglabs/goutils/behaviorlog"
	"github.com/pangpanglabs/goutils/echomiddleware"
	"github.com/sirupsen/logrus"

	"log"
	"net/http"

	configutil "github.com/pangpanglabs/goutils/config"
	"github.com/pangpanglabs/goutils/ctxdb"
	"github.com/pangpanglabs/goutils/kafka"
)

var (
	echoApp          *echo.Echo
	handleWithFilter func(handlerFunc echo.HandlerFunc, c echo.Context) error
)

func TestMain(m *testing.M) {
	db := enterTest()
	code := m.Run()
	exitTest(db)
	os.Exit(code)
}

func enterTest() *xorm.Engine {
	configutil.SetConfigPath("../")
	c := config.Init(os.Getenv("APP_ENV"))
	xormEngine, err := xorm.NewEngine(c.Database.Driver, c.Database.Connection)
	if err != nil {
		panic(err)
	}
	if err = models.InitTable(xormEngine); err != nil {
		panic(err)
	}

	echoApp = echo.New()
	echoApp.Validator = validator.New()
	behaviorlog.SetLogLevel(logrus.InfoLevel)
	behaviorlogger := echomiddleware.BehaviorLogger(c.ServiceName, c.BehaviorLog.Kafka)
	db := ContextDB("test", xormEngine, echomiddleware.KafkaConfig{})

	headerCtx := func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			c.SetRequest(req)
			return next(c)
		}
	}

	handleWithFilter = func(handlerFunc echo.HandlerFunc, c echo.Context) error {
		return behaviorlogger(headerCtx(db(handlerFunc)))(c)
	}
	return xormEngine
}

func exitTest(db *xorm.Engine) {
	// if err := models.DropTables(db); err != nil {
	// 	panic(err)
	// }
}

func ContextDB(service string, xormEngine *xorm.Engine, kafkaConfig kafka.Config) echo.MiddlewareFunc {
	return ContextDBWithName(service, echomiddleware.ContextDBName, xormEngine, kafkaConfig)
}
func ContextDBWithName(service string, contexDBName echomiddleware.ContextDBType, xormEngine *xorm.Engine, kafkaConfig kafka.Config) echo.MiddlewareFunc {
	db := ctxdb.New(xormEngine, service, kafkaConfig)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			ctx := req.Context()

			session := db.NewSession(ctx)
			defer session.Close()

			c.SetRequest(req.WithContext(context.WithValue(ctx, contexDBName, session)))

			switch req.Method {
			case "POST", "PUT", "DELETE", "PATCH":
				if err := session.Begin(); err != nil {
					log.Println(err)
				}
				if err := next(c); err != nil {
					session.Rollback()
					return err
				}
				if c.Response().Status >= 500 {
					session.Rollback()
					return nil
				}
				if err := session.Rollback(); err != nil { // rollback data
					return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
				}
			default:
				return next(c)
			}

			return nil
		}
	}
}
