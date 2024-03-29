module github.com/VishalTanwani/gochat-apiserver

// +heroku goVersion go1.15      <--add to
go 1.15

require (
	github.com/asaskevich/govalidator v0.0.0-20210307081110-f21760c49a8d
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/go-chi/chi/v5 v5.0.7
	github.com/go-chi/cors v1.2.0
	github.com/xhit/go-simple-mail/v2 v2.10.0
	go.mongodb.org/mongo-driver v1.8.1
)
