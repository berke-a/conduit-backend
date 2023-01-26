package main

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Token    string `json:"token"`
	Password string `json:"password"`
	Bio      string `json:"bio"`
	Image    string `json:"image"`
}

type UserWrapper struct {
	User User `json:"user"`
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
var users []User

var currentUser UserWrapper

var secretKey = []byte("secretKey")

func GenerateJWTToken(username string) (string, error) {
	// Create the Claims
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["username"] = username
	claims["authorized"] = true
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	tokenString, err := token.SignedString(secretKey)

	if err != nil {
		fmt.Errorf("JWT token generation failed: %s", err.Error())
		return "", err
	}
	return tokenString, nil
}

func isAuthorized(tokenString string) bool {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Authentication error")
		}
		return secretKey, nil
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error at parsing token: %v", err.Error())
		return false
	}
	if token.Valid {
		return true
	}
	return false
}

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

func getCurrentUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if json.NewEncoder(w).Encode(currentUser) != nil {
		log.Fatal("Error at encoding current user")
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var body UserWrapper
	json.NewDecoder(r.Body).Decode(&body)
	e := json.NewEncoder(w)
	for _, user := range users {
		if body.User.Email == user.Email {
			if body.User.Password == user.Password {
				var tokenString, err = GenerateJWTToken(user.Username)
				if err != nil {
					e.Encode("Error at generating token")
					return
				}
				var responseUser UserWrapper
				responseUser.User.Email = user.Email
				responseUser.User.Token = tokenString
				responseUser.User.Username = user.Username
				responseUser.User.Bio = user.Bio
				responseUser.User.Image = user.Image
				if json.NewEncoder(w).Encode(responseUser) != nil {
					e.Encode("Error at encoding login response")
				}
				currentUser = responseUser
				return
			}
			e.Encode("Invalid password")
			return
		}
	}
	e.Encode("Invalid email")
}

func createUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var body UserWrapper
	json.NewDecoder(r.Body).Decode(&body)
	e := json.NewEncoder(w)
	for _, user := range users {
		if body.User.Username == user.Username {
			e.Encode("Username already exists")
			return
		}
		if body.User.Email == user.Email {
			e.Encode("Email already exists")
			return
		}
	}
	users = append(users, body.User)
	e.Encode(body)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var body UserWrapper
	json.NewDecoder(r.Body).Decode(&body)
	e := json.NewEncoder(w)
	fmt.Println(r.Header.Get("Authorization")[7:])
	if isAuthorized(r.Header.Get("Authorization")[7:]) {
		if body.User.Username != "" {
			currentUser.User.Username = body.User.Username
		}
		if body.User.Email != "" {
			currentUser.User.Email = body.User.Email
		}
		if body.User.Bio != "" {
			currentUser.User.Bio = body.User.Bio
		}
		if body.User.Image != "" {
			currentUser.User.Image = body.User.Image
		}
		if body.User.Password != "" {
			currentUser.User.Password = body.User.Password
		}
		return
	}
	e.Encode("Unauthorized")
}

func initDummyData() {
	profiles = append(profiles, Profile{Username: "berke", Bio: "I am a student.", Image: "", Following: false})
	profiles = append(profiles, Profile{Username: "keskul", Bio: "I am a cat.", Image: "", Following: false})

	articles = append(articles, Article{Slug: "how-to-train-your-dragon", Title: "How to train your dragon", Description: "Ever wonder how?", Body: "It takes a Jacobian", TagList: []string{"dragons", "training"}, CreatedAt: "2016-02-18T03:22:56.637Z", UpdatedAt: "2016-02-18T03:48:35.824Z", Favorited: false, FavoritesCount: 0, Author: &profiles[0]})

	comments = append(comments, Comment{Id: 1, CreatedAt: "2016-02-18T03:22:56.637Z", UpdatedAt: "2016-02-18T03:48:35.824Z", Body: "Nice post, thanks!", Author: &profiles[1], Slug: "how-to-train-your-dragon"})

	articles[0].Comments = append(articles[0].Comments, Comment{Id: 1, CreatedAt: "2016-02-18T03:22:56.637Z", UpdatedAt: "2016-02-18T03:48:35.824Z", Body: "Nice post, thanks!", Author: &profiles[1], Slug: "how-to-train-your-dragon"})

	tags = append(tags, "dragons")
	tags = append(tags, "training")

	users = append(users, User{
		Username: "keskul",
		Email:    "keskul@home.com",
		Password: "123456",
		Bio:      "I am a cat.",
		Image:    "",
	})
	users = append(users, User{
		Username: "berke",
		Email:    "berke.ahlatci@gmail.com",
		Password: "654321",
		Bio:      "I am a student.",
		Image:    "",
	})
}

func initRoutes(r *mux.Router) {
	// GET
	r.HandleFunc("/api/profiles/{username}", getProfile).Methods("GET")
	r.HandleFunc("/api/articles", getArticles).Methods("GET")
	r.HandleFunc("/api/articles/{slug}", getArticle).Methods("GET")
	r.HandleFunc("/api/articles/{slug}/comments", getComments).Methods("GET")
	r.HandleFunc("/api/tags", getTags).Methods("GET")
	r.HandleFunc("/api/user", getCurrentUser).Methods("GET")

	// POST
	r.HandleFunc("/api/users/login", login).Methods("POST")
	r.HandleFunc("/api/users", createUser).Methods("POST")

	//PUT
	r.HandleFunc("/api/user", updateUser).Methods("PUT")
}

func main() {
	r := mux.NewRouter()

	initDummyData()
	initRoutes(r)

	fmt.Println("Starting server at port 8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}
