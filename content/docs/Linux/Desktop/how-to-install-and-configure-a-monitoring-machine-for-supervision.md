---
weight: 999
url: "/How_to_install_and_configure_a_monitoring_machine_for_supervision/"
title: "How to install and configure a monitoring machine for supervision"
description: "A guide on quickly setting up a monitoring machine with Nagios/Shinken for supervision purposes with automatic login and minimal interaction requirements."
categories: ["Monitoring", "Linux", "Debian"]
date: "2013-08-08T07:43:00+02:00"
lastmod: "2013-08-08T07:43:00+02:00"
tags:
  ["Nagios", "Shinken", "Debian", "Monitoring", "Slim", "Awesome", "Firefox"]
toc: true
---

![Nagios](/images/nagios_logo.avif)

{{< table "table-hover table-striped" >}}
|||
|-|-|
| **Operating System** | Debian 7 |
| **Last Update** | 08/08/2013 |
{{< /table >}}

## Introduction

The goal of this documentation is how to quickly setup a monitoring machine on Nagios/Shinken:

![Monito screen](/images/monito_screen.avif)

The idea is to have a machine that nobody never touch once installed. It's always boring to plug a keyboard and a mouse when somethings goes wrong on this kind of machine. That's why we'll see how to setup a minimal installation with all requirements to satisfy those needs.

## Installation

Install a Debian wheezy with a SSH server and additional utilities but without any X server or Window Manager and create an 'monitoring' user. Once done, connect with SSH on the machine install install those packages:

```bash
aptitude install xserver-xorg-core xorg awesome slim iceweasel wmctrl
```

We're installing iceweasel to get all dependencies of Firefox satisfied, but as we don't want any Iceweasel update break our configuration, we're going to install firefox:

```bash
wget ftp://ftp.mozilla.org/pub/firefox/releases/22.0/linux-x86_64/fr/firefox-22.0.tar.bz2
tar -xjf firefox-22.0.tar.bz2
mv firefox /usr/share/
rm firefox-22.0.tar.bz2
ln -s /usr/share/firefox/firefox /usr/bin/
chown -Rf monitoring /usr/share/firefox/firefox
```

Now we're ready for the configuration.

## Configuration

### slim

The configuration of slim is simple: we want to boot awesome without any credentials questions:

```bash {linenos=table,hl_lines=[4,7,11,15]}
# NOTE: if your system does not have bash you need
# to adjust the command according to your preferred shell,
# i.e. for freebsd use:
login_cmd           exec /bin/sh - ~/.xinitrc %session
##login_cmd           exec /bin/bash -login /etc/X11/Xsession %session
[...]
sessions            awesome
[...]
# default user, leave blank or remove this line
# for avoid pre-loading the username.
default_user        monitoring
[...]
# Automatically login the default user (without entering
# the password. Set to "yes" to enable this feature
auto_login          yes
```

### Awesome

Awesome should be configured to launch firefox in fullscreen mode automatically at startup. That's why we need to create a configuration file and then modify it. Login with monitoring user and copy the default config:

```bash
mkdir -p .config/awesome
cp /etc/xdg/awesome/rc.lua .config/awesome/
```

Then edit the configuration file and modify/add those lines:

```lua {linenos=table,hl_lines=[4,"18-29"]}
[...]
layouts =
{
    awful.layout.suit.max.fullscreen,
    awful.layout.suit.floating,
    awful.layout.suit.tile,
    awful.layout.suit.tile.left,
    awful.layout.suit.tile.bottom,
    awful.layout.suit.tile.top,
    awful.layout.suit.fair,
    awful.layout.suit.fair.horizontal,
    awful.layout.suit.spiral,
    awful.layout.suit.spiral.dwindle,
    awful.layout.suit.max,
    awful.layout.suit.magnifier
}
[...]
-- Move mouse to a corner
local safeCoords = {x=0, y=60}
local moveMouseOnStartup = true
local function moveMouse(x_co, y_co)
    mouse.coords({ x=x_co, y=y_co })
end
if moveMouseOnStartup then
        moveMouse(safeCoords.x, safeCoords.y)
end

-- Launch Firefox at startup
awful.util.spawn_with_shell("firefox")
```

As you can see, we're moving the mouse at boot instead of hiding it with slim. It's preferable to have the possibility to use it if we really need it than changing the configuration and reboot instead.

### Firefox

Launch firefox and install a fullscreen extension to autohide everything. You also need to change the "about:config" to set to false this parameter:

```
browser.sessionstore.resume_from_crash
```

To finish, disable automatic update and configure the URL with [Nagios]({{< ref "docs/Servers/Monitoring/Nagios/nagios_installation_and_configuration.md">}}):

```
https://nagios.deimos.fr/cgi-bin/status.cgi?host=all&servicestatustypes=28&hoststatustypes=3&serviceprops=42&sorttype=1&sortoption=6&noheader
```

Or if you prefer [Thruk]({{< ref "docs/Servers/Monitoring/Nagios/thruk_an_advanced_interface_for_nagios_and_mklivestatus.md">}}):

```
https://nagios.deimos.fr/thruk/cgi-bin/status.cgi?serviceprops=42&servicestatustypes=28&sortoption=6&type=detail&sorttype=1&host=all&hostprops=10&minimal=1
```

### Screensaver

Another interesting thing is to disable the screensaver :-D. And add in xinitrc:

```bash
#!/bin/sh

# Disable screen shutdown
xset -dpms
xset s off

# Launch Window Manager
exec awesome >> ~/.cache/awesome/stdout 2>> ~/.cache/awesome/stderr
```

And create the folder that doesn't exist:

```bash
mkdir -p ~/.cache/awesome
```

### Powersaving

As I'm a green IT, I don't want the monitoring machine to be up & running 24/24h 7/7j. That's why in the BIOS, I've set it up to start automatically at the wished hour. And I added this in the root crontab to autoswitch off correctly Firefox (using wmctrl is necessary are there is a bug on Firefox with SIGINT for more than 10 years):

```bash
PATH=/usr/bin:/sbin:/bin
0 20 * * * export DISPLAY=:0 ; su - monitoring -c "wmctrl -c firefox" ; halt
```
