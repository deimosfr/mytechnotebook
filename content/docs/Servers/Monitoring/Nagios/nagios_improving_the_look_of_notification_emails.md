---
weight: 999
url: "/Nagios_\\:_Améliorer_le_look_des_emails_de_notification/"
title: "Nagios: Improving the Look of Notification Emails"
description: "How to improve the appearance of Nagios notification emails using HTML formatting instead of plain text"
categories: ["Monitoring", "Linux"]
date: "2011-04-22T15:28:00+02:00"
lastmod: "2011-04-22T15:28:00+02:00"
tags: ["Nagios", "Email", "Monitoring", "HTML", "Notifications", "Perl"]
toc: true
---

## Introduction

If you've already used standard email notification in Nagios, you'll find it functional but not really attractive. That's why a project called "[Flexible Notifications for Nagios](https://nagios.frank4dd.com/howto/nagios-flexible-notifications.htm)" was initiated.

I will explain here how to use it easily. The goal is to have HTML email notifications like this instead of plain text:

![Nagios HTML alert](/images/nagios_html_alert.avif)

You can also do other things like adding images, carbon copies, different languages... but I'll explain here only what I need.

## Installation

Very simple, create a scripts folder and put the 2 required scripts in it:

```bash
mkdir -p /etc/nagios3/scripts
cd /etc/nagios3/scripts
wget "http://nagios.frank4dd.com/howto/source/nagios_send_host_mail.pl"
wget "http://nagios.frank4dd.com/howto/source/nagios_send_service_mail.pl"
chmod 755 nagios_send_host_mail.pl nagios_send_service_mail.pl
```

Also we'll need to install perl dependencies:

```bash
aptitude install libmail-sendmail-perl librrds-perl
```

## Configuration

### commands.cfg

Here is the file containing default email commands. You just have to comment the ones shown below and add the new lines (`/etc/nagios3/commands.cfg`):

```bash
###############################################################################
# COMMANDS.CFG - SAMPLE COMMAND DEFINITIONS FOR NAGIOS 
################################################################################


################################################################################
# NOTIFICATION COMMANDS
################################################################################

# 'notify-service-by-email' command definition
define command{
  command_name  notify-host-by-email
  command_line  /bin/sleep 1
}

# 'notify-host-by-email' command definition
#define command{
#   command_name    notify-host-by-email
#   command_line    /usr/bin/printf "%b" "***** INTERNAL Nagios *****\n\nNotification Type: $NOTIFICATIONTYPE$\nHost: $HOSTNAME$\nState: $HOSTSTATE$\nAddress: $HOSTADDRESS$\nInfo: $HOSTOUTPUT$\n\nDate/Time: $LONGDATETIME$\n" | /usr/bin/mail -s "** Internal $NOTIFICATIONTYPE$ Host Alert: $HOSTNAME$ is $HOSTSTATE$ **" $CONTACTEMAIL$
#   }

# sends HTML e-mails for hosts
define command{
        command_name    host-notify-by-email
        command_line    /etc/nagios3/scripts/nagios_send_service_mail.pl -c "$CONTACTADDRESS1$" -f html -u -p "Deimos.fr Monitoring Tool"
}

# 'notify-host-by-email' command definition
#define command{
# command_name  host-notify-by-email
# command_line  /usr/bin/printf "%b" "***** Nagios *****\n\nNotification Type: $NOTIFICATIONTYPE$\nHost: $HOSTNAME$\nState: $HOSTSTATE$\nAddress: $HOSTADDRESS$\nInfo: $HOSTOUTPUT$\n\nDate/Time: $LONGDATETIME$\n$HOSTACKAUTHOR$: $NOTIFICATIONCOMMENT$\n\n$HOSTNOTES$\n" | /usr/bin/mail -s "** $NOTIFICATIONTYPE$ Host Alert: $HOSTNAME$ is $HOSTSTATE$ **" $CONTACTEMAIL$
# }

# sends HTML e-mails for services
define command{
        command_name    notify-service-by-email
        command_line    /etc/nagios3/scripts/nagios_send_service_mail.pl -c "$CONTACTADDRESS1$" -f html -u -p "Deimos.fr Monitoring Tool"
}

# 'notify-service-by-email' command definition
#define command{
#   command_name    notify-service-by-email
#   command_line    /usr/bin/printf "%b" "***** INTERNAL Nagios *****\n\nNotification Type: $NOTIFICATIONTYPE$\n\nService: $SERVICEDESC$\nHost: $HOSTALIAS$\nAddress: $HOSTADDRESS$\nState: $SERVICESTATE$\n\nDate/Time: $LONGDATETIME$\n\nAdditional Info:\n\n$SERVICEOUTPUT$" | /usr/bin/mail -s "** $Internal NOTIFICATIONTYPE$ Service Alert: $HOSTALIAS$/$SERVICEDESC$ is $SERVICESTATE$ **" $CONTACTEMAIL$
#   }
```

### nagios_send_host_mail.pl and nagios_send_service_mail.pl

We now need to modify a few things like URL, names etc. Let's open the files and set custom fields:

```bash
 ...
 my $mail_sender        = "Nagios Monitoring <xxx@mycompany.com>";
 my $nagios_cgiurl      = "http://nagios/nagios3/cgi-bin";
 my $o_smtphost         = "127.0.0.1";
 ...
```

Those fields correspond to:

- $mail_sender: Email sender name
- $nagios_cgiurl: The CGI URL to nagios server. It's needed for links in the mail
- $o_smtphost: The SMTP server

That's all. Reload your Nagios now.
