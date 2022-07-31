# This project is under development

# go-saas-kit

Kit for golang sass project

Frontend Repo: https://github.com/go-saas/kit-frontend  
Layout Repo( For creating new service): https://github.com/go-saas/kit-layout

# Architecture
![Architecture](https://github.com/go-saas/kit/blob/main/docs/go-saas-kit.drawio.png?raw=true)

[//]: # (# Demo )

[//]: # (address http://saas.nihaosaoya.com &#40;Shanghai&#41;)

[//]: # (- **Host** Username:admin  Password:123456)

[//]: # (- **Separate Storage Tenant** domain:separateDbDemo Username:admin  Password:123456)

[//]: # (- **Shared Storage Tenant** domain:sharedDbDemo Username:admin  Password:123456)


# Feature

* [x] Saas
* [x] Modularity
* [x] Microservice/Monolithic compatible
* [x] Distributed Eventbus: [kafka](https://kafka.apache.org/), [pulsar](https://pulsar.apache.org/)
* [x] Cache (Redis)
* [x] Background Job: [asynq](https://github.com/hibiken/asynq)

# Quick Start

### For Microservice

```
docker compose -f docker-compose.yml -f docker-compose.ms.yml -f docker-compose.kafka.yml -f docker-compose.tracing.yml up -d
```

Or with build
```
docker compose -f docker-compose.yml -f docker-compose.ms.yml -f docker-compose.kafka.yml  -f docker-compose.tracing.yml up -d --build
```

### For Monolithic

```
docker compose -f docker-compose.yml -f docker-compose.monolithic.yml -f docker-compose.kafka.yml  -f docker-compose.tracing.yml up -d
```
Or with build
```
docker compose -f docker-compose.yml -f docker-compose.monolithic.yml -f docker-compose.kafka.yml  -f docker-compose.tracing.yml up -d --build
```

### Demo

Open `http://localhost:80` to see the web ui

Username: admin  
Password: 123456

# Development

```shell
make init
```
```shell
make all
```
```shell
make build
```

# Modularity

Module design: 

![Minimal](https://github.com/go-saas/kit/blob/main/docs/minimal-module-design.drawio.png?raw=true)


**Api:** Protobuf definition for public/internal service and models

**Event:** Protobuf definition for distributed event bus

**Biz:** Domain layer, definition for all entities and repository interface

**Service:** Business logic, depends on biz repository interface

**Data:** Data access layer, implement biz repository interface, init databases( mysql ,redis), init event bus (kafka ), expose migration function

**Conf:** Protobuf configuration definition

**Server:** Set up http and grpc server. register all services, set up middlewares. set up distributed  event handler, seeding behavior

**Host:** Process entry point, read configuration, set up tracing, logging



For Microservice:

![Minimal](https://github.com/go-saas/kit/blob/main/docs/microservice.drawio.png?raw=true)



For Monolithic:

![Minimal](https://github.com/go-saas/kit/blob/main/docs/monolithic.drawio.png?raw=true)
