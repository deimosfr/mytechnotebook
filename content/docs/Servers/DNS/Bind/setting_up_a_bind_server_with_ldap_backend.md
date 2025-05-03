---
weight: 999
url: "/Mise_en_place_d'un_serveur_Bind_en_backend_LDAP/"
title: "Setting up a Bind server with LDAP backend"
description: "This guide explains how to configure a Bind9 DNS server with LDAP backend storage instead of text files, providing dynamic updates and centralized configuration."
categories: ["Linux", "Database", "Debian"]
date: "2008-02-22T18:02:00+02:00"
lastmod: "2008-02-22T18:02:00+02:00"
tags: ["Bind", "DNS", "LDAP", "Network", "Servers"]
toc: true
---

## Introduction

This document explains the configuration of a Bind9 server as master but with storage in an LDAP server rather than in a text file.

For this, we need to use the `sdb-ldap` patch from [https://www.venaas.no/ldap/bind-sdb/](https://www.venaas.no/ldap/bind-sdb/) for Bind9. FreeBSD includes a port where the patch is directly applied to Bind9 and you just need to install it like any other port.

On Debian, you'll need to use `apt-get source bind9` and then apply the patch manually.

In the rest of this document, everything will be based on FreeBSD.

## Prerequisites

You'll need to have a functional LDAP server. For this, if you want to use OpenLDAP `slapd`, you can follow this document: [Standard SLAPD Configuration](./LDAP_:_Installation_et_configuration_d'un_Annuaire_LDAP.html)

For Bind9, I installed the `bind9-sdb-ldap` port with `WITH_PORT_REPLACES_BASE_BIND9=yes` to replace the system's Bind9.

## The named.conf file

