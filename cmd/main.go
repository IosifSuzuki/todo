package main

import (
	db "github.com/IosifSuzuki/todo/internall/db"
	"github.com/IosifSuzuki/todo/internall/handler"
	"github.com/IosifSuzuki/todo/internall/logger"
	"github.com/IosifSuzuki/todo/internall/middleware"
	"github.com/IosifSuzuki/todo/internall/utility"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func main() {
	rootRouter := configureRouter()
	utility.Setup()
	db.ConnectToDB()
	defer func() {
		_ = db.CloseConnectionToDB()
	}()

	server := http.Server{
		Addr:         ":8080",
		Handler:      rootRouter,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	logger.Info("Server is listening...")
	if err := server.ListenAndServe(); err != nil {
		logger.Fatal("Server failed started with error", zap.Error(err))
	}
}

func configureRouter() http.Handler {
	var amw = middleware.AuthenticationMiddleware{}
	var lms = middleware.LoggerMiddleware{}
	var rootRouter = mux.NewRouter()
	rootRouter.StrictSlash(true)

	var authenticationRouter = rootRouter.PathPrefix("/authentication").Subrouter()
	authenticationRouter.Use(lms.Middleware)
	authenticationRouter.HandleFunc("/sign-in", handler.SignInHandler).Methods(http.MethodPost)
	authenticationRouter.HandleFunc("/sign-up", handler.SignUpHandler).Methods(http.MethodPost)
	authenticationRouter.HandleFunc("/refresh-token", handler.RefreshTokenHandler).Methods(http.MethodPost)

	var todoRouter = rootRouter.PathPrefix("/todo").Subrouter()
	todoRouter.Use(lms.Middleware)
	todoRouter.Use(amw.Middleware)
	todoRouter.HandleFunc("/ping", handler.HomeHandler).Methods(http.MethodGet)

	var accountRouter = rootRouter.PathPrefix("/account").Subrouter()
	accountRouter.Use(lms.Middleware)
	accountRouter.Use(amw.Middleware)
	accountRouter.HandleFunc("/user/{id}}", handler.UserInfoHanlder).Methods(http.MethodGet)

	return rootRouter
}
