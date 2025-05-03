---
weight: 999
url: "/monit-easily-use-triggers-on-your-system/"
title: "Monit: Easily Use Triggers on Your System"
description: "Learn how to install, configure, and use Monit for managing and monitoring processes, programs, files, directories, and filesystems on Unix systems with automatic maintenance and error handling."
categories: ["Linux", "Monitoring", "System Administration"]
date: "2014-05-28T09:11:00+02:00"
lastmod: "2014-05-28T09:11:00+02:00"
tags: ["monit", "monitoring", "system", "triggers", "alerts", "automation", "process management"]
toc: true
---

![Monit](/images/monit_logo.avif)

{{< table "table-hover table-striped" >}}
|||
|-|-|
| **Software version** | 5.4-2 |
| **Operating System** | Debian 7 |
| **Website** | [Monit Website](https://mmonit.com) |
| **Last Update** | 28/05/2014 |
{{< /table >}}

## Introduction

[Monit](https://mmonit.com/monit/documentation/monit.html) is a utility for managing and monitoring processes, programs, files, directories and filesystems on a Unix system. Monit conducts automatic maintenance and repair and can execute meaningful causal actions in error situations. E.g. Monit can start a process if it does not run, restart a process if it does not respond and stop a process if it uses too much resources. You can use Monit to monitor files, directories and filesystems for changes, such as timestamps changes, checksum changes or size changes.

Monit is controlled via an easy to configure control file based on a free-format, token-oriented syntax. Monit logs to syslog or to its own log file and notifies you about error conditions via customizable alert messages. Monit can perform various TCP/IP network checks, protocol checks and can utilize SSL for such checks. Monit provides a http(s) interface and you may use a browser to access the Monit program.

## Installation

To install Monit, this is simple:

```bash
aptitude install monit
```

## Configuration

Regarding the configuration file, you've got a global configuration file where you can adjust some parameters:

(`/etc/monit/monitrc`)

```bash
set daemon 120            # check services at 2-minute intervals
set logfile /var/log/monit.log
set idfile /var/lib/monit/id
set statefile /var/lib/monit/state
set mailserver localhost
set eventqueue
basedir /var/lib/monit/events # set the base directory where events will be stored
slots 100                     # optionally limit the queue size
set alert my@email.com
set httpd port 2812 and
allow localhost        # allow localhost to connect to the server and
include /etc/monit/conf.d/*
```

But here is a configuration file to restart multiple services in a shell script when an URL is not containing a specific content:

(`/etc/monit/conf.d/web`)

```bash
check host blog.deimos.fr with address blog.deimos.fr
    if failed (url http://blog.deimos.fr and content == 'Because human memory can not contain Gb')
    with timeout 20 seconds for 3 cycles
    then exec "/etc/scripts/web_services.sh restart"
    alert my@email.com
```

There are several type of usages with Monit and you can see examples [here](https://mmonit.com/monit/documentation/monit.html#configuration_examples).

### Web interface

Regarding the web interface, you can use Nginx and do a proxy pass to access from outside with credentials:

(`/etc/nginx/sites-enabled/default`)

```bash
    # Monit status
    location /monit/ {
        rewrite ^/monit/(.*) /$1 break;
        proxy_ignore_client_abort on;
        proxy_pass   http://127.0.0.1:2812; 
        # User access
        auth_basic "Please logon"; 
        auth_basic_user_file /etc/nginx/access/htaccess;
    }
```

And create the htaccess file with credentials.

## References

1. [https://mmonit.com/monit/documentation/monit.html](https://mmonit.com/monit/documentation/monit.html)
2. [https://mmonit.com/monit/documentation/monit.html#configuration_examples](https://mmonit.com/monit/documentation/monit.html#configuration_examples)
