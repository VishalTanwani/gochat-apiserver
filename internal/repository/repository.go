package repository

import (
	"github.com/VishalTanwani/gochat-apiserver/internal/models"
)

//DatabaseRepo interface will hold all db functions
type DatabaseRepo interface {
	RegisterUser(user models.User) (string, error)
	GetUserByID(id string) (models.User, error)
	GetUserByEmail(email string) (models.User, error)
	CheckUserAvaiability(email string) error
	UpdateUser(u models.User) (string, error)
	CreateRoom(room models.Room) (string, error)
	GetRoomByID(id string) (models.Room, error)
	GetRoomByName(name string) ([]models.Room, error)
	CheckRoomAvaiability(name string) error
	UpdateRoom(room models.Room) (string, error)
	GetUserRooms(email string) ([]models.Room, error)
	SendMessage(message models.MessageWithToken) (string, error)
	GetMessagesByRoom(room string) ([]models.Message, error)
	CreateStory(id string, userStory models.UserStory) (string, error)
	GetStory(id string) (models.UserStory, error)
	GetLastMeessage(id string) (models.Message, error)
	SetOTP(user models.UserRegister) (string, error)
	ValidateOTP(user models.UserRegister) (bool, error)
}
