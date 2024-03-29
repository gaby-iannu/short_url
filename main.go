package main

import (
	"fmt"
	"short_url/short/repository"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func buildDataSource() repository.DataSource {
	return repository.DataSource{
		Username:"root",
		Password:"root",
		Hostname:"127.0.0.1:3306",
		DBName:"shorturl",
		MaxOpenConns:20,
		MaxIdleConns:20,
		ConnMaxLifetime: time.Minute * 3,
	}
}

func test() {
	db, mock, err := sqlmock.New()
	if err != nil {
		fmt.Println("failed to open sqlmock database:", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "title"}).
		AddRow(0, "two").
		RowError(0, fmt.Errorf("row error"))
	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	rs, _ := db.Query("SELECT")
	defer rs.Close()

	if rs.Next() {
		var id int
		var title string
		rs.Scan(&id, &title)
		fmt.Println("scanned id:", id, "and title:", title)
	}

	if rs.Err() != nil {
		fmt.Println("got rows error:", rs.Err())
	}
}

func main() {
	test()
	// long := "https://medium.com/@sandeep4.verma/system-design-scalable-url-shortener-service-like-tinyurl-106f30f23a82"
	// short := short.New(cache.New("localhost", 6379), repository.New(buildDataSource()))
	// url := model.New(long, "yo")
	
	// shorturl := short.Tiny(url)
	
	// fmt.Printf("Short URL: %s\n", shorturl)

	// longUrl, err:=short.Get(shorturl)
	// if err == nil {
	// 	fmt.Printf("Long URL: %s\n", longUrl)
	// }

	// _, err= short.Get("no existe")
	// if err != nil {
	// 	fmt.Printf("Error:%v\n",err)
	// }
}
