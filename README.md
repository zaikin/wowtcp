# WoWtcp

WoWtcp is a TCP server and client application that demonstrates a simple Proof of Work (PoW) challenge-response mechanism. The server sends a PoW challenge to connecting clients, and once the client solves the challenge by sending the correct nonce, the server responds with an inspirational quote. This approach helps protect the server from DDoS attacks.

## Features

- **TCP Server**: Sends PoW challenges to clients.
- **TCP Client**: Solves PoW challenges and receives quotes. The client is implemented in a very basic way to demonstrate that the algorithm works and is not intended as a full-featured client.
- **Proof of Work (PoW)**: Uses the Hashcash algorithm to secure connections.
- **Challenger Component**: The challenge generation and verification logic is implemented as a separate component. This design facilitates testing and allows for the easy replacement of the challenge algorithm in the future. Although it may not be the optimal solution, it demonstrates how such a component might be integrated into a larger system.
- **Configurable Logging & Challenge Difficulty**: Easily adjusted via environment variables.
- **Dockerized Deployment**: Run the server and client using Docker containers or Docker Compose.

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

### Using Docker Compose

1. Start the services:
```sh
make run
```

2. Wait for the services to initialize.
3. View the logs: 
```sh
make logs
```

### Running without Docker Compose

1. Set up your environment variables (example):
```bash 
export LOGGER_LEVEL=info
export LOGGER_ENABLE_CALLER=true
export LOGGER_ENABLE_CONSOLE=true
export SERVER_PORT=8080
export CHALLENGE_TYPE=hashcash
export CHALLENGE_DIFFICULTY=2
```

2. Start the server:
```sh 
go run cmd/server/main.go
```

3. Start the client (with the required environment variables):
```sh
go run cmd/client/main.go -host localhost -port 8080
```

## Configuration

You can configure the application’s behavior using environment variables:

- `LOGGER_LEVEL`: Logging level (e.g., info, debug)
- `LOGGER_ENABLE_CALLER`: Enable caller information in logs (true or false)
- `LOGGER_ENABLE_CONSOLE`: Enable console logging (true or false)
- `SERVER_PORT`: Port for the TCP server
- `CHALLENGE_DIFFICULTY`: PoW challenge difficulty (number of leading zeros required).

## Proof of Work Algorithm

For the Proof of Work (PoW) algorithm, Hashcash was chosen. Hashcash is a well-known and widely used PoW algorithm that is simple to implement and understand. It involves finding a nonce such that the hash of the challenge and the nonce has a certain number of leading zeros. This makes it computationally expensive to find the correct nonce, providing protection against DDOS attacks.

### Why Hashcash?
- **Simple Implementation**: 
    Unlike more complex algorithms (such as scrypt, Ethash, or Cuckoo Cycle), Hashcash is based on standard hashing (using SHA256 or SHA1) and requires only finding a nonce that produces a hash with a specified number of leading zeros. This greatly simplifies its implementation.
- **Balanced Computational Complexity**:
    While a basic SHA256-based PoW could also be used for DDoS protection, Hashcash, with its inclusion of additional parameters (such as a timestamp and resource identifier), ensures each request is unique and makes the task more challenging for potential attackers. This helps achieve a balance between computational expense and effective protection.
- **Proven Reliability**:
    Hashcash has proven to be an effective method for preventing spam and mitigating DDoS attacks, with its reliability validated through its use in various systems. Its proven track record sets it apart from some newer algorithms that have not yet seen widespread adoption.
- **Ease of Adaptation and Scalability**:
    While algorithms like scrypt or Ethash are designed for scenarios requiring significant computational and/or memory resources, Hashcash remains a lightweight and flexible solution ideal for demonstration projects and systems with moderate security requirements.

## TCP Messages

The TCP server and client communicate using the following messages:
- `quote!`: The client sends this message to request a quote from the server. The server responds with a PoW challenge.
- `challenge: ...`: The server sends this message as a PoW challenge to the client. The client must solve the challenge and send the correct nonce back to the server.
- `nonce: ...`: The client sends this message with the solved nonce. If the nonce is correct, the server responds with a quote.
- `quote: ...`: The server sends this message with a quote if the client's nonce is correct.
- `quit!`: The client sends this message to close the connection.

## Design Decisions

Due to the simplicity of the application, a separate service layer was omitted to avoid unnecessary complexity. This design keeps the code straightforward and easy to understand while demonstrating secure TCP communication using PoW. Additionally, the TCP client is implemented in a very basic way solely to demonstrate that the algorithm works rather than serving as a fully featured client application.