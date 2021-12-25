package main

import (
	"github.com/VishalTanwani/gochat-apiserver/internal/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"net/http"
)

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

// package main

// import (
// 	"github.com/VishalTanwani/gochat-apiserver/internal/handler"
// 	"github.com/gorilla/handlers"
//     "github.com/gorilla/mux"
// 	"net/http"
// )

// func routes() http.Handler {
// 	router := mux.NewRouter()

// 	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
// 	originsOk := handlers.AllowedOrigins([]string{"*"})
// 	methodsOk := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"})

// 	router.HandleFunc("/user/register", handler.Repo.RegisterUser).Methods("POST")
// 	router.HandleFunc("/user/profile", handler.Repo.GetUserProfile).Methods("GET")
// 	router.HandleFunc("/user/rooms", handler.Repo.UserRooms).Methods("GET")
// 	router.HandleFunc("/user/update", handler.Repo.UpdateUser).Methods("POST")
// 	router.HandleFunc("/room/create", handler.Repo.CreateRoom).Methods("POST")
// 	router.HandleFunc("/room/join", handler.Repo.JoinRoom).Methods("POST")
// 	router.HandleFunc("/room/update", handler.Repo.UpdateRoom).Methods("POST")
// 	router.HandleFunc("/room/leave", handler.Repo.LeaveRoom).Methods("POST")

// 	return handlers.CORS(originsOk, headersOk, methodsOk)(router)

// }
