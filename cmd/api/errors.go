package main

import (
	"net/http"
)

func (app *application) intertalServerError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Error("Internal Server Error: %s path:%s error: %s", r.Method, r.URL.Path, err)
	WriteJSONError(w, http.StatusInternalServerError, "The server has a problem :c ")
}

func (app *application) badRequestError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Error("Bad Request Error %s path:%s error: %s", r.Method, r.URL.Path, err)
	WriteJSONError(w, http.StatusBadRequest, "Bad Request")
}

func (app *application) forbiddenResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Error("Forbidden Request Error  path: error:", r.URL.Path)
	WriteJSONError(w, http.StatusForbidden, "Forbidden access")

}

func (app *application) notFoundError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Error("Bad Request Error %s path:%s error: %s", r.Method, r.URL.Path, err)
	WriteJSONError(w, http.StatusNotFound, "Not found")
}

func (app *application) unAuthorizeRequest(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnf("UNATHOIRZE  Error %s path:%s error: %s", r.Method, r.URL.Path, err)

	WriteJSONError(w, http.StatusUnauthorized, "Unauthorized access")

}
func (app *application) rateLimitExceededResponse(w http.ResponseWriter, r *http.Request, retryAfter string) {
	app.logger.Warnf("Rate Limit Exceeded Errorr %s path:%s", r.Method, r.URL.Path)
	w.Header().Set("Retry-After", retryAfter)
	WriteJSONError(w, http.StatusTooManyRequests, "Rate limit Exceeded")
}
func (app *application) unAuthorizeResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnf("UNATHOIRZE  Error %s path:%s error: %s", r.Method, r.URL.Path, err)
	w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)

	WriteJSONError(w, http.StatusUnauthorized, err.Error())
}
