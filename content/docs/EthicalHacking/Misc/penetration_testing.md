---
weight: 999
url: "/Tests_d'intrusion/"
title: "Penetration Testing"
description: "A comprehensive guide to penetration testing methodologies, including information gathering, vulnerability analysis, and exploitation techniques."
categories: ["Security", "Linux", "Network"]
date: "2006-11-28T22:23:00+02:00"
lastmod: "2006-11-28T22:23:00+02:00"
tags: ["penetration testing", "security", "network", "vulnerability", "hacking", "whois", "DNS", "port scan", "fingerprinting"]
toc: true
---

## Introduction

Once a network is installed and configured, it continuously evolves. New systems are added, old and faithful machines disappear - everything changes constantly. Users can also make modifications to the network without the administrator's knowledge.

To verify the state of a network, an administrator must behave like a hacker and attempt to penetrate their own defenses. This article presents the methodology to follow.

It's important to test the security of your own network by putting yourself in the hacker's position to discover potential vulnerabilities. These tests are broken down into several steps:

- **Approach Phase**: Gathering information about the target network
- **Analysis Phase**: Using the results obtained in the previous step to determine potential vulnerabilities and the tools needed to exploit them
- **Attack Phase**: Taking action

Before undertaking the final step, the network administrator must explicitly give their consent, and not just verbal approval. Chapter III of the Criminal Code deals with attacks on automated data processing systems. It contains 7 articles, here are the first three:

- Article 323-1: Fraudulently accessing or maintaining access to all or part of an automated data processing system is punishable by one year imprisonment and a fine of 15,000 Euros.

When this results in either the deletion or modification of data contained in the system, or an alteration of the operation of this system, the penalty is two years imprisonment and a fine of 17,500 Euros.

- Article 323-2: Hindering or distorting the operation of an automated data processing system is punishable by three years imprisonment and a fine of 30,000 euros.
- Article 323-3: Fraudulently introducing data into an automated data processing system or fraudulently deleting or modifying data it contains is punishable by three years imprisonment and a fine of 15,000 Euros.

Article 323-1 concerns the intrusion itself. When it also causes alteration of the system data, the penalty is increased. Article 323-2 deals with damage done to the network (viruses, mail bombing, DoS, etc.). Finally, Article 323-3 punishes changes deliberately made to data present on the network (meaning both high scores in games like xbill as well as network configuration files).

Testing network security generally means "resistance to external threats". However, many malicious operations can easily be conducted from a machine within the network itself (viruses, backdoors, sniffers, etc.). These sources of danger are rarely taken into consideration when evaluating the network. Similarly, most of the attacks presented in Eric Detoisien's article should also be evaluated (spoofing, (D)DoS, etc.).

Various scenarios can be used for penetration testing:

- The tester knows nothing about the target network and has no access to it; this is an external penetration test
- The tester has minimal privileges on the target network (any user account). They try to increase their privileges from inside the network itself (unactivated screen savers, sniffing, exploitation of local vulnerabilities, etc.). In this case, it is an internal penetration test.

The results obtained should allow for better identification and correction of potential or existing problems on the target network. For example, the administrator's behavior in response to intrusion attempts or a successful intrusion should be evaluated (will it be detected? How long will it take? etc.)

## The Information War

We are now assuming the role of someone trying to gather as much information as possible about a target network. This collection is divided into two stages. First, we'll gather all available information without directly accessing the target's resources. Then, when we begin to have a clearer idea of what we're dealing with, we'll directly access the resources provided by the target.

### Indirect Queries

In this category, we include all means that allow us to learn more about the target without directly contacting it. This information is available - you just need to know where to look.

#### Interrogation of Whois Databases

Whois servers, also called nicname (port 43), provide access to the database of information provided when registering a domain name:

- Administrative information such as names, phone numbers and addresses for different contacts (admin-c, tech-c, zone-c, bill-c...)
- Technical information such as DNS name(s), email addresses of the officials mentioned above, IP address ranges allocated to the target...

This database, formerly managed by InterNIC and now by Network Solutions, remains easily accessible to everyone as it allows verification of domain name availability.

Companies that register domain names generally offer an online query service (see Table 1). On Unix, there is also the whois command.

{{< table "table-hover table-striped" >}}
| Name | Meaning | URL | Description |
|------|---------|-----|-------------|
| AFNIC | Association Française pour le Nommage Internet en Coopération | http://www.nic.fr/cgi-bin/whois | all ".fr" domains |
| RIPE | Réseau IP Européen | http://www.ripe.net/cgi-bin/whois | covers Europe, Middle East and some Asian and African countries |
| InterNIC | Stated on their webpage: "InterNIC is a registered service mark of the U.S. Department of Commerce. This site is being hosted by ICANN on behalf of the U.S. Department of Commerce". | http://www.internic.net/whois.html | ".com", ".net", ".edu" and other ".org" domains |
{{< /table >}}

To show you the richness of information contained in this type of database, here's a small example using a domain name (fictitious for confidentiality reasons):

```bash
whois pigeons.fr@whois.nic.fr
```

