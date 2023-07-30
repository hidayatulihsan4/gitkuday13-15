package connection

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

var Conn *pgx.Conn

func DatabaseConnect() {
	var err error

	databaseUrl := "postgres://postgres:123rahasia@localhost:5432/db_project"

	Conn, err = pgx.Connect(context.Background(), databaseUrl)

	if err != nil {
		fmt.Fprintf(os.Stderr, "gagal connect: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("pun connect")
}
