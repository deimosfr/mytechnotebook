---
weight: 999
url: "/Check_MK_\\:_Collecter_facilement_des_infos_Nagios_et_étendez_ses_possibilités/"
title: "Check MK: Easily collect Nagios information and extend its capabilities"
description: "Set up Check MK to extend Nagios functionality with easy data collection and additional features like multisites and clustering"
categories: ["Servers", "Monitoring", "Nagios"]
date: "2012-05-25T07:42:00+02:00"
lastmod: "2012-05-25T07:42:00+02:00"
tags: ["Check MK", "Nagios", "Monitoring", "Livestatus"]
toc: true
---

![Check MK](/images/check_mk_logo.avif)

{{< table "table-hover table-striped" >}}
|||
|-|-|
| **Software version** | 1.1.12p7 |
| **Operating System** | Debian 6 |
| **Website** | [Check MK Website](https://mathias-kettner.de) |
| **Last Update** | 25/05/2012 |
{{< /table >}}

## Introduction

Nagios is great, but sometimes it lacks certain features. Here's an excellent addon that allows you to quickly extract data and enhance capabilities (multisites, cluster...).

In this article, we'll explore some of these possibilities with Check MK (also known as MK Live Status).

## Prerequisites

Here are the necessary packages:

```bash
aptitude install xinetd gcc g++ libc6-dev libstdc++6-dev libapache2-mod-python
```

## Installation

### Check MK

Let's download the latest version of check_mk and the agent:

```bash
wget http://mathias-kettner.de/download/check_mk-1.1.12p7.tar.gz
wget http://mathias-kettner.de/download/check-mk-agent_1.1.12p7-2_all.deb
tar -xzf check_mk-1.1.12p7.tar.gz
cd check_mk-1.1.12p7
```

Now we'll run the installer (basically it's just hitting enter repeatedly, Windows-style):

