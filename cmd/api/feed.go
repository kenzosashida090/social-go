package main

import (
	"fmt"
	"net/http"

	"github.com/kenzosashida090/social/store"
)

//  GetUserFeed godoc

// @Summary		Get the user feed
// @Description	Get a user
// @Tags			users
// @Accept			json
// @Produce		json
// @Param			limit	query		int			false	"Limit search of query"
// @Param			offset	query		int			false	"Offset for each search"
// @Param			sort	query		string		false	"Sort asc or desc"
// @Param			search	query		string		false	"search title or content"
// @Param			since	query		string		false	"search since the post"
// @Param			tags	query		[]string	false	"search post from tags"
// @Success		200		{object}	store.User
// @Failure		400		{object}	error
// @Failure		404		{object}	error
// @Security		ApiKeyAuth
// @Router			/feed [get]
func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	fq := store.PaginationStructQuery{
		Limit:  20,
		Offset: 0,
		Sort:   "desc",
	}
	fq, err := fq.Parse(r)
	fmt.Println("----FQ FROM FEED---- ", fq)
	feed, err := app.store.Posts.GetUserFeed(ctx, int64(320), fq)

	fmt.Println(feed, "-----FEED-----")
	if err != nil {
		fmt.Println("----------", err.Error(), "----------------")
		app.intertalServerError(w, r, err)
		return
	}

	if err := app.responseJson(w, http.StatusOK, feed); err != nil {
		app.intertalServerError(w, r, err)
	}
}
