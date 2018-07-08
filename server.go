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
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

var (
	bufferSize int
	syncDir    string
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

	syncDir = os.Getenv("TEST_DIR")
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
		_, msg, err := conn.Read()
		if err != nil {
			if ce, ok := err.(*websocket.CloseError); ok {
				if ce.Code != websocket.CloseNormalClosure {
					log.Println("error reading client files from connection")
					return
				}
			}
		}

		data := string(msg)
		if data == "" {
			break
		}

		path := filepath.ToSlash(data)
		clientFiles = append(clientFiles, path)
	}

	// Get the list of files from the filesystem
	localFiles, err := synth.Scan(syncDir)
	if err != nil {
		log.Println("error retrieving local file list")
		return
	}

	localFiles = synth.TrimPaths(localFiles, syncDir)

	// Filter out unwanted files and files that are already on the client
	// TODO: Add support for setting filters
	filters := []string{"xyz"}
	filters = append(filters, clientFiles...)

	// Add an empty element to the end of the list to indicate the end
	newFiles := godash.DifferenceSubstr(localFiles, filters)
	newFiles = append(newFiles, "")

	// Send the list of new files to the client
	for _, file := range newFiles {
		conn.Write(file)
	}
}
