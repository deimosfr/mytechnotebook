---
weight: 999
url: "/Mise_en_place_d'un_serveur_FreeRadius/"
title: "Setting up a FreeRadius server"
description: "This guide explains how to install and configure a FreeRadius server on OpenBSD, including basic user setup and verification."
categories: ["Linux", "Network", "Servers"]
date: "2009-01-13T11:10:00+02:00"
lastmod: "2009-01-13T11:10:00+02:00"
tags: ["FreeRADIUS", "Servers", "Network", "OpenBSD", "Authentication"]
toc: true
---

## Introduction

[FreeRADIUS](https://fr.wikipedia.org/wiki/FreeRADIUS) is an open-source RADIUS server.

It offers an alternative to other enterprise RADIUS servers, and is one of the most modular and feature-rich RADIUS servers available today. It is considered the most widely used server in the world.

It is suitable for both embedded systems with limited memory and systems with several million users.

I installed this server on OpenBSD to connect a WiFi access point to it. This is quite practical and currently the most secure approach.

## Installation

On OpenBSD:

```bash
pkg_add -iv freeradius
```

Now we need to add it to the boot process:

```bash
if [ -x /usr/local/sbin/radiusd ]; then
        install -d -o _freeradius /var/run/radiusd
        echo -n ' radiusd';     /usr/local/sbin/radiusd
fi
```

If you want a configuration file example, look at `/usr/local/share/examples/freeradius`.

## Configuration

### client.conf

We will edit the client.conf file to add a test user:

```bash
client 127.0.0.1 {
    secret      = testing123
    shortname   = localhost
}
```

### users

Let's add a simple user for now:

```
"deimos"    Cleartext-Password := "password"
```

## Verification

To verify that the user is working properly, you can use the radtest command:

```bash
radtest deimos test 127.0.0.1 1812 testing123
```

- deimos : the user
- password : the password
- 127.0.0.1 : the Radius server
- 1812 : the Radius server port
- testing123: the additional password
