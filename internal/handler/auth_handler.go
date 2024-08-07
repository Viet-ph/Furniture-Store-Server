package handler

import (
	"log"
	"net/http"

	"github.com/Viet-ph/Furniture-Store-Server/internal/dto"
	"github.com/Viet-ph/Furniture-Store-Server/internal/helper"
	"github.com/Viet-ph/Furniture-Store-Server/internal/service"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(a *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: a,
	}
}

func (a *AuthHandler) UserLogin() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type response struct {
		dto.User
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req, err := helper.Decode[request](r)
		if err != nil {
			log.Printf("Error decoding parameters: %s", err)
			w.WriteHeader(500)
			return
		}

		user, accessToken, refreshToken, err := a.authService.Login(r.Context(), req.Email, req.Password)
		if err != nil {
			helper.RespondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}

		helper.RespondWithJSON(w, http.StatusOK, response{
			User:         user,
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		})
	}
}

func (a *AuthHandler) RefreshAccessToken() http.HandlerFunc {
	type response struct {
		AccessToken string `json:"access_token"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		newAccessToken, err := a.authService.RefreshAccessToken(r.Context(), r)
		if err != nil {
			helper.RespondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}

		helper.RespondWithJSON(w, http.StatusAccepted, response{
			AccessToken: newAccessToken,
		})
	}
}

func (a *AuthHandler) RevokeRefreshToken() http.HandlerFunc {
	type request struct {
		RefreshToken string `json:"refresh_token"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req, err := helper.Decode[request](r)
		if err != nil {
			log.Printf("Error decoding parameters: %s", err)
			w.WriteHeader(500)
			return
		}

		err = a.authService.RevokeRefreshToken(r.Context(), req.RefreshToken)
		if err != nil {
			helper.RespondWithError(w, http.StatusUnauthorized, "Error revoking token: "+err.Error())
		}
	}
}
