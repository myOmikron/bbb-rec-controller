package handler

import (
	"github.com/myOmikron/bbb-rec-controller/models"
	"github.com/myOmikron/bbb-rec-controller/modules/bigbluebutton"
	"github.com/myOmikron/echotools/worker"
)

type Wrapper struct {
	Config       *models.Config
	BBB          *bigbluebutton.BBB
	SeleniumPool worker.Pool
}

type errorResponse struct {
	ReturnCode string   `xml:"returncode"`
	Message    string   `xml:"message"`
	MessageKey string   `xml:"messageKey"`
	XMLName    struct{} `xml:"response"`
}
