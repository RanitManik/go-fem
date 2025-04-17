package api

import (
	"encoding/json"
	"github.com/RanitManik/go-fem/internal/store"
	"github.com/RanitManik/go-fem/internal/tokens"
	"github.com/RanitManik/go-fem/internal/utils"
	"log"
	"net/http"
	"time"
)

type TokenHandler struct {
	tokenStore store.TokenStore
	userStore  store.UserStore
	logger     *log.Logger
}

type createTokenRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewTokenHandler(tokenStore store.TokenStore, userStore store.UserStore, logger *log.Logger) *TokenHandler {
	return &TokenHandler{
		tokenStore: tokenStore,
		userStore:  userStore,
		logger:     logger,
	}
}

func (h *TokenHandler) HandleCreateToken(w http.ResponseWriter, r *http.Request) {
	var req createTokenRequest
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		h.logger.Printf("ERROR: create token request: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}

	user, err := h.userStore.GetUserByUsername(req.Username)
	if err != nil || user == nil {
		h.logger.Printf("ERROR: get user by username: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	passwordsDoMatch, err := user.PasswordHash.Matches(req.Password)

	if err != nil {
		h.logger.Printf("ERROR: PasswordHash.Matches %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
	}

	if !passwordsDoMatch {
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "invalid credentials"})
		return
	}

	token, err := h.tokenStore.CreateNewToken(user.ID, 24*time.Hour, tokens.ScopeAuth)
	if err != nil {
		h.logger.Printf("ERROR: Creating Token: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})

	}

	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"auth_token": token})

}
