package db

import (
	"database/sql"
	"fmt"
	env "main/environment"
	"main/lib/e"

	_ "github.com/lib/pq"
)

type DataBase struct {
	database *sql.DB
}

func (db *DataBase) Get() *sql.DB{
	return db.database
}

func New() (DataBase, error) {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", 
		env.HOST, env.PORT, env.USER, env.PASSWORD, env.DBNAME)

	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
        return DataBase{}, e.WrapIfErr("can`t create instance of database", err)
    }

	return DataBase{database: db}, nil
}

func (db *DataBase) Close() {
	db.database.Close()
}

func (db *DataBase) insert(req string, args ...interface{}) (string, error) {
	id := ""
	err := db.database.QueryRow(req, args...).Scan(&id)
	if err != nil {
		return "", e.Wrap("can`t insert to database", err)
	}
	return id, nil
}

func (db *DataBase) selectRow(req string) (string, error) {
	res := ""
	db.database.QueryRow(req).Scan(&res)
	return res, nil
}

func (db *DataBase) selectRows(req string) (*sql.Rows, error) {
	return db.database.Query(req)
}

func (db *DataBase) makeQuery(req string) error {
	_, err := db.database.Exec(req)
	if err != nil {
		return e.Wrap("can`t make query", err)
	}
	return nil
}