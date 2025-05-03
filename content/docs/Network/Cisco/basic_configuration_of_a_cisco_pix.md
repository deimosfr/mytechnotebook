---
weight: 999
url: "/Configuration_de_base_d'un_Cisco_Pix/"
title: "Basic Configuration of a Cisco PIX"
description: "Learn the basics of Cisco PIX firewall configuration including VPN setup, user administration, and network access rules"
categories: ["Network", "Security", "Cisco"]
date: "2006-11-08T10:33:00+01:00"
lastmod: "2006-11-08T10:33:00+01:00"
tags: ["Cisco", "PIX", "firewall", "VPN", "network security"]
toc: true
---

## Cisco Basics

### Introduction

In graphical mode it's not too complicated, but in command line - oh boy - the PIX is challenging to configure!

It's Cisco, so it uses proprietary commands but with a touch of Unix. For example, the grep command works :-)

Let's get started. Connect and enter your password.

### The Basics

* Switch to **enable** mode to get admin privileges:

```bash
en
```

* Then a **show running-config** to see the current configuration:

```bash
sh run
```

* Once you've viewed the configuration and decided to add something, do a **configure terminal**:

```bash
conf t
```

* To exit the current mode:

```bash
Ctrl+z or exit
```

* Are you sure? Then save to memory (**write memory**):

```bash
wr m
```

### Adding and Removing Commands

* To add a NAT rule for example, copy the existing lines, then copy/paste:

```bash
static (dmz,outside) tcp 193.252.19.3 2099 SRV-FRONT 2099 netmask 255.255.255.255 0 0
```

* To delete this rule:

```bash
no static (dmz,outside) tcp 193.252.19.3 2099 SRV-FRONT 2099 netmask 255.255.255.255 0 0
```

### Adding an Admin User

To add a Cisco admin user:

```bash
username login password pass privilege 15
```

### Creating an Address Pool for the Client

Here's an example of creating an address pool for client access:

```bash
name CLIENT_IP VANLANSCHOT_TEST
name CLIENT_IP VANLANSCHOT_PROD
object-group network VANLANSCHOT_RANGE
 network-object VANLANSCHOT_TEST 255.255.255.255
 network-object VANLANSCHOT_PROD 255.255.255.255
access-list radianz_access_in permit tcp object-group VANLANSCHOT_RANGE host LOCAL_SERVER_IP eq 9024
pdm location VANLANSCHOT_TEST 255.255.255.255 radianz
pdm location VANLANSCHOT_PROD 255.255.255.255 radianz
route radianz VANLANSCHOT_TEST 255.255.255.255 GATEWAY_IP 1
route radianz VANLANSCHOT_PROD 255.255.255.255 GATEWAY_IP 1
```

## VPN Setup

### Introduction

To create a VPN, you need:

* The remote IP of the person and their local IP/local network through which they will connect.
* Ask the person if they want the shared key system (pre-shared key)?
* Provide this type of encryption: DES-MD5 - Group 2

### Creation

```bash
# Access List
# Enter the client's local IPs
access-list inside_outbound_nat0_acl permit ip host OUR_LOCAL_IP host CLIENT_LOCAL_IP
access-list outside_cryptomap_240 permit ip host OUR_LOCAL_IP host CLIENT_LOCAL_IP

# IPSec Encryption
#crypto ipsec transform-set ESP-3DES-MD5 esp-3des esp-md5-hmac 
 
# Crypto Map Configuration 
# This is the line that reads the policy
# Check if 240 exists, otherwise add +
crypto map outside_map 240 ipsec-isakmp
crypto map outside_map 240 match address outside_cryptomap_240
#crypto map outside_map 240 set pfs group2
crypto map outside_map 240 set peer CLIENT_ROUTER_IP
crypto map outside_map 240 set transform-set ESP-3DES-MD5
#crypto map outside_map 240 set security-association lifetime seconds 86400 kilobytes 10000
#crypto map outside_map interface outside

# ISAKMP Pre-Shared Key
# Enter the shared key here
#isakmp enable outside
isakmp key PRE-SHARED_KEY address CLIENT_ROUTER_IP netmask 255.255.255.255 no-xauth no-config-mode 
#isakmp identity address

# ISAKMP Encryption
# Add if it doesn't exist
isakmp policy 160 authentication pre-share
isakmp policy 160 encryption 3des
isakmp policy 160 hash md5
isakmp policy 160 group 2
isakmp policy 160 lifetime 86400
```

In case of conflicts between networks, it may be necessary to NAT our network. For this, don't create the access-list inside_outbound_nat0_acl but add the following lines:

```bash
access-list nat_to_customer permit ip host OUR_LOCAL_IP host CLIENT_LOCAL_IP
static (inside,outside) OUR_NATED_IP access-list nat_to_customer 0 0
```

### Enabling Debug Mode for VPN

```bash
debug crypto isakmp
```

```bash
debug crypto ipsec
```

Add "2" at the end of the line to increase the debug level.

Be sure to disable debug mode when finished as it consumes resources.

```bash
no debug crypto isakmp
```

```bash
no debug crypto ipsec
```

### Closing a VPN Connection

Enter configuration mode and execute the following command:

```bash
clear crypto sa peer ip_address_of_the_remote_host
```
