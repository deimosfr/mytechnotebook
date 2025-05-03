---
weight: 999
url: "/Dnsmasq_and_dhclient\\:_use_a_specific_DNS_for_a_specific_domain/"
title: "Dnsmasq and dhclient: use a specific DNS for a specific domain"
description: "Configure dnsmasq and dhclient to use a specific DNS server for a specific domain, allowing you to resolve local domain names from different networks."
categories: ["Debian", "Linux"]
date: "2014-07-27T13:41:00+02:00"
lastmod: "2014-07-27T13:41:00+02:00"
tags: ["Servers", "DNS", "Networking", "DHCP"]
toc: true
---

{{< table "table-hover table-striped" >}}
| | |
|------------|--------------------------|
| Software version | isc-dhcp-client 4.3.0 dnsmasq 2.71 |
| Operating System | Debian 8 |
| Website | [Debian Website](https://www.debian.org) |
| Last Update | 27/07/2014 |
{{< /table >}}

## Introduction

My use case is specific but not isolated. When I'm at work, I'm connected to my VPN at home. I have a specific DNS at home for my domain in deimos.lan and this is very useful to avoid me to remind all the IP of the services I have.

Sometimes, I want to connect to a home service from the VPN, but my bookmarks are with the local DNS at home which is of course not known from the DNS at work. A solution is to add specifics entries in `/etc/hosts` but it quickly starts to be very boring. That's why I've searched a solution to use my DNS at home only when I try to reach deimos.lan domain.

## Installation

First of all, I need a dhcp client (as I have a DHCP server at home and at work) and dnsmasq to run locally on my laptop:

```bash
aptitude install dnsmasq isc-dhcp-client
```

## Configuration

### Dnsmasq

We're going to setup dnsmasq like this:

```bash
server=/deimos.lan/192.168.0.1
interface=lo
listen-address=127.0.0.1
bind-interfaces
```

Here are explanations:

- server: you need to specify for the domain (deimos.lan) which DNS server should be targeted (192.168.0.1)
- interface: the interface to listen on
- listen-address: only listen to 127.0.0.1 IP address
- bind-interfaces: only bind to specified interfaces (here: lo)

And restart dnsmasq to apply this new configuration.

### DHCP client

If we now change resolv.conf values to point to the DNS 127.0.0.1, we will be correctly redirected. But the problem is everytime the dhcp lease will expire and renewed, it will change the resolv.conf file. To avoid it, we're going to add this line in the dhclient configuration file:

```bash
prepend domain-name-servers 127.0.0.1;
```

This will force the first "nameserver" line in resolv.conf to be 127.0.0.1. To finish, restart the dhclient service to have this new version working.

{{< alert context="info" text="On RedHat OS like, you may need to add this line \"PEERDNS=yes\" in your network configuration file" />}}
