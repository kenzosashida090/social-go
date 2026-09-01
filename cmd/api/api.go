package main

import (
	"context"
	"expvar"

	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/kenzosashida090/social/docs"
	"github.com/kenzosashida090/social/internal/auth"
	"github.com/kenzosashida090/social/internal/ratelimiter"
	"github.com/kenzosashida090/social/mailer"
	"github.com/kenzosashida090/social/store"
	"github.com/kenzosashida090/social/store/cache"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"go.uber.org/zap"
)

type application struct {
	config                  config
	store                   store.Storage
	redisStorage            cache.UserStorage
	redisRateLimiterStorage cache.RateLimitStorage
	env                     string
	version                 string
	logger                  *zap.SugaredLogger
	mailer                  *mailer.MailtrapBody
	authenticator           auth.Authenticator
	rateLimiter             ratelimiter.Limiter
}

type config struct {
	addr             string
	frontEndUrl      string
	db               dbConfig
	mail             mailConfig
	auth             authConfig
	redis            redisConfig
	redisRateLimiter redisConfig
	rateLimitCfg     ratelimiter.Config
}
type authConfig struct {
	basic basicAuth
	token tokenConfig
}
type tokenConfig struct {
	secrete string
	exp     time.Duration
}
type basicAuth struct {
	username string
	password string
}
type mailConfig struct {
	exp time.Duration
}
type redisConfig struct {
	addrs    string
	pass     string
	db       int
	isActive bool
}
type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

func (app *application) mount() http.Handler {

	r := chi.NewRouter()
	r.Use(middleware.Recoverer) //Recovering from a panic
	r.Use(middleware.RequestID)
	r.Use(middleware.ClientIPFromRemoteAddr) // Retrieve ip addrrs client for rateLimit
	r.Use(middleware.Logger)
	r.Use(middleware.ClientIPFromRemoteAddr)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(app.ReateLimiterMiddleware)
	r.Route("/v1", func(r chi.Router) {
		r.Get("/swagger/*", httpSwagger.Handler(
			httpSwagger.URL(":3030/swagger/doc.json"),
		))
		r.With(app.BasicAuthMiddleware()).Get("/health", app.healthCheckHandler)
		r.With(app.BasicAuthMiddleware()).Get("/debug/vars", expvar.Handler().ServeHTTP)
		r.Route("/posts", func(r chi.Router) {
			r.Use(app.AuthUserMiddleware)
			r.Post("/", app.createPostHandler)
			r.Route("/{postId}", func(r chi.Router) {
				r.Use(app.GetPostByIdMiddleware)
				r.Get("/", app.getPostIdHandler) //change the getIdPost to a middleware
				r.Delete("/", app.validatePostOwnership("admin", app.deletePostHandler))
				r.Patch("/", app.validatePostOwnership("moderator", app.updatePostHandler))
			})
		})
		r.Route("/users", func(r chi.Router) {
			r.Route("/{userId}", func(r chi.Router) {
				r.Use(app.AuthUserMiddleware)
				r.Use(app.userContextMiddleware)
				r.Get("/", app.getUserById)
				r.Put("/follow", app.followHandler)
				r.Put("/unfollow", app.unfollowHandler)
			})
			r.Put("/activate/{token}", app.activateUserHandler)
			r.Group(func(r chi.Router) {
				r.Use(app.userContextMiddleware)
				r.Get("/feed", app.getUserFeedHandler)
			})
		})
		r.Route("/authentication", func(r chi.Router) {
			r.Post("/user", app.registerUserHandler)
			r.Post("/token", app.createTokenHandler)
		})
	})

	return r

}

func (app *application) run(mux http.Handler) error {
	docs.SwaggerInfo.Version = "0.0.1"
	//THis is for habndling the incoming url request, try to amtch to our urls
	fmt.Println(app.config.addr, "------------")
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}
	shutdown := make(chan error)
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGTERM)
		s := <-quit
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		app.logger.Infow("signal caught", "signal", s.String())
		shutdown <- srv.Shutdown(ctx)
	}()
	app.logger.Infow("Started server", "addr", srv.Addr, "env", app.env)
	err := srv.ListenAndServe()
	if err != nil {
		return err
	}
	app.logger.Infow("Server has stioed", "addr", app.config.addr, "env", app.env)
	return nil
}
