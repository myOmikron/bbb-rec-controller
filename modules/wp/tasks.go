package wp

import (
	"context"
	"fmt"
	"github.com/tebeka/selenium/firefox"
	"net"
	"os"

	"github.com/myOmikron/bbb-rec-controller/models"
	"github.com/myOmikron/echotools/worker"
	"github.com/tebeka/selenium"
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
	ctx := context.WithValue(context.Background(), "selenium", w.Driver)
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

func CreateSeleniumWorker(conf *models.Config) func() (worker.Worker, error) {
	return func() (worker.Worker, error) {
		port, err := pickUnusedPort()
		if err != nil {
			return nil, err
		}

		opts := []selenium.ServiceOption{
			selenium.Output(os.Stdout),
		}

		if conf.Selenium.UseChromium {
			opts = append(opts, selenium.ChromeDriver(conf.Selenium.ChromedriverPath))
		} else {
			opts = append(opts, selenium.GeckoDriver(conf.Selenium.GeckoDriverPath))
		}

		if !conf.Selenium.DisableHeadless {
			opts = append(opts, selenium.StartFrameBuffer())
		}

		service, err := selenium.NewSeleniumService(conf.Selenium.SeleniumPath, port, opts...)
		if err != nil {
			return nil, err
		}

		var caps selenium.Capabilities
		if conf.Selenium.UseChromium {
			caps = selenium.Capabilities{"browserName": "chrome"}
		} else {
			caps = selenium.Capabilities{"browserName": "firefox"}
			f := firefox.Capabilities{
				Binary: conf.Selenium.FirefoxPath,
			}
			caps.AddFirefox(f)
		}

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
