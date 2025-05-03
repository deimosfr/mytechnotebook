---
weight: 999
url: "/Gestion_des_utilisateurs/"
title: "User Management"
description: "A comprehensive guide to user management in Solaris and Linux systems, including account creation, modification, password management, and troubleshooting login issues."
categories: ["Linux", "Security"]
date: "2006-11-30T18:31:00+02:00"
lastmod: "2006-11-30T18:31:00+02:00"
tags: ["User Management", "Solaris", "passwd", "shadow", "group", "useradd", "usermod", "smuser", "System Administration"]
toc: true
---

## Introduction

Here we're going to see how to manage users.

## passwd

Due to the critical nature of the `/etc/passwd` file, you should refrain from editing this file directly. Instead, you should use the Solaris Management Console or command-line tools to maintain the file.

The following is an example of an `/etc/passwd` file that contains the default system account entries.

```bash
root:x:0:0:Super-User:/:/sbin/sh
daemon:x:1:1::/:
bin:x:2:2::/usr/bin:
sys:x:3:3::/:
adm:x:4:4:Admin:/var/adm:
lp:x:71:8:Line Printer Admin:/usr/spool/lp:
uucp:x:5:5:uucp Admin:/usr/lib/uucp:
nuucp:x:9:9:uucp Admin:/var/spool/uucppublic:/usr/lib/uucp/uucico
smmsp:x:25:25:SendMail Message Submission Program:/:
listen:x:37:4:Network Admin:/usr/net/nls:
gdm:x:50:50:GDM Reserved UID:/:
webservd:x:80:80:WebServer Reserved UID:/:
nobody:x:60001:60001:NFS Anonymous Access User:/:
noaccess:x:60002:60002:No Access User:/:
nobody4:x:65534:65534:SunOS 4.x NFS Anonymous Access User:/:
```

Each entry in the `/etc/passwd` file contains seven fields. A colon separates each field. The following is the format for an entry:

```
loginID:x:UID:GID:comment:home_directory:login_shell
```

## shadow

Due to the critical nature of the `/etc/shadow` file, you should refrain from editing it directly. Instead, maintain the fields of the file by using the Solaris Management Console or command-line tools. Only the root user can read the `/etc/shadow` file.

The following is an example `/etc/shadow` file that contains initial system account entries.

```bash
root:rJrdhjNWQQHoY:6445::::::
daemon:NP:6445::::::
bin:NP:6445::::::
sys:NP:6445::::::
adm:NP:6445::::::
lp:NP:6445::::::
uucp:NP:6445::::::
nuucp:NP:6445::::::
smmsp:NP:6445::::::
listen:*LK*:::::::
gdm:*LK*:::::::
webservd:*LK*:::::::
nobody:*LK*:6445::::::
noaccess:*LK*:6445::::::
nobody4:*LK*:6445::::::
```

Each entry in the `/etc/shadow` file contains nine fields. A colon separates each field.

Following is the format of an entry:

```
loginID:password:lastchg:min:max:warn:inactive:expire:
```

The table defines the requirements for each of the eight fields.

{{< table "table-hover table-striped" >}}
| Field | Description |
|-------|-------------|
| loginID | The user's login name. |
| password | A 13-character encrypted password. The string *LK* indicates a locked account, and the string NP indicates no valid password. Passwords must be constructed to meet the following requirements:<br>Each password must be at least six characters and contain at least two alphabetic characters and at least one numeric or special character. It cannot be the same as the login ID or the reverse of the login ID. |
| lastchg | The number of days between January 1, 1970, and the last password modification date. |
| min | The minimum number of days required between password changes. |
| max | The maximum number of days the password is valid before the user is prompted to enter a new password at login. |
| warn | The number of days the user is warned before the password expires. |
| inactive | The number of inactive days allowed for the user before the user's account is locked. |
| expire | The date (given as number of days since January 1, 1970) when the user account expires. After the date is exceeded, the user can no longer log in. |
| flag | To track failed logins. The count is in low order four bits; the remainder is reserved for future use, set to zero. |
{{< /table >}}

## group

Each user belongs to a group that is referred to as the user's primary group. The GID number, located in the user's account entry within the `/etc/passwd` file, specifies the user's primary group.

