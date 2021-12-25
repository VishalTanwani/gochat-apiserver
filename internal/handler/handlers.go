package handler

import (
	"strconv"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/VishalTanwani/gochat-apiserver/internal/config"
	"github.com/VishalTanwani/gochat-apiserver/internal/driver"
	"github.com/VishalTanwani/gochat-apiserver/internal/helpers"
	"github.com/VishalTanwani/gochat-apiserver/internal/models"
	"github.com/VishalTanwani/gochat-apiserver/internal/repository"
	"github.com/VishalTanwani/gochat-apiserver/internal/repository/dbrepo"
	"github.com/dgrijalva/jwt-go"
	"math/rand"
	"net/http"
	"sort"
	"strings"
	"time"
)

//Repository is repository type
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

//Repo used by the handlers
var Repo *Repository

var key = []byte("gochatjwttoken")

//NewRepo creates new Repository
func NewRepo(a *config.AppConfig, db *driver.DB) {
	Repo = &Repository{
		App: a,
		DB:  dbrepo.NewMongoRepo(db.Mongo, a),
	}
}

//RegisterUser will register the user in our data base
func (m *Repository) RegisterUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	// w.Header().Set("Access-Control-Allow-Origin", "*")
	var temp models.UserRegister
	err := json.NewDecoder(r.Body).Decode(&temp)
	if err != nil {
		m.App.ErrorLog.Println("error at decoding body")
		helpers.ServerError(w, err)
		return
	}
	check := govalidator.IsEmail(temp.Email)
	if !check {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"meesage": "email is not valid" }`))
		return
	}
	if temp.Code == "" {
		rand.Seed(time.Now().UnixNano())
		temp.Code = strconv.Itoa(rand.Intn(1000000)) 
		message,err := m.DB.SetOTP(temp)
		if err != nil {
			m.App.ErrorLog.Println("error at setting otp for user")
			helpers.ServerError(w, err)
			return
		}
		htmlMessage := fmt.Sprintf(`
		<html>
			<head>
				<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
				<title>OTP for login in gochat</title>
			</head>
			<body>
				dear user <br/> 
				this is your otp <b>%s</b><br/>
				this is valid for 10 minutes<br/>
				do not share this with any one
			</body>
		</html>`,temp.Code)
		msg := models.MailData{
			From: "gochat34@gmail.com",
			To: temp.Email,
			Subject: "OTP for login in gochat",
			Content: htmlMessage,
		}
		m.App.MailChan <- msg
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message" : "`+message+`"}`))
		return
	} else {
		check,err := m.DB.ValidateOTP(temp)
		if err!=nil {
			m.App.ErrorLog.Println("error at validating otp for user")
			helpers.ServerError(w, err)
			return
		}
		if !check {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"message" : "otp is invalid"}`))
			return
		}
	}
	u, err := m.DB.GetUserByEmail(temp.Email)
	if err != nil {
		var user models.User
		user.Email = temp.Email
		user.CreatedAt = time.Now().Unix()
		user.UpdatedAt = time.Now().Unix()
		user.Status = "online"
		user.About = "Hey there i am using gochat which is rip off of whatsapp"
		user.LastLogin = append(user.LastLogin, time.Now().Unix())
		rand.Seed(time.Now().UnixNano())
		user.ProfileImage = fmt.Sprintf("https://avatars.dicebear.com/api/avataaars/%v.svg", rand.Intn(1000))
		user.Name = strings.Split(user.Email, "@")[0]
		userID, err := m.DB.RegisterUser(user)
		if err != nil {
			m.App.ErrorLog.Println("error at registering user")
			helpers.ServerError(w, err)
			return
		}
		u, err := m.DB.GetUserByID(userID)
		if err != nil {
			m.App.ErrorLog.Println("error at geting user")
			helpers.ServerError(w, err)
			return
		}
		u.Token, err = generateJWTToken(u)
		if err != nil {
			m.App.ErrorLog.Println("error at generating token")
			helpers.ServerError(w, err)
			return
		}
		json.NewEncoder(w).Encode(u)
	} else {
		u.LastLogin = append(u.LastLogin, time.Now().Unix())
		u.UpdatedAt = time.Now().Unix()
		_, err := m.DB.UpdateUser(u)
		if err != nil {
			m.App.ErrorLog.Println("error at registering user")
			helpers.ServerError(w, err)
			return
		}
		u, err := m.DB.GetUserByID(primtiveObjToString(u.ID))
		if err != nil {
			m.App.ErrorLog.Println("error at geting user")
			helpers.ServerError(w, err)
			return
		}
		u.Token, err = generateJWTToken(u)
		if err != nil {
			m.App.ErrorLog.Println("error at generating token")
			helpers.ServerError(w, err)
			return
		}
		json.NewEncoder(w).Encode(u)
	}
}

//CreateRoom will creat a room in our data base
func (m *Repository) CreateRoom(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	var temp models.RoomWithToken
	err := json.NewDecoder(r.Body).Decode(&temp)
	if err != nil {
		m.App.ErrorLog.Println("error at decoding body")
		helpers.ServerError(w, err)
		return
	}
	check, err := verifyToken(temp.Token)
	if err != nil {
		m.App.ErrorLog.Println("error at verifing token")
		helpers.ServerError(w, err)
		return
	}
	if check {
		err = m.DB.CheckRoomAvaiability(temp.Name)
		if err != nil {
			mapData, err := tokenDecode(temp.Token)
			if err != nil {
				m.App.ErrorLog.Println("error at decoding token")
				helpers.ServerError(w, err)
				return
			}
			var room models.Room
			room.Name = temp.Name
			room.Type = temp.Type
			room.CreatedAt = time.Now().Unix()
			room.UpdatedAt = time.Now().Unix()
			room.GroupIcon = fmt.Sprintf("https://avatars.dicebear.com/api/avataaars/%s.svg", temp.Name)
			room.Description = "Description ..."
			room.CreatedBy = fmt.Sprint(mapData["email"])
			room.Users = append(room.Users, fmt.Sprint(mapData["email"]))

			roomID, err := m.DB.CreateRoom(room)
			if err != nil {
				m.App.ErrorLog.Println("error at creating room")
				helpers.ServerError(w, err)
				return
			}

			room, err = m.DB.GetRoomByID(roomID)
			if err != nil {
				m.App.ErrorLog.Println("error at getting room")
				helpers.ServerError(w, err)
				return
			}
			json.NewEncoder(w).Encode(room)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"meesage": "room is already created" }`))
			return
		}
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"meesage": "token invalidation" }`))
		return
	}
}

//JoinRoom will join a user to room
func (m *Repository) JoinRoom(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	var temp models.RoomWithToken
	err := json.NewDecoder(r.Body).Decode(&temp)
	if err != nil {
		m.App.ErrorLog.Println("error at decoding body")
		helpers.ServerError(w, err)
		return
	}
	check, err := verifyToken(temp.Token)
	if err != nil {
		m.App.ErrorLog.Println("error at verifing token")
		helpers.ServerError(w, err)
		return
	}
	if check {
		room, err := m.DB.GetRoomByID(primtiveObjToString(temp.ID))
		if err == nil {
			mapData, err := tokenDecode(temp.Token)
			if err != nil {
				m.App.ErrorLog.Println("error at decoding token")
				helpers.ServerError(w, err)
				return
			}
			room.UpdatedAt = time.Now().Unix()
			for _, v := range room.Users {
				if v == fmt.Sprint(mapData["email"]) {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"meesage": "already joined" }`))
					return
				}
			}
			room.Users = append(room.Users, fmt.Sprint(mapData["email"]))

			_, err = m.DB.UpdateRoom(room)
			if err != nil {
				m.App.ErrorLog.Println("error at updateing room")
				helpers.ServerError(w, err)
				return
			}

			room, err = m.DB.GetRoomByID(primtiveObjToString(room.ID))
			if err != nil {
				m.App.ErrorLog.Println("error at getting room")
				helpers.ServerError(w, err)
				return
			}
			json.NewEncoder(w).Encode(room)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"meesage": "cannot find room" }`))
			return
		}
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"meesage": "token invalidation" }`))
		return
	}

}