```
[whois.nic.fr]

Tous droits reserves par copyright.
Voir http://www.nic.fr/outils/dbcopyright.html
Rights restricted by copyright.
See http://www.nic.fr/outils/dbcopyright.html

domain:      pigeons.fr
descr:       PIGEON ET CIE
descr:       10 RUE DE PARIS
descr:       75001 Paris 
admin-c:     BPxxx-FRNIC
tech-c:      LPxxx-FRNIC
zone-c:      CPxxx-FRNIC

nserver:     ns1.pigeons.fr
nserver:     ns2.pigeons.fr
nserver:     ns.heberge.fr

mnt-by:      FR-NIC-MNT
mnt-lower:   FR-NIC-MNT
changed:     frnic-dbm-updates@nic.fr 20001229
source:      FRNIC

role:        HEBERGE Hostmaster
address:     Heberge Telecom
address:     10 rue de Gennevilliers
address:     92230 Gennevilliers
phone:       +33 1 41 00 00 00
fax-no:      +33 1 41 00 00 01
e-mail:      hostmaster@heberge.fr
admin-c:     AAxxx-FRNIC
tech-c:      BBxxx-FRNIC
nic-hdl:     CCxxx-FRNIC
notify:      hm-dbm-msgs@ripe.net
mnt-by:      HEBERGE-NOC
changed:     hostmaster@heberge.fr 20000814
changed:     migration-dbm@nic.fr 20001015
source:      FRNIC

person:      Bernard Pigeon
address:     PIGEON ET CIE
address:     10 RUE DE PARIS
address:     75001 Paris
address:     France
phone:       +33 1 53 00 00 00
fax-no:      +33 1 53 00 00 01
e-mail:      bernard.pigeon@pigeons.fr
nic-hdl:     BPxxxx-FRNIC
notify:      bernard.pigeon@pigeons.fr
mnt-by:      HEBERGE-NOC
changed:     bernard.pigeon@pigeons.fr 20001228
source:      FRNIC

person:      Luc Pigeon
address:     PIGEON ET CIE
address:     10 RUE DE PARIS
address:     75001 Paris
address:     France
phone:       +33 1 53 00 00 00
fax-no:      +33 1 53 00 00 01
e-mail:      luc.pigeon@pigeons.fr
nic-hdl:     LPxxx-FRNIC
mnt-by:      HEBERGE-NOC
changed:     aaa@heberge.fr 20001228
source:      FRNIC
```

We learn that the company Pigeon et Cie is actually hosted by Heberge Telecom. The email addresses probably correspond to aliases, but they also provide clues about the existence of accounts on machines: if the corresponding passwords are weak (like first names, birth dates...), this information could be useful.

Also via whois databases, more advanced searches using an IP address belonging to Pigeon et Cie lead us to discover the IP address ranges allocated to them. For this, we use the RIPE whois database which proves to be the most relevant:

```bash
whois 10.51.23.246@whois.ripe.net
```

```
% This is the RIPE Whois server.
% The objects are in RPSL format.
% Please visit http://www.ripe.net/rpsl for more information.
% Rights restricted by copyright.
% See http://www.ripe.net/ripencc/pub-services/db/copyright.html

inetnum:      10.51.23..0 - 10.51.23.255
netname:      PIGEON-CIE
descr:        Pigeon et Cie
country:      FR
admin-c:      OMxxxx-RIPE
tech-c:       OMxxxx-RIPE
status:       ASSIGNED PA
notify:       admin@pigeons.fr
mnt-by:       PC-XXX
changed:      admin@pigeons.fr 20000223
source:       RIPE
```

In addition, cross-searches on recovered names can give us additional information about the target (discovering new IP addresses or new DNS servers).

The outcome of whois database searches is very conclusive as we have retrieved:

- The DNS servers with authority over the pigeons.fr domain
- The contact details of administrative and technical contacts related to Pigeon et Cie
- The IP address ranges allocated to Pigeon et Cie.

We can still continue to gather information on the internet.

#### News Groups

Administrators or developers often face problems. To solve them, they use NewsGroups to ask questions. Unfortunately, they often give too much information about their information system (technology, versions of applications used, code fragments...). To do research, we can ask groups.google.com with search criteria like @pigeons.fr for example. This very powerful search engine will return all messages posted by people from Pigeon et Cie. In addition to technical information, we can obtain information about the personal tastes of certain people in the company. This information will be useful when searching for passwords or improving a potential Social Engineering attempt (explained later).

#### Search Engines

Same principle as for NewsGroups, we try to retrieve other data about the target information system using a search engine. The keywords used are limited only by our imagination. Again, www.google.com is very effective especially thanks to its cache. Nevertheless, meta-engines like www.dogpile.com give relevant results by multiplying searches across several engines.

#### Social Engineering

We're departing a bit from indirect queries since this technique involves direct contact with the target. However, this approach is still not technical which is why it's not part of direct queries. Social engineering is practiced to obtain confidential information (password, technical information, phone number, IP address...) from users of the target information system. All possible and imaginable means are available (telephone, email, fax...). With identity theft and clever use of information previously gathered about people and the company, credibility is achieved along with valuable data.

#### Miscellaneous

Indirect queries are virtually unlimited; a hacker has time on their side and will check the company's website or its subsidiaries. Other data on companies and brands can be found on sites like www.societe.com or in yellow pages (phone numbers or people's names). The results obviously depend on the hacker's creativity.

### Direct Queries

The information gathered so far does not come directly from the target. We will now launch some probes in its direction and see what we can retrieve.

From the target's perspective, we'll also see how to thwart certain queries. Unlike previous steps, the target now controls the data that the tester is looking for. It's up to them to limit it to the minimum.

#### DNS Interrogation

The whois query shows us the DNS servers used. What can these reveal to us?

To get information from these servers, simply query them in a language they understand, namely the DNS protocol. Several queries are at our disposal:

- Retrieving all DNS servers with authority over the domain

```bash
host -v -t ns pigeons.fr ns1.pigeons.fr
```

```
Using domain server:
Name: ns1.pigeons.fr
Address: 10.250.149.163
Aliases:

Trying null domain
rcode = 0 (Success), ancount=3

The following answer is not verified as authentic by the server:
pigeons.fr	172800 IN	NS	ns1.pigeons.fr
pigeons.fr	172800 IN	NS	ns2.pigeons.fr
pigeons.fr	172800 IN	NS	ns.heberge.fr

Additional information:
ns1.pigeons.fr	172800 IN	A	10.250.149.163
ns2.pigeons.fr	172800 IN	A	10.250.149.165
ns.heberge.fr	345317 IN	A	10.51.3.65
```

So we have confirmation of the DNS servers used and their IP addresses.

- Retrieving mail servers (Mail eXchanger) for the domain

```bash
host -v -t mx pigeons.fr ns1.pigeons.fr
```

```
Using domain server:
Name: ns1.pigeons.fr
Address: 10.250.149.163

Aliases:

Trying null domain
rcode = 0 (Success), ancount=1

The following answer is not verified as authentic by the server:
pigeons.fr	172800 IN	MX	0 smtp1.pigeons.fr

For authoritative answers, see:
pigeons.fr	172800 IN	NS	ns1.pigeons.fr
pigeons.fr	172800 IN	NS	ns2.pigeons.fr
pigeons.fr	172800 IN	NS	ns.heberge.fr

Additional information:
smtp1.pigeons.fr	172800 IN	A	10.250.149.35
ns1.pigeons.fr		172800 IN	A	10.250.149.163
ns2.pigeons.fr		172800 IN	A	10.250.149.165
ns.heberge.fr		345239 IN	A	10.51.3.65
```

