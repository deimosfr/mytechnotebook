---
weight: 999
url: "/Réinjection_de_Spams_pour_tests/"
title: "Spam Reinjection for Testing"
description: "How to reinjecting spam messages for testing purposes, including methods to recover spam headers and use them for server testing."
categories: ["Linux"]
date: "2006-10-03T13:36:00+02:00"
lastmod: "2006-10-03T13:36:00+02:00"
tags: ["cd ~", "View source", "search", "What links here", "Servers", "Special pages", "Network", "Development", "Resume", "Solaris"]
toc: true
---

I need to reinject spam for testing on a server... but using the simple `mail` command is not sufficient to reinject messages with their headers...

We'll use the sendmail command (even for Postfix):

```bash
for i in message.*; do cat "$i" | sendmail -f from@domain.tld to@domain.tld ;done
```

The Postfix sendmail command implements the Postfix to Sendmail compatibility interface:

```bash
-f sender
Set the envelope sender address. This is the address where delivery problems are sent to, unless the message contains an Errors-To: message header.
```

To get "fresh" spam samples, you can use [SpamArchive.org](https://www.spamarchive.org/)

```bash
wget ftp://spamarchive.org/pub/archives/submit/679.r2.gz
```

And a small script to split everything:

```perl
cat convert
#!/usr/bin/perl -pl
if ( /^From / ) { close(OUT); open(OUT, ">>message.".$i++) || die "Can't open new file! $i\n"; select(OUT); print STDERR "Opened $i"; }
# ./convert 679.r2
```
