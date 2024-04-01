package main

import (
	"short_url/builder"
	"short_url/handler"
)


func main() {
	handler.InitializeAndRun(builder.Build())
}
