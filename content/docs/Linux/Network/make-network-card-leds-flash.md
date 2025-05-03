---
weight: 999
url: "/Faire_clignoter_les_LEDs_d'une_carte_réseau/"
title: "Make Network Card LEDs Flash"
description: "How to make the LEDs on a network card flash to locate it physically on a server or device."
categories: ["Network", "Linux", "Hardware"]
date: "2012-08-27T06:43:00+02:00"
lastmod: "2012-08-27T06:43:00+02:00"
tags: ["network", "ethtool", "hardware", "troubleshooting"]
toc: true
---

## Introduction

It's sometimes useful to be able to locate a network card when you're communicating with someone on the other side of the planet, who is in front of a machine you're connected to remotely and they need to plug in a cable but don't know which interface to use.

## Configuration

Here's how to make the eth0 interface flash for 10 seconds:

```bash
ethtool -p eth0 10
```
