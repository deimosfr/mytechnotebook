---
weight: 999
url: "/TinyProxy\\:_Mise_en_place_d'un_proxy_simple_et_rapide/"
title: "TinyProxy: Setting up a Simple and Fast Proxy"
description: "Learn how to set up TinyProxy, a lightweight alternative to Squid that simplifies proxy configuration for all your applications."
categories: ["Linux", "Debian", "Network"]
date: "2012-07-27T11:48:00+02:00"
lastmod: "2012-07-27T11:48:00+02:00"
tags: ["Proxy", "Network", "Debian", "Servers"]
toc: true
---

![TinyProxy](/images/tinyproxy_logo.avif)

{{< table "table-hover table-striped" >}}
|||
|-|-|
| **Software version** | 1.8.3 |
| **Operating System** | Debian 7 |
| **Website** | [TinyProxy Website](https://banu.com/tinyproxy/) |
| **Last Update** | 27/07/2012 |
{{< /table >}}

## Introduction

TinyProxy is a proxy server, similar to Squid, but designed to be lightweight. I use it at work on my machine for all my applications. The advantage is having the proxy configuration for wget, apt, etc., pointing to localhost. This way, in the TinyProxy configuration, you only need to configure the proxy server you want to connect to, and all your applications will be redirected without having to reconfigure them every time there's a change :-)

## Installation

On Debian:

```bash
aptitude install tinyproxy
```

## Configuration

I'll skip explaining all the configuration file lines and will show you the ones I've used:

(`/etc/tinyproxy.conf`)

```bash {linenos=table,hl_lines=[7,13,51,61]}
[...]
#
# Port: Specify the port which tinyproxy will listen on.  Please note
# that should you choose to run on a port lower than 1024 you will need
# to start tinyproxy using root.
# We specify here the listening port for Tiny proxy
Port 3128
#
# Listen: If you have multiple interfaces this allows you to bind to
# only one. If this is commented out, tinyproxy will bind to all
# interfaces present.
# The IP on which it should run
Listen 127.0.0.1
[...]
#
# Upstream:
#
# Turns on upstream proxy support.
#
# The upstream rules allow you to selectively route upstream connections
# based on the host/domain of the site being accessed.
#
# For example:
#  # connection to test domain goes through testproxy
#  upstream testproxy:8008 ".test.domain.invalid"
#  upstream testproxy:8008 ".our_testbed.example.com"
#  upstream testproxy:8008 "192.168.128.0/255.255.254.0"
#
#  # no upstream proxy for internal websites and unqualified hosts
#  no upstream ".internal.example.com"
#  no upstream "www.example.com"
#  no upstream "10.0.0.0/8"
#  no upstream "192.168.0.0/255.255.254.0"
#  no upstream "."
#
#  # connection to these boxes go through their DMZ firewalls
#  upstream cust1_firewall:8008 "testbed_for_cust1"
#  upstream cust2_firewall:8008 "testbed_for_cust2"
#
#  # default upstream is internet firewall
#  upstream firewall.internal.example.com:80
#
# The LAST matching rule wins the route decision.  As you can see, you
# can use a host, or a domain:
#  name     matches host exactly
#  .name    matches any host in domain "name"
#  .        matches any host with no domain (in 'empty' domain)
#  IP/bits  matches network/mask
#  IP/mask  matches network/mask
# The proxy server it should point to (the one you normally configure your applications for)
Upstream proxy.deimos.fr:3128
[...]
#
# Allow: Customization of authorization controls. If there are any
# access control keywords then the default action is to DENY. Otherwise,
# the default action is ALLOW.
#
# The order of the controls are important. All incoming connections are
# tested against the controls based on order.
# We authorize localhost
Allow 127.0.0.1
```

All you need to do now is configure your applications to use it :-)
