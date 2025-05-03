---
weight: 999
url: "/Mise_en_place_d'un_Antivirus_(ClamAV_et_Amavis)/"
title: "Setting up an Antivirus (ClamAV and Amavis)"
description: "This guide explains how to set up ClamAV antivirus with Amavis to integrate with Postfix for email scanning."
categories: ["Linux", "Database"]
date: "2008-09-24T11:46:00+02:00"
lastmod: "2008-09-24T11:46:00+02:00"
tags: ["Antivirus", "ClamAV", "Amavis", "Postfix", "Email", "Security", "Servers", "Network"]
toc: true
---

## Introduction

ClamAV is the antivirus component while Amavis is the interface that connects Postfix with add-ons such as antispam and antivirus tools.

## Installation

First, let's install what we need:

```bash
apt-get install amavisd-new clamav clamav-daemon zoo unzip unzoo bzip2
```

At the end of the installation, it will ask you some questions. Here are the answers you should provide:

```bash
Virus database update method: <-- daemon
Local database mirror site: <-- db.fr.clamav.net (France; select the mirror that is closest to you)
HTTP proxy information (leave blank for none): <-- (blank)
Should clamd be notified after updates? <-- Yes
```

## Configuration

Next we'll add these lines to `/etc/postfix/main.cf`:

```bash
# Use Amavis
content_filter = amavis:[127.0.0.1]:10024
receive_override_options = no_address_mappings
```

Now let's edit `/etc/postfix/master.cf` and add:

```bash
amavis unix - - n - 2 smtp
        -o smtp_data_done_timeout=1200
        -o disable_dns_lookups=yes

127.0.0.1:10025 inet n - n - - smtpd
        -o content_filter=
        -o local_recipient_maps=
        -o relay_recipient_maps=
        -o smtpd_restriction_classes=
        -o smtpd_client_restrictions=
        -o smtpd_helo_restrictions=
        -o smtpd_sender_restrictions=
        -o smtpd_recipient_restrictions=permit_mynetworks,reject
        -o mynetworks=127.0.0.0/8
        -o strict_rfc821_envelopes=yes
```

We restart Postfix and it's working :-)

```bash
/etc/init.d/postfix restart
```

You can test it using [the Eicar website](https://www.eicar.org/anti_virus_test_file.htm).

## FAQ

### warning: connect to transport amavis: No such file or directory

Here's the command to reallocate emails:

```bash
postsuper -r ALL
```

Then check with:

```bash
postqueue -p
```

or

```bash
mailq
```

## Resources
- [ASSP: Documentation on implementing a simplification and enhancement tool for SPAM and Virus detection](/pdf/assp_with_embedded_clamav_integrated_into_postfix.pdf)
- [ClamAV, the antivirus from the cold](/pdf/clamav,_l'antivirus_qui_vient_du_froid.pdf)
