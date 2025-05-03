---
weight: 999
url: "/Limesurvey_\\:_Mise_en_place_d'une_solution_de_Sondages/"
title: "Limesurvey: Setting up a Survey Solution"
description: "Learn how to install and configure Limesurvey, a complete survey solution for your web server, allowing you to create and manage sophisticated surveys."
categories: ["Debian", "Database", "Linux"]
date: "2012-11-14T10:10:00+02:00"
lastmod: "2012-11-14T10:10:00+02:00"
tags: ["Limesurvey", "Survey", "PostgreSQL", "Web application", "PHP"]
toc: true
---

![Limesurvey](/images/limesurvey_logo.avif)

{{< table "table-hover table-striped" >}}
|||
|-|-|
| **Software version** | 2.00+ (Build 121115) |
| **Operating System** | Debian 6 + backports |
| **Website** | [Limesurvey Website](https://www.limesurvey.org/) |
| **Last Update** | 14/11/2012 |
{{< /table >}}

## Introduction

LimeSurvey is a web application that is installed on the user's server. After installation, the user can manage LimeSurvey from a web interface. It provides a complete text editor to write questions and messages, and also allows integration of images and videos into surveys. The layout and design of surveys can be modified by changing the template. Templates can be modified using a WYSIWYG (What You See Is What You Get) HTML editor.

Additionally, templates can be easily imported and exported through the template editor. Once a survey is complete, the user can activate it, making it available to all. Similarly, you can import and export questions through the interface editor. LimeSurvey allows you to create as many surveys as desired. There is also no limit on the number of invited participants. Apart from technical and practical constraints, there are no limits on the number of questions each survey can have.

Questions are added by group. Questions in the same group are displayed on the same page. Surveys can contain many different question types: lists, multiple choice, text, numeric, as well as simple "yes" or "no" answers. Questions can be organized with arrows, with options for questions on one axis based on the other axis. Questions can also depend on answers to previous questions. For example, a voter can answer a question about transportation if they answered affirmatively to a question about employment.[^1]

This tutorial is based on the latest stable version of Limesurvey.

## Installation

For the installation, we need the following:

```bash
aptitude install postgresql-8.4 apache2 apache2 libapache2-mod-php5 php5 php5-gd php5-imap php5-ldap php5-pgsql
```

Next, we'll download the latest version of Limesurvey, extract it, and set the proper permissions:

```bash
cd /var/www
wget -O limesurvey.tgz http://www.limesurvey.org/fr/stable-release/finish/25-latest-stable-release/686-limesurvey200plus-build121115targz
tar -xzvf limesurvey.tgz
chown -Rf www-data. limesurvey
```

## Configuration

We'll create a PostgreSQL user and database. To begin, let's configure the authentication part:

```bash {linenos=table,hl_lines=["7-9"]}
[...]
# Database administrative login by UNIX sockets
local   all         postgres                          ident

# TYPE  DATABASE    USER        CIDR-ADDRESS          METHOD

# Limesurvey
local   limesurvey    limesurvey    md5
host    limesurvey    limesurvey    127.0.0.1/32    md5
# "local" is for Unix domain socket connections only
local   all         all                               ident
# IPv4 local connections:
host    all         all         127.0.0.1/32          md5
# IPv6 local connections:
host    all         all         ::1/128               md5
```

Now let's create users, databases, and grant access:

```bash
su postgres
psql
create user limesurvey password 'limesurvey' nosuperuser;
create database limesurvey owner limesurvey;
```

Replace the password part with the password you desire.

Then we restart everything to ensure the new configuration is active:

```bash
service postgresql restart
```

For the Limesurvey configuration part, it's simple, everything is done via the wizard: http://server/limesurvey

Follow the instructions and that's it, all that's left is to use it :-)

## References

[^1]: http://fr.wikipedia.org/wiki/LimeSurvey
