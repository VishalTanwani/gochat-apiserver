package main

import (
	"context"
	"fmt"
	"os"
	"flag"
	"github.com/VishalTanwani/gochat-apiserver/internal/config"
	"github.com/VishalTanwani/gochat-apiserver/internal/driver"
	"github.com/VishalTanwani/gochat-apiserver/internal/handler"
	"github.com/VishalTanwani/gochat-apiserver/internal/models"
	"log"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/xhit/go-simple-mail/v2"
	"time"
	"net/http"
)

var app config.AppConfig

var emailID,emailPass *string

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}
	fmt.Println("api server")
	fmt.Println("server is running at", ":"+port)
	db, err := run()
	if err != nil {
		log.Println("error at run in main", err)
		return
	}

	defer func() {
		if err := db.Mongo.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	defer close(app.MailChan)

	go listenForMail()

	server := &http.Server{
		Addr:    ":"+port,
		Handler: routes(),
	}

	err = server.ListenAndServe()
	if err != nil {
		fmt.Println("error at running server", err)
	}
}

func run() (*driver.DB, error) {

	mailChan := make(chan models.MailData)
	app.MailChan = mailChan

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	//connect to database
	fmt.Println("Connecting to database...")
	emailID = flag.String("email","gochat34@gmail.com","email password")
	emailPass = flag.String("emailpass","eganvzpnpzengtej","email password")
	name := flag.String("dbname","vishal","data base name")
	pass := flag.String("dbpass","0109","data base password")
	flag.Parse()

	connectionString := fmt.Sprintf("mongodb+srv://%s:%s@gochat.gcc8h.mongodb.net/myFirstDatabase?retryWrites=true&w=majority", *name, *pass)
	db, err := driver.ConnectMongo(connectionString)
	if err != nil {
		log.Fatal("cannot connect to database ", err)
		return nil, err
	}
	handler.NewRepo(&app, db)
	return db, nil

}

func routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
	}))

	mux.Post("/user/register", handler.Repo.RegisterUser)
	mux.Post("/user/profile", handler.Repo.GetUserProfile)
	mux.Post("/user/rooms", handler.Repo.UserRooms)
	mux.Post("/user/get", handler.Repo.GetUserByID)
	mux.Post("/user/update", handler.Repo.UpdateUser)
	mux.Post("/room/create", handler.Repo.CreateRoom)
	mux.Post("/room/details", handler.Repo.RoomDetails)
	mux.Post("/room/join", handler.Repo.JoinRoom)
	mux.Post("/room/search", handler.Repo.SearchRoom)
	mux.Post("/room/update", handler.Repo.UpdateRoom)
	mux.Post("/room/leave", handler.Repo.LeaveRoom)
	mux.Post("/message/send", handler.Repo.SendMessage)
	mux.Post("/message/get", handler.Repo.GetMessagesByRoom)
	mux.Post("/message/getLastMessage", handler.Repo.GetLastMessagesOfRoom)
	mux.Post("/story/create", handler.Repo.CreateStoryForUser)
	mux.Post("/story/get", handler.Repo.GetStoryForUser)
	return mux

}

func listenForMail(){
	for {
		msg := <- app.MailChan
		sendMail(msg)
	}
}

func sendMail(m models.MailData) {
	server := mail.NewSMTPClient()
	server.Host = "smtp.gmail.com"
	server.Port = 587
	server.Username = *emailID
	server.Password = *emailPass
	server.Encryption = mail.EncryptionSTARTTLS
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	client,err := server.Connect()
	if err!=nil {
		fmt.Println("error at connecting mail server",err)
	}

	email := mail.NewMSG()
	email.SetFrom(m.From).AddTo(m.To).SetSubject(m.Subject)
	email.SetBody(mail.TextHTML, m.Content)
	
	err = email.Send(client)
	if err != nil {
		fmt.Println("error at sending email",err)
	} else {
		fmt.Println("MailSend")
	}
}
