package main

import (
	"database/sql"
	"log"

	"github.com/jnka9755/go-07SIMPLE-BANK/api"
	db "github.com/jnka9755/go-07SIMPLE-BANK/db/sqlc"
	"github.com/jnka9755/go-07SIMPLE-BANK/util"
	_ "github.com/lib/pq"
)

func main() {

	config, err := util.LoadConfig(".")

	if err != nil {
		log.Fatal("cannot load configurations:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal("Cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)

	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
