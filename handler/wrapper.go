package handler

import (
	"github.com/myOmikron/bbb-rec-controller/models"
	"github.com/myOmikron/bbb-rec-controller/modules/bigbluebutton"
)

type Wrapper struct {
	Config *models.Config
	BBB    *bigbluebutton.BBB
}

type errorResponse struct {
	ReturnCode string   `xml:"returncode"`
	Message    string   `xml:"message"`
	MessageKey string   `xml:"messageKey"`
	XMLName    struct{} `xml:"response"`
}
