---
weight: 999
url: "/traefik_install_and_config/"
title: "Traefik: a modern reverse proxy, installation and configuration"
description: "Traefik is a modern reverse proxy that can be used to manage and route traffic to your applications"
categories: ["Reverse-Proxy", "Web", "Proxy"]
tags: ["Reverse-Proxy", "Web", "Proxy"]
toc: true
---

## Introduction

Traefik is an open-source Application Proxy that makes publishing your services a fun and easy experience. It receives requests on behalf of your system, identifies which components are responsible for handling them, and routes them securely.

What sets Traefik apart, besides its many features, is that it automatically discovers the right configuration for your services. The magic happens when Traefik inspects your infrastructure, where it finds relevant information and discovers which service serves which request.

Traefik is natively compliant with every major cluster technology, such as Kubernetes, Docker Swarm, AWS, and the list goes on; and can handle many at the same time. (It even works for legacy software running on bare metal.)

With Traefik, there is no need to maintain and synchronize a separate configuration file: everything happens automatically, in real time (no restarts, no connection interruptions). With Traefik, you spend time developing and deploying new features to your system, not on configuring and maintaining its working state.

## Prerequisites

In this tutorial, we will use Docker to run Traefik. You need to have Docker and Docker Compose installed on your system.

We'll also use Cloudflare as a DNS provider for automatic SSL certificate generation and renewal. You need to have a Cloudflare account and an API key.

### Prepare your Traefik configuration

Create a directory for your Traefik configuration files. You can name it `traefik` or any other name you prefer.

```bash
mkdir -p ~/traefik/{config,logs}
```

## Docker-compose configuration

With docker-compose, here is a simple example of how to install Traefik (`docker-compose.yml`):

```yaml
version: "3"

services:
  traefik:
    image: traefik:2.11.20
    container_name: "traefik"
    restart: always
    network_mode: "host"
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
      - ./traefik/config:/etc/traefik
      - ./traefik/logs:/logs
    labels:
      traefik.enable: true
      traefik.domain: "mycompany.com"
      traefik.tags: web,lb,traefik
      traefik.frontend.rule: "Host:traefik.docker.local"
      traefik.docker.network: "frontend"
    environment:
      TZ: "Europe/Paris"
      CLOUDFLARE_EMAIL: xxx
      CLOUDFLARE_API_KEY: xxx

networks:
  frontend:
    external: true
```

Here is a breakdown of the important configuration part:

- `volumes`: Mount the Docker socket and Traefik configuration files.
  - `/var/run/docker.sock:/var/run/docker.sock`: This allows Traefik to communicate with the Docker daemon and discover services.
  - `./traefik/config:/etc/traefik`: This is where you can put your Traefik configuration files.
  - `./traefik/logs:/logs`: This is where Traefik will store its logs.
- `labels`: These are Docker labels that Traefik uses to configure routing.
  - `traefik.enable: true`: This enables Traefik for this service.
  - `traefik.domain`: This is the domain name for your service.
  - `traefik.tags`: These are tags that can be used to group services.
  - `traefik.frontend.rule`: This is the rule that Traefik will use to route requests to this service.
  - `traefik.docker.network`: This is the Docker network that Traefik will use to communicate with this service.
- `environment`: These are environment variables that Traefik will use.
  - `TZ`: This sets the timezone for Traefik.
  - `CLOUDFLARE_EMAIL`: This is your Cloudflare email address.
  - `CLOUDFLARE_API_KEY`: This is your Cloudflare API key.

Then you can start Traefik with the following command:

## Traefik Configuration

### Basics

Here is a simple example of a Traefik configuration file (`traefik.yml`) to store in the `config` directory:

