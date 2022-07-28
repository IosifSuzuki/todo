package model

import "fmt"

type DBConfig struct {
	UserName string
	Password string
	DBName   string
	DBHost   string
}

const testMode = true

func (d *DBConfig) URL() string {
	if testMode {
		return fmt.Sprintf("postgres://%s:%s@localhost:5432/%s?sslmode=disable", d.UserName, d.Password, d.DBName)
	} else {
		return fmt.Sprintf("postgres://%s:%s@db:5432/%s?sslmode=disable", d.UserName, d.Password, d.DBName)
	}
}
