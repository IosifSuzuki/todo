package handler

import (
	"encoding/json"
	"github.com/IosifSuzuki/todo/internall/logger"
	"github.com/IosifSuzuki/todo/internall/model"
	"github.com/IosifSuzuki/todo/internall/utility"
	"go.uber.org/zap"
	"net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value(utility.UserIdKey).(int)
	if !ok {
		logger.Error("Cannot retrieve user id")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var pingMessage = model.Ping{
		UserId:  userId,
		Message: "Success made request",
		Code:    http.StatusOK,
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(pingMessage); err != nil {
		logger.Error("Error occurred during encoding", zap.Error(err))
	}
}
