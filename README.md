# GRPC Server

GRPC Server is a Go application that integrates a gRPC server with Gateway for Rest api with Kafka messaging and Jaeger distributed tracing.

## Table of Contents
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Usage](#usage)
- [API Endpoints](#api-endpoints)
  - [Get](#get)
  - [Create](#create)
  - [Delete](#delete)
  - [Update](#update)
  - [Api Documentation](#api-documentation)
- [Linting and Code Quality](#linting-and-code-quality)
  - [Linting Installation](#linting-installation)
  - [Linting Usage](#linting-usage)
- [Testing](#testing)
- [Contributing](#contributing)
- [License](#license)


## Installation

1. Clone this repository
  ```bash
    https://github.com/kemalkochekov/GRPC-Server.git
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

## Linting and Code Quality

This project maintains code quality using `golangci-lint`, a fast and customizable Go linter. `golangci-lint` checks for various issues, ensures code consistency, and enforces best practices, helping maintain a clean and standardized codebase.

### Linting Installation

To install `golangci-lint`, you can use `brew`:

```bash
  brew install golangci-lint
```

### Linting Usage

Once installed, you can run golangci-lint on your project by executing the following command at the root directory of your project:

```bash
  golangci-lint run
```
This command performs linting checks on your entire project and provides a detailed report highlighting any issues or violations found.


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

Contributions are welcome! If you find any issues or have suggestions for improvements, please open an issue or submit a pull request.

## License

This project is licensed under the [MIT License](LICENSE).

