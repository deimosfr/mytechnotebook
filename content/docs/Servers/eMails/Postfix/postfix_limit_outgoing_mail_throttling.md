
---
weight: 999
url: "/Postfix\\:_limit_outgoing_mail_throttling/"
title: "Postfix: Limit Outgoing Mail Throttling"
description: "Learn how to limit outgoing mail throttling in Postfix to prevent blacklisting from MX servers when sending large volumes of email."
categories: ["Debian", "Linux", "Servers", "Network"]
date: "2015-08-06T00:10:00+02:00"
lastmod: "2015-08-06T00:10:00+02:00"
tags: ["Postfix", "Mail", "SMTP", "Throttling", "Mail Server", "Configuration"]
toc: true
---

![Postfix](/images/postfix_logo.avif)

{{< table "table-hover table-striped" >}}
|||
|-|-|
| **Software version** | 2.10 |
| **Operating System** | Debian 7 |
| **Website** | [Postfix Website](https://www.postfix.org/) |
| **Last Update** | 06/08/2015 |
{{< /table >}}

## Introduction

When you have a huge amount of mail to deliver, you can't release the queue at once and let the server maximize the outgoing mail throughput! The result will be: you'll get blacklisted from a lot of MX servers.

That's why you should take care of it and do traffic shaping

## Usage

You can add those lines in your Postfix configuration[^1]:

```bash
smtp_destination_concurrency_limit = 2
smtp_destination_rate_delay = 1s
smtp_extra_recipient_limit = 10
```

* default_destination_concurrency_limit: This means that postfix will up to two concurrent connections per receiving domains. The default value is 20.
* default_destination_rate_delay: Postfix will add a delay between each message to the same receiving domain. It overrides the previous rule and in this example, it will send one email after another with a delay of 1 second. If you want to disable this rule, either delete it or set to 0.
* default_extra_recipient_limit: Limit the number of recipients of each message. If a message had 20 recipients on the same domain, postfix will break it out to two different email messages instead of one.

Then restart your Postfix.

### Limit by domain

You can limit per domain if you want like this:

```bash
transport_maps = hash:/etc/postfix/transport

# Throttle limit policy mail (global)
smtp_destination_concurrency_limit = 4
smtp_extra_recipient_limit = 2

# Polite policy
polite_destination_concurrency_limit = 3
polite_destination_rate_delay = 0
polite_destination_recipient_limit = 5

# Turtle policy
turtle_destination_concurrency_limit = 2
turtle_destination_rate_delay = 1s
turtle_destination_recipient_limit = 2
```

Then add domains with the wished policy:

```bash
gmail.com polite:
yahoo.com polite:
hotmail.com turtle:
live.fr turtle:
orange.fr turtle:
```

Edit master configuration to inform postfix of those config. Add those lines:

```bash
polite unix - - n - - smtp
turtle unix - - n - - smtp
```

Postmap and reload:

```bash
postmap /etc/postfix/transport
service postfix reload
```

You're done :-)

## Resources

[^1]: http://steam.io/2013/04/01/postfix-rate-limiting/
