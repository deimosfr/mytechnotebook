---
weight: 999
url: "/Protocole_BGP_et_sa_sécurité/"
title: "BGP Protocol and its Security"
description: "A detailed explanation of the Border Gateway Protocol (BGP), including its usage between autonomous systems and its security aspects."
categories: ["Linux", "Network"]
date: "2009-02-03T19:16:00+02:00"
lastmod: "2009-02-03T19:16:00+02:00"
tags: ["Réseaux", "Network", "RFC 4271"]
toc: true
---

[Border Gateway Protocol (BGP)](https://fr.wikipedia.org/wiki/Border_Gateway_Protocol) is a routing protocol used particularly on the Internet. Its objective is to exchange networks (IP address + mask) with its neighbors through TCP sessions (on port 179).

BGP is used to transmit networks between autonomous systems (AS) because it is the only protocol that supports very large volumes of data.

BGP supports classless routing and uses route aggregation to limit the size of the routing table. Since 1994, version 4 of the protocol has been used on the Internet, with previous versions being considered obsolete. Its specifications are described in [RFC 4271](https://tools.ietf.org/html/rfc4271) A Border Gateway Protocol 4 (BGP-4).

Besides the Internet, very large private IP networks can use BGP, for example to connect local networks using OSPF.

Most end users of the Internet have only one connection to an Internet service provider. In this case, BGP is unnecessary because a default route is sufficient. However, a company that is redundantly connected to multiple ISPs could obtain its own autonomous system number and establish BGP sessions with each provider.

[BGP Protocol Security](/pdf/sécurité_du_protocole_bgp.pdf)
