---
weight: 999
url: "/Conversion_de_Timestamp_en_Date/"
title: "Converting Unix Timestamp to Date"
description: "How to convert Unix timestamps to human-readable dates using various command-line methods."
categories: ["Linux"]
date: "2009-12-11T21:38:00+02:00"
lastmod: "2009-12-11T21:38:00+02:00"
tags: ["Linux", "Command Line", "Timestamp", "Date Conversion"]
toc: true
---

How to convert a Unix timestamp to a readable date using command line tools.

## Using Perl

```bash
TIMESTAMP=1173279767
perl -e "print scalar(localtime($TIMESTAMP))"
```

## Using Perl with ctime module

```bash
TIMESTAMP=1173279767
perl -e "require 'ctime.pl'; print &ctime($TIMESTAMP);"
```

## Using Shell with awk

```bash
echo 1173279767 | awk '{print strftime("%c",$1)}'
```

## Using Shell with date

```bash
TIMESTAMP=1173279767
date -d "1970-01-01 $TIMESTAMP sec GMT"
```

Or even simpler:

```bash
date -d@1234567890
```
