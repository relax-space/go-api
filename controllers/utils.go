package controllers

import (
	"context"
	"fmt"
	"net/http"
	"reflect"

	"github.com/pangpanglabs/goutils/behaviorlog"

	"go-api/factory"
	"nomni/utils/api"

	"github.com/fatih/structs"
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

type UrlInfo struct {
	ControllerName string
	ApiName        string
	Method         string //GET,POST
	Uri            string
	ResponseStatus interface{}
	Struct         interface{}
	Err            error
}

func PrintApiBehaviorError(c context.Context, urlInfo UrlInfo) {
	logContext := behaviorlog.FromCtx(c)
	if logContext != nil {
		logClone := logContext.Clone()
		if urlInfo.Err != nil {
			logClone.WithError(urlInfo.Err)
		}
		logClone.Controller = urlInfo.ControllerName
		logClone.Params = map[string]interface{}{}
		param := make(map[string]interface{}, 0)
		if urlInfo.Struct != nil && !reflect.ValueOf(urlInfo.Struct).IsNil() {
			s := structs.New(urlInfo.Struct)
			s.TagName = "json"
			param = s.Map()
		}
		var statusCode int
		switch t := urlInfo.ResponseStatus.(type) {
		case int:
			statusCode = t
		case *http.Response:
			statusCode = t.StatusCode
		}

		logClone.WithCallURLInfo(urlInfo.Method, urlInfo.Uri, param, statusCode).Log(urlInfo.ApiName)
		logContext.Params = map[string]interface{}{}
	}
}
