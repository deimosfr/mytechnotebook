---
weight: 999
url: "/Introduction_aux_IDS_et_à_SNORT/"
title: "Introduction to IDS and SNORT"
description: "A comprehensive guide to Intrusion Detection Systems (IDS) and SNORT, covering installation, configuration, and best practices for network security monitoring."
categories: ["Security", "Monitoring", "Networking"]
date: "2008-05-24T16:43:00+02:00"
lastmod: "2008-05-24T16:43:00+02:00"
tags: ["IDS", "SNORT", "Security", "Network", "Monitoring"]
toc: true
---

## Introduction to IDS and SNORT

Intrusion Detection Systems (IDS) are increasingly common and play an important and growing role in the security of today's information systems. SNORT is one of these intrusion detection systems with the distinction of being open source software under the GPL license.

*Note: This document was written in 2006. Some versions and configurations of the software mentioned may differ from those shown here; please refer to the official project websites in case of issues.*

Many commercial solutions have emerged with the appearance of this previously underdeveloped market. However, there are different open source alternatives, which are often the basis for commercial products and whose performance and reliability are just as good as their competitors.

While the technical aspects are sometimes complex, the organizational aspects are no less important in the implementation of these systems.

Indeed, the deployment of IDS requires a qualified team, which can be outsourced to a specialized IT service company for installation and configuration, but also a team that will be responsible for monitoring the alerts raised, refining the configuration to minimize false positives, and implementing possible countermeasures to detected attacks.

This often underestimated aspect leads some companies to deploy IDS in their infrastructure and then leave them in place without monitoring, which makes them useless, or even dangerous for the information system.

There are mainly 2 types of IDS:

- HIDS (Host IDS), or system IDS, which locally analyze the integrity of machines via control of system calls, various event logs, file system integrity, etc.

- NIDS (Network IDS) which are based solely on captured network traffic and which most often operate on pattern matching rules that can be triggered through known signatures.

NIDS analyze flows in real time (allowing for latency) and can reassemble frames (IP fragmentation), reconstruct flows (stream4) and manage states (stateful).

Probes constitute the active agents of the IDS: they isolate the relevant information they capture, the events that are - a priori - suspicious, then send them up to a centralized system: the alert concentrator.

The alert format must be known to the concentrator. This issue gave rise to the IDMEF (Intrusion Detection Message Exchange Format) format, based on XML, which allows the standardization of different alerts raised.

In a large infrastructure, several probes are necessary. Care must then be taken to carefully choose the locations of these on the network in order to avoid unnecessarily overloading the alerts raised.

**Presentation of Snort**

