package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"crypto.com/wallet"
)

func main() {
	go func() {
		s.Handler = wallet.NewRouter()
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("s.ListenAndServe err: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // watch for SIGINT (Ctrl+C)
	<-quit

	s.Close() // shutdown the server
	fmt.Printf("\n\nBye!\n\n")
	os.Exit(0)
}

var s = &http.Server{
	Addr:           ":8080",
	ReadTimeout:    10 * time.Second,
	WriteTimeout:   10 * time.Second,
	MaxHeaderBytes: 1 << 20, // 1 MB
}
