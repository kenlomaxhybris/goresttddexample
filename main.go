package main

import (
	"log"
	"net/http"

	"github.com/kenlomaxhybris/resttddexample/endpoint"
)

func main() {
	r := endpoint.InitRouter()
	log.Fatal(http.ListenAndServe(":8089", r))
}
