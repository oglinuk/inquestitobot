package main

import (
	"io"
	"log"
	"os"
)

func main() {
	logf, err := os.OpenFile("inquestitobot.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err.Error())
	}

	mw := io.MultiWriter(logf, os.Stdout)
	log.SetOutput(mw)

	p := NewProcessor()
	for {
		p.process()
	}
}