```bash
./setup.sh

               ____ _               _        __  __ _  __
              / ___| |__   ___  ___| | __   |  \/  | |/ /
             | |   | '_ \ / _ \/ __| |/ /   | |\/| | ' /
             | |___| | | |  __/ (__|   <    | |  | | . \
              \____|_| |_|\___|\___|_|\_\___|_|  |_|_|\_\
                                       |_____|

   Check_MK setup                                 Version: 1.1.12p7



Welcome to Check_MK. This setup will install Check_MK into user defined
directories. If you run this script as root, installation paths below
/usr will be suggested. If you run this script as non-root user paths
in your home directory will be suggested. You may override the default
values or just hit enter to accept them.

Your answers will be saved to /root/.check_mk_setup.conf and will be
reused when you run the setup of this or a later version again. Please
delete that file if you want to delete your previous answers.

 * Found running Nagios process, autodetected 17 settings.


  1) Installation directories of check_mk


Executable programs
Directory where to install executable programs such as check_mk itself.
This directory should be in your search path ($PATH). Otherwise you
always have to specify the installation path when calling check_mk:
( default  --> /usr/bin):

Check_MK configuration
Directory where check_mk looks for its main configuration file main.mk.
An example configuration file will be installed there if no main.mk is
present from a previous version.:
( default  --> /etc/check_mk):

check_mk checks
check_mk's different checks are implemented as small Python scriptlets
that parse and interpret the various output sections of the agents. Where
shall those be installed:
( default  --> /usr/share/check_mk/checks):

check_mk modules
Directory for main componentents of check_mk itself. The setup will
also create a file 'defaults' in that directory that reflects all settings
you are doing right now:
( default  --> /usr/share/check_mk/modules):

Check_MK Multisite GUI
Directory where Check_mk's Multisite GUI should be installed. Multisite is
an optional replacement for the Nagios GUI, but is also needed for the
logwatch extension.  That directory should not be
in your WWW document root. A separate apache configuration file will be
installed that maps the directory into your URL schema:
( default  --> /usr/share/check_mk/web):

Localization dir
Base directory for gettext localization files. Multisite comes prepared for localzation
but does not ship any language per default.:
( default  --> /usr/share/check_mk/locale):

documentation
Some documentation about check_mk will be installed here. Please note,
however, that most of check_mk's documentation is available only online at
http://mathias-kettner.de/check_mk.html:
( default  --> /usr/share/doc/check_mk):

check manuals
Directory for manuals for the various checks. The manuals can be viewed
with check_mk -M <CHECKNAME>:
( default  --> /usr/share/doc/check_mk/checks):

working directory of check_mk
check_mk will create caches files, automatically created checks and
other files into this directory. The setup will create several subdirectories
and makes them writable by the Nagios process:
( default  --> /var/lib/check_mk):

agents for operating systems
Agents for various operating systems will be installed here for your
conveniance. Take them and install them onto your target hosts:
( default  --> /usr/share/check_mk/agents):



  2) Configuration of Linux/UNIX Agents


extensions for agents
This directory will not be created on the server. It will be hardcoded
into the Linux and UNIX agents. The agent will look for extensions in the
subdirectories plugins/ and local/ of that directory:
( default  --> /usr/lib/check_mk_agent):

configuration dir for agents
This directory will not be created on the server. It will be hardcoded
into the Linux and UNIX agents. The agent will look for its configuration
files here (currently only the logwatch extension needs a configuration file):
( default  --> /etc/check_mk):



  3) Integration with Nagios


Name of Nagios user
The working directory for check_mk contains several subdirectories
that need to be writable by the Nagios user (which is running check_mk
in check mode). Please specify the user that should own those
directories:
( autodetected  --> nagios):

User of Apache process
Check_MK WATO (Web Administration Tool) needs a sudo configuration,
such that Apache can run certain commands as root. If you specify
the correct user of the apache process here, then we can create a valid
sudo configuration for you later::
( autodetected  --> www-data):

Common group of Nagios+Apache
Check_mk creates files and directories while running as nagios.
Some of those need to be writable by the user that is running the webserver.
Therefore a group is needed in which both Nagios and the webserver are
members (every valid Nagios installation uses such a group to allow
the web server access to Nagios' command pipe)::
( default  --> nagios):

Nagios binary
The complete path to the Nagios executable. This is needed by the
option -R/--restart in order to do a configuration check.:
( autodetected  --> /usr/sbin/nagios3):

Nagios main configuration file
Path to the main configuration file of Nagios. That file is always
named 'nagios.cfg'. The default path when compiling Nagios yourself
is /usr/local/nagios/etc/nagios.cfg. The path to this file is needed
for the check_mk option -R/--restart:
( autodetected  --> /etc/nagios3/nagios.cfg):

Nagios object directory
Nagios' object definitions for hosts, services and contacts are
usually stored in various files with the extension .cfg. These files
are located in a directory that is configured in nagios.cfg with the
directive 'cfg_dir'. Please specify the path to that directory
(If the autodetection can find your configuration
file but does not find at least one cfg_dir directive, then it will
add one to your configuration file for your conveniance):
( autodetected  --> /etc/nagios3/conf.d):

Nagios startskript
The complete path to the Nagios startskript is used by the option
-R/--restart to restart Nagios.:
( autodetected  --> /etc/init.d/nagios3):

Nagios command pipe
Complete path to the Nagios command pipe. check_mk needs write access
to this pipe in order to operate:
( default  --> /var/log/nagios/rw/nagios.cmd):

Check results directory
Complete path to the directory where Nagios stores its check results.
Using that directory instead of the command pipe is faster.:
( autodetected  --> /var/lib/nagios3/spool/checkresults):

Nagios status file
The web pages of check_mk need to read the file 'status.dat', which is
regularily created by Nagios. The path to that status file is usually
configured in nagios.cfg with the parameter 'status_file'. If
that parameter is missing, a compiled-in default value is used. On
FHS-conforming installations, that file usually is in /var/lib/nagios
or /var/log/nagios. If you've compiled Nagios yourself, that file
might be found below /usr/local/nagios:
( autodetected  --> /var/cache/nagios3/status.dat):

Path to check_icmp
check_mk ships a Nagios configuration file with several host and
service templates. Some host templates need check_icmp as host check.
That check plugin is contained in the standard Nagios plugins.
Please specify the complete path (dir + filename) of check_icmp:
( autodetected  --> /usr/lib/nagios/plugins/check_icmp):



  4) Integration with Apache


URL Prefix for Web addons
Usually the Multisite GUI is available at /check_mk/ and PNP4Nagios
is located at /pnp4nagios/. In some cases you might want to define some
prefix in order to be able to run more instances of Nagios on one host.
If you say /test/ here, for example, then Multisite will be located
at /test/check_mk/. Please do not forget the trailing slash.:
( default  --> /): /check_mk

Apache config dir
Check_mk ships several web pages implemented in Python with Apache
mod_python. That module needs an apache configuration section which
will be installed by this setup. Please specify the path to a directory
where Apache reads in configuration files.:
( autodetected  --> /etc/apache2/conf.d):

HTTP authentication file
Check_mk's web pages should be secured from unauthorized access via
HTTP authenticaion - just as Nagios. The configuration file for Apache
that will be installed contains a valid configuration for HTTP basic
auth. The most conveniant way for you is to use the same user file as
for Nagios. Please enter your htpasswd file to use here:
( default  --> /etc/nagios/htpasswd.users):

HTTP AuthName
Check_mk's Apache configuration file will need an AuthName. That
string will be displayed to the user when asking for the password.
You should use the same AuthName as for Nagios. Otherwise the user will
have to log in twice:
( autodetected  --> Nagios Access):



  5) Integration with PNP4Nagios 0.6


PNP4Nagios templates
Check_MK ships templates for PNP4Nagios for most of its checks.
Those templates make the history graphs look nice. PNP4Nagios
expects such templates in the directory pnp/templates in your
document root for static web pages:
( autodetected  --> /usr/share/pnp4nagios/html/templates):

RRA config for PNP4Nagios
Check_MK ships RRA configuration files for its checks that
can be used by PNP when creating the RRDs. Per default, PNP
creates RRD such that for each variable the minimum, maximum
and average value is stored. Most checks need only one or two
of these aggregations. If you install the Check_MK's RRA config
files into the configuration directory of PNP, PNP will create
RRDs with the minimum of required aggregation and thus save
substantial amount of disk I/O (and space) for RRDs. The default
is to install the configuration into a separate directory but
does not enable them:
( default  --> /usr/share/check_mk/pnp-rraconf):



  6) Check_MK Livestatus Module


compile livestatus module
This version of Check_mk ships a completely new and experimental
Nagios event broker module that provides direct access to Nagios
internal data structures. This module is called the Check_MK Livestatus
Module. It aims to supersede status.dat and also NDO. Currenty it
is completely experimental and might even crash your Nagios process.
Nevertheless - The Livestatus Module does not only allow extremely
fast access to the status of your services and hosts, it does also
provide live data (which status.dat does not). Also - unlike NDO -
Livestatus does not cost you even measurable CPU performance, does
not need any disk space and also needs no configuration.

Please answer 'yes', if you want to compile and integrate the
Livestatus module into your Nagios. You need 'make' and the GNU
C++ compiler installed in order to do this:
( default  --> yes):

check_mk's binary modules
Directory for architecture dependent binary libraries and plugins
of check_mk:
( default  --> /usr/lib/check_mk):

Unix socket for Livestatus
The Livestatus Module provides Nagios status data via a unix
socket. This is similar to the Nagios command pipe, but allows
bidirectional communication. Please enter the path to that pipe.
It is recommended to put it into the same directory as Nagios'
command pipe:
( default  --> /var/log/nagios/rw/live):

Backends for other systems
Directory where to put backends and configuration examples for
other systems. Currently this is only Nagvis, but other might follow
later.:
( default  --> /usr/share/check_mk/livestatus):


----------------------------------------------------------------------

You have chosen the following directories:

 Executable programs             /usr/bin
 Check_MK configuration          /etc/check_mk
 check_mk checks                 /usr/share/check_mk/checks
 check_mk modules                /usr/share/check_mk/modules
 Check_MK Multisite GUI          /usr/share/check_mk/web
 Localization dir                /usr/share/check_mk/locale
 documentation                   /usr/share/doc/check_mk
 check manuals                   /usr/share/doc/check_mk/checks
 working directory of check_mk   /var/lib/check_mk
 agents for operating systems    /usr/share/check_mk/agents
 extensions for agents           /usr/lib/check_mk_agent
 configuration dir for agents    /etc/check_mk
 Name of Nagios user             nagios
 User of Apache process          www-data
 Common group of Nagios+Apache   nagios
 Nagios binary                   /usr/sbin/nagios3
 Nagios main configuration file  /etc/nagios3/nagios.cfg
 Nagios object directory         /etc/nagios3/conf.d
 Nagios startskript              /etc/init.d/nagios3
 Nagios command pipe             /var/log/nagios/rw/nagios.cmd
 Check results directory         /var/lib/nagios3/spool/checkresults
 Nagios status file              /var/cache/nagios3/status.dat
 Path to check_icmp              /usr/lib/nagios/plugins/check_icmp
 URL Prefix for Web addons       /check_mk
 Apache config dir               /etc/apache2/conf.d
 HTTP authentication file        /etc/nagios/htpasswd.users
 HTTP AuthName                   Nagios Access
 PNP4Nagios templates            /usr/share/pnp4nagios/html/templates
 RRA config for PNP4Nagios       /usr/share/check_mk/pnp-rraconf
 compile livestatus module       yes
 check_mk's binary modules       /usr/lib/check_mk
 Unix socket for Livestatus      /var/log/nagios/rw/live
 Backends for other systems      /usr/share/check_mk/livestatus


Proceed with installation (y/n)? y
(Compiling MK Livestatus.......................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................)
Installation completed successfully.
Please restart Nagios and Apache in order to update/active check_mk's web pages.

You can access the new Multisite GUI at http://localhost/check_mkcheck_mk/
```

