package models_test

import (
	"bytes"
	"context"
	"go-api/config"
	"go-api/factory"
	"go-api/models"

	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	configutil "github.com/pangpanglabs/goutils/config"
	"github.com/pangpanglabs/goutils/echomiddleware"
	"github.com/pangpanglabs/goutils/httpreq"
	"github.com/pangpanglabs/goutils/jwtutil"
)

var ctx context.Context

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
	// xormEngine.ShowSQL(true)
	if err = initData(xormEngine, true); err != nil {
		panic(err)
	}
	ctx = context.WithValue(context.Background(), echomiddleware.ContextDBName, xormEngine)
	return xormEngine
}

func exitTest(db *xorm.Engine) {
	//db.Close()
}

func rollback() {
	db := factory.DB(ctx).(*xorm.Engine)
	if err := initData(db, false); err != nil {
		panic(err)
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
	_, err = db.Import(bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	return nil
}
func writeUrl(url, fileName, jwtToken string) (err error) {
	req := httpreq.New(http.MethodGet, url, nil, func(httpReq *httpreq.HttpReq) error {
		httpReq.RespDataType = httpreq.ByteArrayType
		return nil
	})
	resp, err := req.WithToken(jwtToken).RawCall()
	if err != nil {
		return
	}
	defer resp.Body.Close()

	out, err := os.OpenFile(fileName, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.ModePerm)
	if err != nil {
		return
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	return
}

func getToken() string {
	token, _ := jwtutil.NewTokenWithSecret(map[string]interface{}{
		"aud": "membership", "tenantCode": "pangpang", "iss": "membership",
		"nbf": time.Now().Add(-5 * time.Minute).Unix(),
	}, os.Getenv("JWT_SECRET"))
	return token
}
