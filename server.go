package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
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
	// TODO: Implement
	w.Write([]byte("Trace"))
}
