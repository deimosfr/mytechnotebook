---
weight: 999
url: "/Trickle_\\:_limit_your_application_bandwidth/"
title: "Trickle: Limit Your Application Bandwidth"
description: "Guide on how to use Trickle to limit bandwidth usage for applications in Linux"
categories: ["Linux", "Debian", "Network"]
date: "2013-05-06T13:47:00+02:00"
lastmod: "2013-05-06T13:47:00+02:00"
tags: ["Bandwidth", "Network", "Performance", "Linux"]
toc: true
---

![Linux](/images/poweredbylinux.avif)

{{< table "table-hover table-striped" >}}
|||
|-|-|
| **Software version** | 1.07 |
| **Operating System** | Debian 7 |
| **Website** | [Trickle Website](https://monkey.org/~marius/pages/?page=trickle) |
| **Last Update** | 06/05/2013 |
{{< /table >}}

## Introduction

trickle is a portable lightweight userspace bandwidth shaper. It can run in collaborative mode (together with trickled) or in stand alone mode.

trickle works by taking advantage of the unix loader preloading. Essentially it provides, to the application, a new version of the functionality that is required to send and receive data through sockets. It then limits traffic based on delaying the sending and receiving of data over a socket. trickle runs entirely in userspace and does not require root privileges.

## Installation

```bash
aptitude install trickle
```

## Usage

To use it, it's really easy. Simply add a command prefix with desired max bandwidth:

```bash
trickled -u 10 -d 20 firefox
```

Here, the maximum upload is set to 10k and upload to 20k.
