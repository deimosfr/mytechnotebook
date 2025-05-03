---
weight: 999
url: "/LDAP_\\:_Installation_et_configuration_d'un_Annuaire_LDAP_(secondaire)/"
title: "LDAP: Installation and Configuration of a Secondary LDAP Directory"
description: "Guide on how to configure OpenLDAP for LDAP tree replication between 2 servers using slapd and slurpd."
categories: ["FreeBSD", "Linux", "Debian"]
date: "2008-12-26T18:36:00+02:00"
lastmod: "2008-12-26T18:36:00+02:00"
tags: ["LDAP", "OpenLDAP", "Replication", "Network", "Servers"]
toc: true
---

## Introduction

This document explains how to configure OpenLDAP for LDAP tree replication between 2 servers.

For this, we will use *slapd* and *slurpd* from the OpenLDAP suite.

This applies to both Linux distributions, tested on Debian Sarge, and FreeBSD, tested on 6.1.

## Prerequisites

You will need to configure the *slapd* daemons on both servers. For this, you can follow this document: [Standard SLAPD Configuration](/LDAP_\\:_Installation_et_configuration_d'un_Annuaire_LDAP/)

## The Master Server

On the master server, we will create a dedicated user for replication who will only have read rights but will be able to read the entire tree.

To do this, we modify the *slapd.conf* file to add the line for the user "cn=replication,o=deimos,dc=fr":

```
## Acces Lists
# Admin can change all, watcher and all authentified users can read all
access to *     by dn.regex="cn=manager,o=deimos,dc=fr" write
                by dn.regex="cn=watcher,o=deimos,dc=fr" read
                by dn.regex="cn=replication,o=deimos,dc=fr" read
                by * auth
```

Its LDIF file will be of the form:

```
# replication, deimos, net
dn: cn=replication,o=deimos,dc=fr
objectClass: top
objectClass: person
userPassword: {SSHA}hSixML09eyZsQncyqSebQq5tFpXgXT63
cn: replication
sn: LDAP replication user
```

To add it to our tree, you can use client tools.

Once this is done, we need to explain to the master server who the slave servers are and how to update them.

For this, we will add some directives to the *slapd.conf* configuration file:

```bash
 # Réplication
 # Comment joindre le serveur esclave
 replica
     # L'url pour le serveur esclave
     uri=ldaps://ldap2.deimos.fr
     # Le DN distant que l'on utilise pour mettre à jour l'arbre distant
     binddn="cn=replication,o=deimos,dc=fr"
     # La méthode d'authentification et le mot de passe associé en clair !!
     bindmethod=simple credentials=<mdp_de_replication>
 # le fichier de log de la réplication
 replogfile /var/db/openldap-slurp/replica/replog
```

Be careful with permissions on the `/var/db/openldap-slurpd/` directory!

The *slurpd* daemon must be started automatically when the server starts. There is no specific configuration other than what we added to the *slapd.conf* file.

Our master server is now properly configured.

## The Slave Server

On the slave server, not many changes need to be made to the SLAPD configuration. The main thing is that the user "cn=replication,o=deimos,dc=fr" must have write permissions to update the tree!

Then we will also need to define who the master server is.

For this, we modify the *slapd.conf* file to add the line for the user "cn=replication,o=deimos,dc=fr":

```
## Acces Lists
# Admin can change all, watcher and all authentified users can read all
access to *     by dn.regex="cn=manager,o=deimos,dc=fr" write
                by dn.regex="cn=watcher,o=deimos,dc=fr" read
                by dn.regex="cn=replication,o=deimos,dc=fr" write
                by * auth
```

We also modify the file to declare the master server:

```bash
 # Réplication
 # l'utilisateur qui met à jour l'arbre
 updatedn        "cn=replication,o=deimos,dc=fr"
 # L'URL du serveur maître
 updateref       "ldaps://openldap.deimos.fr"
```

Our slave server is now configured.

## Implementation

To put our servers into service, there are a few simple steps:

* Stop the *slapd* daemon on both servers:

```
/usr/local/etc/rc.d/slapd stop
```

* Copy the contents of the `/var/db/openldap-slapd/` directory from the master server to the `/var/db/openldap-slapd/` directory on the slave server
* Start the *slapd* daemon on the master server:

```
/usr/local/etc/rc.d/slapd start
```

* Start the *slurpd* daemon on the master server:

```
/usr/local/etc/rc.d/slurpd start
```

* Start the *slapd* daemon on the slave server:

```
/usr/local/etc/rc.d/slapd start
```

Our tree is now synchronized between the two servers.

## Final Remarks

Modifications are only possible on the master server. If we try to make modifications on a slave server, even with a DN that has write permissions, we get an error message indicating the URL of the master server.

All modifications made on the master server are replicated "instantly" to the slave servers, depending on the server loads of course!

For additions and modifications, they are done exactly as if the master server was autonomous. For read access, it can be done on all servers.

## Resources
- http://www.free-4ever.net/index.php/Openldap:configuration_slapd_replication
- [Setting up a multimaster OpenLDAP](/pdf/openldap_multimasters.pdf)
