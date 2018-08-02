package main

import (
	"log"
	"math/rand"
	"strconv"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func GetAllBooks(ds *DataStore) []*Book {
	session := ds.Session.Copy()
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	c := session.DB("test").C("books")

	result := []*Book{}
	err := c.Find(nil).All(&result)
	if err != nil {
		log.Fatal(err)
	}
	return result
}

func GetBookById(ds *DataStore, bookId string) *Book {
	session := ds.Session.Copy()
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	c := session.DB("test").C("books")

	var result *Book
	err := c.Find(bson.M{"id": bookId}).One(&result)
	if err != nil {
		log.Fatal(err)
	}
	return result
}

func InsertBook(ds *DataStore, book *Book) *Book {
	session := ds.Session.Copy()
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	c := session.DB("test").C("books")
	book.ID = strconv.Itoa(rand.Intn(10000000)) // mock ID

	err := c.Insert(&book)
	if err != nil {
		log.Fatal(err)
	}
	var result *Book
	err = c.Find(bson.M{"id": book.ID}).One(&result)
	if err != nil {
		log.Fatal(err)
	}
	return result
}

func UpdateBook(ds *DataStore, bookId string, book *Book) *Book {
	session := ds.Session.Copy()
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	c := session.DB("test").C("books")

	err := c.Update(bson.M{"id": bookId}, book)
	if err != nil {
		log.Fatal(err)
	}
	var result *Book
	err = c.Find(bson.M{"id": bookId}).One(&result)
	if err != nil {
		log.Fatal(err)
	}
	return result
}

func DeleteBookById(ds *DataStore, bookId string) bool {
	session := ds.Session.Copy()
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	c := session.DB("test").C("books")

	err := c.Remove(bson.M{"id": bookId})
	if err != nil {
		log.Fatal(err)
		return false
	}
	return true
}
