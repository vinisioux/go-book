package main

import (
	"fmt"
	"go-book-api/src/config"
	"go-book-api/src/router"
	"log"
	"net/http"
)

func main() {
	config.LoadConfigs()

	r := router.Generate()

	fmt.Printf("Server running on port: %d\n", config.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Port), r))
}
