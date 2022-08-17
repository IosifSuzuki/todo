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

func GetTodosBy(userId int) ([]model.Todo, error) {
	rows, err := connectionDB.Query(context.Background(), "SELECT id, title, description, created_on, updated_on, closed "+
		"FROM item INNER JOIN account_item ON account_item.item_id = id "+
		"WHERE account_item.account_id = $1 ORDER BY updated_on", userId,
	)
	defer rows.Close()
	var todos = make([]model.Todo, 0)
	if err != nil {
		return todos, err
	}
	for rows.Next() {
		var todoModel = model.Todo{}
		err = rows.Scan(
			&todoModel.Id,
			&todoModel.Title,
			&todoModel.Description,
			&todoModel.CreatedOn,
			&todoModel.UpdatedOn,
			&todoModel.Closed,
		)
		if err != nil {
			return todos, err
		}
		todos = append(todos, todoModel)
	}
	return todos, err
}

func CreteTodoFor(userId int, todoForm request.TodoForm) (*model.Todo, error) {
	var todo = &model.Todo{}
	tx, err := connectionDB.BeginTx(context.Background(), pgx.TxOptions{})
	if err != nil {
		return todo, err
	}
	defer func() {
		if err != nil {
			tx.Rollback(context.Background())
		} else {
			tx.Commit(context.Background())
		}
	}()
	var todoId int
	err = tx.QueryRow(context.Background(),
		"INSERT INTO item(title, description, closed) VALUES($1, $2, $3) RETURNING id",
		todoForm.Title, todoForm.Description, false,
	).Scan(&todoId)
	if err != nil {
		return todo, err
	}
	_, err = tx.Exec(context.Background(), "INSERT INTO account_item (account_id, item_id) VALUES ($1, $2)",
		userId, todoId)
	if err != nil {
		return todo, err
	}
	err = tx.QueryRow(context.Background(), "SELECT id, title, description, created_on, updated_on, closed "+
		"FROM item WHERE id = $1", todoId,
	).Scan(&todo.Id, &todo.Title, &todo.Description, &todo.CreatedOn, &todo.UpdatedOn, &todo.Closed)
	return todo, err
}

func ToggleTodoFor(todoId int) error {
	var closed bool
	err := connectionDB.QueryRow(context.Background(), "SELECT closed FROM item WHERE id = $1", todoId).Scan(&closed)
	if err != nil {
		return err
	}
	_, err = connectionDB.Exec(context.Background(), "UPDATE item SET closed = $1 WHERE id = $2", !closed, todoId)
	return err
}

func RemoveTodoBy(todoId int) error {
	_, err := connectionDB.Exec(context.Background(), "DELETE FROM item WHERE id = $1", todoId)
	return err
}

func GetTodoBy(todoId int) (*model.Todo, error) {
	var todo = &model.Todo{}
	err := connectionDB.QueryRow(context.Background(), "SELECT id, title, description, created_on, updated_on, closed "+
		"FROM item WHERE id = $1", todoId,
	).Scan(&todo.Id, &todo.Title, &todo.Description, &todo.CreatedOn, &todo.UpdatedOn, &todo.Closed)
	return todo, err
}

func GetAccountBy(userId int) (*model.AccountModel, error) {
	var account = &model.AccountModel{}
	err := connectionDB.QueryRow(
		context.Background(),
		"SELECT id, username, email, created_on FROM account WHERE id = $1",
		userId,
	).Scan(&account.Id, &account.UserName, &account.Email, &account.CreatedAt)
	return account, err
}

func GetAccounts() ([]model.AccountModel, error) {
	var accounts = make([]model.AccountModel, 0)
	rows, err := connectionDB.Query(context.Background(), "SELECT id, username, email, created_on FROM account")
	if err != nil {
		return accounts, err
	}
	for rows.Next() {
		var account = model.AccountModel{}
		err = rows.Scan(
			&account.Id,
			&account.UserName,
			&account.Email,
			&account.CreatedAt,
		)
		if err != nil {
			return accounts, err
		}
		accounts = append(accounts, account)
	}
	return accounts, err
}
