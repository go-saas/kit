# This project is under development

<div align="center">
<p><a href="https://github.com/go-saas/kit" target="_blank"><img src="https://github.com/go-saas/kit/blob/main/docs/static/img/logo.png?raw=true" width="100"></a></p>
<span  style="font-weight: 800;font-size: 20px;">GO-SAAS-KIT </span> <span  style="font-weight: 800;font-size: 16px;">Starter kit for golang sass project</span>
<br> 
<a href="https://go-saas.github.io/kit/" target="_blank"><span  style="font-weight: 600;font-size: 18px;">docs</span></a>
<br>
</div>
Overview

![Overview](https://github.com/go-saas/kit/blob/main/docs/en-US/overview.png?raw=true)

# Architecture
![Architecture](https://github.com/go-saas/kit/blob/main/docs/static/img/go-saas-kit.drawio.png?raw=true)

# Demo 

address http://saas.nihaosaoya.com (Shanghai)

- **Host** Username:admin  Password:123456


# Feature

* [x] Saas
* [x] Modularity
* [x] ACL(Access Control List), RBAC(Role-based Access Control)
* [x] Localization
* [x] Microservice/Monolithic compatible
* [x] Distributed Eventbus: [kafka](https://kafka.apache.org/), [pulsar](https://pulsar.apache.org/)
* [x] Cache (Redis)
* [x] Background Job: [asynq](https://github.com/hibiken/asynq)
* [x] Virtual File System: [vfs](https://github.com/goxiaoy/vfs)
* [x] Distributed Transaction: [dtm](https://dtm.pub/)
* [x] OpenId Connect: [ory](https://www.ory.sh/)
* [x] Logging/Tracing


# Modules
* [x] User Management
* [x] Tenant Management, Tenant Plans and Subscription
* [x] Payments and Orders
* [x] Product Management


# Quick Start

### For Microservice

```
docker compose -f docker-compose.yml -f docker-compose.ms.yml -f docker-compose.kafka.yml -f docker-compose.tracing.yml up -d
```

Or with build
```
docker compose -f docker-compose.yml -f docker-compose.ms.yml -f docker-compose.kafka.yml  -f docker-compose.tracing.yml up -d --build
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

## Create New Service

```shell
kratos new <name> -r https://github.com/go-saas/kit-layout.git
```


Frontend Repo: https://github.com/go-saas/kit-frontend  
Layout Repo( For creating new service): https://github.com/go-saas/kit-layout

