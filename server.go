package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/BrandonWade/contact"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var (
	bufferSize int
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	bufferSize, err = strconv.Atoi(os.Getenv("TRACE_BUFFER_SIZE"))
	if err != nil {
		log.Fatal("error reading buffer size")
	}
}

func main() {
	port := os.Getenv("TRACE_SERVER_PORT")

	r := mux.NewRouter()
	r.HandleFunc("/sync", SyncHandler)

	http.ListenAndServe(port, r)
}

// SyncHandler - Handler for incoming sync requests
func SyncHandler(w http.ResponseWriter, r *http.Request) {
	conn := contact.NewConnection(bufferSize)
	conn.Open(&w, r)
}
