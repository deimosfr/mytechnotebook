---
weight: 999
url: "/SSL_\\:_Gestion_des_certificats/"
title: "SSL: Certificate Management"
description: "Guide on managing SSL certificates, including generation, renewal, and application of certificates for Courier (POP3/IMAP) servers."
categories: ["Linux", "Servers", "Security"]
date: "2007-07-08T21:33:00+02:00"
lastmod: "2007-07-08T21:33:00+02:00"
tags: ["SSL", "Certificates", "Courier", "Security", "POP3", "IMAP"]
toc: true
---

## Problem Statement

After a year of good and loyal service, your **Courier (POP3 or IMAP)** server fails due to a simple SSL problem! Yes, after a year, certificates expire!

## Preparation

We need to generate new certificates. First, let's go to the right location:

```bash
cd /etc/courier/
```

Then, we delete the old one:

```bash
rm pop3d.pem
```

## Generation

### Automatic

If you decide to simply renew this certificate every year, edit the ".cnf" file and fill it out correctly. Here's an example:

```bash
RANDFILE = /usr/lib/courier/pop3d.rand

[ req ]
default_bits = 1024 # Use 2056 if you're paranoid
encrypt_key = yes
distinguished_name = req_dn
x509_extensions = cert_type
prompt = no

[ req_dn ]
C=FR
ST=France
L=Paris
O=Company
OU=Managed by Deimos System Engineer
CN=company.fr
emailAddress=admin@company.fr

[ cert_type ]
nsCertType = server
```

Then, run the certificate regeneration command:

```bash
/usr/lib/courier/mkpop3dcert
```

You should see something like this:

```
generating a 1024 bit RSA private key
...........................++++++
.++++++
writing new private key to '/usr/lib/courier/imapd.pem'
-----
1024 semi-random bytes loaded
Generating DH parameters, 512 bit long safe prime, generator 2
This is going to take a long time
.....................+.........................+..........+.....................
....+...............+........+............................................+..+..
.................................+....+................................+...+....
....................+...........................................................
.+...........................+..........+........................+..............
............+............++*++*++*++*++*++*
```

### Manual

To create your key manually, here's the command that will generate the key:

```bash
openssl genrsa -out pop3d.pem 1024
```

Replace pop3.pem with imap.pem if you're using IMAP (adapt as needed).
1024 corresponds to the number of encryption bits. Increase if necessary.

Then, you have two options:

- Self-signature
- Signature from an authority

#### Self-signature

The **-x509** option is used for self-signing:

```bash
openssl req -new -days 365 -key pop3d.pem -x509 -out pop3d.crt
```

- 365: number of days before expiration
- pop3d.pem: certificate to sign
- pop3d.crt: certificate acting as authority

#### Signature from an authority

Here's an example:

```bash
openssl req -new -days 365 -key pop3d.pem -out pop3d.crt
```

- 365: number of days before expiration
- pop3d.pem: certificate to sign
- **pop3d.crt: authoritative certificate, you should insert the certificate provided by the authority here**

## Applying New Certificates

To apply these new certificates, simply restart the appropriate services. Example:

```bash
/etc/init.d/courier-pop-ssl restart
```

## Modifying the Automatic Certificate Generation Script

As we saw above for automatic certificate generation, we run a script. But if we want to change the content slightly to have, for example, 2 or 3 years of grace period, it's convenient, even if not recommended.

Let's edit the file `/usr/lib/courier/mkpop3dcert`:

```bash
test -x /usr/bin/openssl || exit 0

prefix="/usr"

if test -f /usr/lib/courier/popd.pem
then
       echo "/usr/lib/courier/popd.pem already exists."
       exit 1
fi

umask 077 
cp /dev/null /usr/lib/courier/popd.pem
chmod 600 /usr/lib/courier/popd.pem
chown daemon /usr/lib/courier/popd.pem

cleanup() {
       rm -f /usr/lib/courier/popd.pem
       rm -f /usr/lib/courier/popd.rand
       exit 1
}

cd /usr/lib/courier
dd if=/dev/urandom of=/usr/lib/courier/popd.rand count=1 2>/dev/null
/usr/bin/openssl req -new -x509 -days 365 -nodes \
       -config /etc/courier/popd.cnf -out /usr/lib/courier/popd.pem -keyout /usr/lib/courier/popd.pem || cleanup
/usr/bin/openssl gendh -rand /usr/lib/courier/popd.rand 512 >>/usr/lib/courier/popd.pem || cleanup
/usr/bin/openssl x509 -subject -dates -fingerprint -noout -in /usr/lib/courier/popd.pem || cleanup
rm -f /usr/lib/courier/popd.rand
```

Now that you've reached this point, you should better understand which options to modify.

## Resources
- [Create your own Certificate Authority](/pdf/unix_openssl_et_ac.pdf)
