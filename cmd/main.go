package main

import (
	_ "github.com/IosifSuzuki/todo/docs"
	db "github.com/IosifSuzuki/todo/internall/db"
	"github.com/IosifSuzuki/todo/internall/handler"
	"github.com/IosifSuzuki/todo/internall/logger"
	"github.com/IosifSuzuki/todo/internall/middleware"
	"github.com/IosifSuzuki/todo/internall/utility"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
	"net/http"
	"time"
)

// @title Todo API
// @version 1.0
// @description This is api documentation for todo
// @termsOfService http://swagger.io/terms/
// @schemes http https

// @contact.name API Documentation Support
// @contact.url http://www.swagger.io/support
// @contact.email iosifsuzuki@gmail.com

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host todo-app
// @BasePath /api/v1
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

	var apiRouter = rootRouter.PathPrefix("/api/v1").Subrouter()
	apiRouter.Use(lms.Middleware)
	//apiRouter.PathPrefix("/").Handler(handler.CommonHanlder{})
	apiRouter.StrictSlash(true)

	var accountRouter = apiRouter.PathPrefix("/account").Subrouter()
	accountRouter.Use(amw.Middleware)
	accountRouter.HandleFunc("/user/{id:[0-9]+}", handler.UserInfoHandler).Methods(http.MethodGet)
	accountRouter.HandleFunc("/users", handler.UsersInfoHanlder).Methods(http.MethodGet)

	var authenticationRouter = apiRouter.PathPrefix("/authentication").Subrouter()
	authenticationRouter.Use(lms.Middleware)
	authenticationRouter.HandleFunc("/sign-in", handler.SignInHandler).Methods(http.MethodPost)
	authenticationRouter.HandleFunc("/sign-up", handler.SignUpHandler).Methods(http.MethodPost)
	authenticationRouter.HandleFunc("/refresh-token", handler.RefreshTokenHandler).Methods(http.MethodPost)

	var todoRouter = apiRouter.PathPrefix("/todo").Subrouter()
	todoRouter.Use(amw.Middleware)
	todoRouter.HandleFunc("/ping", handler.HomeHandler).Methods(http.MethodGet)
	todoRouter.HandleFunc("/my/todos", handler.MyTodosHandler).Methods(http.MethodGet)
	todoRouter.HandleFunc("/add", handler.AddTodoHandler).Methods(http.MethodPost)
	todoRouter.HandleFunc("/remove/{id}", handler.RemoveTodoHandler).Methods(http.MethodDelete)
	todoRouter.HandleFunc("/{id}", handler.GetTodoHandler).Methods(http.MethodGet)
	todoRouter.HandleFunc("/toggle/{id}", handler.ToggleTodoHandler).Methods(http.MethodPut)

	rootRouter.PathPrefix("/doc").Handler(httpSwagger.WrapHandler)

	return rootRouter
}
