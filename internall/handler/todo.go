package handler

import (
	"encoding/json"
	"fmt"
	"github.com/IosifSuzuki/todo/internall/db"
	"github.com/IosifSuzuki/todo/internall/logger"
	"github.com/IosifSuzuki/todo/internall/model"
	"github.com/IosifSuzuki/todo/internall/model/request"
	"github.com/IosifSuzuki/todo/internall/utility"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

// HomeHandler docs
// @Summary Check connection to server through bearer token
// @Tags todo
// @ID home-handler
// @Accept   json
// @Produce  json
// @Security ApiKeyAuth
// @param    Authorization header string true "Authorization"
// @Success  200 {object} model.Ping
// @Failure  401 {object} model.ResponseError
// @Failure  500 {object} model.ResponseError
// @Router   /todo/ping [get]
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	userId, ok := r.Context().Value(utility.UserIdKey).(int)
	if !ok {
		err := model.ResponseError{
			Code:    http.StatusUnauthorized,
			Message: "Cannot retrieve user id",
		}
		logger.Error(err.Message)
		w.WriteHeader(err.Code)
		return
	}
	var pingMessage = model.Ping{
		UserId:  userId,
		Message: "Success made request",
		Code:    http.StatusOK,
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(pingMessage); err != nil {
		logger.Error("Error occurred during encoding", zap.Error(err))
	}
}

// MyTodosHandler docs
// @Summary Get my todos list server
// @Tags todo
// @ID my-todos-handler
// @Accept   json
// @Produce  json
// @Security ApiKeyAuth
// @param    Authorization header string true "Authorization"
// @Success  200 {array} model.Todo
// @Failure  401 {object} model.ResponseError
// @Failure  500 {object} model.ResponseError
// @Router   /todo/my/todos [get]
func MyTodosHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	userId, ok := r.Context().Value(utility.UserIdKey).(int)
	if !ok {
		err := model.ResponseError{
			Code:    http.StatusUnauthorized,
			Message: "Cannot retrieve user id",
		}
		logger.Error(err.Message)
		w.WriteHeader(err.Code)
		return
	}
	todoModels, err := db.GetTodosBy(userId)
	if err != nil {
		errResponse := model.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: "Cannot retrieve todo models",
		}
		logger.Error("Cannot retrieve todo models", zap.Error(err))
		w.WriteHeader(errResponse.Code)
		return
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(todoModels); err != nil {
		errResponse := model.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: "Error occurred during encoding",
		}
		logger.Error(errResponse.Message, zap.Error(err))
		w.WriteHeader(errResponse.Code)
		return
	}
}

// RemoveTodoHandler docs
// @Summary Remove todo by id
// @Tags todo
// @ID remove-todo-handler
// @Accept   json
// @Produce  json
// @Security ApiKeyAuth
// @param    Authorization header string true "Authorization"
// @Param    id      path   int     true  "todo id"
// @Success  200 {object} model.Response
// @Failure  401 {object} model.ResponseError
// @Failure  500 {object} model.ResponseError
// @Router   /todo/remove/{id} [delete]
func RemoveTodoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	vars := mux.Vars(r)
	todoId, err := strconv.Atoi(vars["id"])
	if err != nil {
		errResponse := model.ResponseError{
			Code:    http.StatusBadRequest,
			Message: "Cannot retrieve todo id",
		}
		logger.Error(errResponse.Message, zap.Error(err))
		w.WriteHeader(errResponse.Code)
		return
	}
	err = db.RemoveTodoBy(todoId)
	if err != nil {
		errResponse := model.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: "Cannot complete operation remove todo item",
		}
		logger.Error(errResponse.Message, zap.Error(err))
		w.WriteHeader(errResponse.Code)
		return
	}
	var response = model.Response{
		Code:    http.StatusOK,
		Message: fmt.Sprintf("Removed todo by %d", todoId),
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		errResponse := model.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: "Error occurred during encoding",
		}
		logger.Error(errResponse.Message, zap.Error(err))
		w.WriteHeader(errResponse.Code)
		return
	}
}

