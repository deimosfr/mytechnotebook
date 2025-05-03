---
weight: 999
url: "/Collectd_\\:_Installation_et_configuration_de_Collectd/"
title: "Collectd: Installation and Configuration"
description: "Guide to installing and configuring Collectd, a system statistics collection daemon, on various platforms including Debian and Solaris"
categories: ["Servers", "Monitoring", "Linux", "Solaris"]
date: "2010-10-06T08:20:00+02:00"
lastmod: "2010-10-06T08:20:00+02:00"
tags: ["collectd", "monitoring", "statistics", "rrd", "performance"]
toc: true
---

## Introduction

Collectd gathers statistics about the system it is running on and stores this information. Those statistics can then be used to find current performance bottlenecks (i.e. performance analysis) and predict future system load (i.e. capacity planning). Or if you just want pretty graphs of your private server and are fed up with some homegrown solution you're at the right place, too ;).

After this installation, you'll likely want access through a web interface. Feel free to continue with [these documentations](https://wiki.deimos.fr/Serveurs.html#Collectd).

## Installation

Whether it's on the server or client, it's quite convenient, it's the same package that needs to be installed.

### Debian

On Debian, it's always as simple as:

```bash
aptitude install collectd
```

### Client

#### Debian

On Debian, it's always as simple as:

```bash
aptitude install collectd
```

#### Solaris

On Solaris, it's available on SunFreeware:

```bash
wget "http://collectd.org/files/collectd-4.6.2-0-SunOS-5.10-x86_64.pkg"
pkgadd -d collectd-4.6.2-0-SunOS-5.10-x86_64.pkg
```

## Configuration

### Server

For the server, we're going to make sure these modules are activated (uncommented):

```bash
...
LoadPlugin "logfile"
LoadPlugin "network"
LoadPlugin "rrdtool"
...
```

Then we're going to add the IP on which the Collectd server should listen:

```bash
...
# Server
<Plugin "network">
  Listen "192.168.0.25"
</Plugin>
...
```

Unless the Collectd server is a beast of a machine, it's preferable to cache the information before writing it to disk. This will reduce I/O and CPU usage:

```bash
...
<Plugin rrdtool>
    DataDir "/var/lib/collectd/rrd"
    CacheTimeout 485 
    CacheFlush 1800
</Plugin>
...
```

Now restart the collectd server.

### Client

On the client side, choose a few modules to activate:

```bash
...
LoadPlugin "logfile"
LoadPlugin "network"
LoadPlugin "cpu"
LoadPlugin "memory"
...
```

And finally enter the IP address of the server:

```bash
...
# Client
<Plugin "network">
  Server "192.168.0.42"
</Plugin>
...
```

Now restart the collectd client.

If you want to activate other services, just uncomment the ones you're interested in.

## Adding Modules

You can easily add modules by compiling them. Depending on the OS you're on, it's more or less easy.

### Compiling Modules on Solaris

After struggling for a few hours, here's the process to compile collectd with its modules in 32 and 64 bits.

* First, we're going to download and decompress the sources:

```bash
wget http://collectd.org/files/collectd-4.10.0.tar.gz
gtar -xzvf collectd-4.10.0.tar.gz
```

* Next, we'll install some dependencies. You need to download the [Solaris companion CD](https://www.sun.com/software/solaris/freeware/) and install these packages:

```bash
cd /cdrom/companioncd/Packages
pkgadd -d `pwd` SFWaconf
pkgadd -d `pwd` SFWamake
pkgadd -d `pwd` SFWcoreu
pkgadd -d `pwd` SFWcurl
pkgadd -d `pwd` SFWctags
pkgadd -d `pwd` SFWdiffu
pkgadd -d `pwd` SFWgawk
pkgadd -d `pwd` SFWgfind
pkgadd -d `pwd` SFWltool
pkgadd -d `pwd` SFWncur
pkgadd -d `pwd` SFWscrn
```

* Then you need to install [Oracle Solaris Studio Express](https://developers.sun.com/sunstudio/downloads/express/).

You can (for those who want) use GCC:

```bash
pkg-get install libiconv gcc3 gcc3core gmake
```

But I wasn't able to compile 64-bit libraries even following the [docs and FAQ](https://collectd.org/faq.shtml). That's why I recommend Solaris Studio Express.

* Then start the 64-bit compilation with standard options:

```bash
cd collectd-4.10.0
export PATH=$PATH:/opt/solstudioex1006/bin:/opt/SUNWspro/bin:/usr/ccs/bin
./configure CFLAGS="-m64 -mt -D_POSIX_PTHREAD_SEMANTICS"
make
```

This 'CFLAGS="-m64 -mt -D_POSIX_PTHREAD_SEMANTICS"' indicates that we want to use 64-bit. Remove it to do 32-bit.

If you want to compile only a specific plugin, you can do it like this:

```bash
./configure CFLAGS="-m64 -mt -D_POSIX_PTHREAD_SEMANTICS" --disable-all-plugins --enable-zfs_arc
```

Then you just need to do the famous 'make install' when you're ready.

## Creating the Solaris SMF

All necessary files are in the sources (contrib/solaris-smf). The cleanest way is to move the collectd.xml file to the same place as the other xmls:

```bash
mv collectd.xml /var/svc/manifest/application
```

By default the xml looks for the SMF in /lib/svc/method, so we move the script to the right place:

```bash
mv collectd /lib/svc/method
```

Then we register the xml:

```bash
svccfg import /var/svc/manifest/application/collectd.xml
```

## Resources
- http://collectd.org/wiki/index.php/First_steps
- http://collectd.org/wiki/index.php/Networking_introduction