//SearchRoom will search a room by name
func (m *Repository) SearchRoom(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	var temp models.RoomWithToken
	err := json.NewDecoder(r.Body).Decode(&temp)
	if err != nil {
		m.App.ErrorLog.Println("error at decoding body")
		helpers.ServerError(w, err)
		return
	}
	check, err := verifyToken(temp.Token)
	if err != nil {
		m.App.ErrorLog.Println("error at verifing token")
		helpers.ServerError(w, err)
		return
	}
	if check {
		rooms, err := m.DB.GetRoomByName(temp.Name)
		if err == nil {
			json.NewEncoder(w).Encode(rooms)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"meesage": "cannot find room" }`))
			return
		}
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"meesage": "token invalidation" }`))
		return
	}

}

//LeaveRoom will leave a user to room
func (m *Repository) LeaveRoom(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	var temp models.RoomWithToken
	err := json.NewDecoder(r.Body).Decode(&temp)
	if err != nil {
		m.App.ErrorLog.Println("error at decoding body")
		helpers.ServerError(w, err)
		return
	}
	check, err := verifyToken(temp.Token)
	if err != nil {
		m.App.ErrorLog.Println("error at verifing token")
		helpers.ServerError(w, err)
		return
	}
	if check {
		room, err := m.DB.GetRoomByID(primtiveObjToString(temp.ID))
		if err == nil {
			mapData, err := tokenDecode(temp.Token)
			if err != nil {
				m.App.ErrorLog.Println("error at decoding token")
				helpers.ServerError(w, err)
				return
			}
			room.UpdatedAt = time.Now().Unix()
			for i, v := range room.Users {
				if v == fmt.Sprint(mapData["email"]) {
					room.Users = append(room.Users[:i], room.Users[i+1:]...)
					_, err = m.DB.UpdateRoom(room)
					if err != nil {
						m.App.ErrorLog.Println("error at updateing room")
						helpers.ServerError(w, err)
						return
					}

					room, err = m.DB.GetRoomByID(primtiveObjToString(room.ID))
					if err != nil {
						m.App.ErrorLog.Println("error at getting room")
						helpers.ServerError(w, err)
						return
					}
					json.NewEncoder(w).Encode(room)
					return
				}
			}
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"meesage": "You are not in a room" }`))
			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"meesage": "cannot find room" }`))
			return
		}
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"meesage": "token invalidation" }`))
		return
	}

}

//UpdateRoom will join a update a room
func (m *Repository) UpdateRoom(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	var temp models.RoomWithToken
	err := json.NewDecoder(r.Body).Decode(&temp)
	if err != nil {
		m.App.ErrorLog.Println("error at decoding body")
		helpers.ServerError(w, err)
		return
	}
	check, err := verifyToken(temp.Token)
	if err != nil {
		m.App.ErrorLog.Println("error at verifing token")
		helpers.ServerError(w, err)
		return
	}
	if check {
		room, err := m.DB.GetRoomByID(primtiveObjToString(temp.ID))
		if err == nil {
			mapData, err := tokenDecode(temp.Token)
			if err != nil {
				m.App.ErrorLog.Println("error at decoding token")
				helpers.ServerError(w, err)
				return
			}
			for _, v := range room.Users {
				if v == fmt.Sprint(mapData["email"]) {
					if room.Name != temp.Name {
						room.Name = temp.Name
					}
					if room.Description != temp.Description {
						room.Description = temp.Description
					}
					room.UpdatedAt = time.Now().Unix()

					_, err = m.DB.UpdateRoom(room)
					if err != nil {
						m.App.ErrorLog.Println("error at updateing room")
						helpers.ServerError(w, err)
						return
					}

					room, err = m.DB.GetRoomByID(primtiveObjToString(room.ID))
					if err != nil {
						m.App.ErrorLog.Println("error at getting room")
						helpers.ServerError(w, err)
						return
					}
					json.NewEncoder(w).Encode(room)
					return
				}
			}
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"meesage": "you cannot update a room" }`))
			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"meesage": "cannot find room" }`))
			return
		}
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"meesage": "token invalidation" }`))
		return
	}
}

