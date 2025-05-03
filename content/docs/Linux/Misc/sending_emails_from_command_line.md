---
weight: 999
url: "/Envoie_de_mails_en_lignes_de_commandes/"
title: "Sending Emails from Command Line"
description: "Guide to sending emails via command line using mail and mutt on Unix/Linux systems, with examples of attaching files and multiple recipients."
categories: ["Linux"]
date: "2009-09-21T13:26:00+02:00"
lastmod: "2009-09-21T13:26:00+02:00"
tags: ["Servers", "Linux", "CLI", "Email"]
toc: true
---

## Introduction

There can be great utility in sending emails from the command line. Everyone has their own use case, and I'll explain how it works.

## Mail

Mail is one of the most commonly used commands for sending emails:

```bash
echo "My text body" | mail -s "My Subject" "xxx@mycompany.com" "mail@mail2.com" "mail@mail3.com"
```

I think this is clear enough. You can do whatever you want before the "pipe" and then send everything.

If you want to send an attachment, do this:

```bash
echo "see attached file" | mail -a filename -s "subject" email@address
```

or

```bash
cat filename | uuencode filename | mail -s "Email subject" user@example.com
```

## Mutt

Mutt is not specifically dedicated to sending emails, but the advantage is being able to attach files. Here's another small example:

```bash
echo "My text body" | mutt -x -a my_attachment -s "My Subject" "xxx@mycompany.com" "mail@mail2.com" "mail@mail3.com"
```