- Verification of the two previous queries

```bash
host -a pigeons.fr ns1.pigeons.fr
```

```
Using domain server:
Name: ns1.pigeons.fr
Address: 10.250.149.163
Aliases:

Trying null domain
rcode = 0 (Success), ancount=5
The following answer is not verified as authentic by the server:

pigeons.fr	172800 IN   NS	ns1.pigeons.fr
pigeons.fr	172800 IN   SOA	ns1.pigeons.fr dnsmaster.pigeons.fr(
				      2000060601    ;;;serial (version)
				      21600	    ;;refresh period
				      3600	    ;;retry refresh this often
				      3600000	    ;;expiration period
				      172800	    ;;minimum TTL
				     )

pigeons.fr	172800 IN	NS	ns2.pigeons.fr
pigeons.fr	172800 IN	NS	ns.heberge.com
pigeons.fr	172800 IN	MX	0 smtp1.pigeons.fr

For authoritative answers, see:
pigeons.fr	172800 IN	NS	ns1.pigeons.fr
pigeons.fr	172800 IN	NS	ns2.pigeons.fr
pigeons.fr	172800 IN	NS	ns.heberge.fr

Additional information:
ns1.pigeons.fr	172800 IN	A	10.250.149.163
ns2.pigeons.fr	172800 IN	A	10.250.149.165
ns.heberge.fr	345225 IN	A	10.51.3.65
smtp1.pigeons.fr	172800 IN	A	10.250.149.35
```

Now we have all the information we could easily and consistently retrieve. Let's explore the pigeons.fr domain further by looking for information on machines present on the network:

The zone transfer returns the entire configuration of the DNS server

```bash
host -l pigeons.fr ns1.pigeons.fr
```

```
Using domain server:
Name: ns1.pigeons.fr
Address: 10.250.149.163
Aliases:

pigeons.fr name server ns1.pigeons.fr
pigeons.fr name server ns2.pigeons.fr
pigeons.fr name server ns.heberge.fr

m01.pigeons.fr has address 10.51.23.226
m02.pigeons.fr has address 10.51.23.227
www2.pigeons.fr has address 10.51.23.247
m03.pigeons.fr has address 10.51.23.228
m04.pigeons.fr has address 10.51.23.229
m05.pigeons.fr has address 10.51.23.230
m10.pigeons.fr has address 10.51.23.238
m09.pigeons.fr has address 10.51.23.237
m12.pigeons.fr has address 10.250.149.162
m13.pigeons.fr has address 10.250.149.163

m14.pigeons.fr has address 10.250.149.165
m16.pigeons.fr has address 10.51.23.251
m39.pigeons.fr has address 10.51.23.249
w3.pigeons.fr has address 10.101.154.68
w5.pigeons.fr has address 10.101.154.67
w7.pigeons.fr has address 10.101.154.73
w8.pigeons.fr has address 10.101.154.77
w9.pigeons.fr has address 10.101.154.79
w5-private.pigeons.fr has address 10.101.154.70
w3-ccc.pigeons.fr has address 10.101.154.72
w3-bbb.pigeons.fr has address 10.101.154.71
www.pigeons.fr has address 10.51.23.246
```

Luck is on our side: a misconfiguration has allowed us to list all the machines in the pigeons.fr domain. Otherwise, we would have repeated the same operation on the other DNS servers because often the main server is well configured unlike the others.

- Other tests are interesting when zone transfer is impossible. For example, we can check if it's possible to retrieve internal addressing.

```bash
host -a 0.168.192.in-addr.arpa ns1.pigeon2.com
```

```
Using domain server:
Name: ns1.pigeons2.fr
Address: 10.81.144.121
Aliases:

Trying null domain
rcode = 0 (Success), ancount=4
The following answer is not verified as authentic by the server:
0.168.192.in-addr.arpa 3600 IN NS   server2.pigeons2.fr
0.168.192.in-addr.arpa 3600 IN NS   server1.pigeons2.fr
0.168.192.in-addr.arpa 3600 IN NS   echange.pigeons2.fr
0.168.192.in-addr.arpa 3600 IN SOA  server1.pigeons2.fr root.pigeons2.fr(
					      585    ;;serial (version)
					      900    ;;refresh period
					      600    ;;retry refresh this often
					      86400  ;;expiration period
					      3600   ;;minimum TTL
					     )
Additional information:
server2.pigeons2.fr       3600 IN A       192.168.0.1

server1.pigeons2.fr       3600 IN A       192.168.0.2
server1.pigeons2.fr       3600 IN A       10.81.144.121
echange.pigeons2.fr       3600 IN A       192.168.0.3
```

We've used a different target here since the previous zone transfer didn't show us any machines with private addressing (like 192.168.0.1). We see that the zone 0.168.192.in-addr.arpa is managed by the DNS server ns1.pigeons2.fr. It's therefore sufficient to test all machines in the 192.168.0.* network.

```bash
host 192.168.0.1 ns1.pigeons2.fr
```

```
1.0.168.192.IN-ADDR.ARPA        1200 IN PTR     server1.pigeons2.fr
```

All that remains is to create a small script (in Perl for example) that repeats this command for addresses from 192.168.0.1 to 192.168.0.254.

- Finally, in the case of Bind only, we can obtain its version, which is very interesting given this server's long history of remote exploits.

```
nslookup
```

```
Default Server:  ns1.pigeons.fr
Address:  10.250.149.163

> set class=chaos
> set query=txt
> version.bind
Server:  ns1.pigeons.fr
Address:  10.250.149.163

VERSION.BIND    text = "8.2.3-REL"
```

We get the Bind version, which could allow us to find a potential vulnerability in this version (a buffer overflow for example).

#### Using Ping

