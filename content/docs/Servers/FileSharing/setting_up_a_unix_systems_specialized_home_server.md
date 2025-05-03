---
weight: 999
url: "/Mise_en_place_d'un_serveur_de_Home_spécialisé_systèmes_Unix/"
title: "Setting up a Unix systems specialized Home server"
description: "Learn how to set up a dedicated Home server specialized for Unix systems with centralized accounts, shared profiles, and quota management."
categories: ["Linux", "FreeBSD", "Security"]
date: "2008-07-21T12:41:00+02:00"
lastmod: "2008-07-21T12:41:00+02:00"
tags:
  [
    "nfs",
    "samba",
    "quota",
    "ldap",
    "pam",
    "servers",
    "home directory",
    "profile management",
  ]
toc: true
---

## Introduction

This article will show how to have a simultaneous dedicated profile on multiple machines.

Here are the goals of this tutorial:

- Account centralization with LDAP (this should already be done).
- Each user will have a single personal unix profile on every unix server.
- The profile is the same everywhere with a share system.
- They will have paths, aliases and environment vars easier to set via a perl script.
- The profile will be the same for every user (global environment) plus their personal customizations.
- Users on Windows will have their Unix profile on a shared drive.
- A special easy environment for Java is also possible.

## Set global environment

### Global Profile

This file will be loaded every time when a user logs on. You can add anything if you want it to be loaded. Here is the `/etc/profile` file:

```bash
# /etc/profile: system-wide .profile file for the Bourne shell (sh(1))
# and Bourne compatible shells (bash(1), ksh(1), ash(1), ...).

# If not running interactively, don't do anything
[ -z "$PS1" ] && return

# don't put duplicate lines in the history. See bash(1) for more options
export HISTCONTROL=ignoredups

# check the window size after each command and, if necessary,
# update the values of LINES and COLUMNS.
shopt -s checkwinsize

# make less more friendly for non-text input files, see lesspipe(1)
[ -x /usr/bin/lesspipe ] && eval "$(lesspipe)"

# Comment in the above and uncomment this below for a color prompt
PS1='\[\033[01;37m\][\[\033[01;32m\]`date +%D` \[\033[01;35m\]\t\[\033[01;37m\]] - \[\033[01;33m\][\w]\n\[\033[01;34m\]\u\[\033[01;31m\]@\[\033[01;34m\]\h\[\033[00m\] \[\033[01;31m\]\$\[\033[00;37m\] '

#######################################
# Default Variables and Environnement #
#######################################

# Uname system
export MYSYSTEM=`uname`

# Define Path
PATH="/usr/bin:/sbin:/usr/sbin:/usr/local/bin:/usr/local/sbin:/bin"
FULLPATH=$PATH

# Linux
if [ $MYSYSTEM = "Linux" ] ; then
        CORE_NUMBER=`grep processor /proc/cpuinfo | wc | awk '{ print $1 }'`
        RAM_INFO=`free -m | grep "Mem:" | awk '{ print $2 "Mo / Free :", $3"Mo" }'`
fi

# Solaris
if [[ $MYSYSTEM = "SunOS" && -d /usr/openwin/bin/ ]] ; then
        PATH="$PATH:/usr/openwin/bin/:/usr/X11/bin:/opt/csw/bin"
        FULLPATH=$FULLPATH:$PATH
        CORE_NUMBER=`psrinfo | wc | awk '{ print $1 }'`
        RAM_MO=`echo "`vmstat | grep -v [a-z] | awk '{ print $5 }'`/1024" | bc`
        RAM_INFO=`echo `prtconf | grep "Memory size" | awk '{ print $3 }'`"Mo / Free : "$RAM_MO Mo`
fi

# Cygwin
if [ $MYSYSTEM = "CYGWIN*" ] ; then
        export TERM=cygwin
else
        export TERM=xterm-color
fi

export PATH

usernames=( $(cut -d: -f1 /etc/passwd) )
groups=( $(cut -d: -f1 /etc/group) )

case "$TERM" in
    xterm*|rxvt|linux|cygwin)
        ;;
    *)
        nocolor=yes
        ;;
esac

# Set locales
if [ $MYSYSTEM = "SunOS" ] ; then
        export LANGUAGE=fr_FR.ISO8859-15
        export LC_ALL=fr_FR.ISO8859-15
        export LANG=fr_FR.ISO8859-15
else
        export LANGUAGE=fr_FR@euro
        export LC_ALL=fr_FR@euro
        export LANG=fr_FR@euro
fi

export LESSCHARSET=latin9
export MINICOM="-c on"
export LESS="-S -g"
GCHECK=60
WATCHFMT="%n has %a %l from %M"

# CVS
export CVS_RSH=/usr/bin/ssh
export CVSROOT=:ext:user@host:/var/lib/cvs

# History
export HISTSIZE=5000
export HISTFILE=$HOME/.bash_history
export SAVEHIST=1

# Defaut editor
export EDITOR=vim

export LISTMAX=0

# Use most if possible
if [[ -x /usr/bin/most || -x /opt/local/bin/most ]] ; then
        export PAGER=most
else
        export PAGER=more
fi

#################
# Environnement #
#################

# Usefull alias
alias utar="tar -xvzf"
alias ..='cd ..'
alias ...='cd ../..'
alias uu='source /etc/profile &> /dev/null'

# ls
if [ $MYSYSTEM = "Linux" ] ; then
        alias ls='ls --color=auto'
        alias l='ls --color=auto -lg'
        alias ll='ls --color=auto -lag | $PAGER'
fi

# Graphical User Informations
echo ""
echo "     You're connected has $USER on $HOSTNAME machine."
echo "     Connected users : `who -q | grep "=" | awk -F'=' '{ print $2 }'`"
echo "     Host OS : $MYSYSTEM - `uname -m`"
echo ""
echo "     Host CPU Core(s) : $CORE_NUMBER"
echo "     Host RAM : Total : $RAM_INFO"

# Personals mycompany users environnements
if [ -d ~/mycompany_bash ] ; then
        # Generate a correct file in /tmp dir
        for personal_env in `ls ~/mycompany_bash/` ; do
                /etc/mycompany_skel/mycompany_alias_env.pl $personal_env
        done
        # Load the environnement
        . /tmp/$USER"_genprofil" || echo "Error while trying to load your personnal environnement"
        rm -f /tmp/$USER"_genprofil"
else
        # Create the environnement for the new user
        cp -Rf /etc/skel/mycompany_bash ~/mycompany_bash
        chown -Rf $USER:Team ~/mycompany_bash
fi

if [ $JAVA_HOME ] ; then
        export PATH=$PATH_PERSO:$JAVA_HOME:$FULLPATH
else
        export PATH=$PATH_PERSO:$FULLPATH
fi

# enable programmable completion features (you don't need to enable
# this, if it's already enabled in /etc/bash.bashrc and /etc/profile
# sources /etc/bash.bashrc).
if [ -f /etc/bash_completion ]; then
    . /etc/bash_completion
fi

umask 022

### Easier env, alias and path Script

With this script, users could have a folder called mycompany_bash with their own aliases, env or path. They only have to add a file in this folder which should contain:

- alias: for an aliases file
- env: for an environment file
- path: the path file

The syntax on all those files should be like:

```

ll = ls -l

```

For those who are using JAVA_HOME environment, there is a special addition which will auto-create the full path:

```

JVM = 1.5.14

````

This will automatically create the correct path to JVM.
Now you need the following script for this syntax to be effective. Place it in `/etc/global_skel/global_alias_env.pl`:

