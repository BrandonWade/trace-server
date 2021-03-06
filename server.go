package main

import (
	"bufio"
	"io"
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

	syncDir = os.Getenv("SYNC_DIR")
	if syncDir == "" {
		log.Fatal("error reading sync directory")
	}
}

func main() {
	port := ":" + os.Getenv("TRACE_SERVER_PORT")

	r := mux.NewRouter()
	r.HandleFunc("/sync", syncHandler).Methods("GET")
	r.HandleFunc("/download", downloadHandler).Methods("GET")

	http.ListenAndServe(port, r)
}

// syncHandler - handler for incoming sync requests
func syncHandler(w http.ResponseWriter, r *http.Request) {
	conn := contact.NewConnection(bufferSize)

	conn.Open(&w, r)
	defer conn.Close()

	// Get the list of files from the client
	clientFiles := []*synth.File{}
	for {
		file := synth.File{}
		conn.ReadJSON(&file)

		if file.IsEmpty() {
			break
		}

		file.Path = filepath.ToSlash(file.Path)
		clientFiles = append(clientFiles, &file)
	}

	// Get the list of files from the filesystem
	localFiles, err := synth.Scan(syncDir)
	if err != nil {
		log.Println("error retrieving local file list")
		return
	}

	synth.TrimPaths(localFiles, syncDir)

	// Filter out unwanted files and files that are already on the client
	// TODO: Add support for setting filters
	filters := []string{}
	clientFilePaths := synth.GetPathSlice(clientFiles)
	filters = append(filters, clientFilePaths...)

	// Determine which files are new
	newFiles := simpleDiff(&localFiles, &filters)
	newFiles = append(newFiles, &synth.File{})

	// Send the list of new files to the client
	for _, file := range newFiles {
		conn.WriteJSON(file)
	}
}

// downloadHandler - handler for uncoming file download requests
func downloadHandler(w http.ResponseWriter, r *http.Request) {
	conn := contact.NewConnection(bufferSize)
	conn.Open(&w, r)

	// Retrieve file name from the request
	file := r.URL.Query()["file"][0]

	// Send the file to the client
	go sendFile(conn, file)
}

// sendFile - sends the contents of a file over a websocket connection
func sendFile(conn *contact.Connection, file string) {
	defer conn.Close()
	filePath := syncDir + file

	filePtr, err := os.Open(filePath)
	if err != nil {
		log.Printf("error opening file %s\n", filePath)
		return
	}
	defer filePtr.Close()

	// Read the contents of the file and send them over the connection in chunks
	buffer := bufio.NewReader(filePtr)
	for {
		// Read the contents of the file one block at a time
		block := make([]byte, bufferSize)
		_, err := buffer.Read(block)
		if err != nil {
			if err != io.EOF {
				log.Printf("error reading contents from file %s\n", filePath)
				return
			}

			break
		}

		// Write the current block to the client
		conn.WriteBinary(block)
	}

	// Send an empty message to indicate the end of the file
	conn.Write("")
}
