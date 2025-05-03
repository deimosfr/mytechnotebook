---
weight: 999
url: "/Synchronisation_de_boites_mails_IMAP/"
title: "IMAP Mailbox Synchronization"
description: "How to synchronize IMAP mailboxes between two remote servers using imapsync"
categories: ["Linux", "Network", "Servers"]
date: "2006-10-23T11:00:00+02:00"
lastmod: "2006-10-23T11:00:00+02:00"
tags: ["IMAP", "Mail", "Synchronization", "Automation", "imapsync"]
toc: true
---

## IMAP Mailbox Synchronization

How do you synchronize a mailbox between two remote servers? A small tool called **Imapsync** exists for this purpose!

{{< alert context="warning" text="Warning: This software consumes a lot of CPU resources!" />}}

To install it, nothing complicated:

```bash
apt-get install imapsync
```

Next, let's create a folder in our _home_ directory called _.imapsync_.

```bash
mkdir ~/.imapsync
```

In this folder, we'll store the IMAP password for the account(s). The first file corresponds to the first server, and the second file to the second server:

```bash
echo "password1" > ~/.imapsync/secret1
echo "password2" > ~/.imapsync/secret2
```

You can also use the same file if you use the same password (as in my case).

For security reasons, let's change the permissions of this file(s):

```bash
chmod 600 ~/.imapsync/secret*
```

Then, we simply call the command with the appropriate arguments:

```bash
imapsync --syncinternaldates --host1 fire --ssl1 --user1 deimos --passfile1 ~/.imapsync/secret --host2 burnin --ssl2 --user2 deimos --passfile2 ~/.imapsync/secret
```

The parameters in italics should be adapted according to your needs. Here are some brief explanations of the options:

```
--syncinternaldates: fixes date issues with Eudora, Thunderbird...
--ssl1 and --ssl2: allows IMAP connections with SSL support. Remove these arguments if they don't apply to you.
```

Given the CPU load it generates, it's recommended to redefine the process priority. We can combine everything in a crontab entry to automate the process:

```bash
0 */1 * * * nice -n +19 imapsync --syncinternaldates --host1 fire --ssl1 --user1 deimos --passfile1 ~/.imapsync/secret --host2 burnin --ssl2 --user2 deimos --passfile2 ~/.imapsync/secret
```

And that's it, we're all set!
