---
weight: 999
url: "/Scapy_\\:_Trames_et_paquets_de_données/"
title: "Scapy: Data Frames and Packets"
description: "A comprehensive guide to using Scapy for network packet manipulation, inspection and analysis. Learn how to forge, receive and send data packets over a network with this powerful Python tool."
categories: ["Security", "Network", "Linux"]
date: "2007-10-09T06:44:00+02:00"
lastmod: "2007-10-09T06:44:00+02:00"
tags:
  [
    "network",
    "security",
    "python",
    "packet analysis",
    "networking",
    "penetration testing",
  ]
toc: true
---

## Introduction

Scapy is a utility that allows you to forge, receive and send packets or data frames over a network for a multitude of protocols. In this introduction, you'll discover this Python utility that enables traffic capture, network mapping, ARP cache poisoning, VLAN hopping, or passive operating system fingerprinting.

Scapy is a program developed in Python by Philippe Biondi (EADS CCR); it notably allows you to forge, receive and transmit packets and/or data frames via a network to or from an IT infrastructure for a multitude of different network protocols (IP, TCP, UDP, ARP, SNMP, ICMP, DNS, DHCP, ...) with precision and speed.

Scapy comes in the form of a single Python script file - 13,342 lines of code for version 1.1.1 that we'll use throughout this document. Among other notable features of Scapy, we'll note its ability to dissect packets and/or data frames as well as decode certain protocols.

Furthermore, Scapy can also perform network traffic monitoring and capture similar to reading pcap format captures from another traffic analyzer like Wireshark, for example.

It's also possible with Scapy to generate graphs in 2D and/or 3D from packets and/or data frames, or even port scanning like NMAP and passive remote operating system recognition like p0f.

According to its author, Scapy is capable by itself of replacing all of the following utilities: hping, 85% of NMAP, arpspoof, arp-sk, arping, tcpdump, tethereal, p0f and many other system commands (traceroute, ping, route, ...).

For the equivalent of about sixty lines of C code, the combination of Python and Scapy only requires a few lines most of the time to perform these different packet and/or data frame manipulation operations, resulting in considerable time savings for anyone who needs to perform this type of manipulation on a network.

For this, Scapy has many pre-defined functions allowing you to configure the injection of a packet (or frame) into a given network connection; some special functions of Scapy thus make it possible to perform common attacks with great simplicity (non-exhaustive list): network infrastructure mapping, ARP Cache Poisoning, Smurfing, VLAN Hopping as well as IP spoofing and rogue DHCP server setup.

These attacks can be combined with each other (ARP Cache Poisoning + VLAN Hopping for example) to perform security audits specifically adapted to the infrastructure in place whose security level you want to verify.

You can just as well intercept VOIP communications (packet/frame decoding) and this even on a WEP encrypted WIFI wireless connection as long as you know the decryption key associated with these connections (knowing of course that WEP is still secure).

This encryption key can be configured in Scapy, still provided that you have it, so that Scapy can use it during packet or data frame injection operations into the traffic of a WEP-encrypted wireless network (see also the WIFITAP utility developed by Cédric Blancher [EADS CCR] for traffic injection into WIFI connections).

## Installation and configuration

This section concerns the installation of Scapy as well as all the elements necessary for its proper functioning on a GNU/LINUX system. It also includes Scapy's internal configuration system as a first approach to the utility as well as the different network configurations necessary for its use for the rest of this document.

Let's perform the preliminary installation necessary for using Scapy on a Debian/Ubuntu-like Linux operating system:

```bash
root@casper:~# uname -a
Linux casper 2.6.20-15-generic #2 SMP Sun Apr 15 07:36:31 UTC 2007 i686 GNU/Linux

root@casper:~ # apt-get install python python-gnuplot python-pyx python-crypto graphviz imagemagick python-visual
```

First, we test that the Python interpreter is working properly:

```bash
root@casper:~# python
Python 2.5.1 (r251:54863, May 2 2007, 16:56:35)
[GCC 4.1.2 (Ubuntu 4.1.2-0ubuntu4)] on linux2
Type "help", "copyright", "credits" or "license" for more information.
>>>
```

Type "ctrl D" to exit Python.

We access the /etc directory on our system:

```bash
root@casper:~ # cd /etc/
```

Then we retrieve the ethertypes file:

```bash
root@casper:/etc # wget www.secdev.org/projects/scapy/files/ethertypes
02:41:59 (936.05 KB/s) - `ethertypes' saved [1,317/1,317]
```

We then access the personal directory of the current session user (here the root user with /root as the home directory ~):

```bash
root@casper:/etc # cd ~
```

Now we configure the network interfaces present on the operating system of the machine that we are using for all the tests in this document; this information will help better understand the different tests we will be performing:

```bash
root@casper:~# lspci | grep 802.11
08:00.0 Ethernet controller: Atheros Communications, Inc. AR5212 802.11abg NIC (rev 01)
root@casper:~# wlanconfig ath0 destroy
root@casper:~# wlanconfig ath0 create wlandev wifi0 wlanmode adhoc
root@casper:~# iwconfig ath0 essid nat
root@casper:~# ifconfig ath0 192.168.0.2
root@casper:~# route add default gw 192.168.0.1
root@casper:~# ifconfig
ath0 Link encap:Ethernet HWaddr 00:15:6D:53:1E:87
inet addr:192.168.0.2 Bcast:192.168.0.255 Mask:255.255.255.0
inet6 addr: fe80::215:6dff:fe53:1e87/64 Scope:Link
UP BROADCAST RUNNING MULTICAST MTU:1500 Metric:1
RX packets:3867 errors:0 dropped:0 overruns:0 frame:0
TX packets:3719 errors:0 dropped:0 overruns:0 carrier:0
collisions:0 txqueuelen:0
RX bytes:3782362 (3.6 MiB) TX bytes:520446 (508.2 KiB)

lo Link encap:Local Loopback
inet addr:127.0.0.1 Mask:255.0.0.0
inet6 addr: ::1/128 Scope:Host
UP LOOPBACK RUNNING MTU:16436 Metric:1
RX packets:50 errors:0 dropped:0 overruns:0 frame:0
TX packets:50 errors:0 dropped:0 overruns:0 carrier:0
collisions:0 txqueuelen:0
RX bytes:4307 (4.2 KiB) TX bytes:4307 (4.2 KiB)

wifi0 Link encap:UNSPEC HWaddr 00-15-6D-53-1E-87-00-00-00-00-00-00-00-00-00-00
UP BROADCAST RUNNING MULTICAST MTU:1500 Metric:1
RX packets:179191 errors:0 dropped:0 overruns:0 frame:7162
TX packets:4355 errors:147 dropped:0 overruns:0 carrier:0
collisions:0 txqueuelen:199
RX bytes:8345718 (7.9 MiB) TX bytes:668330 (652.6 KiB)
Interrupt:18


root@casper:~# iwconfig
lo no wireless extensions.
eth0 no wireless extensions.
wifi0 no wireless extensions.

ath0 IEEE 802.11g ESSID:"nat" Nickname:""
Mode:Ad-Hoc Frequency:2.462 GHz Cell: 02:15:6D:53:1E:87
Bit Rate:0 kb/s Tx-Power:16 dBm Sensitivity=1/1
Retry:off RTS thr:off Fragment thr:off
Encryption key:off
Power Management:off
Link Quality=26/70 Signal level=-70 dBm Noise level=-96 dBm
Rx invalid nwid:1181 Rx invalid crypt:0 Rx invalid frag:0
Tx excessive retries:0 Invalid misc:0 Missed beacon:0
```

The section above needs to be adapted to each machine configuration depending on the network interfaces present; for the next part we test that network connectivity is working properly:

```bash
root@casper:~# ping -c 1 192.168.0.1
PING 192.168.0.1 (192.168.0.1) 56(84) bytes of data.
64 bytes from 192.168.0.1: icmp_seq=1 ttl=128 time=1.40 ms

--- 192.168.0.1 ping statistics ---
1 packets transmitted, 1 received, 0% packet loss, time 0ms
rtt min/avg/max/mdev = 1.409/1.409/1.409/0.000 ms
```

That's it, our machine with IP address 192.168.0.2 is able to ping the machine with IP address 192.168.0.1; we can now proceed to install Scapy by first downloading it:

```bash
root@casper:~ # wget www.secuobs.com/scapy.py
01:44:27 (61.29 KB/s) - `scapy.py' saved [364,749/364,749]
```

Then simply launch Scapy using the Python interpreter as follows:

```bash
root@casper:/trash# python scapy.py
Welcome to Scapy (v1.1.1 / -)
>>>
```

If we want to get information about the configuration of the Scapy version we're using, simply type conf and press Enter to validate and execute the command:

