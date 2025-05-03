---
weight: 999
url: "/Installation_et_configuration_d\\'un_repository_SVN/"
title: "Installation and Configuration of an SVN Repository"
description: "Comprehensive guide to install, configure and use a Subversion (SVN) repository with Apache on Linux systems."
categories: ["Linux", "Apache"]
date: "2012-07-06T14:40:00+02:00"
lastmod: "2012-07-06T14:40:00+02:00"
tags: ["Subversion", "Apache", "Version Control", "SVN", "Repository"]
toc: true
---

## Introduction

If you want to use [Subversion](https://en.wikipedia.org/wiki/Subversion_(software)), the successor of CVS, you first need to install Apache 2.

## Installation

Here are the packages to install:

```bash
aptitude install subversion libapache2-svn subversion-tools
```

## Configuration of SVN

### Configuration of the module

Edit the "dav_svn" module:

```bash
vi /etc/apache2/mods-enabled/dav_svn.conf
```

Then adapt this configuration for your needs (`/etc/apache2/mods-enabled/dav_svn.conf`):

```apache
 <Location /svn>
      #When the client accesses /svn the URL will be handled
      #by the directives here, thus by subversion 

      #Loading the subversion module 
      DAV svn
      # Path to your repository
      SVNPath ''/usr/local/svn
      # Path if you have multiple repositories
      SVNParentPath ''/usr/local/svn

      #Here we request authentication with password
      #use htpasswd2 to create the file
      AuthType Basic
      AuthName ''"Subversion Repository"
      AuthUserFile ''/etc/apache2/dav_svn.passwd

      #Here we only request authentication for writing
      #operations on the repository.
      <LimitExcept GET PROPFIND OPTIONS REPORT>
          Require valid-user
      </LimitExcept>
 </Location>
```

### Apache Configuration

Add this to your "VirtualHost" in Apache configuration (`/etc/apache2/site-enabled/default`):

```apache
 <Directory /usr/local/svn>
     Options Indexes FollowSymLinks MultiViews
     AllowOverride None
     Order allow,deny
     allow from all
 </Directory>
```

### Access Definition

Next, we will create a file containing users authorized to connect:

```bash
htpasswd -c /etc/apache2/dav_svn.passwd $USER
```

## Setting up the repository

* Creating a repo (repository):

```bash
svnadmin create /usr/local/svn/project
```

* Importing a project:

```bash
svn import /home/$USER/project file:///usr/local/svn/project -m "initial import"
```

If everything went well, you should see this:

```
Committed revision 1.
```

* Verification:

```bash
svn ls file:///usr/local/svn/project
```

## Starting the daemon

To start the daemon:

```bash
svnserve -d -r /usr/local/svn/project --listen-port=3690
```

Set a listening port. 3690 is the default SVN port.

* Log status:

```bash
svn log svn://localhost:3690
```

## Usage

### A repository within another repository

If, for example, you want to have a repository, then a folder inside it pointing to another repository, it's possible by going to the folder that will contain the other repositories, then running this command:

```bash
svn propedit svn:externals .
```

Then, enter the name of the folder you want and the SVN address:

```
lib svn://svnsrv/trunk/library
```

* lib: the folder containing this SVN repository
* svn://: the SVN address

Here's another example:

```
Nagios         http://svn/admin/Production/Nagios
```

You can then save this configuration to make it permanent by doing a commit.

## References

[Setting up an SVN repo](/pdf/svn.pdf)  
[Apache Subversion Documentation](/pdf/howto_apache_subversion.pdf)  
[SVN and auto updatable working copy Documentation](/pdf/subversion_with_auto_updatable_working_copy.pdf)  
[Subversion and Trac as Virtual-Hosts documentation](/pdf/subversion_and_trac_as_virtual_hosts.pdf)
