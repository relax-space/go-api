package models_test

import (
	"bytes"
	"context"
	"io/ioutil"
	"os"
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
	db, err = xorm.NewEngine(os.Getenv("SQL_DRIVER"), os.Getenv("Fruit_CONN"))
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
