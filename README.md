# Golang API Starter

Welcome to the Golang API Starter project! This is a robust template designed for building scalable APIs using Golang. It's packed with powerful tools and configurations to accelerate your development process efficiently.

## Project Structure

- **cmd/api/**: Contains the main application entry point.
- **internal/**: Core application components.
  - **app/**: Core application module.
  - **config/**: Configuration management.
  - **db/**: Database connection and operations.
  - **env/**: Environment settings.
  - **metrics/**: Performance metrics.
  - **security/**: Security features, such as JWT.
  - **server/**: HTTP server setup and routing.
  - **shutdown/**: Graceful shutdown logic.
  - **validation/**: Request validation.
- **docker/**: Docker configuration files.
- **docs/**: API documentation files.
- **pkg/**: Contains reusable packages such as models and utility functions.
- **sql/**: SQL related files and migrations.

## Key Features

- **SQLC with PostgreSQL (pgx)**: Efficient SQL code generation.
- **Gin for HTTP Routing**: High-performance router for handling requests.
- **Graceful Shutdown**: Ensures all resources are cleanly terminated.
- **Viper for Config**: Centralized configuration management with preference for environment variables.
- **Telemetry with Jaeger**: Built-in support for tracing and monitoring.
- **Goose for Migrations**: Manage database migrations seamlessly.
- **Golang FX**: Implements robust dependency injection.
- **Air for Live Reloading**: Hot-reload capabilities for rapid development.

## Scalability

- Modular architecture promotes scalability.
- Efficient database communication with SQLC and pgx.
- Graceful shutdown ensures integrity during scaling actions.

## Deployment Recommendations

- Use Docker for containerized deployments (`docker/docker-compose-dev.yml`).
- Centralize configurations in `application.yml` and environment variables.
- Monitor and trace with Jaeger integrated directly into the environment.

## Setup and Build Instructions

1. **Clone the Repository**: 
   ```bash
   git clone <repository-url>
   cd golang-api-starter
   ```
   
2. **Configuration**:
   - Update `application.yml` with your settings.
   - Use `.env.example` as a template for your environment configurations.

3. **Build and Run**:
   - **Development with Live Reload**: 
     ```bash
     air
     ```
   - **Production**: Use Docker or a suitable CI/CD tool for deployment.

4. **Dependencies**:
   - Required: Golang, PostgreSQL.
   - Use Go modules (`go.mod` file) to install dependencies.

## Recommendations

- Customize middleware and routes within `internal/server/`.
- Leverage Go's concurrency primitives for handling large workloads.
- Regularly update dependencies to improve security and performance.
- Utilize the `Makefile` for automation scripts and tasks.

Embark on creating efficient, scalable applications with this starter as a solid foundation!