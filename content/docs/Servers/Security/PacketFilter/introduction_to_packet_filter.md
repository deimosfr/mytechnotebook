---
weight: 999
url: "/Introduction_à_Packet_Filter/" 
title: "Introduction to Packet Filter"
description: "A comprehensive guide to Packet Filter (PF), the firewall software for OpenBSD and other BSD systems. Learn about installation, configuration, NAT, filtering rules and advanced features."
categories: ["FreeBSD", "Linux", "Network"]
date: "2010-12-13T17:15:00+02:00"
lastmod: "2010-12-13T17:15:00+02:00"
tags: ["Security", "Firewall", "OpenBSD", "BSD", "Network", "VPN"]
toc: true
---

## Introduction

[Packet Filter](https://fr.wikipedia.org/wiki/Packet_Filter) (or PF) is the official software firewall for **OpenBSD**, originally written by Daniel Hartmeier. It is a free Open Source software.

It replaced Darren Reed's IPFilter since OpenBSD version 3.0, due to licensing issues and also Reed's systematic refusal to incorporate code modifications from OpenBSD developers.

It has been ported to DragonflyBSD 1.2 and NetBSD 3.0; it is provided as standard on FreeBSD (version 5.3 and later).

A free port of PF has also been created for Windows 2000 and XP operating systems by the Core FORCE community. However, this port is only a personal firewall: it does not implement PF functions that allow NAT or the use of ALTQ.

## Installation

### FreeBSD

Insert the following entries in `/etc/rc.conf`:

```bash
pf_enable="YES"                 # Load PF
pf_rules="/etc/pf.conf"         # Rules, default path.
pf_flags=""                     # Bonus
pflog_enable="YES"              # start logging
pflog_logfile="/var/log/pflog"  # where to find the log
pflog_flags=""                  # Bonus
```

### OpenBSD

Vive OpenBSD :-), no need to install anything, it's included by default in the system. Just configure and enable it. To configure it, edit the `/etc/pf.conf` file.

## Configuration

We'll look at several examples to learn how to use it.

### Macros & Lists

#### Macros

We can define macros to replace variables that are used frequently (IPs, interfaces...):

```bash
if_ext="ne0"
ssh="192.168.0.1"
```

#### Lists

Lists allow us to group criteria into a variable:

```bash
block out on $if_ext proto tcp from { 192.168.0.1, 10.8.0.6 } to any port { 22 80 }
```

This line allows us to execute 4 rules at once! For beginners, this line will block:

* 192.168.0.1 on port 80
* 192.168.0.1 on port 22
* 10.8.0.6 on port 80
* 10.8.0.6 on port 22

#### Mix of Macros and Lists

To mix everything:

```bash
if_ext="ne0"
trusted="{ 192.168.0.1, 10.8.0.6 }"
pass in on $if_ext inet proto tcp from { 10.10.0.0/24 $trusted } to port 22
```

### Tables

