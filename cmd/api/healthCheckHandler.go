package main

import "net/http"

// add this handler to the application pointer
func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":  "ok",
		"version": app.version,
		"env":     app.env,
	}
	err := writeJSON(w, http.StatusOK, data)
	if err != nil {
		WriteJSONError(w, http.StatusInternalServerError, err.Error())
	}

}
