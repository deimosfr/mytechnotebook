---
weight: 999
url: "/installation_et_configuration_de_cacti/"
title: "Installation and Configuration of Cacti"
description: "How to install and configure Cacti for network monitoring"
categories: ["Monitoring", "Network"]
date: "2007-12-12T11:01:00+02:00"
lastmod: "2007-12-12T11:01:00+02:00"
tags: ["cacti", "monitoring", "snmp", "rrdtool"]
toc: true
---

## Introduction

Before starting, you should know that the installation of [cacti](https://cacti.net/) requires a [MySQL database]({{< ref "docs/Servers/Databases/MySQL-MariaDB">}}), a [web server]({{< ref "docs/Servers/Web/Apache">}}) (apache with PHP) and the SNMP protocol.

Cacti is used for specialized network monitoring since the MRTG team abandoned this project to work on Cacti. That said, Cacti doesn't just monitor networks, you can monitor all types of services.

## Installation

Before installing these programs, enable SNMP on your router. Then, we'll install cacti and rrdtool:

```bash
apt-get install rrdtool cacti
```

In Debian, it will ask you for the name of your MySQL database administrator and password. Then it will ask you for the name of the database that Cacti will use. Choose what you want, then it will ask you to enter a username and password for the cacti user. Your database is ready.

Now you need to create a symbolic link to put it on the website. We'll do this using the ln -s command (in case it's not done automatically):

```bash
ln -s /usr/share/cacti/ /var/www/cacti
```

## Configuration

For Cacti to correctly update the configurations you're going to enter, you need to enter the crontab of the user who manages your apache server (www-data for debian) and add this line. It will update cacti every 5 minutes:

```bash
*/5 * * * * php /usr/share/cacti/cmd.php > /dev/null 2>&1
```

## Launch

Wait a bit for the script to run, then open your browser and type the server address followed by cacti (e.g.: http://127.0.0.1/cacti). It will ask you for a login and password. Enter "admin" for both.

By default, you have some services configured. Up to you if you want to configure others... :-)
