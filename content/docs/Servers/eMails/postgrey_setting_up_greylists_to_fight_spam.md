---
weight: 999
url: "/Postgrey_\\:_Mise_en_place_de_greylists_pour_lutter_contre_le_spam/"
title: "Postgrey: Setting Up Greylists to Fight Spam"
description: "This guide explains how to set up and configure Postgrey, a greylisting implementation for Postfix to effectively combat spam emails."
categories: ["Linux", "Security", "Debian"]
date: "2006-11-14T15:56:00+02:00"
lastmod: "2006-11-14T15:56:00+02:00"
tags: ["Postgrey", "Postfix", "Email", "Spam Protection", "Server Configuration", "Network"]
toc: true
---

## Introduction

The postgrey package is a greylisting implementation for postfix. It is very easy to set up and can stop 99.99% of spam in conjunction with blacklisting.

First back up everything - just in case! ;)

You should let any users know that greylisting may delay some emails, or otherwise pacify the end user community (Tip: offer them free software if they need compensating for delayed email).

## Installation

As root or your preferred way of getting root privileges:

```bash
apt-get install postgrey
```

You should have a policy server running on port 6000 of localhost, check with netstat:

```bash
netstat -anp | grep 60000
```

```
tcp        0      0 127.0.0.1:60000         0.0.0.0:*               LISTEN     18478/postgrey.pid
```

## Configuration

Now you need to add a line to `/etc/postfix/main.cf` to tell Postfix to use the new postgrey policy daemon. Since this is a restriction on SMTP traffic you need to add the following to "smtpd_recipient_restrictions":

```
check_policy_service inet:127.0.0.1:60000
```

Now reload postfix, postfix reload (me being a lazy luddite often runs **/etc/init.d/postfix restart** which is slower, but saves remembering which changes require what sort of restart, and what the syntax for this particular service is).

The only gotcha is that if the default Greylist text isn't right for you, then you probably want to read about this Debian bug 298832 to save some frustrating attempts to edit the text issued. This mostly affects people with many domains, where "postmaster@" isn't the preferred address for support, or to preserve corporate identity.

Postgrey is reliable - very - but you want to know if it is not running as Postfix will defer email if it can't speak to a policy daemon. Set up some sort of test!

Greylisting isn't for everyone. Postgrey by default has some exceptions for email servers that don't behave, and accounts to exempt in /etc/postgrey, and you can of course use these to whitelist servers. After 5 successful deliveries (by default) postgrey will whitelist a server, so you quickly get a whitelist of trusted servers, some people suggest reducing this default from 5. There is a script in /usr/share/doc/postgrey/postgrey_clients_dump which lists the whitelisted servers, it doesn't have execute permissions so you have to:

```bash
perl /usr/share/doc/postgrey/postgrey_clients_dump
```

The documentation directory also has contact details for the mailing list - it is very quiet - I think it just works. Beware greylisting will delay things like password emails for signing up to things like this website, but for many this is a price worth paying for the relief from spam.

## Configuration of Postgrey

* Hand in hand with greylisting is blacklisting. Part of the idea with greylisting is to delay the spammer till others have identified the poor exploited spambots. I looked around and picked the SBL-XBL blacklist. Again a restriction on SMTP so to the "smtpd_client_restrictions" add:

```
reject_rbl_client zen.spamhaus.org
```

You'll have to read the documentation to figure out what order to put your "smtpd_recipient_restrictions", as order can matter terribly - you've been warned.

That single one line change is estimated to stop 60% of spam!

* The other restriction I use is implemented by adding this line to /etc/postfix/main.cf

```
mime_header_checks = regexp:/etc/postfix/mime_header_checks.regexp
```

Where the file `/etc/postfix/mime_header_checks.regexp` has:

```bash
/filename=\&quot;?(.*)\.(bat|chm|cmd|com|cpl|do|exe|hta|jse|rm|scr|pif|vbe|vbs|vxd|xl|zip)\&quot;?$/
   REJECT For security reasons attachments of this type are rejected.
/^\s*Content-(Disposition|Type).*name\s*=\s*&quot;?(.+\.(lnk|cpl|asd|hlp|ocx|reg|bat|c[ho]m|cmd|exe|dll|vxd|pif|scr|hta|jse?|sh[mbs]|vb[esx]|ws[fh]|wav|mov|wmf|xl))&quot;?\s*$/
      REJECT Attachment type not allowed. File &quot;$2&quot; has the unacceptable extension &quot;$3&quot;
```

Note that is only two lines if you cut and paste it.

This is just taken and hacked from earlier examples on the Postfix mailing list, there may be better ways. This hasn't caught much email recently till I added ".zip" to the list, you may well not want the ".zip" (or other types), but I was getting many many viruses with .zip for every genuine attempt.

Currently the Greylisting and RBL stop virtually all the virus email traffic, so the MIME type filtering is almost unneeded, but I figure the few extra viruses caught will save some antivirus program (or end user) a little bit of effort. If the regular expression change worries you, leave it out.

* It may have taken a fair bit of explaining, but with the addition of one package, a three line change to main.cf, and two regular expressions, the advice in this article should stop the vast majority of spam. At least till the spammers start emulating SMTP servers. I still, alas, get the Nigerian scams - seems they just use disposable email accounts on real email servers. Greylisting also doesn't handle backscatter if people spam in your name, but it can spread the pain on your email server.