// AddTodoHandler docs
// @Summary add new todo to my todo list
// @Tags todo
// @ID add-todo-handler
// @Accept   json
// @Produce  json
// @Security ApiKeyAuth
// @param    Authorization header string true "Authorization"
// @Param    body      body   request.TodoForm     true  "form"
// @Success  200 {object} model.Todo
// @Failure  401 {object} model.ResponseError
// @Failure  500 {object} model.ResponseError
// @Router   /todo/add [post]
func AddTodoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	userId, ok := r.Context().Value(utility.UserIdKey).(int)
	if !ok {
		errResponse := model.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: "Cannot retrieve user id",
		}
		logger.Error(errResponse.Message)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var todoForm request.TodoForm
	if err := json.NewDecoder(r.Body).Decode(&todoForm); err != nil {
		errResponse := model.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: "Cannot retrieve todo form from request",
		}
		logger.Error(errResponse.Message, zap.Error(err))
		w.WriteHeader(errResponse.Code)
		return
	}
	todo, err := db.CreteTodoFor(userId, todoForm)
	if err != nil {
		errResponse := model.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: "Cannot complete operation add todo item",
		}
		logger.Error(errResponse.Message, zap.Error(err))
		w.WriteHeader(errResponse.Code)
		return
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(todo); err != nil {
		errResponse := model.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: "Error occurred during encoding",
		}
		logger.Error(errResponse.Message, zap.Error(err))
		w.WriteHeader(errResponse.Code)
	}
}

// GetTodoHandler docs
// @Summary Get todo by id
// @Tags todo
// @ID get-todo-handler
// @Accept   json
// @Produce  json
// @Security ApiKeyAuth
// @param    Authorization header string true "Authorization"
// @Param    id      path   int     true  "todo id"
// @Success  200 {object} model.Todo
// @Failure  401 {object} model.ResponseError
// @Failure  500 {object} model.ResponseError
// @Router   /todo/{id} [get]
func GetTodoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	vars := mux.Vars(r)
	todoId, err := strconv.Atoi(vars["id"])
	if err != nil {
		errResponse := model.ResponseError{
			Code:    http.StatusBadRequest,
			Message: "Cannot retrieve todo id",
		}
		logger.Error(errResponse.Message, zap.Error(err))
		w.WriteHeader(errResponse.Code)
		return
	}
	todoModel, err := db.GetTodoBy(todoId)
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(todoModel); err != nil {
		errResponse := model.ResponseError{
			Code:    http.StatusBadRequest,
			Message: "Error occurred during encoding",
		}
		logger.Error(errResponse.Message, zap.Error(err))
		w.WriteHeader(errResponse.Code)
		return
	}
}

// ToggleTodoHandler docs
// @Summary Toggle todo by id
// @Tags todo
// @ID toggle-todo-handler
// @Accept   json
// @Produce  json
// @Security ApiKeyAuth
// @param    Authorization header string true "Authorization"
// @Param    id      path   int     true  "todo id"
// @Success  200 {object} model.Todo
// @Failure  401 {object} model.ResponseError
// @Failure  500 {object} model.ResponseError
// @Router   /todo/toggle/{id} [put]
func ToggleTodoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	vars := mux.Vars(r)
	todoId, err := strconv.Atoi(vars["id"])
	if err != nil {
		errResponse := model.ResponseError{
			Code:    http.StatusBadRequest,
			Message: "Cannot retrieve todo id",
		}
		logger.Error(errResponse.Message, zap.Error(err))
		w.WriteHeader(errResponse.Code)
		return
	}
	if err := db.ToggleTodoFor(todoId); err != nil {
		errResponse := model.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: "Cannot toggle todo id",
		}
		logger.Error(errResponse.Message, zap.Error(err))
		w.WriteHeader(errResponse.Code)
		return
	}
	todo, err := db.GetTodoBy(todoId)
	if err != nil {
		errResponse := model.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: "Cannot retrieve todo model from db",
		}
		logger.Error(errResponse.Message, zap.Error(err))
		w.WriteHeader(errResponse.Code)
		return
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(todo); err != nil {
		errResponse := model.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: "Error occurred during encoding",
		}
		logger.Error(errResponse.Message, zap.Error(err))
		w.WriteHeader(errResponse.Code)
	}
}
