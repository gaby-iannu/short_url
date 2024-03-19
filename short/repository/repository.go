package repository

import (
	"context"
	"database/sql"
	"fmt"
	"short_url/short/model"
	"time"

	_ "github.com/go-sql-driver/mysql"
)


var urls map[string]model.Url= make(map[string]model.Url)

type Repository interface {
	InsertIfNotExists(tiny string, url model.Url) bool
	Read(tinyUrl string) model.Url
}

type repository struct {}

func New() Repository {
	return &repository{}
}

// Connect to mysql, return db connection
func connectToMysql() *sql.DB {
	username:="root"
	password:="root"
	hostname:="127.0.0.1:3306"
	dbName:="shorturl"
	maxOpenConns:=20
	maxIdleConns:=20
	connMaxLifetime:= time.Minute * 3

	datasource :=fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostname, dbName)
	db, err:=sql.Open("mysql", datasource)
	if err != nil {
		panic(err)
	}
	
	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxLifetime(connMaxLifetime)
	return db
}

const (
	query = "select LONG_URL from URL where TINY_URL = ?"
	query_all = "select LONG_URL, USER_ID from URL where TINY_URL = ?"
	insert = "insert into URL (TINY_URL,LONG_URL,USER_ID) values (?,?,?)"
)

// true si lo inserta 
// false si ya existe y no lo inserta
func (r *repository) InsertIfNotExists(tiny string, url model.Url) bool {
	db := connectToMysql()
	defer db.Close()

	rows, err := db.Query(query, tiny)
	if err != nil {
		panic(err)
	}

	defer rows.Close()
	
	if rows.Next() {
		var longUrl string
		err := rows.Scan(&longUrl)
		if err != nil {
			panic(err)
		}

		return false 
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	stmt, err := db.PrepareContext(ctx, insert)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()
	result, err := stmt.ExecContext(ctx, tiny, url.Long, url.User)
	if err != nil {
		panic(err)
	}

	row, err := result.RowsAffected()
	if err != nil {
		panic(err)
	}

	return row > 0
}

// Return url struct with long url
// and user
func (r *repository) Read(tinyUrl string) model.Url {
	db := connectToMysql()
	defer db.Close()

	rows, err := db.Query(query_all, tinyUrl)
	if err != nil {
		panic(err)
	}

	var longUrl, userId string
	for rows.Next() {
		err = rows.Scan(&longUrl, &userId)
		if err != nil {
			panic(err)
		}
		rows.Close()
		return model.Url{
			Long:longUrl,
			User: userId,
		}
	}

	return model.Url{}
}
