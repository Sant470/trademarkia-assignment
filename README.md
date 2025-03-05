
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

```bash
  add credential in config.yaml
```

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

Use the below collection to test different apis 

```bash
https://universal-spaceship-200537.postman.co/workspace/Flam-Campaigns-Svc~5bb9cd61-1e44-4ff7-91b6-523485c33b4e/collection/8595172-5a29cddc-67e1-46cb-b7e2-fc394e0fabb7?action=share&creator=8595172
```
## Dependency && Installation

It requires go version 1.22.0, you can download it following the guide mentioned below 

```bash
   https://go.dev/dl/
```
    