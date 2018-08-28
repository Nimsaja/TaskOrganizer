package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"cloud.google.com/go/storage"
	"github.com/gorilla/mux"
	"google.golang.org/appengine"
	"google.golang.org/appengine/file"
)

var path = "output.txt"
var inCloud bool

// handle CORS and the OPION method
func corsAndOptionHandler(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
		} else {
			h.ServeHTTP(w, r)
		}
	}
}

// create all used Handler
func handler() http.Handler {
	router := mux.NewRouter()

	url := "/tasks"
	router.HandleFunc(url, taskList).Methods("GET")
	router.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "invalid method: "+r.Method, http.StatusBadRequest)
	}).Methods("DELETE", "PATH", "COPY", "HEAD", "LINK", "UNLINK", "PURGE", "LOCK", "UNLOCK", "VIEW", "PROPFIND")

	url = "/clean"
	router.HandleFunc(url, cleanTaskList).Methods("GET")

	return corsAndOptionHandler(router)
}

func main() {
	inCloud, _ = strconv.ParseBool(os.Getenv("RUN_IN_CLOUD"))

	http.Handle("/", handler())

	appengine.Main()
	// log.Println("Init is ready and start the server on: http://localhost:8080")

	// log.Fatalln(http.ListenAndServe(":8080", nil))
}

func taskList(w http.ResponseWriter, r *http.Request) {
	text := r.FormValue("text")

	if inCloud {
		//in cloud no append is possible, so first read in the file, append the new text and write the whole stuff into the file
		data, err := read(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data = append(data, text)
		writeToClient(w, data)

		err = writeToCloudStorage(r, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {

		appendToOutputFile(text)

		readFile(w)
	}
}

func cleanTaskList(w http.ResponseWriter, r *http.Request) {
	if inCloud {
		data := make([]string, 0)
		writeToCloudStorage(r, data)
	} else {
		//not necessary, just do rm output.txt!!!
	}
}

func appendToOutputFile(t string) {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Printf("Error %s ", err)
		panic(err)
	}

	defer f.Close()

	fmt.Fprintln(f, t)
}

func readFile(w http.ResponseWriter) {
	f, err := os.OpenFile(path, os.O_RDONLY, 0600)
	if err != nil {
		log.Printf("Error %s ", err)
		panic(err)
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		b, err := json.Marshal(scanner.Text())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "%s\n", string(b))
	}
}

func writeToCloudStorage(r *http.Request, text []string) error {
	fileName := "file.txt"

	ctx := appengine.NewContext(r)

	// determine default bucket name
	bucketName, err := file.DefaultBucketName(ctx)
	if err != nil {
		log.Fatalf("failed to get default GCS bucket name: %v", err)
		return err
	}

	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("failed to get default GCS bucket name: %v", err)
		return err
	}
	defer client.Close()

	bucket := client.Bucket(bucketName)

	wc := bucket.Object(fileName).NewWriter(ctx)
	wc.ContentType = "text/plain"

	for _, element := range text {
		if _, err := wc.Write([]byte(element)); err != nil {
			log.Fatalf("createFile: unable to write data to bucket %v, file %q: %v", bucket, fileName, err)
			return err
		}
		if _, err := wc.Write([]byte("\n")); err != nil {
			log.Fatalf("createFile: unable to write data to bucket %v, file %q: %v", bucket, fileName, err)
			return err
		}
	}

	if err := wc.Close(); err != nil {
		log.Fatalf("createFile: unable to close bucket %v, file %q: %v", bucket, fileName, err)
		return err
	}

	return nil
}

func read(w http.ResponseWriter, r *http.Request) ([]string, error) {
	fileName := "file.txt"

	ctx := appengine.NewContext(r)

	// determine default bucket name
	bucketName, err := file.DefaultBucketName(ctx)
	if err != nil {
		log.Fatalf("failed to get default GCS bucket name: %v", err)
		return nil, err
	}

	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("failed to get default GCS bucket name: %v", err)
		return nil, err
	}
	defer client.Close()

	bucket := client.Bucket(bucketName)

	rc, err := bucket.Object(fileName).NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	scanner := bufio.NewScanner(rc)
	stringArray := make([]string, 0)
	for scanner.Scan() {
		stringArray = append(stringArray, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return stringArray, nil
}

func writeToClient(w http.ResponseWriter, data []string) {
	for _, element := range data {
		fmt.Fprintf(w, "%s # \n", element)
	}
}
