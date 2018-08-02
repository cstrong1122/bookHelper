package main

import (
	"log"

	"gopkg.in/mgo.v2"
)

func ClearDatabase(ds *DataStore) {

	session := ds.Session.Copy()
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	c := session.DB("test").C("books")

	_, err := c.RemoveAll(nil)
	if err != nil {
		panic(err)
	}
}

func InitializeDatabase(ds *DataStore) {

	session := ds.Session.Copy()
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	books := []Book{}

	books = append(books, Book{ID: "1", Isbn: "HAHAH", Title: "My Book 1", Author: &Author{FirstName: "Christopher", LastName: "Strong"}})
	books = append(books, Book{ID: "2", Isbn: "abcde", Title: "My Book 2", Author: &Author{FirstName: "Christopher", LastName: "Strong"}})
	books = append(books, Book{ID: "3", Isbn: "51515", Title: "My Book 3", Author: &Author{FirstName: "Christopher", LastName: "Strong"}})
	books = append(books, Book{ID: "4", Isbn: "90255", Title: "My Book 4", Author: &Author{FirstName: "Christopher", LastName: "Strong"}})

	c := session.DB("test").C("books")
	for _, book := range books {
		err := c.Insert(&Book{ID: book.ID, Isbn: book.Isbn, Title: book.Title, Author: book.Author})
		if err != nil {
			log.Fatal(err)
		}
	}
}
