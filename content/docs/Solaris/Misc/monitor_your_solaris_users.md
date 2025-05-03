---
weight: 999
url: "/Monitorer_ses_Solaris_Users/"
title: "Monitor Your Solaris Users"
description: "How to monitor Solaris user activities, track logins, and manage access control on Solaris systems."
categories: ["Monitoring", "Security", "Linux"]
date: "2006-12-01T10:35:00+02:00"
lastmod: "2006-12-01T10:35:00+02:00"
tags: ["su", "logins", "rusers", "Servers", "Network", "Security"]
toc: true
---

## Introduction

All systems should be monitored routinely for unauthorized user access. You can determine who is or who has been logged into the system by executing commands and examining log files.

## who

If a user is logged in remotely, the who command displays the remote host name, or Internet Protocol (IP) address in the last column of the output.

```bash
who
```

```
root       console      Oct 17 08:21	(:0)
root       pts/4        Oct 17 08:21	(:0.0)
root       pts/5        Oct 17 08:21	(:0.0)
user5      pts/6        Oct 17 09:20	(sys-03)
root       pts/7        Oct 17 09:20	(:0.0)
user3      pts/8        Oct 17 09:21	(localhost)
```

## rusers

The rusers command produces output similar to that of the who command, but it displays a list of the users logged in on local and remote hosts. The list displays the user's name and the host's name in the order in which the responses are received from the hosts.

A remote host responds only to the rusers command if its rpc.rusersd daemon is enabled. The rpc.rusersd daemon is the network server daemon that returns the list of users on the remote hosts.

Note: The rusers facility is managed using the Service Management Facility (SMF).

To see whether the rusers facility is online, issue the command:

```bash
# svcs -a | grep rusers
online         17:00:48 svc:/network/rpc/rusers:default
```

The following is the command format for the rusers command:

```bash
rusers -options hostname 
```

The rusers -l command displays a long list of the login names of users who are logged in on local and remote systems. The output displays the name of the system into which a user is logged, the login device (TTY port), the login date and time, the idle time, and the login host name. If the user is not idle, no time is displayed in the idle time field. The term idle means that the user is not actively doing anything at the time on the terminal, which would denote the user is probably at screen lock or away from the terminal.

The following is an example of the rusers command:

```bash
# rusers -l
Sending broadcast for rusersd protocol version 3...
root         sys-02:console            Oct 17 08:21         (:0)
user5        sys-02:pts/6              Oct 17 09:20       1 (sys-03)
user3        sys-02:pts/8              Oct 17 09:21       1 (localhost)
root         fe80::203:baff:f:pts/2    Oct 17 09:18       1 (sys-02)
root         sys-03:pts/2              Oct 17 09:18       1 (sys-02)
Sending broadcast for rusersd protocol version 2...
```

## finger

To display detailed information about user activity that is either local or remote, use the finger command.

The finger command displays:

