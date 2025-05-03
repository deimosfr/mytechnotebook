---
weight: 999
url: "/Fail2ban_\\:_mise_en_place_de_règles_automatisées_iptables_pour_contrer_les_attaques_par_bruteforce/"
title: "Fail2ban: Implementing automated iptables rules to counter bruteforce attacks"
description: "How to implement Fail2ban to automatically create iptables rules that block IP addresses attempting bruteforce attacks on your services."
categories: ["Security", "Network", "Linux"]
date: "2014-08-22T09:12:00+02:00"
lastmod: "2014-08-22T09:12:00+02:00"
tags: ["fail2ban", "iptables", "security", "bruteforce", "wordpress"]
toc: true
---

## Introduction

Fail2ban scans log files (e.g. /var/log/apache/error_log) and bans IPs that show malicious signs -- too many password failures, seeking for exploits, etc. Generally Fail2Ban is then used to update firewall rules to reject the IP addresses for a specified amount of time, although any arbitrary other action (e.g. sending an email) could also be configured. Out of the box Fail2Ban comes with filters for various services (apache, courier, ssh, etc).

Fail2Ban is able to reduce the rate of incorrect authentication attempts however it cannot eliminate the risk that weak authentication presents. Configure services to use only two factor or public/private authentication mechanisms if you really want to protect services.

## Installation

```bash
aptitude install fail2ban
```

## Configuration

You may want to add your own rules. Here are examples.

### Wordpress

I want to block bruteforce on my Wordpress installation. Unfortunately Wordpress does not return 403 errors when an authentication fails. So we have to:

#### Jail

Add this in your jail.conf to check access and error log files (`/etc/fail2ban/jail.conf`):

```bash
[wp-auth-errors]

enabled = true
port = http,https
filter = wp-auth-error
logpath = /var/log/nginx/*error*.log
bantime = 3600
maxretry = 6

[wp-auth-access]

enabled = true
port = http,https
filter = wp-auth-access
logpath = /var/log/nginx/*access*.log
bantime = 3600
maxretry = 6
```

#### Filters

Here is the filter for access. It's a regex to catch the IP address in the log file (`/etc/fail2ban/filter.d/wp-auth-access.conf`):

```bash
# WordPress brute force auth filter
#
# Block IPs trying to auth wp wordpress
#
[Definition]
failregex = ^<HOST> -.*"POST.*(wp-login|xmlrpc)\.php
ignoreregex =
```

And for error logs (`/etc/fail2ban/filter.d/wp-auth-error.conf`):

```bash
# WordPress brute force auth filter
#
# Block IPs trying to auth wp wordpress
#
[Definition]
failregex = ^.*client: <HOST>,.*"POST.*(wp-login|xmlrpc)\.php
ignoreregex =
```

### Validate filters and configuration

You can validate the configuration of your filters like this:

```bash
fail2ban-regex <logfile> <fail2ban rule to test>
```

## Usage

### Unban someone

This solution is to ask iptables to unban an IP. But Fail2ban won't be aware of that and will still think that the attacker is blocked if you do not use solution one, until the maximum blocking retention time is reached.

Get the current chains list:

```bash
> iptables -L | grep ^Chain
Chain INPUT (policy ACCEPT)
Chain FORWARD (policy ACCEPT)
Chain OUTPUT (policy ACCEPT)
Chain fail2ban-nginx-naxsi (2 references)
Chain fail2ban-ssh (1 references)
```

If you do not know on which chain your IP has been blocked, remove the grep command.

Then ask iptables to see the current blocked IPs on a specific chain:

```bash
> iptables -L fail2ban-nginx-naxsi -v -n --line-numbers
Chain fail2ban-nginx-naxsi (1 references)
num   pkts bytes target     prot opt in     out     source               destination         
1      315 75198 RETURN     all  --  *      *       0.0.0.0/0            0.0.0.0/0           
2       16  1704 DROP       all  --  *      *       222.2.5.210          0.0.0.0/0
```

Now I want to remove the second line:

```bash
iptables -D fail2ban-nginx-naxsi 2
```

To finish, inform fail2ban to unban someone:

```bash
fail2ban-client get nginx-naxsi actionunban 222.2.5.210
```

Modify nginx-naxsi with the name of the fail2ban jail name.

## Resources
- [Fail2ban Documentation](/pdf/fail2ban.pdf)
