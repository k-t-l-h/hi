package main

import (
	"RSOI/internal/pkg/middleware"
	"RSOI/internal/pkg/persona/delivery"
	"RSOI/internal/pkg/persona/repository"
	"RSOI/internal/pkg/persona/usecase"
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"os"
	"time"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {

	connection, state := os.LookupEnv("DATABASE_URL")
	if !state {
		log.Print("connection string was not found")
	}

	conn, err := pgxpool.Connect(context.Background(), connection)
	if err != nil {
		log.Fatal("database connection not established")
	}

	port, ok := os.LookupEnv("PORT")

	if ok == false {
		port = "5000"
	}

	name, ok := os.LookupEnv("VER")

	if ok == false {
		name = "ktlh"
	}
	m := middleware.NewMetricsMiddleware()
	m.Register(name)

	pr := repository.NewPRepository(*conn)
	pu := usecase.NewPUsecase(pr)
	pd := delivery.NewPHandler(pu)

	r := mux.NewRouter()
	r.Use(middleware.InternalServerError)
	r.Use(m.LogMetrics)

	r.HandleFunc("/persons/{personID}", pd.Read).Methods("GET")
	r.HandleFunc("/persons", pd.ReadAll).Methods("GET")
	r.HandleFunc("/persons", pd.Create).Methods("POST")
	r.HandleFunc("/persons/{personID}", pd.Update).Methods("PATCH")
	r.HandleFunc("/persons/{personID}", pd.Delete).Methods("DELETE")
	r.PathPrefix("/metrics").Handler(promhttp.Handler())
	http.Handle("/", r)

	srv := &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf(":%s", port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Print("Server running at ", srv.Addr)
	log.Fatal(srv.ListenAndServe())

}
