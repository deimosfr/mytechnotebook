---
weight: 999
url: "/apt-cacher-ng-mise-en-place-d-un-proxy-pour-apt/"
title: "Apt-cacher-ng: Setting Up a Proxy for APT"
description: "How to set up and configure Apt-cacher-ng, a caching proxy for Debian/Ubuntu package repositories"
categories: ["Linux", "Debian", "Server"]
date: "2010-08-31T07:51:00+02:00"
lastmod: "2010-08-31T07:51:00+02:00"
tags: ["Apt", "Proxy", "Caching", "Debian", "Ubuntu"]
toc: true
---

## Introduction

Without having an amazing internet connection and using virtual machines, it quickly becomes tedious to wait ages to download the same packages for all these VMs. So I looked for a solution to have a kind of cache specifically for apt.

I found happiness with apt-cacher-ng which, from what I've read, is the most complete solution available today.

## Installation

As usual, it's simple:

```bash
aptitude install apt-cacher-ng
```

My server has the IP **192.168.100.1**. If you want to modify some small options, I invite you to look at the file `/etc/apt-cacher-ng/acng.conf`.

Then restart the service.

## Configuration

This part is fairly easy, as we'll tell the apt service to use a proxy. We'll set it up on both the server and clients:

```bash
Acquire::http { Proxy "http://192.168.100.1:3142"; };
```

Then we update everything:

```bash
aptitude update
```

## Usage

By default, there's a small graphical interface that allows you to update, import, verify, manage expiration, etc. The address is: http://192.168.100.1:3142/acng-report.html
