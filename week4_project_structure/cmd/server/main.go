package main

import (
	"log"

	"github.com/binderclip/godemo/week4_project_structure/internal/server"
)

func main() {
	srv := server.NewServer()
	err := srv.Run()
	if err != nil {
		log.Fatalf("run server error: %v", err)
	}
}
