---
weight: 999
url: "/Pound_\\:_Installation_et_Configuration_d'un_Reverse_Proxy/"
title: "Pound: Installation and Configuration of a Reverse Proxy"
description: "A guide to installing and configuring Pound as a reverse proxy and load balancer for web servers"
categories: ["Linux", "Servers", "Network"]
date: "2010-05-15T07:13:00+02:00"
lastmod: "2010-05-15T07:13:00+02:00"
tags: ["Reverse Proxy", "Load Balancer", "Pound", "Web Servers", "Apache"]
toc: true
---

![Pound Logo](/images/pound_logo.avif)

## Introduction

A [reverse proxy](https://en.wikipedia.org/wiki/Reverse_proxy) is a type of proxy server, usually placed in front of web servers. It differs in its usage from traditional proxy servers.

The reverse proxy is implemented on the server side of the Internet. Web users go through it to access applications on internal servers. This technique allows, among other things, to protect a web server from attacks from outside.

This technology is used in application security solutions.

There are several recognized applications for reverse proxies:

- Security: The additional layer provided by reverse proxies can bring additional security. Programmable URL rewriting allows masking and controlling, for example, the architecture of an internal website. But this architecture mainly allows filtering access to web resources from a single point.
- SSL Acceleration: The reverse proxy can be used as an "SSL terminator," for example, through dedicated hardware.
- Load Balancing: The reverse proxy can distribute the load of a single site across multiple web application servers. Depending on its configuration, URL rewriting work will therefore be necessary.
- Cache: The reverse proxy can offload web servers from the load of static pages/objects (HTML pages, images) by managing a local cache. The load on web servers is thus generally reduced.
- Compression: The reverse proxy can optimize the compression of site content.

After some research, it appears that Pound is one of the best solutions for reverse proxying. You can also do it with Apache, Lighttpd, Nginx... But apparently, Pound stands out because:

- It is lightweight and efficient (works very well with over 600 connections/sec)
- Configuration can be easily evolved to do load balancing
- It is capable of managing sessions

## Installation

To install it, it's simple:

```bash
apt-get install pound
```

## Configuration

### Default

Configure `/etc/default/pound` if you want it to start automatically:

```bash
startup=1
```

### Basic Reverse Proxy

Here, I have an Apache running locally on port 8080 and I have Pound listening on port 80:

```bash
## Minimal sample pound.cfg
##
## see pound(8) for details


######################################################################
## global options:

User            "www-data"
Group           "www-data"
#RootJail       "/chroot/pound"

## Logging: (goes to syslog by default)
##      0       no logging
##      1       normal
##      2       extended
##      3       Apache-style (common log format)
LogLevel        1

## check backend every X secs:
Alive           30

## use hardware-accelleration card supported by openssl(1):
#SSLEngine      "<hw>"

# poundctl control socket
Control "/var/run/pound/poundctl.socket"


######################################################################
## listen, redirect and ... to:

## redirect all requests on port 8080 ("ListenHTTP") to the local webserver (see "Service" below):
ListenHTTP
        Address 10.101.0.39
        Port    80

        ## allow PUT and DELETE also (by default only GET, POST and HEAD)?:
        xHTTP           0

        Service
                BackEnd
                        Address 127.0.0.1
                        Port    8080
                End
        End
End
```

### Basic Load Balancing

For a configuration, we'll try a redirection with IP or VirtualHost:

```bash
       ListenHTTP
           Address 192.168.0.200
           Port    80

           Service
               HeadRequire "Host: .*www.deimos.fr.*"

               BackEnd
                   Address 192.168.0.1
                   Port    80
               End
           End

           Service
               HeadRequire "Host: .*www.mavro.fr.*"

               BackEnd
                   Address 192.168.0.2
                   Port    80
               End
           End
       End
```

Here, our server listens on port 80 of IP 192.168.0.200. If the VirtualHost deimos.fr is used, there will be a redirection to IP 192.168.0.1:80. Otherwise, if it's mavro.fr, the redirection will be to address 192.168.0.2:80.  
As you can see, it's quite simple.

*Note: The developer of Pound does not recommend using VirtualHosts and suggests letting the lower layer handle it.*

**Important: Be aware that it is impossible to do VirtualHost with HTTPS. This is due to a limitation of the protocol and not specific to Pound.**

## Resources
- http://www.apsis.ch/pound/index_html
- http://www.cyberciti.biz/faq/linux-http-https-reverse-proxy-load-balancer