//UpdateUser will join a update a user
func (m *Repository) UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	var temp models.User
	err := json.NewDecoder(r.Body).Decode(&temp)
	if err != nil {
		m.App.ErrorLog.Println("error at decoding body")
		helpers.ServerError(w, err)
		return
	}
	check, err := verifyToken(temp.Token)
	if err != nil {
		m.App.ErrorLog.Println("error at verifing token")
		helpers.ServerError(w, err)
		return
	}
	if check {
		mapData, err := tokenDecode(temp.Token)
		if err != nil {
			m.App.ErrorLog.Println("error at decodeing token")
			helpers.ServerError(w, err)
			return
		}
		user, err := m.DB.GetUserByID(fmt.Sprint(mapData["_id"]))
		if err == nil {
			user.Name = temp.Name
			user.About = temp.About
			user.UpdatedAt = time.Now().Unix()

			_, err = m.DB.UpdateUser(user)
			if err != nil {
				m.App.ErrorLog.Println("error at updateing room")
				helpers.ServerError(w, err)
				return
			}

			user, err = m.DB.GetUserByID(primtiveObjToString(user.ID))
			if err != nil {
				m.App.ErrorLog.Println("error at getting user")
				helpers.ServerError(w, err)
				return
			}
			json.NewEncoder(w).Encode(user)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"meesage": "cannot find user" }`))
			return
		}
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"meesage": "token invalidation" }`))
		return
	}
}

