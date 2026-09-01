package main

import (
	"fmt"
	"net/http"
	"strconv"
)

type UserIdPayload struct {
	Id string `json:"user_id"`
}

// NewFollow	 godoc
//
//	@Summary		Follow a new user
//	@Description	YOu can follow a user and keep update with their posts
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	body		UserIdPayload	true	"FollowerId"
//	@Param			id	path		int				true	"UserID"
//	@Success		200	{object}	store.Post
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/follow/{id} [post]
func (app *application) followHandler(w http.ResponseWriter, r *http.Request) {
	var userId UserIdPayload
	err := ReadJSON(w, r, &userId)
	if err != nil {
		fmt.Println(err.Error())
		app.badRequestError(w, r, err)
		return
	}
	err = Validate.Struct(userId)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}
	uid, err := strconv.ParseInt(userId.Id, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	followerId := getUserFromContext(r)
	ctx := r.Context()
	err = app.store.Followers.Follow(ctx, followerId.ID, uid)
	if err != nil {
		app.intertalServerError(w, r, err)
	}

	app.responseJson(w, http.StatusOK, "ok")
}

// UnFollow	 godoc
//
//	@Summary		Follow a new user
//	@Description	YOu can follow a user and keep update with their posts
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	body		UserIdPayload	true	"FollowerId"
//	@Param			id	path		int				true	"UserID"
//	@Success		200	{object}	store.Post
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/unfollow/{id} [delete]
func (app *application) unfollowHandler(w http.ResponseWriter, r *http.Request) {
	var userId UserIdPayload
	err := ReadJSON(w, r, &userId)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}
	err = Validate.Struct(userId)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}
	uid, err := strconv.ParseInt(userId.Id, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	followerId := getUserFromContext(r)
	ctx := r.Context()
	err = app.store.Followers.Unfollow(ctx, uid, followerId.ID)
	if err != nil {
		app.intertalServerError(w, r, err)
		return
	}
	app.responseJson(w, http.StatusOK, "ok")

}
