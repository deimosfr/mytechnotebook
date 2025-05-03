---
weight: 999
url: "/Synergy_\\:_Multi_screens_avec_plusieurs_ordinateurs/"
title: "Synergy: Multi-screen Setup with Multiple Computers"
description: "Learn how to set up Synergy to share keyboard and mouse across multiple computers with different operating systems."
categories: ["Linux", "Debian"]
date: "2012-11-07T09:02:00+02:00"
lastmod: "2012-11-07T09:02:00+02:00"
tags: ["Network", "Linux", "Mac"]
toc: true
---

![Synergy](/images/synergy_logo.avif)

{{< table "table-hover table-striped" >}}
|||
|-|-|
| **Software version** | 1.3 |
| **Operating System** | Debian 6<br />Mac OS 10.5+ |
| **Website** | [Synergy Website](https://synergy-foss.org) |
| **Last Update** | 07/11/2012 |
{{< /table >}}

## Introduction

Synergy allows you to easily share your mouse and keyboard between multiple computers. It is free and Open Source. Simply move your mouse from one computer to another by crossing their edges, just like moving between multiple monitors in a multi-screen setup. You can even share clipboards (copy-paste). All you need is a network connection. Synergy is cross-platform (works on Windows, Mac OS X and Linux).

![Synergy schema](/images/synergy_schema.avif)

## Installation

Choose your installation type based on your operating system

### Mac

Download the DMG, then copy the binaries to `/usr/local/bin/`. Go to the folder containing the binaries and run this command:

```bash
sudo cp ./synergyc ./synergys /usr/local/bin/synergyc
mkdir -p ~/Library/synergy
cp synergy.conf ~/Library/synergy
```

There you go, it's now installed :-)

### Windows

Download the installer and simply run it.

### Linux

Let's do it as usual:

```bash
aptitude install synergy
```

## Configuration

The configuration may seem complex, but it's not. Just remain logical. We won't cover how to configure on Windows, as everything is done with clicks and is really simple. We'll focus on the Mac/Unix part.

### The server

#### Mac

The server will determine what each machine does. For clients, there are no configuration files; everything is on the server. I'll take a typical configuration file and explain it. But first, there's something important to understand.

{{< alert context="warning" text="The names to insert in the configuration file correspond to the machine names OR their DNS names" />}}

To edit the configuration file:

```bash
# ~/Library/synergy/synergy.conf
section: screens 
    water:  
    earth:  
end

section: links
    water:  
        left = earth 
    earth:  
        right  = water 
end

section: aliases 
    water:  
        water.deimos.fr
    earth:  
        earth.deimos.fr
 end
```

Here, water corresponds to the server machine, and earth to the client machine. In the *screens* section, you need to declare all machines. I have only 2 here (server + client).

Next, we declare with *links* how the edges of the screens should interact. Here, earth is to the left of water and water is to the right of earth. Just read the configuration lines from right to left and make a sentence to understand how to configure the system.

The last section *aliases* is optional. It allows you to associate a name with a DNS name.

#### Linux

For configuration, please refer to the Mac section above and put this in `~/.synergy.conf`.

### The client

There is no client configuration :-)

## Launching & automation

Let's see how to test and automate these processes.

### Server

#### Mac

Here's the command line to start the server. Don't forget to change the server name. Here *server* should be replaced by *water*. Launch the tests:

```bash
synergys -f --config ~/Library/synergy/synergy.conf --name <server>
```

Now that the tests are complete, we can set up automatic startup so we don't have to type this every time. Let's start by creating what we need:

```bash
sudo mkdir -p /Library/StartupItems/Synergy
sudo chmod 755 /Library/StartupItems/Synergy/Synergy ; sudo touch /Library/StartupItems/Synergy/StartupParameters.plist
```

Then, we'll edit the following file:

```bash
# /Library/StartupItems/Synergy/Synergy
#!/bin/sh
 . /etc/rc.common

 run=(/usr/local/bin/synergys -f --config /Users/deimos/Library/synergy/synergy.conf --name server)

 KeepAlive () {
 proc=${1##*/}
 while [ -x "$1" ] ; do
     if ! ps axco command | grep -q "^${proc}\$" ; then
             "$@"
     fi
     sleep 3
 done
 }

 StartService () {
     ConsoleMessage "Starting Synergy"
     KeepAlive "${run[@]}" &
 }

 StopService () {
     return 0
 }

 RestartService () {
     return 0
 }

 RunService "$1"
```

Don't forget to replace *server* with your machine name and *deimos* with your username.

Then, edit the StartupParameters.plist file:

```bash
# StartupParameters.plist
 {
      Description = "Synergy Client";
      Provides = ("Synergy");
      Requires = ("Network");
      OrderPreference = "None";
 }
```

Restart your computer and it's done :-)

#### Linux

For Linux, once your configuration is done, all that's left is to launch the server:

```bash
synergys
```

### Client

#### Mac

To test connecting to the server, nothing could be simpler, open a terminal and type this command:

```bash
synergyc <server> &
```

Replace **server** with the name of the machine acting as the Synergy server.

Now that the tests are complete, we can set up automatic startup so we don't have to type this every time. Let's start by creating what we need:

```bash
sudo mkdir -p /Library/StartupItems/Synergy
sudo chmod 755 /Library/StartupItems/Synergy/Synergy ; sudo touch /Library/StartupItems/Synergy/StartupParameters.plist
```

Then, we'll edit the file:

```bash
# /Library/StartupItems/Synergy/Synergy
#!/bin/sh
 . /etc/rc.common

 run=(/usr/local/bin/synergyc -n $(hostname -s) -1 -f synergy-server)

 KeepAlive () {
 proc=${1##*/}
 while [ -x "$1" ] ; do
     if ! ps axco command | grep -q "^${proc}\$" ; then
             "$@"
     fi
     sleep 3
 done
 }

 StartService () {
     ConsoleMessage "Starting Synergy"
     KeepAlive "${run[@]}" &
 }

 StopService () {
     return 0
 }

 RestartService () {
     return 0
 }

 RunService "$1"
```

Don't forget to replace *server* with your machine name.

Then, edit the StartupParameters.plist file:

```bash
# StartupParameters.plist
 {
      Description = "Synergy Client";
      Provides = ("Synergy");
      Requires = ("Network");
      OrderPreference = "None";
 }
```

Restart your computer and it's done :-)

#### Linux

For the client, it's simple, just specify the server:

```bash
synergyc <server>
```

That's all :-)

## Graphical Interface

There is now a simple graphical interface to configure Synergy called QuickSynergy:

http://quicksynergy.sourceforge.net/

## FAQ

### My Synergy client keyboard is in English, how do I change it to another language?

You just need to add the "English(US)" language in your keyboard layout. Then restart the Synergy client and it's done :-)
