package main

import (
	"fmt"
	"log"
	"os"
	"short_url/builder"
	"short_url/handler"
)

func configureEnvironment() (string, string, string) {
	shorturlHost, mysqlHost, redisHost := "localhost","127.0.0.1","localhost"
	scope := os.Getenv("SCOPE")
	if len(scope) > 0 && scope == "prod" {
		shorturlHost, mysqlHost, redisHost = os.Getenv("SHORT_HOST"), os.Getenv("BD_HOST") ,os.Getenv("CACHE_HOST")

		if len(shorturlHost) == 0 {
			panic(fmt.Errorf("short host is empty"))
		}

		if len(mysqlHost) == 0 {
			panic(fmt.Errorf("mysql host is empty"))
		}

		if len(redisHost) == 0 {
			panic(fmt.Errorf("redis host is empty"))
		}
	}

	log.Printf("SCOPE: %s, shorturlHost: %s, mysqlHost: %s, redisHost: %s\n", scope, shorturlHost, mysqlHost, redisHost)
	return shorturlHost, mysqlHost, redisHost
}

func main() {
	shorturlHost, mysqlHost, redisHost := configureEnvironment()
	repository, cache, err := builder.Build(mysqlHost, redisHost)
	if err != nil {
		panic(err)
	}
	router := handler.InitializeAndRun(repository, cache)
	router.Run(fmt.Sprintf("%s:%s", shorturlHost, "8080"))
}