The indirect queries (whois) and DNS interrogation have allowed us to retrieve IP addresses and IP address ranges belonging to the target. By pinging each of these IP addresses, we'll know which ones are accessible. However, we must take into account that the presence of a firewall might prevent the machine from responding to pings. Port scanning will solve this problem by detecting the presence of the machine if it has an open port. Nmap performs this task very well:

```bash
nmap -sP 10.51.23.* -n
```

```
Starting nmap V. 2.54BETA25 ( www.insecure.org/nmap/ )
Host  (10.51.23.226) appears to be up.
Host  (10.51.23.227) appears to be up.
Host  (10.51.23.228) appears to be up.
Host  (10.51.23.229) appears to be up.
Host  (10.51.23.230) appears to be up.
Host  (10.51.23.237) appears to be up.
Host  (10.51.23.238) appears to be up.
Host  (10.51.23.246) appears to be up.
Host  (10.51.23.247) appears to be up.
Host  (10.51.23.249) appears to be up.
Host  (10.51.23.251) appears to be up.

Nmap run completed -- 256 IP addresses (11 hosts up) scanned in 2 seconds
```

#### Using Traceroute

The goal here is to obtain the IP address of an access router to the target machines. For this, simply do a traceroute to a machine on the target network:

```bash
traceroute 10.51.23.251
```

```
traceroute 10.101.154.70

 1  10.0.0.1        1.612 ms  1.443 ms  1.532 ms
 2  10.18.23.5      5.790 ms  5.454 ms  5.536 ms
 3  10.20.20.1      5.605 ms  5.453 ms  5.338 ms
 4  10.51.15.1      6.805 ms  6.437 ms  6.552 ms
 5  10.51.192.7     7.783 ms  7.246 ms  7.329 ms
 6  10.51.173.65    7.402 ms  7.246 ms  7.732 ms
 7  10.51.159.33    7.582 ms  7.844 ms  7.935 ms
 8  10.51.23.1      8.202 ms  7.639 ms  7.909 ms
 9  10.51.23.251    7.807 ms  7.633 ms  7.733 ms
```

The machine just before the destination is a router.

#### Port Scan

There is a wide variety of methods for scanning, the description of which is far beyond the scope of this article. The general principle of any scanning method is to send a packet (TCP, UDP, ICMP...) to the target machine and see what happens. Depending on the method used, the tester determines the state of the port (open, closed, filtered).

The purpose of scanning is similar to that of a scout. The tester (or hacker) thus determines the role of machines, available services, supported protocols... At the end of the operation, the following information is obtained:

- The IP addresses of the machines on the network
- The list of available services
- The list of different supported protocols (TCP, UDP, ICMP...)
- For a maximum number of machines, the state of each of its ports.

For an administrator, this step reveals access to their machines. They should also install tools allowing them to detect scans that they have not initiated themselves (iplog or portsentry for example).

There are different solutions to escape this kind of detector, by spoofing IP addresses or distributing the scan from several machines.

For example, when the source address of the packet is spoofed, the test machine remains unknown to the target machine, although it knows it has been scanned:

