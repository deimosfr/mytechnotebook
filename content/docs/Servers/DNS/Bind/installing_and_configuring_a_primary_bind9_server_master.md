---

weight: 999
url: "/Installation_et_configuration_d'un_serveur_Bind9_primaire_(Master)/"
title: "Installing and Configuring a Primary Bind9 Server (Master)"
description: "A comprehensive guide to installing and configuring a primary Bind9 DNS server, including security settings, zone configurations, and troubleshooting tips."
categories: ["Security", "Linux"]
date: "2013-05-07T09:06:00+02:00"
lastmod: "2013-05-07T09:06:00+02:00"
tags: ["DNS", "Bind9", "Server Configuration", "Security", "Linux"]
toc: true

---

## Introduction

[BIND (Berkeley Internet Name Domain)](https://www.isc.org/index.pl?/sw/bind/) is the most widely used DNS server on the Internet, especially on Unix-like systems. It is currently maintained by the Internet Systems Consortium.

A new version of BIND (BIND 9) was rewritten to solve some architectural issues in the initial code and to add support for DNSSEC (DNS Security Extensions).

## Installation

Installing Bind9 is quite simple:

```bash
aptitude install bind9 dnsutils
```

On OpenBSD, bind is installed by default.

## Configuration

### host.conf

Here's how to configure the `/etc/host.conf` file:

```bash
order hosts,bind
multi on
```

We are specifying here that requests made from the server should first check the hosts file and then Bind.

Note: This modification is not needed for OpenBSD.

### rndc.conf

#### Introduction to TSIG Keys

Transaction signatures ("TSIG") are a simpler form of DNS security. They use cryptographic hash functions to generate pseudo-signatures of DNS packets. The hash value is a combination of actual DNS data, timestamps to prevent replay attacks, and a shared secret between client and server. Since both entities involved in the DNS lookup must know the shared secret, TSIG signatures can only really be implemented in environments where systems are under common administrative control and where confidentiality of the shared secret can be absolutely guaranteed. In the case of ENUM, this means they can and should be used among ENUM level 0 name servers. For example, they can be used to validate zone transfers or dynamic update requests, with these functions being restricted to trusted clients because they know the shared secret.

#### Creating a TSIG Key

If you're on OpenBSD, we'll simplify things a bit, and to ensure this tutorial has the same paths everywhere, we'll create a symbolic link:

```bash
ln -s /var/named/etc /etc/bind
```

Let's generate a TSIG key:

```bash
cd /etc/bind
dnssec-keygen -a hmac-md5 -b 512 -n HOST simba
```

* Replace the number of bits with whatever you want and use your hostname.
* Replace **simba** with the name of your server where Bind is installed.

Once generated (this may take a few minutes), you will have 2 files:

* The key that needs to be moved to the bind folder.
* The rndc.conf file to be created.

Here are my 2 generated files:

* Ksimba.+157+18808.key
* Ksimba.+157+18808.private

Let's first move the key and assign it the correct permissions:

```bash
mv Ksimba.+157+18808.key /etc/bind/rndc.key
chmod 640 /etc/bind/rndc.key
```

For Debian:

```bash
chown root:bind /etc/bind/rndc.key
```

For OpenBSD:

```bash
chown root:named /etc/bind/rndc.key
```

If I display the content of the other file, I will see the key that will be used (the one that will be used to fill in the following files):

```bash
$ cat Ksimba.+157+18808.private
Private-key-format: v1.2
Algorithm: 157 (HMAC_MD5)
Key: a4fGtm0fB4zO+4KfqH/zNZ3nPq+ThM5yUCEE7AqzEVI=
Bits: AAA=
```

Then let's create the `/etc/bind/rndc.conf` file and insert the key:

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

### named.conf

```bash
// This is the primary configuration file for the BIND DNS server named.
//
// Please read /usr/share/doc/bind9/README.Debian.gz for information on the 
// structure of BIND configuration files in Debian, *BEFORE* you customize 
// this configuration file.
//
// If you are just adding zones, please do that in /etc/bind/named.conf.local

include "/etc/bind/named.conf.options";
include "/etc/bind/named.conf.local";
include "/etc/bind/named.conf.default-zones";
```

### named.conf.local

The acl section allows you to define reusable access lists in other sections of the configuration file. The following definition defines internal clients:

```bash
//
// Do any local configuration here
//

// Consider adding the 1918 zones here, if they are not used in your
// organization
//include "/etc/bind/zones.rfc1918";

// Acl definition

acl "zoneinterne" {
        // IP range authorized to make DNS requests
        192.168.0.0/24;
        10.8.0.0/24;
        127.0.0.1;
};

acl "srvsecondaires" {
        // My secondary server
        x.x.x.x;
        // Gandi (which offers a secondary DNS)
        217.70.177.40;
};
```

Let's define the logging part:

```bash
// Logs

logging {
        channel xfer-log {
                file "/var/log/bind.log";
                print-category yes;
                print-severity yes;
                print-time yes;
                severity info;
        };  
        category xfer-in { xfer-log; };
        category xfer-out { xfer-log; };
        category notify { xfer-log; };
};
```

Let's add some security to limit remote administration. If you don't want to authorize anything, leave the controls section empty:

```bash
// TSIG Security | RNDC Key

key "rndc-key" {
      algorithm hmac-md5;
      secret "a4fGtm0fB4zO+4KfqH/zNZ3nPq+ThM5yUCEE7AqzEVI=";
};

controls {
      inet 127.0.0.1 port 953
              allow { 127.0.0.1; "srvsecondaires"; } keys { "rndc-key"; };
};
```

### named.conf.default-zones

Ideally, build a zone/view file (to include in named.conf for each of your zones). But for simplicity, we'll leave everything in the same file here.

The view sections define server behaviors based on the IP address of the client sending the request, allowing DNS responses to be differentiated. We define two views:

* One corresponding to clients in the internal and DMZ zone: recursion needs to be re-enabled for these requests, and resolving all possible names (zone ".") needs to be allowed.
* Another corresponding to requests from outside (e.g. Internet). Only authorize requests for the zone where the DNS server has authority:

```bash
view "interne" {

    // These are the clients that see this view
    match-clients {
        zoneinterne;
    };

    // Recursion permited for zoneinterne ACL subnets
    recursion yes;
    allow-recursion {
        zoneinterne;
    };  

    // prime the server with knowledge of the root servers
    zone "." {
        type hint;
            file "/etc/bind/db.root";
        };

    // be authoritative for the localhost forward and reverse zones, and for
    // broadcast zones as per RFC 1912

    zone "deimos.fr" {
        type master;
        notify no;
        allow-update { none; };
        file "/etc/bind/db.deimos.fr.local";
    };

    zone "mavro.fr" {
        type master;
        notify no;
        allow-update { none; };
        file "/etc/bind/db.mavro.fr.local";
    };

   zone "localhost" {
        type master;
        notify no;
        allow-update { none; };
        file "/etc/bind/db.local";
   };

   zone "127.in-addr.arpa" {
        type master;
        notify no;
        allow-update { none; };
        file "/etc/bind/db.127";
   };

   zone "0.in-addr.arpa" {
        type master;
        notify no;
        allow-update { none; };
        file "/etc/bind/db.0";
   };

   zone "255.in-addr.arpa" {
        type master;
        notify no;
        allow-update { none; };
        file "/etc/bind/db.255";
   };

   zone "0.168.192.in-addr.arpa" {
        type master;
        notify no;
        allow-update { none; };
        file "/etc/bind/db.0.168.192.inv.local";
   };
};
```

Now let's move on to our external zone which will be visible to everyone from outside:

```bash
view "externe" {
        match-clients {
                any;
        };
        // Recursion not permited for World Wide Web
        recursion no;

        zone "deimos.fr" {
                type master;
                notify yes;
                allow-update { none; };
                allow-transfer { "rndc-key"; };
                file "/etc/bind/db.deimos.fr";
        };

        zone "mavro.fr" {
                type master;
                notify yes;
                allow-update { none; };
                allow-transfer { "rndc-key"; };
                file "/etc/bind/db.mavro.fr";
        };

        zone "localhost" {
                type master;
                notify yes;
                allow-update { none; };
                file "/etc/bind/db.local";
        };

        zone "127.in-addr.arpa" {
                type master;
                notify yes;
                allow-update { none; };
                file "/etc/bind/db.127";
        };

        zone "0.in-addr.arpa" {
                type master;
                notify yes;
                allow-update { none; };
                file "/etc/bind/db.0";
        };

        zone "255.in-addr.arpa" {
                type master;
                notify yes;
                allow-update { none; };
                file "/etc/bind/db.255";
       };

        zone "0.168.192.in-addr.arpa" {
                type master;
                notify yes;
                allow-update { none; };
                allow-transfer { "rndc-key"; };
                file "/etc/bind/db.0.168.192.inv";
        };
};
```

### named.conf.options

Here's the content:

```bash
options {
    directory "/var/cache/bind";

    // If there is a firewall between you and nameservers you want
    // to talk to, you may need to fix the firewall to allow multiple
    // ports to talk.  See http://www.kb.cert.org/vuls/id/800113

    // If your ISP provided one or more IP addresses for stable 
    // nameservers, you probably want to use them as forwarders.  
    // Uncomment the following block, and insert the addresses replacing 
    // the all-0's placeholder.

    // Force other DNS to answer
    //forwarders {
    //      212.27.32.176;
    //      212.27.32.177;
    //};

    //========================================================================
    // If BIND logs error messages about the root key being expired,
    // you will need to update your keys.  See https://www.isc.org/bind-keys
    //========================================================================
    // Avoid Bind Cache Poisoning
    dnssec-enable yes;
    dnssec-validation auto;

    allow-query {
             any;
    };  

    // Security version
    // Check with: dig -t txt -c chaos VERSION.BIND @<dns.server.com>
    version "Microsoft 2008 DNS Server";

    auth-nxdomain no;    # conform to RFC1035
    listen-on-v6 { any; };
};
```

You can add these types of options to improve security:

```bash
   allow-recursion { none; };      // disabling recursive queries (and thus DNS cache Poisoning)
   allow-transfer { none; };       // disabling zone transfer
   
   notify no;                      // disabling zone change notifications
   zone-statistics no;             // disabling statistics
   interface-interval 0;           // disabling search for new interfaces
   max-cache-size 20M;             // max cache size
```

In case you want to redirect domains (e.g., google.com --> a machine on your local network):

```bash
zone "google.com" {
   type forward;
   forwarders {
      192.168.0.17;      // My primary DNS 
      192.168.0.30;      // My secondary DNS
   };
   forward only;
};
```

You can also add this:

```bash
 zone "128.168.192.in-addr.arpa" {
    type forward;
    forwarders {
       192.168.128.99;
    };
    forward only;
 };
```

### db.deimos.fr

Here is the content of my file that will be available to everyone:

```bash
$TTL    604800
@       IN      SOA     simba.deimos.fr. root.deimos.fr.  (
```

Replace the first field with the FQDN of your machine, and the second corresponds to the admin's email (here root@deimos.fr). The email address is written in a special way, but it works.

```bash
        2010301201 ; Serial (date + incrementation)
```

Note: By convention, the serial is in the form YYYYMMDDNN (YYYY: year, MM: month, DD: day, NN: revision).

```bash
         7200       ; Refresh
         3600       ; Retry
         1209600    ; Expire
         604800     ; Negative Cache TTL
         )

                 NS      mufasa.deimos.fr.
                 A       88.162.130.192
                 NS      ns6.gandi.net.
                 NS      shenzi.deimos.fr.
                 MX      5 mufasa.deimos.fr.
                 MX      10 shenzi.deimos.fr.
                 TXT     "v=spf1 ip4:192.168.0.0/24 a mx ~all exp=getlost.deimos.fr"
 getlost         TXT     "You are not allowed to send a message from this domain"
 _deimos.fr      TXT     "t=y; o=-;"
 m1._deimos.fr   TXT     "g=; k=rsa; p=;...IWWiAyklt5FDmS2U7QIDAQAB..."
```

For the TXT part, I'll let you check out other articles that deal with SPF.

```bash
localhost       A       127.0.0.1
mufasa          A       82.232.191.145
shenzi          A       88.191.130.125
ns6             A       217.70.177.40

mail            CNAME   mufasa
jabber          CNAME   mufasa
server          CNAME   mufasa
serveur         CNAME   mufasa
sftp            CNAME   mufasa
www             CNAME   mufasa
blocnotesinfo   CNAME   mufasa
infos           CNAME   mufasa
webmail         CNAME   mufasa
nagios          CNAME   mufasa
```

All of these correspond to your A records, canonical names, etc.

### db.deimos.fr.local

This is for the internal zone, put whatever you want without restrictions...

```bash
$TTL    604800
@       IN      SOA     simba.deimos.fr. root.deimos.fr.  (
        2010301201 ; Serial (date + incrementation)
        7200       ; Refresh
        3600       ; Retry
        1209600    ; Expire
        604800     ; Negative Cache TTL
        )   

                NS      simba.deimos.fr.
                A       192.168.110.3
                MX      5 mufasa.deimos.fr.
                MX      10 shenzi.deimos.fr.

localhost       A       127.0.0.1
mufasa          A       192.168.0.254
shenzi          A       192.168.20.4
ns6             A       217.70.177.40

rafiki          A       192.168.110.1
simba           A       192.168.110.3
sarabi          A       192.168.99.10

www             CNAME   rafiki
www1            CNAME   rafiki
www2            CNAME   shenzi

blocnotesinfo   CNAME   www 
blog            CNAME   www 
webmail         CNAME   www 
nagios          CNAME   www 
git             CNAME   www 
gitweb          CNAME   www 
phpmyadmin      CNAME   www 

backuppc        CNAME   sarabi
```

### db.0.168.192.inv

If now we want reverse DNS:

```bash
$TTL    604800
@       IN      SOA     simba.deimos.fr. root.deimos.fr.     (
        2010301201 ; Serial (date + incrementation)
        7200       ; Refresh
        3600       ; Retry
        1209600    ; Expire
        604800     ; Negative Cache TTL
        )

                NS      simba.deimos.fr.
                A       88.162.130.192
                A       88.191.31.89
                NS      ns6.gandi.net.
                NS      shenzi.deimos.fr.
                MX      5 simba.deimos.fr.
                MX      10 shenzi.deimos.fr.
                TXT     "v=spf1 ip4:192.168.0.0/24 a mx ~all exp=getlost.deimos.fr"
getlost         TXT     "You are not allowed to send a message from this domain"

1               PTR     localhost
2               PTR     simba
3               PTR     ns6.gandi.net
4               PTR     shenzi.deimos.fr
```

The PTRs should go from smallest to largest based on priorities.

### db.0.168.192.inv.local

Now, the reverse DNS for local:

```bash
$TTL    604800
@       IN      SOA     simba.deimos.fr. root.deimos.fr.     (
        2010301201 ; Serial (date + incrementation)
        3600       ; Refresh
        7200       ; Retry
        1209600    ; Expire
        604800     ; Negative Cache TTL
        )

                NS      simba.deimos.fr.
                MX      5 simba.deimos.fr.

1               PTR     localhost
2               PTR     simba
3               PTR     serveur
4               PTR     earth
5               PTR     wind
6               PTR     water
10              PTR     kiss
11              PTR     imprimante
30              PTR     psp
31              PTR     pocket
32              PTR     ldap
```

## Verification

All that's left is to restart the service and check the logs:

```bash
/etc/init.d/bind9 restart
```

```bash
$ tail -100 /var/log/syslog 
Jun 10 22:44:20 deimos named[9641]: starting BIND 9.4.1 -u bind
Jun 10 22:44:20 deimos named[9641]: found 2 CPUs, using 2 worker threads
Jun 10 22:44:20 deimos named[9641]: loading configuration from '/etc/bind/named.conf'
Jun 10 22:44:20 deimos named[9641]: no IPv6 interfaces found
Jun 10 22:44:20 deimos named[9641]: listening on IPv4 interface eth0, 192.168.0.1#53
Jun 10 22:44:20 deimos named[9641]: listening on IPv4 interface lo, 127.0.0.1#53
Jun 10 22:44:20 deimos named[9641]: listening on IPv4 interface tun0, 10.8.0.1#53
Jun 10 22:44:20 deimos named[9641]: listening on IPv4 interface bond0, 192.168.0.2#53
Jun 10 22:44:20 deimos named[9641]: automatic empty zone: view interne: 254.169.IN-ADDR.ARPA
Jun 10 22:44:20 deimos named[9641]: automatic empty zone: view interne: 2.0.192.IN-ADDR.ARPA
Jun 10 22:44:20 deimos named[9641]: automatic empty zone: view interne: 255.255.255.255.IN-ADDR.ARPA
Jun 10 22:44:20 deimos named[9641]: automatic empty zone: view interne: 0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.IP6.ARPA
Jun 10 22:44:20 deimos named[9641]: automatic empty zone: view interne: 1.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.IP6.ARPA
Jun 10 22:44:20 deimos named[9641]: automatic empty zone: view interne: D.F.IP6.ARPA
Jun 10 22:44:20 deimos named[9641]: automatic empty zone: view interne: 8.E.F.IP6.ARPA
Jun 10 22:44:20 deimos named[9641]: automatic empty zone: view interne: 9.E.F.IP6.ARPA
Jun 10 22:44:20 deimos named[9641]: automatic empty zone: view interne: A.E.F.IP6.ARPA
Jun 10 22:44:20 deimos named[9641]: automatic empty zone: view interne: B.E.F.IP6.ARPA
Jun 10 22:44:20 deimos named[9641]: command channel listening on 127.0.0.1#953
Jun 10 22:44:20 deimos named[9641]: zone 0.in-addr.arpa/IN/interne: loaded serial 1
Jun 10 22:44:20 deimos named[9641]: zone 127.in-addr.arpa/IN/interne: loaded serial 1
Jun 10 22:44:20 deimos named[9641]: zone 0.168.192.in-addr.arpa/IN/interne: loaded serial 2006031801
Jun 10 22:44:20 deimos named[9641]: zone 255.in-addr.arpa/IN/interne: loaded serial 1
Jun 10 22:44:20 deimos named[9641]: zone mavro.fr/IN/interne: loaded serial 2006031801
Jun 10 22:44:20 deimos named[9641]: zone localhost/IN/interne: loaded serial 1
Jun 10 22:44:20 deimos named[9641]: zone deimos.fr/IN/interne: loaded serial 2006110401
Jun 10 22:44:20 deimos named[9641]: zone 0.in-addr.arpa/IN/externe: loaded serial 1
Jun 10 22:44:20 deimos named[9641]: zone 127.in-addr.arpa/IN/externe: loaded serial 1
Jun 10 22:44:20 deimos named[9641]: zone 0.168.192.in-addr.arpa/IN/externe: loaded serial 2006102701
Jun 10 22:44:20 deimos named[9641]: zone 255.in-addr.arpa/IN/externe: loaded serial 1
Jun 10 22:44:21 deimos named[9641]: zone mavro.fr/IN/externe: loaded serial 2007040103
Jun 10 22:44:21 deimos named[9641]: zone localhost/IN/externe: loaded serial 1
Jun 10 22:44:21 deimos named[9641]: zone deimos.fr/IN/externe: loaded serial 2007040101
Jun 10 22:44:21 deimos named[9641]: running
```

No errors here, so everything is good :-)

