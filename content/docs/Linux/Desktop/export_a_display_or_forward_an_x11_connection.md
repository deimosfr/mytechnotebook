---
weight: 999
url: "/Exporter_un_display_ou_forwarder_une_connexion_X11/"
title: "Export a Display or Forward an X11 Connection"
description: "Guide on how to export displays or forward X11 connections to display GUI applications from one user or machine to another."
categories: ["Linux"]
date: "2008-08-31T12:00:00+02:00"
lastmod: "2008-08-31T12:00:00+02:00"
tags: ["Servers", "Linux", "X11", "SSH", "Display"]
toc: true
---

## Introduction

Such complicated names for something quite simple (though very practical).

- Exporting a display allows one user to transfer GUI applications to another user
- Forwarding an X11 connection allows transferring GUI applications from one user to another

In short, you might think they're the same thing! But not exactly.

## Exporting the Display

Let's say I'm connected as the user "deimos" and I launch a shell with the user "hostin". I want to access the GUI of hostin. First, I authorize my deimos user to accept connections from any user:

```bash
$ xhost +
access control disabled, clients can connect from any host
```

Now I want to know what my current display is:

```bash
$ echo $DISPLAY
:0.0
```

Then I connect as hostin and export my display to deimos:

```bash
$ su - hostin
$ export DISPLAY=:0.0
```

Now I just need to test with the hostin user:

```bash
$ xclock
```

Ta-da! The clock appears in my deimos session even though it's launched by hostin :-)

You can also do this with a remote machine by specifying the host (here are the lines to replace for the host deimos.fr):

```bash
deimos: $ xhost + deimos.fr
hostin:$ export DISPLAY=deimos.fr:0.0
```

## Forwarding X

We can forward our little X window through SSH. For this, on the SSH server (`/etc/ssh/sshd_config`):

```bash
X11Forwarding yes
```

This line must be set to yes.

Then, for the client to receive the window, you must connect with the -X argument:

```bash
$ ssh -X xxx@mycompany.com
$ xclock
```

Once again, ta-da! :-)

## References

[Create an X terminal](/pdf/creez_un_terminal_x.pdf)
