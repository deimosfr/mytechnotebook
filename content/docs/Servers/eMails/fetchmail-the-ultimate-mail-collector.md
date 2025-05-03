---
weight: 999
url: "/fetchmail-the-ultimate-mail-collector/"
title: "Fetchmail - The Ultimate Mail Collector"
description: "How to set up Fetchmail to retrieve emails from multiple accounts and consolidate them in one place"
categories: ["Linux", "Email", "System Administration"]
date: "2008-01-07T14:48:00+02:00"
lastmod: "2008-01-07T14:48:00+02:00"
tags: ["Fetchmail", "Email", "POP3", "IMAP"]
toc: true
---

## Introduction

Fetchmail is an IMAP, POP2, POP3, etc. collector.  
It allows you to retrieve emails from different mailboxes and consolidate them in your personal mailbox.

## Installation

For installation, the usual steps:

```bash
apt-get install fetchmail
```

## Configuration

### For a single user

Let's create a .fetchmailrc file in our home directory:

```bash
touch ~/.fetchmailrc
```

And insert these lines (for POP3):

```
poll pop.myprovider.com with proto POP3
user 'address@email.com' there with password 'PASSWORD' is 'USER' here options fetchall
```

### For multiple users

If you want to create a file that will fetch (yes, the verb 😉) emails for multiple users with a single file, create the file "/etc/fetchmailrc" and insert lines like this:

```
poll pop.myprovider.com proto pop3 port 995 user 'username' password 'password' smtpname "unixuser" options fetchall ssl
```

Here the POP3 is with SSL, which is why we have the "ssl" option at the end.

## Launch

Before doing anything, let's just test if the current configuration works:

```bash
fetchmail -c
```

Now that everything is good, let's retrieve our emails:

```bash
fetchmail

1 message for USER at pop.myprovider.com (1801 bytes).
reading message address@email.com:1 of 1 (1801 bytes) . deleted
```

## FAQ

### mail forwarding loop

If you encounter this type of message, and have a line like this:

```
Jan  7 15:40:25 fire postfix/local[24131]: 0AE56422D0: to=<xxx@mycompany.com>, orig_to=<xxx@mycompany.com>, relay=local, delay=0.1,delays=0.07/0.01/0/0.02, dsn=5.4.6, status=bounced (mail forwarding loop for xxx@mycompany.com)
```

You need to add the "dropdelivered" option configured like this:

```
poll x.x.x.x proto pop3 port 995 user 'pmavro' password 'xxx' dropdelivered smtpname "pierre.mavro" options fetchall ssl
```

## Resources
- [Fetchmail Documentation](/pdf/fetchmail.pdf)
