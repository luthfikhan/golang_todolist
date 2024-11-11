package app

import (
	"database/sql"
	"os"
	"todolist/helper"
	"todolist/model"

	"github.com/gin-gonic/gin"
	go_ora "github.com/sijms/go-ora/v2"
)

func GetOracleDb() *sql.DB {
	connectionString := os.Getenv("ORACLE_CONNECTION")
	isTest := os.Getenv("GO_ENV") == "test"

	if isTest {
		connectionString = os.Getenv("ORACLE_CONNECTION_TEST")
	}

	db, err := sql.Open("oracle", connectionString)
	helper.PanicIfError(err)

	err = db.Ping()
	helper.PanicIfError(err)

	go_ora.RegisterType(db, "test", "", model.Task{})

	helper.Log.Info(gin.H{"db": "oracle"}, "Database connected")
	return db
}
