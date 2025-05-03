---
weight: 999
url: "/Faire_du_reverse_Tunelling_avec_OpenSSH/"
title: "Reverse Tunneling with OpenSSH"
description: "How to set up reverse tunneling with OpenSSH to access machines behind NAT and firewalls."
categories: ["Network", "Security", "SSH"]
date: "2012-02-03T16:08:00+02:00"
lastmod: "2012-02-03T16:08:00+02:00"
tags: ["ssh", "tunneling", "nat", "remote access", "openssh"]
toc: true
---

## Introduction

This is going to be really powerful! What am I proposing? Reverse tunneling? Yes! Imagine being able to traverse NAT. You're already starting to salivate, so let's not delay any further!

## Setup scenario

* Here is the machine I want to connect to: 192.168.20.55
* The machine from which I'm going to launch the connection: 138.47.99.99 (WAN IP)

This will give us:
Destination (192.168.20.55) <- |NAT| <- Source (138.47.99.99)

## Configuration

* We'll connect here and use an unused port on our machine (let's say 19999):

```bash
ssh -N -R 19999:localhost:22 sourceuser@138.47.99.99
```

* I can then pass through the tunnel like this (still from the source machine):

```bash
ssh localhost -p 19999
```

* Now it makes sense that I can access 192.168.20.55 via the 138.47.99.99 machine

Destination (192.168.20.55) <- |NAT| <- Source (138.47.99.99) <- Bob's server

From Bob's server:

```bash
ssh sourceuser@138.47.99.99
```

After this connection, you're on the target machine, now you'll need to bounce through the tunnel to the first machine:

```bash
ssh localhost -p 19999
```
