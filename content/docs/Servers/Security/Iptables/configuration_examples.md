---
weight: 999
url: "/Exemples_de_configurations/"
title: "Iptables: Configuration Examples"
description: "Various iptables configuration examples from basic to complex for firewall setup"
categories: ["Apache", "Linux"]
date: "2010-12-13T17:06:00+02:00"
lastmod: "2010-12-13T17:06:00+02:00"
tags: ["Servers", "Firewall", "Security"]
toc: true
---

## Introduction

Iptables is not very intuitive, and examples are almost essential for setting up your configuration. Here are some examples ranging from the simplest to the most complex.

## Example 1

```bash
#!/bin/sh
################################################
#                                              #
#  Basic Firewall Script                       #
#                                              #
################################################

#############
# Variables #
############
    IPTABLES=/sbin/iptables
    IF_EXT=eth0
    IP_SSH=xx.xx.xx.xx


###################
# Clear tables #
##################
   ${IPTABLES} -t mangle -F
   ${IPTABLES} -t nat -F
   ${IPTABLES} -F
   ${IPTABLES} -t mangle -X
   ${IPTABLES} -t nat -X
   ${IPTABLES} -X
   ${IPTABLES} -Z


#####################
# Default rules #
####################
  ## ignore_echo_broadcasts, TCP Syncookies, ip_forward
   echo 1 > /proc/sys/net/ipv4/icmp_echo_ignore_broadcasts

  ## Default Policy
   ${IPTABLES} -P INPUT DROP
   ${IPTABLES} -P OUTPUT DROP
   ${IPTABLES} -P FORWARD DROP

  ## Accept loopback
   ${IPTABLES} -A FORWARD -i lo -o lo -j ACCEPT
   ${IPTABLES} -A INPUT -i lo -j ACCEPT
   ${IPTABLES} -A OUTPUT -o lo -j ACCEPT

  ## REJECT connections pretending to initialize without syn
   ${IPTABLES} -A INPUT -p tcp ! --syn -m state --state NEW,INVALID -j REJECT


####################
# Special rules #
###################
### Create chains
    ${IPTABLES} -N SPOOFED
    ${IPTABLES} -N SERVICES

### Prohibit spoofed packets
    ${IPTABLES} -A SPOOFED -s 127.0.0.0/8 -j DROP
    ${IPTABLES} -A SPOOFED -s 169.254.0.0/12 -j DROP
    ${IPTABLES} -A SPOOFED -s 172.16.0.0/12 -j DROP
    ${IPTABLES} -A SPOOFED -s 192.168.0.0/16 -j DROP
    ${IPTABLES} -A SPOOFED -s 10.0.0.0/8 -j DROP

### Allowed INPUT
    ### ICMP
	## Ping (*)
        ${IPTABLES} -A INPUT -p icmp --icmp-type echo-request -j ACCEPT
    ### TCP
	## SSH (*)
	${IPTABLES} -A SERVICES -p tcp -d ${IP_SSH} --dport 22 -j ACCEPT
   ## MAIL (*)
	${IPTABLES} -A SERVICES -p tcp -d ${IP_SSH} --dport 25 -j ACCEPT


#################################
# Open ports on the firewall #
################################
    ${IPTABLES} -A OUTPUT -j ACCEPT
    ${IPTABLES} -A INPUT -m state --state ESTABLISH,RELATED -j ACCEPT
    ${IPTABLES} -A INPUT -j SPOOFED
    ${IPTABLES} -A INPUT -i ${IF_EXT} -j SERVICES
```

## Example 2

```bash
#!/bin/bash
echo Setting firewall rules...

###### Initialization Start ######

# Block all incoming connections
iptables -t filter -P INPUT DROP
iptables -t filter -P FORWARD DROP
echo - Block all incoming connections: [OK]

# Block all outgoing connections
iptables -t filter -P OUTPUT DROP
echo - Block all outgoing connections: [OK]

# Clear current tables
iptables -t filter -F
iptables -t filter -X
echo - Clearing: [OK]

# Allow SSH
iptables -t filter -A INPUT -p tcp --dport 22 -j ACCEPT
echo - Allow SSH: [OK]

# Don't break established connections
iptables -A INPUT -m state --state RELATED,ESTABLISHED -j ACCEPT
iptables -A OUTPUT -m state --state RELATED,ESTABLISHED -j ACCEPT
echo - Don't break established connections: [OK]

###### End Initialization ######

##### Begin Rules ######

# Allow DNS, FTP, HTTP, NTP requests
iptables -t filter -A OUTPUT -p tcp --dport 21 -j ACCEPT
iptables -t filter -A OUTPUT -p tcp --dport 80 -j ACCEPT
iptables -t filter -A OUTPUT -p tcp --dport 53 -j ACCEPT
iptables -t filter -A OUTPUT -p udp --dport 53 -j ACCEPT
iptables -t filter -A OUTPUT -p udp --dport 123 -j ACCEPT
echo - Allow DNS, FTP, HTTP, NTP requests: [OK]

# Allow loopback
iptables -t filter -A INPUT -i lo -j ACCEPT
iptables -t filter -A OUTPUT -o lo -j ACCEPT
echo - Allow loopback: [OK]

# Allow ping
iptables -t filter -A INPUT -p icmp -j ACCEPT
iptables -t filter -A OUTPUT -p icmp -j ACCEPT
echo - Allow ping: [OK]

# HTTP
iptables -t filter -A INPUT -p tcp --dport 80 -j ACCEPT
iptables -t filter -A INPUT -p tcp --dport 443 -j ACCEPT
iptables -t filter -A INPUT -p tcp --dport 8443 -j ACCEPT
echo - Allow Apache server: [OK]

# FTP
modprobe ip_conntrack_ftp
iptables -t filter -A INPUT -p tcp --dport 20 -j ACCEPT
iptables -t filter -A INPUT -p tcp --dport 21 -j ACCEPT
iptables -t filter -A INPUT -m state --state ESTABLISHED,RELATED -j ACCEPT
echo - Allow FTP server: [OK]

# Mail
iptables -t filter -A INPUT -p tcp --dport 25 -j ACCEPT
iptables -t filter -A INPUT -p tcp --dport 110 -j ACCEPT
iptables -t filter -A INPUT -p tcp --dport 143 -j ACCEPT
iptables -t filter -A OUTPUT -p tcp --dport 25 -j ACCEPT
iptables -t filter -A OUTPUT -p tcp --dport 110 -j ACCEPT
iptables -t filter -A OUTPUT -p tcp --dport 143 -j ACCEPT
echo - Allow Mail server: [OK]

###### End Rules ######

echo Firewall successfully updated!
```

## Example 3

