---
weight: 999
url: "/ZNC\\:_use_a_bouncer_to_get_history/"
title: "ZNC: Use a Bouncer to Get History"
description: "How to install and configure ZNC as an IRC bouncer to maintain connection and keep history when you're offline."
categories: ["Linux", "Debian"]
date: "2014-04-14T10:05:00+02:00"
lastmod: "2014-04-14T10:05:00+02:00"
tags: ["IRC", "Networking", "Servers"]
toc: true
---

![ZNC](/images/znc.avif)

{{< table "table-hover table-striped" >}}
|||
|-|-|
| **Software version** | 0.206-2 |
| **Operating System** | Debian 7 |
| **Website** | [ZNC Website](https://wiki.znc.in/ZNC) |
| **Last Update** | 14/04/2014 |
{{< /table >}}

## Introduction

A [BNC](https://en.wikipedia.org/wiki/BNC_%28software%29)[^1] (short for bouncer) is a piece of software that is used to relay traffic and connections in computer networks, much like a proxy. Using a BNC allows a user to hide the original source of the user's connection, providing privacy as well as the ability to route traffic through a specific location. A BNC can also be used to hide the true target to which a user connects.

I was fed up to launch every day a ssh to connect to an external server, on which I had a [tmux](./tmux_:_le_multiplexeur_de_terminal_remplaçant_de_screen.html) running a [weechat](./weechat_:_a_user_friendly_irc_client.html) to get IRC history. 2 problems here:

1. Sometimes, I forgot to connect on it and missed some messages
2. It was not integrated to Pidgin (what I use for any IM protocol)

So I decided to setup a bounce IRC with ZNC!

## Installation

Easy as it is integrated to Debian:

```bash
aptitude install znc
```

## Configuration

For the configuration, it's not complicated, you just have to follow the wizard on launching znc:

```bash
> znc --makeconf
[ ** ] Building new config
[ ** ]
[ ** ] First let's start with some global settings...
[ ** ]
[ ?? ] What port would you like ZNC to listen on? (1025 to 65535): 12345
[ ?? ] Would you like ZNC to listen using SSL? (yes/no) [no]:
[ ?? ] Would you like ZNC to listen using ipv6? (yes/no) [yes]: no
[ ?? ] Listen Host (Blank for all ips):
[ ok ] Verifying the listener...
```

Up to you to choose what you want for the network part.

```bash

[ ** ]
[ ** ] -- Global Modules --
[ ** ]
[ ** ] +-----------+----------------------------------------------------------+
[ ** ] | Name      | Description                                              |
[ ** ] +-----------+----------------------------------------------------------+
[ ** ] | partyline | Internal channels and queries for users connected to znc |
[ ** ] | webadmin  | Web based administration module                          |
[ ** ] +-----------+----------------------------------------------------------+
[ ** ] And 13 other (uncommon) modules. You can enable those later.
[ ** ]
[ ?? ] Load global module <partyline>? (yes/no) [no]:
[ ?? ] Load global module <webadmin>? (yes/no) [no]: yes
```

Enable at least the webadmin to help you on configuration.

```bash

[ ** ]
[ ** ] Now we need to set up a user...
[ ** ] ZNC needs one user per IRC network.
[ ** ]
[ ?? ] Username (AlphaNumeric): username
[ ?? ] Enter Password:
[ ?? ] Confirm Password:
```

Choose a username and password for you to connect to ZNC.

Now set global IRC configuration credentials and nicknames:

```bash

[ ?? ] Would you like this user to be an admin? (yes/no) [yes]:
[ ?? ] Nick [username]:
[ ?? ] Alt Nick [username_]:
[ ?? ] Ident [test]:
[ ?? ] Real Name [Got ZNC?]:
[ ?? ] Bind Host (optional):
[ ?? ] Number of lines to buffer per channel [50]: 10000
[ ?? ] Would you like to keep buffers after replay? (yes/no) [no]: yes
[ ?? ] Default channel modes [+stn]:
[ ** ]
```

Now choose which modules you want to activate:

```bash

[ ** ] -- User Modules --
[ ** ]
[ ** ] +-------------+------------------------------------------------------------------------------------------------------------+
[ ** ] | Name        | Description                                                                                                |
[ ** ] +-------------+------------------------------------------------------------------------------------------------------------+
[ ** ] | admin       | Dynamic configuration of users/settings through IRC. Allows editing only yourself if you're not ZNC admin. |
[ ** ] | chansaver   | Keep config up-to-date when user joins/parts                                                               |
[ ** ] | keepnick    | Keep trying for your primary nick                                                                          |
[ ** ] | kickrejoin  | Autorejoin on kick                                                                                         |
[ ** ] | nickserv    | Auths you with NickServ                                                                                    |
[ ** ] | perform     | Keeps a list of commands to be executed when ZNC connects to IRC.                                          |
[ ** ] | simple_away | Auto away when last client disconnects                                                                     |
[ ** ] +-------------+------------------------------------------------------------------------------------------------------------+
[ ** ] And 36 other (uncommon) modules. You can enable those later.
[ ** ]
[ ?? ] Load module <admin>? (yes/no) [no]:
[ ?? ] Load module <chansaver>? (yes/no) [no]: yes
[ ?? ] Load module <keepnick>? (yes/no) [no]: yes
[ ?? ] Load module <kickrejoin>? (yes/no) [no]: yes
[ ?? ] Load module <nickserv>? (yes/no) [no]: yes
[ ?? ] Load module <perform>? (yes/no) [no]:
[ ?? ] Load module <simple_away>? (yes/no) [no]:
[ ** ]
```

Now configure the IRC server you want to connect onto, with channels etc...:

```bash

[ ** ] -- IRC Servers --
[ ** ] Only add servers from the same IRC network.
[ ** ] If a server from the list can't be reached, another server will be used.
[ ** ]
[ ?? ] IRC server (host only):
[ ?? ] IRC server (host only): irc.freenode.net
[ ?? ] [irc.freenode.net] Port (1 to 65535) [6667]:
[ ?? ] [irc.freenode.net] Password (probably empty):
[ ?? ] Does this server use SSL? (yes/no) [no]:
[ ** ]
[ ?? ] Would you like to add another server for this IRC network? (yes/no) [no]:
[ ** ]
[ ** ] -- Channels --
[ ** ]
[ ?? ] Would you like to add a channel for ZNC to automatically join? (yes/no) [yes]:
[ ?? ] Channel name: #myassonthecommode
[ ?? ] Would you like to add another channel? (yes/no) [no]:
[ ** ]
[ ?? ] Would you like to set up another user (e.g. for connecting to another network)? (yes/no) [no]:
[ ok ] Writing config [/home/vagrant/.znc/configs/znc.conf]...
[ ** ]
[ ** ] To connect to this ZNC you need to connect to it as your IRC server
[ ** ] using the port that you supplied.  You have to supply your login info
[ ** ] as the IRC server password like this: user:pass.
[ ** ]
[ ** ] Try something like this in your IRC client...
[ ** ] /server <znc_server_ip> 12345 test:<pass>
[ ** ] And this in your browser...
[ ** ] http://<znc_server_ip>:12345/
[ ** ]
[ ?? ] Launch ZNC now? (yes/no) [yes]:
[ ok ] Opening Config [/home/vagrant/.znc/configs/znc.conf]...
[ ok ] Loading Global Module [webadmin]... [/usr/lib/znc/webadmin.so]
[ ok ] Binding to port [12345] using ipv4...
[ ** ] Loading user [username]
[ ok ] Adding Server [irc.freenode.net 6667 ]...
[ ok ] Loading Module [chansaver]... [/usr/lib/znc/chansaver.so]
[ ok ] Loading Module [keepnick]... [/usr/lib/znc/keepnick.so]
[ ok ] Loading Module [kickrejoin]... [/usr/lib/znc/kickrejoin.so]
[ ok ] Loading Module [nickserv]... [/usr/lib/znc/nickserv.so]
[ ok ] Forking into the background... [pid: 4305]
[ ** ] ZNC 0.206+deb2 - http://znc.in
```

Now it is launched! First thing to do is to kill znc to be in debug mode and understand what's wrong when you'll going to add accounts or if it fails to connect:

```bash
znc -D
```

Now open the web interface http://<znc_server_ip>:12345/ and you can finish to configure as it should be. For example in Pidgin, add as IRC server, your ZNC bouncer. Then join a channel that you previously configured into ZNC to get history etc...

If you want to get an idea of what a configuration looks like:

```bash {linenos=table,hl_lines=["12-17","43-47"]}
// WARNING
//
// Do NOT edit this file while ZNC is running!
// Use webadmin or *admin instead.
//
// Buf if you feel risky, you might want to read help on /znc saveconfig and /znc rehash.
// Also check http://en.znc.in/wiki/Configuration

AnonIPLimit  = 10
MaxBufferSize= 500
ProtectWebSessions = true
<User username>
	Pass = <sha256blalbalba>
	Nick = <nick>
	AltNick = <nick1>, <nick2>
	Ident = <username>
	RealName = <realname>
	QuitMsg = ZNC - http://znc.in
	StatusPrefix = *
	ChanModes = +stn
	Buffer = 5000
	KeepBuffer = true
	MultiClients = true
	DenyLoadMod = false
	Admin = true
	DenySetBindHost = false
	TimestampFormat = [%H:%M:%S]
	AppendTimestamp = false
	PrependTimestamp = true
	TimezoneOffset = 0.00
	JoinTries = 10
	MaxJoins = 5
	IRCConnectEnabled = true

	Allow = *

	LoadModule = admin
	LoadModule = chansaver
	LoadModule = kickrejoin
	LoadModule = perform
	LoadModule = buffextras

	Server = irc.freenode.net 6667 <freenode_password>

	<Chan #myassonthecommode>
	Key = <password>
	</Chan>
</User>
```

When finished, kill your znc and relaunch it without debug mode. That's it :-)

## References

[^1]: http://en.wikipedia.org/wiki/BNC_%28software%29
