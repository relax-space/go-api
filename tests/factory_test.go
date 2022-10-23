package tests_test

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"
	"github.com/pangpanglabs/goutils/echomiddleware"
	"github.com/pangpanglabs/goutils/test"
	"github.com/relax-space/go-api/controllers"
)

func TestFactory(t *testing.T) {
	t.Run("DB_panic_ctx", func(t *testing.T) {
		defer func() {
			r := recover()
			test.Equals(t, "DB is not exist in ctx", r)
		}()
		req := httptest.NewRequest(echo.GET, "/v1/fruits/:id", nil)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)
		c.SetRequest(req.WithContext(context.WithValue(req.Context(), echomiddleware.ContextDBName, nil)))
		c.SetParamNames("id")
		c.SetParamValues("1")
		controllers.FruitAPIController{}.GetOne(c)
	})
	t.Run("DB_panic_engine", func(t *testing.T) {
		defer func() {
			r := recover()
			test.Equals(t, "DB is not exist in xorm.Engine", r)
		}()
		req := httptest.NewRequest(echo.GET, "/v1/fruits/:id", nil)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)
		c.SetRequest(req.WithContext(context.WithValue(req.Context(), echomiddleware.ContextDBName, "test")))
		c.SetParamNames("id")
		c.SetParamValues("1")
		controllers.FruitAPIController{}.GetOne(c)
	})
}
