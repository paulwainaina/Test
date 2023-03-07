package test

import (
	"encoding/json"
	"net/http"
	"sort"
	"strings"
	"time"
)

type Post struct {
	Title string `bson:"Title"`
	Body  string `bson:"Body"`
	User  string `bson:"Username"`
	Date  time.Time `bson:"Date"`
}
func (post *Post)UnmarshalJSON(b []byte) (error){
	var x map[string]interface{}
	err:=json.Unmarshal(b,&x)
	if err!=nil{
		return err
	}	
	for k,v :=range x{
		switch strings.ToLower(k){
		case "title":{
			post.Title=v.(string)
		}
		case "body":{
			post.Body=v.(string)
		}
		case "username":{
			post.User=v.(string)
		}
		}
	}
	post.Date=time.Now()
	return nil
}

type PostController struct {
	Posts []*Post //Array of all the posts
}

func NewPostController() *PostController {
	var posts = make([]*Post, 0) //Create empty slice to hold posts
	return &PostController{Posts: posts}
}

func(posts * PostController)MakePost(post Post)(*Post,error){
	posts.Posts=append(posts.Posts, &post)
	return &post,nil
}

func(posts * PostController)GetPost(user string)([]*Post,error){
	r:=make([]*Post,0)
	for _,x:=range posts.Posts{
		if x.User==user{
			r=append(r, x)		
		}
	}
	sort.Slice(r,func(i, j int) bool {
		return r[i].Date.Before(r[j].Date)
	})
	return r,nil
}

func (posts *PostController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		switch r.URL.Path {
		case "/posts":
			{
				u := Post{}
				err := json.NewDecoder(r.Body).Decode(&u)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				Post, err := posts.MakePost(u)
				if err != nil {
					w.Write([]byte(err.Error()))
					return
				}
				json.NewEncoder(w).Encode(Post)
				return
			}
		}
	} else if(r.Method==http.MethodGet){
		switch r.URL.Path {
			case "/posts":{
				cok,err:=r.Cookie("username")
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				Post, err := posts.GetPost(cok.Value)
				if err != nil {
					w.Write([]byte(err.Error()))
					return
				}
				json.NewEncoder(w).Encode(Post)
				return
			}
		}
	}else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}