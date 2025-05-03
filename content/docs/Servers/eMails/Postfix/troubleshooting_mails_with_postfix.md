---
weight: 999
url: "/Réglage_de_problèmes/"
title: "Troubleshooting mails with Postfix"
description: "This documentation provides various methods and solutions for troubleshooting mail server problems, especially in Postfix environments."
categories: ["Servers", "Database", "Linux"]
date: "2010-02-05T11:09:00+02:00"
lastmod: "2010-02-05T11:09:00+02:00"
tags: ["Postfix", "Mail"]
toc: true
---

## Introduction

Postfix is great and handles high loads well, but when you need to debug mail problems, it can become less straightforward. Therefore, I'll document here all the tips and tricks I've found.

## Solutions

### Verify the configuration

In Postfix, it's not always obvious to find the problem, which is why they've been kind enough to include a command that will run checks:

```bash
postfix check
```

### Mail queue congestion

If you find yourself with a mail queue congestion problem (e.g., 30,000 emails in the queue), you can find out how many there are with the following command:

```bash
mailq
```

But if you have to delete them one by one, it's a bit tedious. That's why I created a small script to delete the entire queue:

```perl
#!/usr/bin/perl
# Mailq flusher
# Made by Pierre Mavro

use strict;

# Vars
my $queue_id;
my @queue_ids;
my @mails_addresses;
my $postsuper_command;
my $found_body=0;

# Starting
die "Sorry but you need to give an email address in argument :\n\teg. ./queue_flush_users.pl <viewall|test|delete> <user_mail>\n" if (! defined($ARGV[0]));

# Read the Postfix mailq
open MAIL_QUEUE, "mailq |";

# Put in arrays mails addresses and queue ids
while (<MAIL_QUEUE>) {
        my $mail_address;
        chomp $_;
        if ($_ =~ /\s+(.+)\@(.+)/i) {
                $mail_address="$1\@$2";
        }
        if ($_ =~ /^(\w+|\w+\*).+:\d+\s+(.+)\@(.+)/i) {
                $queue_id=$1;
                $mail_address="$2\@$3";
        }
        if (defined($mail_address)) {
                push @queue_ids, $queue_id;
                push @mails_addresses, $mail_address;
        }
}
close(MAIL_QUEUE);

my $total_mails=@mails_addresses;

for (0..$total_mails) {
        if ($ARGV[1] eq $mails_addresses[$_]) {
                if ($ARGV[0] eq "test") {
                        print "$queue_ids[$_] : $mails_addresses[$_]\n";
                } elsif ($ARGV[0] eq "delete") {
                        $postsuper_command=`postsuper -d $queue_ids[$_]`;
                }
                $found_body++;
        } elsif ($ARGV[0] eq "viewall") {
                print "$queue_ids[$_] : $mails_addresses[$_]\n";
                $found_body++;
        }
}

die "Sorry but nobody has been found.\n" if ($found_body == 0)
```

This script allows you to select a user and delete all emails where they are involved.

### pflogsumm

pflogsumm provides information about emails, including statistics like who receives a lot of emails, who sends a lot, etc. To get logs for the day:

```bash
pflogsumm -d today /var/log/mail.log
```

Get statistics from yesterday:

```bash
pflogsumm -d yesterday /var/log/mail.log.0
```

On the current log file:

```bash
pflogsumm /var/log/mail.log
```

Note:

```bash
pflogsumm /var/log/mail.log.0
```

will not give statistics just for yesterday.

### qshape

For bottleneck analysis, you can use the qshape command:

```bash
qshape -s hold
```

For more info:  
http://postfix.traduc.org/index.php/QSHAPE_README.html

### fatal: open database /etc/aliases.db: No such file or directory

You simply need to regenerate the alias database. The "newaliases" command will be enough. But on Solaris, and apparently older versions of Postfix, you'll need to stop and start the service as well:

```bash
svcadm disable svc:/network/smtp/postfix:default ; newaliases ; postalias ; svcadm enable svc:/network/smtp/postfix:default
```

### exec failed. errno=2.

On Solaris, when trying to send a message with the mail command and it doesn't work, giving this error:

```
exec failed. errno=2.
```

Just create a symbolic link:

```bash
ln -s /opt/csw/sbin/sendmail /usr/lib
```
