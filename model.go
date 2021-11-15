package main

import (
	"bytes"
	"crypto/md5"
	"encoding/gob"
	"encoding/hex"
	"io"
)

// hash generates an md5 hash from given (s)tring
func hash(s string) string {
	h := md5.New()
	io.WriteString(h, s)
	return hex.EncodeToString(h.Sum(nil))
}

// ihash generates an md5 hash from an (i)nterface
func ihash(i interface{}) string {
	var buf bytes.Buffer
	h := md5.New()
	gob.NewEncoder(&buf).Encode(i)
	h.Write(buf.Bytes())
	return hex.EncodeToString(h.Sum(nil))
}

// Document 
type Document struct {
	ID string
	Title string
	URL string
	Description string
	Checksum string
}

// NewDocument constructor
func NewDocument(t, u, d string) *Document {
	doc := &Document{hash(u), t, u, d, ""}
	doc.Checksum = ihash(doc)
	return doc
}