```bash
>>> conf
Version = v1.1.1 / -
ASN1_default_codec = <ASN1Codec BER[1]>
AS_resolver = <__main__.AS_resolver_multi instance at 0x87f612c>
BTsocket = <class __main__.BluetoothL2CAPSocket at 0x8721bfc>
IPCountry_base = 'GeoIPCountry4Scapy.gz'
L2listen = <class __main__.L2ListenSocket at 0x8721aac>
L2socket = <class __main__.L2Socket at 0x8721a4c>
L3socket = <class __main__.L3PacketSocket at 0x8721a1c>
auto_fragment = 1
checkIPID = 0
checkIPaddr = 1
checkIPsrc = 1
check_TCPerror_seqack = 0
color_theme = <DefaultTheme>
countryLoc_base = 'countryLoc.csv'
debug_dissector = 0
debug_match = 0
ethertypes = </etc/ethertypes/ >
except_filter = ''
gnuplot_world = 'world.dat'
histfile = '/root/.scapy_history'
iface = 'ath0'
manufdb = </usr/share/wireshark/wireshark/manuf/ 00:06:A2 00:20:A7 00:06:A5 00:20:A6 00:05:46 00:20:A9 00:05:45 00:10:57 00:05:44 00:12:5B 00:10:55 00:0C:49 00:0C:48 00:0C:45 00:0C:44 00:0C:47 00:0C:46 00:0C:41 00:0C:40 00:0C:43 00:0C:42 00:40:93 00:40:92 00:40:91 00:40:90 00:40:97 00:10:50 00:40:95 00:40:94 00:40:99 00:40:98 47:54:43 00:C0:89 00:C0:88 00:40:9C:F5 ...................... 00:10:F4 00:11:D2 00:11:D3 00:11:D0 00:11:D1 00:11:D6 00:11:D7 00:11:D4 00:11:D5 00:11:D8 00:11:D9 AB-00-04-01-00-00/32 00:06:A9 00:05:49 00:11:DB 00:11:DC 00:11:DA 00:11:DF 00:11:DD 00:11:DE 00:0A:4F 00:0A:4E 00:0A:4D 00:0A:4C 00:0A:4B 00:0A:4A 00:06:A0 00:10:FC 00:10:FB 00:10:FA 00:10:FF 00:10:FE 00:10:FD>
mib = <MIB/ >
nmap_base = '/usr/share/nmap/nmap-os-fingerprints'
noenum = <Resolve []>
p0f_base = '/etc/p0f/p0f.fp'
padding = 1
prog = Version = v1.1.1 / -
display = 'display'
dot = 'dot'
hexedit = 'hexer'
pdfreader = 'acroread'
psreader = 'gv'
tcpdump = 'tcpdump'
tcpreplay = 'tcpreplay'
wireshark = 'wireshark'
promisc = 1
prompt = '>>> '
protocols = </etc/protocols/ pim ip ax_25 esp tcp ah ipv6_opts xtp ipv6_route igmp igp ddp etherip xns_idp ipv6_frag vrrp gre ipcomp encap ipv6 iso_tp4 sctp ipencap rsvp udp ggp hmp idpr_cmtp fc skip st icmp pup isis rdp l2tp ipv6_icmp egp ipip ipv6_nonxt eigrp idrp rspf ospf vmtp>
queso_base = '/etc/queso.conf'
resolve = <Resolve []>
route = Network Netmask Gateway Iface Output IP
127.0.0.0 255.0.0.0 0.0.0.0 lo 127.0.0.1
192.168.0.0 255.255.255.0 0.0.0.0 ath0 192.168.0.2
0.0.0.0 0.0.0.0 192.168.0.1 ath0 192.168.0.2

services_tcp = </etc/services-tcp/ kpop noclog svn cmip_man z3950 rootd ndtp gds_db ftps bacula_fd bgpsim isisd mysql bpdbm xdmcp rtcm_sc104 knetd systat mmcc enbd_cstatd daap radmin_port hylafax distmp3 hostmon snmp_trap isakmp dict ldap enbd_sstatd kshell tempo imaps pawserv afs3_errors x11_6 fax smux spamd krb_prop tproxy auth zebrasrv pop3 pop2 silc x11_5 ssh krbupdate afmbackup saft afbackup zope nntps loc_srv qotd msnp remotefs submission afs3_rmtsys prospero afs3_bos sysrqd webmin munin datametrics supfiledbg rfe kazaa ospf6d cvspserver venus imap2 imap3 afs3_fileserver webster gnutella_svc domain ripngd smsqp l2f afs3_vlserver ospfd ircs proofd rplay wipld omirr radius zebra kermit tfido sunrpc xtel pcrd zserv ftp rsync afs3_update imsp at_zis rmtcfg icpv2 daytime netnews afs3_volser kerberos4 finger x11 nfs irc nntp rpc2portmap kerberos_master rtelnet asp ulistserv moira_db ntp xmpp_server snpp netstat npmp_local ssmtp lotusnote sgi_cad csnet_ns fatserv ggz mrtd postgresql suucp rptp amanda distcc npmp_gui bgp afs3_kaserver vboxd csync2 discard re_mail_ck support tcpmux sip bpcd xinetd rtsp eklogin gopher socks echo ipp radius_acct webcache ipx ldaps rmiregistry conference bprd mdns sieve afpovertcp font_service swat netbios_dgm ripd snmp msp at_nbp frox amidxtape tinc iax https netbios_ns xpilot sane_port ninstall bacula_dir pop3s cfinger nessus remoteping moira_update shell clc_build_daemon xmpp_client courier skkserv codaauth2 openvpn kerberos afs3_prserver dircproxy telnet rje link netbios_ssn gnutella_rtr qmtp omniorb acr_nema x11_7 x11_4 pwdgen x11_2 x11_3 x11_1 gdomap tacacs mtp bacula_sd nextstep klogin www hkp login clearcase sa_msg_port nsca nqs at_echo mailq kamanda kerberos_adm vopied iprop kpasswd bpjava_msvc smtp xtelw sip_tls xtell afs3_callback ftp_data uucp_path wnn6 uucp supfilesrv whois telnets ospfapi unix_status ircd fido linuxconf microsoft_ds cfengine iso_tsap zope_ftp canna ms_sql_m prospero_np printer vnetd isdnlog nut exec ingreslock aol tacacs_ds amandaidx customs venus_se hostnames log_server nameserver ms_sql_s hmmp_ind supdup kx gpsd codasrv poppassd codasrv_se bootps sftp at_rtmp ftps_data binkp chargen time bgpd cmip_agent bootpc mon>

services_udp = </etc/services-udp/ noclog cmip_man z3950 rootd gds_db afs3_rmtsys bacula_fd mysql bpdbm xdmcp rtcm_sc104 x11_5 zephyr_hm daap radmin_port snmp_trap isakmp dict ldap ircs imaps pawserv afs3_errors x11_6 smux pop3 pop2 silc ssh afmbackup saft afbackup nntps loc_srv msnp isdnlog submission mandelspawn prospero afs3_bos sysrqd datametrics rfe kazaa cvspserver venus imap2 imap3 afs3_fileserver gnutella_svc domain smsqp l2f afs3_vlserver nut proofd rplay syslog omirr asp radius kermit sunrpc zserv ingreslock rsync imsp at_zis icpv2 daytime netwall kerberos4 x11 nfs irc rpc2portmap kerberos_master rtelnet aol ulistserv ntp xmpp_server snpp ggz lotusnote npmp_local fatserv hostmon rlp csnet_ns postgresql suucp rptp amanda route npmp_gui bgp afs3_kaserver vboxd timed discard re_mail_ck sgi_crsd gnutella_rtr passwd_server bpcd sgi_gcd rtsp gopher socks echo ipp radius_acct zephyr_clt ipx ldaps rmiregistry bprd mdns afpovertcp font_service netbios_dgm snmp msp at_nbp afs3_volser tinc iax https netbios_ns xpilot ninstall bacula_dir pop3s nessus svn xmpp_client codaauth2 openvpn kerberos afs3_prserver netbios_ssn sgi_cmsd qmtp omniorb acr_nema x11_7 x11_4 pwdgen x11_2 x11_3 x11_1 gdomap tacacs bacula_sd nextstep www hkp webster clearcase sa_msg_port nqs at_echo mailq kamanda vopied predict kpasswd bpjava_msvc distcc sip_tls zephyr_srv afs3_callback ntalk wnn6 fsp sip telnets mmcc microsoft_ds cfengine moira_ureg ms_sql_m prospero_np biff vnetd afs3_update who tacacs_ds customs venus_se tftp ms_sql_s hmmp_ind gpsd codasrv poppassd codasrv_se bootps at_rtmp chargen time cmip_agent bootpc talk mon>

session = ''
sniff_promisc = 1
stealth = 'not implemented'
verb = 2
warning_threshold = 5
wepkey = ''
```

Scapy also allows us to choose to view only part of the configuration if we want, here we want to display only information related to the routing table using the conf.route command:

```bash
>>> conf.route
Network Netmask Gateway Iface Output IP
127.0.0.0 255.0.0.0 0.0.0.0 lo 127.0.0.1
192.168.0.0 255.255.255.0 0.0.0.0 ath0 192.168.0.2
0.0.0.0 0.0.0.0 192.168.0.1 ath0 192.168.0.2
```

By default, this information is equivalent to the routing information delivered by the system command route:

```bash
root@casper:~# route
Kernel IP routing table
Destination Gateway Genmask Flags Metric Ref Use Iface
192.168.0.0 * 255.255.255.0 U 0 0 0 ath0
default 192.168.0.1 0.0.0.0 UG 0 0 0 ath0
```

To add an entry to the routing table we just viewed, use the conf.route.add command:

```bash
>>> conf.route.add(net="192.168.1.0/24",gw="192.168.0.1")
```

We verify that the routing entry has been added for the 192.168.1.0/24 network with the machine whose IP address is 192.168.0.1 as gateway:

```bash
>>> conf.route
Network Netmask Gateway Iface Output IP
127.0.0.0 255.0.0.0 0.0.0.0 lo 127.0.0.1
192.168.0.0 255.255.255.0 0.0.0.0 ath0 192.168.0.2
0.0.0.0 0.0.0.0 192.168.0.1 ath0 192.168.0.2
192.168.1.0 255.255.255.0 192.168.0.1 ath0 192.168.0.2
```

Now we check the system routing table for the second time:

```bash
root@casper:~# route
Kernel IP routing table
Destination Gateway Genmask Flags Metric Ref Use Iface
192.168.0.0 * 255.255.255.0 U 0 0 0 ath0
default 192.168.0.1 0.0.0.0 UG 0 0 0 ath0
```

The routing tables between the system and Scapy are different because Scapy has its own internal routing table.

We delete this entry with the conf.route.delt command, still for the 192.168.1.0/24 network with a machine acting as gateway whose IP address is 192.168.0.1:

```bash
>>> conf.route.delt(net="192.168.1.0/24",gw="192.168.0.1")
```

We verify that the corresponding routing entry has been removed from Scapy's internal routing table:

```bash
>>> conf.route
Network Netmask Gateway Iface Output IP
127.0.0.0 255.0.0.0 0.0.0.0 lo 127.0.0.1
192.168.0.0 255.255.255.0 0.0.0.0 ath0 192.168.0.2
0.0.0.0 0.0.0.0 192.168.0.1 ath0 192.168.0.2
```

Indeed, the entry is no longer present in the internal routing table; if we had wanted to add a route only to a particular machine, we could have used the following syntax (host= instead of net=) to specify that packets/frames destined for the machine whose IP address is 192.168.1.3 must be routed to the machine whose IP is 192.168.0.1, which thus acts as a gateway to access it:

```bash
>>> conf.route.add(host="192.168.1.3",gw="192.168.0.1")
```

We can also change the network interface with which we want to work by default using the conf.iface command:

```bash
>>> conf.iface='eth0'

>>> conf.iface
'eth0'

>>> conf.iface='ath0'

>>> conf.iface
'ath0'
```

The system is now fully functional, and we can move on to the packet/data frame manipulation functions offered by Scapy.

## Basic Usage

This part of the Scapy documentation proposes to give an overview of the basic commands available in Scapy. It also includes different examples of their use that will allow you to become familiar with the internal functioning of this tool in order to better understand its principles.

First, we list all the protocols supported by Scapy using the ls() command:

```bash
>>> ls()
ARP : ARP
ASN1_Packet : None
BOOTP : BOOTP
CookedLinux : cooked linux
DHCP : DHCP options
DNS : DNS
DNSQR : DNS Question Record
DNSRR : DNS Resource Record
Dot11 : 802.11
Dot11ATIM : 802.11 ATIM
Dot11AssoReq : 802.11 Association Request
Dot11AssoResp : 802.11 Association Response
Dot11Auth : 802.11 Authentication
Dot11Beacon : 802.11 Beacon
Dot11Deauth : 802.11 Deauthentication
Dot11Disas : 802.11 Disassociation
Dot11Elt : 802.11 Information Element
Dot11ProbeReq : 802.11 Probe Request
Dot11ProbeResp : 802.11 Probe Response
Dot11ReassoReq : 802.11 Reassociation Request
Dot11ReassoResp : 802.11 Reassociation Response
Dot11WEP : 802.11 WEP packet
Dot1Q : 802.1Q
Dot3 : 802.3
EAP : EAP
EAPOL : EAPOL
Ether : Ethernet
GPRS : GPRSdummy
GRE : GRE
HCI_ACL_Hdr : HCI ACL header
HCI_Hdr : HCI header
HSRP : HSRP
ICMP : ICMP
ICMPerror : ICMP in ICMP
IP : IP
IPerror : IP in ICMP
IPv6 : IPv6 not implemented here.
ISAKMP : ISAKMP
ISAKMP_class : None
ISAKMP_payload : ISAKMP payload
ISAKMP_payload_Hash : ISAKMP Hash
ISAKMP_payload_ID : ISAKMP Identification
ISAKMP_payload_KE : ISAKMP Key Exchange
ISAKMP_payload_Nonce : ISAKMP Nonce
ISAKMP_payload_Proposal : IKE proposal
ISAKMP_payload_SA : ISAKMP SA
ISAKMP_payload_Transform : IKE Transform
ISAKMP_payload_VendorID : ISAKMP Vendor ID
IrLAPCommand : IrDA Link Access Protocol Command
IrLAPHead : IrDA Link Access Protocol Header
IrLMP : IrDA Link Management Protocol
L2CAP_CmdHdr : L2CAP command header
L2CAP_CmdRej : L2CAP Command Rej
L2CAP_ConfReq : L2CAP Conf Req
L2CAP_ConfResp : L2CAP Conf Resp
L2CAP_ConnReq : L2CAP Conn Req
L2CAP_ConnResp : L2CAP Conn Resp
L2CAP_DisconnReq : L2CAP Disconn Req
L2CAP_DisconnResp : L2CAP Disconn Resp
L2CAP_Hdr : L2CAP header
L2CAP_InfoReq : L2CAP Info Req
L2CAP_InfoResp : L2CAP Info Resp
LLC : LLC
MGCP : MGCP
MobileIP : Mobile IP (RFC3344)
MobileIPRRP : Mobile IP Registration Reply (RFC3344)
MobileIPRRQ : Mobile IP Registration Request (RFC3344)
MobileIPTunnelData : Mobile IP Tunnel Data Message (RFC3519)
NBNSNodeStatusResponse : NBNS Node Status Response
NBNSNodeStatusResponseEnd : NBNS Node Status Response
NBNSNodeStatusResponseService : NBNS Node Status Response Service
NBNSQueryRequest : NBNS query request
NBNSQueryResponse : NBNS query response
NBNSQueryResponseNegative : NBNS query response (negative)
NBNSRequest : NBNS request
NBNSWackResponse : NBNS Wait for Acknowledgement Response
NBTDatagram : NBT Datagram Packet
NBTSession : NBT Session Packet
NTP : NTP
NetBIOS_DS : NetBIOS datagram service
NetflowHeader : Netflow Header
NetflowHeaderV1 : Netflow Header V1
NetflowRecordV1 : Netflow Record
NoPayload : None
PPP : PPP Link Layer
PPPoE : PPP over Ethernet
PPPoED : PPP over Ethernet Discovery
Packet : None
Padding : Padding
PrismHeader : Prism header
RIP : RIP header
RIPEntry : RIP entry
RTP : RTP
RadioTap : RadioTap dummy
Radius : Radius
Raw : Raw
SMBMailSlot : SMB Mail Slot Protocol
SMBNegociate_Protocol_Request_Header : SMBNegociate Protocol Request Header
SMBNegociate_Protocol_Request_Tail : SMB Negociate Protocol Request Tail
SMBNegociate_Protocol_Response_Advanced_Security : SMBNegociate Protocol Response Advanced Security
SMBNegociate_Protocol_Response_No_Security : SMBNegociate Protocol Response No Security
SMBNegociate_Protocol_Response_No_Security_No_Key : None
SMBNetlogon_Protocol_Response_Header : SMBNetlogon Protocol Response Header
SMBNetlogon_Protocol_Response_Tail_LM20 : SMB Netlogon Protocol Response Tail LM20
SMBNetlogon_Protocol_Response_Tail_SAM : SMB Netlogon Protocol Response Tail SAM
SMBSession_Setup_AndX_Request : Session Setup AndX Request
SMBSession_Setup_AndX_Response : Session Setup AndX Response
SNAP : SNAP
SNMP : None
SNMPbulk : None
SNMPget : None
SNMPinform : None
SNMPnext : None
SNMPresponse : None
SNMPset : None
SNMPtrapv1 : None
SNMPtrapv2 : None
SNMPvarbind : None
STP : Spanning Tree Protocol
SebekHead : Sebek header
SebekV1 : Sebek v1
SebekV2 : Sebek v3
SebekV2Sock : Sebek v2 socket
SebekV3 : Sebek v3
SebekV3Sock : Sebek v2 socket
Skinny : Skinny
TCP : TCP
TCPerror : TCP in ICMP
TFTP : TFTP opcode
TFTP_ACK : TFTP Ack
TFTP_DATA : TFTP Data
TFTP_ERROR : TFTP Error
TFTP_OACK : TFTP Option Ack
TFTP_Option : None
TFTP_Options : None
TFTP_RRQ : TFTP Read Request
TFTP_WRQ : TFTP Write Request
UDP : UDP
UDPerror : UDP in ICMP
_IPv6OptionHeader : IPv6 not implemented here.
```

