---
weight: 999
url: "/Postfix\\:_hold_outgoing_mail_transport/"
title: "Postfix: hold outgoing mail transport"
description: "A guide on how to hold outgoing mail transport in Postfix without stopping the service, allowing time to analyze and troubleshoot mail infrastructure problems."
categories: ["Debian", "Linux", "Network", "Servers"]
date: "2014-03-04T17:35:00+02:00"
lastmod: "2014-03-04T17:35:00+02:00"
tags: ["Postfix", "Mail", "Queue Management", "Troubleshooting"]
toc: true
---

![Postfix](/images/postfix_logo.avif)

{{< table "table-hover table-striped" >}}
|||
|-|-|
| **Software version** | 2.10 |
| **Operating System** | Debian 7 |
| **Website** | [Postfix Website](https://www.postfix.org/) |
| **Last Update** | 04/03/2014 |
{{< /table >}}

## Introduction

When you manage Postfix and have a trouble with your mail infrastructure, you may want to set in maintenance your Postfix without stopping the service. Here is a way to hold the queue, giving the time to analyze the problem and then release the queue.

## Usage

You need to configure your Postfix as follow (`/etc/postfix/main.cf`):

```bash
defer_transports = hold
default_transport = hold
```

Then restart Postfix. Once you're ready and want to release the queue, remove those two previous lines, restart Postfix and force the queue to release:

```bash
service postfix restart
mailq -q
```

You can then look at the current status of deferred mails:

```bash
> qshape deferred
                                  T  5 10 20 40  80 160   320    640  1280 1280+
                       TOTAL 250218  0  0  0  0 151 138 38983 182693 26188  2065
                  hotmail.fr 116169  0  0  0  0   0   3 18273  88228  9037   628
                 hotmail.com  45086  0  0  0  0   0   5  3119  31012 10646   304
                     live.fr  25418  0  0  0  0   0   1  4187  20547   538   145
                   gmail.com  18519  0  0  0  0   0   4  4448  12874   851   342
                    yahoo.fr  10174  0  0  0  0 105   4  2227   6785   920   133
                     msn.com   5833  0  0  0  0   0   1   427   4138  1239    28
                 laposte.net   3971  0  0  0  0   0   1   956   2633   347    34
                     free.fr   2501  0  0  0  0   0   0   139   1980   360    22
                   orange.fr   2382  0  0  0  0   0   0   605   1578   166    33
                   yahoo.com   1929  0  0  0  0  31   5   188   1428   151   126
                    voila.fr   1662  0  0  0  0   0   0   120   1231   298    13
                      sfr.fr   1373  0  0  0  0   0   0   164   1128    66    15
                     neuf.fr   1295  0  0  0  0   0   0    54   1040   195     6
                  outlook.fr   1037  0  0  0  0   0   1   383    619    15    19
                  wanadoo.fr    985  0  0  0  0   0   0   151    542   286     6
                    live.com    577  0  0  0  0   0   0    81    460    21    15
            club-internet.fr    432  0  0  0  0   0   0     9    309   113     1
```
