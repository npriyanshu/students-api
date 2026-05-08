package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/npriyanshu/students-api/internal/config"
)

// import "fmt"

func main() {

	// load config
	cfg := config.MustLoad()
	// database setup
	// setup router

	router := http.NewServeMux()

	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {

		w.Write([]byte("Welcome to the students api"))

		w.Write([]byte("this is updated response"))
	})
	// setup server
	server := http.Server {
		Addr:    cfg.HTTPServer.Addr,
		Handler : router,
	}

	fmt.Println("Server started successfully on", cfg.HTTPServer.Addr)

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func(){
		err:= server.ListenAndServe()
	if err != nil {
		// panic(err)
		log.Fatalf("Server failed to start: %s", err)
	}
	}()

	<-done

	slog.Info("Shutting down the server")

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	err := server.Shutdown(ctx)

	if err != nil {
		slog.Error("Error occurred while shutting down the server", "error", err)
	} 
	slog.Info("Server stopped gracefully")
	
}
