---
weight: 999
url: "/Lutter_un_peu_plus_contre_le_SPAM_-_Règles_supplémentaires/"
title: "Fight Spam More Effectively - Additional Rules"
description: "A guide on how to enhance spam filtering with additional rules and tools like Razor and Pyzor."
categories: ["Linux", "Servers", "Network"]
date: "2007-11-22T07:15:00+02:00"
lastmod: "2007-11-22T07:15:00+02:00"
tags: ["SpamAssassin", "Email", "Spam", "Security", "Servers", "Linux"]
toc: true
---

## Installation

Install Razor, Pyzor and Configure SpamAssassin.

Razor and Pyzor are spamfilters that use a collaborative filtering network. To install them, run:

```bash
apt-get install razor pyzor
```

## Configuration

Now we have to tell SpamAssassin to use these programs. Edit `/etc/spamassassin/local.cf` so that it looks like this:

```bash
# SpamAssassin Configuration
rewrite_header Subject  *****SPAM*****
use_bayes               1
bayes_auto_learn        1
required_score          5.0
skip_rbl_checks         0
report_safe             0

#pyzor
use_pyzor               1
pyzor_path /usr/bin/pyzor

#razor
use_razor2              1
razor_config /etc/razor/razor-agent.conf

ok_locales              en fr
whitelist_from *@deimos.fr noreply@lists.silicon.fr
blacklist_from *@mandrivaclub.com
```