```bash
# description: Firewall rules with masquerading
# probe: true
#
### BEGIN INIT INFO
# Provides: firewall_passerelle
# Required-Start: $network
# Required-Stop: $network
# Default-Start: 3 5
# Default-Stop:
# Description: Firewall rules with masquerading (configurable)
### END INIT INFO

####################################################################
# INTRODUCTION
####################################################################

## Make sure we are root
if [ ! "`id 2>&1 | egrep 'uid=0' | cut -d '(' -f1`" = "uid=0" ]; then
        echo "This script must be run by the 'root' user"
        exit 1 ## Exit the script
fi

# If iptables utility is not installed, exit with an error
# Note: the path to the IPTABLES utility may vary from one
# system to another
IPT="/sbin/iptables"
[ -x ${IPT} ] || {
        echo "Unable to find the path for iptables"
        exit 1
        }

# Internet connection interface
# This variable is mandatory
OUT="ppp0"

# If the following line is uncommented, the machine
# is not configured in gateway mode and only serves
# as a firewall
IN="eth0"  # private network interface if applicable

# Uncomment the following line to enable protocol filtering
# when using in gateway mode
#FILTRAGE="-p tcp -m multiport --destination-port 6667,5190"

# let's see how we were called
case "$1" in
        start)
                ;;
        stop)
                ${IPT} -t filter -F
                ${IPT} -t nat    -F
                ${IPT} -t filter -X
                ${IPT} -t filter -Z
                ${IPT} -t filter -P INPUT       ACCEPT
                ${IPT} -t filter -P OUTPUT      ACCEPT
                ${IPT} -t filter -P FORWARD     ACCEPT
                /bin/echo "0" > /proc/sys/net/ipv4/ip_forward
                exit 0
                ;;
        restart)
                $0 stop
                $0 start
                ;;
        *)
                echo "Usage: $0 {start|stop|restart}"
                exit 1
esac

# Load modules
modprobe ip_tables
modprobe ip_conntrack
modprobe ip_conntrack_ftp
modprobe ip_conntrack_irc
modprobe ip_conntrack_h323

# Clear all rules and chains
${IPT} -t filter -F
${IPT} -t nat    -F
${IPT} -t filter -X
${IPT} -t filter -Z

# Configure default behavior (Policy)
${IPT} -P INPUT DROP
${IPT} -P OUTPUT DROP
${IPT} -P FORWARD DROP

####################################################################
# Kernel flags
####################################################################

# Enable TCP SYN Cookie protection (repeated connection requests)
#/bin/echo "1" > /proc/sys/net/ipv4/tcp_syncookies

# Ignore ping responses
#/bin/echo "1" > /proc/sys/net/ipv4/icmp_echo_ignore_all

# Disable ICMP broadcast responses
/bin/echo "1" > /proc/sys/net/ipv4/icmp_echo_ignore_broadcasts

# Don't accept source routed packets. Attackers can use source
# routing to generate traffic pretending to be from inside your
# network, but which is routed back along the path which it came,
# namely outside, so attackers can compromise your network.
# Source routing is rarely used for legitimate purposes.
/bin/echo "0" > /proc/sys/net/ipv4/conf/all/accept_source_route

# Disable ICMP Redirect Acceptance
/bin/echo "0" > /proc/sys/net/ipv4/conf/all/accept_redirects

# Enable bad error message protection
/bin/echo "1" > /proc/sys/net/ipv4/icmp_ignore_bogus_error_responses

# To prevent IP SPOOFING, check the source address on all
# interfaces - can cause issues with asymmetric routing
# (packets take different paths in each direction)
for interface in /proc/sys/net/ipv4/conf/*/rp_filter; do
        /bin/echo "1" > ${interface}
done

# Log Spoofed Packets, Source Routed Packets, Redirect Packets
for interface in /proc/sys/net/ipv4/conf/*/log_martians; do
        /bin/echo "1" > ${interface}
done
# For dynamic IP address
echo "1" > /proc/sys/net/ipv4/ip_dynaddr

# Enable IP packet routing
# This is the main command authorizing the gateway function
if [ ${IN} ]; then
        /bin/echo "1" > /proc/sys/net/ipv4/ip_forward
else
        /bin/echo "0" > /proc/sys/net/ipv4/ip_forward
fi

####################################################################
# Rules
####################################################################

# Unlimited traffic on loopback address
${IPT} -A INPUT  -i lo -j ACCEPT
${IPT} -A OUTPUT -o lo -j ACCEPT

# Unlimited traffic on other Ethernet interfaces
# Avoid touching the public network interface
# (connected to the Internet)
for interface in /proc/sys/net/ipv4/conf/eth*; do
        VAL=`echo ${interface} | cut -c 25-`
        if [ ${VAL} != ${OUT} ]; then
                ${IPT} -A INPUT  -i ${VAL} -j ACCEPT
                ${IPT} -A OUTPUT -o ${VAL} -j ACCEPT
        fi
done

# transparent proxy: redirection rule to the proxy
# we consider that eth0 is the private network interface
# and 3128 is the proxy-cache server port
#${IPT} -t nat -A PREROUTING -p tcp -i eth0 --dport 80 -j REDIRECT --to-port 3128

# If gateway function is enabled
if [ ${IN} ]; then

        # Accept forwarding packets on the internal interface
        ${IPT} -A FORWARD -i ${IN} ${FILTRAGE} -j ACCEPT
        ${IPT} -A FORWARD -o ${IN} -j ACCEPT

        # Enable masquerading for traffic from the private
        # subnet (For fixed IP, SNAT is better)
        ${IPT} -t nat -A POSTROUTING -o ${OUT} -j MASQUERADE

fi

# Accept outgoing connections from the
# private subnet
${IPT} -A OUTPUT -o ${OUT} -j ACCEPT

#
# Add rules to authorize certain ports
# Uncomment the lines that interest you
#

## Simultaneous access to a web server and FTP
#${IPT} -A INPUT -i ${OUT} -p tcp -m state --state NEW -m multiport --destination-port 80,20,21 -j ACCEPT

## Access only to a web server
#${IPT} -A INPUT -i ${OUT} -p tcp -m state --state NEW --destination-port 80 -j ACCEPT

## gtk-gnutella
#${IPT} -A INPUT -i ${OUT} -p tcp -m state --state NEW --destination-port 23934 -j ACCEPT
#${IPT} -A INPUT -i ${OUT} -p udp -m state --state NEW --destination-port 23934 -j ACCEPT

## SSH
${IPT} -A INPUT -i ${OUT} -p tcp --destination-port 22 -j ACCEPT

## HTTP
${IPT} -A INPUT -i ${OUT} -p tcp --destination-port 80 -j ACCEPT

# Mldonkey
${IPT} -A INPUT -i ${OUT} -p tcp --destination-port 6666 -j ACCEPT
${IPT} -A INPUT -i ${OUT} -p tcp --destination-port 6682 -j ACCEPT
${IPT} -A INPUT -i ${OUT} -p tcp --destination-port 8155 -j ACCEPT

## Jabber file transfer
#${IPT} -A INPUT -i ${OUT} -p udp --destination-port 8010 -j ACCEPT

#
# End of rule addition
#

# Accept already established incoming connections
${IPT} -A INPUT -i ${OUT} -m state --state ESTABLISHED,RELATED -j ACCEPT

# If gateway function is enabled
if [ ${IN} ]; then

        # Track rejected packets on the FORWARD chain
        ${IPT} -N LOG_FWD
        ${IPT} -A LOG_FWD  -j LOG --log-level info --log-ip-options --log-prefix "Firewall FWD:"
        ${IPT} -A LOG_FWD  -j DROP

        # Log rejected packets on the FORWARD chain
        ${IPT} -A FORWARD -j LOG_FWD

fi

# Initialize tracking for rejected input packets on
# the external interface
${IPT} -N LOG_EXT
${IPT} -A LOG_EXT  -j LOG --log-level info --log-ip-options --log-prefix "Firewall IN:"
${IPT} -A LOG_EXT  -j DROP

# Log rejected input packets on the external interface
${IPT} -A INPUT -i ${OUT} -j LOG_EXT

exit 0
```

## Example 4

