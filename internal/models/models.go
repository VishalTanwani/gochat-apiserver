package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//User is user model
type User struct {
	ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Email        string             `json:"email,omitempty" bson:"email,omitempty"`
	Name         string             `json:"name,omitempty" bson:"name,omitempty"`
	ProfileImage string             `json:"profile_image,omitempty" bson:"profile_image,omitempty"`
	Status       string             `json:"status,omitempty" bson:"status,omitempty"`
	About        string             `json:"about,omitempty" bson:"about,omitempty"`
	Token        string             `json:"token,omitempty" bson:"token,omitempty"`
	LastLogin    []int64            `json:"last_login,omitempty" bson:"last_login,omitempty"`
	CreatedAt    int64              `json:"create_at,omitempty" bson:"create_at,omitempty"`
	UpdatedAt    int64              `json:"update_at,omitempty" bson:"update_at,omitempty"`
}

//UserRegister is for login model
type UserRegister struct {
	ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Email string `json:"email,omitempty" bson:"email,omitempty"`
	Code  string `json:"code,omitempty" bson:"code,omitempty"`
	CreatedAt primitive.DateTime `json:"create_at,omitempty" bson:"create_at,omitempty"`
	ExpireOn  primitive.DateTime `json:"expire_on,omitempty" bson:"expire_on,omitempty"`
}

//Room is room model
type Room struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name,omitempty" bson:"name,omitempty"`
	Type        string             `json:"type,omitempty" bson:"type,omitempty"`
	Description string             `json:"description,omitempty" bson:"description,omitempty"`
	GroupIcon   string             `json:"group_icon,omitempty" bson:"group_icon,omitempty"`
	Users       []string           `json:"users,omitempty" bson:"users,omitempty"`
	CreatedBy   string             `json:"create_by,omitempty" bson:"create_by,omitempty"`
	CreatedAt   int64              `json:"create_at,omitempty" bson:"create_at,omitempty"`
	UpdatedAt   int64              `json:"update_at,omitempty" bson:"update_at,omitempty"`
}

//RoomWithToken is room model
type RoomWithToken struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name,omitempty" bson:"name,omitempty"`
	Type        string             `json:"type,omitempty" bson:"type,omitempty"`
	Token       string             `json:"token,omitempty" bson:"token,omitempty"`
	Description string             `json:"description,omitempty" bson:"description,omitempty"`
	GroupIcon   string             `json:"group_icon,omitempty" bson:"group_icon,omitempty"`
	Users       []string           `json:"users,omitempty" bson:"users,omitempty"`
	CreatedBy   string             `json:"create_by,omitempty" bson:"create_by,omitempty"`
	CreatedAt   int64              `json:"create_at,omitempty" bson:"create_at,omitempty"`
	UpdatedAt   int64              `json:"update_at,omitempty" bson:"update_at,omitempty"`
}

//Message is message model
type Message struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserID    primitive.ObjectID `json:"user_id,omitempty" bson:"user_id,omitempty"`
	RoomID    primitive.ObjectID `json:"room_id,omitempty" bson:"room_id,omitempty"`
	UserName  string             `json:"user_name" bson:"user_name"`
	Body      string             `json:"body" bson:"body"`
	Image     string             `json:"image" bson:"image"`
	Type      string             `json:"type,omitempty" bson:"type,omitempty"`
	Room      string             `json:"room" bson:"room"`
	CreatedAt int64              `json:"create_at,omitempty" bson:"create_at,omitempty"`
}

//MessageWithToken is message model
type MessageWithToken struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserID    primitive.ObjectID `json:"user_id,omitempty" bson:"user_id,omitempty"`
	RoomID    primitive.ObjectID `json:"room_id,omitempty" bson:"room_id,omitempty"`
	Token     string             `json:"token,omitempty" bson:"token,omitempty"`
	UserName  string             `json:"user_name,omitempty" bson:"user_name,omitempty"`
	Body      string             `json:"body,omitempty" bson:"body,omitempty"`
	Image     string             `json:"image" bson:"image"`
	Type      string             `json:"type,omitempty" bson:"type,omitempty"`
	Room      string             `json:"room,omitempty" bson:"room,omitempty"`
	CreatedAt int64              `json:"create_at,omitempty" bson:"create_at,omitempty"`
}

//UserStory is story model
type UserStory struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserID    primitive.ObjectID `json:"user_id,omitempty" bson:"user_id,omitempty"`
	Body      string             `json:"body,omitempty" bson:"body,omitempty"`
	Token     string             `json:"token,omitempty" bson:"token,omitempty"`
	CreatedAt primitive.DateTime `json:"create_at,omitempty" bson:"create_at,omitempty"`
	ExpireOn  primitive.DateTime `json:"expire_on,omitempty" bson:"expire_on,omitempty"`
}

//MailData is our mail model
type MailData struct {
	From     string
	To       string
	Subject  string
	Content  string
}
