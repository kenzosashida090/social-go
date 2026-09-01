package main

import (
	"context"
	"errors"
	"fmt"

	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/kenzosashida090/social/store"
)

type CreatePostPayload struct {
	Title   string   `json:"title" validate:"required,min=8,max=100"`
	Message string   `json:"message" validate:"required,min=20,max=1000"`
	Tags    []string `json:"tags"`
}
type UpdatePostPayload struct {
	Title   string   `json:"title" validate:"max=100"`
	Message string   `json:"message" validate:"max=100"`
	Tags    []string `json:"tags"`
}
type postKey string

const postCtx postKey = "post"

// CreatePost godoc
//
//	@Summary		Create post
//	@Description	Create post
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			body	body		CreatePostPayload	true	"Create post payload"
//	@Success		200		{object}	store.Post
//	@Failure		400		{object}	error
//	@Failure		404		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/posts [post]
func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreatePostPayload
	err := ReadJSON(w, r, &payload)
	user := getUserFromContext(r)
	if err != nil {

		app.badRequestError(w, r, err)
		return
	}
	err = Validate.Struct(payload)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}
	fmt.Println(user, "0--")
	post := &store.Post{
		Title:   payload.Title,
		Content: payload.Message,
		Tags:    payload.Tags,
		UserID:  user.ID,
	}
	ctx := r.Context()

	err = app.store.Posts.Create(ctx, post)
	if err != nil {
		fmt.Println(err)
		app.intertalServerError(w, r, err)
	}
}

// GetPostById	 godoc
//
//	@Summary		Get a post by the id
//	@Description	GEt post by id
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Post ID"
//	@Success		200	{object}	store.Post
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/posts/{id} [get]
func (app *application) getPostIdHandler(w http.ResponseWriter, r *http.Request) {
	postId := chi.URLParam(r, "postId")
	id, err := strconv.ParseInt(postId, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}
	ctx := r.Context()
	post, err := app.store.Posts.GetById(ctx, id)
	fmt.Println("---", post)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundError(w, r, err)
			return
		default:
			app.badRequestError(w, r, err)
			return

		}
	}
	comment, err := app.store.Comments.GetCommentsByUserId(ctx, id)
	if err != nil {
		app.intertalServerError(w, r, err)

		return
	}
	post.Comments = comment
	writeJSON(w, 200, post)

}
func (app *application) GetPostByIdMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		postId := chi.URLParam(r, "postId")
		id, err := strconv.ParseInt(postId, 10, 64)
		if err != nil {
			app.badRequestError(w, r, err)
			return
		}
		ctx := r.Context()
		fmt.Println("00d0d00d", id)
		post, err := app.store.Posts.GetById(ctx, id)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				app.notFoundError(w, r, err)
				return
			default:
				app.badRequestError(w, r, err)
				return

			}
		}

		ctx = context.WithValue(ctx, postCtx, post)

		next.ServeHTTP(w, r.WithContext(ctx))

	})
}

// DeletePostById	 godoc
//
//	@Summary		Delete a post by the id
//	@Description	Delete post by id
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			postId	path		int	true	"Post ID"
//	@Success		200	{object}	store.Post
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/posts/{postId} [delete]
func (app *application) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	postId := chi.URLParam(r, "postId")
	fmt.Println(postId, "from delete")
	id, err := strconv.ParseInt(postId, 10, 64)
	fmt.Println(postId, "-----")
	ctx := r.Context()
	err = app.store.Posts.DeleteById(ctx, id)
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, "Error deleting data")
	}
	writeJSON(w, 200, "ok")
}

// PatchPostById	 godoc
//
//	@Summary		Update a post by the id
//	@Description	Update post by id
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int					true	"Post ID"
//	@Param			body	body		UpdatePostPayload	true	"Update post Payload"
//	@Success		200		{object}	store.Post
//	@Failure		400		{object}	error
//	@Failure		404		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/posts/{id} [patch]
func (app *application) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	var updatePayload UpdatePostPayload
	fmt.Println("METHOD:", r.Method)
	fmt.Println("CONTENT-TYPE:", r.Header.Get("Content-Type"))
	fmt.Println("CONTENT-LENGTH:", r.ContentLength)
	postId := chi.URLParam(r, "postId")
	ctx := r.Context()
	id, err := strconv.ParseInt(postId, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}
	post, err := app.store.Posts.GetById(ctx, id)
	if err != nil {
		fmt.Println(err.Error(), "-------")
	}
	err = ReadJSON(w, r, &updatePayload)
	if err != nil {
		app.intertalServerError(w, r, err)
		return
	}
	if updatePayload.Title != "" {
		post.Title = updatePayload.Title
	} else if updatePayload.Message != "" {
		post.Content = updatePayload.Message
	} else if updatePayload.Tags != nil {
		post.Tags = updatePayload.Tags
	}

	err = Validate.Struct(post)
	if err != nil {
		app.intertalServerError(w, r, err)
		return
	}
	fmt.Println("goasofoas", id)
	updatedPost, err := app.store.Posts.UpdatePost(ctx, id, post)
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, err.Error())
	}
	fmt.Println("updatedpost")
	app.responseJson(w, http.StatusOK, updatedPost)
}

func (app *application) validatePostOwnership(requiredRole string, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := getUserFromContext(r)
		post := getPostFromContext(r)

		if user.ID != post.UserID {
			next.ServeHTTP(w, r)
			return
		}
		allowed, err := app.checkProcedence(r.Context(), user, requiredRole)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		if !allowed {
			app.forbiddenResponse(w, r, err)
			return
		}
		next.ServeHTTP(w, r)
	})
}
func (app *application) checkProcedence(ctx context.Context, user *store.User, roleName string) (bool, error) {
	role, err := app.store.Roles.GetRoleByName(ctx, roleName)
	if err != nil {
		return false, err
	}

	return user.Role.Level >= role.Level, nil

}
func getPostFromContext(r *http.Request) *store.Post {
	post, _ := r.Context().Value(postCtx).(*store.Post)
	return post
}
