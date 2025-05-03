---
weight: 999
url: "/Packet_Filter\\:_Lutter_contre_le_bruteforce/"
title: "Packet Filter: Fighting Against Brute Force Attacks"
description: "How to configure Packet Filter (PF) to protect against brute force attacks on services like SSH by automatically blacklisting suspicious IP addresses."
categories: ["FreeBSD", "Linux", "Network", "Security"]
date: "2010-04-16T07:18:00+02:00"
lastmod: "2010-04-16T07:18:00+02:00"
tags: ["PF", "Firewall", "Security", "SSH", "Network", "FreeBSD", "OpenBSD"]
toc: true
---

## Introduction

You've probably seen brute force connection attempts in your connection logs (sshd, httpd, ftpd, etc.). This is annoying, fills up your logs, and makes your server work harder than it needs to.

Fortunately, Daniel Hartmeier thought of you and added convenient options to his famous PacketFilter firewall, affectionately nicknamed PF. These options are 'max-src-conn-rate' and 'max-src-conn', which are used in combination with 'overload'. These options are available in PF starting with OpenBSD 3.7, FreeBSD 6.0, and NetBSD 2.0.

## PF Configuration

This is configured in the PF configuration file, `/etc/pf.conf`. I'll give an example for SSH, but the principle is the same for other ports.

Previously, to authorize SSH connections from outside, you would have a line that looked like this (with $external being the name of your external network interface):

```pf {linenos=table}
pass in quick on $external inet proto tcp from any to any port ssh flags S/SA keep state
```

Simply replace this line with:

```pf {linenos=table}
table <ssh-bruteforce> persist
block in quick from <ssh-bruteforce>
pass in quick on $external inet proto tcp from any to any port ssh flags S/SA keep state ( max-src-conn-rate 2/10, overload <ssh-bruteforce> flush global)
```

In order:

* We create a table that will store the attackers' IPs
* We block everything coming from these IPs
* We allow SSH connections if there are fewer than 2 connection attempts in 10 seconds
* Otherwise, we register the IP in the table and destroy all connections corresponding to that IP

Obviously, you can customize the frequency of connection attempts and also use 'max-src-conn' to limit the total number of connections from an IP.

That's it - enjoy the tranquility, and say goodbye to mindless attacks!

## Managing Blacklisted IPs

To display the list of blacklisted IPs:

```bash
pfctl -t bruteforce -T show
```

To remove a blacklisted IP or all IPs:

```bash
pfctl -t bruteforce -T delete @IP
pfctl -t bruteforce -T flush
```

## Adding a Whitelist

For those who wish to add a whitelist, here are the lines to add:

```pf {linenos=table}
table <whitelist> persist file "/etc/ssh/whitelist"
pass in on $ext_if proto tcp from <whitelist> to $ext_if port 22 flags S/SA keep state
```

Here, the `/etc/ssh/whitelist` file must be filled with the IPs to whitelist.

## Configuration Example

If this isn't clear enough, here's a configuration example:

```pf {linenos=table}
#       $OpenBSD: pf.conf,v 1.34 2007/02/24 19:30:59 millert Exp $
#
# See pf.conf(5) and /usr/share/pf for syntax and examples.
# Remember to set net.inet.ip.forwarding=1 and/or net.inet6.ip6.forwarding=1
# in /etc/sysctl.conf if packets are to be forwarded between interfaces.

ext_if="trunk0"

set skip on lo

scrub in all 

# Whitelist / Blacklist table
table <blacklist> persist
table <whitelist> persist file "/etc/ssh/whitelist"

# Block SSH bruteforce
pass in on $ext_if proto tcp from !<whitelist> to $ext_if port 22 \
        flags S/SA keep state \
        (max-src-conn-rate 3/60, \
         overload <blacklist> flush global)

# Allow Whitelist
pass in on $ext_if proto tcp from <whitelist> to $ext_if port 22 flags S/SA keep state

# Block the ssh bruteforce bastards
block drop in on $ext_if from <blacklist>
pass in on $ext_if from <whitelist>

# Allow all outbound traffic :
pass out inet proto tcp from $ext_if to any flags S/SA keep state
pass out inet proto { udp, icmp } from $ext_if to any keep state
```

Here, the last rule between whitelist and blacklist is whitelist. This is because the last matching rule takes precedence.
This allows us to connect even if we get blacklisted because we failed to connect after x attempts, as long as we're in the whitelist.

Don't forget to reload the configuration:

```bash
pfctl -f /etc/pf.conf
```

## References

http://www.openbsd.org/faq/pf/fr/filter.html  
http://wiki.gcu.info/doku.php?id=bsd:pf_et_bruteforce&s=ssh
