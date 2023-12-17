# GRPC Server

GRPC Server is a Go application that integrates a gRPC (Google Remote Procedure Calls) server with Gateway for REST API, Kafka messaging, Jaeger distributed tracing, and the `zgo.uber.org/zap` logger.

## Table of Contents
- [Introduction](#introduction)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Usage](#usage)
- [API Endpoints](#api-endpoints)
  - [Get](#get)
  - [Create](#create)
  - [Delete](#delete)
  - [Update](#update)
  - [Api Documentation](#api-documentation)
- [Distributed Tracing with Jaeger](#distributed-tracing-with-jaeger)
- [Linting and Code Quality](#linting-and-code-quality)
  - [Linting Installation](#linting-installation)
  - [Linting Usage](#linting-usage)
- [Testing](#testing)
- [Contributing](#contributing)
- [License](#license)

## Introduction

The GRPC Server is a Go application that provides a gRPC server with a REST API gateway, Kafka messaging, Jaeger distributed tracing, and the zgo.uber.org/zap logger. It allows you to build scalable and high-performance server applications using the gRPC framework while also providing REST API compatibility. Kafka messaging enables asynchronous communication and decoupling of components, and Jaeger distributed tracing helps monitor and identify performance issues in a distributed system. The zgo.uber.org/zap logger ensures efficient logging with minimal impact on performance.

## Prerequisites

Before running this application, ensure that you have the following prerequisites installed:

- Go: [Install Go](https://go.dev/doc/install/)
- Docker: [Install Docker](https://docs.docker.com/get-docker/)
- Docker Compose: [Install Docker Compose](https://docs.docker.com/compose/install/)

## Installation

1. Clone this repository
  ```bash
    git clone https://github.com/kemalkochekov/GRPC-Server.git
  ```
2. Navigate to the project directory:
  ```
    cd GRPC-Server
  ```
3. Build the Docker image:
  ```
    docker-compose build
  ```

## Usage
1. Start the Docker containers:
  ```
    docker-compose up
  ```
2. The application will be accessible at:
  ```
    localhost:8080/v1
  ```

## API Endpoints

### Get

- Method: GET
- Endpoint: /entity
- Query Parameter: id=[entity_id]

Retrieve data from the database based on the provided ID.

- Response:
  - 200 OK: Returns the data in the response body.
  - 400 Bad Request: If the `id` query parameter is missing.
  - 404 Not Found: If the provided ID does not exist in the database.
  - 500 Internal Server Error: If there is an internal server error.

### Create

- Method: POST
- Endpoint: /entity
- Request Body: JSON payload containing the ID and data.

Add new data to the database.

- Response:
  - 200 OK: If the request is successful.
  - 500 Internal Server Error: If there is an internal server error.

### Delete

- Method: DELETE
- Endpoint: /entity
- Query Parameter: id=[entity_id]

Remove data from the database based on the provided ID.

- Response:
  - 200 OK: If the request is successful.
  - 404 Not Found: If the provided ID does not exist in the database.
  - 500 Internal Server Error: If there is an internal server error.

### Update

- Method: PUT
- Endpoint: /entity
- Request Body: JSON payload containing the ID and updated data.

Update existing data in the database based on the provided ID.

- Response:
  - 200 OK: If the request is successful.
  - 404 Not Found: If the provided ID does not exist in the database.
  - 500 Internal Server Error: If there is an internal server error.

### Api Documentation

For detailed API documentation, including examples, request/response structures, and authentication details, please refer to the

<a href="https://documenter.getpostman.com/view/31073105/2s9Ykn8grK" target="_blank">
    <img alt="View API Doc Button" src="https://github.com/kemalkochekov/JWT-Backend-Development-App/assets/85355663/0c231cef-ee76-4cdf-bc41-e900845da493" width="200" height="60"/>
</a>

## Distributed Tracing with Jaeger

This project integrates with Jaeger for distributed tracing. Jaeger allows you to trace the flow of requests across multiple services, providing insights into performance and identifying bottlenecks in the system.

To view the traces, access the Jaeger UI at:
```
http://localhost:16686/
```
For more information on how to use Jaeger for distributed tracing in Go, refer to the Jaeger Go client documentation.

## Linting and Code Quality

This project maintains code quality using `golangci-lint`, a fast and customizable Go linter. `golangci-lint` checks for various issues, ensures code consistency, and enforces best practices, helping maintain a clean and standardized codebase.

### Linting Installation

To install `golangci-lint`, you can use `brew`:

```bash
  brew install golangci-lint
```

### Linting Usage
1. Configuration: 

After installing golangci-lint, create or use a personal configuration file (e.g., .golangci.yml) to define specific linting rules and settings:
```bash
  golangci-lint run --config=.golangci.yml
```
This command initializes linting based on the specified configuration file.

2. Run the linter:

Once configuration is completed, you can execute the following command at the root directory of your project to run golangci-lint:

```bash
  golangci-lint run
```
This command performs linting checks on your entire project and provides a detailed report highlighting any issues or violations found.

3. Customize Linting Rules:

You can customize the linting rules by modifying the .golangci.yml file.

For more information on using golangci-lint, refer to the golangci-lint documentation.

## Testing

- To run unit tests, use the following command:
  ```bash
    go test ./... -cover
  ```
- To run integration tests, use the following command:
  ```bash
    go test -tags integration ./... -cover
  ```

## Contributing

Contributions are welcome! If you would like to contribute to this project, please follow these steps:
1. Fork the repository.
2. Create a new branch for your feature/fix:
3. Make your changes and commit them.
```bash
  git checkout -b feature/your-feature-name
```
4. Push your branch to your forked repository.
5. Create a pull request from your branch to the main repository.

## License

This project is licensed under the [MIT License](LICENSE).