//GetUserProfile will give user profile
func (m *Repository) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	var temp models.User
	err := json.NewDecoder(r.Body).Decode(&temp)
	if err != nil {
		m.App.ErrorLog.Println("error at decoding body")
		helpers.ServerError(w, err)
		return
	}
	check, err := verifyToken(temp.Token)
	if err != nil {
		m.App.ErrorLog.Println("error at verifing token")
		helpers.ServerError(w, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"meesage": "token invalidation" }`))
		return
	}
	if check {
		mapData, err := tokenDecode(temp.Token)
		if err != nil {
			m.App.ErrorLog.Println("error at decodeing token")
			helpers.ServerError(w, err)
			return
		}
		user, err := m.DB.GetUserByID(fmt.Sprint(mapData["_id"]))
		if err == nil {
			json.NewEncoder(w).Encode(user)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"meesage": "cannot find user" }`))
			return
		}
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"meesage": "token invalidation" }`))
		return
	}
}

//GetUserByID will give user profile
func (m *Repository) GetUserByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	var temp models.User
	err := json.NewDecoder(r.Body).Decode(&temp)
	if err != nil {
		m.App.ErrorLog.Println("error at decoding body")
		helpers.ServerError(w, err)
		return
	}

	user, err := m.DB.GetUserByID(primtiveObjToString(temp.ID))
	if err == nil {
		json.NewEncoder(w).Encode(user)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"meesage": "cannot find user" }`))
		return
	}
}