Let's install the agent package:

```bash
dpkg -i check-mk-agent_1.1.12p7-2_all.deb
```

We'll create the socket directory and set the appropriate permissions:

```bash
mkdir -p /var/lib/nagios3/rw/
touch /var/lib/nagios3/rw/live
chown -Rf nagios. /var/lib/nagios3
```

### Xinetd

Configure xinetd:

```bash {linenos=table}
service livestatus
{
    type = UNLISTED
    port = 6557
    socket_type = stream
    protocol = tcp
    wait = no
    cps = 100 3
    instances = 500
    per_source = 250
    flags = NODELAY
    user = nagios
    server = /usr/bin/unixcat
    server_args = /var/lib/nagios3/rw/live
    only_from = 127.0.0.1 # modify this to only allow specific hosts to connect, currenly localhost only
    disable = no
}
```

And restart:

```bash
/etc/init.d/xinetd restart
```

### Apache

First, let's create an htpasswd file that will contain the authorized users. You can later connect it to LDAP or another authentication system:

```bash
> htpasswd -c /etc/nagios/htpasswd.users nagiosadmin
New password:
Re-type new password:
Adding password for user nagiosadmin
```

We'll grant additional permissions to Apache for the web interface to display properly:

```bash
usermod -G nagios -a www-data
mkdir /var/lib/check_mk/web/nagiosadmin
chown nagios:www-data /var/lib/check_mk/web/
chmod ug+rwx /var/lib/check_mk/web/
```

