package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/relax-space/go-api/tests"
	"github.com/pangpanglabs/goutils/test"

	"github.com/labstack/echo"
)


func TestInitEcho(t *testing.T) {
	t.Run("config", func(t *testing.T) {
		exp := struct{
			Driver string
			Conn string
			LoggerKafkaBrokers string
			LoggerKafkaTopic string
			ServiceName string
			HTTPPort string
			
		}{
			"mysql",
			"root:1234@tcp(127.0.0.1:3308)/fruit?charset=utf8&parseTime=True&loc=UTC&multiStatements=true",
			"127.0.0.1:9092",
			"sqllog",
			"go-api",
			"8080",
		}
		_,c,_ := initEcho()
		test.Equals(t,exp.Driver,c.Database.Driver)
		test.Equals(t,exp.Conn,c.Database.Connection)
		test.Equals(t,exp.LoggerKafkaBrokers,c.Database.Logger.Kafka.Brokers[0])
		test.Equals(t,exp.LoggerKafkaTopic,c.Database.Logger.Kafka.Topic)
		test.Equals(t,exp.ServiceName,c.ServiceName)
		test.Equals(t,exp.HTTPPort,c.HTTPPort)
	})
	

	t.Run("echo", func(t *testing.T) {
		echoApp,_,_ := initEcho()
		req := httptest.NewRequest(echo.GET, "/ping", nil)
		rec := httptest.NewRecorder()
		c := echoApp.NewContext(req, rec)
		test.Ok(t, ping(c))
		test.Equals(t, "pong", rec.Body.String())
		test.Equals(t, http.StatusOK, rec.Code)
	})
	
}

