package server

import (
	"net/url"

	"github.com/labstack/echo/v4"
	emw "github.com/labstack/echo/v4/middleware"
	mw "github.com/myOmikron/echotools/middleware"

	"github.com/myOmikron/bbb-rec-controller/models"
)

func initializeMiddleware(e *echo.Echo, conf *models.Config) {

	e.Use(emw.Recover())

	logFormat := "${remote_ip} - - [${time_custom}] \"${method} ${path} ${protocol}\" ${status} ${bytes_out}\n"
	customTimeFormat := "2/Jan/2006:15:04:05 -0700"
	e.Use(emw.LoggerWithConfig(emw.LoggerConfig{
		Format:           logFormat,
		CustomTimeFormat: customTimeFormat,
	}))

	// Security unpacking
	allowedHosts := &mw.SecurityConfig{
		AllowedHosts:            []mw.AllowedHost{},
		UseForwardedProtoHeader: conf.Server.UseForwardedProtoHeader,
	}
	for _, allowedHost := range conf.Server.AllowedHosts {
		u, _ := url.Parse(allowedHost)
		https := false
		if u.Scheme == "https" {
			https = true
		}
		allowedHosts.AllowedHosts = append(allowedHosts.AllowedHosts, mw.AllowedHost{
			Host:  u.Host,
			Https: https,
		})
	}
	e.Use(mw.Security(allowedHosts))

	e.Use(emw.BodyLimit("5MB"))
}
