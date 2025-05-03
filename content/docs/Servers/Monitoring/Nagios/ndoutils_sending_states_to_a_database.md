---
weight: 999
url: "/NDOUtils_\\:_Envoyer_les_états_en_base_de_donnée/"
title: "NDOUtils: Sending States to a Database"
description: "Guide on configuring NDOUtils to send Nagios monitoring states to a MySQL or PostgreSQL database"
categories: ["Database", "MySQL", "PostgreSQL", "Monitoring"]
date: "2010-07-07T08:58:00+02:00"
lastmod: "2010-07-07T08:58:00+02:00"
tags: ["Nagios", "NDOUtils", "MySQL", "PostgreSQL", "Monitoring", "Database"]
toc: true
---

## Introduction

NDO is an additional module allowing Nagios to write the state of the machines and services to be monitored in a database.

NDO is made up of two modules: NDOMOD and NDO2DB.

NDOMOD must be launched on the Nagios server and retrieves the information reported by Nagios to transmit it via TCP (or a Unix socket) to NDO2DB.

NDO2DB is a daemon that listens on a TCP port (or a Unix socket) and writes the received data to a database (MySQL or PgSQL).

## Installation

We'll assume that you already have Nagios installed and configured. We need to install NDOUtils:

```bash
aptitude install ndoutils ndoutils-mysql
```

With this, we'll be able to store data in a MySQL database.

## Configuration

### Default

Edit the `/etc/default/ndoutils` file and check that this value is set to 1:

```bash
...
ENABLE_NDOUTILS=1
..
```

### ndo2db.cfg

Configure the database for ndo2db. Normally all database information was pre-filled during NDOUtils installation. However, we'll make a few small adjustments:

```bash
#####################################################################
# NDO2DB DAEMON CONFIG FILE
#
# Last Modified: 10-29-2007
#####################################################################

# USER/GROUP PRIVILIGES
# These options determine the user/group that the daemon should run as.
# You can specify a number (uid/gid) or a name for either option.

ndo2db_user=nagios
ndo2db_group=nagios

# SOCKET TYPE
# This option determines what type of socket the daemon will create
# an accept connections from.
# Value:
#   unix = Unix domain socket (default)
#   tcp  = TCP socket

#socket_type=unix
socket_type=tcp

# SOCKET NAME
# This option determines the name and path of the UNIX domain 
# socket that the daemon will create and accept connections from.
# This option is only valid if the socket type specified above
# is "unix".

socket_name=/var/cache/nagios3/ndo.sock

# TCP PORT
# This option determines what port the daemon will listen for
# connections on.  This option is only vlaid if the socket type
# specified above is "tcp".

tcp_port=5668

# DATABASE SERVER TYPE
# This option determines what type of DB server the daemon should
# connect to.
# Values:
# 	mysql = MySQL
#       pgsql = PostgreSQL

db_servertype=mysql

# DATABASE HOST
# This option specifies what host the DB server is running on.

db_host=localhost

# DATABASE PORT
# This option specifies the port that the DB server is running on.
# Values:
# 	3306 = Default MySQL port
#	5432 = Default PostgreSQL port

db_port=3306

# DATABASE NAME
# This option specifies the name of the database that should be used.

db_name=ndoutils

# DATABASE TABLE PREFIX
# Determines the prefix (if any) that should be prepended to table names.
# If you modify the table prefix, you'll need to modify the SQL script for
# creating the database!

db_prefix=nagios_

# DATABASE USERNAME/PASSWORD
# This is the username/password that will be used to authenticate to the DB.
# The user needs at least SELECT, INSERT, UPDATE, and DELETE privileges on
# the database.

db_user=ndoutils
db_pass=password

## TABLE TRIMMING OPTIONS
# Several database tables containing Nagios event data can become quite large
# over time.  Most admins will want to trim these tables and keep only a
# certain amount of data in them.  The options below are used to specify the
# age (in MINUTES) that data should be allowd to remain in various tables
# before it is deleted.  Using a value of zero (0) for any value means that
# that particular table should NOT be automatically trimmed.

# Keep timed events for 24 hours
max_timedevents_age=1440

# Keep system commands for 1 week
max_systemcommands_age=10080

# Keep service checks for 1 week
max_servicechecks_age=10080

# Keep host checks for 1 week
max_hostchecks_age=10080

# Keep event handlers for 31 days
max_eventhandlers_age=44640

# DEBUG LEVEL
# This option determines how much (if any) debugging information will
# be written to the debug file.  OR values together to log multiple
# types of information.
# Values: -1 = Everything
#          0 = Nothing
#          1 = Process info
#	   2 = SQL queries

debug_level=0

# DEBUG VERBOSITY
# This option determines how verbose the debug log out will be.
# Values: 0 = Brief output
#         1 = More detailed
#         2 = Very detailed

debug_verbosity=0

# DEBUG FILE
# This option determines where the daemon should write debugging information.

debug_file=@localstatedir@/ndo2db.debug

# MAX DEBUG FILE SIZE
# This option determines the maximum size (in bytes) of the debug file.  If
# the file grows larger than this size, it will be renamed with a .old
# extension.  If a file already exists with a .old extension it will
# automatically be deleted.  This helps ensure your disk space usage doesn't
# get out of control when debugging.

max_debug_file_size=1000000
```

