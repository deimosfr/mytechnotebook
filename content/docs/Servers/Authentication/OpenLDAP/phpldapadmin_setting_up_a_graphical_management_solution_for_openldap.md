---
weight: 999
url: "/PhpLDAPadmin\\:_Mise_en_place_d'une_solution_de_management_graphique_pour_OpenLDAP/"
title: "PhpLDAPadmin: Setting Up a Graphical Management Solution for OpenLDAP"
description: "A guide to install and configure PhpLDAPadmin as a web-based graphical interface for managing OpenLDAP directories."
categories: ["Linux", "Apache", "Servers"]
date: "2009-05-11T11:06:00+02:00"
lastmod: "2009-05-11T11:06:00+02:00"
tags: ["OpenLDAP", "phpLDAPadmin", "Web Interface", "LDAP", "Lighttpd"]
toc: true
---

## Introduction

Managing an OpenLDAP database is not always simple, especially when you don't know all the fields by heart (and there are so many of them).

Here's a fairly simple web-based interface to use. For those curious about other graphical interfaces (non-web), there are:

- The well-known but buggy GQ
- The less known but very good [Apache Directory Studio](https://directory.apache.org/studio/)

## Installation

Let's use the magic command:

```bash
apt-get install phpldapadmin
```

## Configuration

### Minimum Configuration

Edit the file `/etc/phpldapadmin/config.php` and adapt these lines:

```php
$ldapservers->SetValue($i,'server','name','Deimos LDAP Server');
$ldapservers->SetValue($i,'server','host','127.0.0.1');
$ldapservers->SetValue($i,'server','base',array('dc=deimos,dc=fr'));
$ldapservers->SetValue($i,'server','auth_type','session');
$ldapservers->SetValue($i,'login','dn','cn=admin,dc=deimos,dc=fr');
$ldapservers->SetValue($i,'login','pass','le_bon_mot_de_passe_a_entrer');
```

### Disabling Anonymous Account

Edit the file `/etc/phpldapadmin/config.php`. First, we don't need to allow people without accounts to have read access to the LDAP through our new interface. Here's the line to look for:

```php
/* Enable anonymous bind login. */
// $ldapservers->SetValue($i,'login','anon_bind',true);
```

Replace it with:

```php
$ldapservers->SetValue($i,'login','anon_bind',false);
```

### Lighttpd

A little configuration for Lighttpd? Let's go:

```perl
# Alias for phpldapadmin directory
alias.url += (
   "/phpldapadmin" => "/usr/share/phpldapadmin/htdocs"
)
```

I think this configuration is not optimal and certainly not very secure, but at least it works.
