package main

import (
	"context"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"

	"basicapi/ent"
	"basicapi/server"
)

func main() {
	db, err := ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	if err != nil {
		log.Fatal().Err(err).Msg("failed opening connection to sqlite")
		return
	}
	if err := db.Schema.Create(context.Background()); err != nil {
		log.Fatal().Err(err).Msg("failed creating schema resources")
		return
	}

	srv := server.NewServer(db)
	if err := srv.Run(":80"); err != nil {
		panic(err)
	}
}
