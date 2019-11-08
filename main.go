package main

import (
	"log"
	"net/http"
)

func main() {
	r := InitRouter()
	log.Fatal(http.ListenAndServe(":8089", r))
}
