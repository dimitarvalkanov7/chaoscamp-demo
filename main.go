package main

import (
	"github.com/dimitarvalkanov7/chaoscamp-demo/routing"
	"log"
	"net/http"
	"time"

	"github.com/dimitarvalkanov7/chaoscamp-demo/controllers"
	_ "github.com/lib/pq"
)

func init() {
	err := controllers.SeedUsers()
	if err != nil {
		log.Fatal(err)
	}
}

var (
	homepath string = "/home/leron/workspace/src/github.com/dimitarvalkanov7/chaoscamp-demo/"
)

func main() {
	srv := &http.Server{
		Handler: routing.Handlers(),
		Addr:    "127.0.0.1:1234",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
