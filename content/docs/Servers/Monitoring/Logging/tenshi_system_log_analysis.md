---
weight: 999
url: "/Tenshi_\\:_Analyse_des_logs_système/"
title: "Tenshi: System Log Analysis"
description: "How to install and configure Tenshi for analyzing system logs and receiving automated email alerts for specific events"
categories: ["Debian", "Storage", "Database"]
date: "2008-05-09T17:36:00+02:00"
lastmod: "2008-05-09T17:36:00+02:00"
tags:
  ["LogWatch", "Monitoring", "System Administration", "Security", "References"]
toc: true
---

## Introduction and Installation

The analysis of log files is a good compromise to the regular absence of consideration in companies for these files that contain crucial information for protecting systems against intrusions performed by hackers within the networks themselves.

_Note: this document was written in 2006; certain versions and configurations of the software used may be different from those mentioned; please refer to the official project sites in case of problems._

Dedicating a person to this analysis is not economically feasible in many companies, so implementing software that automates this process as much as possible is a lesser evil, provided that the number of reported alerts remains reasonable while maintaining the quality of the alerts that will be escalated; this requires fine-tuning and is therefore often complex to implement.

[LogWatch](https://www.logwatch.org/), [Logsentry](https://sourceforge.net/projects/sentrytools), [Logcheck](https://logcheck.org/), and [Swatch](https://swatch.sourceforge.net/) are the most well-known of these solutions in the open source world; through this document we will introduce [Tenshi](https://dev.inversepath.com/trac/tenshi), a log analyzer formerly known as Wasabi.

Tenshi offers multiple functionalities: email reports, high processing capacity, use of regular expressions (perl regexp), exception handling and filtering.

These features are combined with a very simple installation and configuration flexibility that will delight those most resistant to this type of application: a single file based on a concept of objects or groups of objects, intuitive syntax and internal "crontab".

Also worth noting is the minimal resources needed for this analysis, enabling the processing of multiple files with a large number of entries.

### Debian

With Debian, Tenshi installs as follows:

```bash
apt-get install tenshi
```

### Sources

For installation from the sources:

```bash
root@ns30074:~ # wget http://www.gentoo.org/~lcars/tenshi/tenshi-latest.tar.gz
--01:03:14-- http://www.gentoo.org/~lcars/tenshi/tenshi-latest.tar.gz
=> `tenshi-latest.tar.gz'
Résolution de www.gentoo.org... 38.99.64.201, 66.219.59.46, 66.241.137.77
Connexion vers www.gentoo.org[38.99.64.201]:80...connect
requête HTTP transmise, en attente de la réponse...302 Found
Emplacement: http://www.gentoo.org/~lcars/tenshi/tenshi-latest.tar.gz
--01:03:16-- http://www.gentoo.org/~lcars/tenshi/tenshi-latest.tar.gz
=> `tenshi-latest.tar.gz'
Résolution de dev.gentoo.org... 140.211.166.183
Connexion vers dev.gentoo.org[140.211.166.183]:80...connect
requête HTTP transmise, en attente de la réponse...200 OK
Longueur: 19,220 [application/x-gzip]

100%[==========>] 19,220 71.91K/s

01:03:17 (71.66 KB/s) - tenshi-latest.tar.gz sauvegardé [19220/19220]
```

Tenshi is a program based on the PERL language; for sending reports via email, it requires the installation of the additional "Net::SMTP" module.

This module provides an API for the SMTP protocol; its installation can be done through another module that allows access to the CPAN archive (Comprehensive PERL Archive Network) which groups all homologated PERL modules together with their various documentation.

To start the CPAN connection:

```bash
root@ns30074:~/tenshi-0.4 # perl -e shell -MCPAN

cpan shell -- CPAN exploration and modules installation (v1.7601)
ReadLine support available (try 'install Bundle::CPAN')
```

For newcomers to using the CPAN module, you can retrieve and consult information related to the "Net::SMTP" module with the following command:

```bash
cpan> i Net::SMTP
CPAN: Storable loaded ok
Going to read /root/.cpan/Metadata
Database was generated on Tue, 20 Sep 2005 01:59:21 GMT
CPAN: LWP::UserAgent loaded ok
Fetching with LWP:
ftp://ftp.pasteur.fr/pub/computing/CPAN/authors/01mailrc.txt.gz
Going to read /root/.cpan/sources/authors/01mailrc.txt.gz
CPAN: Compress::Zlib loaded ok

Fetching with LWP:
ftp://ftp.pasteur.fr/pub/computing/CPAN/modules/03modlist.data.gz
Going to read /root/.cpan/sources/modules/03modlist.data.gz
Going to write /root/.cpan/Metadata
Strange distribution name [Net::SMTP]
Module id = Net::SMTP
DESCRIPTION Interface to Simple Mail Transfer Protocol
CPAN_USERID GBARR (Graham Barr )
CPAN_VERSION 2.29
CPAN_FILE G/GB/GBARR/libnet-1.19.tar.gz
DSLI_STATUS adpf (alpha,developer,perl,functions)
MANPAGE Net::SMTP - Simple Mail Transfer Protocol Client
INST_FILE /usr/share/perl/5.8/Net/SMTP.pm
INST_VERSION 2.26
```

The installation of the module is done as follows:

```bash
cpan> install Net::SMTP
Running install for module Net::SMTP
.........
/usr/bin/make install -- OK
```

We can see that the installation was performed correctly; we can now disconnect from CPAN:

```bash
cpan> exit
```

Let's look at the Tenshi archive that we previously retrieved; first decompress and unarchive it:

```bash
root@ns30074:~ # tar -zxvf tenshi-latest.tar.gz
tenshi-0.4/
tenshi-0.4/tenshi.debian-init
tenshi-0.4/Makefile
tenshi-0.4/LICENSE
tenshi-0.4/README
tenshi-0.4/tenshi.conf
tenshi-0.4/tenshi.8
tenshi-0.4/tenshi
tenshi-0.4/tenshi.gentoo-init
tenshi-0.4/INSTALL
tenshi-0.4/tenshi.ebuild
tenshi-0.4/Changelog
tenshi-0.4/COPYING
tenshi-0.4/tenshi.solaris-init
tenshi-0.4/CREDITS
```

Then access the directory that was recently created and contains the files needed to install the application:

```bash
root@ns30074:~ # cd tenshi-0.4/
```

List the directory to see its structure:

```bash
root@ns30074:~/tenshi-0.4 # ls

Changelog CREDITS LICENSE README tenshi.8 tenshi.debian-init tenshi.gentoo-init
COPYING INSTALL Makefile tenshi tenshi.conf tenshi.ebuild tenshi.solaris-init
```

The next two operations consist of creating a "tenshi" user and a group of the same name that will govern the execution of Tenshi on your system:

```bash
root@ns30074:~/tenshi-0.4 # useradd tenshi
root@ns30074:~/tenshi-0.4 # groupadd tenshi
```

Before proceeding with the installation itself, it's necessary, depending on the operating system and distribution, to adjust the configuration of the "Makefile" file; for me it was necessary to change the installation path of the manual pages by modifying line 9 of the file on a "debian-like" system:

```bash
root@ns30074:~ # updatedb | locate "man8/" | more
/usr/share/man/man8/update-passwd.8.gz

root@ns30074:~/tenshi-0.4 # vi Makefile
```

Replace "mandir = /usr/man" with "mandir = /usr/share/man/".

You can now install Tenshi on the system:

```bash
root@ns30074:~/tenshi-0.4 # make install
install -D tenshi /usr/sbin/tenshi
[ -f /etc/tenshi/tenshi.conf ] || \
install -g root -m 0644 -D tenshi.conf /etc/tenshi/tenshi.conf
install -d /usr/share/doc/tenshi-0.4
install -m 0644 README INSTALL CREDITS LICENSE COPYING Changelog /usr/share/doc/tenshi-0.4/
install -g root -m 0644 tenshi.8 /usr/share/man//man8/
install -g root -m 755 -d /var/lib/tenshi
```

## Configuration

Let's examine the configuration file part by part to understand how Tenshi works, which will also introduce us to the mechanisms that govern log analysis.

Once the installation is complete, we can move on to the next step, which is configuration: this is done through a single file called "tenshi.conf" located at `/etc/tenshi/tenshi.conf`

```bash
root@ns30074:~/tenshi-0.4 # vi /etc/tenshi/tenshi.conf
```

At the beginning of the file, we find the "generic" configurations, including the setting of the previously mentioned "tenshi" user and "tenshi" group created on the system that will own the processes related to the future launch of Tenshi:

```bash
##
## tenshi 0.4 sample conf
##

# general settings

set uid tenshi
set gid tenshi
```

The following configuration line requires some adjustments. By default, we find "set pidfile /var/run/tenshi.pid", but the files governing PIDs are located in the `/var/run/` directory, which by default belongs to the "root" user:

```bash
ls -l /var/ | grep run
drwxr-xr-x 17 root root 4096 2006-06-16 07:57 run
```

The following operations must be performed to solve this problem:

- Create a directory specific to the Tenshi PID

```bash
root@ns30074:~/tenshi-0.4 # mkdir /var/run/tenshi
```

- Set the necessary permissions for the future creation of PID files in this directory when Tenshi launches

```bash
root@ns30074:~/tenshi-0.4 # chown tenshi:root /var/run/tenshi/
```

- Verify that the previous operation was successful:

```bash
root@ns30074:~/tenshi-0.4 # ls -l /var/run/ | grep tenshi
drwxr-xr-x 2 tenshi root 4096 2006-06-16 01:49 tenshi
```

Then all that remains is to replace the line "set pidfile /var/run/tenshi.pid" with "set pidfile /var/run/tenshi/tenshi.pid" in the configuration file to avoid warnings related to the mentioned problem during future Tenshi launches.

The rest of the configuration allows you to define the log files that will be monitored and whose alerts will be reported to you; by default, you'll find the following two lines:

```bash
set logfile /var/log/messages
set logfile /var/log/mail.log
```

Ideally, and to adjust according to the context of your distribution regarding the logging management of various events by "syslogd" or others, you can add the following lines to make this list as exhaustive as possible and the analysis more effective:

```bash
set logfile /var/log/syslog
set logfile /var/log/sulog
set logfile /var/log/user.log
set logfile /var/log/auth.log
```

Next come some optimizations that you can leave by default, unlike the definition of the SMTP server that will be used to send reports if you do not want to use the default one running on "localhost":

```bash
set sleep 5
set limit 800
set pager_limit 2
set mask ___
set mailserver mail.deimos.fr
set subject Tenshi report
set hidepid on
```

Tenshi has the unique characteristic of not using the system's "crontab" for its operation; it's therefore possible to configure different automation elements directly in the "tenshi.conf" file according to several criteria.

Elements of type "queue" are to be set up via the following syntax: "set queue []"; the default subject is "tenshi report" if you don't specify specific ones.

Below you can see 6 elements configuring the sending from "webmaster@xxx@mycompany.com" to "xxx@mycompany.com" according to different frequencies and periods throughout the day based on their criticality.

For the "mail" element, scheduled sending at 6:30 PM of a report every day:

```bash
set queue mail webmaster@xxx@mycompany.com xxx@mycompany.com [30 18 * * *]
```

For the "nf" element, scheduled sending every thirty minutes of a report every day:

```bash
set queue nf webmaster@xxx@mycompany.com xxx@mycompany.com [*/30 * * * *]
```

For the "report" element, scheduled sendings every 2 hours in the interval between 9 AM and 5 PM of a report every day:

```bash
set queue report webmaster@xxx@mycompany.com xxx@mycompany.com [0 9-17/2 * * *]
```

For the "misc" element, scheduled sendings under the same conditions as the previous element:

```bash
set queue misc webmaster@xxx@mycompany.com xxx@mycompany.com [0 9-17/2 * * *]
```

For the "critical" element, scheduled sending immediately for each identified occurrence:

```bash
set queue critical webmaster@xxx@mycompany.com xxx@mycompany.com [now]
```

For the "root" element, scheduled sending immediately for each identified occurrence:

```bash
set queue root webmaster@xxx@mycompany.com xxx@mycompany.com [now]
```

Sending reports to multiple users is also possible to configure; simply replace "xxx@mycompany.com" with "xxx@mycompany.com, root@xxx@mycompany.com" for example.

Following in the file are the exceptions that will allow you not to report information that is not of interest; they are sent to a "trash" element that is not part of the report sending "crontab":

```bash
trash ^hub.c
trash ^usb.c
trash ^uhci.c
trash ^sda
trash ^Initializing USB
trash ^scsi0 : SCSI emulation
trash ^Vendor:
trash ^Type:
trash ^Attached scsi removable
trash ^SCSI device sda
trash ^sda: Write
trash ^/dev/scsi
trash ^WARNING: USB
trash ^USB Mass Storage
trash ^/dev
trash ^ISO
trash ^floppy0
trash ^end_request
trash ^Directory
trash ^I/O error: dev 08:(.+), sector
```

Similarly, we remove the indications for repeated entries:

```bash
repeat ^(?:last message repeated|above message repeats) (\\d+) time
```

An example group for the SSHD service; definition of the group:

```bash
group ^sshd(?:\(pam_unix\))?:
```

All occurrences of the string "sshd: fatal: Timeout before authentication for" regardless of the following characters indicating that an SSH session has been closed due to inactivity for too long; this will be recorded in the reports belonging to the "report" element:

```bash
report ^sshd: fatal: Timeout before authentication for (.+)
```

All occurrences of the string "critical ^sshd: Illegal user" indicating that an username not conforming to the system has been used for SSH authentication; this will be recorded in the reports belonging to the "critical" element:

```bash
critical ^sshd: Illegal user
```

All occurrences of the string "sshd: Connection from" regardless of the following characters indicating an established SSH connection; this will be recorded in the reports belonging to the "report" element:

```bash
report ^sshd: Connection from (.+)
```

All occurrences of the string "sshd: Connection closed" regardless of the following characters indicating a closed SSH session; this will be recorded in the reports belonging to the "report" element:

```bash
report ^sshd: Connection closed (.+)
```

All occurrences of the string "sshd: Closing connection" regardless of the following characters indicating the ongoing closure of an SSH connection; this will be reported in the reports belonging to the "report" element:

```bash
report ^sshd: Closing connection (.+)
```

All occurrences of the string "sshd: Found matching (.+) key: (.+)" regardless of the following and intermediate characters indicating that the key proposed for SSH authentication has been found in the list of system keys; this will be recorded in the reports belonging to the "report" element:

```bash
report ^sshd: Found matching (.+) key: (.+)
```

All occurrences of the string "sshd: Accepted publickey" regardless of the following characters indicating that the public key used for opening an SSH session has been accepted; this will be recorded in the reports belonging to the "report" element:

```bash
report ^sshd: Accepted publickey (.+)
```

All occurrences of the string "sshd: Accepted rsa for (.+) from (.+) port (.+)" regardless of the following and intermediate characters indicating that the proposed RSA key has been accepted; this will be recorded in the reports belonging to the "report" element:

```bash
report ^sshd: Accepted rsa for (.+) from (.+) port (.+)
```

All occurrences of the string "sshd\(pam_unix\): session opened for user root by root\(uid=0\)" indicating the opening of an SSH session with the root account from the root account; this will be recorded in the reports belonging to the "root" element:

```bash
root ^sshd\(pam_unix\): session opened for user root by root\(uid=0\)
```

All occurrences of the string "sshd\(pam_unix\): session opened for user root by \(uid=0\)" indicating the opening of an SSH session with the root account; this will be recorded in the reports belonging to the "root" element:

```bash
root ^sshd\(pam_unix\): session opened for user root by \(uid=0\)
```

And finally closing the sshd group:

```bash
group_end
```

On the same principles, several other groups are to be defined according to your system and its applications/services.

The "Sendmail" service with:

```bash
group ^sendmail:
mail ^sendmail: (.+): to=(.+),(.+)relay=(.+),(.+)stat=Sent(.+)
mail ^sendmail: (.+): to=(.+),(.+)relay=(.+),(.+)stat=Sent
mail ^sendmail: (.+): from=(.+),(.+)relay=(.+)
mail ^sendmail: STARTTLS=client(.+)
mail ^sendmail
group_end
```

The "Sm-mta" service:

```bash
group ^sm-mta:
mail ^sm-mta: (.+): to=(.+),(.+)delay=(.+)
mail ^sm-mta: (.+): to=(.+),(.+)relay=(.+),(.+)stat=Sent(.+)
mail ^sm-mta: (.+): to=(.+),(.+)relay=(.+),(.+)stat=Sent
mail ^sm-mta: (.+): to=(.+),(.+)relay=local(.+)stat=Sent(.+)
mail ^sm-mta: (.+): to=(.+),(.+)relay=local(.+)stat=Sent
mail ^sm-mta: (.+): to=(.+),(.+)stat=Sent(.+)
mail ^sm-mta: (.+): to=(.+),(.+)stat=Sent
mail ^sm-mta: (.+): from=(.+),(.+)relay=local(.+)
mail ^sm-mta: (.+): from=(.+),(.+)relay=(.+)
mail ^sm-mta: STARTTLS=server(.+)
mail ^sm-mta: STARTTLS=client(.+)
trash ^sm-mta:.+User unknown
mail ^sm-mta: ETRN
mail ^sm-mta
group_end
```

The "Ipop3d" service with:

```bash
group ^ipop3d:
mail ^ipop3d: Login user=(.+)
mail ^ipop3d: Logout user=(.+)
mail ^ipop3d: pop3s SSL service init from (.+)
mail ^ipop3d: pop3 service init from (.+)
mail ^ipop3d: Auth user=(.+)
mail ^ipop3d: Command stream end of file, while reading
mail ^ipop3d: Command stream end of file while reading
mail ^ipop3d: AUTHENTICATE LOGIN failure host=(.+)
mail ^ipop3d: AUTHENTICATE PLAIN failure host=(.+)
mail ^ipop3d: Login failed
mail,critical ^ipop3d:
group_end
```

The "Imapd" service with:

```bash
group ^imapd:
mail ^imapd: Login user=(.+)
mail ^imapd: Logout user=(.+)
mail ^imapd: port (.+) service init from (.+)
mail ^imapd: imaps SSL service init from (.+)
mail ^imapd: Command stream end of file, while reading
mail ^imapd: Command stream end of file while reading
mail ^imapd: Authenticated user=(.+)
mail ^imapd: AUTHENTICATE LOGIN failure host=(.+)
mail ^imapd: AUTHENTICATE PLAIN failure host=(.+)
mail ^imapd: Autologout(.+)
mail ^imapd: Login failed
mail,critical ^imapd:
group_end
```

The "login\(pam_unix\)" service with:

```bash
group ^login\(pam_unix\):
critical ^login\(pam_unix\): session opened for user root by root\(uid=0\)
critical ^login\(pam_unix\): session opened for user root by \(uid=0\)
report ^login\(pam_unix\): session closed for user (.+)
report ^login\(pam_unix\): session opened for user (.+)
group_end
```

The "su(pam_unix)" service with:

```bash
group ^su\(pam_unix\):
root,report ^su\(pam_unix\): session opened for user root
root,report ^su\(pam_unix\): session closed for user root(.+)
report ^su\(pam_unix\): session opened for user (.+)
report ^su\(pam_unix\): session closed for user (.+)
group_end
```

Then all records of the string "netfilter" for filtering will be sent to the "nf" element:

```bash
nf ^netfilter
```

All records of the string "sudo" for commands executed under an account other than the active SHELL account will be sent to the "critical" element:

```bash
critical ^(?:/usr/bin)?sudo:
```

All records of the string "init" for service launches will be sent to the "critical" element:

```bash
critical ^init
```

All records of the string "passwd\(pam_unix\):" for password changes will be sent to the "report" element:

```bash
report ^passwd\(pam_unix\):
```

All other records not matching the previous identification regular expressions will be sent to the "misc" element:

```bash
misc .*
```

## Usage

With configurations done, we move on to using Tenshi by putting ourselves in concrete alert cases to receive the corresponding alarms on a predefined email provided that the basic configurations have been properly performed beforehand.

Tenshi functions as a service, which notably implies that you can close the terminal where you launched it without stopping it (equivalent to adding the "&" character at the end of a SHELL command).

During a previous "ls", you may have noticed the presence of 3 initialization scripts for Linux operating systems (Debian and Gentoo distributions) and Solaris:

```bash
root@ns30074:~/tenshi-0.4 # ls | grep init

tenshi.debian-init
tenshi.gentoo-init
tenshi.solaris-init
```

We are currently using a "Debian-like" distribution (Ubuntu), so we need to copy the corresponding initialization script into the `/etc/init.d/` directory which contains all the service initialization scripts of the system:

```bash
root@ns30074:~/tenshi-0.4 # cp tenshi.debian-init /etc/init.d/tenshi
```

Test this script:

```bash
root@ns30074:~/tenshi-0.4 # sh /etc/init.d/tenshi
Usage: /etc/init.d/tenshi {start|stop|restart}
```

Add execution rights to the script:

```bash
root@ns30074:~/tenshi-0.4 # chmod +x /etc/init.d/tenshi
```

Start the Tenshi daemon:

```bash
root@ns30074:~/tenshi-0.4 # /etc/init.d/tenshi start
Starting log monitor: tenshi.
```

Test the effective launch of Tenshi:

```bash
root@ns30074:~ # ps auwx | grep tenshi
tenshi 19424 0.0 1.0 8260 5312 ? Ss 01:51 0:00 /usr/bin/perl /usr/sbin/tenshi
tenshi 19425 0.0 0.1 3176 676 ? S 01:51 0:00 /usr/bin/tail -q --follow=name --retry -n 0 /var/log/messages /var/log/mail.log
```

Stop Tenshi:

```bash
root@ns30074:~/tenshi-0.4 # /etc/init.d/tenshi stop
Stopping log monitor: tenshi.
```

Forced stop of Tenshi in case of problems with the scripts:

```bash
root@ns30074:~/tenshi-0.4 # killall -9 tail tenshi
```

Restart Tenshi:

```bash
root@ns30074:~/tenshi-0.4 # /etc/init.d/tenshi restart
Stopping log monitor: tenshi.
Starting log monitor: tenshi.
```

Tenshi can also be launched outside of initialization scripts to, for example, test a configuration file or launch the service with an alternative file. Debug mode is also possible among others, and will be very useful if you encounter a problem particularly with sending emails,

```bash
root@ns30074:~/tenshi-0.4 # tenshi -h
tenshi 0.4
Copyright 2004-2006 Andrea Barisani
and Rob Holland

Usage: /usr/sbin/tenshi [-c conf_file] [-C|-d|-f|-p] [-P pid_file]
-c configuration file
-C test configuration syntax
-d debug mode
-f foreground mode
-p profile mode
-P pid file
-h this help
```

The last point to address, and not the least, are the reports with the alerts that will be sent to you by email. The first is a report from the "report" element for a closure due to inactivity for too long on an active session of the SSH service:

```bash
From: webmaster@xxx@mycompany.com
To: xxx@mycompany.com
Date: Sun, 18 Jun 2006 03:08:33 +0200
X-tenshi-version: 0.4
X-tenshi-hostname: ns30074
Subject: tenshi report [report]


ns30074:
1: sshd: fatal: Timeout before authentication for ___
```

The second example is a report for the "misc" element with an SSH authentication error, as well as the execution of a command by the system "crontab" and the opening of a "pam_unix" session for the root account associated with this previous entry "_/1 _ \* \* \* root /usr/local/rtm/bin/rtm 41 >/dev/null 2>/dev/null" of the "crontab":

```bash
From: webmaster@xxx@mycompany.com
To: xxx@mycompany.com
Date: Sun, 18 Jun 2006 02:59:07 +0200
X-tenshi-version: 0.4
X-tenshi-hostname: ns30074
Subject: tenshi report [misc]


ns30074:
1: sshd: error: PAM: Authentication failure for root from kiko.adsl.nerim.net
1: /USR/SBIN/CRON: (root) CMD (/usr/local/rtm/bin/rtm 41 >/dev/null 2>/dev/null)
1: CRON: (pam_unix) session opened for user root by (uid=0)
```

The log analysis system is functional; all that remains now is to adapt the configuration according to your needs and the time you want to allocate to reading the reports from it.

If you would like additional information, two mailing lists are available:

The first is for general discussions; to subscribe, simply send a message to tenshi-user+subscribe@lists.inversepath.com and messages for the list should be sent to tenshi-user@lists.inversepath.com

The second is dedicated to Tenshi developments; subscribe at tenshi-announce+subscribe@lists.inversepath.com and list messages at tenshi-announce@lists.inversepath.com

## References

http://www.secuobs.com/news/07042008-tenshi.shtml