As we can see, the number of supported protocols is quite substantial; each of these protocols has its own specifications that we can list using the ls() command. Here are the details of the ICMP protocol:

```bash
>>> ls(ICMP)
type : ByteEnumField = (8)
code : ByteField = (0)
chksum : XShortField = (None)
id : XShortField = (0)
seq : XShortField = (0)
```

To access the full list of commands/functions available in Scapy, we must use the lsc() command:

```bash
>>> lsc()
sr : Send and receive packets at layer 3
sr1 : Send packets at layer 3 and return only the first answer
srp : Send and receive packets at layer 2
srp1 : Send and receive packets at layer 2 and return only the first answer
srloop : Send a packet at layer 3 in loop and print the answer each time
srploop : Send a packet at layer 2 in loop and print the answer each time
sniff : Sniff packets
p0f : Passive OS fingerprinting: which OS emitted this TCP SYN?
arpcachepoison : Poison target's cache with (your MAC,victim's IP) couple
send : Send packets at layer 3
sendp : Send packets at layer 2
traceroute : Instant TCP traceroute
arping : Send ARP who-has requests to determine which hosts are up
ls : List available layers, or infos on a given layer
lsc : List user commands
queso : Queso OS fingerprinting
nmap_fp : nmap fingerprinting
report_ports : portscan a target and output a LaTeX table
dyndns_add : Send a DNS add message to a nameserver for "name" to have a new "rdata"
dyndns_del : Send a DNS delete message to a nameserver for "name"
is_promisc : Try to guess if target is in Promisc mode. The target is provided by its ip.
promiscping : Send ARP who-has requests to determine which hosts are in promiscuous mode
```

To display the documentation for a particular function, simply add the .**doc** extension behind the command (without the final ()); so to display the documentation for the arping() command, it will be necessary to type the following syntax:

```bash
>>> arping.__doc__
'Send ARP who-has requests to determine which hosts are up\narping(net, [cache=0,] [iface=conf.iface,] [verbose=conf.verb]) -> None\nSet cache=True if you want arping to modify internal ARP-Cache'
```

It's also possible to get the same documentation result but with a bit more layout using the lsc() command that allowed us to list the available commands in Scapy earlier; just add the name of the command in parentheses that you want documentation for:

```bash
>>> lsc(arping)
Send ARP who-has requests to determine which hosts are up
arping(net, [cache=0,] [iface=conf.iface,] [verbose=conf.verb]) -> None
Set cache=True if you want arping to modify internal ARP-Cache
```

## Data Capturing

Scapy can function like a network traffic analyzer to capture data for later viewing. This part of the documentation offers different examples of captures as well as the many ways available internally to view the results of these captures.

We display the documentation for the sniff() function using the lsc() function:

```bash
>>> lsc(sniff)
Sniff packets
sniff([count=0,] [prn=None,] [store=1,] [offline=None,] [lfilter=None,] + L2ListenSocket args) -> list of packets

count: number of packets to capture. 0 means infinity
store: wether to store sniffed packets or discard them
prn: function to apply to each packet. If something is returned,
it is displayed. Ex:
ex: prn = lambda x: x.summary()
lfilter: python function applied to each packet to determine
if further action may be done
ex: lfilter = lambda x: x.haslayer(Padding)
offline: pcap file to read packets from, instead of sniffing them
timeout: stop sniffing after a given time (default: None)
L2socket: use the provided L2socket
```

Now we launch a sniff on all UDP protocol traffic for the machine with IP address 192.168.0.2 with a maximum of 30 packets collected (using count):

```bash
>>> sniff(filter="udp and host 192.168.0.2", count=30)
<Sniffed: UDP:30 ICMP:0 TCP:0 Other:0>
```

The 30 packets have been collected for the UDP protocol and the machine whose IP address is 192.168.0.2; we can view the results related to this capture by assigning the records to the variable sn via the variable \_ in the following way:

```bash
>>> sn=_
```

If we want to view all these records contained in the variable sn, we need to add .nsummary() after the name of the variable that we chose when assigning the results of a function (here the sniff function and the sn variable):

```bash
>>> sn.nsummary()
0000 Ether / IP / UDP 192.168.0.1:netbios_dgm > 192.168.0.255:netbios_dgm / NBTDatagram / Raw
0001 Ether / IP / UDP 192.168.0.1:netbios_dgm > 192.168.0.255:netbios_dgm / NBTDatagram / Raw
0002 Ether / IP / UDP 192.168.0.1:netbios_dgm > 192.168.0.255:netbios_dgm / NBTDatagram / Raw
0003 Ether / IP / UDP 192.168.0.1:netbios_dgm > 192.168.0.255:netbios_dgm / NBTDatagram / Raw
0004 Ether / IP / UDP 192.168.0.1:netbios_dgm > 192.168.0.255:netbios_dgm / NBTDatagram / Raw
0005 Ether / IP / UDP 192.168.0.1:netbios_dgm > 192.168.0.255:netbios_dgm / NBTDatagram / Raw
0006 Ether / IP / UDP 192.168.0.1:netbios_ns > 192.168.0.255:netbios_ns / NBNSQueryRequest
0007 Ether / IP / UDP 192.168.0.1:netbios_ns > 192.168.0.255:netbios_ns / NBNSQueryRequest
0008 Ether / IP / UDP 192.168.0.1:netbios_ns > 192.168.0.255:netbios_ns / NBNSQueryRequest
0009 Ether / IP / UDP 192.168.0.1:netbios_dgm > 192.168.0.255:netbios_dgm / NBTDatagram / Raw
0010 Ether / IP / UDP 192.168.0.1:1900 > 239.255.255.250:1900 / Raw
0011 Ether / IP / UDP 192.168.0.1:1900 > 239.255.255.250:1900 / Raw
0012 Ether / IP / UDP 192.168.0.1:1900 > 239.255.255.250:1900 / Raw
0013 Ether / IP / UDP 192.168.0.1:1900 > 239.255.255.250:1900 / Raw
0014 Ether / IP / UDP 192.168.0.1:1900 > 239.255.255.250:1900 / Raw
0015 Ether / IP / UDP 192.168.0.1:1900 > 239.255.255.250:1900 / Raw
0016 Ether / IP / UDP 192.168.0.1:1900 > 239.255.255.250:1900 / Raw
0017 Ether / IP / UDP 192.168.0.1:1900 > 239.255.255.250:1900 / Raw
0018 Ether / IP / UDP 192.168.0.1:1900 > 239.255.255.250:1900 / Raw
0019 Ether / IP / UDP 192.168.0.1:1900 > 239.255.255.250:1900 / Raw
0020 Ether / IP / UDP 192.168.0.1:1900 > 239.255.255.250:1900 / Raw
0021 Ether / IP / UDP 192.168.0.1:1900 > 239.255.255.250:1900 / Raw
0022 Ether / IP / UDP 192.168.0.1:1900 > 239.255.255.250:1900 / Raw
0023 Ether / IP / UDP 192.168.0.1:1900 > 239.255.255.250:1900 / Raw
0024 Ether / IP / UDP 192.168.0.1:1900 > 239.255.255.250:1900 / Raw
0025 Ether / IP / UDP 192.168.0.1:1900 > 239.255.255.250:1900 / Raw
0026 Ether / IP / UDP 192.168.0.1:1900 > 239.255.255.250:1900 / Raw
0027 Ether / IP / UDP 192.168.0.1:1900 > 239.255.255.250:1900 / Raw
0028 Ether / IP / UDP 192.168.0.1:1900 > 239.255.255.250:1900 / Raw
0029 Ether / IP / UDP 192.168.0.1:1900 > 239.255.255.250:1900 / Raw
```

Sn is considered here as an object in its own right on which we can apply different operations via appropriate functions like nsummary() here; to view in detail the first record contained in the variable sn, we need to add between [] the number of the record that we want to view in the traffic we just captured and this just behind the name of the variable:

```bash
>>> sn[0]
<Ether dst=ff:ff:ff:ff:ff:ff src=00:0f:b5:0e:89:4a type=0x800 |<IP version=4L ihl=5L tos=0x0 len=229 id=2852 flags= frag=0L ttl=128 proto=udp chksum=0xac93 src=192.168.0.1 dst=192.168.0.255 options='' |<UDP sport=netbios_dgm dport=netbios_dgm len=209 chksum=0xa7d |<NBTDatagram Type=17 Flags=14 ID=32985 SourceIP=192.168.0.1 SourcePort=138 Length=187 Offset=0 SourceName='HEIDI ' SUFFIX1=file server service NULL=0 DestinationName='ARBEITSGRUPPE ' SUFFIX2=browser election service NULL=0 |<Raw load='\xffSMB%\x00\x00\x00\x00\x00\x00\x00\x00
\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00
\x00\x00\x00\x00\x00\x11\x00\x00!\x00\x00\x00\x00\x00\x00
\x00\x00\x00\xe8\x03\x00\x00\x00\x00\x00\x00\x00\x00!\x00V
\x00\x03\x00\x01\x00\x00\x00\x02\x002\x00\\MAILSLOT\\BROWSE
\x00\x0f\x00\x80\xfc\n\x00HEIDI\x00\x00\x00\x00\x00d\x00\x00\
x00\x00\x00\x05\x01\x03\x10\x05\x00\x0f\x01U\xaa\x00'|>
;>>
```

The last packet record (the 30th in our case) has the id 29 since the first record has an id equal to 0:

```bash
>>> sn[29]
<Ether dst=01:00:5e:7f:ff:fa src=00:0f:b5:0e:89:4a type=0x800 |<IP version=4L ihl=5L tos=0x0 len=373 id=2938 flags= frag=0L ttl=1 proto=udp chksum=0xfc5a src=192.168.0.1 dst=239.255.255.250 options='' |<UDP sport=1900 dport=1900 len=353 chksum=0x23e2 |<Raw load='NOTIFY * HTTP/1.1\r\nHost:239.255.255.250:1900\r\nNT:upnp:rootdevice\r
\nNTS:ssdp:alive\r\nLocation:http:[]
/udhisapi.dll?content=uuid:03cda088-f88e-4a2b-96a6-
0bbe0741838e\r\nUSN:uuid:03cda088-f88e-4a2b-96a6-
0bbe0741838e::upnp:rootdevice\r\nCache-Control:max-age=1800
\r\nServer:Microsoft-Windows-NT/5.1 UPnP/1.0 UPnP-Device-Host
/1.0\r\n\r\n' |>>>>
```

We verify that record id 30 does not exist in the sn variable:

```bash
>>> sn[30]
Traceback (most recent call last):
File "<console>", line 1, in <module>
File "/trash/scapy.py", line 2450, in __getitem__
return self.res.__getitem__(item)
IndexError: list index out of range
```

