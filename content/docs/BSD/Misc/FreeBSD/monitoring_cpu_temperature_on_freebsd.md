---
weight: 999
url: "/Monitorer_la_temperature_des_processeurs_sous_FreeBSD/"
title: "Monitoring CPU Temperature on FreeBSD"
description: "How to monitor the temperature of all CPU cores on FreeBSD"
categories: ["FreeBSD", "Linux"]
date: "2010-09-12T06:29:00+02:00"
lastmod: "2010-09-12T06:29:00+02:00"
tags: ["FreeBSD", "CPU", "Monitoring", "Temperature", "System Administration"]
toc: true
---

## Introduction

Here's how to monitor the temperature of all your CPU cores on FreeBSD.

## Configuration

Fortunately, there's nothing to install. We'll simply load the temperature module:

```bash
kldload coretemp
```

If you want to enable it each time your machine boots:

```bash
# /boot/loader.conf
coretemp_load="YES"
```

PS: For AMD users, there are also the k8temp and amdtemp modules available.

## Usage

Now, you can easily check your temperatures:

```bash
> sysctl dev.cpu | grep temperature
dev.cpu.0.temperature: 48.0C
dev.cpu.1.temperature: 50.0C
dev.cpu.2.temperature: 47.0C
dev.cpu.3.temperature: 46.0C
```

## References

http://blog.freelooser.fr/2009/05/temperature-du-cpu-sous-freebsd.html
