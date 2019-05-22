package models_test

import (
	"bytes"
	"context"
	"io/ioutil"
	"runtime"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/pangpanglabs/goutils/echomiddleware"
)

var ctx context.Context
var db *xorm.Engine

func init() {
	runtime.GOMAXPROCS(1)
	var err error
	db, err = xorm.NewEngine("mysql", "root:1234@tcp(localhost:3306)/fruit?charset=utf8&parseTime=True&loc=UTC")
	if err != nil {
		panic(err)
	}
	db.ShowSQL(true)

	resetDb(db, "../test_info/table.sql")

	ctx = context.WithValue(context.Background(), echomiddleware.ContextDBName, db.NewSession())
}

func resetDb(db *xorm.Engine, fileName string) error {
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
