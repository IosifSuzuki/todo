package db

import (
	"context"
	"errors"
	"github.com/IosifSuzuki/todo/internall/logger"
	"github.com/IosifSuzuki/todo/internall/model"
	"github.com/IosifSuzuki/todo/internall/model/request"
	"github.com/IosifSuzuki/todo/internall/utility"
	"github.com/jackc/pgx/v4"
	"go.uber.org/zap"
)

var connectionDB *pgx.Conn

func ConnectToDB() {
	var urlConnection = utility.Config.DB.URL()
	connection, err := pgx.Connect(context.Background(), urlConnection)
	if err != nil {
		logger.Fatal("Error occurred during connection to db", zap.Error(err))
	}
	connectionDB = connection
}

func CloseConnectionToDB() error {
	logger.Info("Will disconnect from data base")
	return connectionDB.Close(context.Background())
}

func CreateAccount(registrationForm request.RegistrationForm) (*model.AccountModel, error) {
	hashPassword, err := utility.HashPassword(registrationForm.Password)
	if err != nil {
		return nil, err
	}
	var accountId int
	err = connectionDB.QueryRow(
		context.Background(),
		"INSERT INTO account (username, hash_password, email) VALUES($1, $2, $3) RETURNING id",
		registrationForm.UserName, hashPassword, registrationForm.Email,
	).Scan(&accountId)
	if err != nil {
		return nil, err
	}
	return GetUserById(accountId)
}

func Authentication(authenticationForm request.AuthenticationForm) (*model.AccountModel, error) {
	var accountModel *model.AccountModel
	var err error
	if accountModel, err = GetUserByUserName(authenticationForm.UserName); err != nil {
		return nil, err
	}
	if utility.CheckHashPassword(authenticationForm.Password, accountModel.HashPassword) {
		return nil, errors.New("access denied")
	}
	return accountModel, nil
}

func GetUserById(id int) (*model.AccountModel, error) {
	var accountModel = new(model.AccountModel)
	err := connectionDB.QueryRow(
		context.Background(),
		"SELECT id, username, email, created_on FROM account WHERE id = $1",
		id,
	).Scan(&accountModel.Id, &accountModel.UserName, &accountModel.Email, &accountModel.CreatedAt)
	return accountModel, err
}

func GetUserByUserName(userName string) (*model.AccountModel, error) {
	var accountModel = new(model.AccountModel)
	err := connectionDB.QueryRow(
		context.Background(),
		"SELECT id, username, email, created_on FROM account WHERE username = $1",
		userName,
	).Scan(&accountModel.Id, &accountModel.UserName, &accountModel.Email, &accountModel.CreatedAt)
	return accountModel, err
}
