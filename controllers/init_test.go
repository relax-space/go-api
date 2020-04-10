package controllers_test

import (
	"context"
	"os"
	"testing"
	"path/filepath"
	"io/ioutil"
	"io"
	"time"
	"log"
	"net/http"

	"github.com/relax-space/go-api/config"
	"github.com/relax-space/go-api/models"

	_ "github.com/go-sql-driver/mysql"
	"github.com/asaskevich/govalidator"
	"github.com/go-xorm/xorm"
	"github.com/labstack/echo"
	"github.com/pangpanglabs/goutils/behaviorlog"
	"github.com/pangpanglabs/goutils/echomiddleware"
	"github.com/sirupsen/logrus"
	configutil "github.com/pangpanglabs/goutils/config"
	"github.com/pangpanglabs/goutils/ctxdb"
	"github.com/pangpanglabs/goutils/kafka"
	"github.com/pangpanglabs/goutils/httpreq"
	"github.com/pangpanglabs/goutils/jwtutil"
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
	if err = initData(xormEngine, true); err != nil {
		panic(err)
	}

	echoApp = echo.New()
	echoApp.Validator = &Validator{}

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

type Validator struct{}

func (v *Validator) Validate(i interface{}) error {
	_, err := govalidator.ValidateStruct(i)
	return err
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



func initData(xormEngine *xorm.Engine, isDownload bool) error {
	if err := models.DropTables(xormEngine); err != nil {
		return err
	}
	if err := models.InitTable(xormEngine); err != nil {
		return err
	}
	if err := loadData(xormEngine, isDownload); err != nil {
		return err
	}
	return nil
}

func loadData(db *xorm.Engine, isDownload bool) (err error) {
	if isDownload {
		urlStr := "https://dmz-staging.p2shop.com.cn/rtc-dmz-api/v1/dbfiles?nsPrefix=pangpang&nsSuffix=&dbName=fruit"
		writeUrl(urlStr, "test.sql", getToken())
	}
	files, err := filepath.Glob("*.sql")
	if err != nil {
		return
	}
	for _, f := range files {
		if err = importFile(db, f); err != nil {
			return
		}
	}
	return
}

func importFile(db *xorm.Engine, fileName string) error {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}
	_, err = db.Exec(string(b))
	if err != nil {
		return err
	}
	return nil
}
func writeUrl(url, fileName, jwtToken string) (err error) {
	if len(jwtToken) ==0{// don't download
		return
	}
	req := httpreq.New(http.MethodGet, url, nil, func(httpReq *httpreq.HttpReq) error {
		httpReq.RespDataType = httpreq.ByteArrayType
		return nil
	})
	resp, err := req.WithToken(jwtToken).RawCall()
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK{
		return
	}
	out, err := os.OpenFile(fileName, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.ModePerm)
	if err != nil {
		return
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	return
}

func getToken() string {
	jwtKey :=os.Getenv("JWT_SECRET")
	if len(jwtKey) ==0{
		return ""
	}
	token, _ := jwtutil.NewTokenWithSecret(map[string]interface{}{
		"aud": "membership", "tenantCode": "pangpang", "iss": "membership",
		"nbf": time.Now().Add(-5 * time.Minute).Unix(),
	}, jwtKey)
	return token
}

