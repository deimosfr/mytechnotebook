---
weight: 999
url: "/DHCP_\\:_Installation_et_configuration_d'un_serveur_DHCP/"
title: "DHCP: Installation and Configuration of a DHCP Server"
description: "Guide to installing and configuring DHCP servers on various operating systems including Debian, FreeBSD, OpenBSD, NetBSD, and Red Hat."
categories: ["Debian", "FreeBSD", "Red Hat"]
date: "2013-03-15T22:36:00+02:00"
lastmod: "2013-03-15T22:36:00+02:00"
tags: ["DHCP", "RFC 2131", "RFC 2132", "Networking", "Configuration"]
toc: true
---

## Introduction

**Dynamic Host Configuration Protocol** (DHCP) is a network protocol designed to automatically configure TCP/IP parameters for a station, particularly by automatically assigning an IP address and subnet mask. DHCP can also configure the default gateway address, DNS name servers, and NBNS name servers (known as WINS servers on Microsoft networks).

The initial design of Internet Protocol (IP) assumed pre-configuration of each computer connected to the network with the appropriate TCP/IP parameters - this is known as *static* addressing. On large or extended networks where changes occur frequently, static addressing creates a heavy maintenance burden and risks of errors. Additionally, assigned addresses cannot be used even when the computer that holds them is not in service. This is a typical problem for Internet Service Providers (ISPs), who generally have more customers than IP addresses available, but never have all customers connected simultaneously.

DHCP offers solutions to these issues:

- Only computers in service use an address from the address space
- Any parameter changes (gateway address, name servers) are applied to stations upon restart
- Modification of these parameters is centralized on DHCP servers

