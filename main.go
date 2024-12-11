package main

import (
	"context"
	"log"

	"github.com/dkpcb/finatext_kadai_2/server"
)

func main() {
	if err := server.Run(context.Background()); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
