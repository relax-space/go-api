package controllers

import (
	"fmt"
	"net/http"

	"github.com/pangpanglabs/goutils/behaviorlog"

	"go-api/factory"
	"nomni/utils/api"

	"github.com/go-xorm/xorm"
	"github.com/labstack/echo"
)

func ReturnApiFail(c echo.Context, status int, apiError api.Error, detail ...interface{}) error {
	behaviorlog.FromCtx(c.Request().Context()).WithError(apiError)
	for _, d := range detail {
		if d != nil {
			apiError.Details = fmt.Sprint(detail...)
		}
	}
	return c.JSON(status, api.Result{
		Success: false,
		Error:   apiError,
	})
}

func ReturnApiSucc(c echo.Context, status int, result interface{}) error {
	req := c.Request()
	if req.Method == "POST" || req.Method == "PUT" || req.Method == "DELETE" {
		var apiErr api.Error
		switch req.Method {
		case "POST":
			apiErr = api.NotCreatedError()
		case "PUT":
			apiErr = api.NotUpdatedError()
		case "DELETE":
			apiErr = api.NotDeletedError()
		}
		if session, ok := factory.DB(req.Context()).(*xorm.Session); ok {
			err := session.Commit()
			if err != nil {
				return ReturnApiFail(c, http.StatusOK, apiErr, err)
			}
		}
	}

	return c.JSON(status, api.Result{
		Success: true,
		Result:  result,
	})
}