- kelly (192.168.1.3) the test machine 
- bosley (192.168.1.2) a quiet machine (i.e. that doesn't generate a lot of traffic)
- charly (192.168.1.1) the target machine.

To detect if a TCP port allows packets to pass through, kelly regularly sends packets to bosley. If bosley generates little traffic on the network, the id field of its packets varies little. At the same time, kelly sends TCP packets to charly with the SYN flag activated (as for a normal connection request), but putting bosley's address as the source address. Thus, charly responds to bosley with SYN-ACK packets if the port is open. bosley, which hasn't requested anything, sends a RST packet to charly to cut the connection. As a result, the id field increases because two packets are emitted (the RST and the response to kelly):

```bash
hping -r bosley
```

```
46 bytes from 192.168.1.2: flags=RA seq=1 ttl=255 id=+1 win=0 rtt=0.3 ms
46 bytes from 192.168.1.2: flags=RA seq=2 ttl=255 id=+1 win=0 rtt=0.3 ms
46 bytes from 192.168.1.2: flags=RA seq=3 ttl=255 id=+1 win=0 rtt=0.4 ms
46 bytes from 192.168.1.2: flags=RA seq=4 ttl=255 id=+1 win=0 rtt=0.3 ms

46 bytes from 192.168.1.2: flags=RA seq=5 ttl=255 id=+2 win=0 rtt=0.4 ms
46 bytes from 192.168.1.2: flags=RA seq=6 ttl=255 id=+2 win=0 rtt=0.4 ms
46 bytes from 192.168.1.2: flags=RA seq=7 ttl=255 id=+3 win=0 rtt=0.3 ms
46 bytes from 192.168.1.2: flags=RA seq=8 ttl=255 id=+2 win=0 rtt=0.4 ms
46 bytes from 192.168.1.2: flags=RA seq=9 ttl=255 id=+2 win=0 rtt=0.3 ms

46 bytes from 192.168.1.2: flags=RA seq=10 ttl=255 id=+1 win=0 rtt=0.4 ms
46 bytes from 192.168.1.2: flags=RA seq=11 ttl=255 id=+1 win=0 rtt=0.4 ms
```

Simultaneously from another terminal:

```bash
hping -a bosley -p 22 -S charly
```

```
eth0 default routing interface selected (according to /proc)
HPING charly (eth0 192.168.1.1): S set, 40 headers + 0 data bytes

--- charly hping statistic ---
6 packets tramitted, 0 packets received, 100% packet loss
```

It's normal that kelly doesn't receive any response from charly since they're sent to bosley.

On the contrary, when the target port is not open, charly doesn't emit any packet. The id field then doesn't vary:

```bash
hping -r bosley
```

```
..
46 bytes from 192.168.1.2: flags=RA seq=61 ttl=255 id=+1 win=0 rtt=0.3 ms
46 bytes from 192.168.1.2: flags=RA seq=62 ttl=255 id=+1 win=0 rtt=0.3 ms
46 bytes from 192.168.1.2: flags=RA seq=63 ttl=255 id=+1 win=0 rtt=0.4 ms
46 bytes from 192.168.1.2: flags=RA seq=64 ttl=255 id=+1 win=0 rtt=0.3 ms
..
```

The target port (hping -a bosley -p 80 -S charly) is therefore closed. Charly's logs contain a connection attempt from bosley.

To mislead a scan, it's also possible to run a honeypot. This looks like a server, tastes like a server, but it's not a real server:

```c
/* fake.c : just a socket bound to a port */
#include <stdlib.h>
#include <netinet/in.h>
#include <sys/types.h>
#include <sys/socket.h>
#include <stdio.h>
#include <unistd.h>

main(int argc,char *argv[]) {

  int port;                 //port number
  struct sockaddr_in sock;  //the socket for the server
  int sd;                   //socket descriptor

  if (argc!=2) exit(EXIT_FAILURE);
  port = htons(atoi(argv[1]));

  if ( (sd = socket(AF_INET, SOCK_STREAM, 0)) == -1) {
    perror("No socket");
    exit(EXIT_FAILURE);
  }

  sock.sin_family = AF_INET;
  sock.sin_port = port;

  sock.sin_addr.s_addr = INADDR_ANY;
  if (bind(sd, (struct sockaddr*)&sock, sizeof(struct sockaddr)) == -1) {
    perror("can't bind");
    exit(EXIT_FAILURE);
  }

  /* Let's go for LISTEN mode */
  if (listen(sd, 2) == -1) {
    perror("Bad listen");
    exit(EXIT_FAILURE);
  }

  while(1) sleep(1);
}
```

You just need to put it on the port of your choice:

```bash
gcc -o fake fake.c
./fake 21 &
```

```
[2] 3373
```

```bash
lsof -ni | grep fake
```

```
fake    3373   root    3u  IPv4 201230       TCP *:ftp (LISTEN)
```

```bash
nmap charly          
```

```
Starting nmap V. 2.54BETA22 ( www.insecure.org/nmap/ )
Interesting ports on charly (192.168.1.1):
(The 1538 ports scanned but not shown below are in state: closed)
Port       State       Service
21/tcp     open        ftp                     
22/tcp     open        ssh                     
6000/tcp   open        X11
```

We launch our fake server on port 21 (ftp). lsof reveals that a server is listening on this port 21. However, nmap (port scanner) is fooled because it just tries to open a connection on port 21. Since it succeeds, it believes it's an ftp server. The illusion works with this type of network scanner because they don't actually try to connect. Any more in-depth connection will reveal the deception, unless the fake server is refined (for example by adding banners to simulate the desired service). Note that the nc command (netcat) produces a similar result (nc -l -p 21 to listen on port 21).

This kind of defense is called a honeypot. Projects like honeynets and honeypots set up networks, or machines, designed to attract hackers to learn their techniques.

Scanning a machine always comes down to sending a packet from the test machine to the target machine, regardless of the method used. Depending on the target machine's resources (i.e. the security expected on it), an attempt with 2 packets per day is enough to detect the scan. Large disks are then needed and all packets arriving on the machine must be recorded for analysis over several days in order to reconstruct the scan.

#### OS Fingerprinting

Thanks to the network scan, we now know the active machines. We refine our knowledge by determining their operating system. This knowledge will allow us, when we have also determined the version of the daemons waiting on the target machine's ports, to search for the exploits necessary for our penetration tests.

Each OS has its own design for managing network protocols. On one hand, some fields are left to the OS (TTL, ToS, Win, DF...). On the other hand, even if RFCs define the essentials, they are not always scrupulously respected. Moreover, while they do prohibit certain packet configurations, they don't specify how to respond to them. For example, what to do with a packet containing flag 64, which is undefined? Each has its own solution.

**Default values in packets:**

By retrieving packets issued by the target, we discover the value of parameters:

- the TTL (time to live) field of outgoing packets;
- the window size;
- the DF bit (Don't Fragment);
- the TOS field (Type Of Service).
- ...

Depending on the OS, all these parameters change. A database containing their default values then facilitates identification. It's sufficient to send different packets to test the responses and then compare them to a signature database to identify the OS.

For example, the id field makes it easy to distinguish between Linux 2.2.x and 2.4.x (the command hping -1 -c 3 sends 3 packets of type 1 i.e. ICMP):

```bash
uname -a
```

```
Linux charly 2.4.4 #4 Wed May 23 10:18:08 CEST 2001 i686 unknown
```

```bash
hping -1 -c 3 charly
```

```
28 bytes from 192.168.1.1: icmp_seq=0 ttl=255 id=0 rtt=0.4 ms
28 bytes from 192.168.1.1: icmp_seq=1 ttl=255 id=0 rt
t=0.3 ms
28 bytes from 192.168.1.1: icmp_seq=2 ttl=255 id=0 rtt=0.3 ms
```

```bash
uname -a
```

```
Linux kelly 2.2.19ow1 #2 Mon May 21 12:29:48 CEST 2001 i686 unknown
```

```bash
hping -1 -c 3 kelly
```

```
28 bytes from 128.93.24.10: icmp_seq=0 ttl=255 id=4901 rtt=0.3 ms
28 bytes from 128.93.24.10: icmp_seq=1 ttl=255 id=4903 rtt=0.2 ms
28 bytes from 128.93.24.10: icmp_seq=2 ttl=255 id=4906 rtt=0.2 ms
```

**The TCP/IP Stack**

However, this method is not very reliable because OSes often allow some of these values to be modified (with sysctl under Linux or in the registry for Windows).

A more effective method is to analyze the target OS's responses to certain packets: the tester then knows the behavior of the target's TCP/IP stack, which is enough to identify the OS if the tests are well chosen.

nmap (again and again ;) uses exactly this approach when the -O option (OS identification) is activated. A database contains the typical responses according to the OS. Thus, the fingerprint of Linux kernels 2.4.0 - 2.4.5 corresponds to:

```
Contributed by  root@dexter.dynu.com
Fingerprint Linux Kernel 2.4.0 - 2.4.5 (X86)
TSeq(Class=RI%gcd=<6%SI=<2983C7E&>3DAF6%IPID=Z%TS=100HZ)
T1(DF=Y%W=16A0|7FFF%ACK=S++%Flags=AS%Ops=MNNTNW)
T2(Resp=N)
T3(Resp=Y%DF=Y%W=16A0|7FFF%ACK=S++%Flags=AS%Ops=MNNTNW)
T4(DF=Y%W=0%ACK=O%Flags=R%Ops=)
T5(DF=Y%W=0%ACK=S++%Flags=AR%Ops=)
T6(DF=Y%W=0%ACK=O%Flags=R%Ops=)
T7(DF=Y%W=0%ACK=S++%Flags=AR%Ops=)
PU(DF=Y|N%TOS=C0|0%IPLEN=164%RIPTL=148%RID=E%RIPCK=E%UCK=E|F%ULEN=134%DAT=E)
```

The tests themselves are described by the Ti lines. Reading Fyodor's article in phrack will detail them for you (phrack 54, file 9/12). However, let's briefly reveal the meaning of each:

- TSeq: describes the nature of the sequence number incrementation
- T1: TCP packet with SYN|64 flag (since 64 doesn't correspond to any flag value, the packet is "syn-bugged") to an open port
- T2: NULL TCP packet, i.e. containing no option or flag, to an open port
- T3: TCP packet with SYN|FIN|URG|PSH flags to an open port
- T4: TCP packet with ACK flag to an open port
- T5: TCP packet with SYN flag to a closed port
- T6: TCP packet with ACK flag to a closed port
- T7: TCP packet with FIN|PSH|URG flags to a closed port
- PU: UDP packet sent to a closed port to retrieve an ICMP "port unreachable" packet.

While it's often possible to modify the values of certain parameters, modifying the complete behavior of the stack is much more difficult, or even impossible with some OSes whose sources are not available.

#### Banners

The objective is simple: to know the version of the application used for a specific service. Most of the time, a simple telnet on the desired port gives us the information. Note the few services that don't deliver this information: finger (port 79), exec (port 512), login (port 513), printer (port 515).

There's also an article on this wiki: [Banners: Hiding Application Banners (Service banner faking)](./banières_:_cacher_les_banières_de_ses_applications_(service_banner_faking).html)

- FTP (port 21):

Often the version is revealed at login:

```bash
telnet ftp.pigeons.fr
```

```
Trying 192.168.14.35...
Connected to ftp.pigeons.fr
Escape character is '^]'.
220 ProFTPD 1.2.0pre9 Server (ProFTPD) [ftp1-1.pigeons.fr
```

However, some servers allow the banner to be hidden. The STAT command can save us:

```bash
telnet ftp.pigeons2.fr 21
```

```
Trying 192.168.96.24...
Connected to ftp.pigeons2.fr.
Escape character is '^]'.
220 ftp.pigeons2.fr FTP server ready.
USER ftp
331 Guest login ok, send your complete e-mail address as password.
PASS raynal@home.net
230 Guest login ok, access restrictions apply.
STAT
211-ftp.pigeons2.fr FTP server status:
   Version wu-2.6.1(1) Fri Feb 16 19:32:14 CET 2001
   Connected to bosley (192.168.1.2)
   Logged in anonymously
   TYPE: ASCII, FORM: Nonprint; STRUcture: File; transfer MODE: Stream
   No data connection
   0 data bytes received in 0 files
   0 data bytes transmitted in 0 files
   0 data bytes total in 0 files
   144 traffic bytes received in 0 transfers
   2502 traffic bytes transmitted in 0 transfers
   2696 traffic bytes total in 0 transfers
211 End of status
```

- telnet (port 23):

Even before the connection is validated by the password, the server returns the information we're looking for:

```bash
telnet 192.168.1.1
```

```
Trying 192.168.1.1...
Connected to charly (192.168.1.1).
Escape character is '^]'.

Red Hat Linux release 7.1 (Seawolf)
Kernel 2.4.4 on an i686
```

If you really want to use telnet, the -h option will only display them once the client is authenticated.

- DNS (port 53):

We've seen that it was quite simple to retrieve the version of a DNS server. However, it's possible to fake this information by modifying the options field in `/etc/named.conf`:

```
# /etc/named.conf
...
options {
           directory "/var/named";
           version "What are you doing, dude !";
        };
```

- HTTP (port 80):

The HEAD command only returns the meta-information constituting the HTTP header:

```bash
telnet minimum 80
```

```
Trying 192.168.1.1...
Connected to charly (192.168.1.1).
Escape character is '^]'.
HEAD / HTTP/1.0

HTTP/1.1 200 OK
Date: Mon, 11 Jun 2001 19:28:57 GMT
Server: Apache/1.3.19 (Unix)  (Red-Hat/Linux)
mod_ssl/2.8.1 OpenSSL/0.9.6 DAV/1.0.2 PHP/4.0.4pl1 mod_perl/1.24_01

Last-Modified: Thu, 29 Mar 2001 17:53:01 GMT
ETag: "731a-b4a-3ac3767d"
Accept-Ranges: bytes
Content-Length: 78208
Connection: close
Content-Type: text/html

Connection closed by foreign host.
```

Adding the line *ServerToken Prod* limits the information to the server name, i.e. Apache.

- portmap (port 111) and RPCs:

As we detail later in this article, the rpcinfo command provides all versions of RPC services running on the target.

- identd (port 113):

Some versions support an extension to RFC 1413: the VERSION command:

```bash
telnet charly 113
```

```
Connected to charly...
VERSION
0 , 0 : X-VERSION : pidentd 3.0.10 for Linux 2.2.5-22smp (Jul 20 2000 15:09:20)
```

Servers that support this command often also have an option to disable it.

#### Information Related to Specific Protocols

We now know exactly what is running on each system (OS, servers, server versions...). We continue our quest for information because many servers still reveal a lot about the network and its users:

- finger provides information about system users:

```bash
finger @charly
```

```
Login     Name              Tty     Idle  Login Time   Office  Office Phone
detoisien Eric Detoisien    pts/7     3d  Jun  5 09:47 (jil)
detoisien Eric Detoisien   *pts/10  160d  Jun  5 11:08 (kelly)
raynal    Frederic Raynal   tty1     10d  May 31 09:57  
raynal    Frederic Raynal   pts/1     3d  Jun  5 09:26 (:0)
raynal    Frederic Raynal   pts/3     3d  May 31 09:58 (:0)
raynal    Frederic Raynal   pts/11    3d  May 31 11:52 (:0)
raynal    Frederic Raynal   pts/7     3d  Jun  6 12:08 (:0)
raynal    Frederic Raynal   pts/2         Jun 10 09:35 (bosley)
root      root              pts/4     5d  May 31 09:58
```

In addition, it's possible to chain queries with the notation finger raynal@hots1@host2.

- The mail server: the SMTP protocol (RFC 821) defines the VRFY and EXPN commands:

VERIFY (VRFY)

This command asks the receiver to confirm that the provided arguments actually designate a user. If it is a username, the user's full name (if known to the receiver) as well as the fully qualified mailbox should be returned to the requester.

EXPAND (EXPN)

This command asks the receiver to confirm whether the associated argument identifies a mailing list, and, if so, to return the list members. The users' full names (if known) and fully qualified mailbox addresses will be returned via a multiline response.

On charly, we obtain the following information:

```
vrfy root
250 system PRIVILEGED account <root@charly.pigeons.fr>
vrfy bin
250 system librarian account <bin@charly.pigeons.fr>
vrfy web
250 Web Server manager <web@charly.pigeons.fr>
vrfy ftp
550 ftp... User unknown
vrfy raynal
250 Frederic Raynal <raynal@charly.pigeons.fr>
expn pigeons
050 pigeons... aliased to detoisien, pappy, raynal
050 /home/detoisien/.forward: line 1: forwarding to detoisien@pigeons.fr
050 /home/raynal/.forward: line 1: forwarding to \raynal@charly.pigeons.fr
050 /home/raynal/.forward: line 2: forwarding to frederic.raynal@linuxmag.fr

250-Eric Detoisien <detoisien@pigeons.fr>
250-Frederic Raynal <pappy@charly.pigeons.fr>
250-Frederic Raynal <\raynal@charly.pigeons.fr>
250-Frederic Raynal <frederic.raynal@linuxmag.fr>
```

Most SMTP servers now allow these to be disabled, which is therefore recommended ;)

- identd (formerly called auth, port 113 - RFC 1413) provides information about the identity of system users. It reveals the holder of a connection, which requires knowing the target and destination ports. For the target port, since we're on our own machine, the netstat -A inet command reveals it. As for the destination port, we've already scanned the target machine! All we need to do now is connect to each of the open ports and then ask identd who is in charge of this connection.

The result of scanning bosley is as follows:

```
7/tcp      open        echo                    
22/tcp     open        ssh                     
80/tcp     open        http                    
113/tcp    open        auth                    
664/tcp    open        unknown                 
1024/tcp   open        kdm                     
1025/tcp   open        listen                  
6000/tcp   open        X11
```

We initialize a connection on port 113 of bosley. Then, for each of the open ports, we connect with a simple telnet client (telnet bosley 664). We then ask identd to give us the desired information (the syntax of requests is <port-on-server>,<port-on-client>):

```bash
telnet bosley 113
```

```
Trying 192.168.1.2...
Connected to bosley.
Escape character is '^]'.
7,32924
7 , 32924 : USERID : OTHER :root
22,32927
22 , 32927 : USERID : OTHER :root
80,32928
80 , 32928 : ERROR : UNKNOWN-ERROR
113, 32926
113 , 32926 : USERID : OTHER :nobody
664,32930
664 , 32930 : USERID : OTHER :root
1024,32931
1024 , 32931 : USERID : OTHER :rpcuser
1025,32932
1025 , 32932 : USERID : OTHER :root
6000,32933
6000 , 32933 : USERID : OTHER :root
Connection closed by foreign host.
```

We see here some limitations of nmap. First, the error obtained on port 80 means that in fact there is no web server on bosley. Then, the kdm running with the rpcuser identity suggests that it is actually an RPC program.

Look carefully at your daemon's configuration instructions. It's often possible to replace the username with its UID, but generally, it's better to disable this server.

- portmap (sunrpc port 111) is the essential server for the proper functioning of services that rely on RPCs (NIS, NFS, rusers, rstat...). The rpcinfo command reveals what's running on a machine:

```bash
rpcinfo -p bosley
```

```
program vers proto   port
100000    2   tcp    111  portmapper
100000    2   udp    111  portmapper
100007    2   udp    661  ypbind
100007    1   udp    661  ypbind
100007    2   tcp    664  ypbind
100007    1   tcp    664  ypbind
100024    1   udp   1024  status
100024    1   tcp   1024  status
100011    1   udp    855  rquotad
100011    2   udp    855  rquotad
100005    1   udp   1025  mountd
100005    1   tcp   1025  mountd
100005    2   udp   1025  mountd
100005    2   tcp   1025  mountd
100003    2   udp   2049  nfs
100003    3   udp   2049  nfs
100021    1   udp   1026  nlockmgr
100021    3   udp   1026  nlockmgr
100021    4   udp   1026  nlockmgr
390113    1   tcp   7937
```

rpcinfo connects to port 111 of the target machine and asks it what's running on it. portmap has not provided control mechanisms. It is therefore advised to block access to this server via the firewall and tcp-wrapper. In this case, all RPC-based queries (like the next 2) will fail.

However, RPC authentication is based on the client's IP address. On a local network, it is very easy to spoof an address and thus access all available RPC services.

- A NIS server controls clients authorized to query it by a mechanism called securenets. By default, everyone can connect to the server. Any machine can then declare itself a client of such a NIS database. Once the name of the NIS database is known (it is often the same name as the server), we declare our test machine (kelly - 192.168.1.3) as a client of the NIS server (charly - 192.168.1.1):

```bash
cat /etc/yp.conf
```

```
domain charly server charly
```

```bash
ypwhich        #what is my NIS server?
```

```
charly        #bingo, it responds ;)
```

```bash
ypcat -k passwd.byname
```

```
...
fguest fguest:B4wLh7jxO1eZA:5555:5555:Compte temporaire:/home/fguest:/bin/bash
raynal raynal:YP5.ojuxdA/6.:10943:21196:Frederic Raynal:/home/raynal:/bin/bash
...
[raynal@kelly intrusion]$ ypcat -k netgroup
angels (charly,,) (bosley,,) (kelly,,) (jil,,) (sabrina,,)
```

- When the target machine exports directories via NFS, it is sometimes possible to know them using the showmount command:

```bash
showmount -e charly
```

```
/var/spool/mail  angels
/home/web/www    (everyone)
/home            angels
/opt/download    jil
```

We find here the netgroup angels discovered previously in the NIS database. Directories exported to everyone (indicated by (everyone)) are then accessible on the test machine via mount -t nfs <target>:<directory> /mnt/target.

- Still among RPCs, here are some less common servers:

- ruserd reveals users connected to a machine:

```bash
rusers -l charly
```

```
        raynal   charly:tty1 Jun 16 18:11     :51 
        raynal   charly:pts/ Jun 16 18:11          (:0)
        raynal   charly:pts/ Jun 16 18:11      :19 (:0)
        raynal   charly:pts/ Jun 16 18:11      :06 (:0)
        raynal   charly:pts/ Jun 16 18:11      :20 (:0)
        raynal   charly:pts/ Jun 16 18:59      :03 (bosley)
```

This reveals connection dates and their origins.

- rstatd generates system statistics, read by the rup command:

```bash
rup -d charly
```

```
charly      19:15pm up      1:05,  load average: 0.09 0.14 0.13
```

#### Email Sending

The header of an email is full of relevant information such as the version of the SMTP server used or even internal addressing. We obtain the path taken by the email.

```
Received: from smtp1.pigeons.fr ([xxx.xxx.xxx.xxx])
	  by front.testeur.fr (8.9.3/No_Relay+No_Spam_MGC990224) with ESMTP id OAA00632
	  for <detoisien@testeur.fr>; Wed, 18 Apr 2001 14:18:11 +0200 (MET DST)
Received: from bpigeon ([10.33.11.153]) by smtp1.pigeons.fr
          (Netscape Messaging Server 3.6)  with SMTP id AAA3A01
          for <detoisien@testeur.fr>; Wed, 18 Apr 2001 14:14:50 +0200
Message-ID: <004401c0c801$319eff40$990b210a@sit.fr>
From: "Bernard Pigeon" bernard.pigeon@pigeons.fr
To: detoisien@testeur.fr
Subject: Test 
Date: Wed, 18 Apr 2001 14:15:01 +0200
MIME-Version: 1.0
Content-Type: text/plain;
	charset="iso-8859-1"
Content-Transfer-Encoding: 8bit
X-Priority: 3
X-MSMail-Priority: Normal
X-Mailer: Microsoft Outlook Express 5.00.2919.6600
X-MimeOLE: Produced By Microsoft MimeOLE V5.00.2919.6600
```

We recover the name and address of an internal machine, the version of the SMTP server and the version of the SMTP client used. It is common to send an email to a non-existent address. Thus, the target's SMTP server automatically sends back a response with most of the information.

#### CGI Scan

The default installation of a Web server and/or poor configuration means that many scripts may be present on a Web server. A significant number of these scripts are the source of vulnerabilities. A CGI scanner allows testing for the presence of these scripts on a target server. The hacker can then use them to attack the target machine. The most well-known scanner is whisker.

```bash
perl whisker.pl -h www.pigeons.fr -i
```

```
-- whisker / v1.3.0a / rain forest puppy / ADM / wiretrip --

= - = - = - = - = - =
= Host: www.pigeons.fr
- Directory index: /

= Server: Microsoft-IIS/4.0

- Appending ::, %2E, or 0x81 to URLs may give script source
- Requesting a bogus .pl may give physical path like .idc bug does
- if perl is installed
- Security settings on directories can be bypassed if you use 8.3
- Warning: Syntax Error: names
- http://www.securityfocus.com/templates/archive.pike?list=1
- &date=1999-08-15&msg=37B5D87E.D6600553@urban-a.net
- Content-Location: http://10.51.23.246/cfdocs/index.htm

+ 200 OK: GET /cfdocs/cfmlsyntaxcheck.cfm

+ 200 OK: GET /cfide/Administrator/startstop.html
- can start/stop the server...w00h00

+ 200 OK: HEAD /iisadmpwd/aexp4b.htr
- gives domain/system name

+ 200 OK: HEAD /msadc/msadcs.dll
- RDS.  See RDS advisory, RDP9902
- Need I remind you, do not abuse, kids?
```

We have a list of potentially vulnerable scripts that could be exploited later.

#### War Dialing

War Dialing is a somewhat separate technique. Indeed, it consists of scanning an entire set of phone numbers.

Software (Toneloc, THC-Scan...) calls each phone number and detects if it's a VMB (Voice Mail Box), a fax, a person, a type of ring (busy, no answer...) or a carrier, meaning a modem (or more generally a Remote Access) that answers. The attack following this discovery focuses on discovering a login/password allowing access to the machine behind the modem.

## Using the Information

After gathering all this information about the target information system, we can plan the rest of the penetration test.

### Vulnerability Research

This step directly uses the previous data. The objective of the analysis phase is to find vulnerabilities at the network, system, and application levels of the target. These flaws can be found in public databases (such as bugtraq, which is the most well-known list) and on hacker group sites. This research results in establishing a list of exploitable vulnerabilities against the target's machines.

In cases where many machines need to be tested, we can use a vulnerability scanner (such as Nessus). This is software that automates vulnerability discovery. These are maintained in a database that can be updated online. This type of application is very useful but has its limitations. Indeed, such a scanner can report false alerts or, conversely, not detect certain vulnerabilities. Nevertheless, it can complement our list of flaws that we'll use in the last step.

### Vulnerability Exploitation

These flaws are exploited using tools available on the Internet or developed for the occasion (mostly in C or Perl). This final phase may lead to the compromise of a machine. Exploitation is very specific and obviously depends on the vulnerabilities discovered.

## Conclusion

This article focused on gathering information about the target with the objective of an external penetration test (via the Internet). This approach phase is substantially the same for each test, unlike the attack itself. Regarding internal penetration tests, the methodology remains identical but the number of vulnerabilities is often greater and attack techniques are more numerous.
