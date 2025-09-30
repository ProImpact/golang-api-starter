package app

import (
	"apistarter/internal/config"
	"apistarter/internal/db"
	"apistarter/internal/server"
	"apistarter/internal/shutdown"
	"apistarter/sql"
	"net/http"

	"go.uber.org/fx"
)

var Api = fx.New(
	config.Mod,
	server.Mod,
	sql.Mod,
	db.Mod,
	shutdown.Mod,
	fx.Invoke(func(*http.Server) {}),
)