The protocol was first presented in October 1993 and is defined by RFC1531, modified and supplemented by [RFC 1534](https://tools.ietf.org/html/rfc1534), [RFC 2131](https://tools.ietf.org/html/rfc2131), and [RFC 2132](https://tools.ietf.org/html/rfc2132).

## Installation

### Debian

For installation, as usual, it's quite simple:

```bash
aptitude install dhcp3-server
```

### OpenBSD

On OpenBSD, nothing to install, just activate it. Edit the rc.conf file and insert the interfaces that should provide DHCP:

```bash
dhcpd_flags="sis1 vr0"
```

Then we'll create the reservations file:

```bash
touch /var/db/dhcpd.leases
```

### FreeBSD

Let's install a DHCP server:

```bash
pkg_add -vr isc-dhcp41-server
```

Then edit the rc.conf file to insert your default configuration with the interfaces on which the DHCP server will listen:

```bash
# Network services
dhcpd_enable="YES"    
dhcpd_flags="-q"    
dhcpd_conf="/usr/local/etc/dhcpd.conf"  
dhcpd_ifaces="vr1 vr2"                        
dhcpd_chuser_enable="YES"               # runs w/o privileges?
dhcpd_withuser="dhcpd"          # user name to run as
dhcpd_withgroup="dhcpd"         # group name to run as
dhcpd_chroot_enable="YES"               # runs chrooted?
dhcpd_devfs_enable="YES"                # use devfs if available?
dhcpd_rootdir="/var/db/dhcpd"   # directory to run in
```

### NetBSD

On NetBSD, the service already exists by default:

```bash
# DHCPd Server
dhcpd=YES
dhcpd_flags="-q vr1 vr2"
```

Specify the interfaces on which you want them to run (vr1 and vr2 here).

### Red Hat

To install a DHCP server on Red Hat:

```bash
yum install dhcp
```

Before modifying the configuration, let's copy a classic configuration:

```bash
cp -f /usr/share/doc/dhcp-*/dhcpd.conf.sample /etc/dhcp/dhcpd.conf
```

## Configuration

#### Example 1: Debian

Edit the /etc/dhcp3/dhcpd.conf file on Debian:

```bash
# dhcpd.conf
# Sample configuration file for ISC dhcpd
 
# option definitions common to all supported networks...
 
# Server authority on the network
authoritative;
 
# Refuse duplicate MAC addresses.
deny duplicates;
 
# Ignore DHCPDECLINE messages from clients, helps prevent
# successive address abandonment.
ignore declines;
 
# Various information is available to configure clients.
# See man dhcp-options for the list. In our case these options
# are the same for all our networks.
# option lpr-servers 192.168.0.7;
option netbios-name-servers 192.168.0.2;
option smtp-server 192.168.0.1;
# option pop-server 192.168.0.2;
 
# Indicate the address of your network or subnet with its mask.
# Parameters for the 192.168.0.0/240 network
subnet 192.168.0.0 netmask 255.255.255.0 {
 
### DNS Options ###
 
# Domain name for this zone.
option domain-name "deimos.fr";
 
# Name or addresses of DNS for all our networks.
option domain-name-servers 192.168.0.2, 212.27.32.176, 212.27.32.177;
 
# DNS update method:
ddns-update-style interim;
 
# updates allowed
ddns-updates on;
 
# here, we force updates by the DHCP server
ignore client-updates;
 
# we also force updates for fixed IPs
update-static-leases on;
 
# Information about your network:
option routers 192.168.0.138;
option subnet-mask 255.255.255.0;
option broadcast-address 192.168.0.255;
 
# static routes that clients will retrieve
option static-routes 
       172.16.0.0   192.168.0.138,
       10.0.0.0     192.168.0.138;
 
# Address ranges covered by DHCP.
Range 192.168.0.15 192.168.0.50;
default-lease-time 21600;
max-lease-time 43200;
 
### Reservations ###
# Computer1
host earth {
       hardware ethernet 24:13:D4:E9:15:56;
       fixed-address 192.168.0.3;
       option host-name "earth";
}
 
# Computer2
host ordiminix {
       hardware ethernet 22:0c:6e:34:80:56;
       fixed-address 192.168.0.4;
       option host-name "flower";
}
 
# include other configuration files
include "/etc/dhcp3/dhcpd.machintruc";
include "/etc/dhcp3/dhcpd.bidulechouette";
}
```

When you're done, restart the service and everything will work:

```bash
/etc/init.d/dhcp3-server restart
```

#### Example 2: Red Hat

Edit the /etc/dhcp/dhcpd.conf file to add the desired configuration. Here I have 2 declared ranges. Each range has its own interface:

```bash
# dhcpd.conf
#
# Sample configuration file for ISC dhcpd
#
 
# option definitions common to all supported networks...
option domain-name "deimos.fr";
option domain-name-servers ns1.deimos.fr, ns2.deimos.fr;
 
default-lease-time 600;
max-lease-time 7200;
 
# Use this to enble / disable dynamic dns updates globally.
ddns-update-style none;
allow booting;
allow bootp;
 
# If this DHCP server is the official DHCP server for the local
# network, the authoritative directive should be uncommented.
authoritative;
 
# Use this to send dhcp log messages to a different log file (you also
# have to hack syslog.conf to complete the redirection).
log-facility local7;
 
# No service will be given on this subnet, but declaring it helps the
# DHCP server to understand the network topology.
 
subnet 10.102.2.32 netmask 255.255.255.224 {
	option routers 10.102.2.63;
	option subnet-mask 255.255.255.224;
	option domain-name-servers 192.168.0.69;
	range 10.102.2.33 10.102.2.62;
	next-server 10.102.2.1;
	filename "pxelinux.0";
}
 
subnet 10.102.2.64 netmask 255.255.255.224 {
	option routers 10.102.2.65;
        option subnet-mask 255.255.255.224;
        option domain-name-servers 192.168.0.69;
        range 10.102.2.66 10.102.2.94;
        next-server 10.102.2.1;
        filename "pxelinux.0";
}
 
# This is a very basic subnet declaration.
 
#subnet 10.254.239.0 netmask 255.255.255.224 {
#  range 10.254.239.10 10.254.239.20;
#  option routers rtr-239-0-1.deimos.fr, rtr-239-0-2.deimos.fr;
#}
 
# This declaration allows BOOTP clients to get dynamic addresses,
# which we don't really recommend.
 
#subnet 10.254.239.32 netmask 255.255.255.224 {
#  range dynamic-bootp 10.254.239.40 10.254.239.60;
#  option broadcast-address 10.254.239.31;
#  option routers rtr-239-32-1.deimos.fr;
#}
 
# A slightly different configuration for an internal subnet.
#subnet 10.5.5.0 netmask 255.255.255.224 {
#  range 10.5.5.26 10.5.5.30;
#  option domain-name-servers ns1.internal.deimos.fr;
#  option domain-name "internal.deimos.fr";
#  option routers 10.5.5.1;
#  option broadcast-address 10.5.5.31;   
#  default-lease-time 600;
#  max-lease-time 7200;
#}
 
# Hosts which require special configuration options can be listed in
# host statements.   If no address is specified, the address will be
# allocated dynamically (if possible), but the host-specific information
# will still come from the host declaration.
 
#host passacaglia {
#  hardware ethernet 0:0:c0:5d:bd:95;
#  filename "vmunix.passacaglia";
#  server-name "toccata.fugue.com";
#}
 
# Fixed IP addresses can also be specified for hosts.   These addresses
# should not also be listed as being available for dynamic assignment.
# Hosts for which fixed IP addresses have been specified can boot using
# BOOTP or DHCP.   Hosts for which no fixed address is specified can only
# be booted with DHCP, unless there is an address range on the subnet
# to which a BOOTP client is connected which has the dynamic-bootp flag
# set.
 
# You can declare a class of clients and then do address allocation
# based on that.   The example below shows a case where all clients
# in a certain class get addresses on the 10.17.224/24 subnet, and all
# other clients get addresses on the 10.0.29/24 subnet.
 
#class "foo" {
#  match if substring (option vendor-class-identifier, 0, 4) = "SUNW";
#}
 
#shared-network 224-29 {
#  subnet 10.17.224.0 netmask 255.255.255.0 {
#    option routers rtr-224.deimos.fr; 
#  }
#  subnet 10.0.29.0 netmask 255.255.255.0 {
#    option routers rtr-29.deimos.fr;  
#  }
#  pool {
#    allow members of "foo";
#    range 10.17.224.10 10.17.224.250;   
#  }
#  pool {
#    deny members of "foo";
#    range 10.0.29.10 10.0.29.230;
#  }
#}
```

Then I'll declare the interfaces on which the dhcpd service should listen:

```bash
# Command line options here
DHCPDARGS="eth1 eth2";
```

As I mentioned above, I have one interface per range, so we'll add the appropriate routes:

```bash
ADDRESS1=10.102.2.32
NETMASK1=255.255.255.224
GATEWAY1=10.102.2.63
```

```bash
ADDRESS2=10.102.2.64
NETMASK2=255.255.255.224
GATEWAY2=10.102.2.94
```

Then I restart the service:

```bash
service restart dhcpd
```

#### Example 3: FreeBSD / NetBSD

On FreeBSD, the configuration is similar to other versions:

```bash
# dhcpd.conf
#
# Sample configuration file for ISC dhcpd
#
 
# option definitions common to all supported networks...
option domain-name "deimos.fr";
option domain-name-servers 8.8.8.8, 192.168.10.138;
 
default-lease-time 600;
max-lease-time 7200;
 
# Use this to enable / disable dynamic dns updates globally.
ddns-update-style none;
 
# If this DHCP server is the official DHCP server for the local
# network, the authoritative directive should be uncommented.
authoritative;
 
# Allow each client to have exactly one lease, and expire
# old leases if a new DHCPDISCOVER occurs
one-lease-per-client true;
 
# Tell the server to look up the host name in DNS
get-lease-hostnames true;
 
# Ping the IP address that is being offered to make sure it isn't
# configured on another node. This has some potential repercussions
# for clients that don't like delays.
ping-check true;
 
# Use this to send dhcp log messages to a different log file (you also
# have to hack syslog.conf to complete the redirection).
log-facility local7;
 
# No service will be given on this subnet, but declaring it helps the 
# DHCP server to understand the network topology.
 
#-----------------------------------------
# Subnet declaration
#-----------------------------------------
subnet 192.168.1.0 netmask 255.255.255.0 {
    range 192.168.1.100 192.168.1.199;
    option domain-name-servers 8.8.8.8;
    option domain-name "deimos.fr";
    option routers 192.168.1.254;
    option broadcast-address 192.168.1.255;
    default-lease-time 600;
    max-lease-time 7200;
}
 
#-----------------------------------------
# Hostname declaration
#-----------------------------------------
#host ipad_deimos {
#    hardware ethernet 00:...;
#    fixed-address 192.168.1.90;
#}
 
#subnet 10.152.187.0 netmask 255.255.255.0 {
#}
 
# This is a very basic subnet declaration.
 
#subnet 10.254.239.0 netmask 255.255.255.224 {
#  range 10.254.239.10 10.254.239.20;
#  option routers rtr-239-0-1.example.org, rtr-239-0-2.example.org;
#}
 
# This declaration allows BOOTP clients to get dynamic addresses,
# which we don't really recommend.
 
#subnet 10.254.239.32 netmask 255.255.255.224 {
#  range dynamic-bootp 10.254.239.40 10.254.239.60;
#  option broadcast-address 10.254.239.31;
#  option routers rtr-239-32-1.example.org;
#}
 
# A slightly different configuration for an internal subnet.
#subnet 10.5.5.0 netmask 255.255.255.224 {
#  range 10.5.5.26 10.5.5.30;
#  option domain-name-servers ns1.internal.example.org;
#  option domain-name "internal.example.org";
#  option routers 10.5.5.1;
#  option broadcast-address 10.5.5.31;
#  default-lease-time 600;
#  max-lease-time 7200;
#}
 
# Hosts which require special configuration options can be listed in
# host statements.   If no address is specified, the address will be
# allocated dynamically (if possible), but the host-specific information
# will still come from the host declaration.
 
#host passacaglia {
#  hardware ethernet 0:0:c0:5d:bd:95;
#  filename "vmunix.passacaglia";
#  server-name "toccata.fugue.com";
#}
 
# Fixed IP addresses can also be specified for hosts.   These addresses
# should not also be listed as being available for dynamic assignment.
# Hosts for which fixed IP addresses have been specified can boot using
# BOOTP or DHCP.   Hosts for which no fixed address is specified can only
# be booted with DHCP, unless there is an address range on the subnet
# to which a BOOTP client is connected which has the dynamic-bootp flag
# set.
#host fantasia {
#  hardware ethernet 08:00:07:26:c0:a5;
#  fixed-address fantasia.fugue.com;
#}
 
# You can declare a class of clients and then do address allocation
# based on that.   The example below shows a case where all clients
# in a certain class get addresses on the 10.17.224/24 subnet, and all
# other clients get addresses on the 10.0.29/24 subnet.
 
#class "foo" {
#  match if substring (option vendor-class-identifier, 0, 4) = "SUNW";
#}
 
#shared-network 224-29 {
#  subnet 10.17.224.0 netmask 255.255.255.0 {
#    option routers rtr-224.example.org;
#  }
#  subnet 10.0.29.0 netmask 255.255.255.0 {
#    option routers rtr-29.example.org;
#  }
#  pool {
#    allow members of "foo";
#    range 10.17.224.10 10.17.224.250;
#  }
#  pool {
#    deny members of "foo";
#    range 10.0.29.10 10.0.29.230;
#  }
#}
```

Since I've made reservations with includes and my DHCP service is chrooted, I'll create a small directory structure that will simplify my life:

```bash
mkdir /var/db/dhcpd/etc/dhcpd.d
cd /usr/local/etc
ln -s /var/db/dhcpd/etc/dhcpd.d .
```

In there I'll have my reservation files.

Then to start the service:

```bash
/usr/local/etc/rc.d/isc-dhcpd start
```

## Resources
- [How To Set Up DHCP Failover](/pdf/how_to_set_up_dhcp_failover.pdf)