Snort ([click](https://www.snort.org/)), which falls into the NIDS category, is a powerful, configurable open source product that meets the main constraints of intrusion detection.

It analyzes traffic in real time and can decode many protocols, perform pattern matching on captured packets, and also detect port scans.

It can send alerts via syslog, to a special file (socket or pipe for example), to a remote database (MySQL, MsSQL or other), or even via direct alerts (WinPopup).

Like all open source projects, it has a large community, allowing for very quick responses to problems encountered or requests for information. The documentation provided is abundant and detailed, allowing you to get the most out of this powerful tool.

**Placement in the architecture**

Probes must be placed carefully on the segments of the network to be audited: they would be useless, for example, on an administration segment where access would be considered secure and where the information returned would "blur" the analysis.

Throughput should be taken into account as the number of alerts is often proportional to the traffic of the segment.

A risk study (what should I protect? against whom? What is the degree of sensitivity of the information passing through this network segment? Etc.) must precede deployment in order to best determine the locations of future probes.

Probes do not require an IP configuration, except if they require remote administration or communication with the alert concentrator.

The classic case of an architecture including a demilitarized zone (DMZ) and an Intranet: The probe downstream of the firewall will collect a large number of alerts and is not always necessary. It could be harmful during alert analysis.

The probe located in the DMZ is the most sensitive and should receive special attention. Indeed, in the event of a machine being compromised from the outside, it will be the one to raise the first alerts generated if detection occurs.

As for the probe placed in the intranet, its usefulness is proportional to the number of users and the degree of trust granted. It can be particularly useful for detecting worms, viruses, or backdoors, but also in the event of user workstations being compromised.

**Physical aspects**

Probes must be able to analyze all of the traffic of the segment on which they act.

Connected to a hub where packets are transmitted on all of its ports, they can be directly connected to any port.

In the case of use on a switch, they must be placed on replication ports (mirroring port).

It is possible to increase the stealth of the probes by making the physical connection to the network unidirectional, thus allowing only passive listening to traffic.

They will then be "transparent" but will require physical access for maintenance needs or for local analysis of alerts.

## Installation and configuration

The installation of Snort, although relatively simple, does not allow the IDS to be used directly without prior modification of the default configuration. However, SNORT can be set up rudimentarily in less than an hour.

*Note: This document was written in 2006. Some versions and configurations of the software used according to these versions may be different from those mentioned; please refer to the official websites of the projects in question in case of problems.*

Network capture management is based on the pcap library (WinPcap under Windows). Snort therefore requires the latter to function correctly.

The latest stable version can be downloaded from ([click](https://winpcap.polito.it/)) and is distributed under a BSD-type license; for more information on this subject (cf. ([click](https://winpcap.polito.it/misc/copyright.htm)).

The trivial installation of this library (Auto-installer) will not be detailed here.

The package to download is "Winpcap auto installer". The developer will turn to the devel package containing the source codes of the library as well as more detailed technical documentation. Depending on the version and type of OS, installation may require a restart.

Then comes the installation of Snort with the installer ([click](https://www.snort.org/dl/binaries/win32/)).

Based on the same type of installer as Winpcap, one of the only modifiable options is the installation path of Snort (by default c:\snort).

The tree structure is presented with the following folders:

bin/ – containing the executable

contrib/ – containing the classic README file of open source projects.

doc/ – containing detailed documentation for installation and configuration. It also contains a FAQ (Frequently Asked Questions)

etc/ – containing the configuration files of the IDS, as well as a highly commented default configuration.

log/ – which will contain log files and any pcap recordings related to alerts.

rules/ – containing rules for different alerts. It is also advised to update them regularly.

schemas/ – containing schemas for different databases (MySQL, MsSQL, Oracle and PostgreSQL). It is indeed possible to send alerts to a remote database for specific processing.

**Configuration**

Most of the configuration is stored in the etc/snort.conf file which will be called when Snort starts.

The file can be taken as is and linearly modified to adapt it to the needs and specificities of the server. Be careful, the "#" are comments and will therefore not be taken into account in the configuration.

The first two variables to define are:

HOME_NET which defines the target network(s) to audit;

EXTERNAL_NET which contains the networks considered hostile. In many cases, it will be set to "any" to not trust any network (case of a probe connected directly to the Internet for example).

Alerts from HOME_NET will also be sent.

To only consider incoming alerts, it is possible to define it as follows:

var EXTERNAL_NET !$HOME_NET

Then we specify the different machines on the network (SMTP, SNMP, HTTP, etc.), and the ports on which services are listening.

RULE_PATH then allows us to define the path of the rules files containing the triggering rules for the different alerts.

Next come the preprocessors which make it possible to follow connections, reassemble packets, decode certain types of protocols, etc.

As indicated in the default configuration, the syntax is:

```
preprocessor <preprocessor name>: <options>
```

Here are the main preprocessors of Snort:

- flow for tracking IP packets (src port, dst port, etc.)

- frag2 for reassembling IP packets (defragmentation)

- stream4 for TCP frame reassembly: stateful

- http_inspect for the HTTP decoder (field normalization, etc.)

- rpc_decode for the normalization and reassembly of RPC packets

- bo for Back Orifice backdoor traffic

- telnet_decode for reassembling Telnet and FTP traffic

- flow-portscan sf Portscan for port scan detection

- arpspoof for detecting L2 attacks such as ARP cache poisoning

Detailed options are available in the default configuration file or on the online documentation ([click](https://www.snort.org/docs/snort_htmanuals/htmanual_232/node11.html) & [click](https://www.snort.org/docs/)).

Outputs allow you to precisely configure the logging of alerts and pcap recordings (network traces).

Several formats are available, the generic syntax is:

```
output <name>: <options>
```

Syslog logging: local or remote escalation of alerts to a syslogd logging server.

Local:

```
output alert_syslog: LOG_AUTH LOG_ALERT
```

Remote:

```
output alert_syslog: host=192.168.0.100:514, LOG_AUTH LOG_ALERT
```

Pcap save: network traces related to alerts in pcap format, readable among others by the tools tcpdump or ethereal:

```
output log_tcpdump: tcpdump.log
```

SQL/Oracle logging: sending alerts to an SQL/Oracle server:

```
output database: log, mysql, user=user password=pwd dbname=db host=localhost
```

Where the server type can be adapted: mysql, mssql, postgresql, oracle or obdc, as well as the user and password.

"Snort" / unified logging: binary format specific to Snort and allowing to increase the overall performance of recordings and limit the size of the output file:

```
output alert_unified: filename snort.alert, limit 128
```

Specific logging: depending on certain alerts. For example for a connect back type backdoor on TCP port 31337:

```
ruletype backdoor { type alert output log_tcpdump: trojan31337.log }
```

With the associated rule:

```
backdoor tcp $HOME_NET any -> $EXTERNAL_NET 31337 (msg:"Backdoor 31337 detected"; flags:A+;)
```

The last point to take into account in this configuration is the activation of the different rule classes: the names of the rules are very explicit.

It is also possible to create your own set of rules: the syntax is very simple and many additional rules are available, for example on mailing lists such as bugtraq or directly on the Snort site.

Example of an alert detecting the exploitation of a bug on a Web/CGI software package:

```
alert tcp $EXTERNAL_NET any -> $HTTP_SERVERS $HTTP_PORTS (msg:"WEB-CGI progiciel_compta exploitation attempt"; flow:to_server,established; uricontent:"/prog_compta"; content:"../"; content:"%00"; classtype:web-application-attack;)
```

The configuration finished, Snort starts simply with:

```
snort -d -l ../log -c ../etc/snort.conf -i [interface]
```

To determine the identifier of the interface, it is possible:

- to use Ethereal to retrieve the identifier of the interface,

- to retrieve the identifier directly in the Windows registry with the key HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Services\Tcpip\Parameters\Interfaces.

Example of launch:

```
snort -d -l ../log -c ../etc/snort.conf -i \Device\NPF_{2D7E7BAE-3442-4FA1-A154-5172CAC0F038}
```

## Example configuration file

The entire file used in this folder for the configuration of the Open-Source intrusion detection system SNORT.

*Note: This document was written in 2006. Some versions and configurations of the software used according to these versions may be different from those mentioned; please refer to the official websites of the projects in question in case of problems.*

```bash
#--------------------------------------------------
# etc/snort.conf
#--------------------------------------------------

# C networks 192.168.0.0 and 192.168.1.0

var HOME_NET [192.168.0.0/24,192.168.1.0/24]

# All traffic except networks considered "safe"

var EXTERNAL_NET !$HOME_NET

# DNS servers

var DNS_SERVERS $HOME_NET

# SMTP servers

var SMTP_SERVERS $HOME_NET

# HTTP servers

var HTTP_SERVERS $HOME_NET

# SQL servers

var SQL_SERVERS $HOME_NET

# Telnet servers

var TELNET_SERVERS $HOME_NET

# SNMP servers

var SNMP_SERVERS $HOME_NET

# Listening ports on HTTP servers

var HTTP_PORTS 80

# Ports on which shellcodes will be monitored
# Port 80 (HTTP) too often causes false positives

var SHELLCODE_PORTS !80

# Oracle Port

var ORACLE_PORTS 1521

# AIM servers (IM)

var AIM_SERVERS [64.12.24.0/23,64.12.28.0/23,64.12.161.0/24]

# Path containing pattern matching rules

var RULE_PATH ../rules

# Preprocessors

preprocessor flow: stats_interval 0 hash 2
preprocessor frag2
preprocessor stream4: disable_evasion_alerts
preprocessor stream4_reassemble
preprocessor http_inspect: global iis_unicode_map unicode.map 1252
preprocessor http_inspect_server: server default profile apache ports { 80 }
preprocessor rpc_decode: 111 32771
preprocessor bo
preprocessor telnet_decode
preprocessor sfportscan: proto { all } memcap { 10000000 } sense_level { low } ignore_scanners { 192.168.1.1 192.168.1.10 }

# Outputs

output log_tcpdump: tcpdump.log
output database: log, mysql, user=user password=pwd dbname=snort host=localhost
output alert_unified: filename snort.alert, limit 128
output log_unified: filename snort.log, limit 128

# Activated rules

include $RULE_PATH/local.rules
include $RULE_PATH/bad-traffic.rules
include $RULE_PATH/exploit.rules
include $RULE_PATH/scan.rules
include $RULE_PATH/ftp.rules
include $RULE_PATH/telnet.rules
include $RULE_PATH/dos.rules
include $RULE_PATH/ddos.rules
include $RULE_PATH/web-cgi.rules
include $RULE_PATH/web-misc.rules
include $RULE_PATH/web-client.rules
include $RULE_PATH/web-php.rules
include $RULE_PATH/sql.rules
include $RULE_PATH/smtp.rules
include $RULE_PATH/imap.rules
include $RULE_PATH/other-ids.rules
include $RULE_PATH/web-attacks.rules
include $RULE_PATH/backdoor.rules
include $RULE_PATH/shellcode.rules
include $RULE_PATH/policy.rules
include $RULE_PATH/virus.rules
```

## Conclusion and webography

This final part of the Open-Source IDS SNORT dossier briefly addresses subjects such as probe security and monitoring. Also find all the links related to this dossier.

### Probe security

A probe is a sensitive element due to its function on the network: an attacker will try to hide the traces of his attack by compromising it.

It must therefore be given special attention regarding updates, especially Winpcap and Snort, and elementary security measures for it.

### Monitoring

Several interfaces allow you to analyze the alerts sent by the probes. ACID (Analysis Console for Intrusion Databases) is one of the most powerful and also allows you to correlate alerts from different probes on the network.

It is accessible via a simple web browser and has many configuration options.

Intrusion detection under Windows is not left behind thanks to ported applications such as Snort and Winpcap.

If the installation and configuration are relatively accessible, the work prior to the deployment of these systems should not be neglected: a precise study (budget/personnel, risk analysis, etc.) is necessary.

IDS, although increasingly efficient, can be bypassed, particularly via hidden channels or even shellcode encodings, and constitute only a complementary means of security!

It would be illusory to think that an IDS could prevent an attack.

If properly configured and used, it can considerably reduce the damage caused by the compromise of machines on the network by reducing the reaction time of competent teams.

The main difficulty remains the analysis of event logs because an attack can quickly be drowned in a flow of false positives.

## Resources
- http://www.secuobs.com/news/11052008-snort_ids.shtml
