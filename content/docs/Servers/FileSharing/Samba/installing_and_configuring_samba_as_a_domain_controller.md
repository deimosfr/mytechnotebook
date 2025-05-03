---
weight: 999
url: "/Installation_et_configuration_de_Samba_en_mode_Contrôleur_de_domaine/"
title: "Installing and configuring Samba as a Domain Controller"
description: "A guide on how to install and configure Samba to function as a domain controller with OpenLDAP backend."
categories: ["Linux", "Security"]
date: "2012-11-07T09:52:00+02:00"
lastmod: "2012-11-07T09:52:00+02:00"
tags: ["Samba", "Domain Controller", "LDAP", "Authentication", "Windows"]
toc: true
---

## Introduction

Samba is very versatile and can emulate a domain controller (similar to Windows NT4).

## Configuration

Here is a typical configuration for this type of environment with an OpenLDAP backend:

(`/etc/samba/smb.conf`):

```bash
[global]
    workgroup = deimos.fr
    netbios name= %h
    server string = Controleur du domaine deimos.fr
    log level = 2
    #log file = /var/log/samba/smbd.log
    log file = /var/log/samba/%m.log
    max log size = 5000
    security = user
    encrypt passwords = yes
    obey pam restrictions = No
    socket options = TCP_NODELAY SO_RCVBUF=8192 SO_SNDBUF=8192
    local master = yes
    os level = 65
    domain master = yes
    preferred master = yes
    domain logons = yes
    logon script = netlogon.vbs
    logon path =
    logon drive =
    logon home =
    wins support = yes
    dns proxy = no
    unix extensions = no

# LDAP
# Pour que Samba puisse lire et écrire dans l'annuaire : smbpasswd -w mypassword
    ldap suffix = dc=deimos.fr,dc=local
    ldap machine suffix = ou=hosts
    ldap user suffix = ou=users
    ldap group suffix = ou=groups
    ldap admin dn = uid=samba,ou=utilisateurs,dc=local
    ldap ssl = Start_tls
    ldap passwd sync = yes
    passdb backend = ldapsam:"ldap://ldap-slave1 ldap://ldap-slave2"

[netlogon]
   comment = Network Logon Service
   path = /mnt/netlogon
   browseable = no
   writable = no
   share modes = no

[homes]
    path = /datas/users/%U
    valid users = %U
    comment = %U personnal folder
    browseable = no
    writable = yes

[partage]
    path = /mnt/partage
    comment = partage
    browseable = yes
    create mask = 0700
    directory mask = 0700
    create mode = 0700
    directory mode = 0700
    writable = yes
    #valid users = @"utilisateurs du domaine"

[commons]
    path = /mnt/commons
    comment = commons
    browseable = yes
    writable = yes
    valid users = @"utilisateurs du domaine"
```

## Resources
- [Documentation on installing a Samba Domain](/pdf/samba_domaincontroller.pdf)
