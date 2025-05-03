---
weight: 999
url: "/Cron-apt_\\:_Installation_des_mises_à_jour_de_sécurité_automatique/"
title: "Cron-apt: Automatic Security Updates Installation"
description: "How to setup automatic security updates on Debian using cron-apt"
categories: ["Security", "Debian", "Linux"]
date: "2014-05-18T11:34:00+02:00"
lastmod: "2014-05-18T11:34:00+02:00"
tags: ["Servers", "Security", "Automation", "Updates"]
toc: true
---

## Introduction

My goal is to install security updates automatically. Obviously this kind of approach is not really recommended, but on a Debian stable system where only security updates are installed, we minimize the risks.

So I started looking for a tool to do this kind of thing and found cron-apt.

## Installation

Simply:

```bash
aptitude install cron-apt
```

## Configuration

Now that it's installed, let's create a file that will contain only the Debian security repositories:

```bash
grep security /etc/apt/sources.list > /etc/apt/security.sources.list
```

Then let's edit the cron-apt configuration file to use aptitude, send update status emails, and specify that we only want security updates:

```bash {linenos=table}
APTCOMMAND=/usr/bin/aptitude
OPTIONS="-o quiet=1 -o Dir::Etc::SourceList=/etc/apt/security.sources.list"
MAILTO="xxx@mycompany.com"
MAILON="always"
```

Then we just need to modify the default actions to perform. By default, it only downloads packages without installing them (because of the -d option on the dist-upgrade line). That's why we are going to modify this file accordingly:

```bash {linenos=table}
autoclean -y
dist-upgrade -y -o APT::Get::Show-Upgraded=true
```

Finally, if you want to change the update time, check this file and adapt it according to your needs:

```bash {linenos=table}
#
# Regular cron jobs for the cron-apt package
#
# Every night at 4 o'clock.
0 4	* * *	root	test -x /usr/sbin/cron-apt && /usr/sbin/cron-apt
# Every hour.
# 0 *	* * *	root	test -x /usr/sbin/cron-apt && /usr/sbin/cron-apt /etc/cron-apt/config2
# Every five minutes.
# */5 *	* * *	root	test -x /usr/sbin/cron-apt && /usr/sbin/cron-apt /etc/cron-apt/config2
```