Restart Apache if needed as you may have new modules (like mod_python) that were just installed.

### Nagios

Verify that your Nagios configuration lines look like this (these lines should have been automatically added to the end of your configuration file):

```bash {linenos=table,hl_lines=[4]}
[...]
# Load Livestatus Module
event_broker_options=-1
broker_module=/usr/lib/check_mk/livestatus.o /var/lib/nagios3/rw/live
[...]
```

Nevertheless, check that the end of the line is correct because we didn't specify it in the setup.

Finally, restart Nagios. In your logs (syslog), you should see this:

```bash {linenos=table,hl_lines=[10]}
[...]
Apr  4 11:28:36 nagios nagios3: Nagios 3.2.1 starting... (PID=17414)
Apr  4 11:28:36 nagios nagios3: Local time is mer. avril 04 11:28:36 CEST 2012
Apr  4 11:28:36 nagios nagios3: LOG VERSION: 2.0
Apr  4 11:28:36 nagios nagios3: livestatus: Livestatus 1.1.12p7 by Mathias Kettner. Socket: '/var/log/nagios/rw/live'
Apr  4 11:28:36 nagios nagios3: livestatus: Please visit us at http://mathias-kettner.de/
Apr  4 11:28:36 nagios nagios3: livestatus: Hint: please try out OMD - the Open Monitoring Distribution
Apr  4 11:28:36 nagios nagios3: livestatus: Please visit OMD at http://omdistro.org
Apr  4 11:28:36 nagios nagios3: livestatus: Finished initialization. Further log messages go to /var/log/nagios3/livestatus.log
Apr  4 11:28:36 nagios nagios3: Event broker module '/usr/lib/check_mk/livestatus.o' initialized successfully.
```

## Usage

Using it is very simple. Here are a few examples to retrieve information:

```bash
echo 'GET services' | unixcat /var/lib/nagios3/rw/live
echo 'GET hosts' | unixcat /var/lib/nagios3/rw/live
```

For more information: http://mathias-kettner.de/checkmk_livestatus.html

And you can use the web interface here: http://nagios/check_mkcheck_mk

## Resources
- http://syslog.tv/2011/10/13/nagios3-mk-livestatus-xinetd/
