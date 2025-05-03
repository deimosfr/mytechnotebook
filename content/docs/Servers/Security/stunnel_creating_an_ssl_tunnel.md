---
weight: 999
url: "/Stunnel_\\:_Fabrication_d'un_tunnel_SSL/"
title: "Stunnel: Creating an SSL Tunnel"
description: "A guide on how to set up and configure Stunnel to create secure SSL tunnels for services that don't natively support encryption."
categories: ["Linux", "Debian", "Network", "Servers", "Security"]
date: "2007-06-26T08:09:00+02:00"
lastmod: "2007-06-26T08:09:00+02:00"
tags: ["stunnel", "ssl", "tunnel", "encryption", "security", "network", "telnet", "vnc"]
toc: true
---

## Introduction

If you have software that doesn't support SSL, and you want to secure network connections, you can encapsulate it in an SSL tunnel. This tunnel will encrypt data from end to end.

## Installation

### Debian

```bash
apt-get install stunnel4
```

### Red-Hat

```bash
wget http://www.stunnel.org/download/stunnel/src/stunnel-4.20.tar.gz
tar -xzvf stunnel-4.20.tar.gz
cd stunnel-4.20
./configure && make && make install
```

### Windows

Download the client: [https://www.stunnel.org/download/binaries.html](https://www.stunnel.org/download/binaries.html)  
*On Windows, all configuration files are in "C:\Program Files\stunnel", so adapt the examples below according to file paths*

## Configuration

### Serveur

Don't modify the `/etc/stunnel/stunnel.conf` file; it's preferable to create a separate file (for example "/etc/stunnel/services.conf").  
Here's an example of the file contents that will forward telnet and a VNC connection (assuming a VNC server is running on port "5901").

```bash
cert = /etc/stunnel/stunnel.pem         # Certificate to use
CAfile = /etc/stunnel/stunnel.pem       # same
verify = 3                              # Certificate verification level

##Service Definitions##
[Telnet]                                # Service Name
accept = 88.191.31.151:12345            # Server address hosting the service: Secure alternative port
connect = 127.0.0.1:23                  # Local server address: Real service port

[VNC]                                   # Service Name
accept = 88.191.31.151:54321            # Server address hosting the service: Secure alternative port
connect = 127.0.0.1:5901                # Local server address: Real service port
```

### Client

As with the server, it's preferable to create a separate configuration file (still "/etc/stunnel/services.conf").
This file will be similar to the server file except that the service logic is reversed and the "Client" option is defined:

```bash
client = yes                            # Indicates this is the client
cert = /etc/stunnel/stunnel.pem         # Certificate to use
CAfile = /etc/stunnel/stunnel.pem       # same
verify = 3                              # Certificate verification level

##Service Definitions##
[Telnet]                                # Service Name
accept = 127.0.0.1:23                   # Server address hosting the service: Secure alternative port
connect = 88.191.31.151:12345           # Local server address: Real service port

[VNC]                                   # Service Name
accept = 127.0.0.1:5901                 # Server address hosting the service: Secure alternative port
connect = 88.191.31.151:54321           # Local server address: Real service port
```

## Generation du Certificat

Create a file "/etc/stunnel/cert.conf" with the following lines:

```bash
[ req ]
default_bits = 1024 # Set 2056 if you're paranoid
encrypt_key = yes
distinguished_name = req_dn
x509_extensions = cert_type
prompt = no

[ req_dn ]
C=FR
ST=France
L=Paris
O=Deimos
OU=Deimos Network Team
CN=deimos.fr
emailAddress=xxx@mycompany.com

[ cert_type ]
nsCertType = server
```

Then, to generate the certificate, navigate to the "/etc/stunnel" directory and type:

```bash
openssl req -new -days 365 -nodes -config /etc/stunnel/cert.conf -out stunnel.pem -x509 -keyout stunnel.pem
```

## Utilisation

Now that the configuration files are created, stunnel is ready to be launched.

Linux:
```bash
stunnel4 /etc/stunnel/services.conf
```

Windows:
```bash
cd "C:\Programe Files\stunnel"
stunnel.exe services.conf
```

The client can then connect to the remote service this way:
```bash
telnet 127.0.0.1
vnc4client 127.0.0.1:5901
```
