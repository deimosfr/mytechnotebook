---
weight: 999
url: "/Mise_en_place_d'un_DHCP_Failover/"
title: "Setting up DHCP Failover"
description: "A guide on how to configure DHCP failover between two servers using ISC DHCP server version 3 on OpenBSD systems."
categories: ["Linux", "Network", "Backup"]
date: "2007-11-06T18:42:00+02:00"
lastmod: "2007-11-06T18:42:00+02:00"
tags: ["DHCP", "failover", "OpenBSD", "ISC", "Network", "Servers", "Special pages"]
toc: true
---

## Introduction

I've been running a DHCP server on my home network for eons now, and today I decided I'd move it on to my OpenBSD firewall cluster. It probably really shouldn't be there but I already run a handful of other internal services there, like DNS, and NTP.

I assume you already have a working [dhcpd configuration](https://www.openbsd.org/faq/faq6.html#DHCP) for your single server, if you don't, then you can get a basic DHCP configuration from the OpenBSD FAQ.

## Installation

You need ISC's DHCP server (at least version 3), to do DHCP failover. As of OpenBSD 4.1, they ship version 2 by default. You can get version 3 out of the packages tree by installing isc-dhcp-server-3.0.4p0.tgz. It installs itself into `/usr/local`, so if you want to view the man pages, you have to do something like:

```bash
export MANPATH=/usr/share/man:/usr/local/man
```

otherwise, you won't see the new man pages for the config files.

## Configuration

Once you have it installed, you need to get it configured to run at startup. I did it by adding the following lines to my `/etc/rc.conf.local` file:

```bash
# turn on dhcpd3
dhcpd3=YES
dhcpd3_flags="-pf /var/run/dhcpd.pid"
```

I then added some stuff to my `/etc/rc.local` file:

```bash
# ISC dhcpd3 with failover configured
if [ X"${dhcpd3}" == X"YES" -a -x /usr/local/sbin/dhcpd -a -f /etc/dhcpd.conf ]; then
        touch /var/db/dhcpd.leases
        if [ -f /etc/dhcpd.interfaces ]; then
                dhcpd_ifs=`stripcom /etc/dhcpd.interfaces`
        fi
        echo -n ' dhcpd'
        /usr/local/sbin/dhcpd ${dhcpd3_flags} ${dhcpd_ifs}
fi
```

This will get the new version of DHCP started at boot time. You ought to remember to disable the other dhcpd by putting this line in your `rc.conf.local`:

```bash
dhcpd_flags=NO          # for normal use: ""
```

Now that you are already to start it up, we need to get the `/etc/dhcpd.conf` file ready. You probably already have one configured and working. If so then just do something like:

```bash
# mv /etc/dhcpd.conf /etc/dhcpd.master
```

and then create a new dhcpd.conf file for your primary node that looks like this:

```bash
#
# dhcpd configuration
#
 
# failover definition
failover peer "dhcp-failover" {
        primary; # declare ourselves primary
        address 192.168.13.6;
        port 520;
        peer address 192.168.13.7;
        peer port 520;
        max-response-delay 10;
        max-unacked-updates 10;
        load balance max seconds 3;
        mclt 1800;
        split 128;
}
 
# include the rest.  This allows us to copy dhcpd.master
# between the two machines safely
include "/etc/dhcpd.master";
```

The method to my madness is simple. The contents of `/etc/dhcpd.master` can be exactly replicated between you two dhcp servers. This is where you will have all your subnets, ranges, mac addresses, etc.etc. Use your favorite method to keep them synched. The contents of `/etc/dhcpd.conf` are different on the primary dhcp server and the secondary. You obviously wouldn't want to be copying them all over the place.

Some comments on the new dhcpd.conf file. The "dhcp-failover" string in the

```
failover peer "dhcp-failover" {
```

line can be whatever you want, but we're going to use it in several other places, and it has to be the same in all of those places. You would of course replace the appropriate address and peer address IP addresses with the ones of the two servers you will be balancing.

```bash
The /etc/dhcpd.conf file on the secondary server looks like this:
 
#
# dhcpd configuration
#
 
# failover definition
failover peer "dhcp-failover" {
        secondary; # declare ourselves secondary
        address 192.168.13.7;
        port 520;
        peer address 192.168.13.6;
        peer port 520;
        max-response-delay 10;
        max-unacked-updates 10;
        load balance max seconds 3;
}
 
# include the rest.  This allows us to copy dhcpd.master
# between the two machines safely
include "/etc/dhcpd.master";
```

Notice the changes from the primary config file. The addresses and peer addresses are swapped, and there a couple of missing config lines, that must not be present.

The final step is to modify our `/etc/dhcpd.master` file so that it knows that it should be failing over. Here is a small snippet from mine:

```bash
subnet 192.168.13.0 netmask 255.255.255.0 {
  option routers 192.168.13.1;
  option broadcast-address 192.168.13.255;
  pool {
    failover peer "dhcp-failover";
    deny dynamic bootp clients;
    range 192.168.13.32 192.168.13.47;
  }
}
```

The only new thing here is the failover peer line. The string there needs to be the same one that we used in our `/etc/dhcpd.conf` file.

That's it. Now to test it out.

## Verification

You can fire up the server and prevent if from forking and logging to standard out by doing something like:

```bash
# /usr/local/sbin/dhcpd -pf /var/run/dhcpd.pid -d -f xl0
```

You would of course replace "xl0" with the interface on your machine you want the server to listen on. Look for error messages, etc. If things are going right, you should see something like this:

```bash
 Internet Systems Consortium DHCP Server V3.0.4
Copyright 2004-2006 Internet Systems Consortium.
All rights reserved.
For info, please visit http://www.isc.org/sw/dhcp/
Wrote 0 deleted host decls to leases file.
Wrote 0 new dynamic host decls to leases file.
Wrote 8 leases to leases file.
Multiple interfaces match the same subnet: xl0 carp0
Multiple interfaces match the same shared network: xl0 carp0
Multiple interfaces match the same subnet: xl0 carp2
Multiple interfaces match the same shared network: xl0 carp2
Multiple interfaces match the same subnet: xl0 carp3
Multiple interfaces match the same shared network: xl0 carp3
Listening on BPF/xl0/00:01:03:d6:82:a1/192.168.13/24
Sending on   BPF/xl0/00:01:03:d6:82:a1/192.168.13/24
Sending on   Socket/fallback/fallback-net
failover peer dhcp-failover: I move from normal to startup
failover peer dhcp-failover: peer moves from normal to communications-interrupted
failover peer dhcp-failover: I move from startup to normal
failover peer dhcp-failover: peer moves from communications-interrupted to normal
pool 80f93200 192.168.15/24 total 16  free 16  backup 0  lts -8
pool 80f93100 192.168.13/24 total 16  free 8  backup 7  lts 0
pool 80f93200 192.168.15/24  total 16  free 16  backup 0  lts 8
```

Now you can try to get a client on your network to request an address, and you should see it happen. If that all works right, then you should try rebooting your machines, making sure that everything comes up properly on startup. You can now experiment with taking one or the other server down, and you should still be able to DHCP properly.
