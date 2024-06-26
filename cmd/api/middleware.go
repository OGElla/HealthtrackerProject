package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/OGElla/Project-API/internal/data"
	"github.com/OGElla/Project-API/internal/validator"
)

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Add("Vary", "Authorization")

		authorizationHeader := r.Header.Get("Authorization")

		if authorizationHeader == "" {
			r = app.contextSetUser(r, data.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}

		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}


		token := headerParts[1]

		v := validator.New()

		if data.ValidateTokenPlaintext(v, token); !v.Valid() {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		user, err := app.models.Users.GetForToken(data.ScopeAuthentication, token)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
				app.invalidAuthenticationTokenResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
			}
			return
		}

		r = app.contextSetUser(r, user)
		data.CurrentToken = token
		data.CurrentUserID, err = app.models.Tokens.CheckForUser(token)
		if err != nil {
			fmt.Println(err)
		}
		next.ServeHTTP(w, r)
	})
}


func (app *application) requireActivatedUser(next http.HandlerFunc) http.HandlerFunc { 
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	user := app.contextGetUser(r)

	if user.IsAnonymous() {
			app.authenticationRequiredResponse(w, r)
	return
	}

	if !user.Activated {
		app.inactiveAccountResponse(w, r)
		return
	}

	next.ServeHTTP(w, r) })
	}

func (app *application) requirePermission(code string, next http.HandlerFunc) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)
	permissions, err := app.models.Permissions.GetAllForUser(user.ID) 
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if !permissions.Include(code) {
		app.notPermittedResponse(w, r)
		return
	}

	next.ServeHTTP(w, r)
}

return app.requireActivatedUser(fn) 
}