```bash
#!/bin/bash

###################################################
## ARCHITECTURE FOR A 4-INTERFACE FIREWALL       ##
##						 ##
##		   INTERNET			 ##
##		      | 			 ##
##	DMZ--------FIREWALL--------SERVER ZONE ##
##		      |				 ##
##		     LAN			 ##
##						 ##
###################################################

###################################################
## REQUIRED IPTABLES MODULES                     ##
###################################################
MODULES_IPTABLES="ip_tables \
                  ipt_string \
		  ip_conntrack \
                  ip_conntrack_ftp \
		  ip_nat_ftp"			# Iptables modules loaded at startup
INTERNET="ppp0" 				# Internet device (multiple devices possible)
INTERNET_NAT="ppp0"                             # Internet device used for NAT (only 1 device possible)
DMZ=""   	        			# DMZ device (public IPs, servers accessible from internet)
ZONE_SERVEURS="eth2"				# Server zone device (private IPs, servers accessible internally)
LAN="eth1 eth3"					# Intranet device (multiple devices possible)

PAQUETS_ICMP_AUTHORISES="0 3 4 5 8 11 12"	# ICMP packets authorized to travel between different networks
PING_FLOOD="1/s"				# Number of PING authorized per second
LOG_FLOOD="1/s"

PROTOCOLES_AUTHORISES="47"			# Protocols authorized to pass through the firewall

MASQ_LAN="YES"                                  # Masquerade the LAN
MASQ_DMZ="NO"                                   # Masquerade the DMZ
MASQ_ZONE_SERVEURS="YES"                        # Masquerade the server zone

PORTS_TCP_INTERNET_AUTHORISES="53"		# TCP ports of the firewall accessible from the internet
PORTS_UDP_INTERNET_AUTHORISES="53"		# UDP ports of the firewall accessible from the internet
PORTS_TCP_DMZ_AUTHORISES=""			# TCP ports of the firewall accessible from the DMZ
PORTS_UDP_DMZ_AUTHORISES=""			# UDP ports of the firewall accessible from the DMZ
PORTS_TCP_ZONE_SERVEURS_AUTHORISES="53 113"	# TCP ports of the firewall accessible from the server zone
PORTS_UDP_ZONE_SERVEURS_AUTHORISES="53 113"	# UDP ports of the firewall accessible from the server zone
PORTS_TCP_LAN_AUTHORISES="53 113 22"		# TCP ports of the firewall accessible from the LAN
PORTS_UDP_LAN_AUTHORISES="53 113 22"		# UDP ports of the firewall accessible from the LAN

PORTS_TCP_SORTIE_REFUSES="6346 \
			  7777 \
			  8888 \
			  6699 \
			  6000"                 # TCP ports forbidden for output from the firewall
PORTS_UDP_SORTIE_REFUSES="6346 \
			  7777 \
			  8888 \
			  6699 \
			  6000"                 # UDP ports forbidden for output from the firewall

RESEAUX_LAN="192.168.10.0/24 \
	     192.168.30.0/24"  			# Networks composing the LAN
RESEAUX_DMZ=""                                  # Networks composing the DMZ
RESEAUX_ZONE_SERVEURS="192.168.50.0/24"         # Networks composing the server zone

NAT_TCP_NET=" 80.13.192.105:80>192.168.50.100:8080 "  # NAT => IP_FIREWALL:PORT_FIREWALL>IP_INTERNAL:PORT_INTERNAL
NAT_UDP_NET=""  				# NAT => IP_FIREWALL:PORT_FIREWALL>IP_INTERNAL:PORT_INTERNAL

MOTS_CLES="root admin"                          # Keywords to log

MOTS_CLES_INTERDITS="mp3>192.168.10.117 \
		     MP3>192.168.10.117 \
		     ogg>192.168.10.117 \
		     OGG>192.168.10.117"       	# Forbidden keywords ;-) KEYWORD>RECIPIENT_IP
IP_INTERDITES=" 66.28.48.0/24 \
		66.28.49.0/24"			# Addresses blocked from entry

###################################################
## SCRIPT VARIABLES (DO NOT EDIT)               ##
###################################################

IPTABLES=`which iptables`
MODPROBE=`which modprobe`
VERT="\033[32m"
JAUNE="\033[33m"
GRAS="\033[1m"
NORMAL="\033[m"
ROUGE="\033[31m"


###################################################
## VERIFY IPTABLES PRESENCE                      ##
###################################################

echo -en "${GRAS}Verifying IPTABLES presence:${NORMAL}"
if [ -z ${IPTABLES} ] ;then
    echo -e "\t\t${ROUGE}FAILED${NORMAL}\n"
    exit 1
else
    echo -e "\t\t${VERT}OK${NORMAL}"
fi


###################################################
## VERIFY MODPROBE PRESENCE                      ##
###################################################

echo -en "${GRAS}Verifying MODPROBE presence:${NORMAL}"
if [ -z ${MODPROBE} ] ;then
    echo -e "\t\t${ROUGE}FAILED${NORMAL}\n"
    exit 1
else
    echo -e "\t\t${VERT}OK${NORMAL}\n"
fi


###################################################
## LOADING IPTABLES MODULES                      ##
###################################################

for module in ${MODULES_IPTABLES} ;do
    echo -e "${GRAS}Loading module ${module}:${NORMAL}\t\t\t${VERT}OK${NORMAL}"
    ${MODPROBE} ${module}
done
echo -e "\n"

###################################################
## BASIC FIREWALL CONFIGURATION USING            ##
## /proc FILESYSTEM                              ##
###################################################

###################################################
## ENABLE IP FORWARDING (routing)                ##
###################################################

echo -en "${GRAS}${JAUNE}Enabling ip forwarding:${NORMAL}"
if [ -e /proc/sys/net/ipv4/ip_forward ] ; then
    echo 1 > /proc/sys/net/ipv4/ip_forward
    echo -e "\t\t\t\t${VERT}OK${NORMAL}"
else
    echo -e "\t\t\t\t${ROUGE}FAILED${NORMAL}\n"
    exit 1
fi

###################################################
## Protection against SYN FLOOD                  ##
###################################################

echo -en "${GRAS}${JAUNE}Protection against SYN/FLOOD:${NORMAL}"
if [ -e /proc/sys/net/ipv4/tcp_syncookies ] ; then
    echo 1 > /proc/sys/net/ipv4/tcp_syncookies
    echo -e "\t\t\t${VERT}OK${NORMAL}"
else
    echo -e "\t\t\t${ROUGE}FAILED${NORMAL}"
fi

###################################################
## Defragment packets before forwarding them     ##
## Useful for masquerading                       ##
###################################################

echo -en "${GRAS}${JAUNE}Packet refragmentation:${NORMAL}"
if [ -e /proc/sys/net/ipv4/ip_always_defrag ] ; then
    echo 1 > /proc/sys/net/ipv4/ip_always_defrag
    echo -e "\t\t\t\t${VERT}OK${NORMAL}"
else
    echo -e "\t\t\t\t${ROUGE}FAILED${NORMAL}"
fi

###################################################
## Don't respond to ICMP packets                 ##
## sent to broadcast                             ##
###################################################

echo -en "${GRAS}${JAUNE}Insensitivity to ICMP packets sent to broadcast:${NORMAL}"
if [ -e /proc/sys/net/ipv4/icmp_echo_ignore_broadcasts ] ; then
    echo 1 > /proc/sys/net/ipv4/icmp_echo_ignore_broadcasts
    echo -e "\t${VERT}OK${NORMAL}"
else
    echo -e "\t${ROUGE}FAILED${NORMAL}"
fi

###################################################
## Ignore ICMP errors from hosts                 ##
## on the network reacting poorly to frames      ##
## sent to what they perceive as                 ##
## the broadcast address                         ##
###################################################

if [ -e /proc/sys/net/ipv4/icmp_ignore_bogus_error_responses ] ; then
    echo 1 > /proc/sys/net/ipv4/icmp_ignore_bogus_error_responses
fi

###################################################
## Reverse Path Filtering                        ##
## Only route packets belonging to               ##
## our networks                                  ##
###################################################

echo -e "${GRAS}${JAUNE}Enabling Reverse Path Filtering:${NORMAL}\t\t\t${VERT}OK${NORMAL}\n"
for f in /proc/sys/net/ipv4/conf/*/rp_filter; do
    echo 1 > $f
done


###################################################
## CLEAR OLD RULES                               ##
###################################################

echo -en "${GRAS}${JAUNE}Clearing old rules:${NORMAL}"
${IPTABLES} -t filter -F INPUT
${IPTABLES} -t filter -F OUTPUT
${IPTABLES} -t filter -F FORWARD
${IPTABLES} -t nat    -F PREROUTING
${IPTABLES} -t nat    -F OUTPUT
${IPTABLES} -t nat    -F POSTROUTING
${IPTABLES} -t mangle -F PREROUTING
${IPTABLES} -t mangle -F OUTPUT
echo -e "\t\t\t${VERT}OK${NORMAL}"

###################################################
## RESET CHAINS                                  ##
###################################################

echo -en "${GRAS}${JAUNE}Resetting chains:${NORMAL}"
${IPTABLES} -t filter -Z
${IPTABLES} -t nat    -Z
${IPTABLES} -t mangle -Z
echo -e "\t\t\t\t${VERT}OK${NORMAL}"

###################################################
## SET DEFAULT POLICY                            ##
###################################################

echo -en "${GRAS}${JAUNE}Setting default policy:${NORMAL}"
${IPTABLES} -t filter -P INPUT   DROP
${IPTABLES} -t filter -P OUTPUT  ACCEPT
${IPTABLES} -t filter -P FORWARD DROP
echo -e "\t\t${VERT}OK${NORMAL}\n"

###################################################
## KEYWORDS TO LOG                               ##
###################################################

if [ "${MOTS_CLES}" != "" ] ;then
    echo -ne "${GRAS}${JAUNE}Enabling keyword-based logging system:${NORMAL}"
    for mot in ${MOTS_CLES} ;do
	${IPTABLES} -A INPUT -m string --string "${mot}" -j LOG --log-level info --log-prefix "${mot}: "
	${IPTABLES} -A FORWARD -m string --string "${mot}" -j LOG --log-level info --log-prefix "${mot}: "
    done
    echo -e "\t\t${VERT}OK${NORMAL}"
fi

###################################################
## Block entry of certain addresses              ##
## via the firewall for tcp and udp              ##
###################################################

if [ "${IP_INTERDITES}" != "" ] ;then
    echo -e "${GRAS}${JAUNE}Blocking entry of certain addresses:${NORMAL}\t\t${VERT}OK${NORMAL}"

    for adr in ${IP_INTERDITES} ;do
	${IPTABLES} -t filter -A FORWARD -p tcp -s ${adr} -j DROP
	${IPTABLES} -t filter -A FORWARD -p udp -s ${adr} -j DROP
    done
fi

###################################################
## Block outgoing of certain ports via           ##
## the firewall for tcp                          ##
###################################################

if [ "${PORTS_TCP_SORTIE_REFUSES}" != "" ] ;then
    echo -e "${GRAS}${JAUNE}Blocking outgoing TCP ports:${NORMAL}\t\t${VERT}OK${NORMAL}"

    for port_no in ${PORTS_TCP_SORTIE_REFUSES} ;do
	${IPTABLES} -t filter -A FORWARD -p tcp --dport ${port_no} -j DROP
	${IPTABLES} -t filter -A OUTPUT -p tcp -o ${INTERNET} --dport ${port_no} -j DROP
    done
fi

###################################################
## Block outgoing of certain ports via           ##
## the firewall for udp                          ##
###################################################

if [ "${PORTS_TCP_SORTIE_REFUSES}" != "" ] ;then
    echo -e "${GRAS}${JAUNE}Blocking outgoing UDP ports:${NORMAL}\t\t${VERT}OK${NORMAL}"

    for port_no in ${PORTS_TCP_SORTIE_REFUSES} ;do
	${IPTABLES} -t filter -A FORWARD -p udp --dport ${port_no} -j DROP
	${IPTABLES} -t filter -A OUTPUT -p udp -o ${INTERNET} --dport ${port_no} -j DROP
    done
fi
###################################################
## Block passage of certain keywords             ##
###################################################

if [ "${MOTS_CLES_INTERDITS}" != "" ] ;then
    echo -e "${GRAS}${JAUNE}Blocking passage of certain keywords:${NORMAL}\t\t${VERT}OK${NORMAL}"

    for mot_cles in ${MOTS_CLES_INTERDITS} ;do
	mot=`echo ${mot_cles} | sed 's/>.*//g'`
	ip=`echo ${mot_cles} | sed 's/.*>//g'`

	${IPTABLES} -A INPUT -m string --string "${mot}" -d ${ip} -j DROP
	${IPTABLES} -A FORWARD -m string --string "${mot}" -d ${ip} -j DROP
    done
fi

###################################################
## Allow ICMP packets                            ##
###################################################

if [ "${PAQUETS_ICMP_AUTHORISES}" != ""  ] ;then
    echo -e "${GRAS}${JAUNE}Allowing certain ICMP packets:${NORMAL}\t\t${VERT}OK${NORMAL}"

    for icmp_no in ${PAQUETS_ICMP_AUTHORISES} ;do
	${IPTABLES} -t filter -A INPUT   -p icmp --icmp-type ${icmp_no} -m limit --limit ${PING_FLOOD} -j ACCEPT
	${IPTABLES} -t filter -A FORWARD -p icmp --icmp-type ${icmp_no} -m limit --limit ${PING_FLOOD} -j ACCEPT
	${IPTABLES} -t filter -A OUTPUT  -p icmp --icmp-type ${icmp_no} -m limit --limit ${PING_FLOOD} -j ACCEPT
    done
fi


###################################################
## Allow certain protocols to pass              ##
###################################################

if [ "${PROTOCOLES_AUTHORISES}" != ""  ] ;then
    echo -e "${GRAS}${JAUNE}Allowing certain protocols:${NORMAL}\t\t\t${VERT}OK${NORMAL}"

    for protocole_no in ${PROTOCOLES_AUTHORISES} ;do
	${IPTABLES} -t filter -A INPUT   -p ${protocole_no} -j ACCEPT
	${IPTABLES} -t filter -A FORWARD -p ${protocole_no} -j ACCEPT
    done
fi


###################################################
## Allow connections already established before  ##
## launch of this script                         ##
###################################################

echo -e "${GRAS}${JAUNE}Allowing already established connections:${NORMAL}\t\t${VERT}OK${NORMAL}"
${IPTABLES} -t filter -A INPUT   -m state --state ESTABLISHED,RELATED -j ACCEPT
${IPTABLES} -t filter -A FORWARD -m state --state ESTABLISHED,RELATED -j ACCEPT
${IPTABLES} -t filter -A OUTPUT  -m state --state ESTABLISHED,RELATED -j ACCEPT


###################################################
## Allow LocalHost connections                    ##
###################################################

echo -e "${GRAS}${JAUNE}Allowing localhost connections:${NORMAL}\t${VERT}OK${NORMAL}"
${IPTABLES} -t filter -A INPUT   -s 127.0.0.1 -d 127.0.0.1 -j ACCEPT
${IPTABLES} -t filter -A FORWARD -s 127.0.0.1 -d 127.0.0.1 -j ACCEPT
${IPTABLES} -t filter -A OUTPUT  -s 127.0.0.1 -d 127.0.0.1 -j ACCEPT

###################################################
## Allow TCP connections on the                  ##
## internet device                              ##
###################################################

if [ "${INTERNET}" != "" ] ;then
    for internet_device in ${INTERNET} ;do
	if [ "${PORTS_TCP_INTERNET_AUTHORISES}" != "" ] ;then
	    echo -e "${GRAS}${JAUNE}TCP connections on internet interface ${internet_device}:${NORMAL}\t\t${VERT}OK${NORMAL}"

	    for port_no in ${PORTS_TCP_INTERNET_AUTHORISES} ;do
		${IPTABLES} -t filter -A INPUT -p tcp -i ${internet_device} --dport ${port_no} -j ACCEPT

		if [ "0${port_no}" == "021" ] ;then
		    ${IPTABLES} -t filter -A INPUT -p tcp -i ${internet_device} --sport 20 --dport 1024:65535 ! --syn -j ACCEPT
		fi
	    done
	fi
    done
fi


###################################################
## Allow UDP connections on the                  ##
## internet device                              ##
###################################################

if [ "${INTERNET}" != "" ] ;then
    for internet_device in ${INTERNET} ;do
	if [ "${PORTS_UDP_INTERNET_AUTHORISES}" != "" ] ;then
	    echo -e "${GRAS}${JAUNE}UDP connections on internet interface ${internet_device}:${NORMAL}\t\t${VERT}OK${NORMAL}"

	    for port_no in ${PORTS_UDP_INTERNET_AUTHORISES} ;do
		${IPTABLES} -t filter -A INPUT -p udp -i ${internet_device} --dport ${port_no} -j ACCEPT
	    done
	fi
    done
fi


###################################################
## Allow TCP connections on the                  ##
## DMZ device                                    ##
###################################################

if [ "${DMZ}" != "" ] ;then
    for dmz_device in ${DMZ} ;do
	if [ "${PORTS_TCP_DMZ_AUTHORISES}" != "" ] ;then
	    echo -e "${GRAS}${JAUNE}TCP connections on DMZ interface ${dmz_device}:${NORMAL}\t\t${VERT}OK${NORMAL}"

	    for port_no in ${PORTS_TCP_DMZ_AUTHORISES} ;do
		${IPTABLES} -t filter -A INPUT -p tcp -i ${dmz_device} --dport ${port_no} -j ACCEPT

		if [ "0${port_no}" == "021" ] ;then
		    ${IPTABLES} -t filter -A INPUT -p tcp -i ${dmz_device} --sport 20 --dport 1024:65535 ! --syn -j ACCEPT
		fi
	    done
	fi
    done
fi

###################################################
## Allow UDP connections on the                  ##
## DMZ device                                    ##
###################################################

if [ "${DMZ}" != "" ] ;then
    for dmz_device in ${DMZ} ;do
	if [ "${PORTS_UDP_DMZ_AUTHORISES}" != "" ] ;then
	    echo -e "${GRAS}${JAUNE}UDP connections on DMZ interface ${dmz_device}:${NORMAL}\t\t${VERT}OK${NORMAL}"

	    for port_no in ${PORTS_UDP_DMZ_AUTHORISES} ;do
		${IPTABLES} -t filter -A INPUT -p udp -i ${dmz_device} --dport ${port_no} -j ACCEPT
	    done
	fi
    done
fi


###################################################
## Allow TCP connections on the                  ##
## server zone device                           ##
###################################################

if [ "${ZONE_SERVEURS}" != "" ] ;then
    for zone_serveurs_device in ${ZONE_SERVEURS} ;do
	if [ "${PORTS_TCP_ZONE_SERVEURS_AUTHORISES}" != "" ] ;then
	    echo -e "${GRAS}${JAUNE}TCP connections on server zone interface ${zone_serveurs_device}:${NORMAL}\t${VERT}OK${NORMAL}"

	    for port_no in ${PORTS_TCP_ZONE_SERVEURS_AUTHORISES} ;do
		${IPTABLES} -t filter -A INPUT -p tcp -i ${zone_serveurs_device} --dport ${port_no} -j ACCEPT

		if [ "0${port_no}" == "021" ] ;then
		    ${IPTABLES} -t filter -A INPUT -p tcp -i ${zone_serveurs_device} --sport 20 --dport 1024:65535 ! --syn -j ACCEPT
		fi
	    done
	fi
    done
fi

###################################################
## Allow UDP connections on the                  ##
## server zone device                           ##
###################################################

if [ "${ZONE_SERVEURS}" != "" ] ;then
    for zone_serveurs_device in ${ZONE_SERVEURS} ;do
	if [ "${PORTS_UDP_ZONE_SERVEURS_AUTHORISES}" != "" ] ;then
	    echo -e "${GRAS}${JAUNE}UDP connections on server zone interface ${zone_serveurs_device}:${NORMAL}\t${VERT}OK${NORMAL}"

	    for port_no in ${PORTS_UDP_ZONE_SERVEURS_AUTHORISES} ;do
		${IPTABLES} -t filter -A INPUT -p udp -i ${zone_serveurs_device} --dport ${port_no} -j ACCEPT
	    done
	fi
    done
fi

###################################################
## Allow TCP connections on the                  ##
## LAN device                                    ##
###################################################

if [ "${LAN}" != "" ] ;then
    for lan_device in ${LAN} ;do
	if [ "${PORTS_TCP_LAN_AUTHORISES}" != "" ] ;then
	    echo -e "${GRAS}${JAUNE}TCP connections on LAN interface ${lan_device}:${NORMAL}\t\t${VERT}OK${NORMAL}"

	    for port_no in ${PORTS_TCP_LAN_AUTHORISES} ;do
		${IPTABLES} -t filter -A INPUT -p tcp -i ${lan_device} --dport ${port_no} -j ACCEPT

		if [ "0${port_no}" == "021" ] ;then
		    ${IPTABLES} -t filter -A INPUT -p tcp -i ${lan_device} --sport 20 --dport 1024:65535 ! --syn -j ACCEPT
		fi
	    done
	fi
    done
fi


###################################################
## Allow UDP connections on the                  ##
## LAN device                                    ##
###################################################

if [ "${LAN}" != "" ] ;then
    for lan_device in ${LAN} ;do
	if [ "${PORTS_UDP_LAN_AUTHORISES}" != "" ] ;then
	    echo -e "${GRAS}${JAUNE}UDP connections on LAN interface ${lan_device}:${NORMAL}\t\t${VERT}OK${NORMAL}"

	    for port_no in ${PORTS_UDP_LAN_AUTHORISES} ;do
		${IPTABLES} -t filter -A INPUT -p udp -i ${lan_device} --dport ${port_no} -j ACCEPT
	    done
	fi
    done
fi
echo -e ""


###################################################
## Masquerade the LAN                            ##
###################################################

if [ "${MASQ_LAN}" = "YES" -o "${MASQ_LAN}" = "yes" ] ;then
    echo -e "${GRAS}${JAUNE}Enabling Masquerading for the LAN:${NORMAL}\t\t${VERT}OK${NORMAL}"

    for reseau in ${RESEAUX_LAN} ;do
	${IPTABLES} -t nat -A POSTROUTING -s ${reseau} -o ${INTERNET} -j MASQUERADE
	${IPTABLES} -t filter -A FORWARD -s ${reseau} -j ACCEPT
    done
fi


###################################################
## Masquerade the DMZ                            ##
###################################################

if [ "${MASQ_DMZ}" = "YES" -o "${MASQ_DMZ}" = "yes" ] ;then
    echo -e "${GRAS}${JAUNE}Enabling Masquerading for the DMZ:${NORMAL}\t\t${VERT}OK${NORMAL}"

    for reseau in ${RESEAUX_DMZ} ;do
	${IPTABLES} -t nat -A POSTROUTING -s ${reseau} -o ${INTERNET} -j MASQUERADE
	${IPTABLES} -t filter -A FORWARD -s ${reseau} -j ACCEPT
    done
fi

###################################################
## Masquerade the server zone                    ##
###################################################

if [ "${MASQ_ZONE_SERVEURS}" = "YES" -o "${MASQ_ZONE_SERVEURS}" = "yes" ] ;then
    echo -e "${GRAS}${JAUNE}Enabling Masquerading for the server zone:${NORMAL}\t${VERT}OK${NORMAL}"

    for reseau in ${RESEAUX_ZONE_SERVEURS} ;do
	${IPTABLES} -t nat -A POSTROUTING -s ${reseau} -o ${INTERNET} -j MASQUERADE
	${IPTABLES} -t filter -A FORWARD -s ${reseau} -j ACCEPT
    done
fi


###################################################
## Enable TCP NAT                                ##
###################################################

if [ "${NAT_TCP_NET}" != "" ] ;then
    echo -e "${GRAS}${JAUNE}Enabling TCP network address translation:${NORMAL}\t\t${VERT}OK${NORMAL}"

    for translation in ${NAT_TCP_NET} ;do
	srcport=`echo ${translation} | sed 's/>.*//g'|cut -d : -f 2`
	srchost=`echo ${translation} | sed 's/:.*//g'`
	desthost=`echo ${translation} | sed 's/.*>//g'| cut -d : -f 1`
	destport=`echo ${translation} | sed 's/.*://g'`

	${IPTABLES} -t nat -A PREROUTING -p tcp -i ${INTERNET_NAT} -d ${srchost} --dport ${srcport} -j DNAT --to ${desthost}:${destport}
	${IPTABLES} -A FORWARD -p tcp -i ${INTERNET_NAT} -d ${desthost} --dport ${destport} -j ACCEPT
    done
fi


###################################################
## Enable UDP NAT                                ##
###################################################

if [ "${NAT_UDP_NET}" != "" ] ;then
    echo -e "${GRAS}${JAUNE}Enabling UDP network address translation:${NORMAL}\t\t${VERT}OK${NORMAL}"

    for translation in ${NAT_UDP_NET} ;do
	srcport=`echo ${translation} | sed 's/>.*//g'|cut -d : -f 2`
	srchost=`echo ${translation} | sed 's/:.*//g'`
	desthost=`echo ${translation} | sed 's/.*>//g'| cut -d : -f 1`
	destport=`echo ${translation} | sed 's/.*://g'`

	${IPTABLES} -t nat -A PREROUTING -p udp -i ${INTERNET_NAT} -d ${srchost} --dport ${srcport} -j DNAT --to ${desthost}:${destport}
	${IPTABLES} -A FORWARD -p udp -i ${INTERNET_NAT} -d ${desthost} --dport ${destport} -j ACCEPT
    done
fi

###################################################
## FUCK nimda and codered:)                     ##
###################################################

echo -e "${GRAS}${JAUNE}Protection against Nimda and codered:${NORMAL}\t\t\t${VERT}OK${NORMAL}"
${IPTABLES} -I INPUT -j DROP -m string -p tcp -s 0.0.0.0/0 --string "c+dir"
${IPTABLES} -I INPUT -j DROP -m string -p tcp -s 0.0.0.0/0 --string "c+tftp"
${IPTABLES} -I INPUT -j DROP -m string -p tcp -s 0.0.0.0/0 --string "cmd.exe"
${IPTABLES} -I INPUT -j DROP -m string -p tcp -s 0.0.0.0/0 --string "default.ida"
${IPTABLES} -I FORWARD -j DROP -m string -p tcp -s 0.0.0.0/0 --string "c+dir"
${IPTABLES} -I FORWARD -j DROP -m string -p tcp -s 0.0.0.0/0 --string "c+tftp"
${IPTABLES} -I FORWARD -j DROP -m string -p tcp -s 0.0.0.0/0 --string "cmd.exe"
${IPTABLES} -I FORWARD -j DROP -m string -p tcp -s 0.0.0.0/0 --string "default.ida"


###################################################
## Enable logging                                ##
###################################################

echo -ne "${GRAS}${JAUNE}Enabling logging system:${NORMAL}"
${IPTABLES} -t filter -A INPUT -p tcp -m limit --limit ${LOG_FLOOD} -j LOG --log-level info --log-prefix "INPUT TCP DROPPED: "
${IPTABLES} -t filter -A INPUT -p udp -m limit --limit ${LOG_FLOOD} -j LOG --log-level info --log-prefix "INPUT UDP DROPPED: "
${IPTABLES} -t filter -A INPUT -p icmp -m limit --limit ${LOG_FLOOD} -j LOG --log-level info --log-prefix "INPUT ICMP DROPPED: "
${IPTABLES} -t filter -A INPUT -f -m limit --limit ${LOG_FLOOD} -j LOG --log-level info --log-prefix "INPUT FRAGMENT DROPPED: "
${IPTABLES} -t filter -A INPUT -p all -m limit --limit ${LOG_FLOOD} -j LOG --log-level info --log-prefix "INPUT PROTOCOL DROPPED: "

${IPTABLES} -t filter -A FORWARD -p tcp -m limit --limit ${LOG_FLOOD} -j LOG --log-level info --log-prefix "FORWARD TCP DROPPED: "
${IPTABLES} -t filter -A FORWARD -p udp -m limit --limit ${LOG_FLOOD} -j LOG --log-level info --log-prefix "FORWARD UDP DROPPED: "
${IPTABLES} -t filter -A FORWARD -p icmp -m limit --limit ${LOG_FLOOD} -j LOG --log-level info --log-prefix "FORWARD ICMP DROPPED: "
${IPTABLES} -t filter -A FORWARD -f -m limit --limit ${LOG_FLOOD} -j LOG --log-level info --log-prefix "FORWARD FRAGMENT DROPPED: "
${IPTABLES} -t filter -A FORWARD -p all -m limit --limit ${LOG_FLOOD} -j LOG --log-level info --log-prefix "FORWARD PROTOCOL DROPPED: "
echo -e "\t\t\t\t${VERT}OK${NORMAL}"
```

