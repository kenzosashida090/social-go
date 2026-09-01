package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"

	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/kenzosashida090/social/store"
)

type RegisterUserPayload struct {
	Username string `json:"username" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,max=255"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

type CreateUserTokenPayload struct {
	Username string `json:"username" `
	Password string `json:"password"`
}
type UserTokenPayload struct {
	*store.User
	Token string `json:"token"`
}

// registerUserHandler godoc
//
//	@Summary		Create a new account
//	@Description	YOu can create an account
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			data	body		RegisterUserPayload	true	"FollowerId"
//	@Success		200		{object}	UserTokenPayload	"User registered"
//	@Failure		400		{object}	error
//	@Failure		404		{object}	error
//	@Router			/authentication/user [post]
func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var userPayload RegisterUserPayload

	err := ReadJSON(w, r, &userPayload)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	err = Validate.Struct(userPayload)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}
	user := &store.User{
		Username: userPayload.Username,
		Email:    userPayload.Email,
		Role: store.Role{
			Name: "moderator",
		},
	}
	err = user.Password.Set(userPayload.Password)
	if err != nil {
		app.intertalServerError(w, r, err)
	}
	ctx := r.Context()
	plainToken := uuid.New().String()
	hash := sha256.Sum256([]byte(plainToken))
	hashToken := hex.EncodeToString(hash[:])

	registerResponse := &UserTokenPayload{
		User:  user,
		Token: plainToken,
	}

	err = app.store.Users.CreateAndInvite(ctx, user, hashToken, app.config.mail.exp)

	if err != nil {
		app.intertalServerError(w, r, err)
		return
	}
	var maxTries = 3
	for i := range maxTries {
		err = app.mailer.SendActivationMail(hashToken, user.Email, user.Username, fmt.Sprintf("%s/activate/%s", app.config.frontEndUrl, plainToken))
		if err == nil {
			fmt.Println("salamaleco", i)
			app.responseJson(w, http.StatusCreated, registerResponse)
			return
		}
		if i < maxTries-1 {
			time.Sleep(time.Second * time.Duration(i+1))
			continue
		}
	}
	err = app.store.Users.Delete(ctx, user.ID)
	if err != nil {
		app.intertalServerError(w, r, err)

	}
	err = app.store.Users.Deactivate(ctx, user.ID)
	if err != nil {
		app.intertalServerError(w, r, err)
	}
}

// registerUserHandler godoc
//
//	@Summary		Create a token
//	@Description	Create a token for a user
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			data	body		CreateUserTokenPayload	true	"User credentials"
//	@Success		200		{string}	string					"Token"
//	@Failure		400		{object}	error
//	@Failure		404		{object}	error
//	@Router			/authentication/token [post]
func (app *application) createTokenHandler(w http.ResponseWriter, r *http.Request) {
	var bodyPayload CreateUserTokenPayload
	err := ReadJSON(w, r, &bodyPayload)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("IS IT FROM HERE")
		app.badRequestError(w, r, err)
		return
	}
	err = Validate.Struct(bodyPayload)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}
	ctx := r.Context()
	queryUser, err := app.store.Users.GetUserByUsername(ctx, bodyPayload.Username)
	user := &store.User{}
	isCorrect := user.Password.Compare(bodyPayload.Password, queryUser.Password)

	if err != nil {

		app.intertalServerError(w, r, err)
		return
	}
	claims := jwt.MapClaims{
		"sub": queryUser.ID,
		"exp": time.Now().Add(app.config.auth.token.exp).Unix(),
		"iat": time.Now().Unix(),
		"aud": "socialgo",
		"nbf": time.Now().Unix(),
		"iss": "socialgo",
	}
	token, err := app.authenticator.GenerateToken(claims)
	if err != nil {
		app.intertalServerError(w, r, err)
	}
	fmt.Println(token, "=token")
	app.responseJson(w, http.StatusAccepted, isCorrect)
}
