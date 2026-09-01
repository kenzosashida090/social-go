package main

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/kenzosashida090/social/store"
)

type CreateUserPayload struct {
	Username string `json:"username" validate:"required,max=40"`
	Email    string `json:"email" validate:"required,max=255"`
	Password string `json:"password" validate:"required,min=8"`
}
type userKey string // for the middleware userId

const userCtx userKey = "user"

// GetUser godoc
//
//	@Summary		Fetches a user profile
//	@Description	Fetches a user profile by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id path int true "User ID"
//	@Success		200	{object}	store.User
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/{id} [get]
func (app *application) getUserById(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userId")
	id, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, err.Error())
	}
	user, err := app.getUser(r.Context(), id)
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, err.Error())
	}
	if err = app.responseJson(w, http.StatusOK, user); err != nil {
		WriteJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
}

func (app *application) userContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId := chi.URLParam(r, "userId")
		id, err := strconv.ParseInt(userId, 10, 64)
		if err != nil {
			app.badRequestError(w, r, err)
			return
		}
		ctx := r.Context()
		user, err := app.store.Users.GetUserById(ctx, id)
		if err != nil {
			WriteJSONError(w, http.StatusBadRequest, err.Error())
			return
		}
		ctx = context.WithValue(ctx, userCtx, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getUserFromContext(r *http.Request) *store.User {
	user, _ := r.Context().Value(userCtx).(*store.User)
	return user
}

// ActivateToken godoc
//
//	@Summary		Activate user by activating the token
//	@Description	The api will confirm the token send previouse on the email
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			token	path		string	true	"Token"
//	@Success		200		{object}	store.User
//	@Failure		400		{object}	error
//	@Failure		404		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/activate/{token} [put]
func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	err := app.store.Users.Activate(r.Context(), token)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}
}
