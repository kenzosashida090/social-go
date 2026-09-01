package main

import (
	"log"

	"github.com/kenzosashida090/social/db"
	"github.com/kenzosashida090/social/internal/env"
	"github.com/kenzosashida090/social/store"
)

func main() {
	addr := env.GetString("DB_ADDR", "")
	conn, err := db.New(addr, 3, 3, "15m")
	if err != nil {
		log.Fatal("Error", err.Error())

	}
	defer conn.Close()
	store := store.NewStorage(conn)
	db.Seed(store)

}
