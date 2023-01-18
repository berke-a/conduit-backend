package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)


type User struct{
	Username string `json:"username"`
	Email string `json:"email"`
	Password string `json:"password"`
	Bio string `json:"bio"`
	Image string `json:"image"`
}

type Profile struct{
	Username string `json:"username"`
	Bio string `json:"bio"`
	Image string `json:"image"`
	Following bool `json:"following"`
}

type Article struct{
	Title string `json:"title"`
	Description string `json:"description"`
	Body string `json:"body"`
	TagList []string `json:"tagList"`
	Comments []Comment 
	Slug string `json:"slug"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	Favorited bool `json:"favorited"`
	FavoritesCount int `json:"favoritesCount"`
	Author *Profile `json:"author"`
}

type Comment struct{
	Id int `json:"id"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	Body string `json:"body"`
	Author *Profile `json:"author"`
	Slug string
}

type MultipleComments struct{
	Comments []Comment `json:"comments"`
}

var profiles []Profile

var comments[]Comment
var tags []string
var articles []Article


func getProfile(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for _, profile := range profiles {
		if profile.Username == params["username"] {
			json.NewEncoder(w).Encode(profile)
			return
		}
	}
}

func getArticles(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type","application/json")
	json.NewEncoder(w).Encode(articles)
}

func getArticle(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	for _,article := range articles{
		if article.Slug == params["slug"]{
			json.NewEncoder(w).Encode(article)
			return
		}
	}
}

func getTags(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type","application/json")
	json.NewEncoder(w).Encode(tags)
}

func getComments(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	
	params := mux.Vars(r)
	
	for _, article := range articles{
		if article.Slug == params["slug"]{

			var coms MultipleComments

			coms.Comments = article.Comments

			json.NewEncoder(w).Encode(coms)

		}
	}
}



func main(){
	r := mux.NewRouter()

	profiles = append(profiles, Profile{Username: "berke", Bio:"I am a student.", Image:"", Following: false})
	profiles = append(profiles, Profile{Username: "keskul", Bio:"I am a cat.", Image:"", Following: false})

	articles = append(articles, Article{Slug: "how-to-train-your-dragon", Title: "How to train your dragon", Description: "Ever wonder how?", Body: "It takes a Jacobian", TagList: []string{"dragons", "training"}, CreatedAt: "2016-02-18T03:22:56.637Z", UpdatedAt: "2016-02-18T03:48:35.824Z", Favorited: false, FavoritesCount: 0, Author: &profiles[0]})

	comments = append(comments, Comment{Id: 1, CreatedAt: "2016-02-18T03:22:56.637Z", UpdatedAt: "2016-02-18T03:48:35.824Z", Body:"Nice post, thanks!", Author: &profiles[1], Slug: "how-to-train-your-dragon"})

	articles[0].Comments = append(articles[0].Comments, Comment{Id: 1, CreatedAt: "2016-02-18T03:22:56.637Z", UpdatedAt: "2016-02-18T03:48:35.824Z", Body:"Nice post, thanks!", Author: &profiles[1], Slug: "how-to-train-your-dragon"})


	tags = append(tags, "dragons")
	tags = append(tags, "training")

	r.HandleFunc("/api/profiles/{username}", getProfile).Methods("GET")
	r.HandleFunc("/api/articles", getArticles).Methods("GET")
	r.HandleFunc("/api/articles/{slug}", getArticle).Methods("GET")
	r.HandleFunc("/api/articles/{slug}/comments", getComments).Methods("GET")
	r.HandleFunc("/api/tags", getTags).Methods("GET")

	fmt.Println("Starting server at port 8000")
	log.Fatal(http.ListenAndServe(":8000",r))
}