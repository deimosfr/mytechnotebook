---
weight: 999
url: "/MFi\\:_install_a_Ubiquiti_server_to_manage_powerstrips/"
title: "MFi: Install a Ubiquiti server to manage powerstrips"
description: "A guide on installing and configuring a Ubiquiti mFi server on Debian Linux to manage powerstrips"
categories: ["Networking", "Linux", "Hardware"]
date: "2014-07-27T13:15:00+02:00"
lastmod: "2014-07-27T13:15:00+02:00"
tags: ["Ubiquiti", "mFi", "powerstrip", "Debian", "LXC"]
toc: true
---

{{< table "table-hover table-striped" >}}
|||
|-|-|
| **Software version** | 2.0.24 |
| **Operating System** | Debian 7 |
| **Website** | [Ubiquiti Website](https://www.ubnt.com) |
| **Last Update** | 27/07/2014 |
{{< /table >}}

## Introduction

Ubiquiti brings software that are easy to install on Windows and Mac OS. However as it is strongly recommended to let this software always up, it's preferable to have a Linux version to run it in a container or a virtual machine. That's why I decided to install it for a powerstrip mPower on Debian inside LXC. The documentation is very poor, that's why I made this one for those who want to do like me.

## Installation

First of all, download the mFi server here. Then install prerequisites:

```bash
aptitude install unzip mongodb openjdk-7-jre-headless openjdk-7-jre
```

Then unzip the archive file and move it to a better folder:

```bash
unzip mFi.unix.zip
mv mFi /usr/share/
```

Now set the service to start at boot:

```bash
cd /usr/share/mFi/ ; java -jar lib/ace.jar start &
```

Now you can start the service by hand to check it works fine:

```bash
cd /usr/share/mFi/ ; java -jar lib/ace.jar start
```

## Configuration

On the server, there is nothing to do especially. However on the powerstrip, you need to access to the web interface and configure them as follow:

- Configuration
  - Controller: Personal
  - Address: **ubiquiti-server-IP or DNS**
  - User: <username>
  - Password: <password>

You need to adapt all fields to fit your server interface.

## Usage

You can now access to the web interface like this https://ubiquiti-server-IP:6443 and you should have an interface to configure your powerstrips:

![Mfi screenshot](/images/mfi_screenshot.avif)