```yaml
# enable the dashboard and API
api:
  dashboard: true
  insecure: true

accessLog:
  filePath: "/logs/access.log"
  bufferingSize: 100

providers:
  # look for docker containers with labels for dynamic configuration
  docker:
    endpoint: "unix:///var/run/docker.sock"
    exposedByDefault: false
    watch: true
  # load additional static configuration from a file
  file:
    filename: /etc/traefik/rules.yml
    watch: true

entryPoints:
  # enable the dashboard on port 80
  local:
    address: ":8080"
  # enable the web entry point on port 80
  web:
    address: ":80"
    # enable HTTP to HTTPS redirection
    forwardedHeaders:
      insecure: true
    http:
      redirections:
        entryPoint:
          to: websecure
          scheme: https
          permanent: true
  # enable the websecure entry point on port 443
  websecure:
    address: ":443"
    forwardedHeaders:
      insecure: true

certificatesResolvers:
  # use the Cloudflare DNS challenge to obtain a certificate from Let's Encrypt
  letsencrypt:
    acme:
      caServer: "https://acme-v02.api.letsencrypt.org/directory"
      #caServer: https://acme-staging-v02.api.letsencrypt.org/directory
      storage: "~/traefik/acme.json"
      dnsChallenge:
        provider: cloudflare
        delayBeforeCheck: 10
        resolvers:
          - "1.1.1.1:53" # Cloudflare DNS
          - "8.8.8.8:53" # Google DNS

metrics:
  prometheus:
    entryPoint: local
```

### Dynamic Docker Configuration

You can use Docker labels to configure Traefik dynamically. Here is an example of how to configure a service with Docker labels in your `docker-compose.yml` file:

```yaml
services:
  wordpress:
    image: wordpress:6.5.2-php8.1-apache
    container_name: wordpress
    labels:
        traefik.enable: true
        traefik.http.routers.blog.rule: Host(`blog.mycompany.com`)
        traefik.http.routers.blog.entrypoints: websecure
        traefik.http.routers.blog.tls: true
        traefik.http.routers.blog.tls.certResolver: letsencrypt
    deploy:
        resources:
        limits:
        memory: 512M
        reservations:
        memory: 512M
    networks:
        - frontend
    ports:
        - "8010:80"
    volumes:
        - ./blog/php/uploads.ini:/usr/local/etc/php/conf.d/uploads.ini
        - ./blog/wordpress:/var/www/html
    environment:
        WORDPRESS_DB_PASSWORD: xxx
        WORDPRESS_DB_USER: xxx
        WORDPRESS_DB_NAME: xxx
    links:
        - mariadb
    depends_on:
        - mariadb
```

Once the service is started, Traefik will automatically discover it and route requests to `blog.mycompany.com` to the WordPress container with a TLS certificate issued by Let's Encrypt.

### Proxy to a service

You can also use Traefik to proxy requests to a service running on a different host. For example, if you have a service running on a host with IP `192.168.0.2` on port `80`, you can configure Traefik this way in your `rules.yml` file:

```yaml
http:
  services:
    home-assistant:
      loadBalancer:
        servers:
          - url: "http://192.168.0.2:80"

  routers:
    wordpress:
      rule: "Host(`wordpress.mycompany.com`)"
      entrypoints: wordpress
      tls:
        certResolver: letsencrypt
      service: "wordpress"
```

### Redirect

You can add a redirect rule to your Traefik configuration file (`rules.yml`) to redirect domain requests to a new domain. This is useful if you want to redirect traffic from an old domain to a new one.

```yaml
http:
  services:
    # fake service to use for the redirect (mandatory)
    fake:
      loadBalancer:
        servers:
          - url: "http://1.2.3.4"

  routers:
    myredirect:
      # redirect based on the host
      rule: "Host(`old.mycompany.com`) || Host(`x.mycompany.com`)"
      entrypoints: websecure
      tls:
        certResolver: letsencrypt
      # use a middleware to redirect the request
      middlewares:
        - "redirect-middleware"
      # use the fake service to avoid 404 errors
      service: fake

  middlewares:
    redirect-middleware:
      redirectRegex:
        # match the old domain and redirect to the new one
        regex: "^https?://(old\\.|x\\.)?mycompany.com"
        replacement: "https://www.mycompany.com"
        permanent: true
```