## Example 5

```bash
#!/bin/bash

#-------------------------------------------------------------------------
# Essentials
#-------------------------------------------------------------------------

IPTABLES='/sbin/iptables';
modprobe nf_conntrack_ftp

#-------------------------------------------------------------------------
# Physical and virtual interfaces definitions
#-------------------------------------------------------------------------

# Interfaces
wan_if="eth0";
vpn_if="tap0";

#-------------------------------------------------------------------------
# Networks definitions
#-------------------------------------------------------------------------

# Networks
wan_ip="x.x.x.x";
lan_net="192.168.90.0/24";
vpn_net="192.168.20.0/24";

# IPs
ed_ip="192.168.90.1";
banzai_ip="192.168.90.2";

#-------------------------------------------------------------------------
# Global Rules input / output / forward
#-------------------------------------------------------------------------

# Flushing tables
$IPTABLES -F
$IPTABLES -X
$IPTABLES -t nat -F

# Define default policy
$IPTABLES -P INPUT DROP
$IPTABLES -P OUTPUT ACCEPT
$IPTABLES -P FORWARD ACCEPT

$IPTABLES -A INPUT -j ACCEPT -d $lan_net;
$IPTABLES -A INPUT -j ACCEPT -m state --state ESTABLISHED,RELATED

#-------------------------------------------------------------------------
# Allow masquerading for VE
#-------------------------------------------------------------------------

# Activating masquerade to get Internet from VE
$IPTABLES -t nat -A POSTROUTING -o $wan_if -s $lan_net -j MASQUERADE

# Activating masquerade to get VPN access from VE
$IPTABLES -t nat -A POSTROUTING -o tap0 -j MASQUERADE

#-------------------------------------------------------------------------
# Allow ports on CT
#-------------------------------------------------------------------------

# Allow ICMP
$IPTABLES -A INPUT -j ACCEPT -p icmp

# SSH access
$IPTABLES -A INPUT -j ACCEPT -p tcp --dport 22

#-------------------------------------------------------------------------
# Redirections for incoming connections (wan)
#-------------------------------------------------------------------------

# HTTP access
$IPTABLES -t nat -A PREROUTING -p tcp --dport 80 -d $wan_ip -j DNAT --to-destination $ed_ip:80

# HTTPS access
$IPTABLES -t nat -A PREROUTING -p tcp --dport 443 -d $wan_ip -j DNAT --to-destination $ed_ip:443
```

