package main

import (
	"flag"
	//"fmt"
	"log"
	"net/http"
	"path"

	"github.com/dimitarvalkanov7/chaoscamp-demo/controllers"
	//"github.com/dimitarvalkanov7/chaoscamp-demo/database"
	_ "github.com/lib/pq"
)

const (
	basePath string = "/home/leron/workspace/src/github.com/dimitarvalkanov7/chaoscamp-demo/"
)

var (
	addr = flag.String("addr", ":1234", "http service address")
)

func main() {
	// db := database.ConnectDB()

	// err := db.Ping()
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println("Successfully connected!")
	// db.Close()
	http.HandleFunc("/register", controllers.Register)

	fs := http.FileServer(http.Dir(path.Join(basePath, "static")))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	log.Println("Starting the HTTP server ...")
	log.Fatal(http.ListenAndServe(*addr, nil))
}