```perl
#!/usr/bin/perl -w
# Made by Pierre Mavro
# Environnement charger for global users

use strict;
use Term::ANSIColor qw(:constants);

# Vars
my $file_to_analyse=$ARGV[0];
my $method;
my $mypath="";
my $counter=0;

# Function to check which type of file does it have to be treated
sub get_enval {
    open FILER, "<$ENV{HOME}/mycompany_bash/$file_to_analyse" or die "Couln't open file : $!\n";
    open FILEW, ">>/tmp/$ENV{USER}_genprofil" or die "Cannot create file on /tmp : $!";
    while (<FILER>) {
        if ($_ =~ /^(\w*)(\ |)\=(\ |)(.*)/) {
            # Keep contents
            my $first_arg=$1;
            my $second_arg=$4;
            # Alias
            if ($method eq "alias") {
                # Save alias in file
                print FILEW "alias $1='$4'\n";
            } elsif ($method eq "environnement") {
                # Env + AutoPATH the JVM for export
                if ($first_arg =~ /(jvm|jdk)/i) {
                    if ($second_arg =~ /(\d\.\d)\.(\d+)/) {
                        # Convert to the good format and save
                        printf FILEW "export JAVA_HOME='/test/jdk$1.0_%#02s'\n", $2;
                    } else {
                        die "Your JDK PATH is in a bad format, please contact ITs";
                    }
                } else {
                    # Convert first arg to uppercase and save
                    print FILEW "export \U$first_arg\E=$second_arg\n";
                }
            } else {
                # If not reconnized die
                die "Problem on charging $method : $_\n";
            }
        } elsif ($method eq "path") {
            chomp $_;
            unless ($_ =~ /^#/) {
                $mypath="$mypath:$_";
            }
        }
    }
    if ($method eq "path") {
        print FILEW "PATH=$ENV{PATH}$mypath\n";
        print FILEW "PATH_PERSO=$mypath\n";
    }
    close FILEW;
    close FILER;
}

sub check_problems {
    if ($file_to_analyse =~ /env/i) {
        $counter++;
    }
    if ($file_to_analyse =~ /alias/i) {
        $counter++;
    }
    if ($file_to_analyse =~ /path/i) {
        $counter++;
    }
    # If more than one path/env/alias... don't analyse this file
    if ($counter ne 1) {
        print CLEAR, RED, "\nWARNING : there is a problem with your $file_to_analyse file. Can't reconize what kind of file it is. Please rename this file correctly.\n", RESET;
    }
}

# Check for problems before starting
&check_problems;
# Alias
if ($file_to_analyse =~ /env/i) {
    $method="environnement";
    &get_enval;
# Environnement
} elsif ($file_to_analyse =~ /alias/i) {
    $method="alias";
    &get_enval;
# Path
} elsif ($file_to_analyse =~ /path/i) {
    $method="path";
    &get_enval;
}
````

Now add execution rights:

```bash
chmod 755 /etc/global_skel/global_alias_env.pl
```

### Skel

In your `/etc/skel`, you need to make some changes. First delete all the content:

```bash
rm -Rf /etc/skel
mkdir /etc/skel
```

Next, we'll add some example files to help users set their custom environment. Create a folder and some files:

```bash
mkdir /etc/skel/mycompany_bash
cd /etc/skel/mycompany_bash
touch alias env path
```

Now put this content in the alias file:

```
## Personal Mycompany Aliases ##
#
# Purpose:
# This file contains your custom aliases,
# these little shorcuts that ease your life.
#
# Instructions:
# - If you want to put comments, add #
# - If you need to create aliases follow this example:
#         myalias = command
# - Your first argument should only be written in lowercase !
#
# Your alias file can be updated with the "uu" command or relog yourself.
#
# If you need other aliases files, simply create in your ~/mycompany_bash directory
# a new file containing "alias" in its name.
#
# Examples :
# javav = java -version
# ll = ls -lah
#
```

Next env file:

```
## Personal Mycompany Environnement ##
#
# Purpose:
# This file contains your custom environment variables
# that can be used in your scripts or aliases
#
# One usefull behavior is to specify some predefined tokens like:
#               JDK = 1.5.0_14
# that will automatically modify your path to add this JDK version and create JAVA_HOME variable.
#
# Instructions:
# - If you want to put comments, add #
# - If you need to create environment variables follow this example:
#               NAME_OF_THE_VARIABLE = value
# - Your first argument (NAME_OF_THE_VARIABLE here) sould only be written in uppercase !
#
# Your environment file can be updated with the "uu" command or relog yourself.
#
# If you need other environment files, simply create in your ~/mycompany_bash directory,
# a new file containing "env" in its name.
#
# Examples :
#   MY_ENV_TEST = test
#   MY_ENV_TEST2 = $MY_ENV_TEST/test2
#   DEV_DIR = ~/dev
#   JDK = 1.5.0_14
#
```

Then the path file:

```
## Personal Mycompany Path ##
#
# Purpose:
# This file contains your custom amendment to your PATH environment variable.
# It is used to add some directories in which binaries can be used everywhere
# without referencing the full path.
#
# Instructions :
# - If you want to put comments, add #
# - If you need to add directories to your path, just simply put them one by one (only one per line).
#
# Your path file can be updated with the "uu" command or relog yourself.
#
# If you need other path files, simply create in your ~/mycompany_bash directory
# a new file containing "path" in its name.
#
# Examples :
# /usr/bin
# /usr/sbin
#
```

## Quotas

First, you should look at this documentation: [Quotas Documentation]({{< ref "docs/Linux/FilesystemsAndStorage/setting_up_quotas_on_linux.md" >}}).

For LDAP:  
If you want to run a solution on FreeBSD, you can look at [pam_quota](https://www.cyberz.org/projects/pam_quota/). Unfortunately, the developer didn't have time to create it for Linux. So we must find another way and look at session scripts: [Pam-script Documentation]({{< ref "docs/Servers/Authentication/PAM/pamscript_execute_scripts_at_authentication_session_open_and_close.md" >}}).

### LDAP

If you have an LDAP, we must take all users and create a home directory for each one of them and add quotas. That's why we will use PAM-script to do it automatically.

Once pam-script has been set up, please add these lines to the `/etc/security/onsessionopen` file:

```bash
#!/bin/sh
# Made by Pierre Mavro

