package main

import (
	"github.com/labstack/echo/v4"

	channel "github.com/emika-team/line-oa-manager/pkg/http/channel"
	message "github.com/emika-team/line-oa-manager/pkg/http/message"
)

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Hello, World!")
	})
	e.POST("/webhook", message.GetContent)
	e.GET("/channels", channel.GetChannel)
	e.Logger.Fatal(e.Start(":1323"))
}
