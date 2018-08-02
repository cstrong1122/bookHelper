package main

import (
	"crypto/tls"
	"encoding/json"
	"log"
	"net"
	"os"
	"time"

	"net/http"

	"github.com/cstrong1122/configHelper"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
)

// Book struct (Model)
type Book struct {
	ID     string  `json:"id"`
	Isbn   string  `json:"isbn"`
	Title  string  `json:"title"`
	Author *Author `json:"author"`
}

// Author ...
type Author struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

// DataStore ...
type DataStore struct {
	Session *mgo.Session
}

var dataStore DataStore

func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	books := GetAllBooks(&dataStore)
	if books != nil {
		json.NewEncoder(w).Encode(&books)
		return
	}
	json.NewEncoder(w).Encode([]Book{})
}

func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	book := GetBookById(&dataStore, params["id"])
	if book != nil {
		json.NewEncoder(w).Encode(book)
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}

func createBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book *Book
	_ = json.NewDecoder(r.Body).Decode(&book)
	ret := InsertBook(&dataStore, book)
	if ret != nil {
		json.NewEncoder(w).Encode(ret)
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var book *Book
	_ = json.NewDecoder(r.Body).Decode(&book)
	ret := UpdateBook(&dataStore, params["id"], book)
	if ret != nil {
		json.NewEncoder(w).Encode(ret)
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	result := DeleteBookById(&dataStore, params["id"])
	if !result {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func main() {
	// init router
	r := mux.NewRouter()

	jsonFilePath := "config.json"
	cfg := configHelper.GetConfigs(jsonFilePath)

	addr1 := cfg.AppSettings["addr1"].(string)
	addr2 := cfg.AppSettings["addr2"].(string)
	addr3 := cfg.AppSettings["addr3"].(string)
	database := cfg.AppSettings["database"].(string)
	source := cfg.AppSettings["source"].(string)
	userName := cfg.AppSettings["userName"].(string)
	passWord := cfg.AppSettings["passWord"].(string)

	dialInfo := &mgo.DialInfo{
		Addrs:    []string{addr1, addr2, addr3},
		Database: database,
		Source:   source,
		Username: userName,
		Password: passWord,
		DialServer: func(addr *mgo.ServerAddr) (net.Conn, error) {
			return tls.Dial("tcp", addr.String(), &tls.Config{})
		},
		Timeout: time.Second * 10,
	}
	s, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		panic(err)
	}

	dataStore = DataStore{Session: s}

	// just for debugging purposes
	ClearDatabase(&dataStore)
	InitializeDatabase(&dataStore)

	// route handlers
	r.HandleFunc("/api/books", getBooks).Methods("GET")
	r.HandleFunc("/api/books/{id}", getBook).Methods("GET")
	r.HandleFunc("/api/books", createBook).Methods("POST")
	r.HandleFunc("/api/books/{id}", updateBook).Methods("PUT")
	r.HandleFunc("/api/books/{id}", deleteBook).Methods("DELETE")

	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	log.Fatal(http.ListenAndServe(":"+port, r))
}
