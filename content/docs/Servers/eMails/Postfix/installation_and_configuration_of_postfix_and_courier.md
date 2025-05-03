---
weight: 999
url: "/Installation_et_configuration_de_Postfix_et_Courrier/"
title: "Installation and Configuration of Postfix and Courier"
description: "A comprehensive guide on installing and configuring Postfix mail server with Courier on various operating systems including Debian, OpenBSD, and FreeBSD."
categories: ["Linux", "Mail", "Security"]
date: "2013-05-07T07:43:00+02:00"
lastmod: "2013-05-07T07:43:00+02:00"
tags:
  ["Postfix", "Courier", "OpenBSD", "FreeBSD", "Debian", "SMTP", "Mail Server"]
toc: true
---

![Postfix Architecture](/images/postfix.avif)

## Introduction

[Postfix](https://www.postfix.org/) is an email server and free software developed by Wietse Venema. It handles the delivery of electronic messages. It was designed as a faster, easier to administer, and more secure alternative to the historical Sendmail.

This software can handle almost all cases of professional use. Used with regexp in a junk file and a public anti-spam list, it prevents many spams without even having to scan message contents. It ideally replaces all kinds of less free solutions. You can find some how-tos on the official Postfix site. To optimize email analysis, Postfix allows delegating email management to an external process, which will determine whether the email is accepted or rejected (very useful in anti-spam systems).

The following diagram describes the internal architecture of postfix:

## Installation

### Debian

To install a Postfix server, here is the minimum to install:

```bash
apt-get install postfix courier-imap-ssl procmail spamc
```

### OpenBSD

On OpenBSD, we'll use the simple packaged version:

```bash
pkg_add -iv postfix
```

After installation, you need to replace postfix in place of sendmail. Just follow the instructions given at the end of the Postfix installation. In particular, you need to delete the sendmail cron tasks and configure the system to use postfix instead of Sendmail.

**This is what the postfix_enable script does!**

### FreeBSD

On FreeBSD, we'll use the packaged version:

```bash
pkg_add -vr postfix
```

After installation, you need to replace postfix in place of sendmail. Just follow the instructions given at the end of the Postfix installation. We'll continue here by adding this line in rc.conf:

```bash
...
postfix_enable="YES"
```

Let's disable sendmail by adding these lines:

```bash
sendmail_enable="NO"
sendmail_submit_enable="NO"
sendmail_outbound_enable="NO"
sendmail_msp_queue_enable="NO"
```

## Configuration

For the configuration, you will need to edit and adapt the configuration file:

```bash
# See /usr/share/postfix/main.cf.dist for a commented, more complete version

# Security
smtpd_banner = fire.deimos.fr - Microsoft Exchange (5.5)
biff = no
disable_vrfy_command = yes
smtpd_helo_required = yes

# Reject unknow domain
reject_unknown_recipient_domain = yes

# appending .domain is the MUA's job.
append_dot_mydomain = no

# Uncomment the next line to generate "delayed mail" warnings
#delay_warning_time = 4h

myhostname = fire.deimos.fr
alias_maps = hash:/etc/aliases
alias_database = hash:/etc/aliases.db
myorigin = /etc/mailname
mydestination = deimos.fr, fire, localhost
relayhost =
mynetworks = 127.0.0.0/8, 192.168.0.0/24, 10.8.0.0
home_mailbox = Maildir/
mailbox_size_limit = 0
recipient_delimiter = +
inet_interfaces = all
mailbox_command = procmail -a "$EXTENSION"

# Masquerade_domains hides hostnames from addresses
masquerade_domains = deimos.fr

# Virtual Domains
# virtual_alias_domains = mavrocordato.com mavro.fr deimos.servehttp.com
# virtual_alias_maps = hash:/etc/postfix/virtual

# Protection against Open Relay
smtpd_client_restrictions = reject_rbl_client bl.spamcop.net

# Protection against Spam
smtpd_recipient_restrictions =  permit_sasl_authenticated,
                                permit_mynetworks,
                                reject_unauth_destination,
                                reject_invalid_hostname,
                                reject_non_fqdn_sender,
                                reject_unknown_sender_domain,
                                reject_non_fqdn_recipient,
                                reject_unknown_recipient_domain,
                                reject_rhsbl_client blackhole.securitysage.com,
                                reject_rhsbl_sender blackhole.securitysage.com,
                                reject_rbl_client relays.ordb.org,
                                reject_rbl_client opm.blitzed.org,
                                reject_rbl_client list.dsbl.org,
                                reject_rbl_client cbl.abuseat.org,
                                reject_rbl_client dul.dnsbl.sorbs.net,
                                permit
smtpd_data_restrictions = reject_unauth_pipelining
mime_header_checks = regexp:/etc/postfix/mime_header_checks.regexp

# Use Amavis
content_filter = amavis:[127.0.0.1]:10024
receive_override_options = no_address_mappings
```

As you can see, the "smtpd_recipient_restrictions" line is quite long. This is because RBLs are integrated into it. Here is a short description:

RBLs aim to provide a list of servers known as major email senders and to list major spammers. It is actually a large generalized blacklist. The principle of use is very simple: when a filter receives an email, it checks if the sending server is contained in an RBL. If so, the email is categorized as spam. The RBLs that a filter uses as sources of servers are usually determined by the system administrator. This method therefore contains its share of controversy, as some RBLs are known to be more effective than others. The choice of RBLs therefore directly influences the effectiveness of the anti-spam system. In addition, some RBLs have looser rules than others regarding adding a server to their list, further complicating the situation. Among the known RBLs, note, among others, [SpamHaus](https://www.spamhaus.org/), [DynaBlock](https://www.njabl.org/dynablock.html), [Sorbs](https://www.sorbs.net/), and [DSBL](https://www.dsbl.org/). It is also possible to associate [ROKSO](https://www.spamhaus.org/rokso/index.lasso) with RBLs. ROKSO (Register of Known Spam Operations) is a list of the most active spammers. In fact, ROKSO members are responsible for nearly 80% of spam on the Net.

The "disable_dns_lookups = yes" option is used to disable DNS requests. When the "relayhost" is between "[ ]", it implies that postfix will not try to resolve the MX.

Then create a file `/etc/postfix/mime_header_checks.regexp`:

```bash
 /filename=\\"?(.*)\.(bat|chm|cmd|com|cpl|do|exe|hta|jse|rm|scr|pif|vbe|vbs|vxd|xl)\\"?$/
   REJECT For security reasons attachments of this type are rejected.
 /^\s*Content-(Disposition|Type).*name\s*=\s*"?(.+\.(lnk|cpl|asd|hlp|ocx|reg|bat|c[ho]m|cmd|exe|dll|vxd|pif|scr|hta|jse?|sh[mbs]|vb[esx]|ws[fh]|wav|mov|wmf|xl))
"?\s*$/
      REJECT Attachment type not allowed. File "$2" has the unacceptable extension "$3"
```

If certain attachments do not pass when sending or receiving, this is where you need to make changes (at the level of extensions).

Edit the `/etc/mailname` file and put your DNS:

```
deimos.fr
```

### OpenBSD

For OpenBSD, there are these additional lines:

```bash
mail_owner = _postfix
inet_protocols = all
unknown_local_recipient_reject_code = 550
debug_peer_level = 2
debugger_command =
         PATH=/bin:/usr/bin:/usr/local/bin:/usr/X11R6/bin
         xxgdb $daemon_directory/$process_name $process_id & sleep 5
sendmail_path = /usr/local/sbin/sendmail
newaliases_path = /usr/local/sbin/newaliases
mailq_path = /usr/local/sbin/mailq
setgid_group = _postdrop
html_directory = /usr/local/share/doc/postfix/html
manpage_directory = /usr/local/man
sample_directory = /etc/postfix
readme_directory = /usr/local/share/doc/postfix/readme
```

Edit the `/etc/mailname` file and put your DNS:

```
deimos.fr
```

## Launch Postfix

### Debian

Once all this is done, you just have to run the service restart command:

```bash
/etc/init.d/postfix restart
```

### OpenBSD

Postfix starts with the postfix start command, but first you need to enable Postfix instead of Sendmail:

```bash
postfix-enable
```

Check the `/var/log/messages` and `/var/log/maillog` files to see if everything went well.

Stopping and starting Postfix, checking the configuration:

- postfix check: basic configuration check
- postfix reload: reload configuration files
- postfix start: start postfix
- postfix stop: stop postfix

Queue management:

- mailq: display queue content
- postqueue -p: display queue content
- postqueue -f: force queue processing
- postfix flush: force queue processing

### FreeBSD

We will kill sendmail and start postfix:

```bash
pkill sendmail
postfix start
```

Postfix is now started.

## Creating mailboxes

To create a mailbox, simply go to your home directory and type:

```bash
maildirmake Maildir
```

Then, for emails to reach their destination, you must place these few lines in a ".procmailrc" file:

```bash
VERBOSE=ON
DROPPRIVS=YES
SHELL=/bin/sh
PATH=/usr/local/bin:/usr/bin:/bin
MAILDIR=$HOME/Maildir/
DEFAULT=$MAILDIR/new
LOGFILE=/var/log/procmail.log

:0fw
* < 256000
        | /usr/bin/spamc -f
        :0e
        {
                EXITCODE=$?
        }
```

## Protecting even more against spam

I recommend a site that lists all interesting addresses. Make good use of it, but be careful not to put too many!

http://spamlinks.net/filter-dnsbl-lists.htm

## Resources

- [Official Postfix Website](https://www.postfix.org)
- [Documentation on setting up Postfix on OpenBSD](https://www.formation.ssi.gouv.fr/stages/documentation/architecture_securisee/postfix.html)
- [Postscreen: fight a little more against spam](https://www.postfix.org/POSTSCREEN_README.html)
- [Quick installation of a secure mail server](https://feedproxy.google.com/~r/BlogInformatiqueByColundrum/~3/bhKbAPMbEPA/installation-rapide-serveur-mail-securise)
- [Advanced Postfix Configuration](/pdf/configuration_avancée_de_postfix.pdf)
- [Fighting spam with Postfix](/pdf/lutter_contre_le_spam_avec_postfix.pdf)
- [SMTP Server Mail Routing with Postfix](/pdf/serveur_smtp_routage_des_mails_avec_postfix.pdf)
- [Antispam Practices](/pdf/antispam_practices.pdf)
- [SMTP reject 554: Eliminate spam at the root](/pdf/273_unix_garden.pdf)
- [Documentation on configuring bounce mail on Postfix](/pdf/how_to_configure_custom_postfix_bounce_messages.pdf)
- [Adaptive and efficient anti-spam with OpenBSD and Spamd](/pdf/antispam_openbsd.pdf)
- [Virtual Users And Domains With Postfix Courier MySQL And SquirrelMail](/pdf/virtual_users_and_domains_with_postfix_courier_mysql_and_squirrelmail.pdf)
- [How To Whitelist Hosts IP Addresses In Postfix](/pdf/how_to_whitelist_hosts_ip_addresses_in_postfix.pdf)
- [Installing a mail server brick by brick](/pdf/Installation_d'un_serveur_mail_brique_par_brique.pdf)
