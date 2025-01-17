# Train Ticketing System

## Overview
The Train Ticketing System is a gRPC-based backend application for purchasing train tickets, viewing receipts, and accessing analytics.

## Features

### Backend Features
- **Ticket Purchase**: Users can purchase tickets with details like origin, destination, name, email, and preferred section.
- **View Receipt**: Retrieve ticket details using the receipt ID.
- **Analytics**: View total tickets sold, total revenue, and section-wise occupancy.
- **User Management**: View users by section, modify seat allocations, and remove users.

## Technologies Used
- **Backend**: Go (Golang), gRPC
- **Database**: In-memory data structure
- **Containerization**: Docker and Docker Compose

## Prerequisites
- Docker and Docker Compose installed
- Go 1.20+ installed (if running locally without Docker)

## Folder Structure
```
.
├── client/              # gRPC Client
├── pkg/                 # Public packages
├── proto/               # Proto definition
├── server/              # gRPC server
├── docker-compose.yml   # Docker Compose configuration
├── Dockerfile           # Dockerfile for backend
├── go.mod               # Go mod file
├── go.sum               # Go sum file
└── README.md            # Project documentation
```

## Setup and Running the Application

### 1. Clone the Repository
```bash
git clone <repository-url>
cd ticket-system
```

### 2. Run the server with Docker Compose
```bash
docker-compose up --build
```
- The backend server will be available on `localhost:50051`.

### 3. Run the client
```bash
go run client/main.go --server="localhost:50051"
```

## Development
### Run Server Locally
1. Navigate to the `server` directory.
2. Run:
   ```bash
   go run main.go --seats=60
   ```
3. For running tests:
   ```bash
   go test ./server -v
   ```
4. (Optional) Number of seats per section (default: `50`)
   ```bash
   --seats=60
   ```

