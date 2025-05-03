---
weight: 999
url: "/Weechat_\\:_a_user_friendly_IRC_client/"
title: "Weechat: A User Friendly IRC Client"
description: "Guide to installing and configuring Weechat IRC client on Linux systems with a focus on Freenode setup and channel autoconnect features."
categories:
  - "Linux"
  - "Debian"
date: "2013-07-17T08:02:00+02:00"
lastmod: "2013-07-17T08:02:00+02:00"
tags:
  - "IRC"
  - "Chat"
  - "Terminal"
toc: true
---

![Weechat](/images/weechat_logo.avif)

{{< table "table-hover table-striped" >}}
|||
|-|-|
| **Software version** | 0.3.8 |
| **Operating System** | Debian 7 |
| **Website** | [Weechat Website](https://www.weechat.org/) |
| **Last Update** | 17/07/2013 |
{{< /table >}}

## Introduction

WeeChat is a fast, light and extensible chat client. It runs on many platforms (including Linux, BSD and Mac OS).

WeeChat is:

- modular: a lightweight core with plugins around
- multi-protocols: IRC and Jabber (other soon)
- extensible: C plugins and scripts (Perl, Python, Ruby, Lua, Tcl and Scheme)
- free software: released under GPLv3 license
- fully documented: user's guide, API, FAQ,.. translated in many languages

Development is very active, and bug fixes are very fast!

## Installation

The installation part is really easy:

```bash
aptitude install weechat-curses
```

To launch it:

```bash
weechat-ncurses
```

## Configuration

### Freenode configuration

If you want to store your freenode configuration, edit and adapt this configuration file (`~/.weechat/irc.conf`):

```bash
[server_default]
[...]
nicks = "username1,username2,username3"
password = "pass"
realname = "realname"
username = "username"
[...]
```

### Autoconnect to channels

If you want to autoconnect to a specific channel (`~/.weechat/irc.conf`):

```bash
[server_default]
[...]
autoconnect = on
autojoin = "#channel1,#channel2"
[...]
```

And if channel 1 requires a password for instance (`~/.weechat/irc.conf`):

```bash
[server_default]
[...]
autoconnect = on
autojoin = "#channel1,#channel2 password"
[...]
```

## References

- [Weechat Website](https://www.weechat.org/)
