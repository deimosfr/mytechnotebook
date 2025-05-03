    ---
weight: 999
url: "/Configuration_d'un_Switch_Catalyst/"
title: "Configuring a Catalyst Switch"
icon: "article"
description: "Learn how to configure and manage Cisco Catalyst switches through console, firmware updates, VLAN setup, and more."
categories: ["Network", "Cisco", "Infrastructure"]
date: "2007-11-20T14:57:00+01:00"
lastmod: "2007-11-20T14:57:00+01:00"
tags: ["Cisco", "Switch", "Catalyst", "VLAN", "Networking"]
toc: true
---

## Connections

Start by connecting the serial port on the PC side and the console port on the switch.

## Connection

Launch HyperTerminal and choose a connection on com1 port, or another com port depending on the PC configuration.
You need to choose 9600 Bauds
The rest doesn't matter much...

## What You Need to Know

When launching the terminal, you'll determine several things:

- The router name (for example: swpar5)
- A secret
- A password
- A virtual terminal password

Then you need to switch to "enable" mode to be able to change the configuration

```bash
enable
```

The prompt changes from swpar5> to swpar5#

In this mode, you can view the router configuration by typing

```bash
sh run
```

## Resources

[Cisco Basic Concepts](https://cisco.goffinet.org/s2/notions_routeurs)

If you want to modify this configuration, you need to switch to "conf" mode

```bash
conf t
```

Result: the prompt changes to swpar5(config)#

The "t" is to indicate that you want to configure the router in terminal mode.

### Configuring the IP Address of an Interface

By running sh run, you can see the different interfaces, like this:

```bash
interface FastEthernet0/24
!
interface GigabitEthernet0/1
!
interface GigabitEthernet0/2
!
interface Vlan1
 no ip address 
 no ip route-cache
!
ip http server
!
control-plane
!
!
line con 0
line vty 0 4
 password *xxxxx*
 no login
line vty 5 15
```

Here we see an interface named Vlan1 that has no IP address, and we notice that the two lines that follow are indented with a space, so we can modify them by typing, in "configuration" mode:

```bash
interface Vlan1
```

The prompt then changes to swpar5(config-if)#
Then type

```bash
ip address 192.168.0.169 255.255.255.0
```

Then return to enable mode (with the prompt swpar5#) and type:

```bash
wr m
```

to save

From there, you can access the router's web interface by typing
```
http://192.168.0.169
```

### Firmware Update

#### Checking Firmware Version

Type

```bash
sh ver
```

#### Downloading the Latest Version

Go to cisco.fr --> support --> download software --> software advisor --> Find software compatible with my hardware --> Select your hardware configuration manually --> then choose the router model and download the archive.
You have the choice between software that supports CRYPTO (secure https connection) or not.

#### Installing the New Firmware

It's a .tar archive that needs to be loaded from the router's web interface. Everything is automatic.

### Configuring the VTP Domain

```bash
swpar5(config)#vtp mode client     (or server or transparent)
swpar5(config)#vtp domain ulnet-fr 
```

http://www.cisco.com/en/US/docs/switches/lan/catalyst4500/12.1/12ew/configuration/guide/vtp.html#wp1032093

### Configuring an Interface

Setting trunk mode allows multiple VLANs to pass through a single port between multiple switches:

```bash
swpar5(config-if)#switchport mode trunk
```

```bash
swpar5(config-if)#switchport nonegotiate
```

```bash
swpar5(config-if)#auto qos voip trust 
```

```bash
swpar5(config-if)#macro description cisco-switch
```

```bash
swpar5(config-if)#spanning-tree link-type point-to-point
```

### Creating a VLAN

Switch to configuration mode, and type 

```bash
swpar5(config)#vlan x
```

where "x" is an "id" for the VLAN (a number), then choose a name for this VLAN. To do this, in conf mode, type "vlan x"

```bash
swpar5(config-vlan)#name vlan_name
```

You can verify the VLAN configuration with the command 

```bash
swpar5#sh vlan
```

Once the VLAN is created, it's better to put it in "no-shutdown" so that it doesn't automatically deactivate:

```bash
swpar5(config-if)#no shutdown
```

http://www.cisco.com/en/US/docs/switches/lan/catalyst2960/software/release/12.2_25_fx/configuration/guide/swvlan.html#wp1273595

#### Assigning an Interface to a VLAN

Enter configuration mode and choose an interface, for example "interface FastEthernet0/1", and type

```bash
swpar5(config-if)#switchport mode access
swpar5(config-if)#switchport access vlan 3 (if the VLAN is called "vlan 3")
```

### Graphical Mode

Once you have assigned an IP address, you can access the switch's graphical interface via http. When asked for a username and password, you only need to enter the password that you set at the beginning.

## Resources

- [Cisco Basic Concepts](https://cisco.goffinet.org/s2/notions_routeurs)
