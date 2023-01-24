package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Token    string `json:"token"`
	Password string `json:"password"`
	Bio      string `json:"bio"`
	Image    string `json:"image"`
}

type Profile struct {
	Username  string `json:"username"`
	Bio       string `json:"bio"`
	Image     string `json:"image"`
	Following bool   `json:"following"`
}

type Article struct {
	Title          string   `json:"title"`
	Description    string   `json:"description"`
	Body           string   `json:"body"`
	TagList        []string `json:"tagList"`
	Comments       []Comment
	Slug           string   `json:"slug"`
	CreatedAt      string   `json:"createdAt"`
	UpdatedAt      string   `json:"updatedAt"`
	Favorited      bool     `json:"favorited"`
	FavoritesCount int      `json:"favoritesCount"`
	Author         *Profile `json:"author"`
}

type Comment struct {
	Id        int      `json:"id"`
	CreatedAt string   `json:"createdAt"`
	UpdatedAt string   `json:"updatedAt"`
	Body      string   `json:"body"`
	Author    *Profile `json:"author"`
	Slug      string
}

type MultipleComments struct {
	Comments []Comment `json:"comments"`
}

var profiles []Profile

var comments []Comment
var tags []string
var articles []Article

func getProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for _, profile := range profiles {
		if profile.Username == params["username"] {
			if json.NewEncoder(w).Encode(profile) != nil {
				log.Fatal("Error at encoding profile")
			}
			return
		}
	}
}

func getArticles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	limit := 20
	offset := 0

	queries := r.URL.Query()
	tag := queries.Get("tag")
	author := queries.Get("author")
	favorited := queries.Get("favorited")
	queryLimit := queries.Get("limit")
	queryOffset := queries.Get("offset")

	encoder := json.NewEncoder(w)

	if queryLimit != "" {
		limit, _ = strconv.Atoi(queryLimit)
	}
	if queryOffset != "" {
		offset, _ = strconv.Atoi(queryOffset)
	}

	if offset+limit >= len(articles) {
		limit = len(articles) - offset
	}

	if tag != "" {
		var newArticles []Article

		for _, article := range articles {
			for _, t := range article.TagList {
				if t == tag {
					newArticles = append(newArticles, article)
				}
			}
		}
		if encoder.Encode(newArticles[offset:offset+limit]) != nil {
			log.Fatal("Error at encoding articles with tag filter")
		}
		return
	}

	if author != "" {
		var newArticles []Article

		for _, article := range articles {
			if article.Author.Username == author {
				newArticles = append(newArticles, article)
			}
		}
		if encoder.Encode(newArticles[offset:offset+limit]) != nil {
			log.Fatal("Error at encoding articles with author filter")
		}
		return
	}

	if favorited != "" {
		var newArticles []Article

		for _, article := range articles {
			if article.Favorited {
				newArticles = append(newArticles, article)
			}
		}
		if encoder.Encode(newArticles[offset:offset+limit]) != nil {
			log.Fatal("Error at encoding articles with favorited filter")
		}
		return
	}

	if encoder.Encode(articles[offset:offset+limit]) != nil {
		log.Fatal("Error at encoding articles")
	}
}

func getArticle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	for _, article := range articles {
		if article.Slug == params["slug"] {
			if json.NewEncoder(w).Encode(article) != nil {
				log.Fatal("Error at encoding single article")
			}
			return
		}
	}
}

func getTags(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if json.NewEncoder(w).Encode(tags) != nil {
		log.Fatal("Error at encoding tags")
	}
}

func getComments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	for _, article := range articles {
		if article.Slug == params["slug"] {

			var coms MultipleComments

			coms.Comments = article.Comments

			if json.NewEncoder(w).Encode(coms) != nil {
				log.Fatal("Error at encoding comments")
			}

		}
	}
}

func login(w http.ResponseWriter, r *http.Request) {

}

func initDummyData() {
	profiles = append(profiles, Profile{Username: "berke", Bio: "I am a student.", Image: "", Following: false})
	profiles = append(profiles, Profile{Username: "keskul", Bio: "I am a cat.", Image: "", Following: false})

	articles = append(articles, Article{Slug: "how-to-train-your-dragon", Title: "How to train your dragon", Description: "Ever wonder how?", Body: "It takes a Jacobian", TagList: []string{"dragons", "training"}, CreatedAt: "2016-02-18T03:22:56.637Z", UpdatedAt: "2016-02-18T03:48:35.824Z", Favorited: false, FavoritesCount: 0, Author: &profiles[0]})

	comments = append(comments, Comment{Id: 1, CreatedAt: "2016-02-18T03:22:56.637Z", UpdatedAt: "2016-02-18T03:48:35.824Z", Body: "Nice post, thanks!", Author: &profiles[1], Slug: "how-to-train-your-dragon"})

	articles[0].Comments = append(articles[0].Comments, Comment{Id: 1, CreatedAt: "2016-02-18T03:22:56.637Z", UpdatedAt: "2016-02-18T03:48:35.824Z", Body: "Nice post, thanks!", Author: &profiles[1], Slug: "how-to-train-your-dragon"})

	tags = append(tags, "dragons")
	tags = append(tags, "training")
}

func initRoutes(r *mux.Router) {
	// GET
	r.HandleFunc("/api/profiles/{username}", getProfile).Methods("GET")
	r.HandleFunc("/api/articles", getArticles).Methods("GET")
	r.HandleFunc("/api/articles/{slug}", getArticle).Methods("GET")
	r.HandleFunc("/api/articles/{slug}/comments", getComments).Methods("GET")
	r.HandleFunc("/api/tags", getTags).Methods("GET")

	// POST
	r.HandleFunc("/api/user/login", login).Methods("POST")
}

func main() {
	r := mux.NewRouter()

	initDummyData()
	initRoutes(r)

	fmt.Println("Starting server at port 8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}
