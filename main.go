package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"tutorial.sqlc.dev/app/api"
	db "tutorial.sqlc.dev/app/db/sqlc"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5050/simple_bank?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)
func main(){
	conn, err := pgxpool.New(context.Background(), dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}