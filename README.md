# This project is under development

# go-saas-kit

Microservice architecture kit for golang sass project
Demo address http://106.75.239.46:8600

# Architecture

![Architecture](https://github.com/goxiaoy/go-saas-kit/blob/main/docs/go-saas-kit.drawio.png?raw=true)

# Quick Start

```
docker-compose -f docker-compose.yml -f docker-compose.tracing.yml up -d
```
With hydra
```
docker-compose -f docker-compose.yml -f docker-compose.hydra.yml up -d
```

Open `http://localhost:8600` to see the web ui

Username: admin
Password: 123456

Open `http://localhost:8600/swagger` to see the api documentation