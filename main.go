package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

const ADDR = ":8080"
const READ_TIMEOUT_COUNT = 10
const WRITE_TIMEOUT_COUNT = 10

func handler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	name := query.Get("name")
	if name == "" {
		name = "my Guest"
	}
	log.Printf("Request received for: %s\n", name)
	w.Write([]byte(fmt.Sprintf("Привет, %s\n", name)))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func main() {
	//Create Server and Route Handlers
	r := mux.NewRouter()

	r.HandleFunc("/", handler)
	r.HandleFunc("/health", healthHandler)
	r.HandleFunc("/readiness", readinessHandler)

	srv := &http.Server{
		Handler:      r,
		Addr:         ADDR,
		ReadTimeout:  READ_TIMEOUT_COUNT * time.Second,
		WriteTimeout: WRITE_TIMEOUT_COUNT * time.Second,
	}

	go func() {
		log.Println("Starting Server - listening at port " + ADDR)
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	//Graceful Shutodwn
	waitForShutdown(srv)
}

func waitForShutdown(srv *http.Server) {
	interruptChan := make(chan os.Signal, 1)
	log.Println("interruptChan")
	log.Println(interruptChan)

	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT)
	log.Println("signalNotify")

	//wait until the signal is received
	log.Println("before interruptChain")
	<-interruptChan
	log.Println("after interruptChain")

	// CReate a deadline to wait for
	log.Println("before ctx :=")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*READ_TIMEOUT_COUNT)
	log.Println("after ctx :=, before defer")
	defer cancel()
	log.Println("after defer, before Shutdown")
	srv.Shutdown(ctx)

	log.Println("Shutting down")
	os.Exit(0)
}
