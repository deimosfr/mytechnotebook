---
weight: 999
url: "/Lancement_d'une_commande_à_la_réception_d'un_mail/"
title: "Running a Command When Receiving an Email"
description: "How to execute commands or scripts automatically when receiving emails on a Linux system."
categories:
  - Linux
date: "2013-05-06T08:32:00+02:00"
lastmod: "2013-05-06T08:32:00+02:00"
tags:
  - Linux
  - Email
  - Automation
  - Mail Server
toc: true
---

## Introduction

There's a very useful feature that can be easily implemented: launching a command or script when receiving an email.

## Usage

### Sending an Email

You can send an email in the following way:

```bash
echo "body" | mail -s "subject" test@test.com
```

If you want to add an attachment:

```bash
mailx bar@foo.com -s "HTML Hello" -a "Content-Type: text/html" < body.htm
```

Or for a binary attachment:

```bash
uuencode archive.tar.gz archive.tar.gz | mail -s "Emailing: archive.tar.gz" user@example.com
```

### Receiving an Email

To use this procedure, edit the aliases configuration and add a line like this (`/etc/aliases`):

```bash
test: "|touch /tmp/test"
```

When you send an email to your server with the recipient test (e.g. test@fqdn), the touch command will be executed.

Note: Don't forget to run newaliases after making changes.