The file will be quite similar to the one in this document: [Bind9 master configuration](Installation_et_configuration_d'un_serveur_Bind9_primaire_(Master).html). The `options` and `logging` sections will be the same. Only the zone declaration changes.

It will be in the following form:

```bash
// A standard zone
zone "deimos.fr" {
        type master;
        database "ldap ldap://127.0.0.1/ou=deimos.fr,ou=dns,o=deimos,dc=fr???? \
        !bindname=cn=dnsadmin%2cou=dns%2co=deimos%2cdc=fr,!x-bindpw=DnsPassword 172800";
        notify yes;
};
 
// A reverse zone
zone "1.168.192.in-addr.arpa" {
        type master;
        database "ldap ldap://127.0.0.1/ou=1.168.192.in-addr.arpa,ou=dns,o=deimos,dc=fr???? \
        !bindname=cn=dnsadmin%2cou=dns%2co=deimos%2cdc=fr,!x-bindpw=DnsPassword 172800";
        notify yes;
};
```

We can see that we have a `database` line that defines access to LDAP rather than the line indicating where the text file is located in a standard configuration.

If we analyze this line, we see:

* `database`: indicates that we are using a database-type storage
* `ldap` indicates that it's LDAP
* `ldap://127.0.0.1/` represents the server to connect to
* `ou=deimos.fr,ou=dns,o=deimos,dc=fr` indicates the DN where all DNS-related information is stored
* `!bindname=cn=dnsadmin%2cou=dns%2co=deimos%2cdc=fr` represents the DN used to connect to the LDAP server. Note that commas are replaced by their code "%2c".
* `x-bindpw=DnsPassword` is the password associated with the DN

With this, our Bind server will know where and how to fetch information from our LDAP server.

## The slapd.conf file

There will be few modifications to make to the `slapd` configuration.

Here is an excerpt from the configuration file:

```bash
[...]
# We add support for DNS attributes and objectclass
include         /usr/local/etc/openldap/schema/dnszone.schema
[...]
# We add indexing of DNS-related parameters
index           relativeDomainName                              eq
index           zoneName                                        eq
[...]
# We add a special user to browse this branch of the tree
# DNS admin ACL
access to dn.subtree="ou=dns,o=deimos,dc=fr"
                by dn.regex="cn=dnsadmin,ou=dns,o=deimos,dc=fr" write
                by * auth

```

Don't forget to restart your `slapd` daemon for these changes to take effect.

## LDAP directory additions

Now we'll look at what to add to our LDAP directory for our DNS to work properly.

For all manipulations with LDAP client commands, all LDIF files should be added using the `ldapadd` command.

### A container for all DNS configurations

We're going to create a branch in the tree to put all the DNS server configurations.

The LDIF will be as follows:

```bash
dn: ou=dns,o=deimos,dc=fr
objectClass: top
objectClass: organizationalUnit
ou: dns
description: All informations about DNS
```

### The dnsadmin user

In the `slapd` configuration file, we added a `dnsadmin` user.

Here is its LDIF file:

```bash
dn: cn=dnsadmin,ou=dns,o=deimos,dc=fr
objectClass: top
objectClass: person
userPassword: {SSHA}J8+mJREWzYkFDmXnZCTalBbQhq17xUzj
cn: dnsadmin
sn: dnsadmin user
```

### A container for the deimos.fr zone

Now we'll add a branch to store all configurations related to the `deimos.fr` zone.

The LDIF will be as follows:

```bash
dn: ou=deimos.fr,ou=dns,o=deimos,dc=fr
objectClass: top
objectClass: organizationalUnit
ou: deimos.fr
description: All informations about deimos.fr zone
```

### Declaration of the deimos.fr zone

Now, we can declare the `deimos.fr` zone in the LDAP tree, using our corresponding branch that we just created.

The LDIF will be as follows:

```bash
dn: zoneName=deimos.fr,ou=deimos.fr,ou=dns,o=deimos,dc=fr
objectClass: top
objectClass: dNSZone
zoneName: deimos.fr
relativeDomainName: deimos.fr
```

Our Bind9 server will now have a correspondence with what we defined in its configuration file.

### The SOA of the deimos.fr zone

This information corresponds to the beginning of a zone configuration file. It includes information about the various TTLs, the DNS servers of the zone, MX servers, etc.

The LDIF will be as follows:

```bash
dn: relativeDomainName=@,ou=deimos.fr,ou=dns,o=deimos,dc=fr
objectClass: top
objectClass: dNSZone
zoneName: deimos.fr
relativeDomainName: @
dNSTTL: 3600
dNSClass: IN
sOARecord: ns.deimos.fr. administrateur.deimos.fr. 2006112801 8H 2H 1W 1D
nSRecord: ns.deimos.fr.
mXRecord: 10 smtp.deimos.fr. 
aRecord: 192.168.1.250
tXTRecord: Zone_principale_deimos.fr
```

The `aRecord` at the end corresponds to an entry that resolves the domain name if a type `A` query is done on it.
The `tXTRecord` serves as a comment for reference... it's not mandatory... but it's always useful!

### Definition of a type A record

Now we'll declare a machine with a type `A` record that simply associates a name with an IP address.

The LDIF will be as follows:

```bash
dn: relativeDomainName=cordelia,ou=deimos.fr,ou=dns,o=deimos,dc=fr
objectClass: top
objectClass: dNSZone
zoneName: deimos.fr
relativeDomainName: orthosie
dNSTTL: 1800
dNSClass: IN
aRecord: 192.168.1.250
tXTRecord: Serveur_principal_applications
```

The parameters speak for themselves - we just declared `orthosie.deimos.fr`, its IP address is `192.168.1.250`.
Same note on the `tXTRecord`, it's always good for reference!

### Definition of a CNAME type record

Now we'll declare a `CNAME` type record that simply associates a name with another existing name. This is commonly called an alias.

The LDIF will be as follows:

```bash
dn: relativeDomainName=ns,ou=deimos.fr,ou=dns,o=deimos,dc=fr
objectClass: top
objectClass: dNSZone
zoneName: deimos.fr
relativeDomainName: ns
dNSTTL: 1800
dNSClass: CNAME
cNAMERecord: cordelia
tXTRecord: Alias_pour_le_DNS
```

The parameters are also self-explanatory here - we just declared `ns.deimos.fr`, its real name is `orthosie.deimos.fr`.
I'll skip the comment on the `tXTRecord`.

### A container for the 1.168.192.in-addr.arpa zone

Now we'll add a branch to store all configurations related to the `1.168.192.in-addr.arpa` zone.

The LDIF will be as follows:

```bash
dn: ou=1.168.192.in-addr.arpa,ou=dns,o=deimos,dc=fr
objectClass: top
objectClass: organizationalUnit
ou: 1.168.192.in-addr.arpa
description: All informations about 1.168.192.in-addr.arpa zone
```

### Declaration of the 1.168.192.in-addr.arpa zone

Now we'll declare the `1.168.192.in-addr.arpa` zone in the LDAP tree, using our corresponding branch that we just created.

The LDIF will be as follows:

```bash
dn: zoneName=1.168.192.in-addr.arpa,ou=1.168.192.in-addr.arpa,ou=dns,o=deimos,dc=fr
objectClass: top
objectClass: dNSZone
zoneName: 1.168.192.in-addr.arpa
relativeDomainName: 1.168.192.in-addr.arpa
```

Our Bind9 server will now also have a correspondence with what we defined in its configuration file.

### Definition of a PTR type record

Now we'll declare a `PTR` type record that simply associates an IP address with a name.

The LDIF will be as follows:

```bash
dn: relativeDomainName=250,ou=1.168.192.in-addr.arpa,ou=dns,o=deimos,dc=fr
objectClass: top
objectClass: dNSZone
zoneName: 1.168.192.in-addr.arpa
relativeDomainName: 250
pTRRecord: cordelia.deimos.fr.
tXTRecord: Serveur_principal_applications
```

Once again, the parameters are self-explanatory. We're associating `192.168.1.250` with the name `orthosie.deimos.fr`. Note the "." at the end of the machine name, which indicates that it's the end of the name, there's nothing to concatenate after it.

## Applying the changes

All that remains is to restart the Bind9 daemon. To do this:

```bash
# /etc/rc.d/named restart
```

In the system logs, you should see Bind9 making queries to the LDAP server.

For additions of entries to already created zones, it's dynamic! There's no need to restart the Bind9 daemon.

However, when creating a new zone, you do need to restart Bind.

## Daily administration

To facilitate the administration of all these entries in LDAP, you can use software like [phpLDAPadmin](https://phpldapadmin.sourceforge.net/), but personally, I don't like it much...

Alternatively, there are also some Perl scripts that I quickly wrote on a virtual desktop corner... You can find them here: [LDAP Administration Scripts for DHCP/DNS daemons](./file:dhcp_dns-ldapscripts.zip.html).

## Conclusion

Our DNS server is now configured with LDAP storage rather than text files.

The biggest advantage is that modifications are applied immediately without needing to restart. I'm sure we can find many other advantages... Centralized administration, etc.

The proper way to have a backup DNS server would be to configure another `slapd` in replication as explained in this document: [SLAPD Configuration in Replication](./LDAP_:_Installation_et_configuration_d'un_Annuaire_LDAP_(secondaire).html), then configure Bind on the machine in question to query the local LDAP server.

On both servers, Bind is configured as `master`!

Another point: dynamic zone updates are not possible because Bind doesn't know how to write to the LDAP tree.

## Resources
- [https://www.free-4ever.net/index.php/Dns:configuration_bind9_backend_ldap](https://www.free-4ever.net/index.php/Dns:configuration_bind9_backend_ldap)
