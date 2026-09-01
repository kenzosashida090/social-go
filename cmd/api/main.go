package main

import (
	"expvar"
	"fmt"

	"log"
	"runtime"
	"time"

	"github.com/kenzosashida090/social/db"
	"github.com/kenzosashida090/social/internal/auth"
	"github.com/kenzosashida090/social/internal/env"
	"github.com/kenzosashida090/social/internal/ratelimiter"
	"github.com/kenzosashida090/social/mailer"
	"github.com/kenzosashida090/social/store/cache"
	"github.com/redis/go-redis/v9"

	"github.com/kenzosashida090/social/store"
	"go.uber.org/zap"
)

//	@title			Social Kib
//	@description	This is an API for a social network.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@BasePath	/v1

// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description
func main() {
	var rediClient *redis.Client
	var redisRateLimiterClient *redis.Client
	cfg := config{
		addr:        env.GetString("ADDR", ":3000"),
		frontEndUrl: env.GetString("FRONTEND_URL_DEV", "http://localhost:4000"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", ""),
			maxOpenConns: env.GetNumber("MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetNumber("MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("MAX_IDLE_TIME", "15m"),
		},
		redis: redisConfig{
			addrs:    env.GetString("REDIS_ADDR", ""),
			pass:     env.GetString("REDIS_PASSWORD", ""),
			db:       0,
			isActive: env.GetBool("REDIS_ACTIVE", false),
		},
		redisRateLimiter: redisConfig{
			addrs:    env.GetString("REDIS_RATE_ADDRS", ""),
			pass:     env.GetString("REDIS_RATE_PASS", ""),
			db:       0,
			isActive: env.GetBool("REDIS_ACTIVE", false),
		},
		mail: mailConfig{
			exp: time.Hour * 24 * 3,
		},
		rateLimitCfg: ratelimiter.Config{
			RequestPerTimeFrame: 20,
			TimeFrame:           time.Second * 4,
			Enabled:             true,
		},
		auth: authConfig{
			basic: basicAuth{
				username: env.GetString("USERNAME_AUTH", "admin"),
				password: env.GetString("PASSWORD_AUTH", "1234"),
			},
			token: tokenConfig{
				secrete: env.GetString("SECRET_JWT", "kibdev"),
				exp:     time.Hour * 24 * 3,
			},
		},
	}
	zap, _ := zap.NewProduction()
	defer zap.Sync()
	logger := zap.Sugar()
	db, err := db.New(
		cfg.db.addr,
		cfg.db.maxIdleConns,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleTime,
	)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()
	/// REDIS CLIENT
	if cfg.redis.isActive {
		logger.Info("Redis Active")
		rediClient = cache.NewConnectionClient(cfg.redis.addrs, cfg.redis.pass, cfg.redis.db)
	}
	fmt.Println(cfg.redis.isActive)
	defer rediClient.Close()
	// REDIS ratelimiter
	if cfg.redisRateLimiter.isActive {
		logger.Info("REDIS RATE LIMITER ACTIVE")
		redisRateLimiterClient = cache.NewConnectionClient(cfg.redisRateLimiter.addrs, cfg.redisRateLimiter.pass, cfg.redisRateLimiter.db)
	}
	defer redisRateLimiterClient.Close()
	logger.Info("DB CONECTION ESTABLISHED")
	store := store.NewStorage(db)
	mailtrap := mailer.ConnectMailTrap()

	//token JWT
	tokenHost := "socialgo"
	jwtAuthenticator := auth.NewJWTAuthenticator(cfg.auth.token.secrete, tokenHost, tokenHost)
	rateLimiter := ratelimiter.NewFixedWindowLimiter(
		cfg.rateLimitCfg.RequestPerTimeFrame,
		cfg.rateLimitCfg.TimeFrame,
	)
	app := &application{
		config:                  cfg,
		store:                   store,
		redisStorage:            cache.NewRedisStorage(rediClient),
		redisRateLimiterStorage: cache.NewRedisLimitRateStorage(redisRateLimiterClient),
		env:                     env.GetString("ENV", "DEV"),
		version:                 env.GetString("VERSION", "1"),
		logger:                  logger,
		mailer:                  mailtrap,
		authenticator:           jwtAuthenticator,
		rateLimiter:             rateLimiter,
	}
	//Metrics
	expvar.NewString("version").Set(app.version)
	expvar.Publish("database", expvar.Func(func() any {
		return db.Stats()
	}))
	expvar.Publish("goroutines", expvar.Func(func() any {
		return runtime.NumGoroutine()
	}))
	//app.mailer.SendActivationMail("panconjamon", "lakakarodriguez@gmail.com")
	mux := app.mount()
	logger.Fatal(app.run(mux))
}
