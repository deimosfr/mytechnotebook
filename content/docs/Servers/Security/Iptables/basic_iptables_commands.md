---
weight: 999
url: "/Les_commandes_de_bases_d\\'Iptables/"
title: "Basic IPTables Commands"
description: "Learn the basic IPTables commands for Linux firewalls - including chain manipulation, rule management, and practical examples."
categories: ["Linux", "Network", "Security"]
date: "2013-05-06T15:42:00+02:00"
lastmod: "2013-05-06T15:42:00+02:00"
tags: ["iptables", "firewall", "security", "networking", "linux"]
toc: true
---

## Basic Commands

There are several operations you can do with iptables. You start with **three default chains INPUT, OUTPUT and FORWARD that you cannot delete**. Let's look at the operations to administer chains:

Creating a new chain:

```bash
iptables -N chain
```

Delete an empty chain:

```bash
iptables -X chain
```

Change the default rule for a starting chain:

```bash
iptables -P chain status
```

Example:

```bash
iptables -P INPUT DROP
```

or

```bash
iptables --policy FORWARD DROP
```

List the rules in a chain:

```bash
iptables -L chain
```

Remove rules from a chain:

```bash
iptables -F chain
```

or

```bash
iptables --flush chain
```

Empty rules from another table (e.g., NAT):

```bash
iptables --table nat --flush chain
```

Reset the bit and packet counters of a chain:

```bash
iptables -Z chain
```

## Manipulating Rules in a Chain

Add a new rule to the chain:

```bash
iptables -A
```

Insert a new rule at a position in the chain:

```bash
iptables -I
```

Replace a rule at a position in the chain:

```bash
iptables -R
```

Delete a rule at a position in the chain:

```bash
iptables -D
```

Delete the first matching rule in a chain:

```bash
iptables -D
```

## Displaying Your IPTables Configuration

Display the entire filter table:

```bash
iptables –L –v
```

Display only the NAT table:

```bash
iptables –t nat –L –v
```

## Usage Examples

To allow packets on the telnet port coming from a local network:

```bash
iptables --append INPUT --protocol tcp --destination-port telnet --source 192.168.13.0/24 --jump ACCEPT
```

To ignore other incoming packets on the telnet port:

```bash
iptables -A INPUT -p tcp --dport telnet -j DROP
```

To reject incoming packets on port 3128, often used by proxies, then add a comment:

```bash
iptables -A INPUT -p tcp --dport 3128 -j REJECT --reject-with tcp-reset -m comment --comment "Rejecting default proxy port"
```

To perform automatic NAT for all packets leaving through the ppp0 interface (often representing the internet connection):

```bash
iptables -t nat -A POSTROUTING -o ppp0 -j MASQUERADE
```

Disable all rules without disconnecting:

```bash
iptables -F && iptables -X && iptables -P INPUT ACCEPT && iptables -OUTPUT ACCEPT
```

## Resources

- [Firewalling under Linux](/pdf/firewalling_sous_linux.pdf)
- [Implementation of an Internet Gateway](/pdf/mise_en_oeuvre_d'une_passerelle_insternet.pdf)
- [IPtables from top to bottom](/pdf/firewall-iptables.pdf)