Each user can also belong to up to 15 additional groups, known as secondary groups. In the `/etc/group` file, you can add users to group entries, thus establishing the user's secondary group affiliations.

The following is an example of the default entries in an `/etc/group` file:

```bash
root::0:
other::1:root
bin::2:root,daemon
sys::3:root,bin,adm
adm::4:root,daemon
uucp::5:root
mail::6:root
tty::7:root,adm
lp::8:root,adm
nuucp::9:root
staff::10:
daemon::12:root
sysadmin::14:
smmsp::25:
gdm::50:
webservd::80:
nobody::60001:
noaccess::60002:
nogroup::65534::
```

Each line entry in the `/etc/group` file contains four fields. A colon character separates each field. The following is the format for an entry:

```
groupname:group-password:GID:username-list
```

The table defines the requirements for each of the four fields.

{{< table "table-hover table-striped" >}}
| Field | Description |
|-------|-------------|
| groupname | Contains the name assigned to the group. Group names contain up to a maximum of eight characters. |
| group-password | Usually contains an empty field or an asterisk. This is a relic of earlier versions of UNIX.<br><br>**Caution:** A group-password is a security hole because it might allow an unauthorized user who is not a member of the group but who knows the group password, to enter the group.<br><br>**Note:** The newgrp command changes a user's primary group association within the shell environment from which it is executed. If this new, active group has a password and the user is not a listed member in that group, the user must enter the password before the newgrp command can continue. |
| GID | Contains the group's GID number. It is unique on the local system and should be unique across the organization. Numbers 0 to 99, 60001, 60002 and 65534 are reserved for system group entries. User-defined groups range from 100 to 60000. |
| username-list | Contains a comma-separated list of user names that represent the user's secondary group memberships. By default, each user can belong to a maximum of 15 secondary groups.<br><br>**Note:** The maximum number of groups is set by the kernel parameter called ngroups_max. You can set this parameter in the `/etc/system` file to allow for a maximum of 32 groups. Not all applications will be able to reference group memberships greater than 16. NFS is a notable example. |
{{< /table >}}

## The Defaults

Set values for the following parameters in the `/etc/default/passwd` file to control properties for all users' passwords on the system:

* MAXWEEKS - Sets the maximum time period (in weeks) that the password is valid.
* MINWEEKS - Sets the minimum time period before the password can be changed.
* PASSLENGTH - Sets the minimum number of characters for a password. Valid entries are 6, 7, and 8.
* WARNWEEKS - Sets the time period prior to a password's expiration to warn the user that the password will expire.

Note: The WARNWEEKS value does not exist by default in the `/etc/default/passwd` file, but it can be added.

The password aging parameters MAXWEEKS, MINWEEKS, and WARNWEEKS are default values. If set in the `/etc/shadow` file, the parameters in that file override those in the `/etc/default/passwd` file for individual users.

The Solaris 10 OS release introduces a number of new controls for password management. These controls are configured by setting values in the `/etc/default/passwd` file.

* NAMECHECK=NO - Sets the password controls to verify that the user is not using their login name as a component of the password.
* HISTORY=26 - Forces the passwd program to log up to 26 changes to the user's password. This prevents the user from reusing the same password for 26 changes. Setting the HISTORY value to zero (0) will case the password log for a user to be removed on the next password change.
* DICTIONLIST= - Causes the passwd program to perform dictionary word lookups.
* DICTIONDBDIR=/var/passwd - Causes the passwd program to perform dictionary word lookups.

Complexity of the password can be controlled using the following parameters:

```bash
#MINDIFF=3
#MINALPHA=2
#MINNONALPHA=1
#MINUPPER=0
#MINLOWER=0
#MAXREPEATS=0
#MINSPECIAL=0
#MINDIGIT=0
#WHITESPACE=YES
```

By default, all of the above parameters are commented out.

Note: By forcing greater complexity of password structure, you may inadvertently cause the users to write down their passwords as they may be too difficult for the user to remember. When setting a password change policy, you must not underestimate the problems that too much complexity may cause.

## Managing users accounts

The Solaris OS provides these command-line tools, defined as follows:

