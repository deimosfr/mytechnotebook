---
weight: 999
url: "/Dumper_les_connections_dune_interface/"
title: "Capturing connections on an interface"
description: "How to capture and monitor network connections on a Cisco interface using access lists and dumps."
categories: ["Linux", "Network"]
date: "2007-05-23T15:45:00+02:00"
lastmod: "2007-05-23T15:45:00+02:00"
tags: ["Cisco", "Networking", "Packet capture", "TCP", "Monitoring"]
toc: true
---

## Introduction

Capturing (or dumping) means to capture network packets. In this article, we'll look at how to capture TCP packets that pass through our Cisco device, particularly through a specific interface. Here's how to proceed.

## Creating the access list

First, create an **access list** called **dumptcp** to allow connections from a host to any of our interfaces. You can specify a particular one if you wish:

```bash
access-list dumptcp permit ip host 192.168.0.104 any
```

Then we do the reverse so that the Cisco can respond:

```bash
access-list dumptcp permit ip any host 192.168.0.104
```

## Creating the Dump

Now that we have the ability to see the traffic, we need to create the dump rule that we'll call dump104, and that we'll use on the **inside** interface:

```bash
capture dump104 access-list dumptcp interface inside
```

Now, we verify that our dump is correctly configured:

```bash
$ show capture
capture dump104 access-list dumptcp interface inside
```

Now that everything is set up, we can view the capture:

```bash
$ show capture dump104
14:25:49.653545 192.168.0.77 > 192.168.0.104: icmp: echo request
14:25:50.650952 192.168.0.77 > 192.168.0.104: icmp: echo request
14:25:51.650967 192.168.0.77 > 192.168.0.104: icmp: echo request
```
