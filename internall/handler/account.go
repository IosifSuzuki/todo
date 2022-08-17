package handler

import (
	"encoding/json"
	_ "github.com/IosifSuzuki/todo/docs"
	"github.com/IosifSuzuki/todo/internall/db"
	"github.com/IosifSuzuki/todo/internall/logger"
	"github.com/IosifSuzuki/todo/internall/model"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

// UserInfoHandler docs
// @Summary Get account info
// @Description get account info by id
// @Tags account
// @ID get-user-info
// @Accept   json
// @Produce  json
// @Security ApiKeyAuth
// @param    Authorization header string true "Authorization"
// @Param    id      path   int     true  "account id"
// @Success  200 {object} model.AccountModel
// @Failure  500 {object} model.ResponseError
// @Failure  400 {object} model.ResponseError
// @Router   /account/user/{id} [get]
func UserInfoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	vars := mux.Vars(r)
	accountId, err := strconv.Atoi(vars["id"])
	if err != nil {
		var errorModel = model.ResponseError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
		w.WriteHeader(errorModel.Code)
		logger.Error("Cannot parse user id from request", zap.Error(err))
	}
	accountModel, err := db.GetAccountBy(accountId)
	if err != nil {
		var errorModel = model.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		w.WriteHeader(errorModel.Code)
		logger.Error("During make query to db", zap.Error(err))
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(accountModel); err != nil {
		var errorModel = model.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		w.WriteHeader(errorModel.Code)
		logger.Error("During encode account model to response", zap.Error(err))
	}
}

// UsersInfoHanlder docs
// @Summary Get accounts info
// @Description get accounts info
// @Tags account
// @ID users-info-hanlder
// @Accept   json
// @Produce  json
// @Security ApiKeyAuth
// @param    Authorization header string true "Authorization"
// @Param    id      path   int     true  "account id"
// @Success  200 {array} model.AccountModel
// @Failure  500 {object} model.ResponseError
// @Failure  400 {object} model.ResponseError
// @Failure  401 {object} model.ResponseError
// @Router   /account/users/ [get]
func UsersInfoHanlder(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	accounts, err := db.GetAccounts()
	if err != nil {
		var errorModel = model.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		w.WriteHeader(errorModel.Code)
		logger.Error("After make quetion to db", zap.Error(err))
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(accounts); err != nil {
		var errorModel = model.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		logger.Error("During encode account model to response", zap.Error(err))
		w.WriteHeader(errorModel.Code)
	}
}
