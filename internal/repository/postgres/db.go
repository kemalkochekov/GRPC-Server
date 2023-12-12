package postgres

import (
	"GRPC_Server/internal/configs"
	"GRPC_Server/pkg/connection"
	"GRPC_Server/pkg/pkgErrors"
	"context"
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/pressly/goose/v3"
)

type TDB struct {
	DB               connection.DBops
	connectionConfig string
}

func NewFromEnv() *TDB {
	dbConfig, err := configs.FromEnv()
	if err != nil {
		if errors.Is(err, pkgErrors.ErrDbConfigNotFound) {
			log.Fatal(pkgErrors.ErrDbConfigNotFound)
		}
		log.Fatalf("could not parse DB_PORT or it is empty: %v", err)
	}
	database, err := connection.NewDB(context.Background(), dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect Database %s", err)
	}
	return &TDB{DB: database, connectionConfig: connection.GenerateDsn(dbConfig)}
}

func (d *TDB) SetUpDatabase(migrationPath string) {

	db, err := sql.Open("postgres", d.connectionConfig)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()
	if err := goose.Up(db, migrationPath); err != nil {
		log.Fatalf("Error setting up the database migrations: %v", err)
	}
}
func (d *TDB) TearDownDatabase(migrationPath string) {
	db, err := sql.Open("postgres", d.connectionConfig)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()
	if err := goose.Down(db, migrationPath); err != nil { // Specify the path to your migrations directory
		log.Fatalf("Error tearing down the database migrations: %v", err)
	}
}
