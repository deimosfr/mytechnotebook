---
weight: 999
url: "/Nagios_\\:_Installation_et_configuration/"
title: "Nagios: Installation and Configuration"
description: "A guide to install and configure Nagios monitoring system on Debian, including server setup, configuration options, and troubleshooting."
categories: ["Debian", "Linux", "Servers"]
date: "2014-09-08T05:57:00+02:00"
lastmod: "2014-09-08T05:57:00+02:00"
tags: ["Nagios", "Monitoring", "NRPE", "System Administration", "Network"]
toc: true
---

![Nagios](/images/nagios_logo.avif)

{{< table "table-hover table-striped" >}}
|||
|-|-|
| **Software version** | Nagios 3 |
| **Operating System** | Debian 6 |
| **Last Update** | 08/09/2014 |
{{< /table >}}

## Introduction

Nagios is a very powerful tool that allows you to test various services such as SMTP, PING, and many other things through plugins. It enables you to know if your platforms are still operational.

In short, it can alert you in different ways. Try it, you'll see, it's magical :-)

## Installation

The installation is quite simple:

```bash
aptitude install nagios3
```

## Configuration

Configuration, on the other hand, gets more complicated. You should know that Nagios is difficult to approach initially, but once you understand it, you save a lot of time and can deploy new services with ease.

### apache2.conf

This is the standard Debian configuration for apache2. It allows you to access Nagios via http://server/nagios3. It's up to you to modify the configuration or not:

```apache
# apache configuration for nagios 3.x
# note to users of nagios 1.x and 2.x:
#	throughout this file are commented out sections which preserve
#	backwards compatibility with bookmarks/config for older nagios versios.
#	simply look for lines following "nagios 1.x:" and "nagios 2.x" comments.

ScriptAlias /cgi-bin/nagios3 /usr/lib/cgi-bin/nagios3
ScriptAlias /nagios3/cgi-bin /usr/lib/cgi-bin/nagios3
# nagios 1.x:
#ScriptAlias /cgi-bin/nagios /usr/lib/cgi-bin/nagios3
#ScriptAlias /nagios/cgi-bin /usr/lib/cgi-bin/nagios3
# nagios 2.x: 
#ScriptAlias /cgi-bin/nagios2 /usr/lib/cgi-bin/nagios3
#ScriptAlias /nagios2/cgi-bin /usr/lib/cgi-bin/nagios3

# Where the stylesheets (config files) reside
Alias /nagios3/stylesheets /etc/nagios3/stylesheets
# nagios 1.x:
#Alias /nagios/stylesheets /etc/nagios3/stylesheets
# nagios 2.x:
#Alias /nagios2/stylesheets /etc/nagios3/stylesheets

# Where the HTML pages live
Alias /nagios3 /usr/share/nagios3/htdocs
# nagios 2.x: 
#Alias /nagios2 /usr/share/nagios3/htdocs
# nagios 1.x:
#Alias /nagios /usr/share/nagios3/htdocs

<DirectoryMatch (/usr/share/nagios3/htdocs|/usr/lib/cgi-bin/nagios3|/etc/nagios3/stylesheets)>
	Options FollowSymLinks

	DirectoryIndex index.php

	AllowOverride AuthConfig
	Order Allow,Deny
	Allow From All

	AuthName "Nagios Access"
	AuthType Basic
	AuthUserFile /etc/nagios3/htpasswd.users
	# nagios 1.x:
	#AuthUserFile /etc/nagios/htpasswd.users
	require valid-user
</DirectoryMatch>

# Enable this ScriptAlias if you want to enable the grouplist patch.
# See http://apan.sourceforge.net/download.html for more info
# It allows you to see a clickable list of all hostgroups in the
# left pane of the Nagios web interface
# XXX This is not tested for nagios 2.x use at your own peril
#ScriptAlias /nagios3/side.html /usr/lib/cgi-bin/nagios3/grouplist.cgi
# nagios 1.x:
#ScriptAlias /nagios/side.html /usr/lib/cgi-bin/nagios3/grouplist.cgi
```