* useradd - Adds a new user account on the local system
* usermod - Modifies a user's account on the local system
* userdel - Deletes a user's account from the local system
* groupadd - Adds a new group entry to the system
* groupmod - Modifies a group entry on the system
* groupdel - Deletes a group entry from the system

In addition to these standard command-line tools, the Solaris 9 and 10 OS has a set of command-line tools that accomplish the same tasks. They are the smuser and smgroup commands.

The smuser command enables you to manage one or more users on the system with the following set of subcommands:

* add - Adds a new user account
* modify - Modifies a user's account
* delete - Deletes a user's account
* list - Lists one or more user entries

The smuser and smgroup commands interact with naming services, can use autohome functionality, and are better suited for remote management.

Note: The smuser and smgroup commands are the command-line interface equivalent to the Solaris Management Console range of operation, and allow you to perform Solaris Management Console actions in scripts. Therefore, the smuser and smgroup commands have numerous subcommands and options designed to function across domains and multiple systems. This module describes only the basic commands.

The smgroup command enables you to manage one or more groups on the system with the following set of subcommands:

* add - Adds a new group entry
* modify - Modifies a group entry
* delete - Deletes a group entry
* list - Lists one or more group entries

Any subcommand to add, modify, list, or delete users with the smuser and smgroup commands requires authentication with the Solaris Management Console server and requires the initialization of the Solaris Management Console. For example, the following is the command format for the smuser command:

```
/usr/sadm/bin/smuser subcommand [auth_args] -- [subcommand_args]
```

The authorization arguments are all optional. However, if you do not specify the authorization argument, the system might prompt you for additional information, such as a password for authentication purposes.

The -- option separates the subcommand-specific options from the authorization arguments.
The -- option must be entered even if an authorization argument is not specified because it must precede the subcommand arguments.

The subcommand arguments are quite numerous. For a complete listing of the subcommands, refer to the smuser man page. It is important to note that descriptions and other arguments that contain white space must be enclosed in double quotation marks.

Use the useradd or smuser add command to add new user accounts to the local system. These commands add an entry for a new user into the `/etc/passwd` and `/etc/shadow` files.

These commands also automatically copy all the initialization files from the `/etc/skel` directory to the user's new home directory.

### useradd

The following is the command format for the useradd command:
useradd [ -u uid ][ -g gid ][ -G gid [,gid,.. ]]

[ -d dir ][ -m ][ -s shell ][ -c comment ] loginname

The table shows the options for the useradd command.

{{< table "table-hover table-striped" >}}
| Option | Definition |
|--------|------------|
| -u uid | Sets the UID number for the new user |
| -g gid | Defines the new user's primary group |
| -G gid [,gid,..] | Defines the new user's secondary group memberships |
| -d dir | Defines the full path name for the user's home directory |
| -m | Creates the user's home directory if it does not already exist |
| -s shell | Defines the full path name for the shell program of the user's login shell |
| -c comment | Specifies any comment, such as the user's full name and location |
| loginname | Defines the user's login name for the user account |
| -D | Displays the defaults that are applied to the useradd command |
{{< /table >}}

The following example uses the useradd command to create an account for a user named newuser1. It assigns 100 as the UID number, adds the user to the group other, creates a home directory in the `/export/home` directory, and sets `/bin/ksh` as the login shell for the user account.

```
# useradd -u 100 -g other -d /export/home/newuser1 -m -s /bin/ksh -c "Regular User Account" newuser1
64 blocks
```

The useradd command has a preset range of default values. These values can be displayed using the useradd -D command. When this command has been used for the first time, the useradd command generates a file called `/var/sadm/defadduser` that contains the default values. If the contents of this file are amended, the new contents become the default values for the next time the useradd command is used.

```bash
# ls -l /usr/sadm/defadduser
/usr/sadm/defadduser: No such file or directory
# useradd -D
group=other,1  project=default,3  basedir=/home 
skel=/etc/skel  shell=/bin/sh  inactive=0  
expire=  auths=  profiles=  roles=  limitpriv=  
defaultpriv=  lock_after_retries=  
# ls -l /usr/sadm/defadduser
-rw-r--r--  1 root   root    286 Oct 17 09:04 /usr/sadm/defadduser
# cat /usr/sadm/defadduser
#	Default values for useradd.  Changed Sun Oct 17 09:04:27 2004
defgroup=1
defgname=other
defparent=/home
defskel=/etc/skel
defshell=/bin/sh
definact=0
defexpire=
defauthorization=
defrole=
defprofile=
defproj=3
defprojname=default
deflimitpriv=
defdefaultpriv=
deflock_after_retries=
```

