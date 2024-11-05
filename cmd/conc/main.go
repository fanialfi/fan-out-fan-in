package main

import (
	"log"
	"time"

	filegeneration "github.com/fanialfi/fan-out-fan-in/fileGeneration"
)

func main() {
	log.Println("start")
	start := time.Now()

	filegeneration.GenerateFileConcurrency()

	duration := time.Since(start)
	log.Printf("done in %.3f seconds\n", duration.Seconds())
}
