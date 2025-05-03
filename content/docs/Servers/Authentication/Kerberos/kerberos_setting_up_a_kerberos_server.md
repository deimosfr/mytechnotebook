---
weight: 999
url: "/Kerberos_\\:_Mise_en_place_d'un_serveur_Kerberos/"
title: "Kerberos: Setting up a Kerberos Server"
description: "A comprehensive guide on setting up Kerberos authentication server on Linux, configuring clients and system authentication using PAM."
categories: ["Red Hat", "Debian", "Linux"]
date: "2012-10-06T01:23:00+02:00"
lastmod: "2012-10-06T01:23:00+02:00"
tags: ["Network", "Servers", "Security", "Authentication", "Windows"]
toc: true
---

![Kerberos](/images/logo_kerberos_consortium.avif)

{{< table "table-hover table-striped" >}}
|||
|-|-|
| **Software version** | 5 |
| **Operating System** | Red Hat 6<br />Debian 6 |
| **Last Update** | 06/10/2012 |
{{< /table >}}

## Introduction

Kerberos is a network authentication protocol that relies on a secret key mechanism (symmetric encryption) and the use of tickets, rather than clear text passwords, thus avoiding the risk of fraudulent interception of user passwords. Created at the Massachusetts Institute of Technology (MIT), it bears the Greek name for Cerberus, guardian of the Underworld. Kerberos was first implemented on Unix systems.

In a simple network using Kerberos, several entities are distinguished:

- The client (C) has its own secret key Kc
- The server (S) also has a secret key Ks
- The ticket-granting service (TGS) has a secret key KTGS and knows the secret key KS of the server
- The key distribution center (KDC) knows the secret keys KC and KTGS

Client C wants to access a service offered by server S.

![Kerberos-simple](/images/kerberos-simple.avif)[^1]

We will first see how to set up a Kerberos server under GNU/Linux. Then in the second part, we will look at client configuration and system authentication via PAM.

## Server Installation

To install Kerberos:

```bash
aptitude install krb5-kdc krb5-admin-server
```

## Server Configuration

The krb5.conf file will need to be configured on all clients. It indicates the different realms and their respective KDC (Key Distribution Center = Kerberos server). Edit /etc/krb5.conf and adapt to your configuration:

```bash
[libdefaults]
        default_realm = EXAMPLE.COM
        ...
[realms]
        EXAMPLE.COM = {
                kdc = localhost
                admin_server = localhost
                default_domain = example.com
        }
        ...
[domain_realm]
        .example.com = EXAMPLE.COM
        example.com = EXAMPLE.COM
        ...
```

The kdc.conf file contains the Kerberos server configuration:

```bash
[realms]
EXMAPLE.COM = {
        ...
}
```

### Creating the Kerberos Database

The creation of the Kerberos database is done via the following command (the -s option allows storage in a file):

```bash
kdb5_util create -s
```

The password requested here will be used to encrypt the database. From now on, we can verify access to the KDC via the kadmin.local command. This is identical to the kadmin command but bypasses the root ACLs (local use only).

```bash
kadmin.local
```

### Creating Accounts

We can already check the main accounts created by default:

```bash
kadmin.local:  listprincs
K/M@EXAMPLE.COM
kadmin/admin@EXAMPLE.COM
kadmin/changepw@EXAMPLE.COM
kadmin/history@EXAMPLE.COM
krbtgt/EXAMPLE.COM@EXAMPLE.COM
```

User creation is done via the ank command: Add New key

```bash
kadmin.local:  ank admin/admin
```

The key must then be stored in a special file called keytab:

```bash
kadmin.local:  ktadd -k /etc/krb5kdc/kadm5.keytab kadmin/admin kadmin/changepw
```

The file /etc/krb5kdc/kadm5.keytab should now contain the corresponding keys.

Finally, set up ACLs to give all privileges to accounts with an admin instance. Edit /etc/krb5kdc/kadm5.acl:

```
*/admin@EXAMPLE.COM *
```

### Server Launch

Start the server as follows:

```bash
/etc/init.d/krb5-admin-server restart
/etc/init.d/krb5-kdc restart
```

## Client Installation

```bash
apt-get install libpam-krb5 krb5-user
```

## Client Configuration

Copy the /etc/krb5.conf file from the server.

### Tests

To test that everything works correctly, you should be able to perform the following sequence:

Obtaining a ticket for the admin principal:

```bash
$ kinit admin/admin@EXAMPLE.COM
Password for admin/admin@EXAMPLE.COM:
```

Display of current tickets:

```bash
$ klist

Ticket cache: FILE:/tmp/krb5cc_0
Default principal: admin/admin@EXAMPLE.COM

Valid starting     Expires            Service principal
06/07/06 11:53:47  06/07/06 21:53:11  krbtgt/EXAMPLE.COM@EXAMPLE.COM
```

Destroying the ticket:

```bash
kdestroy
```

## Setting up System Authentication

### PAM Configuration

On the client, we will use PAM. To do this, add the following lines to the different files. Edit the file /etc/pam.d/common-auth:

```
auth        sufficient    pam_krb5.so use_first_pass
```

Edit /etc/pam.d/common-account:

```
account     [default=bad success=ok user_unknown=ignore service_err=ignore system_err=ignore] pam_krb5.so
```

Edit /etc/pam.d/common-password:

```
password    sufficient    pam_krb5.so use_authtok
```

Edit /etc/pam.d/common-session:

```
session     optional      pam_krb5.so
```

### Adding a User

On the server, create a user named olivier:

```bash
kadmin
kadmin:  ank olivier
```

Now we can do:

```bash
kinit olivier@EXAMPLE.COM
```

Let's now create the user olivier on the client:

```bash
useradd olivier
```

Edit the /etc/shadow file:

```
olivier:*K*:13306:0:99999:7:::
```

The encrypted password here, *K*, is used to indicate that the password comes from Kerberos.

### Test

From a third machine, SSH to the Kerberos client:

```bash
ssh olivier@client
```

By doing a tail -f /var/log/auth.log on the server, you should get:

```bash
Jun  8 10:24:03 192.168.5.7 sshd[18175]: (pam_unix) check pass; user unknown 
Jun  8 10:24:03 192.168.5.7 sshd[18175]: (pam_unix) authentication failure; logname= uid=0 euid=0 tty=ssh ruser= rhost=***  
Jun  8 10:24:03 ldapserver krb5kdc[602]: AS_REQ (7 etypes {18 17 16 23 1 3 2}) 192.168.5.7: NEEDED_PREAUTH: olivier@EXAMPLE.COM for krbtgt/EXAMPLE.COM@EXAMPLE.COM, Additional pre-authentication required
Jun  8 10:24:03 ldapserver krb5kdc[602]: AS_REQ (7 etypes {18 17 16 23 1 3 2}) 192.168.5.7: ISSUE: authtime 1149755043, etypes {rep=16 tkt=16 ses=16}, olivier@EXAMPLE.COM for krbtgt/EXAMPLE.COM@EXAMPLE.COM
Jun  8 10:24:03 192.168.5.7 sshd[18175]: Accepted keyboard-interactive/pam for olivier from 192.168.5.55 port 39932 ssh2 
Jun  8 10:24:03 192.168.5.7 sshd[14434]: (pam_unix) session opened for user olivier by (uid=0)
```

## References

[^1]: http://fr.wikipedia.org/wiki/Kerberos
