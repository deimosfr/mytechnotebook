---
weight: 999
url: "/Lshell_\\:_limiter_les_possibilités_du_shell/"
title: "Lshell: Limiting Shell Capabilities"
description: "Learn how to restrict shell access and commands for users with Lshell on Linux systems. This guide covers installation, configuration, and integration with MySecureShell and sudo."
categories: ["Linux", "Security", "Servers"]
date: "2013-03-28T10:15:00+02:00"
lastmod: "2013-03-28T10:15:00+02:00"
tags: ["ssh", "security", "shell", "access control", "system administration"]
toc: true
---

## Introduction

[Lshell](https://lshell.ghantoos.org/) is a lightweight shell that allows restricting access to various commands and paths on your filesystems. It's ideal for controlling what users can do on your machine.

## Installation

```bash
aptitude install lshell
```

Currently on Debian, it's only available for the Squeeze version.

## Configuration

The configuration is quite simple:

```bash {linenos=table}
# lshell.py configuration file
#
# $Id: lshell.conf,v 1.20 2009/06/09 19:53:46 ghantoos Exp $

[global]
##  log directory (default /var/log/lshell/ )
logpath        : /var/log/lshell/
##  set log level to 0, 1, 2 or 3  (0: no logs, 1: least verbose)
loglevel       : 2
##  configure log file name (default is %u i.e. username.log)
#logfilename    : %y%m%d-%u


[default]
##  Allowed commands
allowed        : ['ls','echo','cd','ll','vi']

##  Forbidden commands
forbidden      : [';', '&', '|','`','>','<']

##  Limit of unauthorized command attempts before being kicked
warning_counter: 2

##  Aliases
aliases        : {'ll':'ls -l', 'vi':'vim'}

##  a value in seconds for the session timer
#timer          : 5

##  Authorized directories for users
path           : ['/home','/etc']

##  set the home folder of your user. If not specified the home_path is set to 
##  the $HOME environment variable
#home_path      : '/home/bla/'

##  update the environment variable $PATH of the user
#env_path       : ':/usr/local/bin:/usr/sbin'

##  allow or forbid the use of scp (set to 1 or 0)
#scp            : 1

##  allow of forbid the use of sftp (set to 1 or 0)
#sftp           : 1

##  list of command allowed to execute over ssh (e.g. rsync, rdiff-backup, etc.)
#overssh        : ['cd']

##  logging strictness. If set to 1, any unknown command is considered as 
##  forbidden, and user's warning counter is decreased. If set to 0, command is
##  considered as unknown, and user is only warned (i.e. *** unknown synthax)
#strict         : 1

##  force files sent through scp to a specific directory
#scpforce       : '/home'
```

As you can see, there are many useful options. However, today (at the time of writing), the options for the scp/sftp part are very limited. Fortunately, I have a solution.

You just need to assign this shell to the appropriate users (in `/etc/passwd`).

### Forcing lshell at login

If you have a PAM authentication via LDAP, it's possible to force a specific shell at login. This will override the information sent by NSS and replace it with the desired shell. Here we'll use lshell for all people connecting via LDAP:

```bash
nss_override_attribute_value loginShell /usr/bin/lshell
```

## Integration with MySecureShell

[MySecureShell](https://docs.services.mozilla.com/howtos/run-sync.html) will allow us to do all the configuration for the SFTP server, and in its configuration, we'll make sure to add lshell for the users we're interested in:

```bash
...
Shell       /usr/bin/lshell
...
```

And that's it - we now have a configuration where both SFTP and shell access are perfectly controlled.

## Integration with sudo

Integration with [sudo](https://docs.services.mozilla.com/howtos/run-sync.html) will be essential for people wanting to use commands that can only be executed as root. You'll need to modify the alias and allowed sections of lshell.

Don't forget to also modify the sudo configuration (otherwise watch out for security vulnerabilities) using [the documentation available here](https://docs.services.mozilla.com/howtos/run-sync.html).

## Resources
- http://lshell.ghantoos.org/
