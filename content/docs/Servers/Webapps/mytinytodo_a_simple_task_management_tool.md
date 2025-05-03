---
weight: 999
url: "/MyTinyTodo_\\:_Un_outil_simple_de_gestion_de_tâches/"
title: "MyTinyTodo: A Simple Task Management Tool"
description: "How to install and configure MyTinyTodo, a simple web-based task management tool that's compatible with smartphones."
categories: ["Linux", "Servers"]
date: "2012-04-20T08:19:00+02:00"
lastmod: "2012-04-20T08:19:00+02:00"
tags: ["Task Management", "Web Applications", "Servers"]
toc: true
---

![MyTinyTodo](/images/logomytodolist.avif)

{{< table "table-hover table-striped" >}}
|||
|-|-|
| **Software version** | 1.4.2 |
| **Operating System** | Debian 6 |
| **Website** | [MyTinyTodo Website](https://www.mytinytodo.net/) |
| **Last Update** | 20/04/2012 |
{{< /table >}}

## Introduction

I've been looking for a long time for a simple tool that I could host myself, web-based and compatible with smartphones. MyTinyTodo is made for that. Setting it up is very easy, as you'll see...

## Installation

```bash
wget http://www.mytinytodo.net/latest.php
unzip latest.php
chown -Rf www-data. mytinytodo
rm -f latest.php
```

## Configuration

You'll need to create a user and database beforehand, then start the installation by connecting to your website, for example: http://www.deimos.fr/mytinytodo

Once the installation is complete, delete the setup file:

```bash
rm -f mytinytodo/setup.php
```

## References

http://www.mytinytodo.net/
