package main

import (
	"fmt"
	"go-book-api/src/router"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Server running")

	r := router.Generate()

	log.Fatal(http.ListenAndServe(":5000", r))
}