With Scapy, we can also capture network traffic and directly display the results while specifying the network interface on which we want to perform this operation if of course it was different from the one defined in the conf.iface variable; this is done as follows:

```bash
>>> sniff(iface="ath0", prn=lambda x: x.summary(),count=30)
Ether / IP / TCP 192.168.0.2:50414 > 213.251.178.32:www PA / Raw
Ether / IP / TCP 213.251.178.32:www > 192.168.0.2:50414 A / Raw
Ether / IP / TCP 192.168.0.2:50414 > 213.251.178.32:www A
Ether / IP / TCP 213.251.178.32:www > 192.168.0.2:50414 A / Raw
Ether / IP / TCP 192.168.0.2:50414 > 213.251.178.32:www A
Ether / IP / TCP 213.251.178.32:www > 192.168.0.2:50414 A / Raw
Ether / IP / TCP 192.168.0.2:50414 > 213.251.178.32:www A
Ether / IP / TCP 213.251.178.32:www > 192.168.0.2:50414 A / Raw
Ether / IP / TCP 192.168.0.2:50414 > 213.251.178.32:www A
Ether / IP / TCP 213.251.178.32:www > 192.168.0.2:50414 A / Raw
Ether / IP / TCP 192.168.0.2:50414 > 213.251.178.32:www A
Ether / IP / TCP 213.251.178.32:www > 192.168.0.2:50414 A / Raw
Ether / IP / TCP 192.168.0.2:50414 > 213.251.178.32:www A
Ether / IP / TCP 213.251.178.32:www > 192.168.0.2:50414 A / Raw
Ether / IP / TCP 192.168.0.2:50414 > 213.251.178.32:www A
Ether / IP / TCP 213.251.178.32:www > 192.168.0.2:50414 A / Raw
Ether / IP / TCP 192.168.0.2:50414 > 213.251.178.32:www A
Ether / IP / TCP 213.251.178.32:www > 192.168.0.2:50414 A / Raw
Ether / IP / TCP 192.168.0.2:50414 > 213.251.178.32:www A
Ether / IP / TCP 213.251.178.32:www > 192.168.0.2:50414 A / Raw
Ether / IP / TCP 192.168.0.2:50414 > 213.251.178.32:www A
Ether / IP / TCP 213.251.178.32:www > 192.168.0.2:50414 A / Raw
Ether / IP / TCP 192.168.0.2:50414 > 213.251.178.32:www A
Ether / IP / TCP 213.251.178.32:www > 192.168.0.2:50414 A / Raw
Ether / IP / TCP 192.168.0.2:50414 > 213.251.178.32:www A
Ether / IP / TCP 213.251.178.32:www > 192.168.0.2:50414 A / Raw
Ether / IP / TCP 192.168.0.2:50414 > 213.251.178.32:www A
Ether / IP / TCP 213.251.178.32:www > 192.168.0.2:50414 A / Raw
Ether / IP / TCP 192.168.0.2:50414 > 213.251.178.32:www A
Ether / IP / TCP 213.251.178.32:www > 192.168.0.2:50414 A / Raw
<Sniffed: UDP:0 ICMP:0 TCP:30 Other:0>
```

We can also view the records in a completely different way by replacing the x.summary() attribute with the x.show() attribute, which this time will allow us to get much more intrinsic details about the traffic collected during the network capture; since the output can be quite long, we'll limit the number of packets we'll capture to 3 with count:

```bash
>>> sniff(iface="ath0", prn=lambda x: x.show(),count=3)

###[ Ethernet ]###
dst= 00:0f:b5:0e:89:4a
src= 00:15:6d:53:1e:87
type= 0x800
###[ IP ]###
version= 4L
ihl= 5L
tos= 0x0
len= 678
id= 37167
flags= DF
frag= 0L
ttl= 64
proto= tcp
chksum= 0x5e5c
src= 192.168.0.2
dst= 213.251.178.32
options= ''
###[ TCP ]###
sport= 34554
dport= www
seq= 1842608767
ack= 4117121824L
dataofs= 8L
reserved= 0L
flags= PA
window= 2003
chksum= 0x1c8f
urgptr= 0
options= [('NOP', None), ('NOP', None), ('Timestamp', (1035373, 3849016752L))]
###[ Raw ]###
load= 'GET / HTTP/1.1\r\nHost: www.secuobs.com\r\nUser-Agent: Mozilla/5.0 (X11; U; Linux i686; en-US; rv:1.8.1.6) Gecko/20061201 Firefox/2.0.0.6 (Ubuntu-feisty)\r\nAccept: text/xml,application/xml,application/xhtml+xml,text/html;
q=0.9,text/plain;q=0.8,image/png,*/*;q=0.5\r\nAccept-Language: fr-fr,en-us;q=0.7,en;q=0.3\r\nAccept-Encoding: gzip,deflate\r
\nAccept-Charset: ISO-8859-1,utf-8;q=0.7,*;q=0.7\r\nKeep-Alive: 300\r\nConnection: keep-alive\r\nCookie: __utma=50890539.1734225999.1190320568.1190636274.1190638025
.62; __utmb=50890539; __utmz=50890539.1190339822.4.2.utmccn=
(organic)|utmcsr=google|utmctr=captive|utmcmd=organic; __
utmc=50890539\r\n\r\n'
###[ Ethernet ]###
dst= 00:15:6d:53:1e:87
src= 00:0f:b5:0e:89:4a
type= 0x800
###[ IP ]###
version= 4L
ihl= 5L
tos= 0x0
len= 1492
id= 38106
flags= DF
frag= 0L
ttl= 52
proto= tcp
chksum= 0x6383
src= 213.251.178.32
dst= 192.168.0.2
options= ''
###[ TCP ]###
sport= www
dport= 34554
seq= 4117121824L
ack= 1842609393
dataofs= 8L
reserved= 0L
flags= A
window= 15024
chksum= 0x3812
urgptr= 0
options= [('NOP', None), ('NOP', None), ('Timestamp', (3849021332L, 1035373))]
###[ Raw ]###
load= 'HTTP/1.1 200 OK\r\nDate: Mon, 24 Sep 2007 13:05:42 GMT\r\nServer: Apache\r\nKeep-Alive: timeout=7200\r\nConnection: Keep-Alive\r\nTransfer-Encoding: chunked\r\nContent-Type: text/html\r\n\r\ncb3\r\n<html>\n<head>\n <title>SecuObs.com -\nL'observatoire de la sécurité internet</title>\n <meta name="description"\n content="L'observatoire de la securite internet - Site d'informations professionnelles francophone sur la sécurité informatique">\n <meta\n content="text/html; charset=ISO-8859-1"\n http-equiv="content-type">\n <title></title>\n'
###[ Ethernet ]###
dst= 00:0f:b5:0e:89:4a
src= 00:15:6d:53:1e:87
type= 0x800
###[ IP ]###
version= 4L
ihl= 5L
tos= 0x0
len= 52
id= 37168
flags= DF
frag= 0L
ttl= 64
proto= tcp
chksum= 0x60cd
src= 192.168.0.2
dst= 213.251.178.32
options= ''
###[ TCP ]###
sport= 34554
dport= www
seq= 1842609393
ack= 4117123264L
dataofs= 8L
reserved= 0L
flags= A
window= 1963
chksum= 0xd48a
urgptr= 0
options= [('NOP', None), ('NOP', None), ('Timestamp', (1035386, 3849021332L))]
<Sniffed: UDP:0 ICMP:0 TCP:3 Other:0>
```

In the same way, we can capture network traffic for a given port or several given ports; here we capture the traffic (5 records) related to port 22 (usually equivalent to encrypted remote connection traffic like SSH) and port 80, which is usually synonymous with traffic to a web server like Apache (or IIS):

```bash
>>> sniff(filter="tcp and ( port 22 or port 80 )", prn=lambda x: x.sprintf("IP Source: %IP.src% ; Port Source: %TCP.sport% ---> IP Destination: %IP.dst% ; Port Destination: %TCP.dport% ; TCP Flags: %TCP.flags% ; Payload: %TCP.payload%"), count=5)
IP Source: 127.0.0.1 ; Port Source: 33917 ---> IP Destination: 127.0.0.1 ; Port Destination: ssh ; TCP Flags: FA ; Payload:
IP Source: 127.0.0.1 ; Port Source: 33917 ---> IP Destination: 127.0.0.1 ; Port Destination: ssh ; TCP Flags: FA ; Payload:
IP Source: 127.0.0.1 ; Port Source: ssh ---> IP Destination: 127.0.0.1 ; Port Destination: 33917 ; TCP Flags: R ; Payload:
IP Source: 127.0.0.1 ; Port Source: ssh ---> IP Destination: 127.0.0.1 ; Port Destination: 33917 ; TCP Flags: R ; Payload:
IP Source: 127.0.0.1 ; Port Source: 33917 ---> IP Destination: 127.0.0.1 ; Port Destination: ssh ; TCP Flags: FA ; Payload:
<Sniffed: UDP:0 TCP:5 ICMP:0 Other:0>
```

The x.sprintf attribute in the example above allows you to format the output results of Scapy's functions according to your preferences using pre-defined variables such as %IP.src% for the source IP address, %IP.dst% for the destination IP address, %TCP.sport% for the TCP source port, %TCP.dport% for the TCP destination port as well as %TCP.flags% for the options activated in the TCP flags as well as %TCP.payload% for the useful payload related to the record being viewed.

The use of these variables and their manipulation accentuate the object-oriented side of the Scapy utility; for example, we can display the information related to a record in the following way:

```bash
>>> IP(dst="192.168.0.2")
<IP dst=192.168.0.2 |>
>>> a.sprintf("IP Source %IP.src% ; IP Destination %IP.dst% ; TTL %IP.ttl%")
'IP Source 192.168.0.2 ; IP Destination 192.168.0.2 ; TTL 64'
```

Note that even if we didn't specify values for the ttl, Scapy takes a default value equivalent to 64 that we find when displaying via the %IP.ttl% variable.

The sequence of commands that will follow presents itself in the form of a capture of traffic equivalent to ports 22 (SSH) and/or 80 (WWW) for a total of 5 records related to the TCP transport protocol; the collected records will be stored in a variable sn for which we can display the first record with sn[0]:

```bash
>>> sniff(filter="tcp and ( port 22 or port 80 )", count=5)
<Sniffed: UDP:0 ICMP:0 TCP:5 Other:0>
>>> sn=_
>>> sn[0]
<Ether dst=00:00:00:00:00:00 src=00:00:00:00:00:00 type=0x800 |<IP version=4L ihl=5L tos=0x0 len=60 id=45016 flags=DF frag=0L ttl=64 proto=tcp chksum=0x8ce1 src=127.0.0.1 dst=127.0.0.1 options='' |<TCP sport=56283 dport=ssh seq=2821237960L ack=0 dataofs=10L reserved=0L flags=S window=32792 chksum=0x5b8a urgptr=0 options=[('MSS', 16396), ('SAckOK', ''), ('Timestamp', (1180432, 0)), ('NOP', None), ('WScale', 5)] |>>>
```

This sequence of commands can also be written in an even more simplified way as follows:

```bash
>>> sn=sniff(filter="tcp and ( port 22 or port 80 )", count=5)
>>> sn[0]
<Ether dst=00:00:00:00:00:00 src=00:00:00:00:00:00 type=0x800 |<IP version=4L ihl=5L tos=0x0 len=60 id=35250 flags=DF frag=0L ttl=64 proto=tcp chksum=0xb307 src=127.0.0.1 dst=127.0.0.1 options='' |<TCP sport=56788 dport=ssh seq=2937645191L ack=0 dataofs=10L reserved=0L flags=S window=32792 chksum=0xb8bc urgptr=0 options=[('MSS', 16396), ('SAckOK', ''), ('Timestamp', (1204533, 0)), ('NOP', None), ('WScale', 5)] |>>>
```

## Traceroute and 2D/3D Visualization

Through different examples of traceroutes performed using Scapy, we'll discover the graphical functionalities of Scapy to generate two and three-dimensional graphs from the results of these traceroutes. You'll also see the different ways you can export these results for visualization.

We now display the information related to Scapy's traceroute command:

```bash
>>> lsc(traceroute)
Instant TCP traceroute
traceroute(target, [maxttl=30,] [dport=80,] [sport=80,] [verbose=conf.verb]) -> None
```