Note: Here is an [automatic SpamAssasin Configuration Generator](https://www.yrex.com/spam/spamconfig.php).

Then, run:

```bash
/etc/init.d/amavis restart
```

## Custom Rules

Now I want to insert some custom rulesets that can be found on the internet into SpamAssassin. I have tested those rulesets, and they make spam filtering a lot more effective.

Create the file `/usr/local/sbin/sa_rules_update.sh`:

```bash
#!/bin/sh
PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
cd /etc/spamassassin/ &> /dev/null && /usr/bin/wget http://www.rulesemporium.com/rules/71_sare_redirect_pre3.0.0.cf -O 71_sare_redirect_pre3.0.0.cf &> /dev/null
cd /etc/spamassassin/ &> /dev/null && /usr/bin/wget http://www.rulesemporium.com/rules/70_sare_bayes_poison_nxm.cf -O 70_sare_bayes_poison_nxm.cf &> /dev/null
cd /etc/spamassassin/ &> /dev/null && /usr/bin/wget http://www.rulesemporium.com/rules/70_sare_html.cf -O 70_sare_html.cf &> /dev/null
cd /etc/spamassassin/ &> /dev/null && /usr/bin/wget http://www.rulesemporium.com/rules/70_sare_html4.cf -O 70_sare_html4.cf &> /dev/null
cd /etc/spamassassin/ &> /dev/null && /usr/bin/wget http://www.rulesemporium.com/rules/70_sare_html_x30.cf -O 70_sare_html_x30.cf &> /dev/null
cd /etc/spamassassin/ &> /dev/null && /usr/bin/wget http://www.rulesemporium.com/rules/70_sare_header0.cf -O 70_sare_header0.cf &> /dev/null
cd /etc/spamassassin/ &> /dev/null && /usr/bin/wget http://www.rulesemporium.com/rules/70_sare_header3.cf -O 70_sare_header3.cf &> /dev/null
cd /etc/spamassassin/ &> /dev/null && /usr/bin/wget http://www.rulesemporium.com/rules/70_sare_header_x30.cf -O 70_sare_header_x30.cf &> /dev/null
cd /etc/spamassassin/ &> /dev/null && /usr/bin/wget http://www.rulesemporium.com/rules/70_sare_specific.cf -O 70_sare_specific.cf &> /dev/null
cd /etc/spamassassin/ &> /dev/null && /usr/bin/wget http://www.rulesemporium.com/rules/70_sare_adult.cf -O 70_sare_adult.cf &> /dev/null
cd /etc/spamassassin/ &> /dev/null && /usr/bin/wget http://www.rulesemporium.com/rules/72_sare_bml_post25x.cf -O 72_sare_bml_post25x.cf &> /dev/null
cd /etc/spamassassin/ &> /dev/null && /usr/bin/wget http://www.rulesemporium.com/rules/99_sare_fraud_post25x.cf -O 99_sare_fraud_post25x.cf &> /dev/null
cd /etc/spamassassin/ &> /dev/null && /usr/bin/wget http://www.rulesemporium.com/rules/70_sare_spoof.cf -O 70_sare_spoof.cf &> /dev/null
cd /etc/spamassassin/ &> /dev/null && /usr/bin/wget http://www.rulesemporium.com/rules/70_sare_random.cf -O 70_sare_random.cf &> /dev/null
cd /etc/spamassassin/ &> /dev/null && /usr/bin/wget http://www.rulesemporium.com/rules/70_sare_oem.cf -O 70_sare_oem.cf &> /dev/null
cd /etc/spamassassin/ &> /dev/null && /usr/bin/wget http://www.rulesemporium.com/rules/70_sare_genlsubj0.cf -O 70_sare_genlsubj0.cf &> /dev/null
cd /etc/spamassassin/ &> /dev/null && /usr/bin/wget http://www.rulesemporium.com/rules/70_sare_genlsubj3.cf -O 70_sare_genlsubj3.cf &> /dev/null
cd /etc/spamassassin/ &> /dev/null && /usr/bin/wget http://www.rulesemporium.com/rules/70_sare_genlsubj_x30.cf -O 70_sare_genlsubj_x30.cf &> /dev/null
cd /etc/spamassassin/ &> /dev/null && /usr/bin/wget http://www.rulesemporium.com/rules/70_sare_unsub.cf -O 70_sare_unsub.cf &> /dev/null
cd /etc/spamassassin/ &> /dev/null && /usr/bin/wget http://www.rulesemporium.com/rules/70_sare_uri.cf -O 70_sare_uri.cf &> /dev/null
cd /etc/spamassassin/ &> /dev/null && /usr/bin/wget http://www.timj.co.uk/linux/bogus-virus-warnings.cf -O bogus-virus-warnings.cf &> /dev/null
cd /etc/spamassassin/ &> /dev/null && /usr/bin/wget http://www.yackley.org/sa-rules/evilnumbers.cf -O evilnumbers.cf &> /dev/null
cd /etc/spamassassin/ &> /dev/null && /usr/bin/wget http://www.stearns.org/sa-blacklist/random.current.cf -O random.current.cf &> /dev/null
cd /etc/spamassassin/ &> /dev/null && /usr/bin/wget http://www.rulesemporium.com/rules/88_FVGT_body.cf -O 88_FVGT_body.cf &> /dev/null
cd /etc/spamassassin/ &> /dev/null && /usr/bin/wget http://www.rulesemporium.com/rules/88_FVGT_rawbody.cf -O 88_FVGT_rawbody.cf &> /dev/null
cd /etc/spamassassin/ &> /dev/null && /usr/bin/wget http://www.rulesemporium.com/rules/88_FVGT_subject.cf -O 88_FVGT_subject.cf &> /dev/null
cd /etc/spamassassin/ &> /dev/null && /usr/bin/wget http://www.rulesemporium.com/rules/88_FVGT_headers.cf -O 88_FVGT_headers.cf &> /dev/null
cd /etc/spamassassin/ &> /dev/null && /usr/bin/wget http://www.rulesemporium.com/rules/88_FVGT_uri.cf -O 88_FVGT_uri.cf &> /dev/null
cd /etc/spamassassin/ &> /dev/null && /usr/bin/wget http://www.rulesemporium.com/rules/99_FVGT_DomainDigits.cf -O 99_FVGT_DomainDigits.cf &> /dev/null
cd /etc/spamassassin/ &> /dev/null && /usr/bin/wget http://www.rulesemporium.com/rules/99_FVGT_Tripwire.cf -O 99_FVGT_Tripwire.cf &> /dev/null
cd /etc/spamassassin/ &> /dev/null && /usr/bin/wget http://www.rulesemporium.com/rules/99_FVGT_meta.cf -O 99_FVGT_meta.cf &> /dev/null
cd /etc/spamassassin/ &> /dev/null && /usr/bin/wget http://www.nospamtoday.com/download/mime_validate.cf -O mime_validate.cf &> /dev/null
/etc/init.d/amavis restart &> /dev/null
exit 0
```

Now you have to set the executable permissions:

```bash
chmod 744 /usr/local/sbin/sa_rules_update.sh
```

Then, add it to Crontab:

```bash
0 1 * * * /usr/local/sbin/sa_rules_update.sh &> /dev/null
```

## Resources
- [Adding And Updating SpamAssassin Rulesets With RulesDuJour](/pdf/adding_and_updating_spamassassin_rulesets_with_rulesdujour.pdf)
