package handler

import "github.com/myOmikron/bbb-rec-controller/models"

type Wrapper struct {
	Config *models.Config
}

type errorResponse struct {
	ReturnCode string   `xml:"returncode"`
	Message    string   `xml:"message"`
	MessageKey string   `xml:"messageKey"`
	XMLName    struct{} `xml:"response"`
}