//UserRooms will give users room
func (m *Repository) UserRooms(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	var temp models.User
	err := json.NewDecoder(r.Body).Decode(&temp)
	if err != nil {
		m.App.ErrorLog.Println("error at decoding body")
		helpers.ServerError(w, err)
		return
	}
	check, err := verifyToken(temp.Token)
	if err != nil {
		m.App.ErrorLog.Println("error at verifing token")
		helpers.ServerError(w, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"meesage": "token invalidation" }`))
		return
	}
	if check {
		mapData, err := tokenDecode(temp.Token)
		if err != nil {
			m.App.ErrorLog.Println("error at decodeing token")
			helpers.ServerError(w, err)
			return
		}
		rooms, err := m.DB.GetUserRooms(fmt.Sprint(mapData["email"]))
		if err == nil {
			json.NewEncoder(w).Encode(rooms)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"meesage": "cannot find user" }`))
			return
		}
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"meesage": "token invalidation" }`))
		return
	}
}

//RoomDetails will give room details
func (m *Repository) RoomDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	var temp models.RoomWithToken
	err := json.NewDecoder(r.Body).Decode(&temp)
	if err != nil {
		m.App.ErrorLog.Println("error at decoding body")
		helpers.ServerError(w, err)
		return
	}
	check, err := verifyToken(temp.Token)
	if err != nil {
		m.App.ErrorLog.Println("error at verifing token")
		helpers.ServerError(w, err)
		return
	}
	if check {
		room, err := m.DB.GetRoomByID(primtiveObjToString(temp.ID))
		if err == nil {
			json.NewEncoder(w).Encode(room)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"meesage": "cannot find room" }`))
			return
		}
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"meesage": "token invalidation" }`))
		return
	}

}

//SendMessage will store message in DB
func (m *Repository) SendMessage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	var temp models.MessageWithToken
	err := json.NewDecoder(r.Body).Decode(&temp)
	if err != nil {
		m.App.ErrorLog.Println("error at decoding body")
		helpers.ServerError(w, err)
		return
	}
	check, err := verifyToken(temp.Token)
	if err != nil {
		m.App.ErrorLog.Println("error at verifing token")
		helpers.ServerError(w, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"meesage": "token invalidation" }`))
		return
	}
	if check {
		room, err := m.DB.GetRoomByID(primtiveObjToString(temp.RoomID))
		if err == nil {
			mapData, err := tokenDecode(temp.Token)
			if err != nil {
				m.App.ErrorLog.Println("error at decoding token")
				helpers.ServerError(w, err)
				return
			}
			if temp.Type == "joinRoom" {
				for _, v := range room.Users {
					if v == fmt.Sprint(mapData["email"]) {
						return
					}
				}
			}
			temp.Token = ""
			temp.CreatedAt = time.Now().Unix()
			res, err := m.DB.SendMessage(temp)
			if err == nil {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"meesage": ` + res + `}`))
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"meesage": "cannot send room" }`))
				return
			}
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"meesage": "cannot find room" }`))
			return
		}

	} else {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"meesage": "token invalidation" }`))
		return
	}
}

//GetMessagesByRoom will give messages of room
func (m *Repository) GetMessagesByRoom(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	var temp models.RoomWithToken
	err := json.NewDecoder(r.Body).Decode(&temp)
	if err != nil {
		m.App.ErrorLog.Println("error at decoding body")
		helpers.ServerError(w, err)
		return
	}
	check, err := verifyToken(temp.Token)
	if err != nil {
		m.App.ErrorLog.Println("error at verifing token")
		helpers.ServerError(w, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"meesage": "token invalidation" }`))
		return
	}
	if check {
		mapData, err := tokenDecode(temp.Token)
		if err != nil {
			m.App.ErrorLog.Println("error at decoding token")
			helpers.ServerError(w, err)
			return
		}
		room, err := m.DB.GetRoomByID(primtiveObjToString(temp.ID))
		if err == nil {
			for _, v := range room.Users {
				if v == fmt.Sprint(mapData["email"]) {
					messages, err := m.DB.GetMessagesByRoom(primtiveObjToString(temp.ID))
					if err == nil {
						sort.Slice(messages, func(i, j int) bool {
							return messages[i].CreatedAt < messages[j].CreatedAt
						})
						json.NewEncoder(w).Encode(messages)
						return
					} else {
						w.WriteHeader(http.StatusInternalServerError)
						w.Write([]byte(`{"meesage": "some error occured" }`))
						return
					}
				}
			}
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"meesage": "You are not in a room" }`))
			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"meesage": "cannot find room" }`))
			return
		}

	} else {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"meesage": "token invalidation" }`))
		return
	}
}

