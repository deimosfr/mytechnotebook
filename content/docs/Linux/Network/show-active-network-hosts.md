---
weight: 999
url: "/afficher-les-machines-allumees-sur-le-reseau-courant/"
title: "Displaying Active Machines on the Current Network"
description: "How to find and display all active machines on your current network using nmap."
categories: ["Networking", "System Administration", "Security"]
date: "2009-12-11T19:05:00+02:00"
lastmod: "2009-12-11T19:05:00+02:00"
tags: ["nmap", "Network", "Security", "Scan", "Linux"]
toc: true
---

## Introduction

There are several solutions to find which hosts are up on your network. A simple solution is to use nmap.

## Usage

Here's how to do it:

```bash
nmap -sP your network/submask | awk "/^Host/"'{ print $2 }'
```
