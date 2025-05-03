---
weight: 999
url: "/Xen_et_vserver_\\:_monitoring_des_VM_sur_une_page_PHP/"
title: "Xen and vserver: monitoring VMs on a PHP page"
description: "A guide to monitoring virtual machines using a PHP page to easily view VM activity without running multiple commands."
categories: ["Monitoring", "Apache", "Linux"]
date: "2007-05-14T10:09:00+02:00"
lastmod: "2007-05-14T10:09:00+02:00"
tags: ["Servers", "PHP", "Xen", "Virtualization", "Monitoring"]
toc: true
---

## Introduction

I used this kind of script because it's simpler to open a small web page to view VM activity rather than launching a bunch of commands.

## Sudo

You need sudo because by default Apache doesn't have the necessary rights to execute these commands:

```bash
apt-get install sudo
```

Then edit `/etc/sudoers` and add the following:

```bash
www-data        ALL=NOPASSWD: /usr/sbin/vserver-stat,/usr/sbin/xm list
```

## Php and script

You obviously need to have Apache and PHP installed for this to work:

```bash
apt-get install apache2 php5
```

Additionally, there's a small JavaScript that refreshes the page every 60 seconds.

Then create a folder and copy this into an index file:

```bash
mkdir /var/www/virtual && vi /var/www/virtual/index.php
```

Here's the content of index.php:

```html
<!DOCTYPE html PUBLIC "-//W3C//DTD HTML 4.01 Transitional//EN">
<html>
  <script language="javascript" type="text/javascript">
    setTimeout("location.reload();", 600000);
  </script>
  <head>
    <meta content="text/html; charset=ISO-8859-1" http-equiv="content-type" />
    <title>Informations for Virtual Machines</title>
  </head>
  <body>
    <?php
       system("date +%c") ;
       echo("<BR><br />--- Vservers informations ---<br />") ;
    echo("<syntaxhighlight lang="text"
      >\n") ; system("sudo /usr/sbin/vserver-stat") ; echo("</syntaxhighlight
    >\n") ; echo("<br />") ; echo("--- Xen informations ---<br />") ;
    echo("<syntaxhighlight lang="text"
      >\n") ; system("sudo /usr/sbin/xm list") ; echo("</syntaxhighlight
    >") ; ?>
  </body>
</html>
```

Now set the proper permissions:

```bash
chown -Rf www-data. /var/www/virtual
```

All that's left is to access it :-)
