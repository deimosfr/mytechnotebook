---
weight: 999
url: "/Mise_en_place_de_certificats_SSL_sous_Apache_2/"
title: "Setting up SSL certificates with Apache 2"
description: "A guide on how to create, configure and implement SSL certificates with Apache 2 on Debian and OpenBSD systems"
categories: ["Debian", "Linux", "Apache"]
date: "2009-09-24T12:28:00+02:00"
lastmod: "2009-09-24T12:28:00+02:00"
tags: ["SSL", "Apache", "Security", "OpenBSD", "Certificates"]
toc: true
---

## Introduction

SSL certificates are not always easy to understand and implement. Nevertheless, I will try to make it clear. For those who wish to be signed by a free certification authority, I invite you to visit the [CACert](https://www.cacert.org/) website.

## Installation

### Debian

Once again, it's quite simple here:

```bash
apt-get install openssl
```

Then you'll download a small program that will make your life easier:

```bash
cd ~/
mkdir ssl
```

[Here is the file to download](https://wiki.deimos.fr/File:Cert_manager.tar.gz.html) and put in your ssl folder.

Once done, decompress it:

```bash
tar -xzvf cert_manager.tar.gz
```

### OpenBSD

For OpenBSD, nothing special to install. Apache is provided as standard with the SSL module.

## Configuration

## Debian

Let's go to the decompressed folder:

```bash
cd CERT_MANAGER
```

If you want to change the number of days for the validity of your certificate, edit the ca_openssl.cnf file and modify line 39:

```bash
default_days    = 365
```

Change 365 to what you want (3650 for 10 years for example).

### Creating a certification authority

We will now create the certification authority for our local network. First, we initialize our certificate management environment:

```bash
./cert_manager.sh --init
```

This command creates the necessary folders and files for the proper functioning of our script and asks you the necessary questions for the configuration of your certification authority:

```bash
You will now be asked to give informations for your certificate authority.
Description du domaine [défaut : Domaine local]: Deimos
Code de votre pays [défaut : FR]:
Nom de votre région [défaut : Ile de France]:
Nom de votre ville [défaut : Paris]:
Nom de votre domaine [défaut : domain.local]: www.deimos.fr
Email de l'administrateur [défaut : root@domain.local]: root@deimos.fr
```

We then create our certification authority:

```bash
./cert_manager.sh --create-ca
```

You now have what you need to sign your own certificates.

### Creating a server certificate for our local network

Now that you have a certification authority, we will create a certification request in order to obtain a certificate signed by our authority. For example, to create a certificate for our HTTPS server:

```bash
./cert_manager.sh --generate-csr=https
```

**Note: https** is used to generate the filename of the request. It is preferable that this value does not contain spaces or special characters.

You must then enter your certificate information:

```bash
You will now be asked to give informations for your certificate authority.
Type de serveur [défault : HTTP server]: HTTPS server
Code de votre pays [défaut : FR]:
Nom de votre région [défaut : Ile de France]:
Nom de votre ville [défaut : Paris]:
Email de l'administrateur [défaut : root@deimos.fr]:
Nom de votre domaine [défaut : domain.local]: www.deimos.fr
Noms de domaines supplémentaires, un par ligne. Finissez par une ligne vide.
SubjectAltName: DNS:doc.deimos.fr
SubjectAltName: DNS:imp.deimos.fr
SubjectAltName: DNS:mail.deimos.fr
SubjectAltName: DNS:
```

Note: As you can see, this tool gives you the possibility to generate certificates valid for several domain names.

At the end of the procedure, the tool displays the created request because you can, if you wish, have this request signed by the [CACert](https://www.cacert.org/) site.

If the certificate is intended for your local network, you can use your certification authority to sign it:

```bash
./cert_manager.sh --sign-csr=https
```

This command displays the information included in the certificate request and asks you if you agree to sign it:

```bash
Sign the certificate? [y/n]:y
1 out of 1 certificate requests certified, commit? [y/n]y
```

You now have 2 files that together form your certificate:

```bash
./CERTIFICATES/https_cert.pem
./PRIVATE_KEYS/https_key.pem
```

All you have to do now is copy these files and the public key of your certification authority to the appropriate place for your server configuration:

```bash
./CERTIFICATE_AUTHORITY/ca-cert.pem
```

### Integration with Apache

For the Apache part, it's quite simple, just copy certain files:

```bash
mkdir -p /etc/apache2/ssl/
cp CERTIFICATE_AUTHORITY/ca-cert.pem /etc/apache2/ssl/deimos.fr.ca.crt
cp CERTIFICATES/https_cert.cert /etc/apache2/ssl/deimos.fr.crt
cp PRIVATE_KEYS/https_key.pem /etc/apache2/ssl/deimos.fr.key
```

Then put this in a VirtualHost in Apache (make a special SSL VirtualHost):

```bash
<VirtualHost *:443>
        ServerAdmin xxx@mycompany.com
        ServerName fire.deimos.fr

        SSLEngine on
        SSLCACertificateFile /etc/apache2/ssl/deimos.fr.ca.crt
        SSLCertificateFile /etc/apache2/ssl/deimos.fr.crt
        SSLCertificateKeyFile /etc/apache2/ssl/deimos.fr.key

        DocumentRoot /var/www/

        <Directory /var/www/>
                Options Indexes FollowSymLinks MultiViews
                AllowOverride None
                Order allow,deny
                allow from all
        </Directory>
</VirtualHost>
```

### OpenBSD

To generate certificates, here's how to proceed:

```bash
cd /etc/ssl
openssl genrsa -out /etc/ssl/private/server.key 2048
openssl req -new -key /etc/ssl/private/server.key -out /etc/ssl/private/server.csr
openssl x509 -req -days 3650 -in /etc/ssl/private/server.csr -signkey /etc/ssl/private/server.key -out /etc/ssl/server.crt
```

## Apache Configuration

### Debian

To finish, add this to `/etc/apache2/ports.conf`:

```bash
Listen 443
```

All that remains is to enable the ssl mod for the configuration:

```bash
a2enmod ssl
```

We can now admire the result by restarting Apache2 :-)

```bash
/etc/init.d/apache2 restart
```

### OpenBSD

## Multi-VirtualHost SSL

Like everyone else, one day you tried to have VHosts in SSL, and [the Apache people informed you that you can't](https://httpd.apache.org/docs/2.0/ssl/ssl_faq.html#vhosts2). The problem is that you can only have one certificate because information like the Host is inside the request, and is not accessible to the layer that decides which certificate to use and handles the encryption.

There were some ugly tricks to avoid the warning in the browser, like putting a certificate for \*.mydomainofdeath.biz but it's dead if you have several domains on the server.

So here's the solution that explains that even if the OpenSSL lib doesn't support the TLS extension that you need (SNI), on the GNU side it has been supported for 2 years. Here's the long-awaited documentation:

- [Documentation on the implementation of multi-SSL with mod_gnutls](/pdf/ssl_mod_gnutls.pdf)
- [Documentation on enabling multiple HTTPS Sites For one IP using TLS extensions](/pdf/enable_multiple_https_sites_for_one_ip_using_tls_extensions.pdf)

## References

- [SSL and Certificates Documentation](/pdf/secure_websites_using_ssl_and_certificates.pdf)
- [https://www.traduc.org/docs/HOWTO/lecture/SSL-Certificates-HOWTO.html](https://www.traduc.org/docs/HOWTO/lecture/SSL-Certificates-HOWTO.html)
- [Documentation for OpenBSD](https://www.openbsd.org/faq/fr/faq10.html#HTTPS)
