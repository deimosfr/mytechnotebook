---
title: "PowerDNS: Creating a Master DNS Server"
slug: powerdns-creating-a-master-dns-server/
description: "How to install and configure a PowerDNS master server using a MySQL backend on Debian"
categories: ["MySQL", "Debian", "Linux"]
date: "2012-05-15T14:45:00+02:00"
lastmod: "2012-05-15T14:45:00+02:00"
tags: ["PowerDNS", "DNS", "MySQL", "Servers", "Network"]
---

![PowerDNS](../../../static/images/powerdns_logo.avif)


|||
|-|-|
| **Software version** | 2.9.22 |
| **Operating System** | Debian 6 |
| **Website** | [PowerDNS Website](https://www.powerdns.com) |


## Introduction

PowerDNS is (as its name suggests) a DNS server. It's a direct competitor to Bind. It aims to be less memory-intensive and offers more flexible configuration options than Bind.

PowerDNS is divided into several roles:
- Master
- Cache

In this guide, we'll cover the master server configuration. If you wish to set up a PowerDNS cache server, please [follow this link](../PowerDNS/powerdns_creating_a_master_dns_server.md).

## Installation

First, we'll install a MySQL database (unless you already have another database you wish to use as a backend):

```bash
aptitude install mysql-server
```

Then we'll install PowerDNS:

```bash
aptitude install pdns-server pdns-backend-mysql
```

## Configuration

### MySQL

First, let's create the database:

```bash
mysqladmin -uroot -p create pdns
```

Then we'll create the tables, indexes and assign the permissions:

```bash
mysql -uroot -p pdns < /usr/share/doc/pdns-backend-mysql/mysql.sql
```

### PowerDNS

Now let's configure PowerDNS. We'll specify that we're going to use a MySQL backend:

``` bash hl_lines="5"
[...]
#################################
# launch        Which backends to launch and order to query them in
#
launch=gmysql
[...]
```

Then we'll provide the previously configured information:

``` bash hl_lines="4 5 6 7 8"
# Here come the local changes the user made, like configuration of 
# the several backends that exist.

# MySQL Configuration
gmysql-host=127.0.0.1
gmysql-user=pdns
gmysql-password=password
gmysql-dbname=pdns
```

Now restart PowerDNS:

```bash
/etc/init.d/pdns restart
```

You can now configure your DNS zones and records. I strongly recommend using a web interface to help you with this. For example, you can use [PowerAdmin](./poweradmin_an_administration_interface_for_powerdns.md).

## References

http://www.debiantutorials.com/installing-powerdns-as-supermaster-with-slaves/
