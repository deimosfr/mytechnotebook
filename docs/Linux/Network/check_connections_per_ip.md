---
title: "Check Connections Per IP"
slug: check-connections-per-ip/
description: "How to check the number of connections per IP address on Linux and BSD systems"
categories: ["Linux", "BSD", "Network", "Security"]
date: "2009-12-06T16:35:00+01:00"
lastmod: "2009-12-06T16:35:00+01:00"
tags: ["connections", "network", "netstat", "monitoring"]
---

## Introduction

Here is a command line to run on your server if you think your server is under attack. It prints out a list of open connections to your server and sorts them by amount.

## Usage

### Linux

```bash
netstat -ntu 
```

### BSD

```bash
netstat -na 
```
