---
weight: 999
url: "/Screen_\\:_Les_commandes_les_plus_utilisées/"
title: "Screen: Most Used Commands"
description: "Guide to using Screen for maintaining persistent terminal sessions, with most commonly used commands and multi-user configuration."
categories: ["Linux", "Security"]
date: "2011-05-11T20:00:00+02:00"
lastmod: "2011-05-11T20:00:00+02:00"
tags: ["Screen", "Terminal", "Linux Commands", "Multi-user", "Configuration"]
toc: true
---

## Introduction

Screen is really great, but if you don't use it every day, you can quickly lose track of the commands. Screen is used to keep a session open, possibly with an application running in it. When you exit the window and come back later, your application is still running and you can completely take control of it again.

I want to warn you, as I have seen many people do this, but screen is not equivalent to nohup!

## Installation

As usual:

```bash
apt-get install screen
```

## Usage

Here are the main commands I use (c-a = Ctrl+a):

```
c-a c   new window
c-a k   close a window
c-a p   previous window
c-a n   next window
c-a d   detach (leave screen running in the background)
```

```
c-a "   display all available windows
c-a [0-9] go to window 0-9
```

```
c-a S Split the screen
c-a <tab> move to the next window
c-a X Close the split window
c-a q Close all split windows
c-a M Monitor a window
```

## Muti-users

### Methode 1

On the first machine, start a screen with user toto:

```bash
screen
```

Then from a second machine, connect via ssh directly with user toto and do:

```bash
screen -x
```

Now both people can interact directly together.
There are also ACLs for screen.

### Methode 2

Using screen in multiuser mode requires screen to be setuid root. If you know about the potential security implications you can enable it by issuing:

```bash
chmod u+s `which screen`
chmod 755 /var/run/screen
```

We need to configure screen to use multiuser mode and change privileges for the guest. Put the following commands into ~/.screenrc. You can also use them in a screen session after pressing CTRL-a:

```bash
multiuser on
aclchg snoopy -x "?"    #Revoke permission to execute any screen command
aclchg snoopy +x "wall" #Allow writing simple messages in the terminal status line
aclumask snoopy-wx      #Default permissions to windows
acladd snoopy           #Enable user snoopy to access screen session
```

Start screen:

```bash
user@localhost $ screen
```

```bash
user@localhost $ screen -ls
There is a screen on:
       11521.pts-4.hostname      (Multi, attached)
1 Socket in /var/run/screen/S-user.
```

Now the guest can attach to the screen:

```bash
screen -r user/11521
```

### ACL

To allow toto to view the session without being able to act on it:

```bash
aclchg toto -w "#?"
```

By default, permissions are "rwx". Here are other very understandable examples:

```bash
aclchg toto -wx "#?"
aclchg toto +x "detatch,wall,colon"
```

The last command only authorizes certain options to be executed. Wall allows you to send messages to all connected screens.

## Configuration

### .screenrc

To see my screen configuration, I invite you to visit my git: [https://git.deimos.fr](https://git.deimos.fr)

For the configuration possibilities of the screenrc file, here's a small reminder:

```
Colors in screenrc
------------------
 0 Black             .    leave color unchanged
 1 Red               b    blue
 2 Green             c    cyan
 3 Brown / yellow    d    default color
 4 Blue              g    green           b    bold
 5 Purple            k    blacK           B    blinking
 6 Cyan              m    magenta         d    dim
 7 White             r    red             r    reverse
 8 unused/illegal    w    white           s    standout
 9 transparent       y    yellow          u    underline
note: "dim" is not mentioned in the manual.

STRING ESCAPES
--------------
 %%      percent sign (the escape character itself)
 %a      either 'am' or 'pm' - according to the current time
 %A      either 'AM' or 'PM' - according to the current time
 %c      current time HH:MM in 24h format
 %C      current time HH:MM in 12h format
 %d      day number - number of current day
 %D      Day's name - the weekday name of the current day
 %f      flags of the window
 %F      sets %? to true if the window has the focus
 %h      hardstatus of the window
 %H      hostname of the system
 %l      current load of the system
 %m      month number
 %M      month name
 %n      window number
 %s      seconds
 %t      window title
 %u      all other users on this window
 %w      all window numbers and names.
 %-w     all window numbers up to the current window
 %+w     all window numbers after the current window
 %W      all window numbers and names except the current one
 %y      last two digits of the year number
 %Y      full year number
```

## References

[File:Au-gnu screen-pdf.pdf](/pdf/au-gnu_screen-pdf.pdf)
