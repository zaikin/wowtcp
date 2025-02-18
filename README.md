# WoWtcp

WoWtcp is a TCP server and client application that demonstrates a simple Proof of Work (PoW) challenge-response mechanism. The server sends a challenge to the client, and the client must solve the PoW challenge and send the correct nonce back to the server to receive a quote.

## Features

- TCP server that sends PoW challenges to clients
- TCP client that solves PoW challenges and receives quotes
- Configurable logging and challenge difficulty

## Requirements

- Go 1.23 or later
- Docker (optional, for containerized deployment)

## Installation

1. Clone the repository:

```sh
git clone https://github.com/azaikin/WoWtcp.git
cd WoWtcp
```

2. Install dependencies:

```sh 
go mod download
```

## Using make Commands

The Makefile provides several useful commands for managing the project:

- `make install-tools`: Installs necessary tools, including the Go linter.
- `make linter`: Runs the linter to check for code issues.
- `make linter-docker`: Runs the linter inside a Docker container.
- `make generate`: Runs code generation tools.
- `make run`: Builds and runs the Docker containers for the server and client.
- `make stop`: Stops the Docker containers.
- `make restart`: Restarts the Docker containers.
- `make logs`: Displays the logs from the Docker containers.

## Usage

### docker-compose 

1. Run Docker Compose:
```sh
make run
```

2. Wait for the services to start.
3. View the logs: 
```sh
make logs
```

### Running without Docker Compose

1. Prepare envirement like 
```sh 
LOGGER_LEVEL=info
LOGGER_ENABLE_CALLER=true
LOGGER_ENABLE_CONSOLE=true

SERVER_PORT=8080

CHALLENGE_TYPE=hashcash
CHALLENGE_DIFFICULTY=2
```

2. Start server 
```sh 
go run cmd/server/main.go
```

3. Start client with some envirements
```sh
go run cmd/client/main.go -host localhost -port 8080
```

## Configuration

The application can be configured using environment variables. The following environment variables are available:

- `LOGGER_LEVEL`: Logging level (e.g., info, debug)
- `LOGGER_ENABLE_CALLER`: Enable caller information in logs (true or false)
- `LOGGER_ENABLE_CONSOLE`: Enable console logging (true or false)
- `SERVER_PORT`: Port for the TCP server
- `CHALLENGE_DIFFICULTY`: Difficulty level for the PoW challenge

## Proof of Work Algorithm

For the Proof of Work (PoW) algorithm, Hashcash was chosen. Hashcash is a well-known and widely used PoW algorithm that is simple to implement and understand. It involves finding a nonce such that the hash of the challenge and the nonce has a certain number of leading zeros. This makes it computationally expensive to find the correct nonce, providing protection against DDOS attacks.

### Why Hashcash?

- *Simplicity*: Hashcash is straightforward to implement and does not require complex cryptographic operations.
- *Effectiveness*: It provides a good balance between computational effort and security, making it suitable for protecting against DDOS attacks.
- *Widely Used*: Hashcash has been used in various applications, including email spam protection and cryptocurrency mining, proving its reliability and effectiveness.

## TCP Messages

The TCP server and client communicate using the following messages:
- `quote!`: The client sends this message to request a quote from the server. The server responds with a PoW challenge.
- `challenge: ...`: The server sends this message as a PoW challenge to the client. The client must solve the challenge and send the correct nonce back to the server.
- `nonce: ...`: The client sends this message with the solved nonce. If the nonce is correct, the server responds with a quote.
- `quote: ...`: The server sends this message with a quote if the client's nonce is correct.
- `quit!`: The client sends this message to close the connection.

## Design Decisions

The application does not include a separate service layer. Given the simplicity of the application, adding a service layer would introduce unnecessary complexity without providing significant benefits. The current design keeps the code straightforward and easy to understand.