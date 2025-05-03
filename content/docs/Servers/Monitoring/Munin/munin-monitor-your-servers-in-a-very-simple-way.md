---
weight: 999
url: "/Munin_\\:_Surveiller_ses_serveurs_de_façons_très_simple/"
title: "Munin: Monitor your servers in a very simple way"
description: "A tutorial on installing and configuring Munin, a simple yet powerful system monitoring tool that provides a web interface and supports client/server architecture."
categories: ["Monitoring", "System Administration", "Linux"]
date: "2009-11-27T21:42:00+02:00"
lastmod: "2009-11-27T21:42:00+02:00"
tags: ["Munin", "Monitoring", "RRDTool", "Graphs", "Server"]
toc: true
---

## Introduction

[Munin](https://munin.projects.linpro.no/wiki/PluginCat) is a relatively unknown monitoring tool, but unlike Cacti, it is very very simple to use and install. It is also full of surprising features.

Here's an overview:

- Simplicity
- Web interface for consultation
- Client/server architecture
- Support for RRDTool for graph generation
- Many plugins available
- Alert sending to Nagios
- SNMP protocol support
- Automatic detection of services on the machine

For a demo: http://munin.ping.uio.no/

## Installation

To install Munin, there are 2 components:

- Munin: The server
- Munin-node: The client

The server should be installed on the main machine as it will contain a web graphical interface.

### Server

To install the server, it's quite simple:

```bash
apt-get install munin
```

Munin is now installed.

### Client

To install the client, run this command:

```bash
apt-get install munin-node
```

Next, we need to configure it.

## Configuration

You found the installation simple? Well, the configuration is the same!

### Server

Edit the "/etc/munin/munin.conf" file and adapt it to your configuration:

```bash
[fire.deimos.fr] # Enter the name of your first machine
   address 127.0.0.1 # If it's also the server, don't change this line
   use_node_name yes 
```

```bash
[burnin.deimos.fr] # Here the name of my client
   address 10.8.0.6 # Here the IP address of my client
   use_node_name yes
```

#### Accelerate data collection

To speed up data collection with munin-update, add to /etc/munin/munin.conf on the master:

```bash
# /etc/munin/munin.conf
fork yes
```

This way, munin will create a fork for each machine to query, rather than querying them one after another.

### Client

The client is now installed, edit the "/etc/munin/munin-node.conf" file:

```bash
allow ^127\.0\.0\.1$
```

Replace this address if the client is not installed on the server. If it's just the client that is installed on this machine, then replace the address with that of the server.

Then restart the client:

```bash
/etc/init.d/munin-node restart
```

## Add-ons

- All available services on the machine are detected by the munin-node-configure command:

```bash
munin-node-configure –suggest | grep yes
```

```
cpu | no | yes
df | no | yes
df_inode | no | yes
entropy | no | yes
exim_mailqueue | no | yes
exim_mailstats | no | yes
forks | no | yes
if_ | no | yes +eth0 +eth1
if_err_ | no | yes +eth0 +eth1
interrupts | no | yes
```

- The Debian package activates plugins for detected services by creating links in the /etc/munin/plugins directory:

```bash
ls -l /etc/munin/plugins
```

```
lrwxrwxrwx 1 root root 28 Nov 20 18:14 cpu -> /usr/share/munin/plugins/cpu
lrwxrwxrwx 1 root root 27 Nov 20 18:14 df -> /usr/share/munin/plugins/df
```

- You can disable a plugin by removing its symlink and enable it by creating a symlink:

```bash
ln -s /usr/share/munin/plugins/apache_volume /etc/munin/plugins/
```

### ZFS

This is a Solaris-only plugin (requires the Kstat module) that monitors I/O on ZFS pools:

```perl
#!/usr/perl5/bin/perl -w
# ZFS munin plugin for (Open)Solaris
# By Nico <nico@gcu.info>

use strict;
use Sun::Solaris::Kstat;

my $Kstat = Sun::Solaris::Kstat->new();
my $bytes_read = ${Kstat}->{unix}->{0}->{vopstats_zfs}->{read_bytes};
my $bytes_write = ${Kstat}->{unix}->{0}->{vopstats_zfs}->{write_bytes};

if($ARGV[0] && $ARGV[0] eq "config") {
	print "graph_title ZFS Read/Write bytes\n";
	print "graph_args --base 1024 -l 0\n";
	print "graph_category disk\n";
	print "read.label Bytes read\n";
	print "read.info Bytes read on the ZFS pools\n";
	print "write.label Bytes written\n";
	print "write.info Bytes written on the ZFS pools\n";
	print "read.type DERIVE\n";
	print "read.min 0\n";
	print "write.type DERIVE\n";
	print "write.min 0\n";
	exit 0;
}

print "read.value ".$bytes_read."\n";
print "write.value ".$bytes_write."\n";
```

## Resources
- [Documentation on Monitoring Multiple Systems With Munin](/pdf/monitoring_multiple_systems_with_munin.pdf)
- [Monitoring with Munin](/pdf/monitoring_with_munin.pdf)
- http://www.rottenbytes.info/?p=79
