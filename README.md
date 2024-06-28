
# Voting Service

This project is a simple voting service built with Golang, utilizing a REST API to handle voting requests and store data in a local DynamoDB instance. It also includes a Docker setup for easy deployment.

## Features

- REST API for voting and retrieving results
- Local DynamoDB for data storage
- Docker Compose setup for local development
- Handles large amounts of data/requests efficiently

## Installation

1. Clone the repository:
    ```bash
    git clone https://github.com/balu6914/voting-service.git
    cd voting-service
    ```

2. Build and run the services using Docker Compose:
    ```bash
    docker-compose up --build
    ```

## Usage

### API Endpoints

- **POST /vote**: Submit a vote
    - Request Body:
      ```json
      {
        "user_id": "string",
        "choice": "string"
      }
      ```

- **GET /results**: Retrieve voting results

### Example Requests

- Submit a vote:
    ```bash
    curl -X POST http://localhost:8001/vote -H "Content-Type: application/json" -d '{"user_id": "user1", "choice": "choiceA"}'
    ```

- Get voting results:
    ```bash
    curl http://localhost:8001/results
    ```

## Code Structure

- `main.go`: Entry point for the application, sets up routes and initializes DynamoDB.
- `models.go`: Contains the data models for the application.
- `cmd/wait-for-dynamodb/wait-for-dynamodb.go`: Script to wait for DynamoDB to be ready before starting the main application.
- `Dockerfile`: Docker configuration for the Go service.
- `docker-compose.yml`: Docker Compose configuration to set up the environment.

## License

This project is licensed under the MIT License.
