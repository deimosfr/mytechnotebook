---
weight: 999
url: "/Monitorer_en_temps_réel_votre_Apache/"
title: "Real-time monitoring of your Apache server"
description: "How to monitor your Apache server in real-time using various tools like mod_status and apachetop."
categories: ["Debian", "Linux", "Ubuntu"]
date: "2012-01-18T13:44:00+02:00"
lastmod: "2012-01-18T13:44:00+02:00"
tags: ["Apache", "Monitoring", "Web Server", "Tools", "Performance"]
toc: true
---

## Introduction

It can sometimes be useful to monitor an Apache server, especially when there is an abnormal load. That's why we're going to look at some tools to help monitor Apache connections and status in real time.

## Mod Status

Mod status is an Apache module that displays very interesting information. It has the advantage of being native to Apache.

### Installation

To activate the module, it's very simple:

```bash
a2enmod info
```

### Configuration

Then we'll add this configuration to Apache's default VirtualHost:

```apache
# Get Extended informations
ExtendedStatus On
<VirtualHost *:80>
<Location /server-info>
    SetHandler server-info
    Order deny,allow
    Deny from all
    Allow from localhost ip6-localhost
#    Allow from .example.com
</Location>
</VirtualHost>
```

You just need to reload the service.

### Utilization

Now you just need to type your server's URL (the permissions here only allow localhost): http://localhost/server-status

