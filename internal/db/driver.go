package db

import (
	"apistarter/internal/config"
	"apistarter/internal/shutdown"
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
)

func NewPostgresqlDrive(cfg *config.Configuration, shutdownFuncs *shutdown.ShutdownManager) *pgx.Conn {
	conn, err := pgx.Connect(context.Background(), fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		cfg.DatabaseConfig.UserName,
		cfg.DatabaseConfig.Password,
		cfg.DatabaseConfig.Host,
		cfg.DatabaseConfig.Port,
		cfg.DatabaseConfig.DatabaseName,
	))
	if err != nil {
		log.Fatal(err)
	}
	shutdownFuncs.CleanupFuncs = append(shutdownFuncs.CleanupFuncs, conn.PgConn().Conn().Close)
	return conn
}

func NewQueries(driver *pgx.Conn) *Queries {
	return New(driver)
}
