package server

import (
	"errors"
	"fmt"
	"github.com/labstack/gommon/log"
	"io/fs"
	"net"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/myOmikron/echotools/color"
	"github.com/myOmikron/echotools/execution"
	"github.com/pelletier/go-toml"

	"github.com/myOmikron/bbb-rec-controller/models"
)

func StartServer(configPath string) {
	conf := &models.Config{}

	if configBytes, err := os.ReadFile(configPath); errors.Is(err, fs.ErrNotExist) {
		color.Printf(color.RED, "Config was not found at %s\n", configPath)
		b, _ := toml.Marshal(conf)
		fmt.Print(string(b))
		os.Exit(1)
	} else {
		if err := toml.Unmarshal(configBytes, conf); err != nil {
			panic(err)
		}
	}

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Logger.SetLevel(log.DEBUG)

	initializeMiddleware(e, conf)

	defineRoutes(e, conf)

	color.Printf(color.PURPLE, "Started listening on http://%s\n", net.JoinHostPort(conf.Server.ListenAddress, strconv.Itoa(int(conf.Server.ListenPort))))
	execution.SignalStart(e, net.JoinHostPort(conf.Server.ListenAddress, strconv.Itoa(int(conf.Server.ListenPort))), &execution.Config{
		ReloadFunc: func() {
			StartServer(configPath)
		},
		StopFunc: func() {

		},
		TerminateFunc: func() {

		},
	})
}
