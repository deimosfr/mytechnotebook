---
weight: 999
url: "/Procmail_\\:_Filtrer_ses_mails_à_la_source/"
title: "Procmail: Filtering emails at the source"
description: "How to use Procmail to filter and sort emails directly at source level on Linux systems."
categories: ["Linux", "Servers", "Email"]
date: "2007-08-06T05:29:00+02:00"
lastmod: "2007-08-06T05:29:00+02:00"
tags: ["Procmail", "Email", "Linux", "Mail filtering", "Spam"]
toc: true
---

## Introduction

Procmail is a very powerful program used to filter emails. With it, you can redirect your mail, sort it, or even protect yourself against spam.

To give instructions to procmail, you need to create a file named `.procmailrc` in your home directory.

## Installation and configuration

To install procmail, as usual:

```bash
apt-get install procmail
```

Then, for documentation, I recommend the following:
[Procmail Documentation](/pdf/procmail1.pdf)

Follow that with my small example and you should be able to do what you want :-)

## Example

```bash
########
# Vars #
########

VERBOSE=ON
DROPPRIVS=YES
SHELL=/bin/sh
PATH=/usr/local/bin:/usr/bin:/bin
MAILDIR=$HOME/Maildir/
DEFAULT=$MAILDIR/new
LOGFILE=/var/log/procmail.log

# Personal Filters
SPAMBOX=$MAILDIR/.Trash/cur # Here I indicate the folder for read mails
CRONDIR=$MAILDIR/.Infos_Serveurs.Crontabs/cur
MLDKDIR=$MAILDIR/.Infos_Serveurs.Mldonkey/new # new corresponds to new mails
MSSBAK=$MAILDIR/.MySecureShell.Sauvegardes/cur
MYBAK=$MAILDIR/.Infos_Serveur.Backups/cur
UGC=$MAILDIR/.Sur_la_toile.UGC/new
EBAY=$MAILDIR/.Sur_la_toile.Ebay/new

# Newsletters
WEBPLANETE=$MAILDIR/.News.Webplanete/new
ZEROUNNET=$MAILDIR/.News.01net/new
CLUBIC=$MAILDIR/.News.Clubic/new
SILICON=$MAILDIR/.News.Silicon/new
PRESENCEPC=$MAILDIR/.News.PresencePC/new
FRSIRT=$MAILDIR/.News.FrSIRT/new
SECUOBS=$MAILDIR/.News.SecuObs/new

:0fw
* < 256000
       | /usr/bin/spamc -f
      :0e
       {
               EXITCODE=$?
       }

####################
# Personal Filters #
####################

# Spam to SPAMBOX
:0
* ^Subject:.*****SPAM***** # Subject starting with *****SPAM***** is sent to $SPAMBOX
$SPAMBOX

# Crontabs
:0
* ^Subject:.Cron
$CRONDIR

:0
* ^From:.root # Sender containing root is sent to the crontab folder
$CRONDIR

:0
* ^From:.arpwatch
$CRONDIR

:0
* ^From:.nagios@deimos.fr
$CRONDIR

# Mldonkey
:0
* ^From:.mldonkey
$MLDKDIR

# MSS Backup
:0
* ^Subject:.(MSSBackup*|MySQL*) # Here emails containing MSSBackup or MySQL are sent to $MSSBAK
$MSSBAK

# UGC
:0
* ^From:.*ugc.fr
$UGC

# Ebay + Paypal
:0
* ^From:(.*eBay.*|.*paypal.*)
$EBAY

#### Newsletters ####

:0
* ^Subject:.*WebPlanete.net*
$WEBPLANETE

:0
* ^From:.*01net
$ZEROUNNET

:0
* ^From:.*clubic
$CLUBIC

:0
* ^From:.*Silicon.fr
$SILICON

:0
* ^Subject:.*Presence PC
$PRESENCEPC

:0
* ^From:.*FrSIRT
$FRSIRT

:0
* ^From:."Secuobs.com"*
$SECUOBS
```
