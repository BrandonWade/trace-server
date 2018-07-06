package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/BrandonWade/contact"
	"github.com/BrandonWade/godash"
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
	clientFiles := []string{}
	for {
		msg := contact.Message{}

		conn.ReadJSON(&msg)
		if msg.IsEmpty() {
			break
		}

		path := filepath.ToSlash(msg.Body)
		clientFiles = append(clientFiles, path)
	}

	// Get the list of files from the filesystem
	localFiles, err := synth.Scan(os.Getenv("TEST_DIR"))
	if err != nil {
		log.Fatal("error retrieving local file list")
	}

	// Filter out unwanted files and files that are already on the client
	// TODO: Add support for setting filters
	filters := []string{"xyz"}
	filters = append(filters, clientFiles...)

	newFiles := godash.DifferenceSubstr(localFiles, filters)

	// Send the list of new files to the client
	for _, file := range newFiles {
		conn.Write(file)
	}
}
