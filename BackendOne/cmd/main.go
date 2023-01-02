package main

import (
	postgres "BackendOne/pkg"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"net/http"
	"time"
)

type Handler struct {
	dbConn *pgxpool.Pool
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//Запрашиваем количество подключений к бд в пуле
	totalConns := h.dbConn.Stat().TotalConns()

	switch r.Method {
	case http.MethodGet:
		fmt.Fprintln(w, "(GET): total connections to db: ", totalConns)
	case http.MethodPost:
		fmt.Fprintln(w, "(POST) : total connections to db: ", totalConns)
	}
}

func main() {
	dbConfig := postgres.Config{
		User:     "postgres",
		Password: "postgres",
		DbName:   "test_db",
		Host:     "postgresdb", //т.к. запускаем в докер-компоуз, создается единая сеть для всех контейнеров, внутри которой
		Port:     5432,         //доступ к контейнеру осуществляется имени соответствущего сервиса (не контейна)
		PoolMax:  10,
	}

	//создаем пул коннектов к базе, как к ней подключиться описываем в dbConfig
	dbConn, err := postgres.Connection(context.Background(), dbConfig)
	if err != nil {
		log.Fatalf("Err in main.main(): %s", err)
	}

	handler := &Handler{
		dbConn: dbConn,
	}

	http.Handle("/", handler)
	srv := &http.Server{
		Addr:         ":8888",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("Err in main.ListenAndServe(): %s", err)
	}
}
