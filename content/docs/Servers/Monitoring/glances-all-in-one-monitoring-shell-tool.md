---
weight: 999
url: "/glances-all-in-one-monitoring-shell-tool/"
title: "Glances: All in One Monitoring Shell Tool"
description: "Glances is a cross-platform curses-based monitoring tool written in Python that provides an all-in-one overview of your system health, replacing the need to run multiple monitoring tools simultaneously."
categories: ["Linux", "Monitoring", "System Administration"]
date: "2013-08-11T19:26:00+02:00"
lastmod: "2013-08-11T19:26:00+02:00"
tags: ["glances", "monitoring", "system", "python", "dashboard", "shell", "debian"]
toc: true
---

![Glances](/images/glances-logo.avif)

{{< table "table-hover table-striped" >}}
|||
|-|-|
| **Software version** | 1.6.1/Latest |
| **Operating System** | Debian 7 |
| **Website** | [Glances Website](https://github.com/nicolargo/glances) |
| **Last Update** | 11/08/2013 |
{{< /table >}}

## Introduction

[Glances](https://github.com/nicolargo/glances) is a cross-platform curses-based monitoring tool written in Python.

It avoids having to run several tools to get an all-in-one overview of your system. For example, when you want to quickly see what's wrong on a system for diagnosis, you'll need to launch top/htop, iostat, vmstat...Glances gives you a large overview of your system health. You can then investigate with the appropriate tool if you want. But you didn't waste your time opening several tools to get the first desired information: where does the problem come from? So Glances answers that question.

![Glances screenshot](/images/glances_screenshot.avif)

## Installation

### Packages

The glances packages are not yet available in Debian wheezy packages. But they are in Jessie!

{{< alert context="warning" text="This will upgrade the libc6!" />}}

That's why we can do APT pinning to use packages:

(`/etc/apt/preferences`)

```bash
Package: *
Pin: release a=wheezy
Pin-priority: 900

Package: *
Pin: release a=jessie
Pin-priority: 100

Package: glances
Pin: release a=jessie
Pin-priority: 1001
```

Add as well the jessie repositories to your current wheezy:

(`/etc/apt/sources.list`)

```bash
deb http://ftp.fr.debian.org/debian/ jessie main contrib non-free
deb-src http://ftp.fr.debian.org/debian/ jessie main contrib non-free
```

Now update and install glances:

```bash
aptitude update
aptitude install glances
```

### Pip

You can install the latest version of Glances using pip. First install dependencies:

```bash
aptitude install python-pip python-dev
```

Then install Glances:

```bash
pip install Glances
```

You're now able to launch glances in command line:

```bash
glances
```

#### Upgrade

To upgrade your Glances version:

```bash
pip install --upgrade glances
```

## Configuration

There is nothing especially to configure as default options could be enough for a large set of users. Anyway, you can add or change several options of the default configuration in `/etc/glances.glances.conf`.

An option that comes with 1.7 version is the possibility to watch a specific software. Let's say Nginx for instance. You can ask glance to look at it by adding these lines:

(`/etc/glances.glances.conf`)

```ini
[monitor]
list_1_description=Web Nginx Server
list_1_regex=.*nginx.*
list_1_command=nginx -v
list_1_countmin=1
list_1_countmax=4
```

* list_X: replace X by 1 to 9, this is the information for additional software (here Nginx)
* description: set the software description (16 chars max)
* regex: regex to group software information
* command: the command to run that will show information in glances
* countmin: minimum number of information to show
* countmax: maximum number of information to show

You can add other software by adding same lines with list_2, list_3...

## Client/Server mode

There is also a client/server mode. You can run a server like this:

```bash
glances -s
```

And you can connect clients:

```bash
glances -c <server>
```

Replace <server> by the server IP.

## References

1. [https://github.com/nicolargo/glances](https://github.com/nicolargo/glances)