### cgi.cfg

I'll spare you all the comments in this file and show you the configuration I use. This is the configuration for the Nagios web interface:

```bash
main_config_file=/etc/nagios3/nagios.cfg
physical_html_path=/usr/share/nagios3/htdocs
url_html_path=/nagios3
show_context_help=1
use_pending_states=1
nagios_check_command=/usr/lib/nagios/plugins/check_nagios /var/cache/nagios3/status.dat 5 '/usr/sbin/nagios3'
use_authentication=1
use_ssl_authentication=0 
authorized_for_system_information=deimos
authorized_for_configuration_information=deimos
authorized_for_system_commands=deimos
authorized_for_all_services=deimos
authorized_for_all_hosts=deimos
authorized_for_all_service_commands=deimos
authorized_for_all_host_commands=deimos
default_statusmap_layout=5
default_statuswrl_layout=4
ping_syntax=/bin/ping -n -U -c 5 $HOSTADDRESS$
refresh_rate=90
escape_html_tags=1
action_url_target=_blank
notes_url_target=_blank
lock_author_names=1
```

### commands.cgi

This file allows you to define special commands for specific checks. For example, if you have developed plugins, you will need to create an associated command:

```bash
###############################################################################
# COMMANDS.CFG - SAMPLE COMMAND DEFINITIONS FOR NAGIOS 
###############################################################################

################################################################################
# NOTIFICATION COMMANDS
################################################################################

# 'notify-host-by-email' command definition
define command{
	command_name	notify-host-by-email
	command_line	/usr/bin/printf "%b" "***** Nagios *****\n\nNotification Type: $NOTIFICATIONTYPE$\nHost: $HOSTNAME$\nState: $HOSTSTATE$\nAddress: $HOSTADDRESS$\nInfo: $HOSTOUTPUT$\n\nDate/Time: $LONGDATETIME$\n" | /usr/bin/mail -s "** $NOTIFICATIONTYPE$ Host Alert: $HOSTNAME$ is $HOSTSTATE$ **" $CONTACTEMAIL$
	}

# 'notify-service-by-email' command definition
define command{
	command_name	notify-service-by-email
	command_line	/usr/bin/printf "%b" "***** Nagios *****\n\nNotification Type: $NOTIFICATIONTYPE$\n\nService: $SERVICEDESC$\nHost: $HOSTALIAS$\nAddress: $HOSTADDRESS$\nState: $SERVICESTATE$\n\nDate/Time: $LONGDATETIME$\n\nAdditional Info:\n\n$SERVICEOUTPUT$" | /usr/bin/mail -s "** $NOTIFICATIONTYPE$ Service Alert: $HOSTALIAS$/$SERVICEDESC$ is $SERVICESTATE$ **" $CONTACTEMAIL$
	}

################################################################################
# HOST CHECK COMMANDS
################################################################################

# On Debian, check-host-alive is being defined from within the
# nagios-plugins-basic package

################################################################################
# PERFORMANCE DATA COMMANDS
################################################################################

# 'process-host-perfdata' command definition
define command{
	command_name	process-host-perfdata
	command_line	/usr/bin/printf "%b" "$LASTHOSTCHECK$\t$HOSTNAME$\t$HOSTSTATE$\t$HOSTATTEMPT$\t$HOSTSTATETYPE$\t$HOSTEXECUTIONTIME$\t$HOSTOUTPUT$\t$HOSTPERFDATA$\n" >> /var/lib/nagios3/host-perfdata.out
	}

# 'process-service-perfdata' command definition
define command{
	command_name	process-service-perfdata
	command_line	/usr/bin/printf "%b" "$LASTSERVICECHECK$\t$HOSTNAME$\t$SERVICEDESC$\t$SERVICESTATE$\t$SERVICEATTEMPT$\t$SERVICESTATETYPE$\t$SERVICEEXECUTIONTIME$\t$SERVICELATENCY$\t$SERVICEOUTPUT$\t$SERVICEPERFDATA$\n" >> /var/lib/nagios3/service-perfdata.out
	}
```

