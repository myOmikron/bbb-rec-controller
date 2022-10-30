package server

import (
	"github.com/labstack/echo/v4"
	"github.com/myOmikron/echotools/execution"
)

func StartServer(configPath string) {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	initializeMiddleware(e)

	defineRoutes(e)

	execution.SignalStart(e, "", &execution.Config{
		ReloadFunc: func() {
			StartServer(configPath)
		},
		StopFunc: func() {

		},
		TerminateFunc: func() {

		},
	})
}
