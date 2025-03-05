
# RBAC and Rate Limiting Service 

This service implements Role-Based Access Control (RBAC) and Rate Limiting to manage API authentication, authorization, and request throttling.


## Run Locally

Clone the project

```bash
  git https://github.com/Sant470/trademarkia-assignment.git
```

Go to the project directory

```bash
  cd trademarkia-assignment
```

Change Credential in config.yaml

### Build and Run with Docker

Build dockerfile 

```bash
  sudo docker build -f build/dev/Dockerfile -t foo .
```

Run the Docker Container 
```bash
  sudo docker run --network=host foo
```

### Run Locally Without Docker

Install dependencies
```bash
  go mod tidy
```

Start the Server
```bash
  go run main.go
```

## Test it locally
Import the collection 'trademarkia.postman_collection.json' in postman.

## Dependency && Installation

It requires go version 1.22.0, you can download it following the guide mentioned below 

```bash
   https://go.dev/dl/
```
    
