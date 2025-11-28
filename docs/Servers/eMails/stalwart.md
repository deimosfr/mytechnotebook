---
title: "Stalwart - A Modern complete suite Email Server"
slug: stalwart-a-modern-complete-suite-email-server/
description: "Stalwart is a modern complete suite email server"
categories: ["Email", "Server", "Linux",]
date: "2025-07-01T00:00:00+00:00"
lastmod: "2025-07-01T00:00:00+00:00"
tags: ["Stalwart", "Email", "Server", "Linux"]
draft: true
---

## Introduction

[Stalwart](https://stalwart.com/) is an open-source mail and collaboration server designed for the modern internet.
Stalwart supports a wide range of protocols including JMAP, IMAP4, POP3, SMTP, CalDAV, CardDAV, and WebDAV, making it a comprehensive solution for managing email, calendars, contacts, file storage, and more. Built in Rust, Stalwart is engineered to be secure, fast, robust, and scalable, capable of running everything from small personal mail servers to large, distributed enterprise deployments.

Stalwart offers a variety of services, storage systems, authentication methods and scalability strategies. We're not covering everything here but just being able to run the service with several usages.

## Installation

Stalwart is available as a Docker image, to make it easy to install, we'll use it:

=== "docker-compose.yaml"

    ``` yaml
    services:
        stalwart-mail:
            image: stalwartlabs/mail-server:latest
            volumes:
                - ./stalwart-mail:/opt/stalwart-mail
            restart: unless-stopped
            ports:
                - 443:443
                - 8080:8080
                - 25:25
                - 587:587
                - 465:465
                - 143:143
                - 993:993
                - 4190:4190
                - 110:110
                - 995:995
    ```

## Quick setup

We're going to setup a simple setup with a single user and a single domain.

### Create a domain

### Setup ACME TLS

### Create a user

## Configuration

### Routes: forward emails

This is a simple feature we can find in [Postfix](./Postfix/virtual_domains.md), unfortunately, Stalwart doesn't support this feature out of the box (yet? in version `0.12`, it's not available). Instead we have to use the `routes` feature and Sieve scripts to forward emails to a specific address.

Here is the way to proceed:

1. Create a user to use the route (`user` here who has the `user@mycompany.com` email address)
2. Prepare the Sieve script to forward emails to a specific address:

```sieve
require ["variables", "copy", "vnd.stalwart.expressions", "envelope", "editheader"];

let "i" "0";
while "i < count(envelope.to)" {
  let "redirected" "false";
  
  if eval "eq_ignore_case(envelope.to[i], 'user@mycompany.com')" {
    addheader "Delivered-To" "user@mycompany.com";
    redirect :copy "user@gmail.com";

    deleteheader :index 1 :is "Delivered-To" "user@mycompany.com";
    let "redirected" "true";
  }

  if eval "!redirected" {
    let "destination" "envelope.to[i]";
    redirect :copy "${destination}";
  }
  let "i" "i+1";
}
discard;
```

3. Update the `routes` configuration to use the Sieve script.


=== 
### Relay host

In order to avoid managing the server reputation, we can use a relay host.


===
### Links

- [Stalwart](https://stalwart.com/)
- [Stalwart SMTP Relay Tutorial](https://gist.github.com/chripede/99b7eaa1101ee05cc64a59b46e4d299f#smtp-relay)