## Example 6

```bash
#!/bin/bash
clear
echo "############################## Firewall Rules ###################################"
# Enable routing
echo 1 > /proc/sys/net/ipv4/ip_forward

echo "Initializing rules"

# Clear all rules
iptables -F
iptables -t nat -F
# Apply basic policies
# Allow internal traffic
iptables -P INPUT ACCEPT
iptables -P OUTPUT ACCEPT
iptables -P FORWARD ACCEPT
# Block all entry and exit
iptables -t nat -P PREROUTING DROP
iptables -t nat -P POSTROUTING DROP

# Internal traffic allowed
echo "Internal traffic"
iptables -t nat -I POSTROUTING -o lo -j ACCEPT
iptables -t nat -I PREROUTING -i lo -j ACCEPT


# Network card definitions
WEB="ppp0"
DMZ="eth2"
COM="eth1"
STA="eth0"
PPP="ppp0"

# IP network definitions
NET_COM="10.0.0.0/8"
NET_STA="192.168.2.0/24"
NET_DMZ="172.16.1.0/24"

# Server definitions for external connection to servers
REMOTE="192.168.2.8:81"
FICS="172.16.1.6/32"
EXC="172.16.1.3/32"
DC="172.16.1.1/32"
MAIL="172.16.1.3:25"
HTTP="172.16.1.4:80"
EMULE="172.16.1.4:5555"
RDP="172.16.1.4:3389"
PPTP="172.16.1.1"
VUE="192.168.2.8/32"
MAILWEB="172.16.1.3/32"
LINUX2="172.16.1.7/32"
LINUX="192.168.2.5/32"
YONI="192.168.2.62/32"
WIFI="192.168.2.7/32"

# Common Rules
# ====================== >>>> Masquerade all networks to the internet
echo "Applying common rules"
# All outgoing traffic to Internet is masqueraded
iptables -t nat -I POSTROUTING -s $NET_STA -d $NET_DMZ -j MASQUERADE

# Squid must always go out to internal clients

iptables -t nat -I POSTROUTING -p tcp --sport 3128 -d $NET_STA -j ACCEPT

iptables -t nat -I POSTROUTING -o $WEB -j MASQUERADE

iptables -t nat -I POSTROUTING -o $COM -j MASQUERADE

iptables -t nat -A POSTROUTING -s $NET_STA -o $COM -j DROP
iptables -t nat -A POSTROUTING -s $NET_DMZ -o $COM -j DROP


iptables -I INPUT -i $WEB -m state --state ESTABLISHED -j ACCEPT
iptables -I OUTPUT -m state --state ESTABLISHED -j ACCEPT

iptables -I INPUT -i $COM -m state --state ESTABLISHED -j ACCEPT

# Allow standard internal routing

# ====================== >>>> DHCP
echo "Allowing DHCP traffic"
iptables -t nat -A PREROUTING -p udp --dport 67:68 -j ACCEPT
iptables -t nat -A POSTROUTING -p udp --sport 67:68 -j ACCEPT

echo "Local DNS to SRV-DC"
iptables -t nat -A PREROUTING -p udp --sport 53 -i $DMZ -s "172.16.1.1/32" -j ACCEPT
iptables -t nat -A POSTROUTING -p udp --dport 53 -o $DMZ -d "172.16.1.1/32" -j ACCEPT

iptables -t nat -A PREROUTING -p tcp --sport 53 -i $DMZ -s "172.16.1.1/32" -j ACCEPT
iptables -t nat -A POSTROUTING -p tcp --dport 53 -o $DMZ -d "172.16.1.1/32" -j ACCEPT

# Network access rules
# 1 --> DMZ
	echo "====================== >>>> Rules for commercial machines"
	echo "Daytime rules"
	echo "Access based on time"

	iptables -t nat -I PREROUTING -i $DMZ -m time --timestart 08:45 --timestop 17:45 \
		--days Mon,Tue,Wed,Thu,Fri -p tcp -m multiport --ports 20,21,80,3128,1863,110,119,25,8080,9000 -j ACCEPT

	echo "Nighttime rules"
	iptables -t nat -I PREROUTING -i $DMZ -m time --timestart 17:46 --timestop 23:59 \
		--days Mon,Tue,Wed,Thu -p tcp -j ACCEPT
        iptables -t nat -I PREROUTING -i $DMZ -m time --timestart 00:00 --timestop 08:44 \
               --days Mon,Tue,Wed,Thu,Fri -p tcp -j ACCEPT

	# No limits on weekends
	echo "No limits on weekends"
       iptables -t nat -I PREROUTING -i $DMZ -m time --timestart 17:46 --timestop 23:59 \
               --days Fri  -p tcp -j ACCEPT
       iptables -t nat -I PREROUTING -i $DMZ -m time --timestart 00:00 --timestop 23:59 \
               --days Sat,Sun  -p tcp -j ACCEPT


	iptables -t nat -I PREROUTING -i $DMZ -p udp --dport 53 -j ACCEPT
	#====>>>>>    Transparent proxy for commercial users
	iptables -t nat -I PREROUTING -p tcp -i $DMZ --dport 80 -j REDIRECT --to-port 3128
	iptables -t nat -A PREROUTING -p tcp -i $DMZ --dport 443 -j ACCEPT
#	iptables -t nat -A PREROUTING -p tcp -i $DMZ --dport 443 -j REDIRECT --to-port 3128

	echo "====================== >>>> Rules for classrooms"
# 2 --> Classroom <-> DMZ
#	A - FICS2
	echo "	Classroom -> SRV-FICS2"
	iptables -t nat -A PREROUTING -p tcp -d $FICS -j ACCEPT
	iptables -t nat -A PREROUTING -p tcp --dport 80 -d $MAILWEB -j ACCEPT
	iptables -t nat -A POSTROUTING -s $NET_DMZ -d $NET_STA -j ACCEPT
# 3 --> Classroom <-> Internet
#	A - HTTP
	echo "	Classroom -> Internet with Squid"
	iptables -t nat -A PREROUTING -p tcp -i $STA --dport 80 -j REDIRECT --to-port 3128
	iptables -t nat -A PREROUTING -p tcp -i $STA --dport 443 -j ACCEPT

#	C - DNS
	iptables -t nat -A PREROUTING -p udp -i $STA --dport 53 -j ACCEPT
#	iptables -t nat -A PREROUTING -p tcp -i $STA --dport 53 -j ACCEPT

	echo "====================== >>>> Rules for Internet to internal network"
# 5 --> Internet <--> DMZ
#	A - SMTP
		echo "	SMTP"
		iptables -t nat -I PREROUTING -i $WEB -p tcp --dport 25 -j DNAT --to-destination $MAIL
		iptables -t nat -I POSTROUTING -o $DMZ -d $EXC -p tcp --dport 25 -j ACCEPT

#	B - WEB
		echo "	WEB"
		iptables -t nat -A PREROUTING -i $WEB -p tcp --dport 80 -j DNAT --to-destination $HTTP
		iptables -t nat -A POSTROUTING -o $DMZ -d "172.16.1.4/32" -p tcp --dport 80 -j MASQUERADE
#	B' - EMULE
		iptables -t nat -A PREROUTING -i $WEB -p tcp --dport 5555 -j DNAT --to-destination $EMULE
		iptables -t nat -A POSTROUTING -o $DMZ -d "172.16.1.4/32" -p tcp --dport 5555 -j MASQUERADE
		iptables -t nat -A PREROUTING -i $WEB -p udp --dport 5555 -j DNAT --to-destination $HTTP
		iptables -t nat -A POSTROUTING -o $DMZ -d "172.16.1.4/32" -p udp --dport 5555 -j MASQUERADE
#	C - PPTP
		echo "	PPTP"
		iptables -t nat -A PREROUTING -i $WEB -p 47 -j DNAT --to-destination $PPTP
		iptables -t nat -A POSTROUTING -o $DMZ -p 47 -j MASQUERADE
		iptables -t nat -A PREROUTING -i $WEB -p tcp --dport 1723 -j DNAT --to-destination $PPTP
		iptables -t nat -A POSTROUTING -o $DMZ -p tcp --dport 1723 -j MASQUERADE

#	D - SSH from outside or only for authorized internal machines
		echo "	SSH from Internet"
		iptables -t nat -A PREROUTING -s 172.16.1.0/24 -p tcp --dport 22 -j ACCEPT

#       E - FTP
                echo "  FTP IS DISABLED!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!"
                #iptables -t nat -A PREROUTING -i $WEB -p tcp --dport 20 -j DNAT --to-destination "172.16.1.4:20"
                #iptables -t nat -A PREROUTING -i $WEB -p tcp --dport 21 -j DNAT --to-destination "172.16.1.4:21"
                #iptables -t nat -A POSTROUTING -o $DMZ -d "172.16.1.4/32" -p tcp --dport 21 -j MASQUERADE

                #iptables -t nat -A POSTROUTING -o $DMZ -d "172.16.1.4/32" -p tcp --dport 20 -j MASQUERADE

#	F - RDP
		echo "  RDP"
		iptables -t nat -A PREROUTING -i $WEB -p tcp --dport 3389 -j DNAT --to-destination $RDP
                iptables -t nat -A POSTROUTING -o $DMZ -d "172.16.1.4/32" -p tcp --dport 3389 -j MASQUERADE
#	G - SNMP
		echo "  SNMP"
		iptables -t nat -A POSTROUTING -p tcp --dport 161 -j ACCEPT
		iptables -t nat -A POSTROUTING -p udp --dport 161 -j ACCEPT
		iptables -t nat -A POSTROUTING -p udp --dport 162 -j ACCEPT


# 6 --> Access by MAC address
	echo "====================== >>>> Special rules for internal users"
	echo "	Yoni"
#	A - Yoni
		iptables -t nat -I PREROUTING -m mac --mac-source '00:00:F0:82:58:AF' -j ACCEPT
		iptables -t nat -I PREROUTING -m mac --mac-source '00:04:23:76:63:10' -j ACCEPT

#	A' - OlivierG
		iptables -t nat -I PREROUTING  -s 192.168.2.69/32 -m mac --mac-source '00:0d:60:75:b8:75' -j ACCEPT
		iptables -t nat -I PREROUTING  -s 192.168.2.39/32 -m mac --mac-source '00:0C:F1:43:14:05' -j ACCEPT

#	B - Olivier all
	echo "	OlivierC"
	iptables -t nat -I PREROUTING -s 192.168.2.63/32 -m mac --mac-source '00:90:F5:1E:51:A1' -j ACCEPT
	iptables -t nat -I PREROUTING -s 172.16.1.63/32 -m mac --mac-source '00:90:F5:1E:51:A1' -j ACCEPT
	# Wifi Olivier
	iptables -t nat -I PREROUTING -m mac --mac-source '00:A0:C5:B1:DD:15' -j ACCEPT
#	C - Steeve all
	echo "	Steeve"
		iptables -t nat -I PREROUTING -s 192.168.2.64/32 -m mac --mac-source '00:08:02:04:fa:d7' -j ACCEPT
		iptables -t nat -I PREROUTING -m mac --mac-source '00:08:02:04:fa:d7' -j ACCEPT
#	D - Portable Compaq
	echo "	Portable Compaq"
		iptables -t nat -I PREROUTING -s 192.168.2.65/32 -m mac --mac-source '00:50:8B:FA:B9:5B' \
		-p tcp -m multiport --ports 443,110,25,119 -j ACCEPT
                iptables -t nat -I PREROUTING -s 192.168.2.65/32 -m mac --mac-source '00:50:8B:FA:B9:5B' \
                -p udp --dport 53 -j ACCEPT
                iptables -t nat -I PREROUTING -s 192.168.2.65/32 -m mac --mac-source '00:50:8B:FA:B9:5B' \
		-d $NET_DMZ -j ACCEPT

#	D' Portable Toshiba
	echo "	Portable Toshiba"
		iptables -t nat -I PREROUTING -s 192.168.2.67/32 -m mac --mac-source '00:01:02:E7:36:E3' \
                -p tcp -m multiport --ports 443,110,25,119 -j ACCEPT
		iptables -t nat -I PREROUTING -s 192.168.2.67/32 -m mac --mac-source '00:01:02:E7:36:E3' \
                -p udp --dport 53 -j ACCEPT
		iptables -t nat -I PREROUTING -s 192.168.2.67/32 -m mac --mac-source '00:01:02:E7:36:E3'  \
                -d $NET_DMZ -j ACCEPT

#	E - VUE Server
	echo "	 VUE Server"
	iptables -t nat -I PREROUTING -s $VUE -m mac --mac-source '00:0c:6e:c5:42:6c' -j ACCEPT
	iptables -t nat -I PREROUTING -i $DMZ -d $VUE -j ACCEPT
#	F- Linux Server Ground Floor
	echo "  Linux Server"
	iptables -t nat -I PREROUTING -i $DMZ -d $LINUX -j ACCEPT

#	F - Quentin Laptop
	echo "   Quentin"
		iptables -t nat -I PREROUTING -s 172.16.1.65/32 -m mac --mac-source '00:0b:db:a1:c2:a5' -j ACCEPT
		iptables -t nat -I PREROUTING -s 192.168.2.65/32 -m mac --mac-source '00:0b:db:a1:c2:a5' -j ACCEPT
		iptables -t nat -I PREROUTING -s 192.168.2.65/32 -m mac --mac-source '00:a0:c5:b1:da:f8' -j ACCEPT

#	F - Eva Laptop
	echo "  Eva is grounded"
		#iptables -t nat -I PREROUTING -m mac --mac-source '00:02:3f:13:bb:21' -j ACCEPT

#	G  - Lionel Laptop
	echo "  Lionel"
		iptables -t nat -I PREROUTING \
		-m mac --mac-source '00:0D:60:2C:12:95' -j ACCEPT

#	H  - WIFI ROUTER
	echo "  WIFI ROUTER"
		iptables -t nat -I PREROUTING \
		-m mac --mac-source '00:0F:66:33:20:12' -j ACCEPT
echo "############################## END ===> Firewall Rules ###################################"
	iptables -t nat -I PREROUTING -s $LINUX -j ACCEPT
	iptables -t nat -I PREROUTING -s $LINUX2 -j ACCEPT

iptables -t nat -I PREROUTING -s 172.16.1.1/32 -j ACCEPT
iptables -t nat -I PREROUTING -s 172.16.1.2/32 -j ACCEPT
iptables -t nat -I PREROUTING -s 172.16.1.3/32 -j ACCEPT
iptables -t nat -I PREROUTING -s 172.16.1.4/32 -j ACCEPT
iptables -t nat -I PREROUTING -s 172.16.1.5/32 -j ACCEPT
iptables -t nat -I PREROUTING -s 192.168.2.95/32 -j ACCEPT
```
