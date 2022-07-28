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

func RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	var credentials model.Credentials
	var validToken bool
	var err error
	var accountModel *model.AccountModel
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		logger.Error("occurred during decode body request", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	validToken, _ = utility.VerifyToken(credentials.AccessToken)
	if validToken {
		logger.Info("access token still valid", zap.String("access token", credentials.AccessToken))
		if err = json.NewEncoder(w).Encode(credentials); err != nil {
			logger.Error("Occurred when encode model to response", zap.Error(err))
		}
		return
	}
	if validToken, err = utility.VerifyToken(credentials.RefreshToken); err != nil {
		logger.Error("occurred during verification token", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else if !validToken {
		logger.Debug("refresh token isn't valid", zap.String("refresh token", credentials.RefreshToken))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	userId, err := utility.GetUserIdByFromToken(credentials.RefreshToken)
	if err != nil {
		logger.Error("error occurred during fetch user id from refresh token",
			zap.Error(err),
		)
		return
	}
	if accountModel, err = db.GetUserById(userId); err != nil {
		logger.Error("error occurred during fetch by user id from db",
			zap.Error(err),
		)
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
	credentials = model.Credentials{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	if err := json.NewEncoder(w).Encode(credentials); err != nil {
		logger.Error("occurred when encode model to response", zap.Error(err))
	}
}