```
Apache Server Status for 127.0.0.1

Server Version: Apache/2.2.14 (Ubuntu) mod_fcgid/2.3.5 mod_python/3.3.1 Python/2.6.5 mod_ssl/2.2.14 OpenSSL/0.9.8k mod_perl/2.0.4 Perl/v5.10.1
Server Built: Nov 3 2011 03:29:23

Current Time: Wednesday, 18-Jan-2012 14:22:34 CET
Restart Time: Wednesday, 18-Jan-2012 13:33:58 CET
Parent Server Generation: 4
Server uptime: 48 minutes 35 seconds
Total accesses: 743 - Total Traffic: 2.2 MB
CPU Usage: u28.08 s1.23 cu0 cs0 - 1.01% CPU load
.255 requests/sec - 782 B/second - 3070 B/request
1 requests currently being processed, 9 idle workers

.____W.._____...................................................
................................................................
................................................................
................................................................

Scoreboard Key:
"_" Waiting for Connection, "S" Starting up, "R" Reading Request,
"W" Sending Reply, "K" Keepalive (read), "D" DNS Lookup,
"C" Closing connection, "L" Logging, "G" Gracefully finishing,
"I" Idle cleanup of worker, "." Open slot with no current process

Srv	PID	Acc	M	CPU 	SS	Req	Conn	Child	Slot	Client	VHost	Request
0-4	-	0/0/109	. 	0.04	1561	0	0.0	0.00	0.43 	127.0.0.1	deimos.fr	OPTIONS * HTTP/1.0
1-4	27635	0/18/71	_ 	1.07	42	0	0.0	0.03	0.18 	shenzi.deimos.fr	deimos.fr	NULL
2-4	27780	0/8/45	_ 	1.38	58	1	0.0	0.01	0.22 	shenzi.deimos.fr	deimos.fr	GET /server-status HTTP/1.1
3-4	27637	0/21/57	_ 	1.06	56	1	0.0	0.03	0.26 	shenzi.deimos.fr	deimos.fr	GET /server-status HTTP/1.1
4-4	28534	0/11/62	_ 	1.11	57	1	0.0	0.00	0.24 	shenzi.deimos.fr	deimos.fr	GET /server-status HTTP/1.1
5-4	27644	0/16/27	W 	11.82	0	0	0.0	0.08	0.09 	88.191.130.125	deimos.fr	GET /server-status HTTP/1.1
6-4	-	0/0/24	. 	0.08	1672	0	0.0	0.00	0.04 	127.0.0.1	deimos.fr	OPTIONS * HTTP/1.0
7-4	-	0/0/43	. 	0.60	1675	0	0.0	0.00	0.14 	127.0.0.1	deimos.fr	OPTIONS * HTTP/1.0
8-4	27647	0/14/62	_ 	0.99	55	1	0.0	0.04	0.10 	shenzi.deimos.fr	deimos.fr	GET /server-status HTTP/1.1
9-4	27633	0/19/58	_ 	2.26	57	1	0.0	0.02	0.12 	shenzi.deimos.fr	deimos.fr	GET /server-status HTTP/1.1
10-4	27634	0/18/57	_ 	2.84	54	1	0.0	0.02	0.10 	shenzi.deimos.fr	deimos.fr	GET /server-status HTTP/1.1
11-4	27648	0/14/54	_ 	0.52	56	1	0.0	0.02	0.13 	shenzi.deimos.fr	deimos.fr	GET /server-status HTTP/1.1
12-4	27782	0/10/30	_ 	1.03	59	1	0.0	0.01	0.05 	shenzi.deimos.fr	deimos.fr	GET /server-status HTTP/1.1
13-4	-	0/0/23	. 	2.13	1671	0	0.0	0.00	0.07 	127.0.0.1	deimos.fr	OPTIONS * HTTP/1.0
14-4	-	0/0/9	. 	1.03	1676	0	0.0	0.00	0.02 	127.0.0.1	deimos.fr	OPTIONS * HTTP/1.0
15-4	-	0/0/4	. 	1.35	1674	0	0.0	0.00	0.00 	127.0.0.1	deimos.fr	OPTIONS * HTTP/1.0
16-4	-	0/0/3	. 	0.00	1686	0	0.0	0.00	0.00 	127.0.0.1	deimos.fr	OPTIONS * HTTP/1.0
17-4	-	0/0/3	. 	0.00	1684	0	0.0	0.00	0.00 	127.0.0.1	deimos.fr	OPTIONS * HTTP/1.0
18-4	-	0/0/1	. 	0.00	1685	0	0.0	0.00	0.00 	127.0.0.1	deimos.fr	OPTIONS * HTTP/1.0
19-4	-	0/0/1	. 	0.00	1677	0	0.0	0.00	0.00 	127.0.0.1	deimos.fr	OPTIONS * HTTP/1.0
Srv	Child Server number - generation
PID	OS process ID
Acc	Number of accesses this connection / this child / this slot
M	Mode of operation
CPU	CPU usage, number of seconds
SS	Seconds since beginning of most recent request
Req	Milliseconds required to process most recent request
Conn	Kilobytes transferred this connection
Child	Megabytes transferred this child
Slot	Total megabytes transferred this slot
mod_fcgid status:
Total FastCGI processes: 0
SSL/TLS Session Cache Status:
cache type: SHMCB, shared memory: 512000 bytes, current sessions: 0
subcaches: 32, indexes per subcache: 133
index usage: 0%, cache usage: 0%
total sessions stored since starting: 0
total sessions expired since starting: 0
total (pre-expiry) sessions scrolled out of the cache: 0
total retrieves since starting: 0 hit, 0 miss
total removes since starting: 0 hit, 0 miss
```

It's normal to see several keep-alives managed by the Apache workers. If you see too many, it means they're keeping connections open for too long. For that, lower the connection reservation time with the KeepAliveTimeout option.

If you see several "." inactive, you should increase the MaxClients value to ensure enough free slots for new connections. This increases the responsiveness, but also increases memory usage.

## Apachetop

Apachetop displays accessed pages similar to the top command.

### Installation

On Debian, it's easy:

```bash
aptitude install apachetop
```

### Utilization

For use:

```bash
> apachetop -f /var/log/apache2/access.log
last hit: 00:00:00         atop runtime:  0 days, 00:01:55             13:41:33
All:            0 reqs (   0.0/sec)          0.0B (    0.0B/sec)       0.0B/req
2xx:       0 ( 0.0%) 3xx:       0 ( 0.0%) 4xx:     0 ( 0.0%) 5xx:     0 ( 0.0%)
R ( 30s):       0 reqs (   0.0/sec)          0.0B (    0.0B/sec)       0.0B/req
2xx:       0 ( 0.0%) 3xx:       0 ( 0.0%) 4xx:     0 ( 0.0%) 5xx:     0 ( 0.0%)
...
```

## Resources
- http://articles.slicehost.com/2010/3/26/enabling-and-using-apache-s-mod_status-on-debian
