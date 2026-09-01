package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/golang-jwt/jwt/v5"
	"github.com/kenzosashida090/social/store"
)

func (app *application) BasicAuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				fmt.Println("BAD CHOICES")
				app.unAuthorizeResponse(w, r, fmt.Errorf("Unavailable to authorize user"))
				return
			}
			parts := strings.Split(authHeader, " ")
			fmt.Println(parts)
			if len(parts) != 2 || parts[0] != "Basic" {
				app.unAuthorizeResponse(w, r, fmt.Errorf("Unavailable to authorize the header"))
				return
			}
			decode, err := base64.StdEncoding.DecodeString(parts[1])
			if err != nil {
				app.unAuthorizeResponse(w, r, fmt.Errorf("Unavailable to authorize"))
				return
			}
			username := app.config.auth.basic.username
			password := app.config.auth.basic.password
			fmt.Println(username, password, "PAAAAAAAAAAAAAAAAAAAAAA")
			creds := strings.SplitN(string(decode), ":", 2)
			fmt.Println(creds, "CREDS")
			if len(creds) != 2 || creds[0] != username || creds[1] != password {
				app.unAuthorizeResponse(w, r, fmt.Errorf("Fail to auth"))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func (app *application) AuthUserMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			app.unAuthorizeRequest(w, r, fmt.Errorf("Not available for this moment."))
			return
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			app.unAuthorizeRequest(w, r, fmt.Errorf("Not available for this moment."))
			return
		}
		jwtToken, err := app.authenticator.ValidateToken(parts[1])
		if err != nil {
			fmt.Println(err.Error())
			app.unAuthorizeRequest(w, r, fmt.Errorf("Unavailable"))
			return
		}
		claims, _ := jwtToken.Claims.(jwt.MapClaims)
		userId, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["sub"]), 10, 64)
		fmt.Println(userId, "---------")
		if err != nil {
			app.unAuthorizeRequest(w, r, fmt.Errorf("Unavailable"))
			return
		}
		ctx := r.Context()
		user, err := app.getUser(ctx, userId)
		ctx = context.WithValue(ctx, userCtx, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) getUser(ctx context.Context, userId int64) (*store.User, error) {
	userRedis, err := app.redisStorage.Users.Get(ctx, userId)
	fmt.Println(userRedis, "user from redi")
	if err != nil {
		fmt.Println("jesus es rey de reyes")
		fmt.Println(err.Error())
		return nil, err
	}

	if userRedis == nil {
		fmt.Println("user from, db", userId)
		userDB, err := app.store.Users.GetUserById(ctx, userId)
		if err != nil {
			return nil, err
		}
		err = app.redisStorage.Users.Set(ctx, userDB)
		if err != nil {
			fmt.Println("is there anb error", err)
			return nil, err
		}
		return userDB, nil
	}
	return userRedis, nil

}

type RateLimiterType struct {
	Addrs string `json:"addrs"`
	Count int64  `json:"count"`
}

func (app *application) ReateLimiterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const limit = 20
		if app.config.rateLimitCfg.Enabled {
			ip := middleware.GetClientIP(r.Context())
			rateType, err := app.redisRateLimiterStorage.RateLimiter.Get(r.Context(), ip)

			if err != nil {
				fmt.Println("jesus es rey de reyes", "redis from rate")
				app.notFoundError(w, r, err)
				return
			}

			fmt.Println("cace ratellimit empty")
			rateType++
			fmt.Printf("rateLimiter: %#v\n", app.rateLimiter)
			if allow, retryAfter := app.rateLimiter.Allow(rateType, 40); !allow {
				app.rateLimitExceededResponse(w, r, retryAfter.String())
				return
			}
			err = app.redisRateLimiterStorage.RateLimiter.Set(r.Context(), ip, rateType)
			if err != nil {
				fmt.Println("is there anb error", err)
				app.intertalServerError(w, r, err)
			}
			next.ServeHTTP(w, r)
			return

		}
		next.ServeHTTP(w, r)
	})
}
