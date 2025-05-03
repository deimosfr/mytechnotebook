---
weight: 999
url: "/EJabberd_\\:_Mise_en_place_d'un_serveur_Jabber_(messagerie_instantanée)/"
title: "EJabberd: Setting up a Jabber server (instant messaging)"
description: "A guide for installing and configuring EJabberd, a Jabber/XMPP server written in Erlang for instant messaging."
categories: ["Linux", "Debian", "Ubuntu", "Servers"]
date: "2007-03-26T20:48:00+02:00"
lastmod: "2007-03-26T20:48:00+02:00"
tags: ["Servers", "Jabber", "Chat", "XMPP", "Communication", "Instant Messaging"]
toc: true
---

## Introduction

ejabberd is a Jabber server written in Erlang, a relatively unknown language but optimized for distributed applications. ejabberd is supported by the French company Process One and is increasingly used. ejabberd handles high loads well and thanks to Erlang, it's easy to create a cluster of ejabberd servers. Its installation and administration are made easy through its web interface.

## Installation and configuration

- First, you need to install ejabberd. (On Ubuntu or Debian, a simple `apt-get install ejabberd` is sufficient)
- Edit the configuration file `/etc/ejabberd/ejabberd.cfg` and add your domain name to the hosts line (around line 94)

```
{hosts, ["example.net"]}.
```

Still in the configuration file, add your username as administrator (around line 9):

```
{acl, admin, {user, "myaccount"}}.
```

- Restart the server:

```bash
/etc/init.d/ejabberd restart
```

- Create a user account with any Jabber client
- Connect to:

```
http://localhost:5280/admin
```

Note: the login is the complete JID, including "@domain"

- That's it, you now have access to the configuration interface
