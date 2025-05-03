---
title: "Docker Compose: create and manage multi-container applications"
date: 2025-05-03T12:00:00+01:00
draft: false
---

## Introduction

Docker Compose is a tool for defining and running multi-container [Docker applications]({{< ref "docs/Servers/Containers/Docker/docker.md" >}}). With Compose, you use a YAML file (typically named `docker-compose.yml`) to configure your application's services, networks, and volumes. Then, with a single command, you create and start all the services defined in your configuration.

Docker Compose is ideal for:

- Development environments
- Automated testing environments
- Single host deployments
- CI/CD workflows

Using Docker Compose involves a three-step process:

1. Define your application's environment in a Dockerfile
2. Define the services that make up your application in a `docker-compose.yml` file
3. Run `docker compose up` to start and run your entire application

## Usage

### Installation

Before using Docker Compose, make sure you have Docker installed and the Docker engine running on your system.

To install Docker Compose:

```bash
sudo apt-get update
sudo apt-get install docker-compose-plugin
```

### Creating a docker-compose.yml file

Here's a basic example of a `docker-compose.yml` file:

```yaml
version: "3"
services:
  web:
    image: nginx:alpine
    ports:
      - "8080:80"
    volumes:
      - ./website:/usr/share/nginx/html
  db:
    image: postgres:13
    environment:
      POSTGRES_PASSWORD: example
      POSTGRES_USER: user
      POSTGRES_DB: mydb
    volumes:
      - postgres_data:/var/lib/postgresql/data

volumes:
  postgres_data:
```

This configuration sets up two services:

- A web server using the nginx:alpine image
- A database using the postgres:13 image

### Basic Commands

#### Start your services

```bash
docker compose up
```

Add `-d` to run in detached mode (background):

```bash
docker compose up -d
```

#### Stop your services

```bash
docker compose down
```

To remove volumes as well:

```bash
docker compose down -v
```

#### View running services

```bash
docker compose ps
```

#### View logs

```bash
docker compose logs
```

Follow logs with:

```bash
docker compose logs -f
```

#### Execute commands in a service container

```bash
docker compose exec web sh
```

### Scaling Services

You can run multiple instances of a service:

```bash
docker compose up -d --scale web=3
```

This starts 3 instances of the web service.

### Environment Variables

You can use environment variables in your docker-compose.yml:

```yaml
services:
  web:
    image: nginx:alpine
    ports:
      - "${NGINX_PORT}:80"
```

Create a `.env` file in the same directory:

```
NGINX_PORT=8080
```

### Working with Networks

Docker Compose automatically creates a network for your application. You can also define custom networks:

```yaml
services:
  web:
    networks:
      - frontend
  db:
    networks:
      - backend
      - frontend

networks:
  frontend:
  backend:
```

### Depends On

Specify dependencies between services:

```yaml
services:
  web:
    depends_on:
      - db
      - redis
  db:
    image: postgres
  redis:
    image: redis
```

### Health Checks

Add health checks to ensure services are properly started:

```yaml
services:
  web:
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost"]
      interval: 1m30s
      timeout: 10s
      retries: 3
      start_period: 40s
```

### Example: Full-Stack Web Application

Here's a more complete example for a web application with frontend, backend, and database:

```yaml
version: "3"

services:
  frontend:
    build: ./frontend
    ports:
      - "3000:3000"
    depends_on:
      - backend
    environment:
      - API_URL=http://backend:8000

  backend:
    build: ./backend
    ports:
      - "8000:8000"
    depends_on:
      - db
    environment:
      - DATABASE_URL=postgres://user:password@db:5432/mydb
      - NODE_ENV=development

  db:
    image: postgres:13
    volumes:
      - postgres_data:/var/lib/postgresql/data
    environment:
      - POSTGRES_USER=user
      - POSTGRES_PASSWORD=password
      - POSTGRES_DB=mydb

volumes:
  postgres_data:
```

This setup demonstrates how Docker Compose can orchestrate a complete application stack with multiple interconnected services.

## Resources

- [Docker]({{< ref "docs/Servers/Containers/Docker/docker.md" >}})
- [Docker Compose Documentation](https://docs.docker.com/compose/)
