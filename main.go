package main

import (
	"log"

	"github.com/alonsoF100/13_14-sprint/pkg/server"
)

func main() {
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