User accounts are locked by default when added with the useradd command. This can be verified by viewing the contents of the `/etc/shadow` file:

```
# grep 'newuser1' /etc/shadow
newuser1:*LK*:12708::::::
```

By convention, a user's login name is also the user's home directory name.

You use the passwd command to create a password for the new account.

```
# passwd newuser1
New Password: 123pass
Re-enter new Password: 123pass
passwd: password successfully changed for newuser1
```

This password setting can be verified by viewing the contents of the `/etc/shadow` file:

```
# grep 'newuser1' /etc/shadow
newuser1:M0/jo1fmSbYio:12708::::::
```

### smuser add

The following is the command format for the smuser add command:

smuser add [auth_args] -- [subcommand_args]

The table shows some of the most common subcommand arguments for the smuser add command.

{{< table "table-hover table-striped" >}}
| Subcommand Argument | Definition |
|--------------------|------------|
| -c comment | A short description of the login, typically the user's name. This string can be up to 256 characters. |
| -d dir | Specifies the home directory of the new user and is limited to 1024 characters. |
| -g group | Specifies the new user's primary group membership. |
| -G group | Specifies the user's secondary group membership. |
| -n name | Specifies the user's login name. |
| -s shell | Specifies the full path name of the user's login shell. |
| -u uid | Specifies the user ID of the user you want to add. If you do not specify this option, the system assigns the next available unique UID greater than 100. |
| -x autohome=Y|N | Sets the home directory to automount if set to Y. |
{{< /table >}}

The following example uses the smuser add command to create an account for a user named newuser2. It designates the login name as newuser2, assigns the UID number 500, adds the user to the group other, creates a home directory in the `/export/home` directory, and sets `/bin/ksh` as the login shell for the user account.

Note: The -x autohome=N option to the smuser command adds the user without automounting the user's home directory. See the man page for automount for more information.

```
# /usr/sadm/bin/smuser add -- -n newuser2 -u 500 -g other -d /export/home/newuser2 -c "Regular User Account 2" -s /bin/ksh -x autohome=N
Authenticating as user: root

Type /? for help, pressing <enter> accepts the default denoted by [ ]
Please enter a string value for: password :: Enter_The_root_Password
Loading Tool: com.sun.admin.usermgr.cli.user.UserMgrCli from sys-02
Login to sys-02 as user root was successful.
Download of com.sun.admin.usermgr.cli.user.UserMgrCli from sys-02 was successful.
```

Users are added without a password by default with the smuser command. This can be verified by viewing the appropriate entry in the `/etc/shadow` file:

```
# grep 'newuser2' /etc/shadow
newuser2::12708::::::
```

Use the passwd command to create a new password for the user.

```
# passwd newuser2
New Password: 123pass
Re-enter new Password: 123pass
passwd: password successfully changed for newuser2
```

Confirm that the password change has been applied by viewing the entry for that user in the `/etc/shadow` file:

```
# grep 'newuser2' /etc/shadow
newuser2:*LK*:12708::::::
```

### smuser & usermod

* The usermod Command Format and Options

The following is the command format for the usermod command:

```
usermod [ -u uid [ -o ] ] [ -g gid ] [ -G gid [ , gid . . . ] ] 
[ -d dir ] [ -m ] [ -s shell ] [ -c comment ]
[ -l newlogname] loginname
```

In general, the options for the usermod command function the same as those for the useradd command.

The smuser modify Command Format and Options

The following is the command format for the smuser modify command:

```
smuser modify [auth_args] -- [subcommand_args]
```

In general, the options for the smuser modify command function the same as for the smuser add command. Refer to the smuser(1M) man page for additional options.

## Others commands

* Use the userdel command or smuser delete command to delete a user's login account from the system.
* To manage groups: smgroup add, groupadd, smgroup modify, groupmod, groupdel, smgroup delete

