# Golang API Starter

Welcome to the Golang API Starter project! 🚀 This is a comprehensive template for building powerful and scalable APIs using Golang. It's packed with essential tools and configurations to jumpstart your development process efficiently.

## Key Features

- **SQL Generation with SQLC**: Seamlessly generate SQL code with sqlc, already configured with pgx for PostgreSQL compatibility.
- **Gin Router**: Utilize Gin for efficient HTTP request routing.
- **Graceful Shutdown Manager**: Includes a shutdown manager to gracefully terminate application components.
- **Configuration with Viper**: Easily manage application settings using Viper, with priorities on environment variables for enhanced security.
- **Application Configuration**: Configurations reside in the `application.yml` file, ensuring central management.
- **Telemetry with Jaeger**: Built-in telemetry support using Jaeger for monitoring and tracing.
- **Database Migrations with Goose**: Pre-configured for database versioning through Goose migrations.
- **Dependency Injection with Golang FX**: Leverages Golang FX for seamless dependency injection.
- **Live Reload with Air**: Develop faster with live reloading provided by Air.

## Why Use This Starter?

🚀 **Boost Productivity**: Kickstart your projects with a stack pre-configured for best practices.

🔧 **Flexibility**: Customizable components to fit various project needs.

🔍 **Observability**: Enhanced monitoring and tracing capabilities out-of-the-box.

💡 **Security**: Prioritizes environment variables for secure configurations.

## Getting Started

1. Clone the repository.
2. Update the `application.yml` with your environment-specific settings.
3. Run the server with live-reload using Air:

   ```bash
   air
   ```

Dive into building your API with confidence, knowing you have a robust foundation! 🌟

## Contributing

Feel free to contribute to the project by submitting issues or pull requests.

## License

This project is licensed under the MIT License.

---

Happy Building! 🎉
