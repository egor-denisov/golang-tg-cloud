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
// Returning instance of database
func (db *DataBase) Get() *sql.DB{
	return db.database
}
// Creating new instance of database
func New() (DataBase, error) {
	// Create request string for connection
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", 
		env.HOST, env.PORT, env.USER, env.PASSWORD, env.DBNAME)
	// Connecting to database
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
        return DataBase{}, e.WrapIfErr("can`t create instance of database", err)
    }
	
	return DataBase{database: db}, nil
}
// Closing database
func (db *DataBase) Close() {
	db.database.Close()
}
// Executing insert request
func (db *DataBase) insert(req string, args ...interface{}) (string, error) {
	id := ""
	// Making a query
	if err := db.database.QueryRow(req, args...).Scan(&id); err != nil {
		return "", e.Wrap("can`t insert to database", err)
	}
	return id, nil
}
// Executing select request, which one row is expected
func (db *DataBase) selectRow(req string) (string, error) {
	res := ""
	// Making a query
	db.database.QueryRow(req).Scan(&res)
	return res, nil
}
// Executing select request, which many rows are expected
func (db *DataBase) selectRows(req string) (*sql.Rows, error) {
	return db.database.Query(req)
}
// Executing some query, which one don`t expect response
func (db *DataBase) makeQuery(req string) error {
	// Making a query
	_, err := db.database.Exec(req)
	return e.WrapIfErr("can`t make query", err)
}