### ndomod.cfg

Then for NDOMOD:

```bash
#####################################################################
# NDOMOD CONFIG FILE
#
# Last Modified: 09-05-2007
#####################################################################


# INSTANCE NAME
# This option identifies the "name" associated with this particular
# instance of Nagios and is used to seperate data coming from multiple
# instances.  Defaults to 'default' (without quotes).

instance_name=default



# OUTPUT TYPE
# This option determines what type of output sink the NDO NEB module
# should use for data output.  Valid options include:
#   file       = standard text file
#   tcpsocket  = TCP socket
#   unixsocket = UNIX domain socket (default)

#output_type=file
output_type=tcpsocket
#output_type=unixsocket



# OUTPUT
# This option determines the name and path of the file or UNIX domain 
# socket to which output will be sent if the output type option specified
# above is "file" or "unixsocket", respectively.  If the output type
# option is "tcpsocket", this option is used to specify the IP address
# of fully qualified domain name of the host that the module should
# connect to for sending output.

#output=/var/cache/nagios2/ndo.dat
output=127.0.0.1
#output=/var/cache/nagios3/ndo.sock



# TCP PORT
# This option determines what port the module will connect to in
# order to send output.  This option is only vlaid if the output type
# option specified above is "tcpsocket".

tcp_port=5668



# OUTPUT BUFFER
# This option determines the size of the output buffer, which will help
# prevent data from getting lost if there is a temporary disconnect from
# the data sink.  The number of items specified here is the number of
# lines (each of variable size) of output that will be buffered.

output_buffer_items=5000



# BUFFER FILE
# This option is used to specify a file which will be used to store the
# contents of buffered data which could not be sent to the NDO2DB daemon
# before Nagios shuts down.  Prior to shutting down, the NDO NEB module
# will write all buffered data to this file for later processing.  When
# Nagios (re)starts, the NDO NEB module will read the contents of this
# file and send it to the NDO2DB daemon for processing.

buffer_file=/var/cache/nagios3/ndomod.tmp



# FILE ROTATION INTERVAL
# This option determines how often (in seconds) the output file is
# rotated by Nagios.  File rotation is handled by Nagios by executing
# the command defined by the file_rotation_command option.  This
# option has no effect if the output_type option is a socket.

file_rotation_interval=14400



# FILE ROTATION COMMAND
# This option specified the command (as defined in Nagios) that is
# used to rotate the output file at the interval specified by the
# file_rotation_interval option.  This option has no effect if the
# output_type option is a socket.
#
# See the file 'misccommands.cfg' for an example command definition
# that you can use to rotate the log file.

#file_rotation_command=rotate_ndo_log



# FILE ROTATION TIMEOUT
# This option specified the maximum number of seconds that the file
# rotation command should be allowed to run before being prematurely
# terminated.

file_rotation_timeout=60



# RECONNECT INTERVAL
# This option determines how often (in seconds) that the NDO NEB
# module will attempt to re-connect to the output file or socket if
# a connection to it is lost.

reconnect_interval=15



# RECONNECT WARNING INTERVAL
# This option determines how often (in seconds) a warning message will
# be logged to the Nagios log file if a connection to the output file
# or socket cannot be re-established.

reconnect_warning_interval=15
#reconnect_warning_interval=900



# DATA PROCESSING OPTION
# This option determines what data the NDO NEB module will process. 
# Do not mess with this option unless you know what you're doing!!!!
# Read the source code (include/ndbxtmod.h) to determine what values
# to use here.  Values from source code should be OR'ed to get the
# value to use here.  A value of -1 will cause all data to be processed.
# Read the source code (include/ndomod.h) and look for "NDOMOD_PROCESS_"
# to determine what values to use here.  Values from source code should
# be OR'ed to get the value to use here.  A value of -1 will cause all
# data to be processed. 

data_processing_options=-1



# CONFIG OUTPUT OPTION
# This option determines what types of configuration data the NDO
# NEB module will dump from Nagios.  Values can be OR'ed together.
# Values: 
# 	  0 = Don't dump any configuration information
#         1 = Dump only original config (from config files)
#         2 = Dump config only after retained information has been restored
#         3 = Dump both original and retained configuration

config_output_options=2
```

### nagios.cfg

Finally, we'll edit the nagios.cfg file:

```bash
event_broker_options=-1
broker_module=/usr/lib/ndoutils/ndomod-mysql-3x.o config_file=/etc/nagios3/ndomod.cfg
```

## Restart

To finish, you need to restart the ndoutils service. Personally, I've had too many annoying errors in the syslog like:

```
nagios3: ndomod: Still unable to connect to data sink.  1325744 items lost, 5000 queued items to flush.
```

Actually, after rebooting the machine, everything started working correctly (the database was populated).

## Resources
- http://blog.nicolargo.com/2009/02/pour-en-finir-avec-ndo.html
