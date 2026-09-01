package main

import (
	"net/http"
	"testing"
)

func TestGetUser(t *testing.T) {
	app := newTestApplication(t)
	mux := app.mount()
	tokenTest, err := app.authenticator.GenerateToken(nil)
	if err != nil {
		t.Fatal("Fail to tokenize")
	}
	t.Run("should not allow unauthenticated requests", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/v1/users/1", nil)
		if err != nil {
			t.Fatal()
		}
		rr := executeRequest(req, mux)

		validateStatusCode(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("should allow authenticated requests", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/v1/users/1", nil)
		req.Header.Set("Authorization", "Bearer "+tokenTest)
		if err != nil {
			t.Fatal()
		}
		rr := executeRequest(req, mux)

		validateStatusCode(t, http.StatusOK, rr.Code)
	})
}
