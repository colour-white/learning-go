package main

import (
	"context"
	"go-status/internal/api"
	"go-status/internal/models"
	"go-status/internal/monitor"
	"go-status/internal/storage"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {

	db, err := storage.InitDatabase("sqlite.db")

	if err != nil {
		panic(err)
	}

	t := &models.Target{Url: "google.com", Interval_sec: 10, Contact_info: "mail.com", Is_active: true, Created_at: time.Now()}
	t, err = models.InsertTarget(db, t)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	manager := monitor.Manager{ActiveWorkers: make(map[int]context.CancelFunc), DB: db, RootCtx: ctx}
	server := api.Server{DB: db, Manager: &manager}
	mux := server.Routes()

	addr := ":8000"
	log.Printf("Listening on %s\n", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
