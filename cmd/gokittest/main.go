package main

import (
	"log"
	"net/http"

	"local/gokit-test/internal"
)

func main() {
	log.Fatal(http.ListenAndServe(":8081", internal.NewHTTPHandler()))
}
