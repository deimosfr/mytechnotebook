---
weight: 999
url: "/Create_a_PKI/"
title: "Create a PKI"
description: "Guide on setting up a Public Key Infrastructure (PKI) for creating, managing, and distributing digital certificates using OpenSSL."
categories: ["Linux", "Database", "Debian"]
date: "2015-02-26T06:57:00+02:00"
lastmod: "2015-02-26T06:57:00+02:00"
tags: ["Solaris", "Development", "Linux", "Security", "Encryption"]
toc: true
---

{{< table "table-hover table-striped" >}}
| | |
|------|------|
| Software version | 1.0.1k |
| Operating System | 8 |
| Website | [Debian Website](https://www.debian.org) |
| Last Update | 26/02/2015 |
{{< /table >}}

## Introduction

A public key infrastructure (PKI)[^1] is a set of hardware, software, people, policies, and procedures needed to create, manage, distribute, use, store, and revoke digital certificates.

In cryptography, a PKI is an arrangement that binds public keys with respective user identities by means of a certificate authority (CA). The user identity must be unique within each CA domain. The third-party validation authority (VA) can provide this information on behalf of the CA. The binding is established through the registration and issuance process. Depending on the assurance level of the binding, this may be carried out by software at a CA or under human supervision. The PKI role that assures this binding is called the registration authority (RA). The RA ensures that the public key is bound to the individual to which it is assigned in a way that ensures non-repudiation.[^2][^3]

{{< alert context="info" text="To know the best recommendation for key encryption, <a href='https://wiki.mozilla.org/Security/Server_Side_TLS#Recommended_configurations'>please look at mozilla wiki</a>." />}}

## Installation

The first thing to do is to ensure you've got openssl installed:

```bash
apt-get install openssl
```

## Generate CA

First, let's create the structure:

```bash
mkdir pki
cd pki
mkdir -p {config,certs,db/ca.db.certs}
echo '01'> db/ca.db.serial
touch db/ca.db.index
```

Create a configuration file:

```bash {linenos=table,hl_lines=[2,4],anchorlinenos=true}
[ ca ]
default_ca      = domain.fqdn

[ domain.fqdn ]
dir              = ./db
certs            = ./db
new_certs_dir    = ./db/ca.db.certs
database         = ./db/ca.db.index
serial           = ./db/ca.db.serial
RANDFILE         = ./db/ca.db.rand
certificate      = ./certs/ca.crt
private_key      = ./certs/ca.key
default_days     = 3650
default_crl_days = 30
default_md       = sha256
preserve         = no
policy           = policy_anything

[ policy_anything ]
countryName             = optional
stateOrProvinceName     = optional
localityName            = optional
organizationName        = optional
organizationalUnitName  = optional
commonName              = supplied
emailAddress            = optional
```

You absolutely need to update the domain.fqdn and to can update some values like the validity of the certificate (default_days).

We're now ready to generate the key certificate:

```bash
> openssl genrsa -des3 -out certs/ca.key 2048
Generating RSA private key, 2048 bit long modulus
.....................+++
...................................................+++
e is 65537 (0x10001)
Enter pass phrase for certs/ca.key:
Verifying - Enter pass phrase for certs/ca.key:
```

Enter a pass phrase. This will be used to generate client certificates.

Self sign it and enter the required informations like in this example:

```bash {linenos=table,hl_lines=["10-16"],anchorlinenos=true}
> openssl req -utf8 -new -x509 -days 3650 -key certs/ca.key -out certs/ca.crt
Enter pass phrase for certs/ca.key:
You are about to be asked to enter information that will be incorporated
into your certificate request.
What you are about to enter is what is called a Distinguished Name or a DN.
There are quite a few fields but you can leave some blank
For some fields there will be a default value,
If you enter '.', the field will be left blank.
-----
Country Name (2 letter code) [AU]:FR
State or Province Name (full name) [Some-State]:France
Locality Name (eg, city) []:Paris
Organization Name (eg, company) [Internet Widgits Pty Ltd]:company
Organizational Unit Name (eg, section) []:section
Common Name (e.g. server FQDN or YOUR name) []:domain.fqdn
Email Address []:email@domain.fqdn
```

## DER certificate

You can also generate a DER certificate to import into web browser:

```bash
openssl x509 -in certs/ca.crt -outform DER -out certs/ca.der
```

## Generate a certificate

We're now ready to create certificates. We can create multiple ones (one by one for each domain) or we can create a wildcard. That's what we're going to do. Generate the key:

```bash
> openssl genrsa -out certs/wildcard.mydomain.fqdn.key 2048
Generating RSA private key, 2048 bit long modulus
..........+++
............................................+++
e is 65537 (0x10001)
```

Then generate the csr. As described above, we want to create a wildcard certificate, so do no forget to add the '\*' character:

```bash {linenos=table,hl_lines=["9-15"],anchorlinenos=true}
> openssl req -days 3650 -new -key certs/wildcard.mydomain.fqdn.key -out certs/wildcard.mydomain.fqdn.csr
You are about to be asked to enter information that will be incorporated
into your certificate request.
What you are about to enter is what is called a Distinguished Name or a DN.
There are quite a few fields but you can leave some blank
For some fields there will be a default value,
If you enter '.', the field will be left blank.
-----
Country Name (2 letter code) [AU]:FR
State or Province Name (full name) [Some-State]:France
Locality Name (eg, city) []:Paris
Organization Name (eg, company) [Internet Widgits Pty Ltd]:company
Organizational Unit Name (eg, section) []:section
Common Name (e.g. server FQDN or YOUR name) []:*.domain.fqdn
Email Address []:user@domain.fqdn

Please enter the following 'extra' attributes
to be sent with your certificate request
A challenge password []:
An optional company name []:
```

Do not enter a challenge password or it will be required each time you'll want to use it.

You can now generate the final self signed certificate:

```bash
> openssl ca -config config/ca.config -out certs/wildcard.mydomain.fqdn.crt -infiles certs/wildcard.mydomain.fqdn.csr
Using configuration from config/ca.config
Enter pass phrase for ./certs/ca.key:
Check that the request matches the signature
Signature ok
The Subject's Distinguished Name is as follows
countryName          :PRINTABLE:'FR'
stateOrProvinceName  :ASN.1 12:'France'
localityName         :ASN.1 12:'Paris'
organizationName     :ASN.1 12:'company'
organizationalUnitName:ASN.1 12:'section'
commonName           :ASN.1 12:'*.domain.fqdn'
emailAddress         :IA5STRING:'user@domain.fqdn'
Certificate is to be certified until Feb 14 22:02:29 2025 GMT (3650 days)
Sign the certificate? [y/n]:y

1 out of 1 certificate requests certified, commit? [y/n]y
Write out database with 1 new entries
Data Base Updated
```

## Check certificates

You can now check SSL certificate like that:

```bash
openssl x509 -text -in certs/wildcard.mydomain.fqdn.crt
```

## References

[^1]: http://en.wikipedia.org/wiki/Public_key_infrastructure
[^2]: http://artisan.karma-lab.net/creer-sa-propre-mini-pki
[^3]: https://developer.mozilla.org/en-US/docs/Mozilla/Security/x509_Certificates
