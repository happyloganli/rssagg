package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/happyloganli/rssagg/internal/database"
	"log"
	"net/http"
	"time"
)

func (apiCfg *apiConfig) CreateUserHandler(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Name string `json:"name"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request payload: %s", err))
		return
	}

	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      params.Name,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Could not create user: %s", err))
	}

	respondWithJSON(w, http.StatusCreated, databaseUserToUser(user))
}

func (apiCfg *apiConfig) GetUserHandler(w http.ResponseWriter, r *http.Request, user database.User) {

	respondWithJSON(w, http.StatusOK, databaseUserToUser(user))
}

func (apiCfg *apiConfig) GetUserPostsHandler(w http.ResponseWriter, r *http.Request, user database.User) {

	log.Printf("user %s get posts", user.Name)
	log.Printf("user id %s", user.ID)

	posts, err := apiCfg.DB.GetPostsForUser(r.Context(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(10),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Could not get posts: %s", err))
		return
	}

	respondWithJSON(w, http.StatusOK, posts)
}
