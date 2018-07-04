package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/BrandonWade/contact"
	"github.com/BrandonWade/synth"
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
	port := ":" + os.Getenv("TRACE_SERVER_PORT")

	r := mux.NewRouter()
	r.HandleFunc("/sync", SyncHandler)

	http.ListenAndServe(port, r)
}

// SyncHandler - Handler for incoming sync requests
func SyncHandler(w http.ResponseWriter, r *http.Request) {
	conn := contact.NewConnection(bufferSize)

	conn.Open(&w, r)
	defer conn.Close()

	// Get the list of files from the client
	clientFiles := make(map[string]bool)
	for {
		msg := contact.Message{}

		conn.ReadJSON(&msg)
		if msg.IsEmpty() {
			break
		}

		relPath := filepath.ToSlash(msg.Body)
		clientFiles[relPath] = true
	}

	// Get the list of files from the filesystem
	localFiles, err := synth.Scan(os.Getenv("TEST_DIR"))
	if err != nil {
		log.Fatal("error retrieving local file list")
	}

	fmt.Printf("%+v", localFiles)
}