//GetLastMessagesOfRoom will give last messages of room
func (m *Repository) GetLastMessagesOfRoom(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	var temp models.MessageWithToken
	err := json.NewDecoder(r.Body).Decode(&temp)
	if err != nil {
		m.App.ErrorLog.Println("error at decoding body")
		helpers.ServerError(w, err)
		return
	}
	check, err := verifyToken(temp.Token)
	if err != nil {
		m.App.ErrorLog.Println("error at verifing token")
		helpers.ServerError(w, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"meesage": "token invalidation" }`))
		return
	}
	if check {
		message, err := m.DB.GetLastMeessage(primtiveObjToString(temp.RoomID))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"meesage": "cannot get last message" }`))
			return
		}
		json.NewEncoder(w).Encode(message)

	} else {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"meesage": "token invalidation" }`))
		return
	}
}

//CreateStoryForUser will create a story for a user
func (m *Repository) CreateStoryForUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	var temp models.UserStory
	err := json.NewDecoder(r.Body).Decode(&temp)
	if err != nil {
		m.App.ErrorLog.Println("error at decoding body")
		helpers.ServerError(w, err)
		return
	}
	check, err := verifyToken(temp.Token)
	if err != nil {
		m.App.ErrorLog.Println("error at verifing token")
		helpers.ServerError(w, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"meesage": "token invalidation" }`))
		return
	}
	if check {
		mapData, err := tokenDecode(temp.Token)
		if err != nil {
			m.App.ErrorLog.Println("error at decodeing token")
			helpers.ServerError(w, err)
			return
		}
		_,err = m.DB.GetStory(fmt.Sprint(mapData["_id"]))
		if err!=nil {
			temp.Token = ""
			res, err := m.DB.CreateStory(fmt.Sprint(mapData["_id"]), temp)
			if err == nil {
				json.NewEncoder(w).Encode(res)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"meesage": "can not set story" }`))
				return
			}
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"meesage": "story is already there" }`))
			return
		}
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"meesage": "token invalidation" }`))
		return
	}
}

//GetStoryForUser will give a user
func (m *Repository) GetStoryForUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	var temp models.UserStory
	err := json.NewDecoder(r.Body).Decode(&temp)
	if err != nil {
		m.App.ErrorLog.Println("error at decoding body")
		helpers.ServerError(w, err)
		return
	}
	check, err := verifyToken(temp.Token)
	if err != nil {
		m.App.ErrorLog.Println("error at verifing token")
		helpers.ServerError(w, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"meesage": "token invalidation" }`))
		return
	}
	if check {
		mapData, err := tokenDecode(temp.Token)
		if err != nil {
			m.App.ErrorLog.Println("error at decodeing token")
			helpers.ServerError(w, err)
			return
		}
		res, err := m.DB.GetStory(fmt.Sprint(mapData["_id"]))
		if err == nil {
			json.NewEncoder(w).Encode(res)
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"meesage": "can not get story" }`))
			return
		}
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"meesage": "token invalidation" }`))
		return
	}
}

func generateJWTToken(user models.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claim := token.Claims.(jwt.MapClaims)
	claim["_id"] = user.ID
	claim["name"] = user.Name
	claim["email"] = user.Email
	claim["exp"] = time.Now().Add(time.Hour * 24 * 30).Unix()
	tokenString, err := token.SignedString(key)
	return tokenString, err
}

func verifyToken(tokenString string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("there was an error")
		}
		return key, nil
	})

	if err != nil {
		return false, err
	}

	if token.Valid {
		return true, nil
	}
	return false, errors.New("User not Found")
}

func tokenDecode(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("there was an error")
		}
		return key, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}

func primtiveObjToString(id interface{}) string {
	ID := fmt.Sprintf("%s", id)
	return strings.Split(ID, "\"")[1]
}
