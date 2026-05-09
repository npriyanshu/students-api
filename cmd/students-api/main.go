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
	"github.com/npriyanshu/students-api/internal/http/handlers/student"
	"github.com/npriyanshu/students-api/internal/storage/sqlite"
)

// import "fmt"

func main() {

	// load config
	cfg := config.MustLoad()

	// database setup
	storage, err := sqlite.New(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize storage: %s", err)
	}

	slog.Info("storage initialized", slog.String("env", cfg.Env), slog.String("version","1.0.0"))

	

	// setup router

	router := http.NewServeMux()

	router.HandleFunc("POST /api/students", student.New(storage))
	router.HandleFunc("GET /api/students/{id}", student.GetById(storage))




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

	serr := server.Shutdown(ctx)

	if serr != nil {
		slog.Error("Error occurred while shutting down the server", "error", serr)
	} 
	slog.Info("Server stopped gracefully")
	
}