### htpasswd.users

Here is the htpasswd for Nagios. To change the default nagiosadmin login, do:

```bash
htpasswd -c htpasswd.users deimos
```

Replace "deimos" with the user you want. Be careful though, this will create a blank file and create this single user.

### nagios.cfg

This is the Nagios server configuration, but not the configuration of hosts and services:

```bash
log_file=/var/log/nagios3/nagios.log
cfg_file=/etc/nagios3/commands.cfg
cfg_dir=/etc/nagios-plugins/config
cfg_dir=/etc/nagios3/conf.d
object_cache_file=/var/cache/nagios3/objects.cache
precached_object_file=/var/lib/nagios3/objects.precache
resource_file=/etc/nagios3/resource.cfg
status_file=/var/cache/nagios3/status.dat
status_update_interval=10
nagios_user=nagios
nagios_group=nagios
check_external_commands=1
command_check_interval=-1
command_file=/var/lib/nagios3/rw/nagios.cmd
external_command_buffer_slots=4096
lock_file=/var/run/nagios3/nagios3.pid
temp_file=/var/cache/nagios3/nagios.tmp
temp_path=/tmp
event_broker_options=-1
log_rotation_method=d
log_archive_path=/var/log/nagios3/archives
use_syslog=1
log_notifications=1
log_service_retries=1
log_host_retries=1
log_event_handlers=1
log_initial_states=0
log_external_commands=1
log_passive_checks=1
service_inter_check_delay_method=s
max_service_check_spread=30
service_interleave_factor=s
host_inter_check_delay_method=s
max_host_check_spread=30
max_concurrent_checks=0
check_result_reaper_frequency=10
max_check_result_reaper_time=30
check_result_path=/var/lib/nagios3/spool/checkresults
max_check_result_file_age=3600
cached_host_check_horizon=15
cached_service_check_horizon=15
enable_predictive_host_dependency_checks=1
enable_predictive_service_dependency_checks=1
soft_state_dependencies=0
auto_reschedule_checks=0
auto_rescheduling_interval=30
auto_rescheduling_window=180
sleep_time=0.25
service_check_timeout=60
host_check_timeout=30
event_handler_timeout=30
notification_timeout=30
ocsp_timeout=5
perfdata_timeout=5
retain_state_information=1
state_retention_file=/var/lib/nagios3/retention.dat
retention_update_interval=60
use_retained_program_state=1
use_retained_scheduling_info=1
retained_host_attribute_mask=0
retained_service_attribute_mask=0
retained_process_host_attribute_mask=0
retained_process_service_attribute_mask=0
retained_contact_host_attribute_mask=0
retained_contact_service_attribute_mask=0
interval_length=60
check_for_updates=1
bare_update_check=0
use_aggressive_host_checking=0
execute_service_checks=1
accept_passive_service_checks=1
execute_host_checks=1
accept_passive_host_checks=1
enable_notifications=1
enable_event_handlers=1
process_performance_data=0
obsess_over_services=0
obsess_over_hosts=0
translate_passive_host_checks=0
passive_host_checks_are_soft=0
check_for_orphaned_services=1
check_for_orphaned_hosts=1
check_service_freshness=1
service_freshness_check_interval=60
check_host_freshness=0
host_freshness_check_interval=60
additional_freshness_latency=15
enable_flap_detection=1
low_service_flap_threshold=5.0
high_service_flap_threshold=20.0
low_host_flap_threshold=5.0
high_host_flap_threshold=20.0
date_format=iso8601
p1_file=/usr/lib/nagios3/p1.pl
enable_embedded_perl=1
use_embedded_perl_implicitly=1
illegal_object_name_chars=`~!$%^&*|'"<>?,()=
illegal_macro_output_chars=`~$&|'"<>
use_regexp_matching=0
use_true_regexp_matching=0
admin_email=root@localhost
admin_pager=pageroot@localhost
daemon_dumps_core=0
use_large_installation_tweaks=0
enable_environment_macros=1
debug_level=0
debug_verbosity=1
debug_file=/var/log/nagios3/nagios.debug
max_debug_file_size=1000000
```

