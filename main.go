package main

import (
	"log"
)

func main() {
	dbi := NewDBInstance()
	id, err := dbi.Insert(NewDocument("test", "test.com", "A test"))
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println(id)
}
