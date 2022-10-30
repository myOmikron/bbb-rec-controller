package bigbluebutton

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/myOmikron/bbb-rec-controller/models"
	"net/url"
)

var (
	ErrInvalidURL = errors.New("the configured bigbluebutton server uri is invalid")
)

type BBB struct {
	Config  *models.BigBlueButton
	BaseUrl *url.URL
	Logger  echo.Logger
}

func New(logger echo.Logger, config *models.BigBlueButton) (*BBB, error) {
	base, err := url.Parse(config.ServerURI)
	if err != nil {
		return nil, ErrInvalidURL
	}
	return &BBB{
		Config:  config,
		BaseUrl: base,
		Logger:  logger,
	}, nil
}
