package config

import (
	"log"
	"github.com/VishalTanwani/gochat-apiserver/internal/models"
)

//AppConfig hold the application config
type AppConfig struct {
	InfoLog  *log.Logger
	ErrorLog *log.Logger
	MailChan      chan models.MailData
}
