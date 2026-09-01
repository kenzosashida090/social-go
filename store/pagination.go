package store

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type PaginationStructQuery struct {
	Limit  int      `json:"limit" validate:"gte=1 lte=100" `
	Offset int      `json:"offset" validate:"gte=0"`
	Sort   string   `json:"sort" validate:"oneof=asc desc"`
	Tags   []string `json:"tags" validate:"max=5"`
	Search string   `json:"search" validate:"max=100"`
	Since  string   `json:"since"`
	Until  string   `json:"until"`
}

func (fq PaginationStructQuery) Parse(r *http.Request) (PaginationStructQuery, error) {
	qs := r.URL.Query()

	limit := qs.Get("limit")
	if limit != "" {
		l, err := strconv.Atoi(limit)
		if err != nil {
			return fq, nil
		}
		fq.Limit = l
	}
	offset := qs.Get("offset")
	if offset != "" {
		o, err := strconv.Atoi(offset)
		if err != nil {
			return fq, nil
		}
		fq.Offset = o
	}
	sort := qs.Get("sort")
	if sort != "" {
		fq.Sort = sort
	}
	tags := qs.Get("tags")
	if tags != "" {

		fq.Tags = strings.Split(tags, ",")
	}

	search := qs.Get("search")
	if search != "" {
		fq.Search = search
	}

	since := qs.Get("since")
	if since != "" {
		date := parseDate(since)
		fq.Since = date
	}
	fmt.Println(fq.Search, "-----")
	return fq, nil
}

func parseDate(date string) string {
	t, err := time.Parse(time.DateTime, date)
	if err != nil {
		return ""
	}
	return t.Format(time.DateTime)
}
