---
weight: 999
url: "/PuppetDB_\\:_Augmentez_les_fonctionnalités_de_votre_Puppet/"
title: "PuppetDB: Enhance Your Puppet Functionality"
description: "Learn how to install and configure PuppetDB to enhance Puppet functionality by collecting data and enabling exported resources."
categories: ["Debian", "PostgreSQL", "Database"]
date: "2013-10-28T10:57:00+02:00"
lastmod: "2013-10-28T10:57:00+02:00"
tags:
  [
    "Puppet",
    "PuppetDB",
    "Debian",
    "PostgreSQL",
    "Configuration Management",
    "Database",
    "Java",
  ]
toc: true
---

![Puppet](/images/puppet-short.avif)

{{< table "table-hover table-striped" >}}
|||
|-|-|
| **Software version** | 1.5.2 |
| **Operating System** | Debian 7 |
| **Website** | [Puppet Website](https://puppetlabs.com/) |
| **Last Update** | 28/10/2013 |
{{< /table >}}

## Introduction

PuppetDB[^1] allows you to retrieve data collected by [Puppet](./puppet_:_solution_de_gestion_de_fichier_de_configuration.html) such as facts and to use exported resources among other things. This data can then be used by other programs, such as the [dashboard](./Puppet_Dashboard_:_Mise_en_place_d'une_interface_graphique_pour_Puppet.html), or your own tools through an API. You can install the PuppetDB server on your PuppetMaster or on a separate server.[^2]

Today it is possible to use 2 backends to store this data:

- HSQLDB: in-memory database, with quite a few limitations including the 100 nodes maximum, but extremely fast (as it's loaded in RAM)
- PostgreSQL: classic database, with less performance (on disk), but with more flexibility and a greater possibility of extension over time (more than 100 Puppet nodes).

We will proceed with the PostgreSQL-based solution.

If you want more information, check out [the different bottlenecks](https://docs.puppetlabs.com/puppetdb/1/scaling_recommendations.html). In summary: if you have more than 100 clients, you'll need a PostgreSQL database, increase your JVM, and increase the number of CPUs/cores.

## Installation

### PuppetDB

Let's start by setting up what we need to install PuppetDB:

```bash
wget http://apt.puppetlabs.com/puppetlabs-release-stable.deb
dpkg -i puppetlabs-release-stable.deb
```

And then update:

```bash
aptitude update
```

Then you need to have puppet installed on the same machine where PuppetDB will be installed.

{{< alert context="warning" text="PuppetDB cannot function without the puppet client!" />}}

It doesn't matter whether it's the same machine as the master or not. For simplicity, we'll install it on the Puppet Master.

```bash
aptitude install puppetdb puppet
```

### PostgreSQL

Since we plan to have more than 100 clients (or even if it's less, we don't want to change the configuration in the future), we'll install PostgreSQL:

```bash
aptitude install postgresql
```

### Terminus

**On the Puppet Master server**, install this package:

```bash
aptitude install puppetdb-terminus
```

## Configuration

### PostgreSQL

We'll create a user and a database:

```bash
su - postgres
createuser -DRSP puppetdb
createdb -O puppetdb puppetdb
exit
```

### PuppetDB

First, configure the database information:

```ini {linenos=table,hl_lines=[5,10,15,18,21]}
[database]
# For the embedded DB: org.hsqldb.jdbcDriver
# For PostgreSQL: org.postgresql.Driver
# Defaults to embedded DB
classname = org.postgresql.Driver

# For the embedded DB: hsqldb
# For PostgreSQL: postgresql
# Defaults to embedded DB
subprotocol = postgresql

# For the embedded DB: file:/path/to/database;hsqldb.tx=mvcc;sql.syntax_pgs=true
# For PostgreSQL: //host:port/databaseName
# Defaults to embedded DB located in <vardir>/db
subname = //localhost:5432/puppetdb

# Connect as a specific user
username = puppetdb

# Use a specific password
password = puppetdb

# How often (in minutes) to compact the database
# gc-interval = 60

# Number of seconds before any SQL query is considered 'slow'; offending
# queries will not be interrupted, but will be logged at the WARN log level.
log-slow-statements = 10
```

Here we define the database connection properties, as well as the credentials we created earlier.

Now let's configure the number of threads:

```ini {linenos=table,hl_lines=[15]}
# See README.md for more thorough explanations of each section and
# option.

[global]
# Store mq/db data in a custom directory
vardir = /var/lib/puppetdb
# Use an external log4j config file
logging-config = /etc/puppetdb/conf.d/../log4j.properties

# Maximum number of results that a resource query may return
resource-query-limit = 20000

[command-processing]
# How many command-processing threads to use, defaults to (CPUs / 2)
threads = 4
```

Adjust the number of threads to your processor count divided by 2.

Then, we tackle the Jetty configuration:

```ini {linenos=table,hl_lines=[3,7]}
[jetty]
# Hostname to list for clear-text HTTP.  Default is localhost
host = 0.0.0.0
# Port to listen on for clear-text HTTP.
port = 8080

ssl-host = 0.0.0.0
ssl-port = 8081
keystore = /etc/puppetdb/ssl/keystore.jks
truststore = /etc/puppetdb/ssl/truststore.jks

key-password = CoaRwY6IL8KQd8H6SfZ7O9hHC
trust-password = CoaRwY6IL8KQd8H6SfZ7O9hHC
```

For the host, add the interface that will listen on ports 8080 and 8081. This notably allows the dashboard to connect to it.

{{< alert context="warning" text="If possible and as a security measure, leave everything on localhost. Obviously, Puppet Master must be on this same machine if host and host_ssl are set to localhost" />}}

Then restart PuppetDB:

```bash
/etc/init.d/puppetdb restart
```

After a few seconds/minutes, you should be able to connect to port 8081 (ssl) or 8080 (non-ssl) (http://<hostname>:8080|https://<hostname>:8081), where you'll have access to a nice interface:

![Puppetdb dashboard](/images/puppetdb_dashboard.avif)

### Puppet Master

On the master, you need to modify its configuration:

```ini {linenos=table,hl_lines=["15-16","21"]}
[main]
logdir=/var/log/puppet
vardir=/var/lib/puppet
ssldir=/var/lib/puppet/ssl
rundir=/var/run/puppet
factpath=$vardir/lib/facter
templatedir=$confdir/templates
pluginsync = true

[master]
# These are needed when the puppetmaster is run by passenger
# and can safely be removed if webrick is used.
ssl_client_header = SSL_CLIENT_S_DN
ssl_client_verify_header = SSL_CLIENT_VERIFY
storeconfigs = true
storeconfigs_backend = puppetdb
report = true

[agent]
server=puppet-prd.deimos.fr
```

{{< alert context="warning" text="Remove the thin_storeconfigs and async_storeconfigs lines if you are using them, or set them to False" />}}

Then we'll set up a file for Puppet's configuration to tell it how to connect to PuppetDB:

```ini {linenos=table,hl_lines=["2"]}
[main]
server = puppet-prd.deimos.fr
port = 8081
```

And finally, a file to define the location of facts:

```yaml {linenos=table,hl_lines=["4"]}
---
master:
  facts:
    terminus: puppetdb
    cache: yaml
```

That's it, your Puppet server now has a working PuppetDB backend! :-)

# FAQ

## I have OutOfMemoryError errors in my logs and PuppetDB responds slowly

To confirm that the problem is indeed due to insufficient memory, check that the file `/var/log/puppetdb/puppetdb-oom.hprof` exists, and make sure the content mentions OOM.

You'll need to increase the Java Heap size (Xmx value) of your PuppetDB which requires more RAM. Increase this value:

```ini {linenos=table,hl_lines=[9]}
###########################################
# Init settings for puppetdb
###########################################

# Location of your Java binary (version 6 or higher)
JAVA_BIN="/usr/bin/java"

# Modify this if you'd like to change the memory allocation, enable JMX, etc
JAVA_ARGS="-Xmx192m -XX:+HeapDumpOnOutOfMemoryError -XX:HeapDumpPath=/var/log/puppetdb/puppetdb-oom.hprof "

# These normally shouldn't need to be edited if using OS packages
USER="puppetdb"
INSTALL_DIR="/usr/share/puppetdb"
CONFIG="/etc/puppetdb/conf.d"
```

For an idea of the value to set (where n represents the number of nodes):

128M + (1M \* n)

# References

[^1]: http://docs.puppetlabs.com/puppetdb/1/index.html
[^2]: http://binbash.fr/2012/11/14/puppet-3-et-puppetdb/
