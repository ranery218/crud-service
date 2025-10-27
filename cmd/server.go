package cmd

import (
	"context"
	id_gen "crud/internal/adapters/id_generator"
	"crud/internal/adapters/password"
	"crud/internal/adapters/repository/postgres"
	redisStore "crud/internal/adapters/session/redis"
	"crud/internal/config"
	"crud/internal/services/user"
	httpapi "crud/internal/transport/http"
	"crud/internal/transport/http/middleware"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func RunServer() error {
	_ = godotenv.Load()
	config, err := config.Load("config.yaml")
	if err != nil {
		return err
	}

	dsnURL := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(config.Postgres.User, config.Postgres.Password),
		Host:   fmt.Sprintf("%s:%d", config.Postgres.Host, config.Postgres.Port),
		Path:   config.Postgres.DB,
	}

	q := dsnURL.Query()
	q.Set("sslmode", config.Postgres.SSLMode)
	dsnURL.RawQuery = q.Encode()

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsnURL.String())
	if err != nil {
		return fmt.Errorf("unable to connect to database: %w", err)
	}
	defer pool.Close()

	repo := postgres.NewUserRepository(pool)
	idGen := id_gen.NewDefaultIDGen()
	hasher := password.NewBcryptHasher(0)

	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d",config.Redis.Host,config.Redis.Port),
		Password: "",
		DB: 0,
	})
	sessionStore := redisStore.NewRedisStore(rdb, config.Session.TTL, idGen.NewID)

	registerService := user.NewRegisterService(repo, hasher, idGen)
	loginService := user.NewLoginService(repo, hasher, sessionStore)
	updateService := user.NewUpdateService(repo, hasher)
	deleteService := user.NewDeleteService(repo)
	logger := log.New(os.Stdout, "[http] ", log.LstdFlags|log.Lshortfile)
	userHandler := httpapi.NewUserHandler(registerService, loginService, updateService, deleteService, logger)
	authHandler := middleware.NewAuthMiddleware(sessionStore)

	router := httpapi.NewRouter(userHandler, authHandler)

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port),
		Handler: router,
	}
	logger.Printf("Starting server on %s:%d", config.Server.Host, config.Server.Port)
	return server.ListenAndServe()
}