### resource.cfg

```bash
...
# Sets $USER1$ to be the path to the plugins
$USER1$=/usr/lib/nagios/plugins
...
```

### conf.d/contacts.cfg

Contacts are used to define contacts and contact groups:

```text
define contact{
        contact_name                    deimos
        alias                           Deimos
        service_notification_period     24x7
        host_notification_period        24x7
        service_notification_options    w,u,c,r
        host_notification_options       d,r
        service_notification_commands   notify-service-by-email
        host_notification_commands      notify-host-by-email
        email                           xxx@mycompany.com
        }
define contactgroup{
        contactgroup_name       admins
        alias                   Nagios Administrators
        members                 deimos
        }
```

### conf.d/custom-commands.cfg

I mentioned custom commands earlier. Here's one I created to test MySQL replication. This will allow me to launch replication via the check_mysqlrep command:

```text
##############################
#   CUSTOM NAGIOS COMMANDS   #
##############################

define command{
    command_name	check_mysqlrep
    command_line	/usr/lib/nagios/plugins/check_mysql -H $HOSTADDRESS$ -u $ARG1$ -p $ARG2$ -w $ARG3$ -c $ARG4$ -S
}
```

### conf.d/generic-host.cfg

This file is used to provide the generic configuration for hosts, notification intervals, etc... I'll let you look at the official documentation:

```text
# Generic host definition template - This is NOT a real host, just a template!

define host{
        name                            generic-host    ; The name of this host template
        notifications_enabled           1       ; Host notifications are enabled
        event_handler_enabled           1       ; Host event handler is enabled
        flap_detection_enabled          1       ; Flap detection is enabled
        failure_prediction_enabled      1       ; Failure prediction is enabled
        process_perf_data               1       ; Process performance data
        retain_status_information       1       ; Retain status information across program restarts
        retain_nonstatus_information    1       ; Retain non-status information across program restarts
	check_command                   check-host-alive
	max_check_attempts              10
	notification_interval           0
	notification_period             24x7
	notification_options            d,u,r
	contact_groups                  admins
        register                        0       ; DONT REGISTER THIS DEFINITION - ITS NOT A REAL HOST, JUST A TEMPLATE!
        _PROCWARN                       150
        _PROCCRIT                       200
        }
```

### conf.d/generic_service.cfg

This file is used to provide the generic configuration for services, notification intervals, etc... I'll let you look at the official documentation:

```text
define service{
        name                            generic-service ; The 'name' of this service template
        active_checks_enabled           1       ; Active service checks are enabled
        passive_checks_enabled          1       ; Passive service checks are enabled/accepted
        parallelize_check               1       ; Active service checks should be parallelized (disabling this can lead to major performance problems)
        obsess_over_service             1       ; We should obsess over this service (if necessary)
        check_freshness                 0       ; Default is to NOT check service 'freshness'
        notifications_enabled           1       ; Service notifications are enabled
        event_handler_enabled           1       ; Service event handler is enabled
        flap_detection_enabled          1       ; Flap detection is enabled
        failure_prediction_enabled      1       ; Failure prediction is enabled
        process_perf_data               1       ; Process performance data
        retain_status_information       1       ; Retain status information across program restarts
        retain_nonstatus_information    1       ; Retain non-status information across program restarts
		notification_interval           0		; Only send notifications on status change by default.
		is_volatile                     0
		check_period                    24x7
		normal_check_interval           3
		retry_check_interval            1
		max_check_attempts              3
		notification_period             24x7
		notification_options            w,u,c,r
		contact_groups                  admins
        register                        0       ; DONT REGISTER THIS DEFINITION - ITS NOT A REAL SERVICE, JUST A TEMPLATE!
        }
```

### conf.d/timeperiods.cfg

This file is used to define templates for notification periods:

```text
###############################################################################
# timeperiods.cfg
###############################################################################

# This defines a timeperiod where all times are valid for checks, 
# notifications, etc.  The classic "24x7" support nightmare.  :-)

define timeperiod{
        timeperiod_name 24x7
        alias           24 Hours A Day, 7 Days A Week
        sunday          00:00-24:00
        monday          00:00-24:00
        tuesday         00:00-24:00
        wednesday       00:00-24:00
        thursday        00:00-24:00
        friday          00:00-24:00
        saturday        00:00-24:00
        }

# Here is a slightly friendlier period during work hours
define timeperiod{
        timeperiod_name workhours
        alias           Standard Work Hours
        monday          09:00-17:00
        tuesday         09:00-17:00
        wednesday       09:00-17:00
        thursday        09:00-17:00
        friday          09:00-17:00
        }

# The complement of workhours
define timeperiod{
        timeperiod_name nonworkhours
        alias           Non-Work Hours
        sunday          00:00-24:00
        monday          00:00-09:00,17:00-24:00
        tuesday         00:00-09:00,17:00-24:00
        wednesday       00:00-09:00,17:00-24:00
        thursday        00:00-09:00,17:00-24:00
        friday          00:00-09:00,17:00-24:00
        saturday        00:00-24:00
        }

# This one is a favorite: never :)
define timeperiod{
        timeperiod_name never
        alias           Never
        }

# end of file
```

### conf.d/hostgroups/unix-srv.cfg

Here is an example configuration for a hostgroup. You will then simply associate a host with this hostgroup to inherit all the services described below:

```text
# Some generic hostgroup definitions

define hostgroup {
    hostgroup_name  unix-srv
    alias           Unix servers
}

# Define a service to check the disk space of the root partition
# on the local machine.  Warning if < 20% free, critical if
# < 10% free space on partition.

define service{
    use                             generic-service
    hostgroup_name                  unix-srv
    service_description             Disk Space
    check_command                   check_nrpe!check_all_disks!20%!10%
}

# Define a service to check the number of currently logged in
# users on the local machine.  Warning if > 20 users, critical
# if > 50 users.

define service{
    use                             generic-service         ; Name of service template to use
    hostgroup_name                  unix-srv
    service_description             Current Users
    check_command                   check_nrpe!check_users!2!3
}

# Define a service to check the number of currently running procs
# on the local machine.  Warning if > 250 processes, critical if
# > 400 processes.

define service{
    use                             generic-service         ; Name of service template to use
    hostgroup_name                  unix-srv
    service_description             Total Processes
    check_command                   check_nrpe!check_procs!$_HOSTPROCWARN$!$_HOSTPROCCRIT$
}

# Check Zombie process
define service{
    use                             generic-service         ; Name of service template to use
    hostgroup_name                  unix-srv
    service_description             Zombie Processes
    check_command                   check_nrpe!check_zombie_procs!2!3
}
# Define a service to check the load on the local machine. 

define service{
    use                             generic-service         ; Name of service template to use
    hostgroup_name                  unix-srv
    service_description             Current Load
    check_command                   check_nrpe!check_load!5.0,4.0,3.0!10.0,6.0,4.0
}

# check that ssh services are running
define service{
    use                             generic-service
    hostgroup_name                  unix-srv
    service_description             SSH Servers
    check_command                   check_ssh
}
```

### conf.d/hosts/serveur.cfg

And here I declare a host and associate unix-srv declared above so it inherits the services above:

```text
define host{
    use                     generic-host
    host_name               server.deimos.fr
    alias                   server
    address                 server.deimos.fr
    hostgroups              unix-srv
    _PROCWARN               280
    _PROCCRIT               350
}
```

