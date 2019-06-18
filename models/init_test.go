package models_test

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	_ "github.com/denisenkom/go-mssqldb"
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

	if getDbTypeName() == "mysql" {
		if err = resetDb(db, "../test_info", "*"); err != nil {
			fmt.Println(err)
		}
	}

	ctx = context.WithValue(context.Background(), echomiddleware.ContextDBName, db.NewSession())
}

func getDbTypeName() (name string) {
	drive := os.Getenv("SQL_DRIVER")
	switch drive {
	case "mssql":
		name = "sqlserver"
	case "mysql":
		name = "mysql"
	}
	return
}

func resetDb(db *xorm.Engine, folderName, databaseName string) (err error) {
	fileName := fmt.Sprintf("%v/%v/%v.sql", folderName, getDbTypeName(), databaseName)
	files, err := filepath.Glob(fileName)
	if err != nil {
		return
	}
	for _, f := range files {
		if err = importView(db, f); err != nil {
			return
		}
	}
	return
}

func importView(db *xorm.Engine, fileName string) error {
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
