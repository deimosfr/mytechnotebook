---
weight: 999
url: "/Mise_en_place_d'un_quorum_server/"
title: "Setting up a quorum server"
description: "Instructions for setting up and configuring a quorum server for a 2-node cluster without a disk array"
categories:
  - "Linux"
  - "Solaris"
date: "2010-07-26T12:07:00+02:00"
lastmod: "2010-07-26T12:07:00+02:00"
tags:
  - "Servers"
  - "Quorum"
  - "Clustering"
  - "Solaris"
toc: true
---

## Introduction

**THIS IS NEEDED ONLY FOR A 2 NODES CLUSTER WITHOUT A DISKS ARRAY!**

You generally need it when you have 2 servers without a SAN, so you need a third machine with this service installed and launched for your quorum.

## Installation

Normally the Quorum Server software is installed with the Sun Cluster Software.

## Configuration

Just configure the file `/etc/scqsd/scqsd.conf` with this line:

```bash
/usr/cluster/lib/sc/scqsd -d /var/scqsd -i quorum_name -p 9001
```

- quorum_name: is the instance name of the quorum server
- 9001: is the default listening port of quorum server

After that, just launch the quorum server instance:

```bash
scadm restart svc:/system/cluster/quorumserver:default
```