Here I use variables (_PROCWARN and _PROCCRIT) that override the default values ([See the documentation for more information](https://nagios.sourceforge.net/docs/3_0/customobjectvars.html)). You must add as in the documentation HOST ($_HOSTPROCWARN$ and $_HOSTPROCCRIT$) only for the command declaration part.

If you want a more complete configuration, I've attached an archive with a more complete Nagios3 configuration: [Nagios configuration](/others/nagios3.tgz)

## Addons

### NRPE

We need to configure the NRPE check to handle multiple arguments. Otherwise, we will be limited to just one. Edit the following file and add the necessary number of arguments:

```text
# this command runs a program $ARG1$ with arguments $ARG2$
define command {
    command_name    check_nrpe
    command_line    /usr/lib/nagios/plugins/check_nrpe -H $HOSTADDRESS$ -c $ARG1$ -a $ARG2$ $ARG3$ $ARG4$ $ARG5$ $ARG6$ $ARG7$ $ARG8$ $ARG9$
}

# this command runs a program $ARG1$ with no arguments
define command {
    command_name    check_nrpe_1arg
    command_line    /usr/lib/nagios/plugins/check_nrpe -H $HOSTADDRESS$ -c $ARG1$
}
```

Here I've gone up to 9, but if I remember correctly, you can go up to 32 arguments.

### Hiding "non-OK" alerts

There is a simple solution to avoid displaying certain alerts that are somewhat annoying, such as if like me you have nearly 1500 alerts, a 107cm screen just for Nagios, and some alerts take up too much space and cannot be resolved quickly (e.g., a weekly backup check that failed).

The solution is done through the graphical interface by acknowledging the services you no longer want to see. This way, they will be hidden and when they return to OK, they will be displayed again and the acknowledgement will disappear.

Then in your browser URL, you'll need to modify it a bit to ask it not to display acknowledged alerts. By looking through the CGI sources, you can find this kind of information:

```bash
grep SERVICE_ include/cgiutils.h.in

#define SERVICE_SCHEDULED_DOWNTIME	1
#define SERVICE_NO_SCHEDULED_DOWNTIME	2
#define SERVICE_STATE_ACKNOWLEDGED	4
#define SERVICE_STATE_UNACKNOWLEDGED	8
#define SERVICE_CHECKS_DISABLED		16
#define SERVICE_CHECKS_ENABLED		32
#define SERVICE_EVENT_HANDLER_DISABLED	64
#define SERVICE_EVENT_HANDLER_ENABLED	128
#define SERVICE_FLAP_DETECTION_ENABLED	256
#define SERVICE_FLAP_DETECTION_DISABLED	512
#define SERVICE_IS_FLAPPING		1024
#define SERVICE_IS_NOT_FLAPPING		2048
#define SERVICE_NOTIFICATIONS_DISABLED	4096
#define SERVICE_NOTIFICATIONS_ENABLED	8192
#define SERVICE_PASSIVE_CHECKS_DISABLED	16384
#define SERVICE_PASSIVE_CHECKS_ENABLED	32768
#define SERVICE_PASSIVE_CHECK           65536
#define SERVICE_ACTIVE_CHECK            131072
#define SERVICE_HARD_STATE		262144
#define SERVICE_SOFT_STATE		524288
```

PS: I took these lines from Nagios 3 sources; there are a few less things for versions 2 and 1.

Here is an example URL with the correspondences:

```
http://nagioshost/cgi-bin/nagios3/status.cgi?host=all&servicestatustypes=28&serviceprops=8
```

* servicestatustypes=28: all states except OK
* serviceprops=8: Removes acknowledged states

And a little bonus now, if you want to hide the information summary at the top of the page (hide the header), here's the '&noheader' option:

```
http://nagioshost/cgi-bin/nagios3/status.cgi?host=all&servicestatustypes=28&serviceprops=8&noheader
```

### Adding a custom CGI

In some cases, you may have certain checks that temporarily store information on the Nagios server and you want to be able to execute actions from the Nagios interface. For this, there's the 'action_url' option where we can give a URL to a CGI that will execute what we want, perhaps with options.

To start, we'll create our CGI. Here's a minimalist example where I delete a temporary file:

```perl
#!/usr/bin/perl 

use CGI;
$query = CGI::new();
$host = $query->param("host");

# Avoid inputing special characters that would crash the program
if ( $h =~ /\`|\~|\@|\#|\$|\%|\^|\&|\*|\(|\)|\:|\=|\+|\"|\'|\;|\<|\>/ ) { 
    print "Illegal special chars detected. Exit\n";
    exit(1);
}

print "Content-type: text/html\n\n";
print "<HTML>\n";
print "<HEAD><Title>Removing $host temporary file</Title>\n";
print "<LINK REL='stylesheet' TYPE='text/css' HREF='/nagios/stylesheets/common.css'><LINK REL='stylesheet' TYPE='text/css' HREF='/nagios/stylesheets/status.css'>\n";
print "</HEAD><BODY>\n";
print "Removing $host Interface Network Flapping temporary file...";
if (-f "/tmp/iface_state_$host.txt")
{
    unlink("/tmp/iface_state_$host.txt") or print "FAIL<br />/tmp/iface_state_$host.txt: $!\n" and exit(1);
    print "OK\n";
}
else
{
    print "FAIL<br />/tmp/iface_state_$host.txt: No such file or directory\n";
}
print "</body></html>\n";
```

And then in the configuration of the service in question, I insert my 'action_url':

```text
define service{
         use                             generic-services-ulsysnet
         hostgroup_name                  network
         service_description             Interface Network Flapping
         check_period                    24x7
         notification_period             24x7
 	 _SNMP_PORT			 161
	 _SNMP_COMMUNITY		         public
         _DURATION			 86400
         check_command                   check_interface_flapping
         # For Thruk & Nagios
	 # action_url			 ../../cgi-bin/nagios3/remove.cgi?host=$HOSTADDRESS$
         # For Nagios only               
         action_url			 remove.cgi?host=$HOSTADDRESS$
}
```

All you have to do now is reload Nagios.

### Sending SMS alerts via Free Mobile

Thanks to Free Mobile for offering the possibility to send SMS via an API (which can be activated from the web interface of your account). How does it work? Simply create a scripts folder in the configuration folder and put the content of this script in it:

```bash {linenos=table,hl_lines=[3]}
#!/bin/bash
message="$(perl -MURI::Escape -e 'print uri_escape($ARGV[0]);' "Naemon $1: $2 on $3 $4 ($5)")"
curl --insecure "https://smsapi.free-mobile.fr/sendmsg?user=<userid>&pass=<password>&msg=$message"
```

Adapt the 'userid' and 'password' fields with your personal settings (these credentials are available from the Free Mobile web interface). Now create the associated contacts and contactgroups:

```text
define contact {
  contact_name                   oncall-sms                        ; Short name of user
  alias                          oncall-sms                       ; Full name of user
  use                            contact-sms                    ; Inherit default values from generic-contact template (defined above)
}

define contactgroup {
  contactgroup_name              admins-sms
  alias                          Nagios Administrators SMS
  members                        oncall-sms
}
```

Add the necessary commands that can call the script:

```text
# 'notify-host-by-sms' command definition
define command {
  command_name                   notify-host-by-sms
  command_line                   /etc/naemon/scripts/send_sms.sh $SHORTDATETIME$ $NOTIFICATIONTYPE$ $HOSTNAME$/$SERVICEDESC$ $HOSTSTATE$ $HOSTOUTPUT$
}

# 'notify-service-by-sms' command definition
define command {
  command_name                   notify-service-by-sms
  command_line                   /etc/naemon/scripts/send_sms.sh $SHORTDATETIME$ $NOTIFICATIONTYPE$ $HOSTNAME$/$SERVICEDESC$ $SERVICESTATE$ $SERVICEOUTPUT$
}
```

```text
define contact {
  name                           contact-sms                    ; The name of this contact template
  host_notification_commands     notify-host-by-sms               ; send host notifications via email
  host_notification_options      d,u,r                          ; send notifications for all host states, flapping events, and scheduled downtime events
  host_notification_period       24x7                               ; host notifications can be sent anytime
  register                       0                                  ; DONT REGISTER THIS DEFINITION - ITS NOT A REAL CONTACT, JUST A TEMPLATE!
  service_notification_commands  notify-service-by-sms            ; send service notifications via email
  service_notification_options   u,c,r                        ; send notifications for all service states, flapping events, and scheduled downtime events
  service_notification_period    24x7                               ; service notifications can be sent anytime
}
```

And finally the service escalation, to notify everyone by SMS if there hasn't been action taken quickly enough:

```text
define serviceescalation {
	host_name              *
	service_description    *
	first_notification     3
	last_notification      4
	notification_interval  0
	contact_groups         admins,admins-sms
}
```

Of course, this is just an example, and you will certainly need to adapt it to your needs.

## FAQ

### No output returned from plugin

#### Solution 1

If you encounter this type of error, it's because you haven't activated the arguments in the NRPE configuration. In the NRPE configuration file `/etc/nagios/nrpe.cfg`, change this value:

```
dont_blame_nrpe=0
```

to 1:

```
dont_blame_nrpe=1
```

Then restart NRPE and you're good to go :-)

#### Solution 2

I had this problem and searched for a while before finding the solution. It comes from NRPE and for the host's configuration in question, replace something like this:

```
check_nrpe!check_disk_c
```

with

```
check_nrpe_1arg!check_disk_c
```

The difference is that:
* check_nrpe: takes arguments
* check_nrpe_1arg: takes no arguments

### Error: Could not stat() command file '/var/lib/nagios3/rw/nagios.cmd'!

#### Solution 1

This kind of error exists simply because you don't have the permissions, so to solve this problem:

```bash
chmod -Rf a+rx /var/lib/nagios3
```

Or:

```bash
chown nagios.www-data /var/lib/nagios3/rw/nagios.cmd
```

#### Solution 2

Usually a permissions issue on Debian, here's the solution! Execute these commands as root:

```bash
dpkg-statoverride --update --add nagios www-data 2710 /var/lib/nagios3/rw
dpkg-statoverride --update --add nagios nagios 751 /var/lib/nagios3
chmod -Rf a+rx /var/lib/nagios3
chown nagios.www-data /var/lib/nagios3/rw/nagios.cmd
```

### I am not receiving emails

You need to check if the mail command is present in command.cfg and especially with the correct PATH. For example with Debian, I encountered a small problem and this simple symbolic link solved the problem:

```bash
ln -s /usr/bin/mail /bin/mail
```

### return code of 127 is out of bounds

This is probably due to a permission error (plugin execution) or because you are not indicating the complete path of the executable you want to launch. Forget $USER1$, it doesn't work very well for me. After a few Nagios reloads, it starts to malfunction, so use the full path.

### CHECK_NRPE: Error - Could not complete SSL handshake.

I bet it's because you don't have the permissions! A small telnet on port 5666 of your server will tell you a lot. Then modify the allowed_hosts line of nrpe.cfg and everything should be back in order :-)

### Monitoring Wordpress

Wordpress has a small peculiarity with check_http, which is that you need to make it follow links (option "follow"). Here's an example check:

```text
define command{
    command_name    check_blog
    command_line    /usr/lib/nagios/plugins/check_http -H 'www.deimos.fr' -u '/blog/index.php' -s 'Parce que la mémoire humaine ne fait pas des Go' -f follow
}
```

### You don't have permission to access /nagios3/ on this server

If like me after a Debian update (5 -> 6) you have this message, you need to replace "index.html" with "index.php":

```apache
 ...
 DirectoryIndex index.php
 ...
```

Then restart the service:

```bash
/etc/init.d/apache2 restart
```

### My pending checks are still there even after a reboot

It's possible to have issues on your Nagios machine and the checks that need to be done remain active, even after a reboot. This is simply because Nagios keeps them in memory to be able to re-execute them later. To purge this queue, modify this in the Nagios configuration, restart it, then change the parameter back to 1:

```text
[...]
retain_state_information=0
[...]
```

## Resources
- [Nagios Documentation on OpenBSD](/pdf/openbsd_nagios.pdf)
- [https://www.mail-archive.com/nagios-users@lists.sourceforge.net/msg04394.html](https://www.mail-archive.com/nagios-users@lists.sourceforge.net/msg04394.html)
- [Nagios Documentation](/pdf/nagios.pdf)
