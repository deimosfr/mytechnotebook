---
weight: 999
url: "/Bacula_\\:_Mise_en_place_d'un_serveur_de_Backup_Performant/"
title: "Bacula: Setting Up a High-Performance Backup Server"
description: "Learn how to set up Bacula, a powerful backup solution with web interface, MySQL plugin, and file verification capabilities."
categories: ["Server", "Backup", "Linux"]
date: "2008-05-14T15:31:00+02:00"
lastmod: "2008-05-14T15:31:00+02:00"
tags: ["bacula", "backup", "mysql", "database"]
toc: true
---

## Introduction

[Bacula](https://www.bacula.org) is a backup server that allows for the purchase of other commercial solutions with plugins for additional features. It also has the advantage of having a nice graphical web interface that creates beautiful graphs, etc.

In short, its competitor BackupPC, which I used for a long time, is very good, but doesn't have all the [advantages of Bacula](https://www.bacula.org/fr/dev-manual/etat_actuel_Bacula.html) such as checksum verification on files or a MySQL plugin.

## Prerequisites

You'll need to set up some services before taking the plunge:

- A [MySQL]({{< ref "docs/Servers/Databases/MySQL-MariaDB/mysql_installation_and_configuration.md" >}}) server
- A [Web server like Apache]({{< ref "docs/Servers/Web/Apache/apache_2_installation_and_configuration.md" >}}) with PHP

## Installation

We'll do it under Debian, let's stick to our good habits:

```bash
apt-get install bacula bacula-client bacula-common bacula-console bacula-director-common bacula-director-mysql bacula-fd bacula-sd bacula-sd-mysql bacula-server
```

By default, Bacula installs with SQLite to make it quicker to set up, but the packages above are for MySQL. I chose to use this database server here because for the web graphical interface, you absolutely need to use MySQL or PostgreSQL. For simplicity, I chose MySQL :-)

During the installation of the packages, you'll be asked for the SQL server password so it can create everything it needs (database + user).

## Configuration

Here's a diagram explaining the positioning of the configuration files. This helps to better plan how to arrange the structure of these files:

![Bacula-objects](/images/bacula-objects.avif)

And here's how the files behave between themselves:

![Conf-Diagram](/images/conf-diagram.avif)

### Language

The configuration is not always obvious, which is why a little explanation does no harm to understand this language:

- FileSet: What to backup?
- Client: Who to backup?
- Schedule: When to backup?
- Pool: Where to backup? (i.e., On which volume?)
- Volume: a simple physical medium (cartridge, or simple file) on which Bacula writes your backup data
- Pools: group Volumes so that a backup is not restricted to the capacity of a single Volume

Although the basic options are specified in the Pool resource in the Director's configuration file, the actual Pool is managed by the Bacula Catalog. It contains information from the Pool resource (bacula-dir.conf) as well as information about all Volumes that have been added to the Pool. Volumes are normally added manually from the Console using the label command.

For Bacula to read or write to a physical Volume, it must be software-labeled so that Bacula is assured that the correct Volume is mounted. This is normally done manually from the Console using the label command.

- Director: To define the name of the Director and its password for authentication of the Console program. There should be only one Director resource definition in the configuration file. If you have either /dev/random or bc on your machine, Bacula will generate a random password during the configuration process, otherwise, it will be left blank.
- Job: To define backup and restore Jobs, and to link Client, FileSet, and Schedules resources to be used together for each Job.
- JobDefs: Optional resource designed to provide default values for Job resources.
- Schedule: To define when a Job should be automatically launched by Bacula's internal scheduler.
- FileSet: To define the set of files to be backed up for each client.
- Client: To define which Client is to be backed up.
- Storage: To define on which physical device volumes will be mounted.
- Pool: To define which pool of volumes can be used for a given Job
- Catalog: To define the database to keep lists of backed up files and volumes where they were backed up.
- Messages: To define the recipients (or log files) of error and information messages.

### Director

We'll modify the catalog here so it can connect to the MySQL database. Edit your file /etc/bacula/bacula-dir.conf and modify these lines to adapt them to your needs:

```
Catalog {
  Name = MySQL
  dbname = bacula
  user = bacula
  password = "bacula_password"
  DB Address = 'localhost
  DB Port = 3306
}
```

To get all the information for the director, I invite you to access this page: [https://www.bacula.org/fr/rel-manual/Configurer_le_Director.html](https://www.bacula.org/fr/rel-manual/Configurer_le_Director.html)