Tables allow you to store a large number of addresses (50 or 50,000, it's the same), which are then used directly in filtering/NAT/redirection rules. Searching for an address in a memory table is much faster and less CPU/memory intensive than searching through a set of rules each corresponding to a value in an address list.

There are several keywords for tables with different functions:

* **cont**: used when you want the table to not be modifiable
* **persist**: tells PF not to delete a table that isn't referenced by a rule

The advantage of a **non**-const table compared to lists is that you can add/remove addresses or subnets on-the-fly, **useful for temporarily blocking a spammer's address, a script-kiddie or managing redirection to a set of high-availability servers.**

Finally, as the icing on the cake, you can initialize a table with a file containing a list of addresses:

```bash
table <privateip> const { 192.168.0.0/24, 10.8.0.0/8 }
table <spammers> persist file "/etc/spammers"
block in on $if_ext from { <privateip>, <spammer> } to any
```

And now, if we want to add or remove IPs:

```bash
pfctl -t spammers -T add 218.70.0.0/16
pfctl -t spammers -T delete 218.70.0.12
```

To delete everything:

```bash
pfctl -t spammers -T flush
```

### NAT

I won't explain here what NAT does, but how to use it with PF. First, we activate packet forwarding by adding this to `/etc/sysctl.conf`:

```bash
net.inet.ip.forwarding=1
```

This will make the NAT persistent.  
Remember that packets will pass through the packet filter after being modified, unless the **pass** keyword is used. This also applies to **RDR** which we'll see below. Here is a NAT rule:

```bash
nat [pass] on interface [address_family] from src_addr [port src_port] to dst_addr [port dst_port] -> ext_addr
```

If we break it down (in brackets: variable, italic: optional):

* nat: indicates that this is a nat rule
* pass: the packet is NATed and sent directly without going through the packet filter
* on *interface*: the packet arrived on this network interface ($if_ext, ne0...)
* *address_family*: **inet or niet6**, this is a detail that could be important
* from *src_addr*: the packet comes from this address. For the address, you can specify many things:
  * an IP address
  * a CIDR
  * a DNS that will be resolved by PF when loading the rules
  * a network interface
  * a [Table](#tables) or a [List](#lists)
  * any of these notations, preceded by a ! to signify negation
  * finally **any** for any address
* port *src_port*: if you want to NAT only a certain port or range of ports...rarely used
* to *dst_addr*: the packet is destined for this address. Same possibilities as for src_addr
* -> ext: replace the source address with this address. The return will be handled automatically. And if this address changes (assigned by DHCP), you can specify the name of the network interface in parentheses (rl0), and the address will be automatically updated in the rule.

Now a small example. If you want to share your internet connection with your local network:

```bash
nat on r10 inet from ne3: network to any -> (r10)
```

That's it, it's not the same as with Iptables!

### Packet Redirections

RDR is NAT's hidden little brother, in that it does exactly the opposite: it takes packets coming from the outside to redirect them to the local network. Here's an example of syntax:

```bash
rdr pass on interface [address_family] from src_addr [port src_port] to dst_addr [port dst_port] -> int_addr [port int_port]
```

It's the same syntax as for NAT with a few exceptions like: I redirect what was destined for **dst_addr:dst_port** to **int_addr:int_port** and as with NAT, the return will be handled automatically.

For example, if I want to access the SSH of one of my machines on the local network from the outside:

```bash
rdr pass on $external proto tcp to port 35422 -> $diane port ssh
```

And now, you can access the machine in my LAN from outside via my external IP address and port 35422.

You'll notice that I used ssh as a name. All names present in `/etc/services` can be used.

### Filtering Rules

Filtering rules are evaluated sequentially from the first to the last (from top to bottom in the rules files used). This means that each packet will be evaluated by each rule, and the last rule matching the packet wins the decision (**block** or **pass**). On the contrary, if we use the **quick** keyword, the evaluation stops as soon as a rule matches the packet. The first (implicit) rule is "let everything through" so that if no rule applies to a packet, it is accepted. This is why the first explicit rule is usually a **block all**.

Here's the syntax:

```bash
action [direction] [log] [quick] [on interface] [address_family] [proto protocol] [from src_addr [port src_port]] [to dst_addr [port dst_port]] [flags tcp_flags] [state]
```

* action: choose between **block** or **pass**. The policy for blocked packets will be **drop** or **return** depending on the **block-policy** option. By default it's drop.
* direction: **in** or **out**, if you want to filter incoming or outgoing traffic on the interface. If nothing is specified, the rule will be evaluated for both directions.
* log: if this flag is present, we record the decision made by this rule concerning the packet. To analyze this, **pflog** will be our friend.
* quick: I've already mentioned this - if this flag is also present and the packet matches the rule, it will no longer be analyzed/manipulated, and the decision made by this rule will be final.
* proto protocol: a level 4 protocol, generally tcp, udp or icmp, but we can also encounter any level 4 protocol referenced in `/etc/protocols`. We can even call it by its little number!
* port dst_port: Here you can specify a complex range of ports with operators <, <=, >=, >, <> and :, see man pf.conf.
* flags tcp_flags_check/mask: you can specify additional checks on the flags of a TCP packet, for example to handle TCP session openings. We often use flag S/SA which I would translate as "this rule applies to TCP packets which, on the two SYN and ACK flags, only have SYN set". If both flags are set, the packet will not match the rule. For other flags, RTFM a bit.
* state: here we generally use two possibilities:
  * keep state: used when we want to create an entry in the connection state table when a packet matches the rule, and apply the same policy to subsequent packets taking part in the connection. All these packets are therefore attached to this entry, and we can also check that the TCP packet sequence is respected.
  * synproxy state is used when we want PF to act as a TCP proxy for establishing a connection. In this case, PF will handle the request in place of the recipient and will only forward the packets to the latter afterwards. No packets are forwarded to the recipient before the client has completed the initial exchange. **This technique helps protect the recipient from TCP SYN flood attacks, where a large number of connection openings are requested in order to cause a denial of service.**

Here's a small example to clarify all this:

```bash
pass in quick on $external inet proto tcp from any to any port { http, https, smtp, imaps } flags S/SA keep state
```

I allow all TCP/IPv4 packets arriving on the external interface destined for http/https/smtp/imaps ports to pass. I check that these are TCP connection openings, I record their state in the table, and I stop the analysis of these packets at this rule (quick).

```bash
block in log on $external from { 192.168.0.0/16, 172.16.0.0/12 } to any.
```

Here, I block packets arriving on the external interface with a private source address, and I log the information of the blocked packet. **This helps prevent certain spoofing attacks** where spoofed packets are sent in order to mislead network equipment.

Of course, there are still plenty of detailed options and particularities (such as anchors, **scrubbing**, antispoofing...) that I haven't mentioned. For more information, refer to the pf.conf documentation.

## Usage

The binary for using pf is **pfctl**. Before starting, you need to enable PF at the kernel level. 2 solutions:

```bash
pfctl -e
```

or add pf=YES to `/etc/rc.conf.local`.

To disable pf, simply do:

```bash
pfctl -d
```

If you want to do a syntactic analysis of the file without loading the filtering rules, use the -n argument:

```bash
pfctl -n
```

If you want to optimize PF, add the -O option, which will remove duplicates and reorder rules.

### On-the-fly Modifications

* **-T (kill/flush/add/delete/show/test..)**: used with -t table, allows you to manipulate a table: delete it, empty it, add an address, delete, display it, check if it's in the table. Example:

```bash
pfctl -t blocked-hosts -T show
```

This will display the addresses of all machines that have been added to the blocked-hosts table, declared a little earlier in `/etc/pf.conf`.

* **-F (nat/rules/state/Tables/..)**: resets NAT rules, filtering rules, states of open connections or tables, respectively. Useful if you want to clean up a bit, reset counters or connections, disable NAT, delete all entries from all tables, etc...
* **-k (host/network)**: Allows you to kill all entries in the state table concerning connections from a machine/network. If you use this option twice, you delete the states of connections from the first address to the second. Example:

```bash
pfctl -k 192.168.1.0/24 -k 172.16.0.0/16
```

This will delete all states of connections between these 2 subnets.

* **-s modifier**: This option allows you to get a lot of information about the status of PF. If you use it with -r, PF will do reverse-dns lookup for the addresses it displays. The most interesting values for modifier are:
  * rules: display the loaded filtering rules in memory
  * nat: NAT rules
  * state: open connections
  * info: global statistics on PF
  * all: will display everything PF has to tell us

e.g:

```bash
pfctl -sr
```

* **-v, -vv, -g, -x, -q**: for more verbose modes, and even debug mode (-v -> -x). The -q will put it back in quiet mode.

## Practical Examples

### Simple

Requirements:

* We have a machine connected to the internet via a network interface with a fixed IP called bge0
* It must be accessible from outside via 3 trusted IPs
* The Apache server must be accessible to everyone
* It must respond to pings from the outside
* We want to be able to access the internet from this machine
* We will silently block all other packets

```bash {linenos=table}
# Our network interface
iface = "bge0"

# Trusted machines for SSH connections:
trusted_hosts = "{ 131.25.4.12, 88.12.74.5, 207.124.20.9 }"

# We don't worry about the loopback interface used by several internal services on the machine:
set skip on lo

# We enable packet normalization on input. PF will reassemble fragmented packets and perform additional checks on them:
match in all scrub (no-df)

# By default, we block all packets:
block all

# We allow ICMP packets of type "echo request" for pings from outside, and echo reply/time exceeded/destination
# unreachable for responses to pings that we had initiated towards the outside:
pass in inet proto icmp form any to $iface icmp-type { echoreq, echorep, timex, unreach }

# We allow connections to the Apache server and record them in the state table:
pass in inet proto tcp from any to $iface port www flags S/SA keep state

# We only allow SSH connections from trusted machines:
pass in inet proto tcp from $trusted_hosts to $iface port ssh flags S/SA keep state

# Finally we allow all outgoing traffic:
pass out inet proto tcp from $iface to any flags S/SA keep state
pass out inet proto { udp, icmp } from $iface to any keep state
```

Now we activate our new configuration:

```bash
pfctl -e -f /etc/pf.conf
```

This configuration is quite restrictive in the sense that generally, we put all instead of "from any to $iface" and we use "quick" extensively.

### Advanced

Requirements:

* Our gateway (renton) has 2 interfaces, r10 on the internet side and ne3 on the local network side
* The first card has a fixed IP (that of our ISP)
* The second interface on the local network side is configured with the address 192.168.1.1 on the local network in 192.168.1.0/24
* We will protect all of this from the outside world, and since we're paranoid, we'll do a bit of logging
* We want to be able to access the internet from our local network without limitations
* A particular machine (diane) must be accessible via SSH and HTTPS
* We want to access Free Multipost on the machine (sickboy), it's just a standard RTSP/UDP stream
* Our main machine (tommy) would like to make clicka-compliant things like Jabber file transfer work
* We will allow SSH from the outside

```bash {linenos=table}
# /etc/pf.conf

# Network interface declaration
ext="r10"
int="ne3"

# Declaration of ports not to log
ports_not_logged = "{ netbios-ssn, microsoft-ds, epmap, ms-sql-s, 5900 }"

# Declaration of hosts on my local network
diane = "192.168.1.2"
tommy = "192.168.1.20"
sickboy = "192.168.1.21"

# Declaration of the port used on the outside to access diane's SSH
ssh_diane = "65322"

set skip on lo
match in all scrub (no-df)

# First, the main NAT rule for the local network. Here, I don't put "pass" because later I explicitly allow all
# outgoing traffic. The :network suffix is used to say "the subnet corresponding to the address of this interface"
nat on $ext from $int:network to any -> $ext

# Next, the redirection rules. I added the 'pass' keyword to not do additional filtering on these connections,
# otherwise I would have had to add the corresponding ports in the filtering rule for traffic from the outside.

# The rtsp stream from Free's multiposte is UDP coming from freeplayer.freebox.fr
rdr pass on $ext protoudp from 212.27.38.253 -> $sickboy

# We redirect port 8010 used by Jabber for opening incoming file transfers
rdr pass on $ext proto tcp to port 8010 -> $tommy

# Finally, we redirect SSH (on the specific port used from the outside) and https to diane
rdr pass on $ext proto tcp to port $ssh_diane -> $diane port ssh
rdr pass on $ext proto tcp to port https -> $diane

# We enable antispoofing on the external interface to block packets coming from outside trying to fraudulently
# use our address to pass through the filter
antispoof for $ext

# We don't filter packets on the internal interface
pass quick on $int

# Now the filtering rules
# By default we block and log all packets from the outside
block in on $ext log all

# We don't want to fill our logs in 5 minutes with the few well-known worms floating around the net. So we block
# certain packets without logging them
block in on $ext inet proto tcp from any to any port $ports_not_logged

# We allow pings from the outside
pass in on $ext inet proto icmp from any to any icmp-type { chorep, echoreq, timex, unreach }

# We allow SSH from the outside
pass in on $ext inet $proto tcp from any to any port ssh flags S/SA keep state

# We allow all outgoing traffic (the NATed from the local network will go through these rules)
pass out inet proto tcp from $iface to any flags S/SA keep state
pass out inet proto { udp, icmp } from $iface to any keep state
```

### My Config

I have a Soekris with several network interfaces:

* Wan
* DMZ
* Lan
* Wifi
* VPN

And I want the configuration to be as follows:

* Wifi and Lan have access to everything
* The DMZ only has access to Wan
* There is NAT on all interfaces except wan obviously
* Protection against SSH bruteforce
* VPN is accessible on the gateway from the outside

```bash {linenos=table}
#       $OpenBSD: pf.conf,v 1.37 2008/05/09 06:04:08 reyk Exp $
#
# See pf.conf(5) for syntax and examples.
# Remember to set net.inet.ip.forwarding=1 and/or net.inet6.ip6.forwarding=1
# in /etc/sysctl.conf if packets are to be forwarded between interfaces.

#-------------------------------------------------------------------------
# Physical and virtual interfaces definitions
#-------------------------------------------------------------------------

# Physical Interfaces
wan_if="sis0"
lan_if="vr0"
dmz_if="sis1"
wifi_if="vr1"
# VLAN Interfaces
vlan99_if="vlan99"
vlan110_if="vlan110"
vlan120_if="vlan120"
# TUN Interfaces
openvpn_if="tun0"
sshvpn_if="tun1"

#-------------------------------------------------------------------------
# Networks definitions
#-------------------------------------------------------------------------

# Wan
wan_net="192.168.10.0/24"

# Trusted networks
lan_net="192.168.0.0/24"
vlan99_net="192.168.99.0/24"
wifi_net="192.168.200.0/24"

# Remote trusted networks
openvpn_net="192.168.20.0/24"
openvpn_net_nat="10.0.0.0/24"
sshvpn_net="192.168.30.0/24"

# DMZ
dmz_net="192.168.100.0/24"
vlan110_net="192.168.110.0/24"
vlan120_net="192.168.120.0/24"

# OpenVPN Shenzi network
openvpn_shenzi_net="192.168.90.0/24"

#-------------------------------------------------------------------------
# IP definitions
#-------------------------------------------------------------------------

# Router IP interfaces
wan_sks_ip="192.168.10.254"
dmz_sks_ip="192.168.100.254"
vlan110_sks_ip="192.168.110.254"
vlan120_sks_ip="192.168.120.254"

# Others IP
dedibox_ip="x.x.x.x"
work_ip="x.x.x.x"
freebox_tv_ip="212.27.38.253"

# Services IP
dmz_mail_ip="192.168.100.3"
dmz_web_ip="192.168.110.2"
dmz_dns_ip="192.168.100.3"
dmz_sftp_ip="192.168.100.6"
apt_cacher_ip="192.168.120.2"

#-------------------------------------------------------------------------
# Ports definitions and options
#-------------------------------------------------------------------------

# Port descriptions
imaps_ports="143, 993"
smtps_ports="25, 465"
ssh_ports="22, 222"
dns_port="53"
webs_ports="80, 443"
openvpn_port="1194"
proxy_port="3128"
apt_cacher_port="3142"
mysql_port="3306"
git_port="9418"
free_multiposte="31336, 31337"

# Whitelist / Blacklist table
table <blacklist> persist
table <whitelist> persist file "/etc/ssh/whitelist"

# Do not touch lo interface
set skip on lo0

#-------------------------------------------------------------------------
# Packet Normalization: reassemble fragments
#-------------------------------------------------------------------------

match in all scrub (no-df)

#-------------------------------------------------------------------------
# Nat for all internal interfaces
#-------------------------------------------------------------------------

match out on $wan_if from !($wan_if) nat-to ($wan_if)

#-------------------------------------------------------------------------
# block in all with no usurpation
#-------------------------------------------------------------------------

block in log all
block in log quick from urpf-failed

#-------------------------------------------------------------------------
# Redirections for incoming connections (wan)
#-------------------------------------------------------------------------

# From WAN
pass in on $wan_if proto tcp from any to $wan_if port 25 rdr-to $dmz_mail_ip port 25
pass in on $wan_if proto udp from any to $wan_if port $dns_port rdr-to $dmz_dns_ip port $dns_port
pass in on $wan_if proto tcp from any to $wan_if port $dns_port rdr-to $dmz_dns_ip port $dns_port
pass in on $wan_if proto tcp from any to $wan_if port 80 rdr-to $dmz_web_ip port 80
pass in on $wan_if proto tcp from any to $wan_if port 143 rdr-to $dmz_mail_ip port 143
pass in on $wan_if proto tcp from any to $wan_if port 222 rdr-to $dmz_sftp_ip port 22
pass in on $wan_if proto tcp from any to $wan_if port 443 rdr-to 127.0.0.1 port 443
pass in on $wan_if proto tcp from any to $wan_if port 465 rdr-to $dmz_mail_ip port 465
pass in on $wan_if proto tcp from any to $wan_if port 993 rdr-to $dmz_mail_ip port 993
pass in on $wan_if proto tcp from any to $wan_if port 9418 rdr-to $dmz_web_ip port 9418

# Apt-cacher-ng
pass in on $vlan120_if proto tcp from any to $vlan120_if port $apt_cacher_port rdr-to $apt_cacher_ip port $apt_cacher_port

pass in on $wan_if proto tcp from $dedibox_ip to $wan_if port $mysql_port rdr-to $dmz_web_ip port $mysql_port

pass in on $wan_if proto udp from $freebox_tv_ip to $wan_if rdr-to 192.168.0.100
pass in on $wan_if proto tcp from $freebox_tv_ip to $wan_if rdr-to 192.168.0.100

#-------------------------------------------------------------------------
# Global Rules pass and block
#-------------------------------------------------------------------------

# Allow Free Multiposte
pass in quick on $wan_if proto udp from $freebox_tv_ip to 192.168.0.100
pass out quick on $wan_if proto udp from $freebox_tv_ip to 192.168.0.100
pass in quick on $wan_if proto tcp from $freebox_tv_ip to 192.168.0.100
pass out quick on $wan_if proto tcp from $freebox_tv_ip to 192.168.0.100

# Allow all outgoing from $lan_net, $wifi_net and $openvpn_net
pass in on { $lan_if, $vlan99_if, $wifi_if, $openvpn_if, $sshvpn_if } from { $lan_net, $vlan99_net, $wifi_net, $openvpn_net, $sshvpn_net } to any
pass out on { $lan_if, $vlan99_if, $wifi_if, $openvpn_if, $sshvpn_if } from { $lan_net, $vlan99_net, $wifi_net, $openvpn_net, $sshvpn_net } to any
pass out on $wan_if from $wan_net to any
antispoof quick for { $wan_if, $dmz_if, $vlan110_if, $vlan120_if }

# block all incoming on lan_if, wifi_if and openvpn_if
block out log on { $lan_if, $vlan99_if, $wifi_if, $openvpn_if, $sshvpn_if } from { !$lan_if, !$vlan99_if, !$wifi_if, !($openvpn_if), !($sshvpn_if) } to any

#-------------------------------------------------------------------------
# VPN Access
#-------------------------------------------------------------------------

# Allow to access shenzi VE
pass out on $openvpn_if from any to $openvpn_shenzi_net

#-------------------------------------------------------------------------
# Specific ports on dmz
#-------------------------------------------------------------------------

# DNS
pass in on $dmz_if proto tcp from $dmz_dns_ip to $dmz_sks_ip port $dns_port
pass in on $dmz_if proto udp from $dmz_dns_ip to $dmz_sks_ip port $dns_port

# Apt-cacher
pass out on $vlan120_if proto tcp to ($vlan120_if) port $apt_cacher_port
pass in on $vlan120_if proto tcp to ($vlan120_if) port $apt_cacher_port

# DMZ and Vlan autorisations
#pass in on $vlan110_if from $vlan110_net to {  !$lan_net,  !$wifi_net,  !$openvpn_net,  !$sshvpn_net,  !$vlan120_net,  !$vlan110_sks_ip }
# Arrive pas a avoir juste any -> 3142
pass in on $vlan110_if from $vlan110_net to { !$lan_net, !$wifi_net, !$openvpn_net, !$sshvpn_net, !$vlan110_sks_ip }
pass in on $vlan120_if from $vlan120_net to { !$lan_net, !$wifi_net, !$openvpn_net, !$sshvpn_net, !$vlan120_sks_ip }
pass in on $dmz_if from $dmz_net to { !$lan_net, !$wifi_net, !$openvpn_net, !$sshvpn_net, !$dmz_sks_ip }

#-------------------------------------------------------------------------
# Specific Rules pass and block
#-------------------------------------------------------------------------

# Allow all incoming ICMP
pass in on $wan_if proto icmp to any

# Autoblacklist on SSH
pass in on $wan_if proto tcp from !<whitelist> to ($wan_if) port { $ssh_ports } \
        flags S/SA keep state \
        (max-src-conn-rate 3/60, \
        overload <blacklist> flush global)
pass in on $wan_if proto tcp from <whitelist> to $wan_if port { $ssh_ports } flags S/SA keep state

# block the ssh bruteforce bastards
block drop in on $wan_if from <blacklist>
pass in on $wan_if proto tcp from <whitelist> port { $ssh_ports }

# Allow OpenVPN
pass in on $wan_if proto tcp to $wan_if port $openvpn_port

# Allow DNS TCP from DMZ (Bind SRV) to secondary (Soekris)
#pass in quick on $dmz_if proto tcp from $dmz_dns_ip to $dmz_sks_ip port $dns_port

# Allow on wan interface from wan for tcp
pass out on $wan_if proto tcp to ($wan_if) port { $ssh_ports, $smtps_ports, $imaps_ports, $dns_port, $webs_ports, $git_port, $mysql_port }
# Allow on dmz interface from wan for udp
pass out on $wan_if proto udp to ($dmz_if) port { $dns_port }

# Allow all outbound traffic
pass out inet from !($wan_if) to any flags S/SA keep state
```

## Logs

When PF wants to report something, it will send binary information (to make it more fun, it's standard PCAP/TCPdump) to a pseudo interface (pflog0), and one of its good friends **pflogd** will store everything in `/var/log/pflog`.

First, you need to enable/start the pflogd daemon. Normally it should start automatically if PF is enabled when the machine boots. If not:

```bash
ifconfig pflog0 up && pflogd
```

We can check that everything is working properly:

```bash
ps waux 
```

Now that PF is talking on pflog0, when a packet matches a rule where the log keyword is used, let's move on to tcpdump. It can be used in 2 modes:

* Interactive:

```bash
tcpdump -n -e -ttt -i pflog0
```

It will directly read what's happening live on pflog0, so pflogd will be useful.

* Passive:

```bash
tcpdump -r /var/log/pflog
```

It will read what has been recorded by pflogd in its output file.

You can also pass an expression to tcpdump for it to filter its output according to specific criteria:

```bash
tcpdump -ttt -r /var/log/pflog port 80 and host 192.168.1.2
```

Finally, tcpdump also understands PF configuration syntax. So we can ask it things like:

```bash
tcpdump -ttt -i pflog0 inbound and action pass and on wi0
```

With this command, it will only show packets allowed to pass through, logged and incoming on the wi0 interface. tcpdump can also read information such as passive OS fingerprint.

## Advanced PF Functions

If you want to block Windows 95/98 machines for example:

```bash
block in on $ext proto tcp from any os { "Windows 95", "Windows 98" } to any port smtp
```

To get the list of OS that PF recognizes:

```bash
pfctl -so
```

To do OS analysis with tcpdump:

```bash
tcpdump -o -r /var/log/pflog
```

Obviously these techniques are CPU intensive and not infallible.

## FAQ

### no IP address found for tun0

If you get these kinds of messages, it's because an interface (here tun0) is trying to be initialized with PF, while the associated service (supposed to create this device) hasn't started yet.

To avoid chaos, just try to put your devices in parentheses (!($vpn_if)) eg:

```bash
block out quick on { $wifi_if } from { !$lan_if, !$wifi_if, !($vpn_if) } to any
```

And if parentheses already exist, try to remove them.

## Resources
- [Documentation for Securing your BSD with PF](/pdf/secure_bsd_pf.pdf)
- [Firewalling with OpenBSD's PF packet Filter](/pdf/pf-firewall.pdf)
- [Packet filtering with IPFilter software](/pdf/filtrage_des_paquets_avec_le_logiciel_ipfilter.pdf)
- [https://www.openbsd.org/faq/pf/filter.html](https://www.openbsd.org/faq/pf/filter.html)
