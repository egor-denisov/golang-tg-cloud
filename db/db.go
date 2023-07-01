package db

import (
	"database/sql"
	"fmt"
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
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
        return DataBase{}, e.WrapIfErr("can`t create instance of database", err)
    }

	return DataBase{database: db}, nil
}

func (db *DataBase) Close() {
	db.database.Close()
}

func (db *DataBase) Insert(req string, args ...interface{}) (string, error) {
	id := ""
	err := db.database.QueryRow(req, args...).Scan(&id)
	if err != nil {
		return "", e.Wrap("can`t insert to database", err)
	}
	return id, nil
}

func (db *DataBase) GetId(req string) (string, error) {
	id := ""
	db.database.QueryRow(req).Scan(&id)
	return id, nil
}