## FAQ

### What folders and permissions for Bind's chroot?

Here are the commands to execute:

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

### named: invalid command from 127.0.0.1: bad auth

How is it that all my conf files contain the same RNDC key, and yet I get this type of error?

The reason is simple: Bind must already be running. So a quick check with:

```bash
netstat -auntp 
```

And there we should realize that it's already running. So kill the corresponding PIDs and restart the Bind service:

```bash
pkill named
```

### Why can't I make records with "_"?

Bind versions from 9.3 now incorporate more precise control over domain name validity. They can no longer contain _ (0x5f in the [ASCII table](https://www.lookuptables.com/)), as stipulated by [RFC 1035](https://tools.ietf.org/html/rfc1035). This is however quite unfortunate for me because I have several domains containing _.

Error message produced:

```
Mar  6 07:48:08 dns3 named[25459]: pri/rags.ch.hosts:25: wisteria_lane.rags.ch: bad owner name (check-names)
Mar  6 07:48:08 dns3 named[25459]: zone rags.ch/IN: loading master file pri/rags.ch.hosts: bad owner name (check-names)
```

Here's a small patch that solves the problem:

```bash
name_with_underscore.patch:

--- lib/dns/name.c.orig 2006-03-06 17:44:30.000000000 +0100
+++ lib/dns/name.c      2006-03-06 17:45:07.000000000 +0100
@@ -261,7 +261,7 @@
        return (ISC_FALSE);
 }

-#define hyphenchar(c) ((c) == 0x2d)
+#define hyphenchar(c) ((c) == 0x2d || (c) == 0x5f)
 #define asterchar(c) ((c) == 0x2a)
 #define alphachar(c) (((c) >= 0x41 && (c) <= 0x5a) \
                      || ((c) >= 0x61 && (c) <= 0x7a))
```

### zone 'deimos.fr' allows updates by IP address, which is insecure

You may have log errors like this:

```
Oct  2 17:29:03 star1 named[5120]: zone 'deimos.fr' allows updates by IP address, which is insecure 
```

This simply means that your RNDC key is not being used. To use it with ACLs, just add your secondary servers to the "controls" section and only allow updates (at zone level) with the RNDC key:

```bash
...
controls {
      inet 127.0.0.1 port 953 
              allow { 127.0.0.1; secondaryinternaldns; } keys { "rndc-key"; };
};
...
       zone "deimos.fr" {
               type master;
               notify yes;
               allow-update { key "rndc-key"; };
               file "/etc/bind/db.deimos.fr";
       };
...
```

Now restart your DNS server and there are no more problems, exchanges are encrypted :-)

### too many timeouts resolving 'mycompany.com/MX' (in 'eu'?): disabling EDNS

This error can appear for name resolution that takes too long due to UDP packet size. This can be very annoying, especially for email reception which can take a few hours. The solution is to add this line:

```bash
...
edns-udp-size 1460;
...
```

Then restart the bind service. If the problem persists, try lowering the size (change from 1460 to 500 for example). There are ways to test all this using the dig command.  
Use it like this until you no longer have timeouts:

```bash
dig +norec +dnssec mycompany.com MX @my_dns_server
dig +dnssec +norec +ignore dnskey MX @my_dns_server
```

### zone deimos.fr/IN/internalview: journal rollforward failed: journal out of sync with zone

The server might have been shut down abruptly, or the zone file was manually edited without the zone being frozen. To solve this problem:

* Stop the bind server
* Delete the journal files (*.jnl in "/etc/bind")
* Restart the bind server

### How do I clear my cache?

To clear your bind cache:

```bash
rndc flush
```

## Resources
- [DNS Server on OpenBSD](/pdf/openbsd_dns_server.pdf)
- [Installing An Ubuntu DNS Server With BIND](/pdf/installing_an_ubuntu_dns_server_with_bind.pdf)
- http://fr.wikipedia.org/wiki/BIND
- http://fr.gentoo-wiki.com/HOWTO_Bind
- http://www.zytrax.com/books/dns/
- http://brocas.org/blog/post/2006/06/22/14-de-la-securite-d-une-architecture-dns-d-entreprise
- http://groups.google.com/group/comp.protocols.dns.bind/browse_thread/thread/cfa8c63ec6bd08d6?pli=1
