---
weight: 999
url: "/Syslog-ng_\\:_Installation_et_configuration_de_Syslog-ng/"
title: "Syslog-ng: Installation and Configuration of Syslog-ng"
description: "A guide to install and configure Syslog-ng, a powerful new generation system log manager with advanced features for centralizing, sorting, and securing log data across networks."
categories: ["Security", "Debian", "Database"]
date: "2011-12-12T20:18:00+02:00"
lastmod: "2011-12-12T20:18:00+02:00"
tags: ["Syslog", "Logging", "MySQL", "Network", "Security", "Servers", "Linux"]
toc: true
---

## Introduction

[Syslog-ng](https://www.balabit.com/network-security/syslog-ng/opensource-logging-system/) is a new generation system log manager. It allows you to centralize logs from machines in a computer network with incredible ease, and to sort them just as easily. This daemon is therefore essential for network administrators concerned about the performance of their machines.

Free and working on many systems such as Linux, FreeBSD, HP-UX, Solaris or AIX, Syslog-Ng effectively replaces the basic syslogd daemon. It thus allows to overcome its many shortcomings:

- Recognized portability
- Advanced basic configuration
- Ability to "chroot" its environment
- Export of received logs to a MySQL server
- Possible use of macros for log names
- Ability to encrypt logs sent via SSL technology
- Use of UDP and TCP protocols when transporting logs
- Ability to sort logs according to contents, origin or facility, with the possibility of using regular expressions

With it, you can easily retrieve logs from a given machine according to your own criteria, with increased security during their transmission thanks to the power of SSL encryption. Allowing the most advanced processing of logs upon reception, it will prevent you from using Bash or Perl scripts for flat file processing, as was the case when using syslogd.

## Installation

### Debian

For installation on Debian:

```bash
apt-get install syslog-ng
```

### From Source

Note: commands starting with # must be executed in root mode.

First, we need to download the archives needed to run Syslog-NG, its own archive as well as the "libol" library.

Downloading syslog-ng and libol packages with Wget:

```bash
wget http://www.balabit.com/downloads/libol/0.3/libol-0.3.16.tar.gz http://www.balabit.com/downloads/syslog-ng/1.6/src/syslog-ng-1.6.8.tar.gz
```

Unpacking:

```bash
$ gzip -dc syslog-ng-1.6.8.tar.gz
```

Compiling unpacked packages:

```bash
$ cd /usr/local/src/libol-0.3.16
$ ./configure
# make && make install

$ cd /usr/local/src/syslog-ng-1.6.8
$ ./configure

# make && make install
```

Your system is now equipped with Syslog-ng. You now need to start syslog-ng at boot time, in background, as well as the syslog_mysql.sh script, which we will study in the following pages. To do this, you simply need to create a corresponding service.

Startup script `/etc/init.d/syslog-ng`:

```bash
#! /bin/sh

# /etc/init.d/syslog-ng

# Description: Syslog-ng is the system that
# will be used by the machine to monitor
# logs from the entire network

NAME=syslog-ng
DAEMON=/usr/local/sbin/syslog-ng
PIDFILE=/var/run/$NAME.pid
case "$1" in

start)
 echo "Starting $NAME..."
 start-stop-daemon --start --pidfile $PIDFILE \
 --exec $DAEMON
 PID=`more /var/run/syslog-ng.pid`
 echo "MySQL sync system for $NAME..."
 /scripts/syslog_mysql.sh $PID &
 ;;

stop)
 echo "Stopping $NAME"
 start-stop-daemon --stop --pidfile $PIDFILE \
 --oknodo --exec $DAEMON
 rm –f /var/log/mysql.pipe
 ;;

*)
 echo "Usage of syslog-Ng must be start or stop"
 exit 1

 esac
 exit 1
```

You must then link this script to the usual Run Level of your system.

```bash
$ ln -s /etc/init.d/syslog-ng /etc/rc2.d/S99syslog-ng
```

## Configuration

Syslog-ng loads a configuration file at startup, located in `/etc/syslog-ng/syslog-ng.conf`. However, you can use another configuration file using the -f option. The file is divided into 5 major sections:

- Options, to define the general parameters of the log server
- Sources, to define the different possible sources
- Destination, to define where the log will be stored
- Filter, to define the filters operating on logs (content, facility, etc.)
- Log, to define the actions and handling of logs

Logs are processed according to their sources, destinations, and forms (filters). You can indeed ignore a good part of the logs that you receive by setting restrictive filter rules when processing the logs.

### The options section

Here is the list of possible parameters to customize the daemon's operation:

- log_msg_size(): Maximum size of a message in bytes
- sync(): Number of events before writing to logs
- log_fifo_size(): The event processing stack, allows storing x lines in memory

- time_reap(): Closes a log file after x seconds
- time_reopen(): Number of seconds to wait if the connection is down

- create_dirs(): Create log directories if necessary (Yes or No)
- perm() and dir_perm(): Log file and log directory permissions
- group() and dir_group(): Group owner of logs and log directories
- owner() and dir_owner(): User owner of logs and log directories

- use_dns(): Use of DNS servers to resolve names
- long_hostnames(): Use of long DNS names (On or Off)
- check_hostname(): Checks that client DNS name is valid (Yes or No)
- use_fqdn(): Use of FQDNs in log names (Yes or No)
- keep_hostname(): Specifies whether to "trust" the hostname (Yes or No)
- dns_cache(), dns_cache_size(), dns_cache_expire(): Activation of DNS cache for x Hosts during X Seconds

- use_time_recvd(): Use local time instead of timestamp (Yes or No)
- gc_busy_threshold(): Launch of Garbage Collector after X events when syslog-ng is active
- gc_idle_threshold(): Launch of Garbage Collector after X events when syslog-ng is inactive

### The source section

We define the different possible sources of logs. We can define several at once, in order to group several "sources" into a single "virtual" one.

Example to retrieve logs arriving from 10.0.1.5 in TCP:

```bash
source ip780 {
   tcp (ip ("10.0.1.5") );
};
```

- Possible sources:
  - file() = opens a given file and reads it
  - internal() = internal messages from syslog-ng
  - pipe(), fifo() = opens a given fifo file and reads it
  - udp(), tcp() = listens on the specified udp or tcp port
  - sun-stream(), sun-streams() = listens for Solaris systems
  - unix-dgram() = reads UNIX sockets in SOCK_DGRAM mode (BSD)
  - unix-stream() = reads UNIX sockets in SOCK_STREAM mode (Linux)

If you don't want to bother, and accept all tcp and udp requests:

```bash
source s_all {
   ...
   udp();
   tcp();
};
```

### The destination section

We define the different possible destinations for logs. We can define several at once, to send logs to several files for example.

Example of sending the message via a program, here a perl script, and to a file:

```bash
destination mailfile {
   program("perl /scripts/log_mailsend.pl $MSG ");
   file ("/var/log/syslog-ng/$YEAR.$MONTH.$DAY/$HOST/logfile.log");
};
```

Note that filenames support Macros. This is one of the most interesting aspects of syslog-ng: you can easily sort your log files using easily recognizable and customizable names. Thus, macros allow you to name the log file according to the machine, its IP, the date, etc...  
In our example, the file could for instance be named `/var/log/syslog-ng/2005.10.04/MACHINE/logfile.log`.

- Possible destinations:

  - file() = writes to a given file
  - usertty() = sends the log to a given tty
  - fifo(), pipe() = writes to a given fifo file
  - udp(),tcp() = sends the message on the specified port
  - program() = launches a given program and sends it the message
  - unix-dgram() = sends a message containing the log in SOCK_DGRAM
  - unix-stream() = sends a message containing the log in SOCK_STREAM

- Usable Macros:
  - FACILITY: log facility
  - LEVEL or PRIORITY: log level
  - DATE: the date during which the log was sent
  - DAY: the day during which the log was sent
  - HOST: the name of the machine that sent the log
  - YEAR: the year during which the log was sent
  - HOUR: the hour during which the log was sent
  - MIN: the minute during which the log was sent
  - MONTH: the month during which the log was sent
  - SEC: the second during which the log was sent
  - PROGRAM: the name of the program that sent the log
  - FULLDATE: the entire date of the log: with hour, minute and second
  - WEEKDAY: the first three letters of the day during which the log was sent (example: "Wed")

### The filter section

We define filters to limit the logs taken into account. This is a crucial part of the syslog-ng server, as it allows direct processing of logs upon arrival, without having to parse or process them after reception. This is the section you should mainly focus on.

Example to retrieve authentication errors:

```bash
filter auth_errors { level(error) and facility(auth); };
```

Note the ability to invert them with NOT and use the AND and OR operands

- Possible filters:
  - facility() = the facility of the log
  - level(), priority() = the level of the log
  - filter() = evaluates the log according to another filter
  - netmask() = checks the subnet mask
  - match() = message pattern (supports Regex)
  - host() = the machine presenting the log (supports Regex)
  - program() = the program that generated the log (supports Regex)

Caution, Syslog-ng developers recommend using regular expressions as little as possible to avoid overloading the CPU of machines running the daemon. It is therefore better to use several "match" filters, coupled with "AND" operands, rather than using a single "match" filter using Regex. Think carefully before using them!

### The log section

This last part simply allows you to "build" your log captures from the sections we defined earlier. You can therefore completely configure the application according to your needs.

Example of log section:

```bash
log exemple {
   source(ip780);
   destination(mailfile);
   filter(auth_errors);
};
```

This example allows you to retrieve authentication errors (filter auth_errors) coming from machine 10.0.1.5 (source ip780) and send them to the file defined in the destination (mailfile).

### Example configuration file

To better understand the general system of Syslog-ng, here is an example of a commented basic configuration file:

```bash
#
# Configuration file for syslog-ng under Debian
#
# attempts at reproducing default syslog behavior

# the standard syslog levels are (in descending order of priority):
# emerg alert crit err warning notice info debug
# the aliases "error", "panic", and "warn" are deprecated
# the "none" priority found in the original syslogd configuration is
# only used in internal messages created by syslogd


######
# options

options {
        # disable the chained hostname format in logs
        # (default is enabled)
        chain_hostnames(0);

        # the time to wait before a died connection is re-established
        # (default is 60)
        time_reopen(10);

        # the time to wait before an idle destination file is closed
        # (default is 60)
        time_reap(360);

        # the number of lines buffered before written to file
        # you might want to increase this if your disk isn't catching with
        # all the log messages you get or if you want less disk activity
        # (say on a laptop)
        # (default is 0)
        #sync(0);

        # the number of lines fitting in the output queue
        log_fifo_size(2048);

        # enable or disable directory creation for destination files
        create_dirs(yes);

        # default owner, group, and permissions for log files
        # (defaults are 0, 0, 0600)
        #owner(root);
        group(adm);
        perm(0640);

        # default owner, group, and permissions for created directories
        # (defaults are 0, 0, 0700)
        #dir_owner(root);
        #dir_group(root);
        dir_perm(0755);

        # enable or disable DNS usage
        # syslog-ng blocks on DNS queries, so enabling DNS may lead to
        # a Denial of Service attack
        # (default is yes)
        use_dns(yes);

        # maximum length of message in bytes
        # this is only limited by the program listening on the /dev/log Unix
        # socket, glibc can handle arbitrary length log messages, but -- for
        # example -- syslogd accepts only 1024 bytes
        # (default is 2048)
        #log_msg_size(2048);

        #Disable statistic log messages.
        stats_freq(0);
};

# Possible sources

# Local events
source s_localhost
{
    pipe ("/proc/kmsg" log_prefix("kernel: "));
    unix-stream ("/dev/log");
    internal();
};


# All network udp logs
source s_network { udp( port(514) ); };

# Possible destinations
# For all local logs
destination d_localhost {
file ("/var/log/syslogng/$YEAR.$MONTH.$DAY/localhost/$FACILITY.log");
};

# For all network logs
destination d_network {
# Possibility to use mySQL via a fifo file

pipe ("/tmp/mysql.pipe" template ("INSERT INTO logs(host, \
       facility, priority, level, tag, datetime, program,\
       msg) VALUES ('$HOST', '$FACILITY', '$PRIORITY', \
       '$LEVEL', '$TAG', '$YEAR-$MONTH-$DAY $HOUR:$MIN:$SEC', \
       '$PROGRAM',  '$MSG');\n") template-escape(yes) );

file ("/var/log/syslog-ng/$YEAR.$MONTH.$DAY/$HOST/reseau.log");

};


# Filters
filter f_local6 { facility(user); };


# Log processing itself
log {
    source(s_localhost);
    destination(d_localhost);
};


log {
    source(s_network);
    filter(f_local6);
    destination(d_network);
};
```

The server configuration is now complete. But this is obviously not enough, you now need to specify to the target machines that they must transmit their syslog frames to this server. We will now explain how the basic syslogd daemon works, and configure it to properly redirect its frames.

## MySQL

For integration of logs in MySQL, here is an example of destination to add in the configuration file:

```bash
destination d_mysql {
       program(
               "mysql -usyslogng -psyslogng syslogng -B > /dev/null"
               template("INSERT INTO logs (host, facility, priority,
                       level, tag, date, time, program, msg) VALUES ( '$HOST', '$FACILITY', '$PRIORITY', '$LEVEL',
                       '$TAG', '$YEAR-$MONTH-$DAY $HOUR:$MIN:$SEC', '$PROGRAM', '$MSG' );\n")
               template-escape(yes)
       );
};
```

- The name of the destination here is d_mysql.
- -u: for the mysql username.
- -p: the mysql password.
- sysloggn: this is the database name here.
- logs: the name of the table that will receive the logs.

Also add the logs table record in the syslogng database:

```bash
$ mysql -uroot -p
create database syslogng;
use syslogng;
create table logs (
host varchar(32) default NULL,
facility varchar(10) default NULL,
priority varchar(10) default NULL,
level varchar(10) default NULL,
tag varchar(10) default NULL,
datetime time default NULL,
program varchar(15) default NULL,
msg text, seq int(10) unsigned NOT NULL auto_increment,
PRIMARY kEY (seq))
TYPE=MyISAM;
```

Now, let's go back to the configuration file to configure at least one log:

```bash
log {
       source(s_all);
       filter(f_daemon);
       destination(d_mysql);
       destination(df_daemon);
};
```

I've added my destination line here. This means that in addition to writing to a file located in /var/log, I'm going to send the logs to mysql. You may not need both lines. It's up to you.

For proper reading and interpretation of logs, I invite you to follow the article explaining the [setup of php-syslog-ng](https://wiki.deimos.fr/Php-syslog-ng_:_Interpr%C3%A9tation_des_logs_Syslog-ng_dans_une_interfa%C3%A7e_web.html).

To test proper operation, I recommend proceeding as follows:

- Restart the ssh service on the client machine
- Check your log file `/var/log/auth.log` on the server. You should see logs from the remote machine arriving.
- In case you use name resolution and the DNS server does not resolve, insert your machine in `/etc/hosts` (this happens if you don't see your SQL table incrementing)

## Clients

For clients, there are several syslogs and several OSes. We will therefore see how to do with almost everything.

### Linux with Syslog

On Linux and most distributions, Syslog is the default installed log server. That's why we're going to see how to send from syslog to syslog-ng the logs we're interested in. Edit the `/etc/syslog.conf` file:

```bash
auth,authpriv.*                 @192.168.0.193
auth,authpriv.*                 /var/log/auth.log
*.*;auth,authpriv.none          @192.168.0.193
*.*;auth,authpriv.none          -/var/log/syslog
#cron.*                         /var/log/cron.log
daemon.*                        @192.168.0.193
daemon.*                        -/var/log/daemon.log
kern.*                          @192.168.0.193
kern.*                          -/var/log/kern.log
lpr.*                           @192.168.0.193
lpr.*                           -/var/log/lpr.log
mail.*                          @192.168.0.193
mail.*                          -/var/log/mail.log
user.*                          @192.168.0.193
user.*                          -/var/log/user.log
uucp.*                          @192.168.0.193
uucp.*                          /var/log/uucp.log
```

Here my Syslogng server has the IP address: 192.168.0.193. I specify a "@" to indicate redirection followed by the IP. I added lines for the server, because I want logs to be local and on the Syslogng server. If you don't want them locally, delete the lines that are not in bold above. If you don't want to go into detail and send everything, put this line:

```bash
*.*            @ip_server_syslog
```

The "-" sign in front of log files means that you want to record logs asynchronously, which avoids bottlenecks.

Then restart the service:

```bash
/etc/init.d/sysklogd restart
```

### Linux with Syslogng

For a Linux client with a Syslog-ng server, simply define another destination:

```bash
destination dtcp_syslog_server { tcp("192.168.0.193" port(514)); };
```

Here, I chose to send my logs in tcp to server 192.168.0.193 and port 514. However, this is just the destination, I'm not sending anything yet. That's why we need to configure the logs:

```bash
...
log {
       source(s_all);
       filter(f_auth);
       destination(df_auth);
       destination(dtcp_syslog_server);
};

# *.*;auth,authpriv.none          -/var/log/syslog
log {
        source(s_all);
        filter(f_syslog);
        destination(df_syslog);
        destination(dtcp_syslog_server);
};

# daemon.*                        -/var/log/daemon.log
log {
       source(s_all);
       filter(f_daemon);
       destination(df_daemon);
       destination(dtcp_syslog_server);
};

# kern.*                          -/var/log/kern.log
log {
       source(s_all);
       filter(f_kern);
       destination(df_kern);
       destination(dtcp_syslog_server);
};
...
```

As you can see in the example above, you just need to add the recently configured destinations in the log. All that remains is to restart the service and that's it :-)

## FAQ

### I/O error occurred while writing; fd='6', error='Broken pipe (32)'

This is most likely due to a connection problem to the database. Check your permissions and the credentials used in the configuration file.

## Resources
- http://www.supinfo-projects.com/fr/2005/syslogng_2005/1/
- [Install a Syslog Server](/pdf/installer_un_serveur_syslog.pdf)
