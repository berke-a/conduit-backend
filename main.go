package main

import (
	"encoding/json"
	"errors"
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

type Error struct {
	Errors ErrorBody `json:"errors"`
}

type ErrorBody struct {
	Body []string `json:"body"`
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
			return nil, fmt.Errorf("authentication error")
		}
		return secretKey, nil
	})
	if err != nil {
		_, err := fmt.Fprintf(os.Stderr, "Error at parsing token: %v", err.Error())
		if err != nil {
			return false
		}
		return false
	}
	if token.Valid {
		return true
	}
	return false
}

func isContains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func getProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	var tmpProfile Profile

	if params["username"] == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Error{Errors: ErrorBody{Body: []string{"username is required"}}})
		return
	}

	for _, profile := range profiles {
		if profile.Username == params["username"] {
			tmpProfile = profile
			break
		}
	}

	json.NewEncoder(w).Encode(tmpProfile)

}

func getArticles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	limit := 20
	offset := 0
	var err error

	queries := r.URL.Query()
	tag := queries.Get("tag")
	author := queries.Get("author")
	favorited := queries.Get("favorited")
	queryLimit := queries.Get("limit")
	queryOffset := queries.Get("offset")

	encoder := json.NewEncoder(w)

	if queryLimit != "" {
		limit, err = strconv.Atoi(queryLimit)
		if err != nil {
			panic("Error at converting limit query")
			limit = 20
		}
	}
	if queryOffset != "" {
		offset, err = strconv.Atoi(queryOffset)
		if err != nil {
			panic("Error at converting offset query")
			offset = 0
		}
	}

	if offset+limit >= len(articles) {
		limit = len(articles) - offset
	}

	if tag != "" {
		var newArticles []Article

		for _, article := range articles {
			if isContains(article.TagList, tag) {
				newArticles = append(newArticles, article)
			}
		}

		if encoder.Encode(newArticles[offset:offset+limit]) != nil {
			panic("Error at encoding articles with tag query")
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
			panic("Error at encoding articles with author query")
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
			panic("Error at encoding articles with favorited query")
		}

		return
	}

	if encoder.Encode(articles[offset:offset+limit]) != nil {
		panic("Error at encoding articles without query")
	}

}

func getArticle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	if params["slug"] == "" {
		w.WriteHeader(http.StatusBadRequest)
		if json.NewEncoder(w).Encode(Error{Errors: ErrorBody{Body: []string{"article slug is required"}}}) != nil {
			panic("Error at encoding error at getArticle")
		}
		return
	}

	var tmpArticle Article

	for _, article := range articles {
		if article.Slug == params["slug"] {
			tmpArticle = article
			break
		}
	}

	if json.NewEncoder(w).Encode(tmpArticle) != nil {
		panic("Error at encoding article")
	}
	return

}

func getTags(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if json.NewEncoder(w).Encode(tags) != nil {
		panic("Error at encoding tags")
	}
}

func getComments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	if params["slug"] == "" {
		w.WriteHeader(http.StatusBadRequest)
		if json.NewEncoder(w).Encode(Error{Errors: ErrorBody{Body: []string{"article slug is required"}}}) != nil {
			panic("Error at encoding error at getComments")
		}
		return
	}

	var coms MultipleComments

	for _, article := range articles {
		if article.Slug == params["slug"] {

			coms.Comments = article.Comments

			break
		}
	}

	if json.NewEncoder(w).Encode(coms) != nil {
		panic("Error at encoding comments")
	}

}

func getCurrentUser(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if json.NewEncoder(w).Encode(currentUser) != nil {
		panic("Error at encoding current user")
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var body UserWrapper

	if json.NewDecoder(r.Body).Decode(&body) != nil {
		panic("Error at decoding login request")
	}
	encoder := json.NewEncoder(w)

	if body.User.Email == "" {
		w.WriteHeader(http.StatusBadRequest)
		if encoder.Encode(Error{Errors: ErrorBody{Body: []string{"email is required"}}}) != nil {
			panic("Error at encoding error at login")
		}
		return
	}

	if body.User.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		if encoder.Encode(Error{Errors: ErrorBody{Body: []string{"password is required"}}}) != nil {
			panic("Error at encoding error at login")
		}
		return
	}

	for _, user := range users {

		if body.User.Email == user.Email {

			if body.User.Password == user.Password {

				// Token creation
				var tokenString, err = GenerateJWTToken(user.Username)
				if err != nil {
					if encoder.Encode(errors.New("token generation failed")) != nil {
						panic("Error at encoding token error")
					}
					return
				}

				// Set response body
				var responseUser UserWrapper
				responseUser.User.Email = user.Email
				responseUser.User.Token = tokenString
				responseUser.User.Username = user.Username
				responseUser.User.Bio = user.Bio
				responseUser.User.Image = user.Image
				if json.NewEncoder(w).Encode(responseUser) != nil {
					panic("Error at encoding login response")
				}

				// Update current user
				currentUser = responseUser
				return
			}
			if encoder.Encode(errors.New("invalid password")) != nil {
				panic("Error at encoding invalid password")
			}
			return
		}
	}
	if encoder.Encode(errors.New("invalid email")) != nil {
		panic("Error at encoding invalid email")
	}
}

func createUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var body UserWrapper
	if json.NewDecoder(r.Body).Decode(&body) != nil {
		panic("Error at decoding user creation request")
	}
	encoder := json.NewEncoder(w)

	if body.User.Username == "" {
		w.WriteHeader(http.StatusBadRequest)
		if encoder.Encode(Error{Errors: ErrorBody{Body: []string{"username is required"}}}) != nil {
			panic("Error at encoding error at createUser")
		}
		return
	}

	if body.User.Email == "" {
		w.WriteHeader(http.StatusBadRequest)
		if encoder.Encode(Error{Errors: ErrorBody{Body: []string{"email is required"}}}) != nil {
			panic("Error at encoding error at createUser")
		}
		return
	}

	if body.User.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		if encoder.Encode(Error{Errors: ErrorBody{Body: []string{"password is required"}}}) != nil {
			panic("Error at encoding error at createUser")
		}
		return
	}

	for _, user := range users {

		if body.User.Username == user.Username {
			if encoder.Encode(errors.New("username already exists")) != nil {
				panic("Error at encoding username already exists error")
			}
			return
		}
		if body.User.Email == user.Email {
			if encoder.Encode(errors.New("email already exists")) != nil {
				panic("Error at encoding email already exists error")
			}
			return
		}

	}

	users = append(users, body.User)
	if encoder.Encode(body) != nil {
		panic("Error at encoding user creation response")
	}
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var body UserWrapper
	if json.NewDecoder(r.Body).Decode(&body) != nil {
		panic("Error at decoding user update request")
	}
	encoder := json.NewEncoder(w)

	if body.User.Username == "" && body.User.Email == "" && body.User.Bio == "" && body.User.Image == "" && body.User.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		if encoder.Encode(Error{Errors: ErrorBody{Body: []string{"at least one field is required"}}}) != nil {
			panic("Error at encoding error at updateUser")
		}
		return
	}

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

		if encoder.Encode(currentUser) != nil {
			panic("Error at encoding user update response")
		}

		return
	}

	if encoder.Encode(errors.New("authorization failed")) != nil {
		panic("Error at encoding authorization failed")
	}
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
