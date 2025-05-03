---
weight: 999
url: "/Domaines_Virtuels/"
title: "Virtual Domains"
description: "How to set up virtual domains in Postfix for handling multiple domains with aliases and advanced mail routing."
categories: ["Linux", "Servers"]
date: "2009-07-31T14:14:00+02:00"
lastmod: "2009-07-31T14:14:00+02:00"
tags: ["Servers", "Postfix", "Email", "Configuration"]
toc: true
---

## Introduction

Most Postfix systems are the final destination for only a few domain names. This includes the machine name, the [IP address] of the machine, and sometimes the parent domain name. The rest of this document will refer to these domains as canonical domains. They usually correspond to Postfix's "local domain" address class.

In addition to canonical domains, Postfix can be configured to be the final destination for many other domains. These domains are called "hosted" because they are not directly associated with the machine name. These hosted domains generally correspond to Postfix's virtual alias domain address class and/or the virtual mailbox domain address class.

But wait, there's more! Postfix can be configured to be the backup MX server for other domains. In this case, Postfix is not the final destination for these domains. It keeps mail when the main MX server is not working and forwards the mail as soon as it works again.

Finally, Postfix can be configured as a mail relay machine across the Internet. Obviously, Postfix is not the final destination for these messages.

## Configuration

With the approach described in this section, each hosted domain can have its own information, email addresses, etc. However, it still uses UNIX system accounts for local deliveries.

With virtual alias domains, each hosted address is an alias of a UNIX system account or an external address. The following example shows how to use this mechanism for the example.com domain.

{{< table "table-striped table-hover" >}}
| File: `/etc/postfix/main.cf` |
| --- |
```
virtual_alias_domains = example.com ...other hosted domains...
virtual_alias_maps = hash:/etc/postfix/virtual
```
{{< /table >}}

{{< table "table-striped table-hover" >}}
| File: `/etc/postfix/virtual` |
| --- |
```
postmaster@example.com postmaster
info@example.com       joe
sales@example.com      jane
# Uncomment the following entry to implement a catch-all address
# @example.com         jim
# ...virtual aliases for other domains...
```
{{< /table >}}

Notes:

- Line 2: The virtual_alias_domains parameter tells Postfix that the example.com domain is a virtual alias domain. If you forget this, Postfix will reject mail (relay denied) or won't know how to deliver it (mail for example.com will be sent back to the machine itself).

**Never list a virtual alias domain in the mydestination domain list!**

- Lines 3-8: The `/etc/postfix/virtual` file contains the virtual aliases. With this example, mail from **postmaster@example.com** is delivered to the local postmaster while mail from sales@example.com is sent to the UNIX account jane. Mail from all other addresses in the **example.com** domain is rejected with the error message "User unknown".

- Line 10: The commented entry (text after #) shows how to implement a catch-all address that receives all mail from addresses in the example.com domain not listed in the virtual alias file. This is not without risk. Spammers try to send mail that appears to come from and is destined for any possible name. A catch-all address is likely to receive many spam messages or notifications of messages sent with an address anything@example.com.

## Applying Changes

Run the following command after modifying the virtual file:

```bash
postmap /etc/postfix/virtual
```

To check if one of your addresses works correctly, type this:

```bash
postmap -q user@example.com /etc/postfix/virtual
```

then run this command after modifying the main.cf file:

```bash
postfix reload
```

Note: Virtual aliases can correspond to a local address, an external address, or both. They don't necessarily have to correspond to UNIX system accounts on your machine.

Virtual aliases solve one problem: they allow each domain to have its own email addresses. But there's one left: each virtual address corresponds to a UNIX account. With each new address, you increase the system's UNIX accounts.

## Advanced Usage

### Redirecting All Emails to One Person

You may want to redirect all emails from a machine or simply all emails that hit this server to a specific email address for testing purposes. This is feasible. To do this, edit the Postfix configuration file:

{{< table "table-striped table-hover" >}}
| File: `/etc/postfix/main.cf` |
| --- |
```
# we match emails we'll accept with regular expressions
virtual_alias_maps = regexp:/etc/postfix/virtual
```
{{< /table >}}

Let's create a virtual file:

{{< table "table-striped table-hover" >}}
| File: `/etc/postfix/virtual` |
| --- |
```
# redirect everything to someone who can read the emails
/@/	moi@mycompany.com
```
{{< /table >}}

## FAQ

### User unknown in virtual alias table

#### Rejection from the internal server

If you receive an email like:

```
User unknown in virtual alias table
```

If it's a response from your internal server, there's a problem with the virtual file. Check that it's the same one in "main.cf" and that you've correctly executed the postmap command.

#### Rejection from the external server

If you receive an email like:

```
User unknown in virtual alias table
```

It's probably Amavis causing issues. To fix this problem, comment out this line in /etc/postfix/main.cf:

{{< table "table-striped table-hover" >}}
| File: `/etc/postfix/main.cf` |
| --- |
```
# receive_override_options = no_address_mappings
```
{{< /table >}}

## Resources
- [Secure Virtual Mailserver: Postfix + OpenLDAP + Dovecot + Jamm + SASL + SquirrelMail](/pdf/secure_virtual_mailserver.pdf)
- http://linux-attitude.fr/post/Vers-un-trou-de-ver