soft_limit=90M
hard_limit=100M
quota_folder=/home

userid=$1
service=$2

if [ $userid != "root" ] ; then
        if [ -x /usr/sbin/quotatool ] ; then
                /usr/sbin/quotatool -u $userid -bq $soft_limit -l $hard_limit $quota_folder
        else
                echo "/usr/sbin/quotatool doesn't exist. Can't set user quota"
                exit 1
        fi
fi
```

You just need to set the 3 variables as you wish and it will automatically be applied at the next user connection.

### PAM

If you have PAM and create your own users, just edit adduser.conf and configure this line:

```
# If QUOTAUSER is set, a default quota will be set from that user with
# `edquota -p QUOTAUSER newuser'
QUOTAUSER=""
```

Add your desired quota user in this line. Now create a new folder:

```bash
mkdir -p /etc/scripts/
```

And add this script for already created users to have quotas set:

```bash
#!/bin/sh
# Set quotas for every users
# Made by Pierre Mavro

# Set soft limit (in Mo)
soft_limit="90M"
# Set hard limit (in Mo)
hard_limit="100M"
# Quota directoy (where users are located)
home_users="/home"
# Quotatool binary
quotatool_bin="/usr/sbin/quotatool"

if [ -x $quotatool_bin ] ; then
    for username in `ls $home_users | grep -v aquota.user` ; do
        quotatool -u $username -bq $soft_limit -l $hard_limit $home_users
    done
else
    echo "Sorry but couldn't locate or can't execute $quotatool_bin"
    exit 1
fi
```

Adapt this as you wish and add permissions to make it executable:

```bash
chmod 740 /etc/scripts/set_quotas.sh
```

Now add this to your crontab to run every 5 minutes for example.

## Share Servers Profiles

### NFS

For a NFS server, read the [NFS Server Documentation]({{< ref "docs/Servers/FileSharing/nfs_setting_up_an_nfs_server.md" >}}).

When installed, just configure the `/etc/export` file:

```
/home (rw)
```

After that restart NFS services.

### Samba

For Samba, please follow this: [Samba and OpenLDAP Documentation](<Installation_et_configuration_de_Samba_en_mode_"ADS"_(Authentification_sur_un_serveur_OpenLDAP).html>).

## Guests Servers

Now I need to set up a Linux server. When users log into it, the Home profile should be automatically mounted and unmounted at logout.

### Linux

We will use [pam_mount](https://pam-mount.sourceforge.net/) to automate the NFS mount. Please [follow this documentation](./pam_mount_:_monter_des_partages_réseaux_au_login.html).

Edit the `/etc/security/pam_mount.conf` file and configure what you need. Here we would like users to have NFS home share mounted from the server at logon. Add this line at the end of the file:

```
volume * nfs my_nfs_server /home/& ~ - - -
```

Edit the `/etc/pam.d/ssh` file and add these lines:

```
...
auth       required     pam_env.so envfile=/etc/default/locale
auth       required     pam_mount.so
@include   common-auth
...
@include   common-account
session    required     pam_mount.so
@include   common-session
...
```
