package main

import (
	"fmt"
	"go-book-api/src/router"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Running API")

	r := router.Generate()

	log.Fatal(http.ListenAndServe(":5000", r))
}
