package main

import (
	"fmt"
	"github.com/happyloganli/rssagg/internal/auth"
	"github.com/happyloganli/rssagg/internal/database"
	"net/http"
)

type authedHandler func(w http.ResponseWriter, r *http.Request, user database.User)

func (apiCfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetAPIKey(r.Header)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Auth error: %s", err))
			return
		}

		user, err := apiCfg.DB.GerUserByApiKey(r.Context(), apiKey)

		if err != nil {
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Could not get user: %s. Error: %s", apiKey, err))
			return
		}

		handler(w, r, user)
	}
}
