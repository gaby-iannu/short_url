package main

import (
	"short_url/builder"
	"short_url/handler"
)


func main() {
	router := handler.InitializeAndRun(builder.Build())
	router.Run("localhost:8080")
}
