---
weight: 999
url: "/Installation_et_configuration_de_SpamAssassin/"
title: "SpamAssassin Installation and Configuration"
description: "A guide for installing and configuring SpamAssassin to filter spam emails, including adding Spam and Ham messages for training the system."
categories: 
  - Linux
date: "2010-03-15T10:17:00+02:00"
lastmod: "2010-03-15T10:17:00+02:00"
tags:
  - Spam
  - Email
  - Security
  - Filtering
toc: true
---

## Installation

To install SpamAssassin, it's very simple:

```bash
apt-get install spamassassin libmail-spf-query-perl
```

## Configuration

Here's my configuration file that you can adapt to your needs (`/etc/spamassassin/local.cf`):

```bash
# SpamAssassin Configuration
rewrite_header Subject  *****SPAM*****
use_bayes               1
bayes_auto_learn        1
required_score          5.0
skip_rbl_checks         0
report_safe             0

#pyzor
#use_pyzor               1
#pyzor_path /usr/bin/pyzor

#razor
#use_razor2              1
#razor_config /etc/razor/razor-agent.conf

ok_locales              en fr
whitelist_from *@deimos.fr noreply@lists.silicon.fr
blacklist_from *@mandrivaclub.com
```

Now, we need to enable SpamAssassin to start automatically. For this, in the file `/etc/default/spamassassin`, change from:

```bash
ENABLED=0
```

to:

```bash
ENABLED=1
```

Then, restart SpamAssassin:

```bash
/etc/init.d/spamassassin restart
```

There is also a website that allows you to [generate a SpamAssassin configuration](https://www.yrex.com/spam/spamconfig.php).

## Adding Spam and Ham

To add Ham or Spam, we'll insert this into the crontab of the person(s) who want to manage this:

```bash
sa-learn --spam --dir ~/Maildir/.Spam/cur && mv ~/Maildir/.Spam/cur/* ~/Maildir/.Trash/cur/
sa-learn --ham --dir ~/Maildir/.NoSpam/cur && mv ~/Maildir/.NoSpam/cur/* ~/Maildir/cur/
```

Alternatively, a small script can also do the job (`~/.antispam.sh`):

```bash
#!/bin/sh

sa-learn --spam --dir ~/Maildir/.Spam/cur
if [ `ls ~/Maildir/.Spam/cur/ | wc | awk '{ print $1 }'` != 0 ] ; then
	mv ~/Maildir/.Spam/cur/* ~/Maildir/.Trash/cur/ 2>&1 /dev/null
fi

sa-learn --ham --dir ~/Maildir/.NoSpam/cur
if [ `ls ~/Maildir/.NoSpam/cur/ | wc | awk '{ print $1 }'` != 0 ] ; then
	mv ~/Maildir/.NoSpam/cur/* ~/Maildir/cur/ 2>&1 /dev/null
fi
```

This creates two new folders in your mailbox (one for desired emails and one for undesired emails):

- If spam is found in a folder and it is not detected as spam, put it in the Spam folder.
- If an email arrives as spam when it is not, put it in the NoSpam folder to make it valid. This way, SpamAssassin will analyze the email so that next time, it arrives without being detected as spam.

Finally, if your Spam and NoSpam folders don't exist in your mailboxes:

```bash
mkdir -p ~/Maildir/.Spam/cur/ ~/Maildir/.NoSpam/cur/
```

But I recommend creating these folders with your regular email client.

## FAQ

### How can I test if the SPF module is working properly?

Put a simple valid email in sample-nonspam.txt and run this command:

```bash
spamassassin -D < sample-nonspam.txt
```

You should see something like this:

```
....
debug: registering glue method for check_for_spf_helo_pass
(Mail::SpamAssassin::Plugin::SPF=HASH(0x8d21990))
....
```
