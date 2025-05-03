---
weight: 999
url: "/Colorisations_dans_les_consoles/"
title: "Console Colorization"
description: "How to add color to your Linux console for better readability of man pages and log files"
categories: ["Linux", "Terminal"]
date: "2006-10-03T14:44:00+02:00"
lastmod: "2006-10-03T14:44:00+02:00"
tags: ["console", "terminal", "most", "ccze", "logs"]
toc: true
---

## Introduction

Here we'll see how to colorize certain elements in console terminals.

## Man Pages

To make the output of the man command more readable by adding color, you can do the following which involves using the most command instead of less for man pages:

```bash
sudo apt-get install most
export PAGER=`which most`
```

(use the Alt Gr+7 characters to execute the command)

To permanently set the PAGER value:

```bash
sudo vi /etc/security/pam_env.conf
```

Then modify to have these lines:

```bash
#PAGER DEFAULT=less
PAGER DEFAULT=most
```

## Log Files

I installed CCZE which allows me to better view Postfix logs:

```bash
apt-get install ccze
```

I'll let you read the man page and visit the website for all the features... but I already like the very practical:

```bash
tail -f /var/log/mail.log | ccze
```

Alternatively, there's also the possibility to save the colorized output as an HTML page... Just need to set up Apache :)

```bash
ccze -h < /var/log/mail.log > output.htm
```

Note that it doesn't just colorize Postfix, but also fetchmail, exim, apache, procmail, proftpd, squid, vsftpd, syslog, and others.
