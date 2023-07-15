package main

import (
	"context"

	_ "github.com/mattn/go-sqlite3"

	"basicapi/ent"
	"basicapi/server"
)

func main() {
	db, err := ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	if err != nil {
		panic(err)
	}
	_ = db.Schema.Create(context.Background())

	srv := server.NewServer(db)
	if err := srv.Run(":80"); err != nil {
		panic(err)
	}
}
