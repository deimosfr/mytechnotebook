---
weight: 999
url: "/Utilisation_de_MySQL/"
title: "Using MySQL"
description: "Tips and tricks for MySQL database management including charset conversion and table prefix manipulation"
categories: ["MySQL", "Database", "Linux"]
date: "2013-01-25T16:43:00+02:00"
lastmod: "2013-01-25T16:43:00+02:00"
tags: ["MySQL", "Database", "UTF8", "Table Prefix", "Command Line"]
toc: true
---

## Introduction

Here are some solutions that will hopefully save you time.

## Converting a Latin1 Database to UTF8

Here's the magic command:

```bash
mysqldump --add-drop-table -uroot -p "DB_name" | replace CHARSET=latin1 CHARSET=utf8 | iconv -f latin1 -t utf8 | mysql -uroot -p "DB_name"
```

## Adding a Prefix to All Tables in a Database

Here's how to add a prefix to all tables in a database:

```sql
SELECT Concat('ALTER TABLE ', TABLE_NAME, ' RENAME TO my_prefix_', TABLE_NAME, ';') FROM information_schema.TABLES WHERE table_schema = 'my_database'
```

You just need to replace:
- `my_prefix`: with your desired prefix
- `my_database`: with the desired database

## References

1. [https://steindom.com/articles/adding-prefix-all-tables-mysql-database](https://steindom.com/articles/adding-prefix-all-tables-mysql-database)
