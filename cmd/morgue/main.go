package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type readOp struct {
	uri  string
	resp chan string
}

type writeOp struct {
	uri  string
	body string
	resp chan bool
}

var reads = make(chan readOp)
var writes = make(chan writeOp)

func bodiesStorage() {

	var bodies = make(map[string]string)

	for {
		select {
		case read := <-reads:
			read.resp <- bodies[read.uri]
		case write := <-writes:
			bodies[write.uri] = write.body
			write.resp <- true
		}
	}
}

func showBody(w http.ResponseWriter, r *http.Request) {

	requestURI := r.URL.RequestURI()
	log.Printf("show %s", requestURI)

	w.Header().Set("Content-Type", "application/json")

	read := readOp{
		uri:  requestURI,
		resp: make(chan string)}
	reads <- read

	io.WriteString(w, <-read.resp)
}

func saveBody(w http.ResponseWriter, r *http.Request) {

	requestURI := r.URL.RequestURI()
	log.Printf("save %s", requestURI)

	w.Header().Set("Content-Type", "application/json")

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("%v", err)
	}

	write := writeOp{
		uri:  requestURI,
		body: string(body),
		resp: make(chan bool)}
	writes <- write
	<-write.resp

	io.WriteString(w, `{"message": "body saved"}`)
}

func NewRouter() *mux.Router {

	router := mux.NewRouter()

	router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		saveBody(w, r)
	}).Methods("PUT")

	router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		showBody(w, r)
	})

	return router
}

func main() {

	port := "8080"
	if portEnv := os.Getenv("PORT"); portEnv != "" {
		port = portEnv
	}

	go bodiesStorage()

	router := NewRouter()

	log.Printf(fmt.Sprintf("listening on port %s ...", port))
	log.Fatal(http.ListenAndServe(":"+port, router))
}
