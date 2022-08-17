package handler

import (
	"encoding/json"
	"github.com/IosifSuzuki/todo/internall/db"
	"github.com/IosifSuzuki/todo/internall/logger"
	"github.com/IosifSuzuki/todo/internall/model"
	"github.com/IosifSuzuki/todo/internall/model/request"
	"github.com/IosifSuzuki/todo/internall/utility"
	"go.uber.org/zap"
	"net/http"
)

// SignInHandler docs
// @Summary Sign in flow
// @Tags authentication
// @ID sign-in-handler
// @Accept   json
// @Produce  json
// @Param    body    body   request.AuthenticationForm     true  "form"
// @Success  200 {object} model.AccountModel
// @Failure  500 {object} model.ResponseError
// @Failure  401 {object} model.ResponseError
// @Failure  400 {object} model.ResponseError
// @Router   /authentication/sign-in [post]
func SignInHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	var authenticationForm *request.AuthenticationForm
	var err error
	if err = json.NewDecoder(r.Body).Decode(&authenticationForm); err != nil {
		logger.Error("occurred during decode body request", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var accountModel *model.AccountModel
	if accountModel, err = db.Authentication(*authenticationForm); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		logger.Error("occurred during check authentication", zap.Error(err))
		return
	}
	accessToken, err := utility.GenerateAccessToken(accountModel.Id, accountModel.UserName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Error("occurred during generate access token", zap.Error(err))
		return
	}
	refreshToken, err := utility.GenerateRefreshToken(accountModel.Id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Error("occurred during generate refresh token", zap.Error(err))
		return
	}
	credentials := model.Credentials{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(credentials); err != nil {
		logger.Error("occurred during encode response", zap.Error(err))
		return
	}
}

// SignUpHandler docs
// @Summary Sign up flow
// @Tags authentication
// @ID sign-up-handler
// @Accept   json
// @Produce  json
// @Param    body    body   request.RegistrationForm     true  "form"
// @Success  200 {object} model.AccountModel
// @Failure  500 {object} model.ResponseError
// @Failure  400 {object} model.ResponseError
// @Router   /authentication/sign-up [post]
func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	var registrationForm request.RegistrationForm
	if err := json.NewDecoder(r.Body).Decode(&registrationForm); err != nil {
		logger.Error("Occurred during decode body request", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !registrationForm.IsValidated() {
		logger.Error("Form request is not validated")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	logger.Info("Before create Account: ",
		zap.String("username", registrationForm.UserName),
		zap.String("email", registrationForm.Email),
		zap.String("password", registrationForm.Password),
	)
	accountModel, err := db.CreateAccount(registrationForm)
	if err != nil {
		logger.Error("During create account", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(accountModel); err != nil {
		logger.Error("Occurred when encode model to response", zap.Error(err))
	}
}

// RefreshTokenHandler docs
// @Summary Refresh token flow
// @Tags authentication
// @ID sign-token-handler
// @Accept   json
// @Produce  json
// @Param    body    body   model.Credentials     true  "form"
// @Success  200 {object} model.AccountModel
// @Failure  500 {object} model.ResponseError
// @Failure  400 {object} model.ResponseError
// @Router   /authentication/refresh-token [post]
func RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	var credentials model.Credentials
	var validToken bool
	var err error
	var accountModel *model.AccountModel
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		var errorModel = model.ResponseError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
		w.WriteHeader(errorModel.Code)
		logger.Error("occurred during decode body request", zap.Error(err))
		return
	}
	validToken, _ = utility.VerifyToken(credentials.AccessToken)
	if validToken {
		logger.Info("access token still valid", zap.String("access token", credentials.AccessToken))
		if err = json.NewEncoder(w).Encode(credentials); err != nil {
			logger.Error("Occurred when encode model to response", zap.Error(err))
		}
		w.WriteHeader(http.StatusOK)
		return
	}
	if validToken, err = utility.VerifyToken(credentials.RefreshToken); err != nil {
		var errorModel = model.ResponseError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
		w.WriteHeader(errorModel.Code)
		logger.Error("occurred during verification token", zap.Error(err))
		return
	} else if !validToken {
		logger.Debug("refresh token isn't valid", zap.String("refresh token", credentials.RefreshToken))
		var errorModel = model.ResponseError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
		w.WriteHeader(errorModel.Code)
		return
	}
	userId, err := utility.GetUserIdByFromToken(credentials.RefreshToken)
	if err != nil {
		var errorModel = model.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		w.WriteHeader(errorModel.Code)
		logger.Error("error occurred during fetch user id from refresh token",
			zap.Error(err),
		)
		return
	}
	if accountModel, err = db.GetUserById(userId); err != nil {
		var errorModel = model.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		w.WriteHeader(errorModel.Code)
		logger.Error("error occurred during fetch by user id from db",
			zap.Error(err),
		)
		return
	}

	accessToken, err := utility.GenerateAccessToken(accountModel.Id, accountModel.UserName)
	if err != nil {
		var errorModel = model.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		w.WriteHeader(errorModel.Code)
		logger.Error("occurred during generate access token", zap.Error(err))
		return
	}
	refreshToken, err := utility.GenerateRefreshToken(accountModel.Id)
	if err != nil {
		var errorModel = model.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		w.WriteHeader(errorModel.Code)
		logger.Error("occurred during generate refresh token", zap.Error(err))
		return
	}
	credentials = model.Credentials{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	if err := json.NewEncoder(w).Encode(credentials); err != nil {
		var errorModel = model.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		w.WriteHeader(errorModel.Code)
		logger.Error("occurred when encode model to response", zap.Error(err))
	}
}
