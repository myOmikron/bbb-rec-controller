package wp

import (
	"context"
	"fmt"
	"net"

	"github.com/labstack/echo/v4"
	"github.com/myOmikron/echotools/worker"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/firefox"

	"github.com/myOmikron/bbb-rec-controller/models"
)

type w struct {
	queue   chan worker.Task
	quit    chan bool
	Driver  selenium.WebDriver
	Service *selenium.Service
}

func (w *w) SetQueue(c chan worker.Task) {
	w.queue = c
}

func (w *w) Start() {
	ctx := context.WithValue(context.Background(), "conn", "todo!")
	for {
		select {
		case <-w.quit:
			return
		case t := <-w.queue:
			t.ExecuteWithContext(ctx)
		}
	}
}

func (w *w) Stop() {
	go func() {
		w.Driver.Quit()
		w.Service.Stop()
		w.quit <- true
	}()
}

func pickUnusedPort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	port := l.Addr().(*net.TCPAddr).Port
	if err := l.Close(); err != nil {
		return 0, err
	}
	return port, nil
}

func CreateSeleniumWorker(conf *models.Config, logger echo.Logger) func() (worker.Worker, error) {
	return func() (worker.Worker, error) {
		port, err := pickUnusedPort()
		if err != nil {
			return nil, err
		}

		opts := []selenium.ServiceOption{
			selenium.GeckoDriver(conf.Selenium.GeckoDriverPath),
			selenium.Output(logger.Output()),
		}

		if conf.Selenium.DisableHeadless {
			opts = append(opts, selenium.StartFrameBuffer())
		}

		service, err := selenium.NewSeleniumService(conf.Selenium.SeleniumPath, port, opts...)
		if err != nil {
			return nil, err
		}

		caps := selenium.Capabilities{"browserName": "firefox"}
		f := firefox.Capabilities{
			Binary: conf.Selenium.FirefoxPath,
		}
		caps.AddFirefox(f)

		wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
		if err != nil {
			return nil, err
		}

		return &w{
			quit:    make(chan bool),
			Driver:  wd,
			Service: service,
		}, nil
	}
}