We then perform a simple traceroute to define the path taken by network traffic to go from the test machine to the IP address associated at the DNS level to the FQDN (Full Qualified Domain Name) and this with a latency time equal to 10 jiffies (number of clock periods); recall that TCP/IP networks are fragmentation and packet-switched networks, so it's normal to find a different packet routing result between two traceroutes launched at different intervals:

```bash
>>> traceroute(["www.google.com"],maxttl=10)
Begin emission:
****Finished to send 10 packets.***

Received 7 packets, got 7 answers, remaining 3 packets
64.233.183.104:tcp80
1 192.168.0.1 11
2 192.168.2.1 11
3 192.168.1.1 11
4 81.62.144.1 11
7 195.186.0.234 11
9 72.14.198.57 11
10 64.233.174.34 11
(<Traceroute: UDP:0 ICMP:7 TCP:0 Other:0>, <Unanswered: UDP:0 ICMP:0 TCP:3 Other:0>)
```

We then perform a multiple traceroute to the FQDNs "www.free.fr" and "www.exoscan.net" with a ttl still equal to 10; the results are stored in the variable sn for records that had a response while those without a response are stored in the variable unans:

```bash
>>> sn,unans=traceroute(["www.free.fr","www.exoscan.net"],maxttl=10)
Begin emission:
****************Finished to send 20 packets.

Received 16 packets, got 16 answers, remaining 4 packets
212.27.48.10:tcp80 213.186.41.29:tcp80
1 192.168.0.1 11 192.168.0.1 11
2 62.4.16.248 11 62.4.16.248 11
3 62.4.16.7 11 62.4.16.7 11
4 194.79.131.146 11 194.79.131.146 11
5 213.228.3.225 11 194.79.131.129 11
6 - 213.186.32.141 11
7 212.27.57.213 11 213.186.32.30 11
8 212.27.50.6 11 -
9 212.27.48.10 SA -
10 212.27.48.10 SA -
```

If we want to display only the no-response records, we can use the display() attribute behind the unans variable as follows:

```bash
>>> unans.display()
0000 IP / TCP 192.168.0.2:16324 > 213.186.41.29:www S
0001 IP / TCP 192.168.0.2:9718 > 213.186.41.29:www S
0002 IP / TCP 192.168.0.2:17318 > 213.186.41.29:www S
0003 IP / TCP 192.168.0.2:15117 > 212.27.48.10:www S
```

While to display only the records representing a response, we'll use the same display() attribute but this time on the sn variable as follows:

```bash
>>> sn.display()
212.27.48.10:tcp80 213.186.41.29:tcp80
1 192.168.0.1 11 192.168.0.1 11
2 62.4.16.248 11 62.4.16.248 11
3 62.4.16.7 11 62.4.16.7 11
4 194.79.131.146 11 194.79.131.146 11
5 213.228.3.225 11 194.79.131.129 11
6 - 213.186.32.141 11
7 212.27.57.213 11 213.186.32.30 11
8 212.27.50.6 11 -
9 212.27.48.10 SA -
10 212.27.48.10 SA -
```

For this time, we perform a multiple traceroute to the following FQDNs "www.google.com", "www.aol.com", "www.free.fr" and "www.secuobs.com" with a ttl still equal to 10:

```bash
>>> sn,unans=traceroute(["www.google.com","www.aol.com",
"www.free.fr","www.secuobs.com"],maxttl=10)
Begin emission:
************************
Finished to send 40 packets.

Received 24 packets, got 24 answers, remaining 16 packets
212.27.48.10:tcp80 213.251.178.32:tcp80 64.12.89.12:tcp80 64.233.183.104:tcp80
1 192.168.0.1 11 192.168.0.1 11 192.168.0.1 11 -
2 192.168.2.1 11 - 192.168.2.1 11 192.168.2.1 11
4 81.62.144.1 11 81.62.144.1 11 81.62.144.1 11 -
5 195.186.123.1 11 195.186.123.1 11 195.186.123.1 11 195.186.123.1 11
6 - - 195.186.0.229 11 195.186.0.229 11
8 138.187.129.45 11 138.187.129.45 11 138.187.129.74 11 -
9 138.187.130.73 11 138.187.130.73 11 194.42.48.4 11 72.14.198.57 11
10 - - 212.74.84.242 11 64.233.174.34 11
```

The multiple traceroute results are stored in the sn variable for records representing a response; we can just as well display the content of the sn variable in the following way using the show() attribute:

```bash
>>> sn.show()
212.27.48.10:tcp80 213.251.178.32:tcp80 64.12.89.12:tcp80 64.233.183.104:tcp80
1 192.168.0.1 11 192.168.0.1 11 192.168.0.1 11 -
2 192.168.2.1 11 - 192.168.2.1 11 192.168.2.1 11
4 81.62.144.1 11 81.62.144.1 11 81.62.144.1 11 -
5 195.186.123.1 11 195.186.123.1 11 195.186.123.1 11 195.186.123.1 11
6 - - 195.186.0.229 11 195.186.0.229 11
8 138.187.129.45 11 138.187.129.45 11 138.187.129.74 11 -
9 138.187.130.73 11 138.187.130.73 11 194.42.48.4 11 72.14.198.57 11
10 - - 212.74.84.242 11 64.233.174.34 11
```

We can similarly view the records contained in the sn variable by placing the nsummary() attribute instead of the show() attribute for the following result:

```bash
>>> sn.nsummary()
0000 IP / TCP 192.168.0.2:56554 > 64.12.89.12:www S ==> IP / ICMP 192.168.0.1 > 192.168.0.2 time-exceeded 0 / IPerror / TCPerror
0001 IP / TCP 192.168.0.2:38001 > 212.27.48.10:www S ==> IP / ICMP 192.168.0.1 > 192.168.0.2 time-exceeded 0 / IPerror / TCPerror
0002 IP / TCP 192.168.0.2:33842 > 213.251.178.32:www S ==> IP / ICMP 192.168.0.1 > 192.168.0.2 time-exceeded 0 / IPerror / TCPerror
0003 IP / TCP 192.168.0.2:55277 > 64.233.183.104:www S ==> IP / ICMP 192.168.2.1 > 192.168.0.2 time-exceeded 0 / IPerror / TCPerror
0004 IP / TCP 192.168.0.2:62887 > 64.12.89.12:www S ==> IP / ICMP 192.168.2.1 > 192.168.0.2 time-exceeded 0 / IPerror / TCPerror
0005 IP / TCP 192.168.0.2:amanda > 212.27.48.10:www S ==> IP / ICMP 192.168.2.1 > 192.168.0.2 time-exceeded 0 / IPerror / TCPerror
0006 IP / TCP 192.168.0.2:55421 > 64.12.89.12:www S ==> IP / ICMP 81.62.144.1 > 192.168.0.2 time-exceeded 0 / IPerror / TCPerror
0007 IP / TCP 192.168.0.2:46998 > 212.27.48.10:www S ==> IP / ICMP 81.62.144.1 > 192.168.0.2 time-exceeded 0 / IPerror / TCPerror
0008 IP / TCP 192.168.0.2:34599 > 213.251.178.32:www S ==> IP / ICMP 81.62.144.1 > 192.168.0.2 time-exceeded 0 / IPerror / TCPerror
0009 IP / TCP 192.168.0.2:58715 > 64.233.183.104:www S ==> IP / ICMP 195.186.123.1 > 192.168.0.2 time-exceeded 0 / IPerror / TCPerror
0010 IP / TCP 192.168.0.2:9093 > 64.12.89.12:www S ==> IP / ICMP 195.186.123.1 > 192.168.0.2 time-exceeded 0 / IPerror / TCPerror
0011 IP / TCP 192.168.0.2:14089 > 212.27.48.10:www S ==> IP / ICMP 195.186.123.1 > 192.168.0.2 time-exceeded 0 / IPerror / TCPerror
0012 IP / TCP 192.168.0.2:2258 > 213.251.178.32:www S ==> IP / ICMP 195.186.123.1 > 192.168.0.2 time-exceeded 0 / IPerror / TCPerror
0013 IP / TCP 192.168.0.2:60402 > 64.233.183.104:www S ==> IP / ICMP 195.186.0.229 > 192.168.0.2 time-exceeded 0 / IPerror / TCPerror
0014 IP / TCP 192.168.0.2:32073 > 64.12.89.12:www S ==> IP / ICMP 195.186.0.229 > 192.168.0.2 time-exceeded 0 / IPerror / TCPerror
0015 IP / TCP 192.168.0.2:18250 > 64.12.89.12:www S ==> IP / ICMP 138.187.129.74 > 192.168.0.2 time-exceeded 0 / IPerror / TCPerror
0016 IP / TCP 192.168.0.2:13695 > 212.27.48.10:www S ==> IP / ICMP 138.187.129.45 > 192.168.0.2 time-exceeded 0 / IPerror / TCPerror
0017 IP / TCP 192.168.0.2:19761 > 213.251.178.32:www S ==> IP / ICMP 138.187.129.45 > 192.168.0.2 time-exceeded 0 / IPerror / TCPerror
0018 IP / TCP 192.168.0.2:23264 > 64.233.183.104:www S ==> IP / ICMP 72.14.198.57 > 192.168.0.2 time-exceeded 0 / IPerror / TCPerror
0019 IP / TCP 192.168.0.2:32236 > 64.12.89.12:www S ==> IP / ICMP 194.42.48.4 > 192.168.0.2 time-exceeded 0 / IPerror / TCPerror
0020 IP / TCP 192.168.0.2:20611 > 212.27.48.10:www S ==> IP / ICMP 138.187.130.73 > 192.168.0.2 time-exceeded 0 / IPerror / TCPerror
0021 IP / TCP 192.168.0.2:40775 > 213.251.178.32:www S ==> IP / ICMP 138.187.130.73 > 192.168.0.2 time-exceeded 0 / IPerror / TCPerror
0022 IP / TCP 192.168.0.2:57456 > 64.233.183.104:www S ==> IP / ICMP 64.233.174.34 > 192.168.0.2 time-exceeded 0 / IPerror / TCPerror / Padding
0023 IP / TCP 192.168.0.2:56892 > 64.12.89.12:www S ==> IP / ICMP 212.74.84.242 > 192.168.0.2 time-exceeded 0 / IPerror / TCPerror
```

Visualization can also be done this way with the summary() attribute still in place of the show() attribute:

```bash
>>> sn.summary()
IP / TCP 192.168.0.2:56554 > 64.12.89.12:www S ==> IP / ICMP 192.168.0.1 > 192.168.0.2 time-exceeded 0 / IPerror / TCPerror
IP / TCP 192.168.0.2:38001 > 212.27.48.10:www S ==> IP / ICMP 192.168.0.1 > 192.168.0.2 time-exceeded 0 / IPerror / TCPerror
IP / TCP 192.168.0.2:33842 > 213.251.178.32:www S ==> IP / ICMP 192.168.0.1 > 192.168.0.2 time-exceeded 0 / IPerror / TCPerror
IP / TCP 192.168.0.2:55277 > 64.233.183.104:www S ==> IP / ICMP 192.168.2.1 > 192.168.0.2 time-exceeded 0 / IPerror / TCPerror
IP / TCP 192.168.0.2:62887 > 64.12.89.12:www S ==> IP / ICMP 192.168.2.1 > 192.168.0.2 time-exceeded 0 / IPerror / TCPerror
IP / TCP 192.168.0.2:amanda > 212.27.48.10:www S ==> IP / ICMP 192.168.2.1 > 192.168.0.2 time-exceeded 0 / IPerror / TCPerror
IP / TCP 192.168.0.2:55421 > 64.12.89.12:www S ==> IP / ICMP 81.62.144.1 > 192.168.0.2 time-exceeded 0 / IPerror / TCPerror
IP / TCP 192.168.0.2:46998 > 212.27.48.10:www S ==> IP / ICMP 81.62.144.1 > 192.168.0.2 time-exceeded 0 / IPerror / TCPerror
IP / TCP 192.168.0.2:34599 > 213.251.178.32:!www S ==> IP / ICMP 81.62.144.1 > 192.168.0.2 time-exceeded 0 / IPerror / TCPerror
IP / TCP 192.168.0.2:58715 > 64.233.183.104:www S ==> IP / ICMP 195.186.123.1 > 192.168.0.2 time-exceeded 0 / IPerror / TCPerror
IP / TCP 192.168.0.2:9093 > 64.12.89.12:www S ==> IP / ICMP 195.186.123.1 > 192.168.0.2 time-exceeded 0 / IPerror / TCPerror
IP / TCP 192.168.0.2:14089 > 212.27.48.10:www S ==> IP / ICMP 195.186.123.1 > 192.168.0.2 time-exceeded 0 / IPerror / TCPerror
IP / TCP 192.168.0.2:2258 > 213.251.178.32:www S ==> IP / ICMP 195.186.123.1 > 192.168.0.2 time-exceeded 0 / IPerror / TCPerror
IP / TCP 192.168.0.2:60402 > 64.233.183.104:www S ==> IP / ICMP 195.186.0.229 > 192.168.0.2 time-exceeded 0 / IPerror / TCPerror
IP / TCP 192.168.0.2:32073 > 64.12.89.12:www S ==> IP / ICMP 195.186.0.229 > 192.168.0.2 time-exceeded 0 / IPerror / TCPerror
IP / TCP 192.168.0.2:18250 > 64.12.89.12:www S ==> IP / ICMP 138.187.129.74 > 192.168.0.2 time-exceeded 0 / IPerror / TCPerror
IP / TCP 192.168.0.2:13695 > 212.27.48.10:www S ==> IP / ICMP 138.187.129.45 > 192.168.0.2 time-exceeded 0 / IPerror / TCPerror
IP / TCP 192.168.0.2:19761 > 213.251.178.32:www S ==> IP / ICMP 138.187.129.45 > 192.168.0.2 time-exceeded 0 / IPerror / TCPerror
IP / TCP 192.168.0.2:23264 > 64.233.183.104:www S ==> IP / ICMP 72.14.198.57 > 192.168.0.2 time-exceeded 0 / IPerror / TCPerror
IP / TCP 192.168.0.2:32236 > 64.12.89.12:www S ==> IP / ICMP 194.42.48.4 > 192.168.0.2 time-exceeded 0 / IPerror / TCPerror
IP / TCP 192.168.0.2:20611 > 212.27.48.10:www S ==> IP / ICMP 138.187.130.73 > 192.168.0.2 time-exceeded 0 / IPerror / TCPerror
IP / TCP 192.168.0.2:40775 > 213.251.178.32:www S ==> IP / ICMP 138.187.130.73 > 192.168.0.2 time-exceeded 0 / IPerror / TCPerror
IP / TCP 192.168.0.2:57456 > 64.233.183.104:www S ==> IP / ICMP 64.233.174.34 > 192.168.0.2 time-exceeded 0 / IPerror / TCPerror / Padding
IP / TCP 192.168.0.2:56892 > 64.12.89.12:www S ==> IP / ICMP 212.74.84.242 > 192.168.0.2 time-exceeded 0 / IPerror / TCPerror
```

