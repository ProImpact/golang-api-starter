package sql

import (
	"apistarter/internal/config"
	"embed"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

//go:embed schema/*.sql
var embedMigrations embed.FS

func MigrateTo(driver *pgx.Conn, cfg *config.Configuration) {

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatal(err)
	}

	sqlDB := stdlib.OpenDB(*driver.Config().Copy())
	defer sqlDB.Close()

	if err := goose.UpTo(sqlDB, "schema", int64(cfg.DatabaseConfig.Version)); err != nil {
		log.Fatal(err)
	}

}
