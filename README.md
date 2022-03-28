# This project is under development

# go-saas-kit

Microservice architecture kit for golang sass project

# Architecture

![Architecture](https://github.com/goxiaoy/go-saas-kit/blob/main/docs/go-saas-kit.drawio.png?raw=true)

# Quick Start

```
docker-compose -f docker-compose.yml -f docker-compose.dev.yml -f docker-compose.tracing.yml up -d --build
```
With hydra
```
docker-compose -f docker-compose.yml -f docker-compose.dev.yml -f docker-compose.hydra.yml up -d --build
```

Open `http://localhost:8600` to see the web ui

Username: admin
Password: youShouldChangeThis@123

Open `http://localhost:8600/swagger` to see the api documentation