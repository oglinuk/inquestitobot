package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var (
	dbname = "inquestitobot.db"
)

// DBInstance
type DBInstance struct {
	DB *sql.DB
}

// NewDBInstance constructor
func NewDBInstance() *DBInstance {
	db, err := sql.Open("sqlite3", dbname)
	if err != nil {
		log.Fatal(err.Error())
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS tblDocuments(
		id TEXT PRIMARY KEY,
		title TEXT,
		URL TEXT,
		description TEXT,
		checksum TEXT
	);`)
	if err != nil {
		log.Fatal(err.Error())
	}

	return &DBInstance{
		DB: db,
	}
}


// Insert (d)ocument into the database
func (dbi *DBInstance) Insert(d *Document) (string, error) {
	var err error
	dbi.DB, err = sql.Open("sqlite3", dbname)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer dbi.DB.Close()

	_, err = dbi.DB.Exec(`INSERT INTO tblDocuments(id, title, URL, description, checksum)
	VALUES(?, ?, ?, ?, ?)`, d.ID, d.Title, d.URL, d.Description, d.Checksum)
	if err != nil {
		return "", err
	}

	return d.ID, nil
}
