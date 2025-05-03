---
weight: 999
url: "/Trac_\\:_Mise_en_place_d'une_solution_de_Tracking/"
title: "Trac: Setting up a Tracking Solution"
description: "A guide for setting up Trac, a web-based wiki and issue tracking system that integrates with Subversion for project management."
categories: ["Linux", "Apache", "Development"]
date: "2008-01-30T17:15:00+02:00"
lastmod: "2008-01-30T17:15:00+02:00"
tags: ["Trac", "Subversion", "Apache", "Project Management", "Tracking", "Development", "Web Application"]
toc: true
---

## Introduction

[Trac](https://trac.edgewall.org/) is a wiki with tracking capabilities in a web environment. It works with Subversion and manages tickets. It also has many other very practical features to discover.

You need to install [Subversion (SVN)](Installation_et_configuration_d'un_repository_SVN.html) first. Only then can you use Trac.

## Installation

Installing Trac is quite straightforward:

```bash
aptitude install trac enscript python-docutils libapache2-mod-python
```

## Creating a Trac Project

I place my various Trac instances in ~/trac. Run the following commands to create the Trac environment:

```bash
mkdir ~/trac
cd trac
trac-admin project initenv
```

And answer the few questions asked:

* Project Name [My Project]> Project
* Database connection string [sqlite:db/trac.db]> press enter to validate
* Repository type [svn]> press enter to validate
* Path to repository [/var/svn/test]> /root/svn/project
* Templates directory [/usr/share/trac/templates]> press enter to validate

## Apache2 Configuration

Now that the Trac project has been created, Apache needs to be configured to make it accessible. Edit your virtualhost file and add the following:

```bash
Alias /trac "/usr/share/trac/htdocs"
ScriptAlias /projet /usr/share/trac/cgi-bin/trac.cgi
<location /projet>
    SetEnv TRAC_ENV "/root/trac/projet"
</location>
```

## File Permission Modifications

To make the files accessible by Apache, you need to modify the permissions, similar to what was done for the Subversion repository:

```bash
cd ~/trac
chown -R www-data:www-data projet
chmod 775 projet -R
```

## Creating Trac Users

To connect to Trac, you need to create Trac users. Like for Subversion users, we use the htpasswd2 command:

```bash
cd /etc/apache2
htpasswd -cm dav_svn.passwd pmavro
New password:********
Re-type new password:********
Adding password for user pmavro
```

And to add another user:

```bash
htpasswd -m dav_svn.passwd anonymous
New password:
Re-type new password:
Adding password for user anonymous
```

You need to edit your virtualhost file again and add the following lines:

```bash
<location /projet/login>
    AuthType Basic
    AuthName "Trac: login"
    AuthUserFile /etc/apache2/dav_svn.passwd
    Require valid-user
</location>
```

## Finalization

To finish, let's add some security to the files containing passwords:

```bash
cd /var/www/trac
chown www-data:www-data *
chmod go-rwx *
```

## Resources
- [Documentation on setting up Subversion and Trac](/pdf/trac.pdf)
