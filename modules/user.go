package test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username string `bson:"Username"`
	Password string `bson:"Password"`
}

// Constructor for new user
func NewUser(username, password string) *User {
	return &User{Username: username, Password: password}
}

type UserController struct {
	Users []*User //Array of all the users
}

func NewUserController() *UserController {
	var users = make([]*User, 0) //Create empty slice to hold users
	return &UserController{Users: users}
}

func (users *UserController) Register(user User) (*User, error) {
	for _, x := range users.Users { //Check if  similar username exist
		if user.Username == x.Username {
			return nil, fmt.Errorf("username %v exists", x.Username)
		}
	}
	p, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14) // encrypt password
	if err != nil {
		return nil, err
	}
	user.Password = string(p)
	users.Users = append(users.Users, &user) //add created user to slice
	return &user, nil
}

func (users *UserController) Login(user User) (*User, error) {
	for _, x := range users.Users {
		if x.Username == user.Username { // first get the user
			err := bcrypt.CompareHashAndPassword([]byte(x.Password), []byte(user.Password)) //Compare the hash
			if err != nil {
				return nil, fmt.Errorf("Wrong credentials")
			}
			return x, nil
		}
	}
	return nil, fmt.Errorf("user %v not found", user.Username)
}

func (users *UserController) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		switch r.URL.Path {
		case "/login":
			{
				w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:8001")
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				w.Header().Set("Access-Control-Expose-Headers", "*")
				u := User{}
				err := json.NewDecoder(r.Body).Decode(&u)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				user, err := users.Login(u)
				if err != nil {
					res := struct{ Error string }{Error: err.Error()}
					json.NewEncoder(w).Encode(res)
					return
				}
				http.SetCookie(w, &http.Cookie{
					Name:     "username",
					Value:    u.Username,
					Expires:  time.Now().Add(time.Minute * 30),
					Path:     "/",
					SameSite: http.SameSiteStrictMode,
				})
				json.NewEncoder(w).Encode(user)
				return
			}
		case "/register":
			{
				u := User{}
				err := json.NewDecoder(r.Body).Decode(&u)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				user, err := users.Register(u)
				if err != nil {
					res := struct{ Error string }{Error: err.Error()}
					json.NewEncoder(w).Encode(res)
					return
				}
				json.NewEncoder(w).Encode(user)
				return
			}
		}
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
