package main

import "apistarter/internal/app"

// @title Golang API Starter
// @version 1.0.0
// @description A production-ready starter template for building REST APIs in Go.
// Includes JWT authentication, structured error handling, request validation,
// PostgreSQL integration (via sqlc), Redis/Valkey caching, modular middleware,
// observability (logging, metrics, distributed tracing), and OpenAPI documentation.
//
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
//
// @host localhost:8080
// @BasePath /api/v1
//
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description JWT token. Example: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
//
// @schemes http https
//
// @accept json
// @produce json
func main() {
	app.Api.Run()
}
