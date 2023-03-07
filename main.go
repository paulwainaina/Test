package main

import (
	"fmt"
	"net/http"
	"os"
	"example.com/test"
)

func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {	
		if r.URL.Path == "/posts" {
			_, err := r.Cookie("username")
			if err != nil {
				http.Redirect(w, r, "/loginPage", http.StatusMovedPermanently)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	pos:=test.NewPostController()
	http.Handle("/posts", middleware(http.HandlerFunc(pos.ServeHTTP)))
	user:=test.NewUserController()
	http.Handle("/register",http.HandlerFunc(user.ServeHTTP))
	http.Handle("/login",http.HandlerFunc(user.ServeHTTP))
	err := http.ListenAndServe("127.0.0.1:8001", nil)
	if err == http.ErrServerClosed {
		fmt.Println("Backend server closed")
	} else if err != nil {
		fmt.Println("Backend server:Error occured " + err.Error())
		os.Exit(1)
	}
}