We launch a new multiple traceroute but this time just to the FQDNs "www.free.fr" and "www.secuobs.com" with a ttl still equal to 10, the positive responses are still stored in the sn variable while the records without response go into the unans variable (for unanswered):

```bash
>>> sn,unans=traceroute(["www.free.fr","www.secuobs.com"],maxttl=10)
Begin emission:
***********Finished to send 20 packets.
****
Received 15 packets, got 15 answers, remaining 5 packets
212.27.48.10:tcp80 213.251.178.32:tcp80
1 192.168.0.1 11 -
2 192.168.2.1 11 -
3 192.168.1.1 11 -
4 - 81.62.144.1 11
5 195.186.123.1 11 -
6 195.186.0.229 11 195.186.0.229 11
7 195.186.0.234 11 195.186.0.234 11
8 138.187.129.45 11 138.187.129.45 11
9 138.187.130.73 11 138.187.130.73 11
10 138.187.129.10 11 138.187.129.10 11
```

From these results and using the graph() attribute, we can generate a graph of the routing taken by the data flows to reach these two different destinations from the test machine used for this document:

```bash
>>> sn.graph()
```

We can observe on the graph below these different routing results:

![Scapy.jpg](/images/scapy.avif)

We could also have saved the generated graph directly to an image file in our operating system's file system (here in the /tmp directory for an svg image format and an image file named graph.svg):

```bash
>>> sn.graph(target="> /tmp/graph.svg")
```

We check the existence and creation date of the graph.svg file contained in the /tmp directory:

```bash
root@casper:~# ls -l /tmp/graph.svg
-rw-r--r-- 1 root root 11103 2007-09-24 15:56 /tmp/graph.svg
```

We can view this graph.svg image file as follows with the display command which is part of the ImageMagick suite:

```bash
root@casper:~# display /tmp/graph.svg
```

To generate a graph of these results to a JPEG format image file, it is necessary to use the following command:

```bash
>>> sn.graph(target="> /tmp/graph.jpg")
```

We check that the graph generated in JPEG format to the graph.jpg image file in the /tmp directory is indeed present:

```bash
root@casper:~# ls -l /tmp/graph.jpg
-rw-r--r-- 1 root root 20997 2007-09-24 16:46 /tmp/graph.jpg
```

We can view this /tmp/graph.jpg image file still with the display function from Image Magick:

```bash
root@casper:~# display /tmp/graph.jpg
```

We can also directly print these graphs provided that a printer is connected to the machine and configured at the operating system level allowing it to function (here a postscript type printer):

```bash
>>> sn.graph(type="ps",target="| lp")
```

We now perform a new simple traceroute to the FQDN "www.free.fr" with a ttl equivalent to 10:

```bash
>>> sn,unans=traceroute(["www.free.fr"],maxttl=10)
Begin emission:
********Finished to send 10 packets.
**
Received 10 packets, got 10 answers, remaining 0 packets
212.27.48.10:tcp80
1 192.168.0.1 11
2 192.168.2.1 11
3 192.168.1.1 11
4 81.62.144.1 11
5 195.186.123.1 11
6 195.186.0.229 11
7 195.186.0.234 11
8 138.187.129.45 11
9 138.187.130.73 11
10 138.187.129.10 11
```

The representation of the results of this traceroute in a graph, this time in 3D, can be done as follows using vpython's visual plugin:

```bash
>>> sn.trace3D()
```

The vpython visual plugin we're using here allows you to zoom on this graph using the 3rd mouse button; if this button is emulated, simply press the left and right mouse or touchpad buttons simultaneously. You can also move the graphical representation with the left mouse button and tilt it with the right button:

![Scapy2.jpg](/images/scapy2.avif)

This time we'll perform a multiple traceroute to the FQDNs "www.free.fr", "www.google.fr" and "www.microsoft.com" with a ttl equivalent to 10; the responses are still stored in the sn variable while the records corresponding to no responses are contained in the unans variable:

```bash
>>> sn,unans=traceroute(["www.free.fr","www.google.fr",
"www.microsoft.com"],maxttl=10)
Begin emission:
************************Finished to send 30 packets.
*
Received 25 packets, got 25 answers, remaining 5 packets
207.46.19.254:tcp80 212.27.48.10:tcp80 64.233.183.99:tcp80
1 192.168.0.1 11 192.168.0.1 11 192.168.0.1 11
2 192.168.2.1 11 192.168.2.1 11 192.168.2.1 11
3 192.168.1.1 11 192.168.1.1 11 192.168.1.1 11
4 81.62.144.1 11 - -
5 195.186.123.1 11 195.186.123.1 11 195.186.123.1 11
6 195.186.0.229 11 195.186.0.229 11 195.186.0.229 11
7 195.186.0.254 11 195.186.0.234 11 195.186.0.254 11
8 138.187.159.54 11 - 138.187.129.46 11
9 138.187.159.18 11 - 138.187.129.74 11
10 164.128.236.38 11 138.187.129.10 11 -
```

We trace the 3D graph, using the trace3D() attribute, of the results contained in the sn variable that correspond to this multiple traceroute:

```bash
>>> sn.trace3D()
```

We get the following result as a 3D graph of the results:

![Scapy3.jpg](/images/scapy3.avif)

## Packet and Frame Manipulation

Scapy is above all a utility for forging, receiving, and sending packets and frames of data over a network. This section covers the many functions available for this purpose through different examples such as performing an ICMP ping or a port scan.

We now display the specifications of the IP protocol using the ls() command:

```bash
>>> ls(IP)
version : BitField = (4)
ihl : BitField = (None)
tos : XByteField = (0)
len : ShortField = (None)
id : ShortField = (1)
flags : FlagsField = (0)
frag : BitField = (0)
ttl : ByteField = (64)
proto : ByteEnumField = (0)
chksum : XShortField = (None)
src : Emph = (None)
dst : Emph = ('127.0.0.1')
options : IPoptionsField = ('')
```

We also display the specifications of the ICMP protocol still with the ls() command:

```bash
>>> ls(ICMP)
type : ByteEnumField = (8)
code : ByteField = (0)
chksum : XShortField = (None)
id : XShortField = (0)
seq : XShortField = (0)
```

We perform a ping to the machine with IP address 192.168.0.2; first we visualize the corresponding ICMP packet:

```bash
>>> IP(dst='192.168.0.2')/ICMP()
<IP frag=0 proto=icmp dst=192.168.0.2 |<ICMP |>>
```

We display the documentation of the sr1() function using the lsc() function:

```bash
>>> lsc(sr1)
Send packets at layer 3 and return only the first answer
nofilter: put 1 to avoid use of bpf filters
retry: if positive, how many times to resend unanswered packets
if negative, how many times to retry when no more packets are answered
timeout: how much time to wait after the last packet has been sent
verbose: set verbosity level
multi: whether to accept multiple answers for the same stimulus
filter: provide a BPF filter
iface: listen answers only on the given interface
```

We now send the same ICMP packet and accept to receive only one packet in response using the sr1() command:

```bash
>>> sr1(IP(dst='192.168.0.1')/ICMP())
Begin emission:
.Finished to send 1 packets.
*
Received 2 packets, got 1 answers, remaining 0 packets
<IP version=4L ihl=5L tos=0x0 len=28 id=31775 flags= frag=0L ttl=128 proto=icmp chksum=0x3d6e src=192.168.0.1 dst=192.168.0.2 options='' |<ICMP type=echo-reply code=0 chksum=0xffff id=0x0 seq=0x0 |>>
```

The ping is effective and is similar to the results of the system ping command:

```bash
root@casper:~# ping -c 1 192.168.0.1
PING 192.168.0.1 (192.168.0.1) 56(84) bytes of data.
64 bytes from 192.168.0.1: icmp_seq=1 ttl=128 time=1.40 ms

--- 192.168.0.1 ping statistics ---
1 packets transmitted, 1 received, 0% packet loss, time 0ms
rtt min/avg/max/mdev = 1.409/1.409/1.409/0.000 ms
```

We note with the sr1() function that we only receive and record one packet in response. This same function can also be used to perform a port scan; we send a TCP/IP packet to www.secuobs.com to port 80 with the TCP Syn flag activated which indicates that we are asking to establish a connection without preliminaries:

```bash
>>> sr1(IP(dst="www.secuobs.com")/TCP(dport=80, flags="S"))
Begin emission:
Finished to send 1 packets.
*
Received 1 packets, got 1 answers, remaining 0 packets
<IP version=4L ihl=5L tos=0x0 len=44 id=0 flags=DF frag=0L ttl=52 proto=tcp chksum=0xfe05 src=213.251.178.32 dst=192.168.0.2 options='' |<TCP sport=www dport=ftp_data seq=466528407 ack=1 dataofs=6L reserved=0L flags=SA window=5840 chksum=0x73bd urgptr=0 options=[('MSS', 1452)] |>>
```

We note here that we received a packet in response, we can conclude that port 80 (www) is open on the server whose IP address is assigned to the FQDN "www.secuobs.com"; we make a second attempt on port 22 (SSH):

```bash
>>> sr1(IP(dst="www.secuobs.com")/TCP(dport=22, flags="S"))
Begin emission:
.Finished to send 1 packets.
*
Received 2 packets, got 1 answers, remaining 0 packets
<IP version=4L ihl=5L tos=0x0 len=44 id=0 flags=DF frag=0L ttl=52 proto=tcp chksum=0xfe05 src=213.251.178.32 dst=192.168.0.2 options='' |<TCP sport=ssh dport=ftp_data seq=771588344 ack=1 dataofs=6L reserved=0L flags=SA window=5840 chksum=0x8967 urgptr=0 options=[('MSS', 1452)] |>>
```

We also receive a response, port 22 is open; we make one last attempt on port 53 (DNS):

```bash
>>> sr1(IP(dst="www.secuobs.com")/TCP(dport=53, flags="S"))
Begin emission:
Finished to send 1 packets.
..............................
Received 30 packets, got 0 answers, remaining 1 packets
```

We conclude that port 53 is not open on this server, it doesn't necessarily mean that a DNS server isn't running on this port, we can just conclude that filtering rules are present to prevent external connections to this port 53.

We now display, using the lsc() command, the details of the options of the sr() function which allows it to receive and store more than one packet in response to those previously sent:

```bash
>>> lsc(sr)
Send and receive packets at layer 3
nofilter: put 1 to avoid use of bpf filters
retry: if positive, how many times to resend unanswered packets
if negative, how many times to retry when no more packets are answered
timeout: how much time to wait after the last packet has been sent
verbose: set verbosity level
multi: whether to accept multiple answers for the same stimulus
filter: provide a BPF filter
iface: listen answers only on the given interface
```

We display the details of the send() function using the lsc() command:

```bash
>>> lsc(send)
Send packets at layer 3
send(packets, [inter=0], [loop=0], [verbose=conf.verb]) -> None
```

We can also use Scapy's srloop function to execute a loop on sending a data packet (here ICMP to the machine whose IP address is equal to 192.168.0.2):

```bash
>>> srloop(IP(dst='192.168.0.2')/ICMP())
RECV 1: IP / ICMP 192.168.0.1 > 192.168.0.2 echo-reply 0
RECV 1: IP / ICMP 192.168.0.1 > 192.168.0.2 echo-reply 0
RECV 1: IP / ICMP 192.168.0.1 > 192.168.0.2 echo-reply 0
RECV 1: IP / ICMP 192.168.0.1 > 192.168.0.2 echo-reply 0
RECV 1: IP / ICMP 192.168.0.1 > 192.168.0.2 echo-reply 0
RECV 1: IP / ICMP 192.168.0.1 > 192.168.0.2 echo-reply 0
RECV 1: IP / ICMP 192.168.0.1 > 192.168.0.2 echo-reply 0
RECV 1: IP / ICMP 192.168.0.1 > 192.168.0.2 echo-reply 0
RECV 1: IP / ICMP 192.168.0.1 > 192.168.0.2 echo-reply 0

Sent 9 packets, received 9 packets. 100.0% hits.
(<Results: UDP:0 TCP:0 ICMP:9 Other:0>, <PacketList: UDP:0 TCP:0 ICMP:0 Other:0>)
```

Press Ctrl+C to stop srloop, here we have sent and received 9 packets; we can also assign to srloop() a limited send number with the count parameter, here 10 packets:

```bash
>>> srloop(IP(dst='192.168.0.2')/ICMP(),count=10)
RECV 1: IP / ICMP 192.168.0.1 > 192.168.0.2 echo-reply 0
RECV 1: IP / ICMP 192.168.0.1 > 192.168.0.2 echo-reply 0
RECV 1: IP / ICMP 192.168.0.1 > 192.168.0.2 echo-reply 0
RECV 1: IP / ICMP 192.168.0.1 > 192.168.0.2 echo-reply 0
RECV 1: IP / ICMP 192.168.0.1 > 192.168.0.2 echo-reply 0
RECV 1: IP / ICMP 192.168.0.1 > 192.168.0.2 echo-reply 0
RECV 1: IP / ICMP 192.168.0.1 > 192.168.0.2 echo-reply 0
RECV 1: IP / ICMP 192.168.0.1 > 192.168.0.2 echo-reply 0
RECV 1: IP / ICMP 192.168.0.1 > 192.168.0.2 echo-reply 0
RECV 1: IP / ICMP 192.168.0.1 > 192.168.0.2 echo-reply 0

Sent 10 packets, received 10 packets. 100.0% hits.
(<Results: UDP:0 TCP:0 ICMP:10 Other:0>, <PacketList: UDP:0 TCP:0 ICMP:0 Other:0>)
```

The functions sr, sr1, send and srloop send and/or receive packets at the network level, which is layer 3 of the OSI model, that is, the one present just before (emission) or after (reception) the level corresponding to the data link which is in second position in this 7-layer model.

The data link layer is therefore just before (emission) or after (reception) the level corresponding to the physical layer of the OSI model.

The equivalent of these network level functions for the data link level in Scapy are the functions srp, srp1, srploop and sendp which forge, receive and/or send data frames and not just data packets on the network; note that at the physical layer of the OSI model we speak in bytes and no longer in packets or data frames.

We display the documentation of srp():

```bash
>>> lsc(srp)
Send and receive packets at layer 2
nofilter: put 1 to avoid use of bpf filters
retry: if positive, how many times to resend unanswered packets
if negative, how many times to retry when no more packets are answered
timeout: how much time to wait after the last packet has been sent
verbose: set verbosity level
multi: whether to accept multiple answers for the same stimulus
filter: provide a BPF filter
iface: work only on the given interface
```

Consider the following frame:

```bash
>>> Ether()/IP(dst="www.secuobs.com")/TCP(dport=[80,443])/"GET / HTTP/1.0 \n\n"
Ether()/IP(dst="www.secuobs.com")/TCP(dport=[80,443])/"GET / HTTP/1.0 \n\n"
<Ether type=IPv4 |<IP frag=0 proto=tcp dst=Net('www.secuobs.com') |<TCP dport=['www', 'https'] |<Raw load='GET / HTTP/1.0 \n\n' |>>>>
```

We assign this frame to the variable sn:

```bash
>>> sn=Ether()/IP(dst="www.secuobs.com")/TCP(dport=[80,443])/"GET / HTTP/1.0 \n\n"

>>> sn
<Ether type=0x800 |<IP frag=0 proto=tcp dst=Net('www.secuobs.com') |<TCP dport=['www', 'https'] |<Raw load='GET / HTTP/1.0 \n\n' |>>>>
```

Using the srp function, we send the frame contained in the sn variable:

```bash
>>> srp(sn)
Begin emission:
.Finished to send 2 packets.
**
Received 3 packets, got 2 answers, remaining 0 packets
(<Results: UDP:0 TCP:2 ICMP:0 Other:0>, <Unanswered: UDP:0 TCP:0 ICMP:0 Other:0>)
```

We could also have assigned certain values to the sending of the frame like the number of times they should be resent after a negative response thanks to the retry parameter (here 3 times), the time interval between the sending of two frames with inter (here 1) or the time limit (here 2) to wait after sending the last frame with timeout (all of these parameters also being available for sending packets via the sr() function):

```bash
>>> srp(Ether()/IP(dst="www.secuobs.com")/TCP(dport=[80,443])/"GET / HTTP/1.0 \n\n",inter=1,retry=-3,timeout=2)
Begin emission:
...Finished to send 2 packets.
*..Begin emission:
..Finished to send 1 packets.
*
Received 9 packets, got 2 answers, remaining 0 packets
(<Results: UDP:0 TCP:2 ICMP:0 Other:0>, <Unanswered: UDP:0 TCP:0 ICMP:0 Other:0>)
```

## Object Orientation and Representation

Scapy provides great flexibility in handling the various commands available. It's thus possible to manage each field of a packet or a function in the manner of an object that can be defined in a variable and then represented and visualized in several ways.

Display the specifications of the DNS protocol:

```bash
>>> ls(DNS)
id : ShortField = (0)
qr : BitField = (0)
opcode : BitEnumField = (0)
aa : BitField = (0)
tc : BitField = (0)
rd : BitField = (0)
ra : BitField = (0)
z : BitField = (0)
rcode : BitEnumField = (0)
qdcount : DNSRRCountField = (None)
ancount : DNSRRCountField = (None)
nscount : DNSRRCountField = (None)
arcount : DNSRRCountField = (None)
qd : DNSQRField = (None)
an : DNSRRField = (None)
ns : DNSRRField = (None)
ar : DNSRRField = (None)
```

We retrieve a list of DNS servers by displaying the contents of the /etc/resolv.conf file:

```bash
root@casper:~# cat /etc/resolv.conf
nameserver 62.4.17.69
nameserver 62.4.16.70
```

We perform a UDP DNS request using the sr1 send function on the DNS server 62.4.16.70 for the FQDN exoscan.net:

```bash
<IP version=4L ihl=5L tos=0x0 len=135 id=13021 flags= frag=0L ttl=62 proto=udp chksum=0x3a95 src=62.4.16.70 dst=192.168.0.2 options='' |<UDP sport=domain dport=domain len=115 chksum=0x7504 |<DNS id=0 qr=1L opcode=QUERY aa=0L tc=0L rd=1L ra=1L z=0L rcode=ok qdcount=1 ancount=1 nscount=2 arcount=1 qd=<DNSQR qname='exoscan.net.' qtype=A qclass=IN |> an=<DNSRR rrname='exoscan.net.' type=A rclass=IN ttl=3469 rdata='213.186.41.29' |> ns=<DNSRR rrname='exoscan.net.' type=NS rclass=IN ttl=28243 rdata='ns7.gandi.net.' |<DNSRR rrname='exoscan.net.' type=NS rclass=IN ttl=28243 rdata='custom2.gandi.net.' |>> ar=<DNSRR rrname='ns7.gandi.net.' type=A rclass=IN ttl=643 rdata='217.70.177.44' |> |>>>
```

We check the value of \_ by adding the summary() attribute to it:

```bash
>>> _.summary()
'IP / UDP / DNS Ans "213.186.41.29" '
```

We display the result in more detail with the display() attribute instead of summary():

```bash
>>> _.display()
###[ IP ]###
version= 4L
ihl= 5L
tos= 0x0
len= 135
id= 13021
flags=
frag= 0L
ttl= 62
proto= udp
chksum= 0x3a95
src= 62.4.16.70
dst= 192.168.0.2
options= ''
###[ UDP ]###
sport= domain
dport= domain
len= 115
chksum= 0x7504
###[ DNS ]###
id= 0
qr= 1L
opcode= QUERY
aa= 0L
tc= 0L
rd= 1L
ra= 1L
z= 0L
rcode= ok
qdcount= 1
ancount= 1
nscount= 2
arcount= 1
\qd\
|###[ DNS Question Record ]###
| qname= 'exoscan.net.'
| qtype= A
| qclass= IN
\an\
|###[ DNS Resource Record ]###
| rrname= 'exoscan.net.'
| type= A
| rclass= IN
| ttl= 3469
| rdlen= 0
| rdata= '213.186.41.29'
\ns\
|###[ DNS Resource Record ]###
| rrname= 'exoscan.net.'
| type= NS
| rclass= IN
| ttl= 28243
| rdlen= 0
| rdata= 'ns7.gandi.net.'
|###[ DNS Resource Record ]###
| rrname= 'exoscan.net.'
| type= NS
| rclass= IN
| ttl= 28243
| rdlen= 0
| rdata= 'custom2.gandi.net.'
\ar\
|###[ DNS Resource Record ]###
| rrname= 'ns7.gandi.net.'
| type= A
| rclass= IN
| ttl= 643
| rdlen= 0
| rdata= '217.70.177.44'
```

We now launch a capture on all traffic, still using the sniff() command, asking to capture only one record (still with the count option) which we will place as usual in the sn variable:

```bash
>>> sn=sniff(count=1)
```

We display this record using the display() attribute applied to the sn variable:

```bash
>>> sn.display()
0000 Ether / IP / TCP 127.0.0.1:50018 > 127.0.0.1:microsoft_ds S
```

As we've already noted previously, Scapy follows an object orientation by offering significant flexibility in the operations of manipulating the different functions that are present there but also in the parameters and attributes that refer to these functions; we now display again the specifications of the IP protocol using the ls() command:

```bash
>>> ls(IP)
version : BitField = (4)
ihl : BitField = (None)
tos : XByteField = (0)
len : ShortField = (None)
id : ShortField = (1)
flags : FlagsField = (0)
frag : BitField = (0)
ttl : ByteField = (64)
proto : ByteEnumField = (0)
chksum : XShortField = (None)
src : Emph = (None)
dst : Emph = ('127.0.0.1')
options : IPoptionsField = ('')
```

Each field of a packet or frame can also be considered as an object and its value can then be defined by a variable just as it is possible to do in the same way for parts or sub-parts of functions and this with a direct rebound effect on the values (of these fields) that were defined by default:

```bash
>>> IP()
<IP |>
>>> a=IP(dst="192.168.0.2")
>>> a
<IP dst=192.168.0.2 |>
>>> ls(a)
version : BitField = 4 (4)
ihl : BitField = None (None)
tos : XByteField = 0 (0)
len : ShortField = None (None)
id : ShortField = 1 (1)
flags : FlagsField = 0 (0)
frag : BitField = 0 (0)
ttl : ByteField = 64 (64)
proto : ByteEnumField = 0 (0)
chksum : XShortField = None (None)
src : Emph = '192.168.0.2' (None)
dst : Emph = '192.168.0.2' ('127.0.0.1')
options : IPoptionsField = '' ('')
>>> a.dst
'192.168.0.2'
>>> a.ttl
64
>>> a.ttl=10
>>> a
<IP ttl=10 dst=192.168.0.2 |>
>>> a
<IP ttl=10 dst=192.168.0.2 |>
>>> del a.ttl
>>> a
<IP dst=192.168.0.2 |>
>>> a.ttl
64
>>> del a.dst
>>> a.dst
'127.0.0.1'
```

