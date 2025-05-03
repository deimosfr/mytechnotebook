---
weight: 999
url: "/PNP4Nagios\\:_Grapher_ses_alertes_Nagios/"
title: "PNP4Nagios: Graph Your Nagios Alerts"
description: "Guide to install and configure PNP4Nagios to generate graphs from Nagios performance data"
categories: ["Linux", "Monitoring", "Apache"]
date: "2011-08-26T15:43:00+02:00"
lastmod: "2011-08-26T15:43:00+02:00"
tags: ["Nagios", "PNP4Nagios", "Monitoring", "RRD", "Performance Graphs"]
toc: true
---

![PNP4Nagios](/images/pnp4nagios_logo.avif)

## Introduction

PNP is an addon to Nagios which analyzes performance data provided by plugins and stores them automatically into RRD-databases (Round Robin Databases, see RRD Tool).

## Prerequisites

First, create a folder where to unzip the installation files. Here, we are going to use `/etc/nagios3/pnp4nagios`:

```bash
mkdir /etc/nagios3/pnp4nagios
cd /etc/nagios3/pnp4nagios
```

Then get the file, unzip and go into the folder:

```bash
wget http://downloads.sourceforge.net/project/pnp4nagios/PNP-0.6/pnp4nagios-0.6.0.tar.gz?use_mirror=freefr
cd pnp4nagios-0.6.0
```

Activate the Apache2 rewrite module:

```bash
a2enmod rewrite
/etc/init.d/apache2 restart
```

Edit the php5 conf file in `/etc/php5` on this line:

```php
magic_quotes_gpc = Off
```

If not installed yet, install GCC (C compiler)

```bash
aptitude install gcc
```

## Installation

Configure the installation, using the folder you want to install pnp4nagios in:

```bash
./configure --prefix=/etc/nagios3/pnp4nagios --with-nagios-user=nagios --with-nagios-group=nagios
```

Launch every make command:

```bash
make all
make install
make install-webconf
make install-config
make install-init
```

## Configuration

In the nagios generic services file add:

```bash
...
define command {
       command_name    process-service-perfdata
       command_line    /usr/bin/perl /usr/local/pnp4nagios/libexec/process_perfdata.pl
}
```

And in the generic hosts file add:

```bash
define command {
       command_name    process-host-perfdata
       command_line    /usr/bin/perl /usr/local/pnp4nagios/libexec/process_perfdata.pl -d HOSTPERFDATA
}
```

Edit the nagios.cfg file and be sure those line are uncommented:

```bash
process_performance_data=1
service_perfdata_command=process-service-perfdata
host_perfdata_command=process-host-perfdata
```

In any service or host configuration file, add (use generic-*.cfg for example):

```bash
define host {
   name       host-pnp
   action_url /pnp4nagios/index.php/graph?host=$HOSTNAME$&srv=_HOST_' class='tips' rel='/pnp4nagios/index.php/popup?host=$HOSTNAME$&srv=_HOST_
   register   0
}

define service {
   name       srv-pnp
   action_url /pnp4nagios/index.php/graph?host=$HOSTNAME$&srv=$SERVICEDESC$' class='tips' rel='/pnp4nagios/index.php/popup?host=$HOSTNAME$&srv=$SERVICEDESC$
   register   0
}
```

Now for any host and service you want to use pnp4nagios with, on its "use" line, add this for a host:

```
host-pnp
```

and for a service:

```
srv-pnp
```

Example:

```bash
...
define service{
        use                     generic-services,srv-pnp
```

You can add it on a generic host or service used by a group of hosts or services. This will be inherited.

### Onmouseover graphs in cgi interface

Copy the ssi files in the pnp4nagios contrib folder (pnp4nagios-0.6.0/contrib/ssi/*) in the nagios ssi folder (/usr/share/nagios3/htdocs/ssi/):

```bash
cp -Rf /etc/nagios3/pnp4nagios/pnp4nagios-0.6.0/contrib/ssi/*.ssi /usr/share/nagios3/htdocs/ssi/
rm -Rf /etc/nagios3/pnp4nagios/pnp4nagios-0.6.0/
```

Finally, restart apache2 and nagios and everything will work fine. You can see the PNP4Nagios interface to: http://nagios-srv/pnp4nagios
