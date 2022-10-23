package tests_test

import (
	"testing"

	"go-api/config"

	"github.com/pangpanglabs/goutils/test"
)

func TestConfig(t *testing.T) {
	t.Run("Init_unknown_appEnv", func(t *testing.T) {
		c := config.Init("test11")
		test.Equals(t, "", c.ServiceName)
	})
	t.Run("Init_options", func(t *testing.T) {
		c := config.Init("", func(d *config.C) {
			d.ServiceName = "test1"
		})
		test.Equals(t, "test1", c.ServiceName)
	})

}
