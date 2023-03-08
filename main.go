package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"

	"example.com/test"
)
var tpl *template.Template
func init() {
	tpl = template.Must(template.ParseGlob("./templates/*.html"))
}

type Page struct {
	Body  []byte
	Title string
	Data  interface{}
	Error error
}

func LoadPage(file string) (*Page, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	var body []byte
	_, err = f.Read(body)
	if err != nil {
		return nil, err
	}
	return &Page{Body: body}, nil
}

func RenderTemplate(w http.ResponseWriter, file string, page *Page) {
	err := tpl.ExecuteTemplate(w, file, page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {	
		w.Header().Set("Cache-Control", "no-cache")
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:8001")
			w.Header().Set("Access-Control-Allow-Methods", "POST,GET")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			return
		}
		if r.URL.Path == "/posts" || r.URL.Path=="/index" {
			_, err := r.Cookie("username")
			if err != nil {
				http.Redirect(w, r, "/loginPage", http.StatusMovedPermanently)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	file := "index.html"
	filePath := "templates/" + file
	pageName := "Home Page"
	page, err := LoadPage(filePath)
	if err != nil {
		page = &Page{Title: pageName}
	}
	page.Title = pageName
	RenderTemplate(w, file, page)
}
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	file := "login.html"
	filePath := "templates/" + file
	pageName := "Login Page"
	page, err := LoadPage(filePath)
	if err != nil {
		page = &Page{Title: pageName}
	}
	page.Title = pageName
	RenderTemplate(w, file, page)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	file := "register.html"
	filePath := "templates/" + file
	pageName := "Register Page"
	page, err := LoadPage(filePath)
	if err != nil {
		page = &Page{Title: pageName}
	}
	page.Title = pageName
	RenderTemplate(w, file, page)
}


func main() {
	pos:=test.NewPostController()
	http.Handle("/posts", middleware(http.HandlerFunc(pos.ServeHTTP)))
	user:=test.NewUserController()
	http.Handle("/register",middleware(http.HandlerFunc(user.ServeHTTP)))
	http.Handle("/login",middleware(http.HandlerFunc(user.ServeHTTP)))


	http.Handle("/loginPage",middleware(http.HandlerFunc(LoginHandler)))
	http.Handle("/index",middleware(http.HandlerFunc(IndexHandler)))
	http.Handle("/registerPage",middleware(http.HandlerFunc(RegisterHandler)))

	err := http.ListenAndServe("127.0.0.1:8001", nil)
	if err == http.ErrServerClosed {
		fmt.Println("Backend server closed")
	} else if err != nil {
		fmt.Println("Backend server:Error occured " + err.Error())
		os.Exit(1)
	}
}