### Middlewares

You can use Traefik middleware to modify requests and responses. There are two middlewares I like to use:

- **[sablier](https://github.com/sablierapp/sablier)**: An free and open-source software to start workloads on demand and stop them after a period of inactivity.
- **[crowdsec](https://www.crowdsec.net/)**: CrowdSec provides open source solution for detecting and blocking malicious IPs, safeguarding both infrastructure and application security.

#### Sablier

First of all, you need to install Sablier. You can do this by adding the following to your `docker-compose.yml` file, next to your Traefik service:

```yaml
version: "3"

services:
  sablier:
    image: sablierapp/sablier:1.8.5
    container_name: "sablier"
    restart: always
    ports:
      - "10000:10000"
    deploy:
      resources:
        limits:
          memory: 48M
        reservations:
          memory: 48M
    networks:
      - frontend
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
      - ./sablier/config/sablier.yaml:/etc/sablier/sablier.yaml
```

Create the configuration file for Sablier in the `sablier/config` directory (to match the config file name in the docker-compose config). You can name it `sablier.yaml`. Here is an example configuration:

```yaml
provider:
  name: docker
server:
  port: 10000
  base-path: /
storage:
  file:
sessions:
  default-duration: 5m
  expiration-interval: 20s
logging:
  level: info
strategy:
  dynamic:
    show-details-by-default: trace
    default-theme: shuffle
    default-refresh-frequency: 5s
  blocking:
    default-timeout: 1m
```

You now have to run it:

```bash
docker-compose up -d sablier
```

Then, you need to add the Sablier plugin to your Traefik configuration. You can do this by adding the following to your `traefik.yml` file:

```yaml
experimental:
  plugins:
    sablier:
      moduleName: "github.com/sablierapp/sablier"
      version: "v1.8.5"
```

{{< alert context="warning" text="Ensure you're using the same version in the docker-compose image **AND** `traefik.yml` file" />}}

From the Docker labels, you can add the following to your service configuration in the `docker-compose.yml` file. Here is an example for pgAdmin:

```yaml {linenos=table,hl_lines=[13,14],anchorlinenos=true}
pgadmin:
  container_name: pgadmin
  image: dpage/pgadmin4:latest
  restart: unless-stopped
  mem_limit: 256m
  healthcheck:
    test: wget --no-verbose --tries=1 --spider http://localhost:5050
  environment:
    PGADMIN_DEFAULT_EMAIL: xxx
    PGADMIN_DEFAULT_PASSWORD: xxx
    PGADMIN_LISTEN_PORT: 5050
  labels:
    sablier.enable: true
    sablier.group: pgadmin
    traefik.enable: true
    traefik.http.routers.pgadmin.rule: Host(`pgadmin.mycompany.com`)
    traefik.http.services.pgadmin.loadbalancer.server.port: 5050
    traefik.http.routers.pgadmin.entrypoints: local
  ports:
    - 5050:5050
  links:
    - postgresql
  depends_on:
    - postgresql
  volumes:
    - ./pgadmin/data:/var/lib/pgadmin:rw
```

Add in the `rules.yml` file, add the following middleware configuration (this configuration can be mutualized for all your services):

```yaml {linenos=table,hl_lines=[7],anchorlinenos=true}
middlewares:
  pgadmin:
    plugin:
      sablier:
        dynamic:
          displayName: "PgAdmin"
        group: pgAdmin
        sablierUrl: http://127.0.0.1:10000
        sessionDuration: 15m
```

The options are:

- `dynamic`: The dynamic configuration for the Sablier plugin.
- `displayName`: The name of the service to display in Sablier during loading.
- `group`: The group name for the service.
- `sablierUrl`: The URL of the Sablier service (same as the one in the docker-compose file).
- `sessionDuration`: The duration of the session in Sablier before it stops the service.

## References

- [Traefik Documentation](https://doc.traefik.io/traefik/)
