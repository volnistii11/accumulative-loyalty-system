package auth

import "net/http"

type Auth struct {
}

func NewAuth() *Auth {
	return &Auth{}
}

func (a *Auth) RegisterUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO:
	}
}

func (a *Auth) AuthenticateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO:
	}
}
