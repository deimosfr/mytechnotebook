---
weight: 999
url: "/Mise_en_place_d'une_solution_de_Syncronisation_entre_ActiveDirectory_et_OpenLDAP/"
title: "Setting up a Synchronization Solution between ActiveDirectory and OpenLDAP"
description: "This guide explains how to set up a synchronization solution between Active Directory and OpenLDAP, including server configuration, client setup, and password synchronization mechanisms."
categories: ["Red Hat", "Debian", "Security"]
date: "2010-04-12T08:05:00+02:00"
lastmod: "2010-04-12T08:05:00+02:00"
tags:
  [
    "3.2.4 Solaris",
    "cd ~",
    "View source",
    "search",
    "What links here",
    "Servers",
    "3.1 Serveur",
    "Network",
    "3.2.2.2 Méthode 2",
    "Schemas Microsoft",
  ]
toc: true
---

## Introduction

This documentation is quite technically advanced, which is why **you need some basic knowledge** before attempting it. I recommend reading this [documentation](./LDAP_:_Installation_et_configuration_d'un_Annuaire_LDAP) first.

For authentication, we want to use an OpenLDAP directory as a **cache/backup of an Active Directory (AD) server**. Authentication will be performed from different applications and Unix-type machines. The goal is also to be able to continue authenticating in case of an AD crash.

Note: Because objectively, in case of a crash, here are the approximate reinstallation and restore times:

- Linux: ~30 min
- Windows: ~4 h

The choice is quickly made, and this solution truly holds up. (Pure Windows users who say Windows never crashes can leave now!)

Some other information can be stored locally in OpenLDAP.

To simplify installation on the Unix side (Linux or Solaris), Microsoft SFU is installed on AD, which allows maintaining Unix account information such as gid/uid/homedirectory, shell... (this will modify the AD schema).

To access the user passwords contained in AD (unicodePwd), **the connection between the OpenLDAP server and AD will be encrypted via SSL**. This is **the only way to send requests to AD that touch this attribute**. PS: This does not allow reading it, only modifying it.

This document will present the installation for the most advanced configuration. This allows you to have an independent OpenLDAP server with content synchronized in different ways: by a Python script for all attributes of all entries, by SSOD for Unix and Samba passwords, and by PAM for Unix and Samba passwords.

## Installation

### Server

For the installation, we'll use a Debian system:

```bash
apt-get install ldap-server ldap-client python-ldap
```

For SSOD:

```bash
apt-get install libpam-ldap
```

For Samba:

```bash
apt-get install libpam-smbpass libnss-ldap
```

During installation, it will ask you what root password you want for OpenLDAP.

SSOD is a utility written by Microsoft with source code downloadable from these pages:

- http://www.google.fr/search?q=+ssod.tar.gz
- http://www.microsoft.com/windowsserver2003/R2/unixcomponents/idmu.mspx

However, the code provided in this file is not portable and does not work on 64-bit systems... I had to modify these sources to make it portable.

### Windows

For Windows, you need a DNS server, a configured AD, and the small [SFU](https://www.google.fr/search?q=microsoft+SFU) that you will need to download and install.

During installation, you can perform a custom installation and check only:

- Password synchronization
- NIS server

Then start the installation.

- Launch the Unix Services Configuration utility and configure it like this:

![Sfu 1.png](/images/sfu_1.avif)
Click apply.

- Add your OpenLDAP server here:

![Sfu 2.png](/images/sfu_2.avif)
Click apply.

- Then, if you already have a MASTER (like DELL1800 here), connect to the utility on the master and add it. Otherwise, if this server must serve as the master, add it in this MMC:

![Sfu 3.png](/images/sfu_3.avif)
Click apply.

Then you'll need to reboot and you're done.

### Clients

## Configuration

### Server

#### slapd.conf

For the following information

Edit the configuration file `/etc/ldap/slapd.conf` and add these lines:

```bash
# Schema and objectClass definitions
include         /etc/ldap/schema/core.schema
include         /etc/ldap/schema/cosine.schema
include         /etc/ldap/schema/nis.schema
include         /etc/ldap/schema/inetorgperson.schema
include         /etc/ldap/schema/microsoft.schema
include         /etc/ldap/schema/microsoft.sfu.schema
include         /etc/ldap/schema/microsoft.exchange.schema
# For samba
include         /etc/ldap/schema/samba.schema
```

So here we have the schemas so that OpenLDAP can recognize the different Microsoft attributes. We also include Samba (because we also want to integrate Samba with OpenLDAP).

The Samba schema is in `/usr/share/doc/samba-doc/examples/LDAP` and should be copied to `/etc/ldap/schema`, but first you need to install a package for Samba:

```bash
apt-get install samba-doc
cp /usr/share/doc/samba-doc/examples/LDAP/samba.schema /etc/ldap/schema/samba.schema
```

We specify where the AD synchronization base will be stored:

```bash
# Where the database file are physically stored for database #1
directory       "/var/lib/ldap/copy-ad"
```

We add indexes on uid, sn, cn, gid... fields to speed up searches for these fields:

```bash
# Indexing options for database #1
index           objectClass eq
index           cn,sn,uid,mail  pres,eq,sub
index           mailnickname,userprincipalname,proxyaddresses  pres,eq,sub
```

Now configure your AD domain and the encrypted password (see [basic ldap documentation](./LDAP_:_Installation_et_configuration_d'un_Annuaire_LDAP)):

```bash
# The base of your directory in database #1
suffix          "dc=openldap,dc=mydomain,dc=local"

# rootdn directive for specifying a superuser on the database. This is needed
# for syncrepl.
rootdn          "cn=admin,dc=openldap,dc=mydomain,dc=local"
rootpw          {SSHA}V2c83+XHO/DNrUjeNjyTwAA9W+yKm/4h

access to attrs=userPassword,shadowLastChange
        by dn="cn=admin,dc=mydomain,dc=local" write
        by anonymous auth
        by self write
        by * none

access to *
        by dn="cn=admin,dc=openldap,dc=mydomain,dc=local" write
        by * read
```

Which gives something like this:

```bash
include         /etc/ldap/schema/core.schema
include         /etc/ldap/schema/cosine.schema
include         /etc/ldap/schema/nis.schema
include         /etc/ldap/schema/inetorgperson.schema
include         /etc/ldap/schema/microsoft.schema
include         /etc/ldap/schema/microsoft.sfu.schema
include         /etc/ldap/schema/microsoft.exchange.schema
include         /etc/ldap/schema/samba.schema

pidfile         /var/run/slapd/slapd.pid

argsfile        /var/run/slapd/slapd.args

loglevel        0

modulepath      /usr/lib/ldap
moduleload      back_bdb

sizelimit 500

tool-threads 1

backend         bdb
checkpoint 512 30

database        bdb

suffix          "dc=openldap,dc=mydomain,dc=local"
rootdn          rootdn "cn=admin,dc=openldap,dc=mydomain,dc=local"
rootpw          {SSHA}V2c83+XHO/DNrUjeNjyTwAA9W+yKm/4h

directory       "/var/lib/ldap/copy-ad"

dbconfig set_cachesize 0 2097152 0
dbconfig set_lk_max_objects 1500
dbconfig set_lk_max_locks 1500
dbconfig set_lk_max_lockers 1500

index       objectClass eq
index       cn,sn,uid,mail  pres,eq,sub
index       mailnickname,userprincipalname,proxyaddresses  pres,eq,sub

lastmod         on

access to attrs=userPassword,shadowLastChange
        by '''dn="cn=admin,dc=openldap,dc=mydomain,dc=local" write
        by anonymous auth
        by self write
        by * none

access to dn.base="" by * read

access to *
        by dn="cn=admin,dc=openldap,dc=mydomain,dc=local" write
        by * read
```

#### slapd.conf proxy

In case we want our LDAP server to serve as a proxy, we need to load the back_meta module:

```bash
# Where the dynamically loaded modules are stored
modulepath      /usr/lib/ldap
moduleload      back_bdb
moduleload      back_meta
```

And configure properly to contain the information from the Active Directory server:

```bash
#######################################################################
database        meta
suffix          "dc=ad,dc=mydomain,dc=local"
uri             "ldap://192.168.0.30/dc=ad,dc=mydomain,dc=local"
suffixmassage   "dc=ad,dc=mydomain,dc=local" "dc=mydomain,dc=local"
rootdn          "cn=admin,dc=ad,DC=mydomain,DC=local" # local admin account
rootpw          {SSHA}V2c83+XHO/DNrUjeNjyTwAA9W+yKm/4h # local password
#acl-authcDN    "CN=Administrateur,CN=Users,DC=mydomain,DC=local"
#acl-passwd     {SSHA}V2c83+XHO/DNrUjeNjyTwAA9W+yKm/4h

access to attrs=userPassword,shadowLastChange
        by dn="CN=Administrateur,CN=Users,DC=ad,DC=mydomain,DC=local" write
        by anonymous auth
        by self write
        by * none

access to dn.base="" by * read

# The admin dn has full write access, everyone else
# can read everything.
access to *
        by dn="CN=Administrateur,CN=Users,DC=ad,DC=mydomain,DC=local" write
        by * read

# Save the time that the entry gets modified, for database #1
lastmod         on
cachesize 20
directory /var/lib/ldap/real-ad
index       objectClass eq
index       cn,sn,uid,mail  pres,eq,sub
```

#### Schemas

We'll copy the schemas to the right place. Download the archive:  
[Microsoft Schemas](/others/microsoft_shema.tgz)

Now we decompress:

```bash
mv Microsoft_shema.tgz /etc/ldap/schema
cd /etc/ldap/schema
tar -xzvf Microsoft_shema.tgz
rm Microsoft_shema.tgz
```

We'll create the folders associated with the location where the AD base will be stored:

```bash
mkdir -p /var/lib/ldap/copy-ad
chown -Rf openldap. /var/lib/ldap/copy-ad
```

And now we can restart the LDAP server:

```bash
/etc/init.d/slapd restart
```

#### SSOD

The advantage of SSOD is that password synchronization is done in real time.  
This synchronization can be done via PAM, which allows updating any backend or password type (LDAP, shadow, Samba, Kerberos, ...)

This last point is very interesting for synchronizing passwords for Samba servers that use their own LDAP attributes to store passwords.  
Indeed, we can use both pam_smbpass and SSOD to fill the sambaNTPassword, sambaLMPassword and other associated attributes when a password is changed on the AD server.

SSOD compilation
The SSOD utility provided by Microsoft does not work in 64 bits. I therefore had to modify the source code, and it is this modified version that should be used.

##### Configuration

To compile SSOD, the packages g++ and libpam0g-dev must be installed:

```bash
apt-get install g++ libpam0g-dev
```

Download these sources [available here](/others/ssod-src.zip):

```bash
unzip ssod-src.zip
cd sfu/tripldes
make -f make3des.debian clean
make -f make3des.debian
cd ../ssod
make -f makessod.debian clean
make -f makessod.debian
```

We get a binary in bin/ssod.l52.

##### Installation

Copy the binary obtained after compiling to `/usr/bin/ssod` on the server where OpenLDAP is installed.  
Download the file [ssod-conf.tgz](/others/ssod-conf.tgz) to `/tmp` and decompress it while in the root directory:

```bash
cd /
tar zxvf /tmp/ssod-conf.tgz
```

This will install the startup and shutdown file in `/etc/init.d`, the corresponding links to the different runlevels in `/etc/rc?.d` and finally the daemon configuration file in `/etc/sso.conf`.  
The various possible options are available on this page:  
http://technet2.microsoft.com/windowsserver/en/library/3f2ac52d-e9b3-4c8a-bc1d-a4e3adde91191033.mspx?mfr=true

By default, the configuration is as follows:

```
SYNC_HOSTS=(192.168.0.30)
FILE_PATH=/etc
ENCRYPT_KEY=ABCDZ#efgh$12345
CASE_IGNORE_NAME=0
USE_NIS=0
USE_SHADOW=0
TEMP_FILE_PATH=/tmp
SYNC_USERS=ALL
SYNC_RETRIES=5
SYNC_DELAY=0
PORT_NUMBER=6677
NIS_UPDATE_PATH=bidon
```

- SYNC_HOSTS must contain the name or IP address of the AD server with which you want to synchronize.
- ENCRYPT_KEY must be identical to what is indicated in the SFU configuration on the AD server.
- PORT_NUMBER must contain the port number used by SFU on the AD server.

With this configuration file, SSOD is configured to update passwords using PAM. So you need to create a file corresponding to the SSOD service in the PAM configuration directory (`/etc/pam.d`). A priori, this file is a copy of `/etc/pam.d/passwd`. So:

```bash
cp /etc/pam.d/passwd /etc/pam.d/ssod
```

`/etc/pam.d/ssod` normally contains an include directive for the file `/etc/pam.d/common-password-ldap`. The content of the latter varies according to the desired functionality. A priori it should reference the pam-ldap and possibly pam-smbpass modules.

Once the installation of SSOD and the various PAM modules necessary for your configuration is complete, you can verify that the SSOD daemon is working correctly by modifying a password on AD and checking the file `/var/log/auth.log` on the server where SSOD is running. You should see a line like:

```
Successfully updated password via PAM User:  dummyuser
```

for each changed password.

##### Configuration libpam-ldap and libnss-pam

Refer to the chapter of the same name in the client configuration part concerning the files `/etc/nsswitch.conf`, `/etc/libnss-ldap.*` and `/etc/pam_ldap.conf`

##### Configuration/use of libpam-smbpass

For this module to work, you must configure libnss-pam, otherwise the execution of pam-smbpass.so will end with a core dump :(. This module must be able to find the information of users and/or groups that are in the LDAP directory.

The file `/etc/samba/smb.conf` must be configured to use the LDAP backend. See the chapter relating to the Samba configuration.

Libpam-smbpass allows authentication using the sambaNTPassword and/or sambaLMPassword attributes. It also allows updating these attributes when the password is changed or when a user logs in after being authenticated by a third-party module (pam_ldap for example). It is these last two functionalities that particularly interest us. Indeed, to be able to authenticate with Samba servers, these attributes must be filled and these are a priori the only ways to fill them.

To synchronize Samba passwords from those contained in the LDAP directory when a user authenticates, you must modify the file `/etc/pam.d/common-auth-ldap`. It should contain:

```bash
auth    required        pam_nologin.so
auth    sufficient      pam_unix.so nullok_secure
auth    optional        pam_smbpass.so migrate use_first_pass
auth    sufficient      pam_ldap.so use_first_pass
auth    required        pam_deny.so
```

To modify Samba passwords when a user changes their password on AD or on a Unix machine, you must modify the file `/etc/pam.d/common-password-ldap`. It should contain:

```bash
password   required   pam_smbpass.so migrate
password   sufficient pam_ldap.so try_first_pass
password   sufficient pam_unix.so try_first_pass nullok obscure min=4 max=8 md5
password   required   pam_deny.so
```

#### Synchronization script

Now at the crontab level, we'll add a small script that will allow us to do the sync. Edit the script and modify it according to your needs:  
[AD import script](/others/ad_import.tgz)

Put the script archive somewhere and decompress it:

```bash
mv Ad_import.tgz /etc/scripts/
cd /etc/scripts/
tar -xzvf Ad_import.tgz
chmod 755 ad_import.py
rm Ad_import.tgz
```

You can edit the script and modify the following variables:

- LOGLEVEL: contains an integer. The larger the integer, the more verbose the script is
- GENLDIF: 0 or 1. If it's 1, the script generates an LDIF file containing the modified LDAP entries intended for the OpenLDAP server
- SYNCHRONIZE: 0 or 1. If it's 1, the script directly modifies the entries on the OpenLDAP server
- LOGSTDOUT: 0 or 1. If it's 1, the script displays its messages on standard output.
- FULLSYNC: 0 or 1. If it's 1, the script retrieves all entries contained in AD, otherwise it only retrieves those that have been modified since the last synchronization (incremental). If you change the source AD server (or do a complete sync), you must either set this variable to 1, or set the uSNChanged attribute of the base DN in your OpenLDAP directory to 0.
- src\*: various parameters to connect to the source server.
- srcuri: URI of the source server
- srcadmin: DN of the user used for replication
- srcpassword: his password
- srcbasedn: the base of the hierarchy to synchronize.
- dst\*: same as above but for the target server
- dstsamba: where to store info about the samba domain.

Then we add it to the crontab and run it every 10 min:

```bash
*/10 * * * * python /etc/scripts/ad_import.py
```

### Clients

In general, Unix clients can authenticate using the information contained in the LDAP directory via PAM and NSS. So you just have to install and configure these components on the different systems.

#### Debian

You have the choice of installing libpam-ldapd or libpam-ldap. libpam-ldapd is newer and avoids some bugs seen in libpam-ldap. It's up to you to decide what you want :)

##### libpam-ldapd

If you opt for libnss-ldapd, then you simply need to install this and answer the questions.

```bash
aptitude install libnss-ldapd
```

##### libpam-ldap

If you have chosen to install libpam-ldap instead of libpam-ldapd, you will need to do this manually.

By default, debian creates two different configuration files for libpam-ldap and libnss-ldap. This is unnecessary since these two files will contain the same thing. You need to delete the libpam-ldap configuration files and create links from those of libnss-ldap to those of libpam-ldap:

```bash
rm /etc/pam_ldap.*
ln -s /etc/libnss-ldap.conf /etc/pam_ldap.conf
ln -s /etc/libnss-ldap.secret /etc/pam_ldap.secret
```

Edit `/etc/libnss-ldap.conf` and put in it (it should only contain these lines):

```bash
uri ldap://ldap.mydomain.local/
base dc=openldap,dc=mydomain,dc=local
ldap_version 3
rootbinddn cn=admin,dc=openldap,dc=mydomain,dc=local
scope sub
nss_paged_results yes
pagesize 1000
nss_base_passwd dc=openldap,dc=mydomain,dc=local?sub?&(&(objectClass=posixAccount)(!(objectClass=computer)))
nss_base_shadow dc=openldap,dc=mydomain,dc=local?sub
```

In `/etc/libnss-ldap.secret`, indicate the password of the user indicated on the rootdn line of the file `/etc/libnss-ldap.conf`.

You must then modify the file `/etc/nsswitch.conf` to indicate that the search will be done among others in the LDAP directory for the different services. This gives for example:

```bash
passwd:         files ldap
group:          files ldap
shadow:         files ldap

hosts:          files dns
networks:       files

protocols:      db files
services:       db files
ethers:         db files
rpc:            db files

netgroup:       nis
```

##### libpam-ldapd and libpam-ldap

Finally, you need to modify the PAM chains so that they allow authentication via the LDAP directory. This is done by copying the different files `/etc/pam.d/common-*` to `/etc/pam.d/common-*-ldap`:

```bash
cd /etc/pam.d
cp common-account{,-ldap}
cp common-auth{,-ldap}
cp common-pammount{,-ldap}
cp common-password{,-ldap}
cp common-session{,-ldap}
```

Edit the different files `/etc/pam.d/common-*-ldap` to have:

- `/etc/pam.d/common-account-ldap`:

```bash
account sufficient      pam_ldap.so
account sufficient      pam_unix.so use_first_pass
account required        pam_deny.so
```

- `/etc/pam.d/common-auth-ldap`

```bash
auth    required        pam_nologin.so
auth    sufficient      pam_unix.so nullok_secure
auth    optional        pam_mount.so debug use_first_pass
auth    optional        pam_smbpass.so migrate use_first_pass
auth    sufficient      pam_ldap.so use_first_pass
auth    required        pam_deny.so
```

- `/etc/pam.d/common-password-ldap`

```bash
password   required   pam_smbpass.so migrate
password   sufficient pam_ldap.so try_first_pass
password   sufficient pam_unix.so try_first_pass nullok obscure min=4 max=8 md5
password   required   pam_deny.so
```

- `/etc/pam.d/common-session-ldap`

```bash
session required        pam_unix.so
session required        pam_ldap.so
session required        pam_mkhomedir.so skel=/etc/skel/ umask=0022
session optional        pam_mount.so
```

- `/etc/pam.d/common-session`

```bash
session required        pam_mkhomedir.so skel=/etc/skel/ umask=0022
session sufficient      pam_ldap.so
session required        pam_unix.so
```

Once these common-\*-ldap files are created, you can edit the files of the different services for which you want to authorize authentication by LDAP... If for example you want to allow users contained in the LDAP directory to connect via SSH to the machine, you edit the file `/etc/pam.d/ssh` and replace the common-qqc with common-qqc-ldap.

For automounting partitions according to the user, edit the file `/etc/security/pam_mount.conf`:

```bash
volume * cifs 192.168.0.30 &$ /media/windows/& sfu - -
```

If you still cannot connect, restart the nscd service:

```bash
/etc/init.d/nscd restart
```

The following command should work correctly:

```bash
getent passwd mon_user
```

##### Authorize a particular LDAP group

One common method is to only authorize one or certain LDAP groups to access a machine. For this, the groups must have the **posixGroup** attribute named **login**.

###### Debian

Install libpam-modules if not already done:

```bash
aptitude install libpam-modules
```

Then add this line to the file `/etc/pam.d/common-auth`:

```bash
...
auth required    pam_access.so
```

###### Red Hat

Install this package:

```bash
yum install pam-devel
```

And add this line to the service you want (sshd for example), the restriction:

```bash {linenos=table,hl_lines=[2]}
auth       include      password-auth
account    required     pam_access.so
account    required     pam_nologin.so
```

###### Configuration

This will allow us to use the file `/etc/security/access.conf`. And here's the kind of line that needs to be added:

```bash
...
# disallow all except people in the login group and root
-:ALL EXCEPT root (sysadmin):ALL EXCEPT LOCAL
```

This allows disabling all accounts except:

- root
- The sysadmin group (not the user thanks to parentheses)
- LOCAL: local users

#### Red Hat

There are 2 methods. The first uses a Red Hat script that will do everything for us, while the second is the manual solution.

##### Method 1

To configure PAM with LDAP, use this command and adapt it to your needs:

```bash
authconfig --enableldap --enableldapauth --ldapserver=ldap://openldap-server.deimos.fr:389 --ldapbasedn="dc=openldap,dc=deimos,dc=fr" --enableldaptls --ldaploadcacer=http://serveur-web/deimosfr.crt --enablemkhomedir --update
```

- --ldapserver: enter the address of your web server
- --ldapbasedn: your server's DN
- --enableldaptls: if you use secure LDAP connections
- --ldaploadcacer: the certificate to use (if you don't have a way to retrieve it this way, look at the procedure a bit below)

Or a version without ssl/tls:

```bash
authconfig --enableldap --enableldapauth --disablenis --disableshadow --enablecache  --passalgo=sha512 --disableldaptls --disableldapstarttls --disablesssdauth --enablemkhomedir --enablepamaccess --enablecachecreds --enableforcelegacy --disablefingerprint  --ldapserver=192.168.0.1 --ldapbasedn=dc=openldap,dc=deimos,dc=fr --updateall
```

To retrieve the SSL certificate requested above, here is a solution:

```bash {linenos=table,hl_lines=[1,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24]}
> openssl s_client -connect openldap-server.deimos.fr:636
CONNECTED(00000003)
depth=0 C = FR, ST = IDF, L = Paris, O = DEIMOS, CN = openldap-server.deimos.fr, emailAddress = xxx@mycompany.com
verify error:num=18:self signed certificate
verify return:1
depth=0 C = FR, ST = IDF, L = Paris, O = DEIMOS, CN = openldap-server.deimos.fr, emailAddress = xxx@mycompany.com
verify return:1
---
Certificate chain
 0 s:/C=FR/ST=IDF/L=Paris/O=DEIMOS/CN=openldap-server.deimos.fr/emailAddress=xxx@mycompany.com
   i:/C=FR/ST=IDF/L=Paris/O=DEIMOS/CN=openldap-server.deimos.fr/emailAddress=xxx@mycompany.com
---
Server certificate
-----BEGIN CERTIFICATE-----
MIIDpTCCAw6gAwIBAgIJAJJUJLhNM1/XMA0GCSqGSIb3DQEBBQUAMIGUMQswCQYD
VQQGEwJGUjEMMAoGA1UECBMDSURGMQ4wDAYDVQQHEwVQYXJpczEPMA0GA1UEChMG
VUxMSU5LMREwDwYDVQQLEwh1bHN5c25ldDEcMBoGA1UEAxMTdGFzbWFuaWEMdWxs
aW5rLmxhbjElMCMGCSqGSIb3DQEJARYWaW503XJuYWwtaXRAdWxsaW5rLmNvbTAe
Fw0xMTEyMDUxMjQzMzVaFw0yMTEyMDIxMjQzMzVaMIGUMQswCQYDVQQGEwJGUjEM
MAoGA1UECBMDSURGMR4wDAYDVQQHEwVQYXJpczEPMA0GA1UEChMGVUxMSU5LMREw
DwYDVQQLEwh1bHN5c25ldDEcMBoGA1UEAxMTdGFzbWFuaWEudWxsaW5rLmxhbjEl
MCMGCSqGSIb3DQEJARYWaW50ZXJuYWwtaXRAdWxsaW5rLmNvbTCBnzANBgkqhkiG
9w0BAQEFAAOBjQAwgYkCgYEA4QoXFn39LhMW7mlA9r3NOX6iTHCCSlZjVQi0mQ5k
BVysN8KMFfC0E4vOeG1Z11AYwW7xCOb4Pl+LgfgfdgfgfdJIn92LX0meJcsgWKOh
qVAsZNkWn2ss8oDw3t5NEOjKFZ5BKVR2fL4Yj23DmFOAwew5PR5xhxGV5LJ9VErS
Ks0CAwEAAaOB/DCB+TAdBgNVHQ4EFgQUn5Ig2hFtROXcG3vxux7izNqcUd4wgckG
A1UdIwSBwTCBvoAUn5Ig2hFtROXcG3vxux7izNqcUd6hgZqkgZcwgZQxCzAJBgNV
BAYTAkZSMQwwCgYDVQQIEwNJREYxDjAMBgNVBAcTBVBhcmlzMQ8wDQYDVQQKEwZV
TExJTksxETAPBgNVBAsTCHVsc3lzbmV0MRwwGgYDVQQDExN0YXNtYW5pYS51bGxp
bmsubGFuMSUwIwYJKoZIhvcNAQkBFhZpbnRlcm5hbC1pdEB1bGxpbmsuY29tggkA
klQkuE0zX9cwDAYCVR0TBAUwAwEB/zANBgkqhkiG9w0BAQUFAAOBgQAbjjAbcBez
dKyq+Tlf3/DURW0BJhHKyY7UW7L39m/KZRIB2lbgFjslrAL4yNnFgipJ6aKlJFfV
BYEu7MhKH2pJZBYFpzuHOdKvDq+Kmn/wGvxeOvzh1GzQPGhQv4cClm2PJNMh/jrK
ZWNzqyLWYtWAoLu6N6gMER1Bd1Z5uzHl3A==
-----END CERTIFICATE-----
subject=/C=FR/ST=IDF/L=Paris/O=DEIMOS/CN=openldap-server.deimos.fr/emailAddress=xxx@mycompany.com
issuer=/C=FR/ST=IDF/L=Paris/O=DEIMOS/CN=openldap-server.deimos.fr/emailAddress=xxx@mycompany.com
---
No client certificate CA names sent
---
SSL handshake has read 1291 bytes and written 311 bytes
---
New, TLSv1/SSLv3, Cipher is AES256-SHA
Server public key is 1024 bit
Secure Renegotiation IS NOT supported
Compression: NONE
Expansion: NONE
SSL-Session:
    Protocol  : TLSv1
    Cipher    : AES256-SHA
    Session-ID: 91E6398F6DE9FBDC1B7EBDF890FE818B09EB79555C9FC1CF64EDC284F7A23B2A
    Session-ID-ctx:
    Master-Key: 51408932336792F4E8F5339BD12F312005022A4B20E6A5FBC56239BC0DD514344449531973B9A8395B1E799196D8F411
    Key-Arg   : None
    Krb5 Principal: None
    PSK identity: None
    PSK identity hint: None
    Start Time: 1327491823
    Timeout   : 300 (sec)
    Verify return code: 18 (self signed certificate)
---
```

In case the certificate is retrieved manually, copy it to `/etc/openldap/cacerts/ldap.crt`, then execute the following command:

```bash
cacertdir_rehash /etc/openldap/cacerts
```

##### Method 2

Modify `/etc/ldap.conf`. This file is the equivalent of `/etc/libnss_pam.conf` on debian. You can therefore put the same thing in it.

Modify the file `/etc/pam.d/system_auth`: it's the equivalent of the different common-\* under debian. This gives for example:

```bash
auth        required      /lib/security/$ISA/pam_env.so
auth        sufficient    /lib/security/$ISA/pam_unix.so likeauth nullok
auth        sufficient    /lib/security/$ISA/pam_ldap.so use_first_pass
auth        required      /lib/security/$ISA/pam_deny.so

account     sufficient    /lib/security/$ISA/pam_unix.so
account     sufficient    /lib/security/$ISA/pam_ldap.so
account     sufficient    /lib/security/$ISA/pam_succeed_if.so uid < 100 quiet
account     required      /lib/security/$ISA/pam_permit.so

password    requisite     /lib/security/$ISA/pam_cracklib.so retry=3
password    sufficient    /lib/security/$ISA/pam_unix.so nullok use_authtok md5 shadow
password    required      /lib/security/$ISA/pam_deny.so

session     optional      /lib/security/$ISA/pam_mkhomedir.so skel=/etc/skel/ umask=0077
session     required      /lib/security/$ISA/pam_limits.so
session     required      /lib/security/$ISA/pam_unix.so
```

As on debian, you must also modify the file `/etc/nsswitch.conf`.

#### Force a shell at login

If you have PAM authentication via LDAP, it is possible to force a particular shell at login. It will override the information sent by NSS and replace it with the desired shell. We will use lshell here for all people connecting via LDAP:

```bash
nss_override_attribute_value loginShell /usr/bin/lshell
```

#### Solaris

- Configure the file `/etc/pam.conf`:

For each line:

```bash
service  auth required           pam_unix_auth.so.1
```

replace "required" with "sufficient" and add behind the line:

```bash
service auth sufficient pam_ldap.so.1 try_first_pass
```

Which should give something like this:

```bash
#
#ident  "@(#)pam.conf   1.28    04/04/21 SMI"
#
# Copyright 2004 Sun Microsystems, Inc.  All rights reserved.
# Use is subject to license terms.
#
# PAM configuration
#
# Unless explicitly defined, all services use the modules
# defined in the "other" section.
#
# Modules are defined with relative pathnames, i.e., they are
# relative to /usr/lib/security/$ISA. Absolute path names, as
# present in this file in previous releases are still acceptable.
#
# Authentication management
#
# login service (explicit because of pam_dial_auth)
#
login   auth requisite          pam_authtok_get.so.1
login   auth required           pam_dhkeys.so.1
login   auth required           pam_unix_cred.so.1
login   auth sufficient         pam_ldap.so.1 try_first_pass
login   auth sufficient         pam_unix_auth.so.1
login   auth required           pam_dial_auth.so.1
#
# rlogin service (explicit because of pam_rhost_auth)
#
rlogin  auth sufficient         pam_rhosts_auth.so.1
rlogin  auth requisite          pam_authtok_get.so.1
rlogin  auth required           pam_dhkeys.so.1
rlogin  auth required           pam_unix_cred.so.1
rlogin  auth sufficient         pam_ldap.so.1 try_first_pass
rlogin  auth sufficient         pam_unix_auth.so.1
#
# Kerberized rlogin service
#
krlogin auth required           pam_unix_cred.so.1
krlogin auth binding            pam_krb5.so.1
krlogin auth sufficient         pam_ldap.so.1
krlogin auth sufficient         pam_unix_auth.so.1
#
# rsh service (explicit because of pam_rhost_auth,
# and pam_unix_auth for meaningful pam_setcred)
#
rsh     auth sufficient         pam_rhosts_auth.so.1
rsh     auth required           pam_unix_cred.so.1
#
# Kerberized rsh service
#
krsh    auth required           pam_unix_cred.so.1
krsh    auth binding            pam_krb5.so.1
krsh    auth sufficient         pam_ldap.so.1
krsh    auth sufficient         pam_unix_auth.so.1
#
# Kerberized telnet service
#
ktelnet auth required           pam_unix_cred.so.1
ktelnet auth binding            pam_krb5.so.1
ktelnet auth sufficient         pam_ldap.so.1
ktelnet auth sufficient         pam_unix_auth.so.1
#
# PPP service (explicit because of pam_dial_auth)
#
ppp     auth requisite          pam_authtok_get.so.1
ppp     auth required           pam_dhkeys.so.1
ppp     auth required           pam_unix_cred.so.1
ppp     auth sufficient         pam_ldap.so.1
ppp     auth sufficient         pam_unix_auth.so.1
ppp     auth required           pam_dial_auth.so.1
#
# Default definitions for Authentication management
# Used when service name is not explicitly mentioned for authentication
#
other   auth requisite          pam_authtok_get.so.1
other   auth required           pam_dhkeys.so.1
other   auth required           pam_unix_cred.so.1
other   auth sufficient         pam_ldap.so.1
other   auth sufficient         pam_unix_auth.so.1
#
# passwd command (explicit because of a different authentication module)
#
passwd  auth required           pam_passwd_auth.so.1
#
# cron service (explicit because of non-usage of pam_roles.so.1)
#
cron    account required        pam_unix_account.so.1
#
# Default definition for Account management
# Used when service name is not explicitly mentioned for account management
#
other   account requisite       pam_roles.so.1
other   account required        pam_unix_account.so.1
#
# Default definition for Session management
# Used when service name is not explicitly mentioned for session management
#
other   session required        pam_unix_session.so.1
#
# Default definition for  Password management
# Used when service name is not explicitly mentioned for password management
#
other   password required       pam_dhkeys.so.1
other   password requisite      pam_authtok_get.so.1
other   password requisite      pam_authtok_check.so.1
other   password required       pam_authtok_store.so.1
#
# Support for Kerberos V5 authentication and example configurations can
# be found in the pam_krb5(5) man page under the "EXAMPLES" section.
#
```

- Configure the file `/etc/nsswitch.ldap`

Leave "ldap" only where it's useful: for now on the passwd: and group: lines.  
For the rest, put the content of the file `/etc/nsswitch.dns`.  
Which gives:

```bash
#
# Copyright 2006 Sun Microsystems, Inc.  All rights reserved.
# Use is subject to license terms.
#

#
# /etc/nsswitch.dns:
#
# An example file that could be copied over to /etc/nsswitch.conf; it uses
# DNS for hosts lookups, otherwise it does not use any other naming service.
#
# "hosts:" and "services:" in this file are used only if the
# /etc/netconfig file has a "-" for nametoaddr_libs of "inet" transports.

# DNS service expects that an instance of svc:/network/dns/client be
# enabled and online.

passwd:     files ldap
group:      files ldap

# You must also set up the /etc/resolv.conf file for DNS name
# server lookup.  See resolv.conf(4).
hosts:      files dns

# Note that IPv4 addresses are searched for in all of the ipnodes databases
# before searching the hosts databases.
ipnodes:   files dns

networks:   files
protocols:  files
rpc:        files
ethers:     files
netmasks:   files
bootparams: files
publickey:  files
# At present there isn't a 'files' backend for netgroup;  the system will
#   figure it out pretty quickly, and won't use netgroups at all.
netgroup:   files
automount:  files
aliases:    files
services:   files
printers:       user files

auth_attr:  files
prof_attr:  files
project:    files

tnrhtp:     files
tnrhdb:     files
```

Once that is done, we will be able to set up the configuration. **Warning: if you are in a cluster environment, adapt to the initial configuration!**:

```bash
cp /etc/nsswitch.ldap /etc/nsswitch.conf
```

- Launch the LDAP client configuration

Just type the command:

```bash
ldapclient manual -v -a authenticationMethod=simple -a proxyDN=cn=admin,dc=openldap,dc=mydomain,dc=local -aproxyPassword=bidon -a defaultSearchBase=dc=openldap,dc=mydomain,dc=local -a defaultServerList=ldap.mydomain.local -a serviceSearchDescriptor=passwd:dc=openldap,dc=mydomain,dc=local?sub -a serviceSearchDescriptor=shadow:dc=openldap,dc=mydomain,dc=local?sub -a serviceSearchDescriptor=group:dc=openldap,dc=mydomain,dc=local?sub -a serviceAuthenticationMethod=pam_ldap:simple<nowiki>=</nowiki>
```

Attention: it seems that the ldapclient command is buggy and requires the proxyDN and proxyPassword parameters even if they are unused! (and even if they contain anything).

- Pay attention to the home directory, you must configure `/etc/auto_home` (http://www.solaris-fr.org/home/docs/base/utilisateurs). For my part, this gives:

```bash
#
# Copyright 2003 Sun Microsystems, Inc.  All rights reserved.
# Use is subject to license terms.
#
# ident "@(#)auto_home  1.6     03/04/28 SMI"
#
# Home directory map for automounter
#
+auto_home
* localhost:/export/home/&
```

In case you want to automatically create the home directory, you must port the pam_mkhomedir from Linux:

- http://mega.ist.utl.pt/~filipe/pam_mkhomedir-sol/?C=D;O=A
- http://www.keutel.de/pam_mkhomedir/index.html

A good idea would also be to automatically mount the home from an NFS server.

User accounts in the LDAP directory must have in their objectClass list the class "shadowAccount" to be taken into account by Solaris.

## Resources
- [Authentication with Linux Documentation](/pdf/ldap_authentication_in_linux.pdf)
- [Using Kerberos to Authenticate a Solaris 10 OS LDAP Client With Microsoft Active Directory](/pdf/kerberos_s10.pdf)
