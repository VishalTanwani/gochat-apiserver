package dbrepo

import (
	"github.com/VishalTanwani/gochat-apiserver/internal/config"
	"github.com/VishalTanwani/gochat-apiserver/internal/repository"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoDBRepo struct {
	App *config.AppConfig
	DB  *mongo.Client
}

//NewMongoRepo will return mongo DB
func NewMongoRepo(conn *mongo.Client, a *config.AppConfig) repository.DatabaseRepo {
	return &mongoDBRepo{
		App: a,
		DB:  conn,
	}
}
