package main

import (
	"fmt"
	"short_url/short"
	"short_url/short/cache"
	"short_url/short/model"
	"short_url/short/repository"
)


func main() {
	long := "https://medium.com/@sandeep4.verma/system-design-scalable-url-shortener-service-like-tinyurl-106f30f23a82"
	short := short.New(cache.New("localhost", 6379), repository.New())
	url := model.New(long, "yo")
	
	shorturl := short.Tiny(url)
	
	fmt.Printf("Short URL: %s\n", shorturl)

	longUrl, err:=short.Get(shorturl)
	if err == nil {
		fmt.Printf("Long URL: %s\n", longUrl)
	}

	_, err= short.Get("no existe")
	if err != nil {
		fmt.Println(err)
	}
}
