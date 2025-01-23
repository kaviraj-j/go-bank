package main

import (
	"database/sql"
	"log"

	"github.com/kaviraj-j/go-bank/api"
	db "github.com/kaviraj-j/go-bank/db/sqlc"
	"github.com/kaviraj-j/go-bank/util"
	_ "github.com/lib/pq"
)

func main() {

	config, err := util.LoadConfig(".")

	if err != nil {
		log.Fatal("cannot read env variables", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db", err)
	}
	store := db.NewStore(conn)
	server := api.NewServer(store)
	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server", err)
	}
}