## Problems

Some of the most common problems you might encounter as a system administrator are user login problems. There are two categories of login problems: login problems when the user logs in at the command line and login problems when the user logs in from the Common Desktop Environment (CDE).

The CDE uses more configuration files, so there are more potential problems associated with logging in from the CDE. When you troubleshoot a login problem, first determine whether you can log in from the command line. Attempt to log in from another system by using either the telnet command or the rlogin command, or click Options from the CDE login panel and select Command Line Login. If you can log in successfully at the command line, then the problem is with the CDE configuration files. If you cannot log in at the command line, then the problem is more serious and involves key configuration files.

### Login Problems at the Command Line

The table presents an overview of common login problems that occur when the user logs in at the command line.

{{< table "table-hover table-striped" >}}
| Login Problem | Description |
|--------------|-------------|
| Login incorrect | This message occurs when there are problems with the login information. The most common cause of an incorrect login message is a mistyped password. Make sure the that correct password is being used, and then attempt to enter it again. Remember that passwords are case-sensitive, so you cannot interchange uppercase letters and lowercase letters. In the same way, the letter "o" is not interchangeable with the numeral "0" nor is the letter "l" interchangeable with the numeral "1." |
| Permission denied | This message occurs when there are login, password, or NIS+ security problems. Most often, an administrator has locked the user's password or the user's account has been terminated. |
| Password will not work at lockscreen | A common error is to have the Caps Lock key on, which causes all letters to be uppercase. This does not work if the password contains lowercase letters. |
| No shell | This message occurs when the user's shell does not exist, is typed incorrectly, or is wrong in the `/etc/passwd` file. |
| No directory! Logging in with home=/ | This message occurs when the user cannot access the home directory for one of the following reasons: An entry in the `/etc/passwd` file is incorrect, or the home directory has been removed or is missing, or the home directory exists on a mount point that is currently unavailable. |
| Choose a new password (followed by the New password: prompt) | This message occurs the first time a user logs in and chooses an initial password to access the account. |
| Couldn't fork a process! | This message occurs then the server could not fork a child process during login. The most common cause of this message is that the system has reached its maximum number of processes. You can either kill some unneeded processes (if you are already logged into that system as root) or increase the number of processes your system can handle. |
{{< /table >}}

### Login Problems in the CDE

Problems associated with logging into the CDE range from a user being unable to login (and returning to the CDE login screen), to the custom environment not loading properly. In general, the system does not return error messages to the user from the CDE. The following is a list of files and directories that provide troubleshooting information about the CDE:

* `/usr/dt/bin/Xsession`

This file is the configuration script for the login manager. This file should not be edited. The first user-specific file that the Xsession script calls is the $HOME/.dtprofile file. 

* `$HOME/.dtprofile`

By default, the file does not contain much content, except for examples. It contains a few echo statements for session logging purposes, and the DTSOURCEPROFILE variable is set. But it also contains information about how it might be edited. The user can edit this file to add user-specific environment variables. 

* `DTSOURCEPROFILE=true`

This line allows the user's $HOME/.login file (for csh users) or the $HOME/.profile (for other shell users) to be sourced as part of the startup process.

Sometimes a .login or .profile file contains problem commands that cause the shell to crash. If the .dtprofile file is set to source a .login or .profile file that has problem commands, desktop startup might fail.

Consequently, no desktop appears. Instead, the system redisplays the Solaris OS CDE login screen. Startup errors from the .login or .profile file are usually noted in the $HOME/.dt/startlog file. Use a Failsafe login Session or a command-line login to debug problem commands in the .login or .profile files.

* `$HOME/.dt/sessions`

This directory structure contains files and directories that configure the display of the user's custom desktop and determine the applications that start when the user logs in. Look for recent changes to files and for changes to the directory structure. For example, examine the home directory and the home.old directory or a current directory and the current.old directory. Compare the changes. The changes could provide information on a new application or on changes in the saved desktop that cause the user's login to fail.

* `$HOME/.dt`

Upon removing the entire .dt directory structure, log out, and log back in again for the system to rebuild a default .dt file structure. This action allows the user to get back into the system if the problem with the CDE files cannot be resolved.