- The user's login name
- The home directory path
- The login time
- The login device name
- The data contained in the comment field of the `/etc/passwd` file (usually the user's full name)
- The login shell
- The name of the host, if the user is logged in remotely, and any idle time

The following is the command format for the finger command:

```bash
finger [-bfhilmpqsw] [username...]
finger [-l] [ username@hostname1 [ @hostname ]]
```

The -m option matches arguments only on username (not the first or last name that might appear in the comment field of `/etc/passwd`).

To display information for usera, perform the command:

```bash
# finger -m user5
Login name: user5     			
Directory: /export/home/user5       	Shell: /bin/ksh
On since Oct 17 09:20:43 on pts/6 from sys-03
1 minute 50 seconds Idle Time
No unread mail
No Plan.
```

If users create the standard ASCII files .plan or .project in their home directories, the content of those files is shown as part of the output of the finger command.

These files are traditionally used to outline a user's current plans or projects and must be created with file access permissions set to 644 (rw-r--r--).

Note: You get a response from the finger command only if the network/finger service is enabled.

```bash
# inetadm | grep finger
enabled   online         svc:/network/finger:default
```

## last

Use the last command to display a record of all logins and logouts with the most recent activity at the top of the output. The last command reads the binary file `/var/adm/wtmpx`, which records all logins, logouts, and reboots.

Each entry includes the user name, the login device, the host that the user is logged in from, the date and time that the user logged in, the time of logout, and the total login time in hours and minutes, including entries for system reboot times.

The output of the last command can be extremely long. Therefore, you might want to use it with the -n number option to specify the number of lines to display.

The following is an example of the last command:

```bash
# last
user3    pts/8      localhost      Sun Oct 17 09:21  still logged in
root     console   	:0             Sun Oct 17 08:21  still logged in
reboot   system boot               Sun Oct 17 08:00 
wtmp begins Fri Oct 15 11:36 
 
(output truncated)
```

You can use the last command also to display information about an individual user if you supply the user's login name as an argument.

```bash
# last user5
user5     pts/6       sys-03       Sun Oct 17 09:20   still logged in
user5     pts/7       localhost    Sun Oct 17 09:13 - 09:15  (00:02)
(output truncated)
```

To view the last five system reboot times only, perform the command:

```bash
# last -5 reboot
reboot    system boot                   Sun Oct 17 08:00 
reboot    system down                   Sun Oct 17 03:27 
reboot    system boot                   Sun Oct 17 03:16 
reboot    system down                   Sun Oct 17 03:27 
reboot    system boot                   Sun Oct 17 03:16
```

## logins

When a user logs in to a system either locally or remotely, the login program consults the `/etc/passwd` and the `/etc/shadow` files to authenticate the user. It verifies the user name and password entered.

If the user provides a login name that is in the `/etc/passwd` file and the correct password for that login name, the login program grants access to the system.

If the login name is not in the `/etc/passwd` file or the password is not correct for the login name, the login program denies access to the system.

You can log failed login attempts in the `/var/adm/loginlog` file. This is a useful tool if you want to determine if attempts are being made to break into a system.

By default, the loginlog file does not exist. To enable logging, you should create this file with read and write permissions for the root user only, and it should belong to the sys group.

```bash
# touch /var/adm/loginlog
# chown root:sys /var/adm/loginlog
# chmod 600 /var/adm/loginlog
```

All failed command-line login activity is written to this file automatically after five consecutive failed attempts.

The loginlog file contains one entry for each of the failed attempts. Each entry contains the user's login name, login device (TTY port), and time of the failed attempt.

If there are fewer than five consecutive failed attempts, no activity is logged to this file. This value is configured by setting the appropriate syslog_failed_login parameter in the `/etc/default/login` file:

```bash
# tail -15 /etc/default/login
 
# RETRIES determines the number of failed logins that will be
# allowed before login exits. Default is 5 and maximum is 15.
# If account locking is configured (user_attr(4)/policy.conf(4))
# for a local user's account (passwd(4)/shadow(4)), that account
# will be locked if failed logins equals or exceeds RETRIES.
#
#RETRIES=5
#
# The SYSLOG_FAILED_LOGINS variable is used to determine how many failed
# login attempts will be allowed by the system before a failed login
# message is logged, using the syslog(3) LOG_NOTICE facility.  For example,
# if the variable is set to 0, login will log -all- failed login attempts.
#
#SYSLOG_FAILED_LOGINS=5
```

## su

For security reasons, you must monitor who has been using the su command, especially those users who are trying to gain root access on the system. You can initiate the monitoring by setting two variables in the `/etc/default/su` file.

Note: There are many variables in the `/etc/default/su` file. This course presents only a small subset of the variables.

**Contents of the /etc/default/su File**

To display the contents of the `/etc/default/su` file, perform the command:

```bash
# cat /etc/default/su
#ident	"@(#)su.dfl	1.6	93/08/14 SMI"	/* SVr4.0 1.2	*/
 
# SULOG determines the location of the file used to log all su attempts
#
SULOG=/var/adm/sulog
 
# CONSOLE determines whether attempts to su to root should be logged
# to the named device
#
#CONSOLE=/dev/console
(output edited for brevity)
SYSLOG=YES
```

In the preceding example, unsuccessful attempts to use the su command to access the root account are logged to the `/var/adm/messages` file. The following is an example entry from that file:

```
Oct 16 12:35:47 sys-02 su: [ID 810491 auth.crit] 'su root' failed for user3 on /dev/pts/2
```

**The CONSOLE Variable in the /etc/default/su File**

By default, the system ignores the CONSOLE variable in the `/etc/default/su` file because of the preceding comment (#) symbol. All attempts to use the su command are logged to the console, regardless of success or failure. Here is an example of output to the console:

```bash
Feb 2 09:50:09 host1 su: 'su root' failed for user1 on /dev/pts/4
Feb 2 09:50:33 host1 su: 'su user3' succeeded for user1 on /dev/pts/4
```

When the comment symbol is removed, the value of the CONSOLE variable is defined for the `/dev/console` file. Subsequently, an additional line of output for each successful attempt to use the su command to access the root account is logged to the console. Here is an example of logged su command activity:

```bash
Feb 2 11:20:07 host1 su: 'su root' succeeded for user1 on /dev/pts/4
SU 02/02 11:20 + pts/4 user1-root
```

**The SULOG Variable in the /etc/default/su File**

The SULOG variable in the `/etc/default/su` file specifies the name of the file in which all attempts to use the su command to switch to another user are logged. If the variable is undefined, the su command logging is turned off.

The `/var/adm/sulog` file is a record of all attempts by users on the system to execute the su command. Each time the su command is executed, an entry is added to the sulog file.

The entries in this file include the date and time the command was issued, whether it was successful (shown by the plus (+) symbol for success or the hyphen (-) symbol for failure), the device from which the command was issued, and, finally, the login and the effective identity.

The following is an example of entries from the `/var/adm/sulog` file:

```bash
# more /var/adm/sulog
SU 10/17 02:51 + ??? root-uucp
SU 10/17 09:26 + pts/10 user3-root
SU 10/17 09:27 + pts/10 user3-user5
SU 10/17 09:28 + pts/10 user3-user5
SU 10/17 09:28 + pts/10 user3-root
SU 10/17 09:29 - pts/10 user3-user4
```

## Controling System Access

Note: There are many variables in the `/etc/default/login` file. This course, presents only a small subset of the variables.

The `/etc/default/login` file establishes default parameters for users when they log in to the system. The `/etc/default/login` file gives you the ability to protect the root account on a system. You can restrict root access to a specific device or to a console, or disallow root access altogether.

To display the contents of the `/etc/default/login` file, perform the command:

```bash
# cat /etc/default/login
(output edited for brevity)
# If CONSOLE is set, root can only login on that device.
# Comment this line out to allow remote login by root.
#
CONSOLE=/dev/console
 
# PASSREQ determines if login requires a password.
#
PASSREQ=YES
```

The CONSOLE Variable in the `/etc/default/login` File

You can set the CONSOLE variable in the `/etc/default/login` file to specify one of three possible conditions that restrict access to the root account:

- If the variable is defined as CONSOLE=/dev/console, the root user can log in only at the system console. Any attempt to log in as root from any other device generates the error message:

```bash
# rlogin host1
Not on system console
Connection closed.
```

- If the variable is not defined, such as #CONSOLE=/dev/console, the root user can log in to the system from any device across the network, through a modem, or using an attached terminal.

Caution: If the variable does not have a value assigned to it (for example CONSOLE= ) then the root user cannot log in from anywhere, not even the console. The only way to become the root user on the system is to log in as a regular user and then become root by using the su command.

Note: You can confine root logins to a particular port with the CONSOLE variable. For example, CONSOLE=/dev/term/a permits the root user to log in to the system only from a terminal that is connected to Serial Port A.
The PASSREQ Variable in the `/etc/default/login` File

When the PASSREQ variable in the `/etc/default/login` file is set to the default value of YES, then all users who had not been assigned passwords when their accounts were created are required to enter a new password as they log in for the first time. If this variable is set to NO, then null passwords are permitted. This variable does not apply to the root user.

For regular users, the `/etc/hosts.equiv` file identifies remote hosts and remote users who are considered to be trusted.

While the `/etc/hosts.equiv` file applies system-wide access for non-root users, the .rhosts file applies to a specific user.

All users, including the root user, can create and maintain their own .rhosts files in their home directories.

For example, if you run an rlogin process from a remote host to gain root access to a local host, the /.rhosts file is checked in the root home directory on the local host.

If the remote host name is listed in this file, it is a trusted host, and, in this case, root access is granted on the local host. The CONSOLE variable in the `/etc/default/login` file must be commented out for remote root logins.

The $HOME/.rhosts file does not exist by default. You must create it in the user's home directory.

## Get informations

The groups command displays group memberships for the user.

The command format for the groups command is:

```bash
groups [username]
```

For example, to see which groups you are a member of, perform the command:

```bash
# groups
```

other root bin sys adm uucp mail tty lp nuucp daemon

To list the groups to which a specific user is a member, use the groups command with the user's name, such as user5, as an argument.

```bash
# groups user5
staff class sysadmin
```

You use the id command to further identify users by listing their UID number, user name, GID number, and group name. This information is useful when you are troubleshooting file access problems for users.

The id command also returns the EUID number and name, and the EGID number and login name. For example, if you logged in as user1 and then used the su command to become user4, the id command reports the information for the user4 account.

The command format for the id command is:

```bash
id options username
```

To view your effective user account, perform the command:

```bash
$ id
uid=101(user1) gid=300(class)
```

To view account information for a specific user, use a user login name with the id command:

```bash
$ id user1
uid=101(user1) gid=300(class)
```

To view information about the secondary groups of a user, use the -a option and a user login name, such as user1:

```bash
$ id -a user1
uid=101(user1) gid=300(class) groups=14(sysadmin)
```

## Set uid, gid & Sticky bits

Three types of special permissions are available for executable files and directories. These are:

- The setuid permission
- The setgid permission
- The Sticky Bit permission

### The setuid Permission on Executable Files

When the set-user identification (setuid) permission is set on an executable file, a user or process that runs this executable file is granted access based on the owner of the file (usually the root user), instead of on who started the executable.

This setting allows a user to access files and directories that are typically accessible only by the owner of the executable. Note that many executable programs must be run by the root user, or by sys or bin to work properly.

Use the ls command to check the setuid permission.

```bash
# ls -l /usr/bin/su
-r-sr-xr-x   1 root     sys        22292 Jan 15 17:49 /usr/bin/su
```

The setuid permission displays as an "s" in the owner's execute field.

Note: If a capital "S" appears in the owner's execute field, it indicates that the setuid bit is on, and the execute bit "x" for the owner of the file is off or denied.

The root user and the owner can set the setuid permissions on an executable file by using the chmod command and the octal value 4###.

For example:

```bash
# chmod 4555 executable_file
```

Caution: Except for those setuid executable files that exist by default in the Solaris OS, you should disallow the use of setuid programs or at least restrict their use.

To search for files with setuid permissions and to display their full path names, perform the command:

```bash
# find / -perm -4000
```

### The setgid Permission on Executable Files

The set-group identification (setgid) permission is similar to the setuid permission, except that when the process runs, it runs as if it were a member of the same group in which the file is a member. Also, access is granted based on the permissions assigned to that group.

For example, the write program has a setgid permission that allows users to send messages to other users' terminals.

Use the ls command to check the setgid permission.

```bash
# ls -l /usr/bin/write
-r-xr-sr-x   1 root     tty        11484 Jan 15 17:55 /usr/bin/write
```

The setgid permission displays as an "s" in the group's execute field.

Note: If a lowercase letter "l" appears in the group's execute field, it indicates that the setgid bit is on, and the execute bit for the group is off or denied. This indicates that mandatory file and record locking occurs during file access for those programs that are written to request locking.

The root user and the owner can set setgid permissions on an executable file by using the chmod command and the octal value 2###. Here is the command-line format:

```bash
# chmod 2555 executable_file
```

The setgid Permission on Directories

The setgid permission is a useful feature for creating shared directories.

When a setgid permission is applied to a directory, files created in the directory belong to the group of which the directory is a member.

For example, if a user has write permission in the directory and creates a file there, that file is a member of the same group as the directory and not the user's group.

To create a shared directory, you must set the setgid bit using symbolic mode. Here is the format for that mode:

```bash
# chmod g+s shared_directory
```

To search for files with setgid permissions and display their full path names, perform the command:

```bash
# find / -perm -2000
```

### Sticky Bit Permission on Public Directories

The Sticky Bit is a special permission that protects the files within a publicly writable directory.

If the directory permissions have the Sticky Bit set, a file can be deleted only by the owner of the file, the owner of the directory, or by the root user. This prevents a user from deleting other users' files from publicly writable directories.

Use the ls command to determine if a directory has the Sticky Bit permission set.

```bash
# ls -ld /tmp
drwxrwxrwt  6  root   sys      719 May 31 03:30    /tmp
```

The Sticky Bit displays as the letter "t" in the execute field for other.

Note: If a capital "T" appears in the execute field for other, it indicates that the Sticky Bit is on; however, the execute bit is off or denied.

The root user and the owner can set the Sticky Bit permission on directories by using the chmod command and the octal value 1###. Here is the command-line format:

```bash
# chmod 1777 public_directory
```

To search for directories that have Sticky Bit permissions and display their full path names, execute the following command:

```bash
# find / -type d -perm -1000
```

Note: For more detailed information on the Sticky Bit, execute the man sticky command.
