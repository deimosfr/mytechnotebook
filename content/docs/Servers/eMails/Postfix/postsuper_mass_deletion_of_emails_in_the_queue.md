---
weight: 999
url: "/Postsuper_\\:_Suppression_massive_de_mails_dans_la_queue/"
title: "Postsuper: Mass Deletion of Emails in the Queue"
description: "Guide to removing emails from Postfix queue, both complete and selective deletions using postsuper utility."
categories: ["Linux"]
date: "2008-01-15T19:45:00+02:00"
lastmod: "2008-01-15T19:45:00+02:00"
tags: ["Postfix", "Mail", "System Administration", "Servers"]
toc: true
---

## Introduction

Postsuper is a Postfix utility that allows you to delete emails in the queue. I needed it to delete around 40,000 emails that were error messages from the MAILER-DAEMON user, with approximately 300 emails in the batch that needed to be delivered.

## Complete Deletion

To delete all emails in the queue:

```bash
postsuper -d ALL
```

## Partial Deletion

To delete all emails from the MAILER-DAEMON user, here's the script (replace MAILER-DAEMON with another name if needed):

```bash
#!/bin/sh
supp=`mailq | grep MAILER-DAEMON | awk -F"*" '{ print $1 }'`

for i in $supp ; do
        postsuper -d $i  
done
```

Or alternatively, here's another example:

```bash
mailq | tail +2 | awk 'BEGIN { RS = "" } / user@mydomain\.com$/ { print $1 }' | tr -d '*!' | postsuper -d -
```
