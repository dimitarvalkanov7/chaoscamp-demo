package main

import (
	"flag"
	"github.com/dimitarvalkanov7/chaoscamp-demo/routing"
	//"fmt"
	"log"
	"net/http"
	//"path"
	"time"

	"github.com/dimitarvalkanov7/chaoscamp-demo/controllers"
	//"github.com/dimitarvalkanov7/chaoscamp-demo/database"
	_ "github.com/lib/pq"
)

func init() {
	err := controllers.SeedUsers()
	if err != nil {
		log.Fatal(err)
	}
}

const (
	basePath string = "/home/leron/workspace/src/github.com/dimitarvalkanov7/chaoscamp-demo/"
)

func main() {
	// db := database.ConnectDB()

	// err := db.Ping()
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println("Successfully connected!")
	// db.Close()
	// http.HandleFunc("/register", controllers.Register)

	// fs := http.FileServer(http.Dir(path.Join(basePath, "static")))
	// http.Handle("/static/", http.StripPrefix("/static/", fs))
	// log.Println("Starting the HTTP server ...")
	// log.Fatal(http.ListenAndServe(*addr, nil))
	var dir string

	flag.StringVar(&dir, "dir", ".", "the directory to serve files from. Defaults to the current dir")
	flag.Parse()
	r := routing.Handlers()

	// This will serve files under http://localhost:1234/static/<filename>
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(dir))))

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:1234",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
