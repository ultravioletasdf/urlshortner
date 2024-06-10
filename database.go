package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/tursodatabase/go-libsql"
)

func turso(config *Config) (*sql.DB, string, error) {
	dir, err := os.MkdirTemp("", "libsql-*")
	if err != nil {
		return nil, "", err
	}
	fmt.Println("Created temporary directory " + dir)

	path := filepath.Join(dir, config.TursoDatabaseName)
	connector, err := libsql.NewEmbeddedReplicaConnector(path, config.TursoUrl, libsql.WithAuthToken(config.TursoToken))
	if err != nil {
		return nil, "", err
	}

	return sql.OpenDB(connector), dir, nil
}
func ConnectToDatabase(config *Config) (*sql.DB, string) {
	turso, dir, err := turso(config)
	if err != nil {
		log.Fatalln("Failed to connect to turso database", err.Error())
	}
	return turso, dir
}