As we just saw by changing the values of the variables a.dst and a.ttl the modifications are directly impacted on the variable a and the values of the ttl and dst parameters are replaced by the changes made or they are simply deleted in the same way as the parameter when it has been previously deleted.

The different parameters of the function however return to default values when deleted as we can see with the value of a.dst which is equal to 127.0.0.1 in the end or that of a.ttl which is equivalent to 64.

In Scapy we can also generate a representation of a packet and/or frame type object with the str function whose specifications are as follows:

```bash
>>> lsc(str)
str(object) -> string

Return a nice string representation of the object.
If the argument is a string, the return value is the same object.
```

We generate the representation of the IP() function with all default values:

```bash
>>> IP()
<IP |>

>>> str(IP())
'E\x00\x00\x14\x00\x01\x00\x00@\x00|\xe7\x7f\x00\x00
\x01\x7f\x00\x00\x01'
```

We can also retrieve the initial packet thanks to the IP() function as follows:

```bash
>>> IP(_)
<IP version=4L ihl=5L tos=0x0 len=20 id=1 flags= frag=0L ttl=64 proto=ip chksum=0x7ce7 src=127.0.0.1 dst=127.0.0.1 |>
```

We generate an IP packet whose destination is www.exoscan.net that we assign to the variable sn:

```bash
>>> sn=IP(dst="www.exoscan.net")
>>> sn
<IP dst=Net('www.exoscan.net') |>
```

We apply the str function on this variable whose result is the following representation:

```bash
>>> str(sn)
'E\x00\x00\x14\x00\x01\x00\x00@\x00\xbbg\xc0\xa8\x00
\x02\xd5\xba)\x1d'
```

The same packet for the TCP protocol specifically:

```bash
>>> IP(dst="www.exoscan.net")/TCP()
<IP frag=0 proto=tcp dst=Net('www.exoscan.net') |<TCP |>>

>>> str(IP(dst="www.exoscan.net")/TCP())
'E\x00\x00(\x00\x01\x00\x00@\x06\xbbM\xc0\xa8\x00\x02
\xd5\xba)\x1d\x00\x14\x00P\x00\x00\x00\x00\x00\x00\
x00\x00P\x02 \x00\xcf\xfc\x00\x00'
```

We now generate an Ethernet frame (level 2 data link layer in the OSI model for the TCP/IP protocol suite) corresponding to this packet which coupled the IP() and TCP() functions and now the Ether() function:

```bash
>>> Ether()/IP(dst="www.exoscan.net")/TCP()
<Ether type=0x800 |<IP frag=0 proto=tcp dst=Net('www.exoscan.net') |<TCP |>>>

>>> str(Ether()/IP(dst="www.exoscan.net")/TCP())
'\x00\x13\xf7x\xcf\xea\x00\x0f\xb5\x0e\x89J\x08\x00E
\x00\x00(\x00\x01\x00\x00@\x06\xbbM\xc0\xa8\x00\x02
\xd5\xba)\x1d\x00\x14\x00P\x00\x00\x00\x00\x00\x00
\x00\x00P\x02 \x00\xcf\xfc\x00\x00'
```

We generate the same frame except that we specify the destination ports on ports 80 and 443:

```bash
>>> Ether()/IP(dst="www.exoscan.net")/TCP(dport=[80,443])
<Ether type=0x800 |<IP frag=0 proto=tcp dst=Net('www.exoscan.net') |<TCP dport=['www', 'https'] |>>>

>>> str(Ether()/IP(dst="www.exoscan.net")/TCP(dport=[80,443]))
'\x00\x13\xf7x\xcf\xea\x00\x0f\xb5\x0e\x89J\x08\x00E
\x00\x00(\x00\x01\x00\x00@\x06\xbbM\xc0\xa8\x00\x02
\xd5\xba)\x1d\x00\x14\x00P\x00\x00\x00\x00\x00\x00
\x00\x00P\x02 \x00\xcf\xfc\x00\x00'
```

We can retrieve the initial frame from this representation thanks to the Ether() function:

```bash
>>> Ether(_)
<Ether dst=00:13:f7:78:cf:ea src=00:0f:b5:0e:89:4a type=0x800 |<IP version=4L ihl=5L tos=0x0 len=40 id=1 flags= frag=0L ttl=64 proto=tcp chksum=0xbb4d src=192.168.0.2 dst=213.186.41.29 options='' |<TCP sport=ftp_data dport=www seq=0 ack=0 dataofs=5L reserved=0L flags=S window=8192 chksum=0xcffc urgptr=0 |>>>
```

We now generate the same frame but with a request to download the index page of the "www.exoscan.net" website:

```bash
>>> Ether()/IP(dst="www.exoscan.net")/TCP(dport=[80,443])/"GET / HTTP/1.0 \n\n"
<Ether type=0x800 |<IP frag=0 proto=tcp dst=Net('www.exoscan.net') |<TCP dport=['www', 'https'] |<Raw load='GET / HTTP/1.0 \n\n' |>>>>

>>> str(Ether()/IP(dst="www.exoscan.net")/TCP(dport=[80,443])/"GET / HTTP/1.0 \n\n")
'\x00\x0f\xb5\x0e\x89J\x00\x15mS\x1e\x87\x08\x00E\x00
\x009\x00\x01\x00\x00@\x06\xbb<\xc0\xa8\x00\x02\xd5
\xba)\x1d\x00\x14\x00P\x00\x00\x00\x00\x00\x00\x00\x00P
\x02 \x00\xe1U\x00\x00GET / HTTP/1.0 \n\n'

>>> Ether(_)
<Ether dst=00:0f:b5:0e:89:4a src=00:15:6d:53:1e:87 type=IPv4 |<IP version=4L ihl=5L tos=0x0 len=57 id=1 flags= frag=0L ttl=64 proto=tcp chksum=0xbb3c src=192.168.0.2 dst=213.186.41.29 options='' |<TCP sport=ftp_data dport=www seq=0 ack=0 dataofs=5L reserved=0L flags=S window=8192 chksum=0xe155 urgptr=0 options=[] |<Raw load='GET / HTTP/1.0 \n\n' |>>>>
```

Note that all fields with a default value are visible, it's possible to remove them with the hide*default() attribute applied to the current result * (or to a variable):

```bash
>>> _.hide_defaults()
```

We verify the value of \_:

```bash
>>> _
<Ether dst=00:0f:b5:0e:89:4a src=00:15:6d:53:1e:87 type=IPv4 |<IP ihl=5L len=57 frag=0 proto=tcp chksum=0xbb3c src=192.168.0.2 dst=213.186.41.29 |<TCP dataofs=5L chksum=0xe155 options=[] |<Raw load='GET / HTTP/1.0 \n\n' |>>>><>
```

It's also possible with scapy to operate a hexdump to get a hexadecimal version of the records we receive (or those we send):

```bash
>>> hexdump(sn)
0000 00 0F B5 0E 89 4A 00 15 6D 53 1E 87 08 00 45 00 .....J..mS....E.
0010 00 39 00 01 00 00 40 06 BB 3C C0 A8 00 02 D5 BA .9....@..<......
0020 29 1D 00 14 00 50 00 00 00 00 00 00 00 00 50 02 )....P........P.
0030 20 00 E1 55 00 00 47 45 54 20 2F 20 48 54 54 50 ..U..GET / HTTP
0040 2F 31 2E 30 20 0A 0A /1.0 ..

>>> hexedit(_)
'\x00\x0f\xb5\x0e\x89J\x00\x15mS\x1e\x87\x08\x00E\x00
\x009\x00\x01\x00\x00@\x06\xbb<\xc0\xa8\x00\x02\xd5
\xba)\x1d\x00\x14\x00P\x00\x00\x00\x00\x00\x00\x00\x00P
\x02 \x00\xe1U\x00\x00GET / HTTP/1.0 \n\n'
```

## Supplements and Webography

This part of the Scapy documentation allows you to find all the web addresses related to this documentation and to the use of Scapy. It also includes an example of a custom tool developed using Scapy and the Python language, in our case a TCP port scanner.

You can find different PDF slides about Scapy including those from the PacSec Core05 conference, Hack.lu 2005, Summerschool Applied IT Security 2005, T2 2005, CanSecWest Core05 and LSM 2003.

For users of Microsoft Windows operating systems, you can get the Windows port of Scapy (while OpenBSD system users can read this document for installation).

If you want to create specific scripts with Scapy, you can consult the document provided for this purpose on the official Scapy website, SecDev.org. For example, if you want to create a TCP port scanner, the corresponding Python script should look more or less like this (here for ports 22 to 25 inclusive):

```python
#!/usr/bin/env python

import sys
from scapy import *

target = sys.argv[1]
fl = 22

while (fl<=25):
   p=sr1(IP(dst=target)/TCP(dport=fl, flags="S"),retry=0,timeout=1)
   if p:
     print "\n Port " + str(fl) + " TCP is open on " + str(target) + "\n"
   else:
     print "\n Port " + str(fl) + " TCP is not open on " + str(target) + "\n"
   fl = fl + 1
```

We save this script under the name scan.py (note that the file must be in the same directory as the scapy.py file) which we will run with the Python interpreter by passing the IP address or FQDN of the machine we want to define as the target of the port scan as follows:

```bash
root@casper:~# python scan.py www.secuobs.com

Begin emission:
.......................Finished to send 1 packets.
.*
Received 25 packets, got 1 answers, remaining 0 packets

Port 22 TCP is open on www.secuobs.com

Begin emission:
..Finished to send 1 packets.
.............................
Received 31 packets, got 0 answers, remaining 1 packets

Port 23 TCP is not open on www.secuobs.com

Begin emission:
..Finished to send 1 packets.
................................
Received 34 packets, got 0 answers, remaining 1 packets

Port 24 TCP is not open on www.secuobs.com

Begin emission:
..Finished to send 1 packets.
...*
Received 6 packets, got 1 answers, remaining 0 packets

Port 25 TCP is open on www.secuobs.com
```

The same script to test the first 1024 ports:

```python
#!/usr/bin/env python

import sys
from scapy import *

target = sys.argv[1]
fl = 1

while (fl<=1024):
   p=sr1(IP(dst=target)/TCP(dport=fl, flags="S"),retry=0,timeout=1)
   if p:
     print "\n Port " + str(fl) + " TCP is open on " + str(target) + "\n"
   else:
     print "\n Port " + str(fl) + " TCP is not open on " + str(target) + "\n"
   fl = fl + 1
```

Given the number of ports to test and the length of the result, it's preferable to launch it in the following way:

```bash
root@casper:~# python scanfull.py www.secuobs.com > portscan_full_secuobs
```

We can still view the result in real time as follows:

```bash
root@casper:/trash# tail -f portscan_full_secuobs
Received 1 packets, got 0 answers, remaining 1 packets

Port 1 TCP is not open on www.secuobs.com


Received 0 packets, got 0 answers, remaining 1 packets

Port 2 TCP is not open on www.secuobs.com


Received 0 packets, got 0 answers, remaining 1 packets

Port 3 TCP is not open on www.secuobs.com
```

Checking open ports:

```bash
root@casper:~# grep " is open" portscan_full_secuobs
Port 22 TCP is open on www.secuobs.com
Port 25 TCP is open on www.secuobs.com
Port 80 TCP is open on www.secuobs.com
Port 443 TCP is open on www.secuobs.com
Port 995 TCP is open on www.secuobs.com
```

The result with NMAP:

```bash
root@casper:~# nmap -P0 -sS www.secuobs.com

Starting Nmap 4.20 at 2007-10-02 17:19 CEST
Interesting ports on ns21533.ovh.net (213.251.178.32):
Not shown: 1691 filtered ports
PORT STATE SERVICE
22/tcp open ssh
25/tcp open smtp
80/tcp open http
443/tcp open https
995/tcp open pop3s
```

Note that it may be necessary to adjust the timeout value in the line "p=sr1(IP(dst=target)/TCP(dport=fl, flags="S"),retry=0,timeout=1)" depending on the locations of the different machines present (source and destination) in order to obtain satisfactory results.

A whole chapter dedicated to Scapy was written by the author himself, Philippe Biondi; this chapter is published in the book Security Power Tools; 856 pages published by O'Reilly Media (ISBN-10: 0596009631; ISBN-13 978-0596009632).

## Resources
- http://www.secuobs.com/news/01102007-scapy1.shtml
