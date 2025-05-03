---
weight: 999
url: "/Installation_et_configuration_d'un_serveur_Bind9_secondaire_(Slave)/"
title: "Installation and Configuration of a Bind9 Secondary (Slave) Server"
description: "Guide for setting up a secondary Bind9 DNS server including installation, configuration, and troubleshooting tips."
categories: ["Linux", "Security"]
date: "2009-10-04T19:25:00+02:00"
lastmod: "2009-10-04T19:25:00+02:00"
tags: ["DNS", "Bind9", "Server", "Configuration", "Security"]
toc: true
---

## Introduction

When you want to set up your own DNS server, you must have a secondary server. If you don't have other machines, you can use Gandi, otherwise follow this guide.

**Don't forget to declare the secondary server on the primary server.**

## Installation

On the secondary server:

```bash
apt-get install bind9
```

For OpenBSD, nothing to do, it's already installed by default.

## Configuration

Before beginning, **declare your secondary servers in your domain name records (update NS records + ACL in named.conf)**, otherwise notifications won't work.

### Configuration of Permissions

If you want a chrooted environment, proceed as follows, otherwise skip this step:

```bash
mkdir -p /var/lib/named/etc
mkdir /var/lib/named/dev
mkdir -p /var/lib/named/var/cache/bind
mkdir -p /var/lib/named/var/run/bind/run
mv /etc/bind /var/lib/named/etc
ln -s /var/lib/named/etc/bind /etc/bind
mknod /var/lib/named/dev/null c 1 3
mknod /var/lib/named/dev/random c 1 8
chmod 666 /var/lib/named/dev/null /var/lib/named/dev/random
chown -R bind:bind /var/lib/named/var/*
chown -R bind:bind /var/lib/named/etc/bind
```

Then edit the `/etc/default/bind9` file to tell it to chroot at service startup:

```bash
OPTIONS="-u bind -t /var/lib/named"
...
```

### host.conf

Here's how to configure the `/etc/host.conf` file:

```bash
order hosts,bind
multi on
```

We are saying here to first look in the hosts file and then in Bind when queries are made from the server.

### named.conf

This file must contain the same information for RNDC and external zone (to replicate):

```bash
include "/etc/bind/named.conf.options";

// be authoritative for the localhost forward and reverse zones, and for
// broadcast zones as per RFC 1912

// TSIG Security | RNDC Key

key "rndc-key" {
      algorithm hmac-md5;
      secret "a4fGtm0fB4zO+4KfqH/zNZ3nPq+ThM5yUCEE7AqzEVI=";
};

controls {
      inet 127.0.0.1 port 953
              allow { 127.0.0.1; } keys { "rndc-key"; };
};

zone "127.in-addr.arpa" {
        type master;
        notify no;
        file "/etc/bind/db.127";
};

// Starting secondary transfert zones
view "externe" {
        match-clients {
                any;
        };
        // Recursion not permited for World Wide Web
        recursion no;

        zone "deimos.fr" {
                type slave;
                file "/etc/bind/db.deimos.fr";
                notify yes;
                masters {
                        88.162.130.192;
                };
        };

        zone "mavro.fr" {
                type slave;
                file "/etc/bind/db.mavro.fr";
                notify yes;
                masters {
                      88.162.130.192;
                };
        };

        zone "0.168.192.in-addr.arpa" {
                type slave;
                file "/etc/bind/db.0.168.192.inv.local";
                notify yes;
                masters {
                        88.162.130.192;
                };
        };
};
```

### named.conf.options

Same file as for the primary server:

```bash
options {
        directory "/var/cache/bind";
        pid-file "/var/run/bind/run/named.pid";

        // If there is a firewall between you and nameservers you want
        // to talk to, you might need to uncomment the query-source
        // directive below.  Previous versions of BIND always asked
        // questions using port 53, but BIND 8.1 and later use an unprivileged
        // port by default.

        query-source address * port 53;

        // If your ISP provided one or more IP addresses for stable 
        // nameservers, you probably want to use them as forwarders.  
        // Uncomment the following block, and insert the addresses replacing 
        // the all-0's placeholder.

        // forwarders {
        //      0.0.0.0;
        // };

        // For dial up connections
        dialup yes;

        allow-query {
                any;
        };  

        // Security version
        version "Microsoft 2000 DNS Server";

        auth-nxdomain no;    # conform to RFC1035
};
```

### rndc.conf

We're using exactly the same file again as for the primary server:

```bash
 key "rndc-key" {
         algorithm hmac-md5;
         secret "a4fGtm0fB4zO+4KfqH/zNZ3nPq+ThM5yUCEE7AqzEVI=";
 };

 options {
         default-key "rndc-key";
         default-server 127.0.0.1;
         default-port 953;
 };
```

Don't forget to copy the **rndc.key** file from the primary server to the **/etc/bind** directory and assign it the correct permissions.

## Verification

All that's left is to restart the bind server and the domain name configuration files will automatically appear in /etc/bind:

```bash
/etc/init.d/bind9 restart
```

## FAQ

### No files appear and the logs say there are no permissions

The solution is simple:

```bash
chown -Rf root:bind /etc/bind
```

### Why do I have this message in my logs: refused notify from non-master

Here's the type of message you might encounter when you're not lucky:

```bash
Oct  2 16:43:50 tasmania named[7978]: zone deimos.local/IN/internalview: refused notify from non-master: 192.168.0.27#37097
```

This is due to a Bind update (9.3). You simply need to authorize the server to notify itself in your named.conf by adding allow-notify:

```bash
...
        zone "deimos.local" {
               type slave;
               file "/etc/bind/db.deimos.local";
               notify yes;
               masters { 192.168.0.69; };
               allow-notify { 192.168.0.27; };
       };
...
```

### refresh: failure trying master 53 (source 0.0.0.0#0): operation canceled

This is certainly due to a firewalling problem, make sure your port 53 is open for both TCP and UDP on both sides.

### could not open entropy source /dev/random: file not found

If you encounter this problem, you are probably working in a vserver. To work around the issue, add a bcapabilities:

```bash
echo CAP_MKNOD >> /etc/vservers/ed/bcapabilities
```

This will avoid problems like these in the logs:

```bash
24530 Oct  4 21:12:27 ed named[8642]: could not open entropy source /dev/random: file not found
24531 Oct  4 21:12:27 ed named[8642]: using pre-chroot entropy source /dev/random
```

## Resources
- http://fr.wikibooks.org/wiki/Linux_VServer
