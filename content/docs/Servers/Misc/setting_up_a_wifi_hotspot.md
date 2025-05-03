---
weight: 999
url: "/Monter_un_Hotspot_Wifi/"
title: "Setting up a WiFi Hotspot"
description: "Guide for setting up a fully functional WiFi hotspot with OpenWRT, captive portal, and QoS"
categories: ["Linux", "Network"]
date: "2008-05-23T07:05:00+02:00"
lastmod: "2008-05-23T07:05:00+02:00"
tags: ["Wifi", "OpenWRT", "Network", "Hotspot", "QoS", "pf", "squid", "Firewall", "DHCP"]
toc: true
---

## Introduction

One of the things, if not THE thing I was dying to do in my new apartment, was to set up a real hotspot, as I explained in the two news posts below. Well, it's now up and running. I don't know yet if I'll make a full article about it or if I'll just give tips as I go along, but in the meantime, here's what it looks like:

![Empire Network](/images/empirenet.avif)

## Installation and configuration

Two OpenWRT devices configured as simple bridges allow guests to connect to the VLAN dedicated to the public wireless network:

```bash
# /etc/config/wireless
config wifi-device  wifi0
        option type     atheros
        option disabled 0
        option mode             11b
        option distance         10000
        option diversity        0
        option txantenna        1
        option rxantenna        1
        option channel          6

config wifi-iface
        option device   wifi0
        option ifname   ath0
        option network  lan
        option mode     ap
        option ssid     Empire-Network
        option encryption none
        option txpower 18
```

```bash
# /etc/config/network
config interface loopback
        option ifname   lo
        option proto    static
        option ipaddr   127.0.0.1
        option netmask  255.0.0.0

config interface lan
        option ifname   eth0 ath0
        option type     bridge
        option proto    static
        option ipaddr   192.168.200.253
        option netmask  255.255.255.0
```

A DHCP server provides an IP in the appropriate subnet. A simple pf rule redirects all HTTP requests to a captive portal that explains to guests what information to enter in their browser to be able to use HTTP, HTTPS, and FTP (note that so far, only one out of about 30 has managed to complete this highly technical operation...). Some QoS rules ensure that guests don't consume all my bandwidth:

```bash
# /etc/pf.conf
int="fxp0"

table <empire_guests> { 192.168.200.0/24, ! 192.168.200.254, ! 192.168.200.253, ! 192.168.200.252 }

altq on $int cbq bandwidth 28Mb queue { empirenet_in, empirenet_out }
queue empirenet_in bandwidth 2Mb priority 1 cbq(default)
queue empirenet_out bandwidth 128Kb priority 7

rdr on $int inet proto tcp from any to <empire_guests> port www -> 127.0.0.1 port 80

pass in on $int from any to <empire_guests> queue empirenet_in
pass out on $int from <empire_guests> to any queue empirenet_out
```

The user then goes through Squid, and their activity is filtered by squidGuard, in which I've blocked the categories !aggressive !violence !hacking !ads !porn !warez !suspect.
I apply port-based access lists on the switch that only allow HTTP, SSH, and DHCP protocols.

```bash
[...]
ip access-list extended wifiout
 permit ip any host 192.168.200.254
 permit tcp any any eq http
 permit udp any any eq bootps
 permit tcp any any eq ssh
[...]
interface e 18
 ip access-group wifiout in
[...]
```

Everything is graphed by Cacti, notably thanks to the dhcpd-snmp extension of Net-SNMP and the associated Cacti template.

If you're in the 17th arrondissement of Paris, look for the SSID "Empire-Network" :)

## Resources
- http://imil.net/wp/?p=196
