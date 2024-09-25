package main

import (
	"fmt"
	"log"
	"os"
	"runtime/trace"
	"time"
)

func foo(channel chan string) {
	// Infinite loop to send pings and receive pongs
	for {
		fmt.Println("Foo is sending: ping")
		channel <- "ping"
		fmt.Println("Foo is receiving:", <-channel)
	}
}

func bar(channel chan string) {
	// Infinite loop to receive pings and send pongs
	for {
		fmt.Println("Bar is receiving:", <-channel)
		fmt.Println("Bar is sending: pong")
		channel <- "pong"
	}
}

func pingPong() {
	gameChan := make(chan string)
	go foo(gameChan) // Nil is similar to null. Sending or receiving from a nil chan blocks forever.
	go bar(gameChan)
	time.Sleep(500 * time.Millisecond)
}

func main() {
	f, err := os.Create("trace.out")
	if err != nil {
		log.Fatalf("failed to create trace output file: %v", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Fatalf("failed to close trace file: %v", err)
		}
	}()

	if err := trace.Start(f); err != nil {
		log.Fatalf("failed to start trace: %v", err)
	}
	defer trace.Stop()

	pingPong()
}
