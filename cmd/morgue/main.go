package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type readBody struct {
	uri  string
	resp chan Response
}

type storeBody struct {
	response Response
	resp     chan bool
}

type Response struct {
	uri       string
	body      string
	sleepTime int
}

var reads = make(chan readBody)
var writes = make(chan storeBody)

func bodiesStorage() {

	var bodies = make(map[string]Response)

	for {
		select {
		case read := <-reads:
			read.resp <- bodies[read.uri]
		case write := <-writes:
			bodies[write.response.uri] = write.response
			write.resp <- true
		}
	}
}

func showBody(w http.ResponseWriter, r *http.Request) {

	requestURI := r.URL.RequestURI()
	log.Printf("show %s", requestURI)

	w.Header().Set("Content-Type", "application/json")

	read := readBody{
		uri:  requestURI,
		resp: make(chan Response)}
	reads <- read

	response := <-read.resp

	time.Sleep(time.Duration(response.sleepTime) * time.Millisecond)

	io.WriteString(w, response.body)
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

	sleepTime := 0
	auxSleepTime, err := strconv.Atoi(r.Header.Get("morgue-set-sleep-time"))
	if err == nil {
		sleepTime = auxSleepTime
	}

	write := storeBody{
		response: Response{
			uri:       requestURI,
			sleepTime: sleepTime,
			body:      string(body),
		},
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
