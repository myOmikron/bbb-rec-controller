package bigbluebutton

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/myOmikron/bbb-rec-controller/models"
	"net/url"
)

var (
	ErrInvalidURL      = errors.New("the configured bigbluebutton server uri is invalid")
	ErrInvalidUsername = errors.New("the configured username mustn't be empty")
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
	if len(config.Username) == 0 {
		return nil, ErrInvalidUsername
	}
	return &BBB{
		Config:  config,
		BaseUrl: base,
		Logger:  logger,
	}, nil
}
